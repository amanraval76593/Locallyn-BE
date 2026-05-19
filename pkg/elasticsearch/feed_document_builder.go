package elasticsearch

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseWKTPoint(value string) (GeoPoint, error) {
	value = strings.TrimSpace(value)
	upperValue := strings.ToUpper(value)

	if !strings.HasPrefix(upperValue, "POINT") {
		return GeoPoint{}, fmt.Errorf("unsupported WKT type: %q", value)
	}

	openIndex := strings.Index(value, "(")
	closeIndex := strings.LastIndex(value, ")")
	if openIndex == -1 || closeIndex == -1 || closeIndex <= openIndex {
		return GeoPoint{}, fmt.Errorf("invalid WKT point: %q", value)
	}

	parts := strings.Fields(value[openIndex+1 : closeIndex])
	if len(parts) != 2 {
		return GeoPoint{}, fmt.Errorf("invalid WKT point coordinates: %q", value)
	}

	lon, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return GeoPoint{}, fmt.Errorf("parse longitude: %w", err)
	}

	lat, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return GeoPoint{}, fmt.Errorf("parse latitude: %w", err)
	}

	return GeoPoint{
		Lat: lat,
		Lon: lon,
	}, nil
}
