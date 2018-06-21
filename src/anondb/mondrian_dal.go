package anondb

import (
	"anonmodel"
	"fmt"

	"github.com/globalsign/mgo/bson"
)

type dbBoundary interface {
	setMatch(string, *[]bson.M)
	setAggregation(string, bson.M, bson.M)
	getResult(string, bson.M, bson.M) interface{}
}

func convertBoundary(boundary anonmodel.Boundary) (dbBoundary, error) {
	switch boundary.(type) {
	case *anonmodel.NumericBoundary:
		return &numericBoundary{boundary.(*anonmodel.NumericBoundary)}, nil
	case *anonmodel.PrefixBoundary:
		return &prefixBoundary{boundary.(*anonmodel.PrefixBoundary)}, nil
	default:
		return nil, fmt.Errorf("Unrecognized boundary type: '%T'", boundary)
	}
}

// GetDimensionStatistics gets the statistics about the
func GetDimensionStatistics(anonCollectionName string, partition anonmodel.Partition) (int, map[string]interface{}, error) {
	session := globalSession.Copy()
	defer session.Close()
	anon := session.DB("anondb").C(anonCollectionName)

	match, err := getMatch(partition)
	if err != nil {
		return 0, nil, err
	}

	facet, err := getAggregation(partition)
	if err != nil {
		return 0, nil, err
	}

	pipeline := []bson.M{
		bson.M{"$match": match},
		bson.M{"$facet": facet},
	}
	var queryResult bson.M
	if err := anon.Pipe(pipeline).One(&queryResult); err != nil {
		return 0, nil, err
	}

	mainGroupResultArray := queryResult["mainGroup"].([]interface{})
	if len(mainGroupResultArray) != 1 {
		return 0, nil, nil
	}
	mainGroupResult := mainGroupResultArray[0].(bson.M)
	count := mainGroupResult["count"].(int)
	result, err := getResult(partition, mainGroupResult, queryResult)
	return count, result, err
}

// GetCount return the number of documents in the specified partition
func GetCount(anonCollectionName string, partition anonmodel.Partition) (int, error) {
	session := globalSession.Copy()
	defer session.Close()
	anon := session.DB("anondb").C(anonCollectionName)

	match, err := getMatch(partition)
	if err != nil {
		return 0, err
	}

	return anon.Find(match).Count()
}

// Generalize generalizes the specified partition
func Generalize(anonCollectionName string, partition anonmodel.Partition) error {
	session := globalSession.Copy()
	defer session.Close()
	anon := session.DB("anondb").C(anonCollectionName)

	match, err := getMatch(partition)
	if err != nil {
		return err
	}

	set := bson.M{"__anonymized": true}
	for fieldName, boundary := range partition {
		set["__anon_"+fieldName] = boundary.GetGeneralizedValue()
	}

	_, err = anon.UpdateAll(match, bson.M{"$set": set})
	return err
}

func getMatch(partition anonmodel.Partition) (bson.M, error) {
	match := []bson.M{
		bson.M{"__anonymized": false},
	}
	for fieldName, boundary := range partition {
		dbBound, err := convertBoundary(boundary)
		if err != nil {
			return nil, err
		}

		dbBound.setMatch(fieldName, &match)
	}

	return bson.M{"$and": match}, nil
}

func getAggregation(partition anonmodel.Partition) (bson.M, error) {
	mainGroup := bson.M{
		"_id":   nil,
		"count": bson.M{"$sum": 1},
	}
	facets := bson.M{
		"mainGroup": []bson.M{
			bson.M{"$group": mainGroup},
		},
	}
	for fieldName, boundary := range partition {
		dbBound, err := convertBoundary(boundary)
		if err != nil {
			return nil, err
		}

		dbBound.setAggregation(fieldName, mainGroup, facets)
	}

	return facets, nil
}

func getResult(partition anonmodel.Partition, mainGroupResult bson.M, queryResult bson.M) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	for fieldName, boundary := range partition {
		dbBound, err := convertBoundary(boundary)
		if err != nil {
			return nil, err
		}

		result[fieldName] = dbBound.getResult(fieldName, mainGroupResult, queryResult)
	}
	return result, nil
}
