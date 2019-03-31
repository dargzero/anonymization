package anonbll

import (
	"bitbucket.org/dargzero/k-anon/generalization"
	"bitbucket.org/dargzero/k-anon/model"
	"errors"
	"github.com/dargzero/anonymization/anonmodel"
	"strconv"
)

const modeQid = "modeQid"
const typeNumeric = "typeNumeric"
const typePrefix = "prefix"
const optGeneralizer = "generalizer"
const optMin = "min"
const optMax = "max"

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
		if err := validate(field); err != nil {
			return nil, err
		}

		//var g generalization.Generalizer
		//switch field.Type {
		//case typeNumeric:
		//	g = generalization.NewIntRangeGeneralizer(0, 100)
		//	break
		//case typePrefix:
		//	break
		//}
	}
	return schema, nil
}

func validate(field *anonmodel.FieldAnonymizationInfo) error {
	if field.Mode != modeQid {
		return errors.New("unexpected field mode: " + field.Mode)
	}
	if field.Type != typeNumeric && field.Type != typePrefix {
		return errors.New("unexpected field type: " + field.Type)
	}
	return nil
}

func createGeneralizer(opts map[string]string) (generalization.Generalizer, error) {
	generalizer, err := getOption(opts, optGeneralizer)
	if err != nil {
		return nil, err
	}
	switch generalizer {
	case "int-range":
		return createGeneralizer(opts)
	case "float-range":
		return nil, errors.New("not implemented")
		break
	case "prefix":
		return nil, errors.New("not implemented")
		break
	}
	return nil, errors.New("unknown generalizer: " + generalizer)
}

func createIntGeneralizer(opts map[string]string) (generalization.Generalizer, error) {
	minVal, err := getOption(opts, optMin)
	maxVal, err := getOption(opts, optMax)
	if err != nil {
		return nil, err
	}
	min, err := strconv.Atoi(minVal)
	max, err := strconv.Atoi(maxVal)
	if err != nil {
		return nil, err
	}
	return generalization.NewIntRangeGeneralizer(min, max), nil
}

func getOption(opts map[string]string, key string) (string, error) {
	val := opts[key]
	if val == "" {
		return "", errors.New("missing field-option: " + key)
	}
	return val, nil
}
