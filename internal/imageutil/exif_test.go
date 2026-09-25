package imageutil

import (
	"encoding/binary"
	"testing"
)

type testEntry struct {
	tag   uint16
	typ   uint16
	count uint32
	data  []byte
}

func ascii(s string) testEntry {
	b := append([]byte(s), 0)
	return testEntry{typ: typeASCII, count: uint32(len(b)), data: b}
}

func short(bo binary.ByteOrder, v uint16) testEntry {
	b := make([]byte, 2)
	bo.PutUint16(b, v)
	return testEntry{typ: typeShort, count: 1, data: b}
}

func rational(bo binary.ByteOrder, values ...[2]uint32) testEntry {
	b := make([]byte, 0, 8*len(values))
	for _, v := range values {
		chunk := make([]byte, 8)
		bo.PutUint32(chunk, v[0])
		bo.PutUint32(chunk[4:], v[1])
		b = append(b, chunk...)
	}
	return testEntry{typ: typeRational, count: uint32(len(values)), data: b}
}

func withTag(tag uint16, e testEntry) testEntry {
	e.tag = tag
	return e
}

func encodeIFD(bo binary.ByteOrder, entries []testEntry, base uint32) []byte {
	fixed := 2 + 12*len(entries) + 4
	buf := make([]byte, fixed)
	bo.PutUint16(buf, uint16(len(entries)))

	var ext []byte
	for i, e := range entries {
		p := 2 + 12*i
		bo.PutUint16(buf[p:], e.tag)
		bo.PutUint16(buf[p+2:], e.typ)
		bo.PutUint32(buf[p+4:], e.count)
		if len(e.data) <= 4 {
			copy(buf[p+8:], e.data)
			continue
		}
		bo.PutUint32(buf[p+8:], base+uint32(fixed+len(ext)))
		ext = append(ext, e.data...)
	}

	return append(buf, ext...)
}

func ifdSize(entries []testEntry) int {
	size := 2 + 12*len(entries) + 4
	for _, e := range entries {
		if len(e.data) > 4 {
			size += len(e.data)
		}
	}
	return size
}

func buildTIFF(bo binary.ByteOrder, ifd0, exif, gps []testEntry) []byte {
	pointers := 0
	if len(exif) > 0 {
		pointers++
	}
	if len(gps) > 0 {
		pointers++
	}

	base := uint32(8)
	next := base + uint32(ifdSize(ifd0)+12*pointers)

	if len(exif) > 0 {
		offset := make([]byte, 4)
		bo.PutUint32(offset, next)
		ifd0 = append(ifd0, testEntry{tag: 0x8769, typ: typeLong, count: 1, data: offset})
		next += uint32(ifdSize(exif))
	}
	if len(gps) > 0 {
		offset := make([]byte, 4)
		bo.PutUint32(offset, next)
		ifd0 = append(ifd0, testEntry{tag: 0x8825, typ: typeLong, count: 1, data: offset})
	}

	header := make([]byte, 8)
	if bo == binary.LittleEndian {
		copy(header, "II")
	} else {
		copy(header, "MM")
	}
	bo.PutUint16(header[2:], 42)
	bo.PutUint32(header[4:], base)

	out := append(header, encodeIFD(bo, ifd0, base)...)
	if len(exif) > 0 {
		out = append(out, encodeIFD(bo, exif, uint32(len(out)))...)
	}
	if len(gps) > 0 {
		out = append(out, encodeIFD(bo, gps, uint32(len(out)))...)
	}

	return out
}

func wrapJPEG(tiff []byte) []byte {
	payload := append([]byte("Exif\x00\x00"), tiff...)
	out := []byte{0xFF, 0xD8, 0xFF, 0xE1}
	size := make([]byte, 2)
	binary.BigEndian.PutUint16(size, uint16(len(payload)+2))
	out = append(out, size...)
	out = append(out, payload...)
	return append(out, 0xFF, 0xDA)
}

func wrapPNG(tiff []byte) []byte {
	out := []byte("\x89PNG\r\n\x1a\n")
	size := make([]byte, 4)
	binary.BigEndian.PutUint32(size, uint32(len(tiff)))
	out = append(out, size...)
	out = append(out, []byte("eXIf")...)
	out = append(out, tiff...)
	return append(out, 0, 0, 0, 0)
}

func sampleTIFF(bo binary.ByteOrder) []byte {
	ifd0 := []testEntry{
		withTag(0x010F, ascii("Canon")),
		withTag(0x0110, ascii("Canon EOS R6")),
		withTag(0x0112, short(bo, 6)),
		withTag(0x9999, short(bo, 1)),
	}
	exif := []testEntry{
		withTag(0x829A, rational(bo, [2]uint32{1, 250})),
		withTag(0x829D, rational(bo, [2]uint32{28, 10})),
		withTag(0x8827, short(bo, 400)),
		withTag(0x9209, short(bo, 25)),
		withTag(0x920A, rational(bo, [2]uint32{50, 1})),
	}
	gps := []testEntry{
		withTag(0x0001, ascii("N")),
		withTag(0x0002, rational(bo, [2]uint32{44, 1}, [2]uint32{48, 1}, [2]uint32{45, 1})),
		withTag(0x0003, ascii("W")),
		withTag(0x0004, rational(bo, [2]uint32{20, 1}, [2]uint32{27, 1}, [2]uint32{0, 1})),
	}
	return buildTIFF(bo, ifd0, exif, gps)
}

func lookup(tags []Tag, name string) (string, bool) {
	for _, t := range tags {
		if t.Name == name {
			return t.Value, true
		}
	}
	return "", false
}

func TestReadEXIF(t *testing.T) {
	orders := map[string]binary.ByteOrder{
		"little endian": binary.LittleEndian,
		"big endian":    binary.BigEndian,
	}

	expected := map[string]string{
		"Make":           "Canon",
		"Model":          "Canon EOS R6",
		"Orientation":    "Rotate 90 CW",
		"ExposureTime":   "1/250 s",
		"FNumber":        "f/2.8",
		"ISO":            "400",
		"Flash":          "fired, auto (0x19)",
		"FocalLength":    "50 mm",
		"GPSLatitude":    `44° 48' 45.00" N`,
		"GPSLongitude":   `20° 27' 0.00" W`,
		"GPSCoordinates": "44.812500, -20.450000",
	}

	for name, bo := range orders {
		t.Run(name, func(t *testing.T) {
			tags, err := ReadEXIF(wrapJPEG(sampleTIFF(bo)), false)
			if err != nil {
				t.Fatalf("ReadEXIF() error = %v", err)
			}

			for tag, want := range expected {
				got, ok := lookup(tags, tag)
				if !ok {
					t.Errorf("missing tag %q", tag)
					continue
				}
				if got != want {
					t.Errorf("%s = %q, want %q", tag, got, want)
				}
			}
		})
	}
}

func TestFormatSeconds(t *testing.T) {
	tests := []struct {
		input    float64
		expected string
	}{
		{0.004, "1/250 s"},
		{0.25, "1/4 s"},
		{0.6, "0.6 s"},
		{1.3, "1.3 s"},
		{30, "30 s"},
		{0, ""},
	}

	for _, tt := range tests {
		if got := formatSeconds(tt.input); got != tt.expected {
			t.Errorf("formatSeconds(%v) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestReadEXIFFractionalGPS(t *testing.T) {
	bo := binary.BigEndian
	ifd0 := []testEntry{withTag(0x010F, ascii("Canon"))}
	gps := []testEntry{
		withTag(0x0001, ascii("N")),
		withTag(0x0002, rational(bo, [2]uint32{44, 1}, [2]uint32{4875, 100}, [2]uint32{0, 1})),
		withTag(0x0003, ascii("E")),
		withTag(0x0004, rational(bo, [2]uint32{204612, 10000}, [2]uint32{0, 1}, [2]uint32{0, 1})),
	}

	tags, err := ReadEXIF(wrapJPEG(buildTIFF(bo, ifd0, nil, gps)), false)
	if err != nil {
		t.Fatalf("ReadEXIF() error = %v", err)
	}

	expected := map[string]string{
		"GPSLatitude":    `44° 48' 45.00" N`,
		"GPSLongitude":   `20° 27' 40.32" E`,
		"GPSCoordinates": "44.812500, 20.461200",
	}
	for tag, want := range expected {
		if got, _ := lookup(tags, tag); got != want {
			t.Errorf("%s = %q, want %q", tag, got, want)
		}
	}
}

func TestReadEXIFUnknownTags(t *testing.T) {
	data := wrapJPEG(sampleTIFF(binary.LittleEndian))

	tags, err := ReadEXIF(data, false)
	if err != nil {
		t.Fatalf("ReadEXIF() error = %v", err)
	}
	if _, ok := lookup(tags, "Tag-0x9999"); ok {
		t.Error("unknown tag listed without all")
	}

	tags, err = ReadEXIF(data, true)
	if err != nil {
		t.Fatalf("ReadEXIF() error = %v", err)
	}
	if _, ok := lookup(tags, "Tag-0x9999"); !ok {
		t.Error("unknown tag missing with all")
	}
}

func TestReadEXIFContainers(t *testing.T) {
	tiff := sampleTIFF(binary.BigEndian)

	tests := []struct {
		name string
		data []byte
	}{
		{"jpeg", wrapJPEG(tiff)},
		{"png", wrapPNG(tiff)},
		{"tiff", tiff},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tags, err := ReadEXIF(tt.data, false)
			if err != nil {
				t.Fatalf("ReadEXIF() error = %v", err)
			}
			if got, _ := lookup(tags, "Make"); got != "Canon" {
				t.Errorf("Make = %q, want %q", got, "Canon")
			}
		})
	}
}

func TestReadEXIFErrors(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{"empty", nil},
		{"not an image", []byte("module orbx")},
		{"jpeg without exif", []byte{0xFF, 0xD8, 0xFF, 0xDA, 0x00, 0x02}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := ReadEXIF(tt.data, false); err == nil {
				t.Error("expected error, got nil")
			}
		})
	}
}
