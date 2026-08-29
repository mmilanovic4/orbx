package imageutil

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func (r *exifReader) ints(e entry) []int64 {
	size := typeSizes[e.typ]
	if size == 0 || len(e.raw) == 0 {
		return nil
	}

	var out []int64
	for i := 0; i+int(size) <= len(e.raw); i += int(size) {
		b := e.raw[i : i+int(size)]
		switch e.typ {
		case typeByte, typeUndefined:
			out = append(out, int64(b[0]))
		case typeSByte:
			out = append(out, int64(int8(b[0])))
		case typeShort:
			out = append(out, int64(r.bo.Uint16(b)))
		case typeSShort:
			out = append(out, int64(int16(r.bo.Uint16(b))))
		case typeLong:
			out = append(out, int64(r.bo.Uint32(b)))
		case typeSLong:
			out = append(out, int64(int32(r.bo.Uint32(b))))
		default:
			return nil
		}
	}
	return out
}

func (r *exifReader) rationals(e entry) [][2]int64 {
	if e.typ != typeRational && e.typ != typeSRational {
		return nil
	}

	var out [][2]int64
	for i := 0; i+8 <= len(e.raw); i += 8 {
		num := r.bo.Uint32(e.raw[i : i+4])
		den := r.bo.Uint32(e.raw[i+4 : i+8])
		if e.typ == typeSRational {
			out = append(out, [2]int64{int64(int32(num)), int64(int32(den))})
		} else {
			out = append(out, [2]int64{int64(num), int64(den)})
		}
	}
	return out
}

func (r *exifReader) floats(e entry) []float64 {
	var out []float64
	switch e.typ {
	case typeFloat:
		for i := 0; i+4 <= len(e.raw); i += 4 {
			out = append(out, float64(math.Float32frombits(r.bo.Uint32(e.raw[i:i+4]))))
		}
	case typeDouble:
		for i := 0; i+8 <= len(e.raw); i += 8 {
			out = append(out, math.Float64frombits(r.bo.Uint64(e.raw[i:i+8])))
		}
	case typeRational, typeSRational:
		for _, v := range r.rationals(e) {
			if v[1] == 0 {
				out = append(out, 0)
				continue
			}
			out = append(out, float64(v[0])/float64(v[1]))
		}
	default:
		for _, v := range r.ints(e) {
			out = append(out, float64(v))
		}
	}
	return out
}

func (r *exifReader) float(e entry) (float64, bool) {
	v := r.floats(e)
	if len(v) == 0 {
		return 0, false
	}
	return v[0], true
}

func (r *exifReader) str(e entry) string {
	s := string(e.raw)
	if i := strings.IndexByte(s, 0); i >= 0 {
		s = s[:i]
	}
	return strings.TrimSpace(s)
}

func isPrintable(b []byte) bool {
	if len(b) == 0 || len(b) > 64 {
		return false
	}
	for _, c := range b {
		if c != 0 && (c < 0x20 || c > 0x7E) {
			return false
		}
	}
	return true
}

func formatFloat(f float64) string {
	s := strconv.FormatFloat(f, 'f', 4, 64)
	s = strings.TrimRight(s, "0")
	return strings.TrimSuffix(s, ".")
}

func formatRational(num, den int64) string {
	switch {
	case den == 0:
		return "undefined"
	case num%den == 0:
		return strconv.FormatInt(num/den, 10)
	default:
		return formatFloat(float64(num) / float64(den))
	}
}

func formatExposureTime(r *exifReader, e entry) string {
	v, ok := r.float(e)
	if !ok {
		return ""
	}
	return formatSeconds(v)
}

func formatSeconds(v float64) string {
	if v <= 0 {
		return ""
	}
	if v < 1 {
		return fmt.Sprintf("1/%d s", int(math.Round(1/v)))
	}
	return formatFloat(v) + " s"
}

func suffix(unit string) func(*exifReader, entry) string {
	return func(r *exifReader, e entry) string {
		v, ok := r.float(e)
		if !ok {
			return ""
		}
		return formatFloat(v) + unit
	}
}

func prefix(text string) func(*exifReader, entry) string {
	return func(r *exifReader, e entry) string {
		v, ok := r.float(e)
		if !ok {
			return ""
		}
		return text + formatFloat(v)
	}
}

func signed(unit string) func(*exifReader, entry) string {
	return func(r *exifReader, e entry) string {
		v, ok := r.float(e)
		if !ok {
			return ""
		}
		if v > 0 {
			return "+" + formatFloat(v) + unit
		}
		return formatFloat(v) + unit
	}
}

func formatVersion(r *exifReader, e entry) string {
	s := string(e.raw)
	if len(s) != 4 {
		return ""
	}
	major := strings.TrimLeft(s[:2], "0")
	if major == "" {
		major = "0"
	}
	return major + "." + s[2:]
}

func formatFlash(r *exifReader, e entry) string {
	v := r.ints(e)
	if len(v) != 1 {
		return ""
	}

	f := v[0]
	parts := []string{"did not fire"}
	if f&0x01 != 0 {
		parts = []string{"fired"}
	}
	switch (f >> 3) & 0x03 {
	case 1:
		parts = append(parts, "compulsory")
	case 2:
		parts = append(parts, "suppressed")
	case 3:
		parts = append(parts, "auto")
	}
	if f&0x20 != 0 {
		parts = append(parts, "no flash function")
	}
	if f&0x40 != 0 {
		parts = append(parts, "red-eye reduction")
	}
	if f&0x01 != 0 && f&0x04 != 0 {
		if f&0x02 != 0 {
			parts = append(parts, "return detected")
		} else {
			parts = append(parts, "no return detected")
		}
	}

	return fmt.Sprintf("%s (0x%02X)", strings.Join(parts, ", "), f)
}

func formatUserComment(r *exifReader, e entry) string {
	if len(e.raw) <= 8 {
		return ""
	}

	code := strings.TrimRight(string(e.raw[:8]), "\x00 ")
	body := e.raw[8:]
	if code != "ASCII" && code != "" {
		return fmt.Sprintf("<%s, %d bytes>", code, len(body))
	}

	return r.str(entry{typ: typeASCII, raw: body})
}

func formatLensSpec(r *exifReader, e entry) string {
	v := r.floats(e)
	if len(v) != 4 {
		return ""
	}

	focal := formatFloat(v[0]) + "-" + formatFloat(v[1]) + "mm"
	if v[0] == v[1] {
		focal = formatFloat(v[0]) + "mm"
	}
	if v[2] == 0 && v[3] == 0 {
		return focal
	}

	aperture := "f/" + formatFloat(v[2]) + "-" + formatFloat(v[3])
	if v[2] == v[3] {
		aperture = "f/" + formatFloat(v[2])
	}

	return focal + " " + aperture
}

func formatAPEXShutter(r *exifReader, e entry) string {
	v, ok := r.float(e)
	if !ok {
		return ""
	}
	return formatSeconds(math.Pow(2, -v))
}

func formatAPEXAperture(r *exifReader, e entry) string {
	v, ok := r.float(e)
	if !ok {
		return ""
	}
	return "f/" + formatFloat(math.Round(math.Pow(2, v/2)*10)/10)
}

func formatComponents(r *exifReader, e entry) string {
	names := map[int64]string{0: "-", 1: "Y", 2: "Cb", 3: "Cr", 4: "R", 5: "G", 6: "B"}

	var parts []string
	for _, v := range r.ints(e) {
		name, ok := names[v]
		if !ok {
			return ""
		}
		if name != "-" {
			parts = append(parts, name)
		}
	}
	if len(parts) == 0 {
		return ""
	}

	return strings.Join(parts, ", ")
}
