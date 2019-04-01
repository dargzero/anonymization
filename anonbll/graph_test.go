package anonbll

import (
	"fmt"
	"github.com/dargzero/anonymization/anonmodel"
	"testing"
)

func TestGraphAnonymizer_ValidateExtraOptions(t *testing.T) {

	g := simpleAnonymizer("numeric")
	field := firstField(g)

	t.Run("unknown sub-type", func(t *testing.T) {
		field.Opts.Type = "unknown"
		err := g.validateExtraOptions()
		if err == nil {
			t.Errorf("expected error, got none")
		}
	})

	t.Run("missing min", func(t *testing.T) {
		field.Opts.Type = "float"
		err := g.validateExtraOptions()
		if err == nil {
			t.Errorf("expected error, got none")
		}
	})

	t.Run("missing max", func(t *testing.T) {
		field.Opts.Type = "float"
		min := 0.5
		field.Opts.Min = &min
		err := g.validateExtraOptions()
		if err == nil {
			t.Errorf("expected error, got none")
		}
	})

	t.Run("no validation errors", func(t *testing.T) {
		field.Opts.Type = "float"
		min := 0.5
		max := 1.5
		field.Opts.Min = &min
		field.Opts.Max = &max
		err := g.validateExtraOptions()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func TestGraphAnonymizer_GetSchema(t *testing.T) {

	g := getTestAnonymizer()
	schema := g.getSchema()

	for i, col := range schema.Columns {
		actual := col.GetName()
		expected := fmt.Sprintf("Col%d", i)
		if expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		gen := col.GetGeneralizer()
		if gen == nil {
			t.Errorf("expected generalizer, got nil")
		}
	}

}

func simpleAnonymizer(fieldType string) *graphAnonymizer {
	return &graphAnonymizer{
		qidFields: []anonmodel.FieldAnonymizationInfo{
			{
				Name: "Col0",
				Mode: "qid",
				Type: fieldType,
				Opts: anonmodel.ExtraFieldOptions{},
			},
		},
	}
}

func firstField(anonymizer *graphAnonymizer) *anonmodel.FieldAnonymizationInfo {
	return &anonymizer.qidFields[0]
}

func getTestAnonymizer() *graphAnonymizer {
	min := 0.5
	max := 1.0
	return &graphAnonymizer{
		qidFields: []anonmodel.FieldAnonymizationInfo{
			{
				Name: "Col0",
				Mode: "qid",
				Type: "numeric",
				Opts: anonmodel.ExtraFieldOptions{
					Type: "float",
					Min:  &min,
					Max:  &max,
				},
			},
			{
				Name: "Col1",
				Mode: "qid",
				Type: "numeric",
				Opts: anonmodel.ExtraFieldOptions{
					Min: &min,
					Max: &max,
				},
			},
			{
				Name: "Col2",
				Mode: "qid",
				Type: "numeric",
				Opts: anonmodel.ExtraFieldOptions{
					Type: "int",
					Min:  &min,
					Max:  &max,
				},
			},
			{
				Name: "Col3",
				Mode: "qid",
				Type: "prefix",
				Opts: anonmodel.ExtraFieldOptions{},
			},
			{
				Name: "Col4",
				Mode: "qid",
				Type: "prefix",
				Opts: anonmodel.ExtraFieldOptions{
					Max: &max,
				},
			},
		},
	}
}
