package anonbll

import (
	"log"
	"testing"
)

func TestCoordStrings(t *testing.T) {
	format := FindFormat("-47°, 85°")
	if format != "DD" {
		log.Printf("DD_ERR")
	}
	format = FindFormat("N 47° 22' 23\", W 85°  22' 23\"")
	if format != "DMS" {
		log.Printf("DMS_ERR")
	}

	lat, lon, _ := ReadCordsValue("-47°, 85°", "DD")
	if lat != -47 || lon != 85 {
		log.Printf("ERR")
	}
}
