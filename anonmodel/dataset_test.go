package anonmodel

import "testing"

func TestDataset_Validate(t *testing.T) {

	t.Run("invalid settings", func(t *testing.T) {
		d := newDataSet()
		d.Settings = newInvalidSettings()
		err := d.Validate()
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})

	t.Run("invalid field", func(t *testing.T) {
		d := newDataSet()
		d.Fields = append(d.Fields, newValidField())
		d.Fields = append(d.Fields, newInvalidField())
		err := d.Validate()
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})

	t.Run("valid settings", func(t *testing.T) {
		d := newDataSet()
		err := d.Validate()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func newInvalidField() FieldAnonymizationInfo {
	return FieldAnonymizationInfo{
		Name: "_id",
	}
}

func newValidField() FieldAnonymizationInfo {
	return FieldAnonymizationInfo{
		Name: "test field",
		Mode: "drop",
	}
}

func newInvalidSettings() AnonymizationSettings {
	return AnonymizationSettings{K: -1}
}

func newDataSet() *Dataset {
	d := &Dataset{
		Settings: AnonymizationSettings{
			K:         3,
			Algorithm: "mondrian",
			Mode:      "single",
		},
		Fields: []FieldAnonymizationInfo{},
	}
	return d
}
