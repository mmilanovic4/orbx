package imageutil

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func (r *exifReader) decodeGPS(entries []entry, all bool) []Tag {
	byTag := make(map[uint16]entry, len(entries))
	for _, e := range entries {
		byTag[e.tag] = e
	}

	var tags []Tag
	for _, e := range entries {
		name, known := gpsTags[e.tag]
		if !known {
			if !all {
				continue
			}
			name = fmt.Sprintf("Tag-0x%04X", e.tag)
		}

		tags = append(tags, Tag{Name: name, Value: r.formatGPS(e, byTag, known)})
	}

	// a single decimal pair is what maps and search engines expect
	if lat, ok := r.coordinate(byTag, 0x0002, 0x0001); ok {
		if lon, ok := r.coordinate(byTag, 0x0004, 0x0003); ok {
			tags = append(tags, Tag{
				Name:  "GPSCoordinates",
				Value: fmt.Sprintf("%.6f, %.6f", lat, lon),
			})
		}
	}

	return tags
}

func (r *exifReader) formatGPS(e entry, byTag map[uint16]entry, known bool) string {
	if !known {
		return r.formatRaw(e)
	}

	switch e.tag {
	case 0x0000: // GPSVersionID
		var parts []string
		for _, v := range r.ints(e) {
			parts = append(parts, strconv.FormatInt(v, 10))
		}
		return strings.Join(parts, ".")
	case 0x0002: // GPSLatitude
		return r.formatDMS(e, byTag[0x0001])
	case 0x0004: // GPSLongitude
		return r.formatDMS(e, byTag[0x0003])
	case 0x0007: // GPSTimeStamp
		v := r.floats(e)
		if len(v) != 3 {
			break
		}
		return fmt.Sprintf("%02d:%02d:%05.2f UTC", int(v[0]), int(v[1]), v[2])
	case 0x001B, 0x001C: // GPSProcessingMethod, GPSAreaInformation
		return formatUserComment(r, e)
	}

	if enum, ok := gpsEnums[e.tag]; ok {
		if v := r.ints(e); len(v) == 1 {
			if name, ok := enum[v[0]]; ok {
				return name
			}
			return fmt.Sprintf("Unknown (%d)", v[0])
		}
	}

	if f, ok := gpsFormatters[e.tag]; ok {
		if v := f(r, e); v != "" {
			return v
		}
	}

	return r.formatRaw(e)
}

func (r *exifReader) formatDMS(e entry, refEntry entry) string {
	v := r.floats(e)
	if len(v) != 3 {
		return r.formatRaw(e)
	}

	ref := ""
	if refEntry.raw != nil {
		ref = " " + r.str(refEntry)
	}

	return fmt.Sprintf(`%d° %d' %.2f"%s`, int(v[0]), int(v[1]), v[2], ref)
}

func (r *exifReader) coordinate(byTag map[uint16]entry, tag, refTag uint16) (float64, bool) {
	e, ok := byTag[tag]
	if !ok {
		return 0, false
	}

	v := r.floats(e)
	if len(v) != 3 {
		return 0, false
	}

	deg := math.Abs(v[0]) + v[1]/60 + v[2]/3600
	ref := strings.ToUpper(r.str(byTag[refTag]))
	if ref == "S" || ref == "W" {
		deg = -deg
	}

	return deg, true
}
