package imageutil

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
)

type Tag struct {
	Name  string
	Value string
}

const (
	typeByte      = 1
	typeASCII     = 2
	typeShort     = 3
	typeLong      = 4
	typeRational  = 5
	typeSByte     = 6
	typeUndefined = 7
	typeSShort    = 8
	typeSLong     = 9
	typeSRational = 10
	typeFloat     = 11
	typeDouble    = 12
)

var typeSizes = map[uint16]uint32{
	typeByte: 1, typeASCII: 1, typeShort: 2, typeLong: 4, typeRational: 8,
	typeSByte: 1, typeUndefined: 1, typeSShort: 2, typeSLong: 4, typeSRational: 8,
	typeFloat: 4, typeDouble: 8,
}

var pointerTags = map[uint16]bool{
	0x8769: true, // Exif IFD
	0x8825: true, // GPS IFD
	0xA005: true, // Interoperability IFD
}

type entry struct {
	tag   uint16
	typ   uint16
	count uint32
	raw   []byte
}

type exifReader struct {
	tiff []byte
	bo   binary.ByteOrder
}

func ReadEXIF(data []byte, all bool) ([]Tag, error) {
	tiff, err := findEXIF(data)
	if err != nil {
		return nil, err
	}
	if len(tiff) < 8 {
		return nil, fmt.Errorf("EXIF block is truncated")
	}

	var bo binary.ByteOrder
	switch string(tiff[0:2]) {
	case "II":
		bo = binary.LittleEndian
	case "MM":
		bo = binary.BigEndian
	default:
		return nil, fmt.Errorf("invalid TIFF byte order marker")
	}

	r := &exifReader{tiff: tiff, bo: bo}
	if r.bo.Uint16(tiff[2:4]) != 42 {
		return nil, fmt.Errorf("invalid TIFF header")
	}

	ifd0, err := r.parseIFD(r.bo.Uint32(tiff[4:8]))
	if err != nil {
		return nil, err
	}

	tags := r.decodeIFD(ifd0, ifd0Tags, all)

	if off, ok := r.pointer(ifd0, 0x8769); ok {
		if sub, err := r.parseIFD(off); err == nil {
			tags = append(tags, r.decodeIFD(sub, exifTags, all)...)
		}
	}

	if off, ok := r.pointer(ifd0, 0x8825); ok {
		if gps, err := r.parseIFD(off); err == nil {
			tags = append(tags, r.decodeGPS(gps, all)...)
		}
	}

	if len(tags) == 0 {
		return nil, fmt.Errorf("no EXIF tags found")
	}

	return tags, nil
}

func findEXIF(data []byte) ([]byte, error) {
	switch {
	case len(data) >= 4 && data[0] == 0xFF && data[1] == 0xD8:
		return findEXIFInJPEG(data)
	case len(data) >= 8 && bytes.HasPrefix(data, []byte("\x89PNG\r\n\x1a\n")):
		return findEXIFInPNG(data)
	case len(data) >= 12 && bytes.HasPrefix(data, []byte("RIFF")) && string(data[8:12]) == "WEBP":
		return findEXIFInWebP(data)
	case len(data) >= 4 && (bytes.HasPrefix(data, []byte("II\x2a\x00")) || bytes.HasPrefix(data, []byte("MM\x00\x2a"))):
		return data, nil // bare TIFF (also covers most raw formats)
	}
	return nil, fmt.Errorf("unsupported file format (expected JPEG, TIFF, PNG or WebP)")
}

func findEXIFInJPEG(data []byte) ([]byte, error) {
	i := 2
	for i+4 <= len(data) {
		if data[i] != 0xFF {
			return nil, fmt.Errorf("malformed JPEG at offset %d", i)
		}

		marker := data[i+1]
		// standalone markers carry no payload
		if marker == 0x01 || (marker >= 0xD0 && marker <= 0xD8) {
			i += 2
			continue
		}
		// start of scan or end of image: no metadata beyond this point
		if marker == 0xDA || marker == 0xD9 {
			break
		}

		size := int(binary.BigEndian.Uint16(data[i+2 : i+4]))
		if size < 2 || i+2+size > len(data) {
			return nil, fmt.Errorf("malformed JPEG segment at offset %d", i)
		}

		payload := data[i+4 : i+2+size]
		if marker == 0xE1 && bytes.HasPrefix(payload, []byte("Exif\x00\x00")) {
			return payload[6:], nil
		}

		i += 2 + size
	}
	return nil, fmt.Errorf("no EXIF data found")
}

func findEXIFInPNG(data []byte) ([]byte, error) {
	i := 8
	for i+8 <= len(data) {
		size := int(binary.BigEndian.Uint32(data[i : i+4]))
		name := string(data[i+4 : i+8])
		start := i + 8
		if size < 0 || start+size > len(data) {
			return nil, fmt.Errorf("malformed PNG chunk at offset %d", i)
		}
		if name == "eXIf" {
			return data[start : start+size], nil
		}
		if name == "IEND" {
			break
		}
		i = start + size + 4 // skip payload and CRC
	}
	return nil, fmt.Errorf("no EXIF data found")
}

func findEXIFInWebP(data []byte) ([]byte, error) {
	i := 12
	for i+8 <= len(data) {
		name := string(data[i : i+4])
		size := int(binary.LittleEndian.Uint32(data[i+4 : i+8]))
		start := i + 8
		if size < 0 || start+size > len(data) {
			return nil, fmt.Errorf("malformed WebP chunk at offset %d", i)
		}
		if name == "EXIF" {
			payload := data[start : start+size]
			// some encoders prepend the JPEG-style "Exif\0\0" header
			payload = bytes.TrimPrefix(payload, []byte("Exif\x00\x00"))
			return payload, nil
		}
		i = start + size + size%2 // chunks are padded to an even size
	}
	return nil, fmt.Errorf("no EXIF data found")
}

func (r *exifReader) parseIFD(offset uint32) ([]entry, error) {
	if int(offset)+2 > len(r.tiff) {
		return nil, fmt.Errorf("IFD offset out of range")
	}

	count := int(r.bo.Uint16(r.tiff[offset : offset+2]))
	pos := int(offset) + 2
	if pos+count*12 > len(r.tiff) {
		return nil, fmt.Errorf("IFD entries out of range")
	}

	entries := make([]entry, 0, count)
	for range count {
		e := entry{
			tag:   r.bo.Uint16(r.tiff[pos : pos+2]),
			typ:   r.bo.Uint16(r.tiff[pos+2 : pos+4]),
			count: r.bo.Uint32(r.tiff[pos+4 : pos+8]),
		}
		pos += 12

		size, ok := typeSizes[e.typ]
		if !ok {
			continue // unknown type, nothing sane to decode
		}

		total := size * e.count
		field := r.tiff[pos-4 : pos]
		if total <= 4 {
			e.raw = field[:total]
		} else {
			start := r.bo.Uint32(field)
			if int(start)+int(total) > len(r.tiff) {
				continue // value points outside the block
			}
			e.raw = r.tiff[start : start+total]
		}

		entries = append(entries, e)
	}

	return entries, nil
}

func (r *exifReader) pointer(entries []entry, tag uint16) (uint32, bool) {
	for _, e := range entries {
		if e.tag == tag {
			if v := r.ints(e); len(v) > 0 {
				return uint32(v[0]), true
			}
		}
	}
	return 0, false
}

func (r *exifReader) decodeIFD(entries []entry, names map[uint16]string, all bool) []Tag {
	var tags []Tag
	for _, e := range entries {
		if pointerTags[e.tag] {
			continue
		}

		name, known := names[e.tag]
		if !known {
			if !all {
				continue
			}
			name = fmt.Sprintf("Tag-0x%04X", e.tag)
		}

		tags = append(tags, Tag{Name: name, Value: r.format(e, known)})
	}
	return tags
}

func (r *exifReader) format(e entry, known bool) string {
	if known {
		if f, ok := formatters[e.tag]; ok {
			if v := f(r, e); v != "" {
				return v
			}
		}
		if enum, ok := enums[e.tag]; ok {
			if v := r.ints(e); len(v) == 1 {
				if name, ok := enum[v[0]]; ok {
					return name
				}
				return fmt.Sprintf("Unknown (%d)", v[0])
			}
		}
	}
	return r.formatRaw(e)
}

func (r *exifReader) formatRaw(e entry) string {
	switch e.typ {
	case typeASCII:
		return r.str(e)
	case typeUndefined:
		if isPrintable(e.raw) {
			return r.str(e)
		}
		return fmt.Sprintf("<%d bytes>", len(e.raw))
	case typeRational, typeSRational:
		var parts []string
		for _, v := range r.rationals(e) {
			parts = append(parts, formatRational(v[0], v[1]))
		}
		return strings.Join(parts, ", ")
	case typeFloat, typeDouble:
		var parts []string
		for _, v := range r.floats(e) {
			parts = append(parts, formatFloat(v))
		}
		return strings.Join(parts, ", ")
	default:
		var parts []string
		for _, v := range r.ints(e) {
			parts = append(parts, strconv.FormatInt(v, 10))
		}
		return strings.Join(parts, ", ")
	}
}
