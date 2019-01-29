package anondb

import (
	"anonmodel"
	"log"

	"github.com/globalsign/mgo/bson"
)

type gpsBoundary struct {
	*anonmodel.GPSBoundary
}

//GetCoordStatistics ...
func GetCoordStatistics(anonCollectionName string, fieldName string, partition anonmodel.Partition, count int) (float64, error) {
	session := globalSession.Copy()
	defer session.Close()
	anon := session.DB("anondb").C(anonCollectionName)

	match, err := getMatch(partition)
	if err != nil {
		return 0, err
	}

	var result bson.M
	if err = anon.Find(match).Sort(fieldName).Skip(count / 2).Limit(1).One(&result); err != nil {
		return 0, err
	}
	return result[fieldName].(float64), nil
}

func (b *gpsBoundary) setMatch(fieldName string, match *[]bson.M) {
	session := globalSession.Copy()
	defer session.Close()
	anon := session.DB("anondb").C("data_test_set__")
	c, _ := anon.Find(bson.M{"$and": *match}).Count()
	c1, _ := anon.Find(nil).Count()
	log.Printf("%v / %v\n", c, c1)
	lat, _ := convertBoundary(&(b.Latitude))
	lat.setMatch(fieldName+".latitude", match)
	lon, _ := convertBoundary(&(b.Longitude))
	lon.setMatch(fieldName+".longitude", match)
}

func (b *gpsBoundary) setAggregation(fieldName string, mainGroup bson.M, facets bson.M) {
	mainGroup["__min_"+fieldName+"_latitude"] = bson.M{"$min": "$" + fieldName + ".latitude"}
	mainGroup["__max_"+fieldName+"_latitude"] = bson.M{"$max": "$" + fieldName + ".latitude"}
	mainGroup["__min_"+fieldName+"_longitude"] = bson.M{"$min": "$" + fieldName + ".longitude"}
	mainGroup["__max_"+fieldName+"_longitude"] = bson.M{"$max": "$" + fieldName + ".longitude"}

}

//GetEqDistribution ...
func GetEqDistribution(anonCollectionName string, fieldName string, partition anonmodel.Partition, n int, min float64, max float64) ([]anonmodel.Bucket, error) {
	session := globalSession.Copy()
	defer session.Close()
	defer session.Close()
	anon := session.DB("anondb").C(anonCollectionName)

	match, err := getMatch(partition)
	if err != nil {
		return nil, err
	}

	bounds := make([]float64, n+1)
	buckets := make([]anonmodel.Bucket, n)
	for i := 0; i <= n; i++ {
		bounds[i] = min + (max-min)/float64(n)*float64(i)
	}
	for i := 0; i < n; i++ {
		buckets[i].Min = min + (max-min)/float64(n)*float64(i)
		buckets[i].Max = min + (max-min)/float64(n)*float64(i+1)
	}
	bucketSettings := bson.M{"$bucket": bson.M{"groupBy": "$" + fieldName, "boundaries": bounds, "default": "other"}}

	pipe := anon.Pipe([]interface{}{match, bucketSettings})
	items := pipe.Iter()
	i := 0
	var result bson.M
	for items.Next(&result) {
		val, _ := result["count"].(int)
		if i < n {
			buckets[i].Count = int64(val)
		}
		i = i + 1
	}
	return buckets, nil
}

func (b *gpsBoundary) getResult(fieldName string, mainGroupResult bson.M, queryResult bson.M) interface{} {
	minLat, ok := mainGroupResult["__min_"+fieldName+"_latitude"].(float64)
	if !ok {
		log.Println("failed result")
	}
	maxLat, ok := mainGroupResult["__max_"+fieldName+"_latitude"].(float64)
	if !ok {
		log.Println("failed result")
	}
	minLon, ok := mainGroupResult["__min_"+fieldName+"_longitude"].(float64)
	if !ok {
		log.Println("failed result")
	}
	maxLon, ok := mainGroupResult["__max_"+fieldName+"_longitude"].(float64)
	if !ok {
		log.Println("failed result")
	}
	return anonmodel.GPSArea{
		Latitude:  anonmodel.NumericRange{Min: minLat, Max: maxLat},
		Longitude: anonmodel.NumericRange{Min: minLon, Max: maxLon},
	}
}
