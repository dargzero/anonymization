package anonbll

import (
	"bitbucket.org/dargzero/k-anon"
	"bitbucket.org/dargzero/k-anon/generalization"
	"bitbucket.org/dargzero/k-anon/model"
	"errors"
	"fmt"
	"github.com/dargzero/anonymization/anondb"
	"github.com/dargzero/anonymization/anonmodel"
	"github.com/globalsign/mgo/bson"
)

const typeNumeric = "numeric"
const typePrefix = "prefix"
const numericTypeInt = "int"
const numericTypeFloat = "float"
const defaultMaxWords = 100
const missingRequiredOption = "option %v is required for numeric field when using the graph algorithm"

// graphAnonymizer is a wrapper for the graph based k-anon anonymization library
type graphAnonymizer struct {
	dataset   *anonmodel.Dataset
	qidFields []anonmodel.FieldAnonymizationInfo
}

func (g *graphAnonymizer) initialize(dataset *anonmodel.Dataset, name string, qidFields []anonmodel.FieldAnonymizationInfo) {
	g.dataset = dataset
	g.qidFields = qidFields
}

func (g *graphAnonymizer) anonymize() (err error) {
	if err = g.validateExtraOptions(); err != nil {
		return
	}

	schema := g.getSchema()
	var table *model.Table
	table, err = g.getTable(schema)
	if err != nil {
		return
	}

	anonymizer := &k_anon.Anonymizer{
		K:     g.dataset.Settings.K,
		Table: table,
	}

	anonymizer.Anonymize()

	// TODO: persist data

	return errors.New("not implemented")
}

func (g *graphAnonymizer) getTable(schema *model.Schema) (table *model.Table, err error) {
	table = model.NewTable(schema)
	var data []bson.M
	if data, err = anondb.FetchUnanonymizedData(g.dataset.Name); err != nil {
		return
	}
	for _, doc := range data {
		var row []interface{}
		for _, col := range schema.Columns {
			val := doc[col.GetName()]
			row = append(row, val)
		}
		table.AddRow(row...)
	}
	return
}

func (g *graphAnonymizer) getSchema() *model.Schema {
	schema := &model.Schema{}
	for _, field := range g.qidFields {
		gen := createGeneralizer(&field)
		var column *model.Column
		if field.Type == typePrefix {
			column = model.NewWeightedColumn(field.Name, gen, 2.0)
		} else {
			column = model.NewColumn(field.Name, gen)
		}
		schema.Columns = append(schema.Columns, column)
	}
	return schema
}

func createGeneralizer(field *anonmodel.FieldAnonymizationInfo) generalization.Generalizer {
	if field.Type == typeNumeric {
		return createNumericGeneralizer(field.Opts)
	}
	return createPrefixGeneralizer(field.Opts)
}

func createPrefixGeneralizer(opts anonmodel.ExtraFieldOptions) generalization.Generalizer {
	var maxWords int
	if opts.Max == nil {
		maxWords = defaultMaxWords
	} else {
		maxWords = int(*opts.Max)
	}
	return &generalization.PrefixGeneralizer{MaxWords: maxWords}
}

func createNumericGeneralizer(opts anonmodel.ExtraFieldOptions) generalization.Generalizer {
	if opts.Type == numericTypeInt {
		min := int(*opts.Min)
		max := int(*opts.Max)
		return generalization.NewIntRangeGeneralizer(min, max)
	}
	return generalization.NewFloatRangeGeneralizer(*opts.Min, *opts.Max)
}

func (g *graphAnonymizer) validateExtraOptions() error {
	for _, field := range g.qidFields {
		if field.Type == typeNumeric {
			subType := field.Opts.Type
			if subType != numericTypeFloat && subType != numericTypeInt && subType != "" {
				return errors.New("unknown numeric type: " + subType)
			}
			if field.Opts.Min == nil {
				return errors.New(fmt.Sprintf(missingRequiredOption, "min"))
			}
			if field.Opts.Max == nil {
				return errors.New(fmt.Sprintf(missingRequiredOption, "max"))
			}
		}
	}
	return nil
}
