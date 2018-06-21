package anondb

import (
	"anonmodel"

	"github.com/globalsign/mgo/bson"
)

type numericBoundary struct {
	*anonmodel.NumericBoundary
}

func (b *numericBoundary) setMatch(fieldName string, match *[]bson.M) {
	if b.LowerBound == nil && b.UpperBound == nil {
		return
	}
	if b.LowerBound != nil && b.UpperBound != nil && *b.LowerBound == *b.UpperBound {
		*match = append(*match, bson.M{fieldName: *b.LowerBound})
		return
	}

	var result = bson.M{}
	if b.LowerBound != nil {
		if b.LowerBoundInclusive {
			result["$gte"] = *b.LowerBound
		} else {
			result["$gt"] = *b.LowerBound
		}
	}
	if b.UpperBound != nil {
		if b.UpperBoundInclusive {
			result["$lte"] = *b.UpperBound
		} else {
			result["$lt"] = *b.UpperBound
		}
	}
	*match = append(*match, bson.M{fieldName: result})
}

func (b *numericBoundary) setAggregation(fieldName string, mainGroup bson.M, facets bson.M) {
	mainGroup["min_"+fieldName] = bson.M{"$min": "$" + fieldName}
	mainGroup["max_"+fieldName] = bson.M{"$max": "$" + fieldName}
}

func (b *numericBoundary) getResult(fieldName string, mainGroupResult bson.M, queryResult bson.M) interface{} {
	min := mainGroupResult["min_"+fieldName].(float64)
	max := mainGroupResult["max_"+fieldName].(float64)
	return anonmodel.NumericRange{Min: min, Max: max}
}

// GetMedian gets the median value of the specified numeric field
func GetMedian(anonCollectionName string, fieldName string, partition anonmodel.Partition, count int) (float64, error) {
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
