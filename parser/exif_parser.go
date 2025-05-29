package parser

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type ExifPayload struct {
	Make                  string `json:"make,omitempty"`
	Model                 string `json:"model,omitempty"`
	Orientation           string `json:"orientation,omitempty"`
	Software              string `json:"software,omitempty"`
	DateTime              string `json:"date_time,omitempty"`
	ExposureTime          string `json:"exposure_time,omitempty"`
	FNumber               string `json:"f_number,omitempty"`
	ExifOffset            string `json:"exif_offset,omitempty"`
	ISOSpeedRatings       string `json:"iso_speed_ratings,omitempty"`
	DateTimeOriginal      string `json:"date_time_original,omitempty"`
	ShutterSpeedValue     string `json:"shutter_speed_value,omitempty"`
	ApertureValue         string `json:"aperture_value,omitempty"`
	Flash                 string `json:"flash,omitempty"`
	PixelXDimension       string `json:"pixel_x_dimension,omitempty"`
	PixelYDimension       string `json:"pixel_y_dimension,omitempty"`
	FocalLengthIn35mmFilm string `json:"focal_length_in_35mm_film,omitempty"`
	ImageDescription      string `json:"image_description,omitempty"`
	XResolution           string `json:"x_resolution,omitempty"`
	YResolution           string `json:"y_resolution,omitempty"`
	ResolutionUnit        string `json:"resolution_unit,omitempty"`
	Artist                string `json:"artist,omitempty"`
	YCbCrPositioning      string `json:"y_cb_cr_positioning,omitempty"`
	Copyright             string `json:"copyright,omitempty"`
	GPSInfoIFDPointer     string `json:"gps_info_ifd_pointer,omitempty"`
	ExifVersion           string `json:"exif_version,omitempty"`
	DateTimeDigitized     string `json:"date_time_digitized,omitempty"`
	BrightnessValue       string `json:"brightness_value,omitempty"`
	ExposureBiasValue     string `json:"exposure_bias_value,omitempty"`
	MaxApertureValue      string `json:"max_aperture_value,omitempty"`
	MeteringMode          string `json:"metering_mode,omitempty"`
	FocalLength           string `json:"focal_length,omitempty"`
	GPSVersionID          string `json:"gps_version_id,omitempty"`
	GPSLatitudeRef        string `json:"gps_latitude_ref,omitempty"`
	GPSLatitude           string `json:"gps_latitude,omitempty"`
	GPSLongitudeRef       string `json:"gps_longitude_ref,omitempty"`
	GPSLongitude          string `json:"gps_longitude,omitempty"`
	GPSAltitudeRef        string `json:"gps_altitude_ref,omitempty"`
	GPSAltitude           string `json:"gps_altitude,omitempty"`
	GPSMapDatum           string `json:"gps_map_datum,omitempty"`
}

var exifTagNames = map[uint16]string{
	0x010F: "Make",  // Camera manufacturer
	0x0110: "Model", // Camera model
	0x0112: "Orientation",
	0x0131: "Software",         // Software used
	0x0132: "DateTime",         // Date and time
	0x829A: "ExposureTime",     // Exposure time
	0x829D: "FNumber",          // F number (aperture)
	0x8769: "ExifOffset",       // Offset to Exif SubIFD
	0x8827: "ISOSpeedRatings",  // ISO speed
	0x9003: "DateTimeOriginal", // Original date and time
	0x9201: "ShutterSpeedValue",
	0x9202: "ApertureValue",
	0x9209: "Flash",
	0xA002: "PixelXDimension", // Image width
	0xA003: "PixelYDimension", // Image height
	0xA405: "FocalLengthIn35mmFilm",
	0x010E: "ImageDescription",
	0x011A: "XResolution",
	0x011B: "YResolution",
	0x0128: "ResolutionUnit",
	0x013B: "Artist",
	0x0213: "YCbCrPositioning",
	0x8298: "Copyright",
	0x8825: "GPSInfoIFDPointer",
	// EXIF SubIFD
	0x9000: "ExifVersion",
	0x9004: "DateTimeDigitized",
	0x9203: "BrightnessValue",
	0x9204: "ExposureBiasValue",
	0x9205: "MaxApertureValue",
	0x9207: "MeteringMode",
	0x920A: "FocalLength",

	// GPS Related data
	0x0000: "GPSVersionID",
	0x0001: "GPSLatitudeRef",
	0x0002: "GPSLatitude",
	0x0003: "GPSLongitudeRef",
	0x0004: "GPSLongitude",
	0x0005: "GPSAltitudeRef",
	0x0006: "GPSAltitude",
	0x0012: "GPSMapDatum",
}

func ParseImageFile(imagePath string) (*ExifPayload, error) {
	log.Printf("Processing image: %s", imagePath)
	file, err := os.Open(imagePath)
	if err != nil {
		return nil, fmt.Errorf("error opening image: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("error reading image file: %w", err)
	}

	if !isJPEG(data) {
		return nil, fmt.Errorf("image is not a JPEG file")
	}

	markerResult, pos := getAPP1Marker(data)
	if pos == 0 {
		return nil, fmt.Errorf("APP1 marker not found in image")
	}

	fieldLength := calculateLengthField(markerResult)
	segment := data[pos:]
	if len(segment) < int(fieldLength) {
		return nil, fmt.Errorf("invalid segment length: expected at least %d bytes, got %d", fieldLength, len(segment))
	}
	block := readEXIFDataBlock(segment, fieldLength)

	return parseEXIFInfo(block)
}

func readEXIFDataBlock(data []byte, fieldLength uint16) []byte {
	var exifBlock []byte
	for i := range int(fieldLength) {
		exifBlock = append(exifBlock, data[i])
	}

	return exifBlock[2:]
}
func getAPP1Marker(data []byte) ([]byte, int) {
	var result []byte
	var markerStartPos int
	for i, b := range data {
		if b == 0xFF {
			app1 := data[i+1]
			// it means we are at the APP1 Segment of the JPEG file
			if app1 == 0xE1 {
				markerStartPos = i + 2
				result = append(result, data[i:i+4]...)
			}
		}
	}

	return result, markerStartPos
}

func calculateLengthField(data []byte) uint16 {

	//TODO: check the length of data, it has to be at least 4 bytes
	// if len(data) < 2 {
	// 	return 0, fmt.Errorf("data slice too short")
	// }
	lastTwo := data[len(data)-2:]
	dataLength := binary.BigEndian.Uint16(lastTwo)
	return dataLength

}

func isJPEG(data []byte) bool {
	return len(data) >= 2 && data[0] == 0xFF && data[1] == 0xD8
}

func parseEXIFInfo(data []byte) (*ExifPayload, error) {
	payload := &ExifPayload{}

	if len(data) < 14 {
		return nil, fmt.Errorf("EXIF block too small")
	}

	exifID := data[0:6]
	if string(exifID) != "Exif\x00\x00" {
		return nil, fmt.Errorf("invalid EXIF header")
	}

	tiffStart := 6
	tiffHeader := data[tiffStart : tiffStart+8]
	byteOrder := tiffHeader[0:2]

	var endian binary.ByteOrder
	switch {
	case byteOrder[0] == 0x49 && byteOrder[1] == 0x49:
		endian = binary.LittleEndian
	case byteOrder[0] == 0x4D && byteOrder[1] == 0x4D:
		endian = binary.BigEndian
	default:
		return nil, fmt.Errorf("unknown byte order")
	}

	checkValue := endian.Uint16(tiffHeader[2:4])
	if checkValue != 42 {
		return nil, fmt.Errorf("invalid TIFF magic number: got %d", checkValue)
	}

	ifdOffset := endian.Uint32(tiffHeader[4:])
	firstIFDIndex := tiffStart + int(ifdOffset)
	if firstIFDIndex+2 > len(data) {
		return nil, fmt.Errorf("IFD offset out of bounds")
	}

	ifdData := data[firstIFDIndex:]
	entryCount := endian.Uint16(ifdData[0:2])

	for i := 0; i < int(entryCount); i++ {
		entryStart := 2 + i*12
		if entryStart+12 > len(ifdData) {
			log.Printf("IFD entry %d out of bounds", i)
			continue
		}

		entry := ifdData[entryStart : entryStart+12]
		tagID := endian.Uint16(entry[0:2])
		dataFormat := endian.Uint16(entry[2:4])
		componentCount := endian.Uint32(entry[4:8])
		valueOrOffset := endian.Uint32(entry[8:12])

		valueSize := getValueSize(dataFormat, componentCount)
		if valueSize == 0 {
			continue
		}

		var tagValue []byte
		if valueSize <= 4 {
			tagValue = entry[8 : 8+valueSize]
		} else {
			absoluteOffset := tiffStart + int(valueOrOffset)
			if absoluteOffset+valueSize > len(data) {
				log.Printf("Tag value offset out of bounds for tag 0x%X", tagID)
				continue
			}
			tagValue = data[absoluteOffset : absoluteOffset+valueSize]
		}

		tagName := getTagName(tagID)
		if tagName == "" {
			continue
		}

		value := decodeTagValue(dataFormat, tagValue, endian)
		if value == "" {
			continue
		}

		assignTagValue(payload, tagName, value)
	}

	return payload, nil
}

func getTagName(tagID uint16) string {
	if name, ok := exifTagNames[tagID]; ok {
		return name
	}

	return fmt.Sprintf("Unknown(0x%X)", tagID)
}

func getValueSize(dataFormat uint16, count uint32) int {
	switch dataFormat {
	case 1, 2, 6, 7:
		return int(count)
	case 3, 8:
		return int(count) * 2
	case 4, 9:
		return int(count) * 4
	case 5, 10:
		return int(count) * 8
	default:
		log.Printf("Unsupported data format: %d", dataFormat)
		return 0
	}
}

func decodeTagValue(format uint16, data []byte, endian binary.ByteOrder) string {
	switch format {
	case 2:
		return strings.TrimRight(string(data), "\x00")
	case 3:
		if len(data) >= 2 {
			return fmt.Sprintf("%d", endian.Uint16(data))
		}
	case 4:
		if len(data) >= 4 {
			return fmt.Sprintf("%d", endian.Uint32(data))
		}
	case 5:
		if len(data) >= 8 {
			num := endian.Uint32(data[0:4])
			den := endian.Uint32(data[4:8])
			if den != 0 {
				return fmt.Sprintf("%v", float64(num)/float64(den))
			}
		}
	}
	return ""
}

func assignTagValue(payload *ExifPayload, tagName, value string) {
	switch tagName {
	case "Make":
		payload.Make = value
	case "Model":
		payload.Model = value
	case "Orientation":
		payload.Orientation = value
	case "Software":
		payload.Software = value
	case "DateTime":
		payload.DateTime = value
	case "ExposureTime":
		payload.ExposureTime = value
	case "FNumber":
		payload.FNumber = value
	case "ExifOffset":
		payload.ExifOffset = value
	case "ISOSpeedRatings":
		payload.ISOSpeedRatings = value
	case "DateTimeOriginal":
		payload.DateTimeOriginal = value
	case "ShutterSpeedValue":
		payload.ShutterSpeedValue = value
	case "ApertureValue":
		payload.ApertureValue = value
	case "Flash":
		payload.Flash = value
	case "PixelXDimension":
		payload.PixelXDimension = value
	case "PixelYDimension":
		payload.PixelYDimension = value
	case "FocalLengthIn35mmFilm":
		payload.FocalLengthIn35mmFilm = value
	case "ImageDescription":
		payload.ImageDescription = value
	case "XResolution":
		payload.XResolution = value
	case "YResolution":
		payload.YResolution = value
	case "ResolutionUnit":
		payload.ResolutionUnit = value
	case "Artist":
		payload.Artist = value
	case "YCbCrPositioning":
		payload.YCbCrPositioning = value
	case "Copyright":
		payload.Copyright = value
	case "GPSInfoIFDPointer":
		payload.GPSInfoIFDPointer = value
	case "ExifVersion":
		payload.ExifVersion = value
	case "DateTimeDigitized":
		payload.DateTimeDigitized = value
	case "BrightnessValue":
		payload.BrightnessValue = value
	case "ExposureBiasValue":
		payload.ExposureBiasValue = value
	case "MaxApertureValue":
		payload.MaxApertureValue = value
	case "MeteringMode":
		payload.MeteringMode = value
	case "FocalLength":
		payload.FocalLength = value
	case "GPSVersionID":
		payload.GPSVersionID = value
	case "GPSLatitudeRef":
		payload.GPSLatitudeRef = value
	case "GPSLatitude":
		payload.GPSLatitude = value
	case "GPSLongitudeRef":
		payload.GPSLongitudeRef = value
	case "GPSLongitude":
		payload.GPSLongitude = value
	case "GPSAltitudeRef":
		payload.GPSAltitudeRef = value
	case "GPSAltitude":
		payload.GPSAltitude = value
	case "GPSMapDatum":
		payload.GPSMapDatum = value
	}
}
