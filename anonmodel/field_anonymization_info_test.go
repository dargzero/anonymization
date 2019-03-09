package anonmodel

import "testing"

func TestFieldAnonymizationInfo_Validate_Name(t *testing.T) {
	tests := []struct {
		name  string
		valid bool
	}{
		{"_id", false},
		{"__test", false},
		{"te.$t", false},
		{"test", true},
		{"field", true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := newFieldAnonymizationInfo()
			f.Name = test.name
			assertFieldValidity(t, f, test.valid)
		})
	}
}

func TestFieldAnonymizationInfo_Validate_Mode(t *testing.T) {
	tests := []struct {
		mode  string
		valid bool
	}{
		{"id", true},
		{"qid", true},
		{"keep", true},
		{"drop", true},
		{"test", false},
		{"other", false},
	}
	for _, test := range tests {
		t.Run(test.mode, func(t *testing.T) {
			f := newFieldAnonymizationInfo()
			f.Mode = test.mode
			assertFieldValidity(t, f, test.valid)
		})
	}
}

func TestFieldAnonymizationInfo_Validate_Type(t *testing.T) {

	t.Run("non id field, no restriction", func(t *testing.T) {
		f := newFieldAnonymizationInfo()
		f.Mode = "keep"
		f.Type = "any"
		assertFieldValidity(t, f, true)
	})

	t.Run("qid numeric", func(t *testing.T) {
		f := newFieldAnonymizationInfo()
		f.Mode = "qid"
		f.Type = "numeric"
		assertFieldValidity(t, f, true)
	})

	t.Run("qid prefix", func(t *testing.T) {
		f := newFieldAnonymizationInfo()
		f.Mode = "qid"
		f.Type = "prefix"
		assertFieldValidity(t, f, true)
	})

	t.Run("qid coords", func(t *testing.T) {
		f := newFieldAnonymizationInfo()
		f.Mode = "qid"
		f.Type = "coords"
		assertFieldValidity(t, f, true)
	})

	t.Run("qid other", func(t *testing.T) {
		f := newFieldAnonymizationInfo()
		f.Mode = "qid"
		f.Type = "other"
		assertFieldValidity(t, f, false)
	})

}

func TestGetQuasiIdentifierFields(t *testing.T) {
	fields := newFieldAnonymizationSlice()
	fields[0].Mode = "id"
	fields[1].Mode = "qid"
	fields[2].Mode = "keep"
	fields[3].Mode = "drop"
	fields[4].Mode = "qid"

	qids := GetQuasiIdentifierFields(fields)

	if len(qids) != 2 &&
		qids[0] != fields[1] &&
		qids[1] != fields[4] {
		t.Errorf("invalid result: %v", qids)
	}
}

func TestGetSuppressedFields(t *testing.T) {
	fields := newFieldAnonymizationSlice()
	fields[0].Mode = "id"
	fields[0].Name = "test1"
	fields[1].Mode = "qid"
	fields[1].Name = "test2"
	fields[2].Mode = "keep"
	fields[2].Name = "test3"
	fields[3].Mode = "drop"
	fields[3].Name = "test4"
	fields[4].Mode = "drop"
	fields[4].Name = "test5"

	suppressed := GetSuppressedFields(fields)

	if len(suppressed) != 3 &&
		suppressed[0] != "test1" &&
		suppressed[1] != "test4" &&
		suppressed[2] != "test5" {
		t.Errorf("invalid result: %v", suppressed)
	}
}

func assertFieldValidity(t *testing.T, f FieldAnonymizationInfo, expected bool) {
	err := f.Validate()
	actual := err == nil
	if actual != expected {
		t.Errorf("expected: %v, got %v, err: %v", expected, actual, err)
	}
}

func newFieldAnonymizationSlice() []FieldAnonymizationInfo {
	fields := []FieldAnonymizationInfo{
		newFieldAnonymizationInfo(),
		newFieldAnonymizationInfo(),
		newFieldAnonymizationInfo(),
		newFieldAnonymizationInfo(),
		newFieldAnonymizationInfo(),
	}
	return fields
}

func newFieldAnonymizationInfo() FieldAnonymizationInfo {
	return FieldAnonymizationInfo{
		Name: "field",
		Mode: "drop",
		Type: "numeric",
	}
}
