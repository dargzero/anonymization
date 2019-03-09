package anondb

import (
	"anonmodel"
	"regexp"
	"strings"

	"github.com/globalsign/mgo/bson"
)

type prefixBoundary struct {
	*anonmodel.PrefixBoundary
}

func (p *prefixBoundary) setMatch(fieldName string, match *[]bson.M) {
	if p.Prefix != "" {
		*match = append(*match, bson.M{
			fieldName: bson.M{"$regex": "^" + getRegexForPrefix(p.Prefix)},
		})
	}

	for filter := range p.Filters {
		regexString := "^(?!" + getRegexForPrefix(filter) + ")"
		*match = append(*match, bson.M{
			fieldName: bson.M{"$regex": regexString},
		})
	}
}

func getRegexForPrefix(prefix string) string {
	return regexp.QuoteMeta(prefix) + "(\\s|$)"
}

func (p *prefixBoundary) setAggregation(fieldName string, mainGroup bson.M, facets bson.M) {
	facets[fieldName] = []bson.M{
		bson.M{
			"$group": bson.M{
				"_id":   "$" + fieldName,
				"count": bson.M{"$sum": 1},
			},
		},
	}
}

func (p *prefixBoundary) getResult(fieldName string, mainGroupResult bson.M, queryResult bson.M) interface{} {
	prefixGroups := queryResult[fieldName].([]interface{})

	prefixes := make(map[string]int)
	for _, prefixGroup := range prefixGroups {
		group := prefixGroup.(bson.M)
		prefixes[strings.TrimRight(group["_id"].(string), " ")] = group["count"].(int)
	}

	return prefixes
}
