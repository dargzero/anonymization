package anonmodel

import "strconv"

//TypeConversionFunc used for converting types
type TypeConversionFunc func(interface{}) (interface{}, error)

//GPSBoundary represents a boundary for a coordinate type field
type GPSBoundary struct {
	Latitude  NumericBoundary
	Longitude NumericBoundary
}

//GPSArea  holds the range for a gps range
type GPSArea struct {
	Latitude  NumericRange
	Longitude NumericRange
}

//GetRelativeArea calculates the relative area of coords
func (a *GPSArea) GetRelativeArea(original *GPSArea) float64 {
	return a.Latitude.GetNormalizedRange(&original.Latitude) *
		a.Longitude.GetNormalizedRange(&original.Longitude)
}

//Clone ...
func (b *GPSBoundary) Clone() Boundary {
	return &GPSBoundary{Latitude: b.Latitude,
		Longitude: b.Longitude}
}

//GetGeneralizedValue gv aaaa
func (b *GPSBoundary) GetGeneralizedValue() string {
	result := strconv.FormatFloat(*b.Latitude.UpperBound, 'f', -1, 64) + ":" + strconv.FormatFloat(*b.Latitude.LowerBound, 'f', -1, 64) + ", "
	result += strconv.FormatFloat(*b.Longitude.UpperBound, 'f', -1, 64) + ":" + strconv.FormatFloat(*b.Longitude.LowerBound, 'f', -1, 64)
	return result
}

//Bucket ...
type Bucket struct {
	Count int64
	Min   float64
	Max   float64
}
