package anonbll

import (
	"github.com/dargzero/anonymization/anonmodel"
	"testing"
)

func TestGraphAnonymizer_GetSchema(t *testing.T) {

	t.Run("unknown field type", func(t *testing.T) {
		a := newAnonymizer(newField("col", "qid", "unknown"))
		_, err := a.getSchema()
		if err == nil {
			t.Errorf("expected error, got none")
		}
	})

	t.Run("numeric column type", func(t *testing.T) {

		col1 := newField("col1", "qid", "numeric")
		a := newAnonymizer(col1)

		t.Run("unknown numeric type", func(t *testing.T) {
			addOption(col1, "type", "unknown")
			_, err := a.getSchema()
			if err == nil {
				t.Errorf("expected error, got none")
			}
		})

		tests := []string{"int", "float"}

		for _, test := range tests {

			col1.Opts = make(map[string]string, 0) // clear options

			t.Run("column type "+test, func(t *testing.T) {

				addOption(col1, "type", test)

				t.Run("missing min option", func(t *testing.T) {
					_, err := a.getSchema()
					if err == nil {
						t.Errorf("expected error, got none")
					}
				})

				t.Run("missing max option", func(t *testing.T) {
					addOption(col1, "min", "10")
					_, err := a.getSchema()
					if err == nil {
						t.Errorf("expected error, got none")
					}
				})

				t.Run("cannot parse min", func(t *testing.T) {
					addOption(col1, "min", "X")
					_, err := a.getSchema()
					if err == nil {
						t.Errorf("expected error, got none")
					}
				})

				t.Run("cannot parse max", func(t *testing.T) {
					addOption(col1, "min", "10")
					addOption(col1, "max", "X")
					_, err := a.getSchema()
					if err == nil {
						t.Errorf("expected error, got none")
					}
				})

				t.Run("proper int column", func(t *testing.T) {
					addOption(col1, "min", "10")
					addOption(col1, "max", "20")
					_, err := a.getSchema()
					if err != nil {
						t.Errorf("unexpected error: %v", err)
					}
				})

			})
		}

		t.Run("default column type", func(t *testing.T) {
			addOption(col1, "min", "0.5")
			addOption(col1, "max", "1.0")
			schema, _ := a.getSchema()
			col1 := schema.Columns[0]
			if col1.GetName() != "col1" {
				t.Errorf("expected %v, got %v", "col1", col1.GetName())
			}
			if col1.GetGeneralizer() == nil {
				t.Errorf("expected generalizer instance, got nil")
			}
		})
	})

}

func newAnonymizer(fields ...*anonmodel.FieldAnonymizationInfo) *graphAnonymizer {
	a := &graphAnonymizer{}
	for _, field := range fields {
		a.qidFields = append(a.qidFields, field)
	}
	return a
}

func newField(name, mode, ft string) *anonmodel.FieldAnonymizationInfo {
	return &anonmodel.FieldAnonymizationInfo{
		Name: name,
		Mode: mode,
		Type: ft,
		Opts: make(map[string]string, 0),
	}
}

func addOption(field *anonmodel.FieldAnonymizationInfo, key, value string) {
	field.Opts[key] = value
}
