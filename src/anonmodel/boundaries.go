package anonmodel

import (
	"strconv"
)

// Boundary represents a boundary for a field
type Boundary interface {
	Clone() Boundary
	GetGeneralizedValue() string
}

// Partition represents a partition
type Partition map[string]Boundary

// Clone deep copies a partition
func (p Partition) Clone() Partition {
	copy := make(map[string]Boundary)
	for fieldName, boundary := range p {
		copy[fieldName] = boundary.Clone()
	}
	return copy
}

// NumericBoundary represents a boundary for a numeric type field
type NumericBoundary struct {
	LowerBound          *float64
	LowerBoundInclusive bool
	UpperBound          *float64
	UpperBoundInclusive bool
}

// Clone copies a boundary
func (b *NumericBoundary) Clone() Boundary {
	return &NumericBoundary{
		LowerBound:          b.LowerBound,
		LowerBoundInclusive: b.LowerBoundInclusive,
		UpperBound:          b.UpperBound,
		UpperBoundInclusive: b.UpperBoundInclusive,
	}
}

// GetGeneralizedValue gets the string representation of the given bound
func (b *NumericBoundary) GetGeneralizedValue() string {
	if b.LowerBound != nil && b.UpperBound != nil && *b.LowerBound == *b.UpperBound {
		return strconv.FormatFloat(*b.LowerBound, 'f', -1, 64)
	}

	result := ""
	if b.LowerBound == nil || !b.LowerBoundInclusive {
		result += "]"
	} else {
		result += "["
	}

	if b.LowerBound != nil {
		result += strconv.FormatFloat(*b.LowerBound, 'f', -1, 64)
	} else {
		result += "-inf"
	}
	result += ", "
	if b.UpperBound != nil {
		result += strconv.FormatFloat(*b.UpperBound, 'f', -1, 64)
	} else {
		result += "inf"
	}

	if b.UpperBound == nil || !b.UpperBoundInclusive {
		result += "["
	} else {
		result += "]"
	}
	return result
}

// NumericRange holds the range for a numeric dimension
type NumericRange struct {
	Min float64
	Max float64
}

// GetNormalizedRange calculates the normalized range
func (r *NumericRange) GetNormalizedRange(original *NumericRange) float64 {
	originalDiff := original.Max - original.Min
	if originalDiff == 0 {
		return 0
	}

	diff := r.Max - r.Min
	return diff / originalDiff
}
