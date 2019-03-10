package anonbll

import (
	"fmt"
	"testing"
)

func TestParseCoordinate(t *testing.T) {

	t.Run("parse DD coordinates", func(t *testing.T) {
		lat, lon, _ := ParseCoordinate("-47.35°, 85.55°")
		if lat != -47.35 || lon != 85.55 {
			t.Errorf("got (%v, %v)", lat, lon)
		}
	})

	t.Run("parse DMS coordinates", func(t *testing.T) {
		lat, lon, _ := ParseCoordinate("S 47° 25' 75\", W 85° 25' 75\"")
		if lat != -47.625 || lon != -85.625 {
			t.Errorf("got (%v, %v)", lat, lon)
		}
	})

	t.Run("parse DMS partial", func(t *testing.T) {
		lat, lon, _ := ParseCoordinate("S 47°, W 85°")
		if lat != -47 || lon != -85 {
			t.Errorf("got (%v, %v)", lat, lon)
		}
	})

	t.Run("parse DMS semi-partial", func(t *testing.T) {
		lat, lon, _ := ParseCoordinate("S 47° 25', W 85° 25'")
		if lat != -47.416666666666664 || lon != -85.41666666666667 {
			t.Errorf("got (%v, %v)", lat, lon)
		}
	})

	t.Run("parse error", func(t *testing.T) {
		_, _, err := ParseCoordinate("invalid format")
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})

	overflows := []string{
		"%s° 50°",
		"50° %s°",
		"S %s° 25' 75\", W 85° 25' 75\"",
		"S 47° %s' 75\", W 85° 25' 75\"",
		"S 47° 25' %s\", W 85° 25' 75\"",
		"S 47° 25' 75\", W %s° 25' 75\"",
		"S 47° 25' 75\", W 85° %s' 75\"",
		"S 47° 25' 75\", W 85° 25' %s\"",
	}
	for _, test := range overflows {
		t.Run("convert error", func(t *testing.T) {
			coordinate := fmt.Sprintf(test, getBigNumberString())
			lat, lon, err := ParseCoordinate(coordinate)
			if err == nil {
				t.Errorf("expected error, got (%f, %f)", lat, lon)
			}
		})
	}
}

func getBigNumberString() string {
	s := "1"
	for i := 0; i < 100; i++ {
		s += "00000"
	}
	return s
}
