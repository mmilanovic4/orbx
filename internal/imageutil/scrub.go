package imageutil

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
)

type Removed struct {
	Name  string
	Bytes int
}

func Scrub(data []byte, stripICC bool) ([]byte, []Removed, error) {
	switch detect(data) {
	case formatJPEG:
		return scrubJPEG(data, stripICC)
	case formatPNG:
		return scrubPNG(data, stripICC)
	case formatWebP:
		return scrubWebP(data, stripICC)
	case formatTIFF:
		return nil, nil, fmt.Errorf("TIFF is not supported: its metadata is part of the image structure")
	case formatHEIF:
		return nil, nil, fmt.Errorf("HEIC and AVIF are not supported, convert the image to JPEG first")
	}
	return nil, nil, fmt.Errorf("unsupported file format (expected JPEG, PNG or WebP)")
}

func scrubJPEG(data []byte, stripICC bool) ([]byte, []Removed, error) {
	out := make([]byte, 0, len(data))
	out = append(out, data[:2]...)

	var removed []Removed
	keptOrientation := false
	i := 2
	for i+4 <= len(data) {
		if data[i] != 0xFF {
			return nil, nil, fmt.Errorf("malformed JPEG at offset %d", i)
		}

		marker := data[i+1]
		// any marker may be preceded by 0xFF fill bytes, they carry nothing
		if marker == 0xFF {
			i++
			continue
		}
		if marker == 0x01 || (marker >= 0xD0 && marker <= 0xD8) {
			out = append(out, data[i:i+2]...)
			i += 2
			continue
		}
		// the scan data after SOS is copied untouched, which is what keeps this lossless
		if marker == 0xDA || marker == 0xD9 {
			end := i + 2
			if marker == 0xDA {
				end = jpegImageEnd(data, i)
			}
			if end < len(data) {
				removed = append(removed, Removed{Name: "Trailing data", Bytes: len(data) - end})
			}
			return append(out, data[i:end]...), removed, nil
		}

		size := int(binary.BigEndian.Uint16(data[i+2 : i+4]))
		if size < 2 || i+2+size > len(data) {
			return nil, nil, fmt.Errorf("malformed JPEG segment at offset %d", i)
		}

		segment := data[i : i+2+size]
		payload := segment[4:]
		if marker == 0xE1 && bytes.HasPrefix(payload, []byte("Exif\x00\x00")) && !keptOrientation {
			if o := exifOrientation(payload[6:]); o > 1 {
				ycbcr := exifIFD0Short(payload[6:], tagYCbCrPositioning)
				if ycbcr == 0 {
					ycbcr = 1 // centered, the default
				}
				exif := append([]byte("Exif\x00\x00"), orientationEXIF(o, ycbcr)...)
				kept := []byte{0xFF, 0xE1, 0, 0}
				binary.BigEndian.PutUint16(kept[2:], uint16(len(exif)+2))
				kept = append(kept, exif...)

				out = append(out, kept...)
				removed = append(removed, Removed{Name: exifOrientationKept, Bytes: len(segment) - len(kept)})
				keptOrientation = true
				i += 2 + size
				continue
			}
		}
		if name, drop := jpegSegment(marker, payload, stripICC); drop {
			removed = append(removed, Removed{Name: name, Bytes: len(segment)})
		} else {
			out = append(out, segment...)
		}

		i += 2 + size
	}

	return out, removed, nil
}

// jpegImageEnd returns the offset just past the EOI marker that ends the
// image whose first scan starts at sos. What follows was appended to the
// file (MPF images such as HDR gain maps, Motion Photo videos, vendor
// trailers) and can carry metadata of its own.
func jpegImageEnd(data []byte, sos int) int {
	i := sos
	for i+1 < len(data) {
		if data[i] != 0xFF {
			next := bytes.IndexByte(data[i:], 0xFF)
			if next < 0 {
				break
			}
			i += next
			continue
		}

		switch marker := data[i+1]; {
		case marker == 0xD9: // EOI
			return i + 2
		case marker == 0xD8: // SOI of an appended image, the EOI is missing
			return i
		case marker == 0xFF: // fill byte before a marker
			i++
		case marker == 0x00 || marker == 0x01 || (marker >= 0xD0 && marker <= 0xD7):
			// stuffed 0xFF byte, TEM or restart marker within the scan data
			i += 2
		default:
			// a segment between scans: DHT, DQT, DRI, the next SOS, ...
			if i+4 > len(data) {
				return len(data)
			}
			i += 2 + int(binary.BigEndian.Uint16(data[i+2:i+4]))
		}
	}
	return len(data)
}

func jpegSegment(marker byte, payload []byte, stripICC bool) (string, bool) {
	switch {
	case marker == 0xFE:
		return "Comment", true
	case marker == 0xE0:
		// JFIF holds pixel density, JFXX only a thumbnail
		return "JFIF thumbnail", bytes.HasPrefix(payload, []byte("JFXX\x00"))
	case marker == 0xE1:
		if bytes.HasPrefix(payload, []byte("Exif\x00\x00")) {
			return "EXIF", true
		}
		if bytes.HasPrefix(payload, []byte("http://ns.adobe.com/xap/")) {
			return "XMP", true
		}
		return "APP1", true
	case marker == 0xE2:
		if bytes.HasPrefix(payload, []byte("ICC_PROFILE\x00")) {
			return "ICC profile", stripICC
		}
		return "APP2", true
	case marker == 0xED:
		return "IPTC", true
	case marker == 0xEE:
		// the Adobe marker defines the color transform, dropping it breaks CMYK files
		return "", false
	case marker >= 0xE0 && marker <= 0xEF:
		return fmt.Sprintf("APP%d", marker-0xE0), true
	}
	return "", false
}

func scrubPNG(data []byte, stripICC bool) ([]byte, []Removed, error) {
	out := make([]byte, 0, len(data))
	out = append(out, data[:8]...)

	var removed []Removed
	keptOrientation := false
	i := 8
	for i+8 <= len(data) {
		size := int(binary.BigEndian.Uint32(data[i : i+4]))
		name := string(data[i+4 : i+8])
		end := i + 12 + size
		if size < 0 || end > len(data) {
			return nil, nil, fmt.Errorf("malformed PNG chunk at offset %d", i)
		}

		if name == "eXIf" && !keptOrientation {
			if o := exifOrientation(data[i+8 : i+8+size]); o > 1 {
				kept := pngChunkBytes("eXIf", orientationEXIF(o, 0))
				out = append(out, kept...)
				removed = append(removed, Removed{Name: exifOrientationKept, Bytes: end - i - len(kept)})
				keptOrientation = true
				i = end
				continue
			}
		}

		if label, drop := pngChunk(name, data[i+8:i+8+size], stripICC); drop {
			removed = append(removed, Removed{Name: label, Bytes: end - i})
		} else {
			out = append(out, data[i:end]...)
		}

		i = end

		// nothing after IEND is part of the image
		if name == "IEND" {
			break
		}
	}

	if i < len(data) {
		removed = append(removed, Removed{Name: "Trailing data", Bytes: len(data) - i})
	}

	return out, removed, nil
}

func pngChunkBytes(name string, payload []byte) []byte {
	out := binary.BigEndian.AppendUint32(nil, uint32(len(payload)))
	out = append(out, name...)
	out = append(out, payload...)
	return binary.BigEndian.AppendUint32(out, crc32.ChecksumIEEE(out[4:]))
}

func pngChunk(name string, payload []byte, stripICC bool) (string, bool) {
	switch name {
	case "eXIf":
		return "EXIF", true
	case "iCCP":
		return "ICC profile", stripICC
	case "tIME":
		return "Timestamp", true
	case "tEXt", "zTXt", "iTXt":
		if bytes.HasPrefix(payload, []byte("XML:com.adobe.xmp")) {
			return "XMP", true
		}
		return "Text metadata", true
	}
	return "", false
}

const (
	webpFlagICC  = 0x20
	webpFlagEXIF = 0x08
	webpFlagXMP  = 0x04
)

func scrubWebP(data []byte, stripICC bool) ([]byte, []Removed, error) {
	body := make([]byte, 0, len(data))
	var removed []Removed
	var flags byte
	keptOrientation := false

	// the RIFF header says where the image ends, anything after that was
	// appended to the file
	riffEnd := 8 + int(binary.LittleEndian.Uint32(data[4:8]))
	if riffEnd < 12 {
		return nil, nil, fmt.Errorf("malformed WebP header")
	}
	limit := min(riffEnd, len(data))

	i := 12
	for i+8 <= limit {
		name := string(data[i : i+4])
		size := int(binary.LittleEndian.Uint32(data[i+4 : i+8]))
		end := i + 8 + size + size%2
		if size < 0 || i+8+size > limit {
			return nil, nil, fmt.Errorf("malformed WebP chunk at offset %d", i)
		}
		if end > limit {
			end = limit
		}

		switch {
		case name == "EXIF":
			payload := bytes.TrimPrefix(data[i+8:i+8+size], []byte("Exif\x00\x00"))
			if o := exifOrientation(payload); o > 1 && !keptOrientation {
				exif := orientationEXIF(o, 0) // even length, no padding needed
				kept := append([]byte("EXIF"), binary.LittleEndian.AppendUint32(nil, uint32(len(exif)))...)
				kept = append(kept, exif...)
				body = append(body, kept...)
				removed = append(removed, Removed{Name: exifOrientationKept, Bytes: end - i - len(kept)})
				keptOrientation = true
				break
			}
			removed = append(removed, Removed{Name: "EXIF", Bytes: end - i})
			flags |= webpFlagEXIF
		case name == "XMP ":
			removed = append(removed, Removed{Name: "XMP", Bytes: end - i})
			flags |= webpFlagXMP
		case name == "ICCP" && stripICC:
			removed = append(removed, Removed{Name: "ICC profile", Bytes: end - i})
			flags |= webpFlagICC
		default:
			body = append(body, data[i:end]...)
		}

		i = end
	}

	if limit < len(data) {
		removed = append(removed, Removed{Name: "Trailing data", Bytes: len(data) - limit})
	}

	// VP8X advertises which optional chunks follow, so its flags have to match
	if keptOrientation {
		flags &^= webpFlagEXIF
	}
	if flags != 0 && len(body) >= 18 && string(body[:4]) == "VP8X" {
		body[8] &^= flags
	}

	out := make([]byte, 0, len(body)+12)
	out = append(out, "RIFF"...)
	out = binary.LittleEndian.AppendUint32(out, uint32(len(body)+4))
	out = append(out, "WEBP"...)

	return append(out, body...), removed, nil
}

// exifOrientationKept labels an EXIF block that was replaced by one holding
// only the Orientation tag.
const exifOrientationKept = "EXIF except orientation"

// exifOrientation returns the Orientation tag (1-8) of a TIFF-structured
// EXIF block, or 0 if it has none.
func exifOrientation(tiff []byte) int {
	if o := exifIFD0Short(tiff, tagOrientation); o >= 1 && o <= 8 {
		return o
	}
	return 0
}

// exifIFD0Short returns the value of a single-value integer tag in IFD0 of
// a TIFF-structured EXIF block, or 0 if it is missing.
func exifIFD0Short(tiff []byte, tag uint16) int {
	if len(tiff) < 8 {
		return 0
	}

	var bo binary.ByteOrder
	switch string(tiff[:2]) {
	case "II":
		bo = binary.LittleEndian
	case "MM":
		bo = binary.BigEndian
	default:
		return 0
	}
	if bo.Uint16(tiff[2:4]) != 42 {
		return 0
	}

	r := &exifReader{tiff: tiff, bo: bo}
	ifd0, err := r.parseIFD(bo.Uint32(tiff[4:8]))
	if err != nil {
		return 0
	}
	for _, e := range ifd0 {
		if e.tag == tag {
			if v := r.ints(e); len(v) == 1 && v[0] > 0 && v[0] <= 0xFFFF {
				return int(v[0])
			}
		}
	}
	return 0
}

const (
	tagOrientation      = 0x0112
	tagYCbCrPositioning = 0x0213
)

// orientationEXIF builds the smallest EXIF block that still carries the
// orientation: a TIFF header and an IFD0 with the Orientation tag. Without
// it a photo taken in portrait shows up rotated, and the tag says nothing
// about who took the photo or where. The EXIF spec also requires
// YCbCrPositioning in JPEG files, it is added when ycbcr is not 0.
func orientationEXIF(orientation, ycbcr int) []byte {
	type field struct{ tag, value uint16 }
	fields := []field{{tagOrientation, uint16(orientation)}}
	if ycbcr != 0 {
		fields = append(fields, field{tagYCbCrPositioning, uint16(ycbcr)})
	}

	out := []byte("MM\x00\x2a\x00\x00\x00\x08") // big-endian, IFD0 at offset 8
	out = binary.BigEndian.AppendUint16(out, uint16(len(fields)))
	for _, f := range fields { // sorted by tag, as TIFF requires
		out = binary.BigEndian.AppendUint16(out, f.tag)
		out = binary.BigEndian.AppendUint16(out, typeShort)
		out = binary.BigEndian.AppendUint32(out, 1)       // count
		out = binary.BigEndian.AppendUint16(out, f.value) // left-aligned in the 4-byte field
		out = append(out, 0, 0)
	}
	return append(out, 0, 0, 0, 0) // no next IFD
}
