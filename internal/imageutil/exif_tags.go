package imageutil

var ifd0Tags = map[uint16]string{
	0x0100: "ImageWidth",
	0x0101: "ImageHeight",
	0x0102: "BitsPerSample",
	0x0103: "Compression",
	0x0106: "PhotometricInterpretation",
	0x010E: "ImageDescription",
	0x010F: "Make",
	0x0110: "Model",
	0x0112: "Orientation",
	0x0115: "SamplesPerPixel",
	0x011A: "XResolution",
	0x011B: "YResolution",
	0x011C: "PlanarConfiguration",
	0x0128: "ResolutionUnit",
	0x0131: "Software",
	0x0132: "DateTime",
	0x013B: "Artist",
	0x0213: "YCbCrPositioning",
	0x8298: "Copyright",
}

var exifTags = map[uint16]string{
	0x829A: "ExposureTime",
	0x829D: "FNumber",
	0x8822: "ExposureProgram",
	0x8824: "SpectralSensitivity",
	0x8827: "ISO",
	0x8830: "SensitivityType",
	0x8832: "RecommendedExposureIndex",
	0x9000: "ExifVersion",
	0x9003: "DateTimeOriginal",
	0x9004: "DateTimeDigitized",
	0x9010: "OffsetTime",
	0x9011: "OffsetTimeOriginal",
	0x9012: "OffsetTimeDigitized",
	0x9101: "ComponentsConfiguration",
	0x9102: "CompressedBitsPerPixel",
	0x9201: "ShutterSpeedValue",
	0x9202: "ApertureValue",
	0x9203: "BrightnessValue",
	0x9204: "ExposureBiasValue",
	0x9205: "MaxApertureValue",
	0x9206: "SubjectDistance",
	0x9207: "MeteringMode",
	0x9208: "LightSource",
	0x9209: "Flash",
	0x920A: "FocalLength",
	0x9286: "UserComment",
	0x9290: "SubSecTime",
	0x9291: "SubSecTimeOriginal",
	0x9292: "SubSecTimeDigitized",
	0xA000: "FlashpixVersion",
	0xA001: "ColorSpace",
	0xA002: "PixelXDimension",
	0xA003: "PixelYDimension",
	0xA20E: "FocalPlaneXResolution",
	0xA20F: "FocalPlaneYResolution",
	0xA210: "FocalPlaneResolutionUnit",
	0xA215: "ExposureIndex",
	0xA217: "SensingMethod",
	0xA300: "FileSource",
	0xA301: "SceneType",
	0xA401: "CustomRendered",
	0xA402: "ExposureMode",
	0xA403: "WhiteBalance",
	0xA404: "DigitalZoomRatio",
	0xA405: "FocalLengthIn35mmFilm",
	0xA406: "SceneCaptureType",
	0xA407: "GainControl",
	0xA408: "Contrast",
	0xA409: "Saturation",
	0xA40A: "Sharpness",
	0xA40C: "SubjectDistanceRange",
	0xA420: "ImageUniqueID",
	0xA430: "CameraOwnerName",
	0xA431: "BodySerialNumber",
	0xA432: "LensSpecification",
	0xA433: "LensMake",
	0xA434: "LensModel",
	0xA435: "LensSerialNumber",
}

var gpsTags = map[uint16]string{
	0x0000: "GPSVersionID",
	0x0001: "GPSLatitudeRef",
	0x0002: "GPSLatitude",
	0x0003: "GPSLongitudeRef",
	0x0004: "GPSLongitude",
	0x0005: "GPSAltitudeRef",
	0x0006: "GPSAltitude",
	0x0007: "GPSTimeStamp",
	0x0008: "GPSSatellites",
	0x0009: "GPSStatus",
	0x000A: "GPSMeasureMode",
	0x000B: "GPSDOP",
	0x000C: "GPSSpeedRef",
	0x000D: "GPSSpeed",
	0x000E: "GPSTrackRef",
	0x000F: "GPSTrack",
	0x0010: "GPSImgDirectionRef",
	0x0011: "GPSImgDirection",
	0x0012: "GPSMapDatum",
	0x001B: "GPSProcessingMethod",
	0x001C: "GPSAreaInformation",
	0x001D: "GPSDateStamp",
	0x001E: "GPSDifferential",
}

var enums = map[uint16]map[int64]string{
	0x0103: { // Compression
		1: "Uncompressed", 6: "JPEG (thumbnail)", 7: "JPEG", 8: "Adobe Deflate",
	},
	0x0106: { // PhotometricInterpretation
		0: "WhiteIsZero", 1: "BlackIsZero", 2: "RGB", 3: "Palette", 6: "YCbCr",
	},
	0x0112: { // Orientation
		1: "Horizontal (normal)",
		2: "Mirror horizontal",
		3: "Rotate 180",
		4: "Mirror vertical",
		5: "Mirror horizontal and rotate 270 CW",
		6: "Rotate 90 CW",
		7: "Mirror horizontal and rotate 90 CW",
		8: "Rotate 270 CW",
	},
	0x011C: { // PlanarConfiguration
		1: "Chunky", 2: "Planar",
	},
	0x0128: { // ResolutionUnit
		1: "None", 2: "inches", 3: "cm",
	},
	0x0213: { // YCbCrPositioning
		1: "Centered", 2: "Co-sited",
	},
	0x8822: { // ExposureProgram
		0: "Not defined", 1: "Manual", 2: "Program AE", 3: "Aperture priority",
		4: "Shutter priority", 5: "Creative (slow speed)", 6: "Action (high speed)",
		7: "Portrait", 8: "Landscape", 9: "Bulb",
	},
	0x8830: { // SensitivityType
		0: "Unknown", 1: "Standard output sensitivity", 2: "Recommended exposure index",
		3: "ISO speed", 4: "SOS and REI", 5: "SOS and ISO", 6: "REI and ISO", 7: "SOS, REI and ISO",
	},
	0x9207: { // MeteringMode
		0: "Unknown", 1: "Average", 2: "Center-weighted average", 3: "Spot",
		4: "Multi-spot", 5: "Multi-segment", 6: "Partial", 255: "Other",
	},
	0x9208: { // LightSource
		0: "Unknown", 1: "Daylight", 2: "Fluorescent", 3: "Tungsten", 4: "Flash",
		9: "Fine weather", 10: "Cloudy", 11: "Shade", 12: "Daylight fluorescent",
		13: "Day white fluorescent", 14: "Cool white fluorescent", 15: "White fluorescent",
		17: "Standard light A", 18: "Standard light B", 19: "Standard light C",
		20: "D55", 21: "D65", 22: "D75", 23: "D50", 24: "ISO studio tungsten", 255: "Other",
	},
	0xA001: { // ColorSpace
		1: "sRGB", 2: "Adobe RGB", 0xFFFF: "Uncalibrated",
	},
	0xA210: { // FocalPlaneResolutionUnit
		1: "None", 2: "inches", 3: "cm", 4: "mm", 5: "um",
	},
	0xA217: { // SensingMethod
		1: "Not defined", 2: "One-chip color area", 3: "Two-chip color area",
		4: "Three-chip color area", 5: "Color sequential area", 7: "Trilinear",
		8: "Color sequential linear",
	},
	0xA300: { // FileSource
		1: "Film scanner", 2: "Reflection print scanner", 3: "Digital camera",
	},
	0xA301: { // SceneType
		1: "Directly photographed",
	},
	0xA401: { // CustomRendered
		0: "Normal", 1: "Custom",
	},
	0xA402: { // ExposureMode
		0: "Auto", 1: "Manual", 2: "Auto bracket",
	},
	0xA403: { // WhiteBalance
		0: "Auto", 1: "Manual",
	},
	0xA406: { // SceneCaptureType
		0: "Standard", 1: "Landscape", 2: "Portrait", 3: "Night",
	},
	0xA407: { // GainControl
		0: "None", 1: "Low gain up", 2: "High gain up", 3: "Low gain down", 4: "High gain down",
	},
	0xA408: { // Contrast
		0: "Normal", 1: "Low", 2: "High",
	},
	0xA409: { // Saturation
		0: "Normal", 1: "Low", 2: "High",
	},
	0xA40A: { // Sharpness
		0: "Normal", 1: "Soft", 2: "Hard",
	},
	0xA40C: { // SubjectDistanceRange
		0: "Unknown", 1: "Macro", 2: "Close", 3: "Distant",
	},
}

var formatters = map[uint16]func(*exifReader, entry) string{
	0x829A: formatExposureTime,
	0x829D: prefix("f/"),
	0x9000: formatVersion,
	0x9101: formatComponents,
	0x9201: formatAPEXShutter,
	0x9202: formatAPEXAperture,
	0x9203: suffix(" EV"),
	0x9204: signed(" EV"),
	0x9205: formatAPEXAperture,
	0x9206: suffix(" m"),
	0x9209: formatFlash,
	0x920A: suffix(" mm"),
	0x9286: formatUserComment,
	0xA000: formatVersion,
	0xA405: suffix(" mm"),
	0xA432: formatLensSpec,
}

var gpsFormatters = map[uint16]func(*exifReader, entry) string{
	0x0006: suffix(" m"),
	0x000D: suffix(""),
	0x000F: suffix("°"),
	0x0011: suffix("°"),
}

var gpsEnums = map[uint16]map[int64]string{
	0x0005: { // GPSAltitudeRef
		0: "Above sea level", 1: "Below sea level",
	},
	0x001E: { // GPSDifferential
		0: "No correction", 1: "Differential corrected",
	},
}
