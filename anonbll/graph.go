package anonbll

import (
	"bitbucket.org/dargzero/k-anon/model"
	"errors"
	"github.com/dargzero/anonymization/anonmodel"
)

const modeQid = "modeQid"
const typeNumeric = "typeNumeric"
const typePrefix = "prefix"

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
		column := &model.Column{
			Name:        field.Name,
			Generalizer: nil, // TODO
		}
		schema.Columns = append(schema.Columns, column)
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
