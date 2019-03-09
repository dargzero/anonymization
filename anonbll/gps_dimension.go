package anonbll

import (
	"anonmodel"
	"log"
)

type gpsDimension struct {
	anonCollectionName string
	fieldName          string
	originalRange      anonmodel.GPSArea
	currentRange       anonmodel.GPSArea
}

func (d *gpsDimension) initialize(anonCollectionName string, fieldName string) {
	d.anonCollectionName = anonCollectionName
	d.fieldName = fieldName
}

func (d *gpsDimension) getLatitude() numericDimension {
	return numericDimension{
		anonCollectionName: d.anonCollectionName,
		fieldName:          d.fieldName + ".latitude",
		originalRange:      d.originalRange.Latitude,
		currentRange:       d.currentRange.Latitude,
	}
}

func (d *gpsDimension) getLongitude() numericDimension {
	return numericDimension{
		anonCollectionName: d.anonCollectionName,
		fieldName:          d.fieldName + ".longitude",
		originalRange:      d.originalRange.Longitude,
		currentRange:       d.currentRange.Longitude,
	}
}

func (d *gpsDimension) getInitialBoundaries() anonmodel.Boundary {
	bound := anonmodel.NumericBoundary{
		LowerBound:          nil,
		LowerBoundInclusive: false,
		UpperBound:          nil,
		UpperBoundInclusive: false,
	}
	return &anonmodel.GPSBoundary{
		Latitude:  bound,
		Longitude: bound,
	}
}

func (d *gpsDimension) getDimensionForStatistics(stat interface{}, firstRun bool) mondrianDimension {
	if firstRun {
		d.originalRange = stat.(anonmodel.GPSArea)
	}
	return &gpsDimension{
		anonCollectionName: d.anonCollectionName,
		fieldName:          d.fieldName,
		originalRange:      d.originalRange,
		currentRange:       stat.(anonmodel.GPSArea),
	}
}

func prep(r *anonmodel.NumericRange, b *anonmodel.NumericBoundary) {
	if b.LowerBound == nil || r.Min > *b.LowerBound {
		b.LowerBound = &r.Min
		b.LowerBoundInclusive = true
	}
	if b.UpperBound == nil || r.Max < *b.UpperBound {
		b.UpperBound = &r.Max
		b.UpperBoundInclusive = true
	}
}

func (d *gpsDimension) prepare(partition anonmodel.Partition, count int) {
	boundary := partition[d.fieldName].(*anonmodel.GPSBoundary)
	prep(&d.currentRange.Latitude, &boundary.Latitude)
	prep(&d.currentRange.Longitude, &boundary.Longitude)
}

func (d *gpsDimension) getNormalizedRange() float64 {
	return d.currentRange.GetRelativeArea(&d.originalRange)
}

func (d *gpsDimension) wrapLongitude(yes bool, left anonmodel.Partition, right anonmodel.Partition, err error) (bool, anonmodel.Partition, anonmodel.Partition, error) {
	if !yes || err != nil {
		return yes, nil, nil, err
	}
	left[d.fieldName].(*anonmodel.GPSBoundary).Longitude = *(left[d.fieldName+".longitude"].(*anonmodel.NumericBoundary))
	right[d.fieldName].(*anonmodel.GPSBoundary).Longitude = *(right[d.fieldName+".longitude"].(*anonmodel.NumericBoundary))
	delete(left, d.fieldName+".longitude")
	delete(right, d.fieldName+".longitude")
	return yes, left, right, err
}

func (d *gpsDimension) wrapLatitude(yes bool, left anonmodel.Partition, right anonmodel.Partition, err error) (bool, anonmodel.Partition, anonmodel.Partition, error) {
	if !yes || err != nil {
		return yes, nil, nil, err
	}
	left[d.fieldName].(*anonmodel.GPSBoundary).Latitude = *(left[d.fieldName+".latitude"].(*anonmodel.NumericBoundary))
	right[d.fieldName].(*anonmodel.GPSBoundary).Latitude = *(right[d.fieldName+".latitude"].(*anonmodel.NumericBoundary))
	delete(left, d.fieldName+".latitude")
	delete(right, d.fieldName+".latitude")
	return yes, left, right, err
}

func (d *gpsDimension) tryGetAllowableCut(k int, partition anonmodel.Partition, count int) (bool, anonmodel.Partition, anonmodel.Partition, error) {
	if partition == nil {
		log.Printf("no partition to anonymize")
	}
	if d.currentRange.Latitude.GetNormalizedRange(&anonmodel.NumericRange{Min: -90, Max: 90}) > d.currentRange.Longitude.GetNormalizedRange(&anonmodel.NumericRange{Min: -90, Max: 90}) {
		yes, left, right, err := d.cutLatitude(k, partition, count)
		if yes {
			return yes, left, right, err
		}
		return d.cutLongitude(k, partition, count)
	}
	yes, left, right, err := d.cutLongitude(k, partition, count)
	if yes {
		return yes, left, right, err
	}
	return d.cutLatitude(k, partition, count)
}

func (d *gpsDimension) cutLatitude(k int, partition anonmodel.Partition, count int) (bool, anonmodel.Partition, anonmodel.Partition, error) {
	n := d.getLatitude()
	p1 := partition.Clone()
	p1[d.fieldName+".latitude"] = &(partition[d.fieldName].(*anonmodel.GPSBoundary).Latitude)
	yes, left, right, err := n.tryGetAllowableCut(k, p1, count)
	return d.wrapLatitude(yes, left, right, err)
	//return false, nil, nil, nil*/

	/*
		if d.currentRange.Latitude.Max == d.currentRange.Latitude.Min {
			return false, nil, nil, nil
		}
		median, err := anondb.GetMedian(d.anonCollectionName, d.fieldName+".latitude", partition, count)
		if err != nil {
			return false, nil, nil, err
		}

		left := partition.Clone()
		leftBoundary := left[d.fieldName+".latitude"].(*anonmodel.GPSBoundary)
		leftBoundary.Latitude.UpperBound = &median
		leftBoundary.Latitude.UpperBoundInclusive = false
		if count, err := anondb.GetCount(d.anonCollectionName, left); err != nil {
			return false, nil, nil, err
		} else if count < k {
			return false, nil, nil, nil
		}

		right := partition.Clone()
		rightBoundary := right[d.fieldName+".latitude"].(*anonmodel.GPSBoundary)
		rightBoundary.Latitude.LowerBound = &median
		if leftBoundary.Latitude.LowerBound != nil && *leftBoundary.Latitude.LowerBound == *leftBoundary.Latitude.UpperBound {
			rightBoundary.Latitude.LowerBoundInclusive = false
		} else {
			rightBoundary.Latitude.LowerBoundInclusive = true
		}
		if count, err := anondb.GetCount(d.anonCollectionName, right); err != nil {
			return false, nil, nil, err
		} else if count < k {
			return false, nil, nil, nil
		}

		return true, left, right, nil
		return false, nil, nil, nil*/
}

func (d *gpsDimension) cutLongitude(k int, partition anonmodel.Partition, count int) (bool, anonmodel.Partition, anonmodel.Partition, error) {
	n := d.getLongitude()
	p1 := partition.Clone()
	p1[d.fieldName+".longitude"] = &(partition[d.fieldName].(*anonmodel.GPSBoundary).Longitude)
	yes, left, right, err := n.tryGetAllowableCut(k, p1, count)
	return d.wrapLongitude(yes, left, right, err)
	/*if d.currentRange.Longitude.Max == d.currentRange.Longitude.Min {
		return false, nil, nil, nil
	}
	median, err := anondb.GetMedian(d.anonCollectionName, d.fieldName+".longitude", partition, count)
	if err != nil {
		return false, nil, nil, err
	}

	left := partition.Clone()
	leftBoundary := left[d.fieldName+".longitude"].(*anonmodel.GPSBoundary)
	leftBoundary.Longitude.UpperBound = &median
	leftBoundary.Longitude.UpperBoundInclusive = false
	if count, err := anondb.GetCount(d.anonCollectionName, left); err != nil {
		return false, nil, nil, err
	} else if count < k {
		return false, nil, nil, nil
	}

	right := partition.Clone()
	rightBoundary := right[d.fieldName+".longitude"].(*anonmodel.GPSBoundary)
	rightBoundary.Longitude.LowerBound = &median
	if leftBoundary.Longitude.LowerBound != nil && *leftBoundary.Longitude.LowerBound == *leftBoundary.Longitude.UpperBound {
		rightBoundary.Longitude.LowerBoundInclusive = false
	} else {
		rightBoundary.Longitude.LowerBoundInclusive = true
	}
	if count, err := anondb.GetCount(d.anonCollectionName, right); err != nil {
		return false, nil, nil, err
	} else if count < k {
		return false, nil, nil, nil
	}

	return true, left, right, nil
	return false, nil, nil, nil*/
}

/**/
