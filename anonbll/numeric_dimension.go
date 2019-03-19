package anonbll

import (
	"github.com/dargzero/anonymization/anondb"
	"github.com/dargzero/anonymization/anonmodel"
)

type NumericDimension struct {
	anonCollectionName string
	fieldName          string
	originalRange      anonmodel.NumericRange
	currentRange       anonmodel.NumericRange
}

func (d *NumericDimension) initialize(anonCollectionName string, fieldName string) {
	d.anonCollectionName = anonCollectionName
	d.fieldName = fieldName
}

func (d *NumericDimension) getInitialBoundaries() anonmodel.Boundary {
	return &anonmodel.NumericBoundary{
		LowerBound:          nil,
		LowerBoundInclusive: false,
		UpperBound:          nil,
		UpperBoundInclusive: false,
	}
}

func (d *NumericDimension) getDimensionForStatistics(stat interface{}, firstRun bool) mondrianDimension {
	if firstRun {
		d.originalRange = stat.(anonmodel.NumericRange)
	}

	return &NumericDimension{
		anonCollectionName: d.anonCollectionName,
		fieldName:          d.fieldName,
		originalRange:      d.originalRange,
		currentRange:       stat.(anonmodel.NumericRange),
	}
}

func (d *NumericDimension) prepare(partition anonmodel.Partition, count int) {
	boundary := partition[d.fieldName].(*anonmodel.NumericBoundary)

	if boundary.LowerBound == nil || d.currentRange.Min > *boundary.LowerBound {
		boundary.LowerBound = &d.currentRange.Min
		boundary.LowerBoundInclusive = true
	}

	if boundary.UpperBound == nil || d.currentRange.Max < *boundary.UpperBound {
		boundary.UpperBound = &d.currentRange.Max
		boundary.UpperBoundInclusive = true
	}
}

func (d *NumericDimension) getNormalizedRange() float64 {
	return d.currentRange.GetNormalizedRange(&d.originalRange)
}

func (d *NumericDimension) tryGetAllowableCut(k int, partition anonmodel.Partition, count int) (bool, anonmodel.Partition, anonmodel.Partition, error) {
	if d.currentRange.Max == d.currentRange.Min {
		return false, nil, nil, nil
	}

	median, err := anondb.GetMedian(d.anonCollectionName, d.fieldName, partition, count)
	if err != nil {
		return false, nil, nil, err
	}

	left := partition.Clone()
	leftBoundary := left[d.fieldName].(*anonmodel.NumericBoundary)
	leftBoundary.UpperBound = &median
	leftBoundary.UpperBoundInclusive = false
	if count, err := anondb.GetCount(d.anonCollectionName, left); err != nil {
		return false, nil, nil, err
	} else if count < k {
		return false, nil, nil, nil
	}

	right := partition.Clone()
	rightBoundary := right[d.fieldName].(*anonmodel.NumericBoundary)
	rightBoundary.LowerBound = &median
	if leftBoundary.LowerBound != nil && *leftBoundary.LowerBound == *leftBoundary.UpperBound {
		rightBoundary.LowerBoundInclusive = false
	} else {
		rightBoundary.LowerBoundInclusive = true
	}
	if count, err := anondb.GetCount(d.anonCollectionName, right); err != nil {
		return false, nil, nil, err
	} else if count < k {
		return false, nil, nil, nil
	}

	return true, left, right, nil
}
