package anonbll

import (
	"bitbucket.org/dargzero/k-anon/generalization"
	"bitbucket.org/dargzero/k-anon/model"
	"errors"
	"github.com/dargzero/anonymization/anonmodel"
	"strconv"
)

const typeNumeric = "numeric"
const typePrefix = "prefix"
const numericTypeInt = "int"
const numericTypeFloat = "float"
const optType = "type"
const optMin = "min"
const optMax = "max"
const defaultMaxWords = 100

// graphAnonymizer is a wrapper for the graph based k-anon anonymization library
type graphAnonymizer struct {
	dataset   *anonmodel.Dataset
	qidFields []*anonmodel.FieldAnonymizationInfo
}

func (g *graphAnonymizer) initialize(dataset *anonmodel.Dataset, name string, qidFields []anonmodel.FieldAnonymizationInfo) {
	g.dataset = dataset
	g.qidFields = make([]*anonmodel.FieldAnonymizationInfo, len(qidFields))
	for i, field := range qidFields {
		g.qidFields[i] = &field
	}
}

func (g *graphAnonymizer) anonymize() error {
	return errors.New("not implemented")
}

func (g *graphAnonymizer) getSchema() (*model.Schema, error) {
	schema := &model.Schema{}
	for _, field := range g.qidFields {
		g, err := createGeneralizer(field)
		if err != nil {
			return nil, err
		}
		column := model.NewColumn(field.Name, g)
		schema.Columns = append(schema.Columns, column)
	}
	return schema, nil
}

func createGeneralizer(field *anonmodel.FieldAnonymizationInfo) (generalization.Generalizer, error) {
	switch field.Type {
	case typeNumeric:
		return createNumericGeneralizer(field.Opts)
	case typePrefix:
		return createPrefixGeneralizer(field.Opts)
	}
	return nil, errors.New("unexpected field type: " + field.Type)
}

func createPrefixGeneralizer(opts map[string]string) (g generalization.Generalizer, err error) {
	var maxWords int
	_, maxDeclared := opts[optMax]
	if maxDeclared {
		if maxWords, err = getInt(opts, optMax); err != nil {
			return
		}
	} else {
		maxWords = defaultMaxWords
	}
	g = &generalization.PrefixGeneralizer{MaxWords: maxWords}
	return
}

func createNumericGeneralizer(opts map[string]string) (generalization.Generalizer, error) {
	numericType := opts[optType]
	switch numericType {
	case numericTypeInt:
		return createIntGeneralizer(opts)
	case numericTypeFloat:
		return createFloatGeneralizer(opts)
	case "":
		return createFloatGeneralizer(opts)
	}
	return nil, errors.New("unknown numeric type: " + numericType)
}

func createIntGeneralizer(opts map[string]string) (g generalization.Generalizer, err error) {
	var min, max int
	if min, err = getInt(opts, optMin); err != nil {
		return
	}
	if max, err = getInt(opts, optMax); err != nil {
		return
	}
	g = generalization.NewIntRangeGeneralizer(min, max)
	return
}

func createFloatGeneralizer(opts map[string]string) (g generalization.Generalizer, err error) {
	var min, max float64
	if min, err = getFloat(opts, optMin); err != nil {
		return
	}
	if max, err = getFloat(opts, optMax); err != nil {
		return
	}
	g = generalization.NewFloatRangeGeneralizer(min, max)
	return
}

func getInt(opts map[string]string, key string) (val int, err error) {
	var s string
	if s, err = getOpt(opts, key); err != nil {
		return
	}
	return strconv.Atoi(s)
}

func getFloat(opts map[string]string, key string) (val float64, err error) {
	var s string
	if s, err = getOpt(opts, key); err != nil {
		return
	}
	return strconv.ParseFloat(s, 64)
}

func getOpt(opts map[string]string, key string) (val string, err error) {
	var ok bool
	if val, ok = opts[key]; !ok {
		return val, errors.New("missing required option: " + key)
	}
	return
}
