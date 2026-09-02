package imageutil

import (
	"bytes"
	"encoding/binary"
	"fmt"
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
	}
	return nil, nil, fmt.Errorf("unsupported file format (expected JPEG, PNG or WebP)")
}

func scrubJPEG(data []byte, stripICC bool) ([]byte, []Removed, error) {
	out := make([]byte, 0, len(data))
	out = append(out, data[:2]...)

	var removed []Removed
	i := 2
	for i+4 <= len(data) {
		if data[i] != 0xFF {
			return nil, nil, fmt.Errorf("malformed JPEG at offset %d", i)
		}

		marker := data[i+1]
		if marker == 0x01 || (marker >= 0xD0 && marker <= 0xD8) {
			out = append(out, data[i:i+2]...)
			i += 2
			continue
		}
		// the scan data after SOS is copied untouched, which is what keeps this lossless
		if marker == 0xDA || marker == 0xD9 {
			return append(out, data[i:]...), removed, nil
		}

		size := int(binary.BigEndian.Uint16(data[i+2 : i+4]))
		if size < 2 || i+2+size > len(data) {
			return nil, nil, fmt.Errorf("malformed JPEG segment at offset %d", i)
		}

		segment := data[i : i+2+size]
		if name, drop := jpegSegment(marker, segment[4:], stripICC); drop {
			removed = append(removed, Removed{Name: name, Bytes: len(segment)})
		} else {
			out = append(out, segment...)
		}

		i += 2 + size
	}

	return out, removed, nil
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
	i := 8
	for i+8 <= len(data) {
		size := int(binary.BigEndian.Uint32(data[i : i+4]))
		name := string(data[i+4 : i+8])
		end := i + 12 + size
		if size < 0 || end > len(data) {
			return nil, nil, fmt.Errorf("malformed PNG chunk at offset %d", i)
		}

		if label, drop := pngChunk(name, data[i+8:i+8+size], stripICC); drop {
			removed = append(removed, Removed{Name: label, Bytes: end - i})
		} else {
			out = append(out, data[i:end]...)
		}

		i = end
	}

	return out, removed, nil
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

	i := 12
	for i+8 <= len(data) {
		name := string(data[i : i+4])
		size := int(binary.LittleEndian.Uint32(data[i+4 : i+8]))
		end := i + 8 + size + size%2
		if size < 0 || i+8+size > len(data) {
			return nil, nil, fmt.Errorf("malformed WebP chunk at offset %d", i)
		}
		if end > len(data) {
			end = len(data)
		}

		switch {
		case name == "EXIF":
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

	// VP8X advertises which optional chunks follow, so its flags have to match
	if flags != 0 && len(body) >= 18 && string(body[:4]) == "VP8X" {
		body[8] &^= flags
	}

	out := make([]byte, 0, len(body)+12)
	out = append(out, "RIFF"...)
	out = binary.LittleEndian.AppendUint32(out, uint32(len(body)+4))
	out = append(out, "WEBP"...)

	return append(out, body...), removed, nil
}
