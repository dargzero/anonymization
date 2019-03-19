package anonbll

import (
	"github.com/dargzero/anonymization/anonmodel"
	"testing"
)

func TestMondrianDimensions_Len(t *testing.T) {

	t.Run("zero length", func(t *testing.T) {
		expected := 0
		dims := newMondrianDimensions(expected)
		actual := dims.Len()
		if actual != expected {
			t.Errorf("expected %v, got %v", expected, actual)
		}
	})

	t.Run("arbitrary length", func(t *testing.T) {
		expected := 55
		dims := newMondrianDimensions(expected)
		actual := dims.Len()
		if actual != expected {
			t.Errorf("expected %v, got %v", expected, actual)
		}
	})
}

func TestMondrianDimensions_Swap(t *testing.T) {
	dims := newMondrianDimensions(10)
	d1 := &MockDim{}
	d2 := &MockDim{}
	dims[3] = d1
	dims[7] = d2

	dims.Swap(3, 7)

	if dims[3] != d2 || dims[7] != d1 {
		t.Errorf("incorrect swap")
	}
}

func TestMondrianDimensions_Less(t *testing.T) {
	dims := newMondrianDimensions(10)
	dims[0] = &MockDim{normalizedRange: 10}
	dims[1] = &MockDim{normalizedRange: 5}
	if !dims.Less(0, 1) {
		t.Errorf("expected true")
	}
	if dims.Less(1, 0) {
		t.Errorf("expected false")
	}
}

func newMondrianDimensions(length int) mondrianDimensions {
	return make([]mondrianDimension, length)
}

type MockDim struct {
	normalizedRange float64
}

func (*MockDim) initialize(anonCollectionName string, fieldName string) {
	panic("implement me")
}

func (*MockDim) getInitialBoundaries() anonmodel.Boundary {
	panic("implement me")
}

func (*MockDim) getDimensionForStatistics(interface{}, bool) mondrianDimension {
	panic("implement me")
}

func (*MockDim) prepare(partition anonmodel.Partition, count int) {
	panic("implement me")
}

func (d *MockDim) getNormalizedRange() float64 {
	return d.normalizedRange
}

func (*MockDim) tryGetAllowableCut(int, anonmodel.Partition, int) (bool, anonmodel.Partition, anonmodel.Partition, error) {
	panic("implement me")
}
