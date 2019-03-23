package anonmodel

import "testing"

func TestGPSBoundary_Clone(t *testing.T) {
	boundary := newGPSBoundary()
	clone := boundary.Clone().(*GPSBoundary)
	if boundary == clone ||
		boundary.Latitude != clone.Latitude ||
		boundary.Longitude != clone.Longitude {
		t.Errorf("cloned object should be equal but not the same")
	}
}

func TestGPSBoundary_GetGeneralizedValue(t *testing.T) {
	boundary := newGPSBoundary()
	actual := boundary.GetGeneralizedValue()
	expected := "10.953:0.159, 8.376:5.558"
	if expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestGPSArea_GetRelativeArea(t *testing.T) {

	t.Run("0 wide range", func(t *testing.T) {
		a1 := newGPSArea()
		a2 := newGPSArea()
		a2.Longitude.Min = 5
		a2.Longitude.Max = 5

		area := a1.GetRelativeArea(a2)

		if area != 0 {
			t.Errorf("expected 0, got %v", area)
		}
	})

	t.Run("equal ranges", func(t *testing.T) {
		a1 := newGPSArea()
		a2 := newGPSArea()

		area := a1.GetRelativeArea(a2)

		if area != 1 {
			t.Errorf("expected 1, got %v", area)
		}
	})

	t.Run("normal ranges", func(t *testing.T) {
		a1 := newGPSArea()
		a2 := newGPSArea()
		a2.Longitude.Min = 2
		a2.Longitude.Max = 3
		a2.Latitude.Min = 4
		a2.Latitude.Max = 5

		area := a1.GetRelativeArea(a2)

		if area != 22.604625 {
			t.Errorf("expected 22.604625, got %v", area)
		}
	})
}

func newGPSBoundary() *GPSBoundary {
	var latL = 0.159
	var latU = 10.953
	var lonL = 5.558
	var lonU = 8.376
	return &GPSBoundary{
		Latitude: NumericBoundary{
			LowerBound:          &latL,
			LowerBoundInclusive: true,
			UpperBound:          &latU,
			UpperBoundInclusive: false,
		},
		Longitude: NumericBoundary{
			LowerBound:          &lonL,
			LowerBoundInclusive: true,
			UpperBound:          &lonU,
			UpperBoundInclusive: false,
		},
	}
}

func newGPSArea() *GPSArea {
	return &GPSArea{
		Latitude: NumericRange{
			Min: 0.570,
			Max: 5.895,
		},
		Longitude: NumericRange{
			Min: 3.748,
			Max: 7.993,
		},
	}
}
