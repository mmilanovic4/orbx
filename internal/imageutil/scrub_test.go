package imageutil

import (
	"bytes"
	"encoding/binary"
	"testing"
)

var scanData = []byte{0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0x11, 0x22, 0x33}

func jpegSeg(marker byte, payload []byte) []byte {
	out := []byte{0xFF, marker, 0, 0}
	binary.BigEndian.PutUint16(out[2:], uint16(len(payload)+2))
	return append(out, payload...)
}

func buildJPEG(segments ...[]byte) []byte {
	out := []byte{0xFF, 0xD8}
	for _, s := range segments {
		out = append(out, s...)
	}
	out = append(out, 0xFF, 0xDA, 0x00, 0x02)
	out = append(out, scanData...)
	return append(out, 0xFF, 0xD9)
}

func pngChunkOf(name string, payload []byte) []byte {
	out := make([]byte, 4)
	binary.BigEndian.PutUint32(out, uint32(len(payload)))
	out = append(out, name...)
	out = append(out, payload...)
	return append(out, 0, 0, 0, 0)
}

func buildPNG(chunks ...[]byte) []byte {
	out := []byte("\x89PNG\r\n\x1a\n")
	out = append(out, pngChunkOf("IHDR", make([]byte, 13))...)
	for _, c := range chunks {
		out = append(out, c...)
	}
	out = append(out, pngChunkOf("IDAT", []byte{1, 2, 3})...)
	return append(out, pngChunkOf("IEND", nil)...)
}

func webpChunk(name string, payload []byte) []byte {
	out := append([]byte(name), 0, 0, 0, 0)
	binary.LittleEndian.PutUint32(out[4:], uint32(len(payload)))
	out = append(out, payload...)
	if len(payload)%2 == 1 {
		out = append(out, 0)
	}
	return out
}

func buildWebP(flags byte, chunks ...[]byte) []byte {
	vp8x := make([]byte, 10)
	vp8x[0] = flags

	body := append([]byte("WEBP"), webpChunk("VP8X", vp8x)...)
	body = append(body, webpChunk("VP8 ", []byte{1, 2, 3, 4})...)
	for _, c := range chunks {
		body = append(body, c...)
	}

	out := append([]byte("RIFF"), 0, 0, 0, 0)
	binary.LittleEndian.PutUint32(out[4:], uint32(len(body)))
	return append(out, body...)
}

func removedNames(removed []Removed) []string {
	names := make([]string, 0, len(removed))
	for _, r := range removed {
		names = append(names, r.Name)
	}
	return names
}

func contains(values []string, want string) bool {
	for _, v := range values {
		if v == want {
			return true
		}
	}
	return false
}

func jpegScan(t *testing.T, data []byte) []byte {
	t.Helper()

	i := bytes.Index(data, []byte{0xFF, 0xDA})
	if i < 0 {
		t.Fatal("no SOS marker found")
	}
	return data[i:]
}

func sampleJPEG() []byte {
	return buildJPEG(
		jpegSeg(0xE0, []byte("JFIF\x00\x01\x02")),
		jpegSeg(0xE1, append([]byte("Exif\x00\x00"), sampleTIFF(binary.BigEndian)...)),
		jpegSeg(0xE1, []byte("http://ns.adobe.com/xap/1.0/\x00<x:xmpmeta/>")),
		jpegSeg(0xE2, []byte("ICC_PROFILE\x00\x01\x02\x03")),
		jpegSeg(0xED, []byte("Photoshop 3.0\x00")),
		jpegSeg(0xEE, []byte("Adobe\x00\x64\x00\x00\x00\x00\x02")),
		jpegSeg(0xFE, []byte("a comment")),
	)
}

func TestScrubJPEG(t *testing.T) {
	data := sampleJPEG()

	cleaned, removed, err := Scrub(data, false)
	if err != nil {
		t.Fatalf("Scrub() error = %v", err)
	}

	for _, want := range []string{"EXIF", "XMP", "IPTC", "Comment"} {
		if !contains(removedNames(removed), want) {
			t.Errorf("%s was not removed", want)
		}
	}
	if contains(removedNames(removed), "ICC profile") {
		t.Error("ICC profile removed without stripICC")
	}

	if !bytes.Contains(cleaned, []byte("ICC_PROFILE")) {
		t.Error("ICC profile segment missing from output")
	}
	if !bytes.Contains(cleaned, []byte("Adobe\x00")) {
		t.Error("Adobe color transform segment missing from output")
	}
	if !bytes.Contains(cleaned, []byte("JFIF\x00")) {
		t.Error("JFIF segment missing from output")
	}
	if bytes.Contains(cleaned, []byte("Exif\x00\x00")) {
		t.Error("EXIF segment still present in output")
	}

	if got, want := jpegScan(t, cleaned), jpegScan(t, data); !bytes.Equal(got, want) {
		t.Error("scan data changed, output is not lossless")
	}
	if _, err := ReadEXIF(cleaned, false); err == nil {
		t.Error("EXIF still readable after scrub")
	}
}

func TestScrubJPEGStripICC(t *testing.T) {
	cleaned, removed, err := Scrub(sampleJPEG(), true)
	if err != nil {
		t.Fatalf("Scrub() error = %v", err)
	}

	if !contains(removedNames(removed), "ICC profile") {
		t.Error("ICC profile was not removed with stripICC")
	}
	if bytes.Contains(cleaned, []byte("ICC_PROFILE")) {
		t.Error("ICC profile segment still present in output")
	}
}

func TestScrubPNG(t *testing.T) {
	data := buildPNG(
		pngChunkOf("eXIf", sampleTIFF(binary.BigEndian)),
		pngChunkOf("tEXt", []byte("Author\x00me")),
		pngChunkOf("iTXt", []byte("XML:com.adobe.xmp\x00<x:xmpmeta/>")),
		pngChunkOf("tIME", make([]byte, 7)),
		pngChunkOf("iCCP", []byte("profile\x00\x00data")),
		pngChunkOf("gAMA", []byte{0, 1, 2, 3}),
	)

	cleaned, removed, err := Scrub(data, false)
	if err != nil {
		t.Fatalf("Scrub() error = %v", err)
	}

	for _, want := range []string{"EXIF", "Text metadata", "XMP", "Timestamp"} {
		if !contains(removedNames(removed), want) {
			t.Errorf("%s was not removed", want)
		}
	}

	for _, chunk := range []string{"IHDR", "IDAT", "IEND", "gAMA", "iCCP"} {
		if !bytes.Contains(cleaned, []byte(chunk)) {
			t.Errorf("chunk %s missing from output", chunk)
		}
	}
	for _, chunk := range []string{"eXIf", "tEXt", "iTXt", "tIME"} {
		if bytes.Contains(cleaned, []byte(chunk)) {
			t.Errorf("chunk %s still present in output", chunk)
		}
	}
}

func TestScrubWebP(t *testing.T) {
	data := buildWebP(
		webpFlagEXIF|webpFlagXMP|webpFlagICC,
		webpChunk("ICCP", []byte{9, 9}),
		webpChunk("EXIF", sampleTIFF(binary.LittleEndian)),
		webpChunk("XMP ", []byte("<x:xmpmeta/>")),
	)

	cleaned, removed, err := Scrub(data, false)
	if err != nil {
		t.Fatalf("Scrub() error = %v", err)
	}

	if got, want := removedNames(removed), []string{"EXIF", "XMP"}; len(got) != len(want) {
		t.Fatalf("removed = %v, want %v", got, want)
	}

	size := binary.LittleEndian.Uint32(cleaned[4:8])
	if int(size) != len(cleaned)-8 {
		t.Errorf("RIFF size = %d, want %d", size, len(cleaned)-8)
	}

	if flags := cleaned[20]; flags != webpFlagICC {
		t.Errorf("VP8X flags = 0x%02X, want 0x%02X", flags, webpFlagICC)
	}
	if !bytes.Contains(cleaned, []byte("VP8 ")) || !bytes.Contains(cleaned, []byte("ICCP")) {
		t.Error("image or ICC chunk missing from output")
	}
}

func TestScrubNothingToRemove(t *testing.T) {
	_, removed, err := Scrub(buildJPEG(jpegSeg(0xE0, []byte("JFIF\x00\x01\x02"))), false)
	if err != nil {
		t.Fatalf("Scrub() error = %v", err)
	}
	if len(removed) != 0 {
		t.Errorf("removed = %v, want none", removedNames(removed))
	}
}

func TestScrubErrors(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{"tiff", sampleTIFF(binary.BigEndian)},
		{"not an image", []byte("module orbx")},
		{"empty", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, _, err := Scrub(tt.data, false); err == nil {
				t.Error("expected error, got nil")
			}
		})
	}
}
