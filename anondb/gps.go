package anondb

import (
	"fmt"
	"math"
	"regexp"
	"strings"
)

//GPSCoord represend a coordinate pair
type GPSCoord struct {
	Latitude  float64 `json:"latitude" bson:"latitude"`
	Longitude float64 `json:"longitude" bson:"longitude"`
	Format    string  `json:"format" bson:"format"`
}

//PreprocessCoord ...
func PreprocessCoord(coordIf interface{}) (coordI interface{}, err error) {
	coordStr, yes := coordIf.(string)
	if !yes {
		return
	}
	var coord GPSCoord
	coord.Format = FindFormat(coordStr)
	Latitude, Longitude, err := ReadCordsValue(coordStr, coord.Format)
	coord.Latitude = Latitude
	coord.Longitude = Longitude
	coordI = coord
	return
}

//FindFormat ... (coord string) (format string)
func FindFormat(coord string) (format string) {
	format = "ERR"
	DDFormat := "[-+]?[1-9][0-9]*[.]?[0-9]*°?,?[-+]?[1-9][0-9]*[.]?[0-9]*°?"
	DMSFormat := "[NS][ ]*[1-9]\\d*°[ ]*[1-9]\\d*['][ ]*[1-9]\\d*([,.]\\d*)?(\"|'')[ ,]*[WE][ ]*[1-9]\\d*°[ ]*[1-9]\\d*['][ ]*[1-9]\\d*([,.]\\d*)?(\"|'')[ \n,.]*"
	match, _ := regexp.MatchString(DDFormat, coord)
	if match {
		format = "DD"
	}
	match, _ = regexp.MatchString(DMSFormat, coord)
	if match {
		format = "DMS"
	}
	return
}

func readDMS(coords string) (Latitude float64, Longitude float64, err error) {
	var Lat, Lon uint8
	var dLat, mLat, sLat, dLon, mLon, sLon int
	coords = strings.ToUpper(coords)
	strings.Replace(coords, "''", "\"", -1)
	strings.Replace(coords, " ", "", -1)
	strings.Replace(coords, ",", "", -1)
	fmt.Sscanf(coords, "%c%d°%d'%d\"%c%d°%d'%d\"", &Lat, &dLat, &mLat, &sLat, &Lon, &dLon, &mLon, &sLon)

	Latitude = float64(dLat) + float64(mLat)/60 + float64(sLat)/360
	if Lat == 'W' {
		Latitude = -Latitude
	} else if Lat != 'E' {
		err = fmt.Errorf("Bad coord string")
		return
	}

	Longitude = float64(dLon) + float64(mLon)/60 + float64(sLon)/360
	if Lat == 'S' {
		Longitude = -Longitude
	} else if Lat != 'N' {
		err = fmt.Errorf("Bad coord string")
		return
	}
	err = nil
	return
}

func readDD(coords string) (Latitude float64, Longitude float64, err error) {
	strings.Replace(coords, " ", "", -1)
	//strings.Replace(coords, ",", "", -1)
	fmt.Sscanf(coords, "%f°,%f°", &Latitude, &Longitude)
	return
}

func readERROR(coords string) (Latitude float64, Longitude float64, err error) {
	Latitude = math.NaN()
	Longitude = math.NaN()
	err = fmt.Errorf("Unrecognised coord format")
	return
}

//ReadCordsValue a
func ReadCordsValue(coords string, format string) (Latitude float64, Longitude float64, err error) {
	if format == "DD" {
		return readDD(coords)
	} else if format == "DMS" {
		return readDMS(coords)
	} else {
		return readERROR(coords)
	}
}
