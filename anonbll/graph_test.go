package anonbll

import (
	"bitbucket.org/dargzero/k-anon/generalization"
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

	t.Run("prefix column type", func(t *testing.T) {

		col1 := newField("col1", "qid", "prefix")
		a := newAnonymizer(col1)

		t.Run("prefix column with defaults", func(t *testing.T) {
			schema, err := a.getSchema()
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			g := schema.Columns[0].GetGeneralizer().(*generalization.PrefixGeneralizer)
			if g.MaxWords != 100 {
				t.Errorf("expected %v, got %v", 100, g.MaxWords)
			}
		})

		t.Run("prefix column with max words", func(t *testing.T) {
			addOption(col1, "max", "1000")
			schema, err := a.getSchema()
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			g := schema.Columns[0].GetGeneralizer().(*generalization.PrefixGeneralizer)
			if g.MaxWords != 1000 {
				t.Errorf("expected %v, got %v", 1000, g.MaxWords)
			}
		})

		t.Run("prefix column with invalid max option", func(t *testing.T) {
			addOption(col1, "max", "invalid")
			_, err := a.getSchema()
			if err == nil {
				t.Errorf("expected error, got none")
			}
		})

	})

	t.Run("numeric column type", func(t *testing.T) {

		t.Run("unknown numeric type", func(t *testing.T) {
			col1 := newField("col1", "qid", "numeric")
			addOption(col1, "type", "unknown")
			a := newAnonymizer(col1)
			_, err := a.getSchema()
			if err == nil {
				t.Errorf("expected error, got none")
			}
		})

		tests := []string{"int", "float"}

		for _, test := range tests {

			t.Run("column type "+test, func(t *testing.T) {

				col1 := newField("col1", "qid", "numeric")
				addOption(col1, "type", test)
				a := newAnonymizer(col1)

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
			col1 := newField("col1", "qid", "numeric")
			addOption(col1, "min", "0.5")
			addOption(col1, "max", "1.0")
			a := newAnonymizer(col1)
			schema, _ := a.getSchema()
			actual := schema.Columns[0].GetName()
			if actual != "col1" {
				t.Errorf("expected %v, got %v", "col1", actual)
			}
			if schema.Columns[0].GetGeneralizer() == nil {
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
