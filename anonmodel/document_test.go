package anonmodel

import (
	"errors"
	"testing"
)

func TestDocument_Validate(t *testing.T) {

	t.Run("valid document", func(t *testing.T) {
		assertValidity(t, Document{
			"test":  1,
			"other": "dummy",
		}, true)
	})

	t.Run("contains id key", func(t *testing.T) {
		assertValidity(t, Document{
			"_id":   1,
			"other": "dummy",
		}, false)
	})

	t.Run("contains invalid key", func(t *testing.T) {
		assertValidity(t, Document{
			"id":     1,
			"__test": "dummy",
		}, false)
	})

	t.Run("contains invalid string", func(t *testing.T) {
		assertValidity(t, Document{
			"id":    1,
			"te.$t": "dummy",
		}, false)
	})

}

func TestDocuments_Validate(t *testing.T) {

	t.Run("empty slice", func(t *testing.T) {
		docs := Documents{}
		err := docs.Validate()
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})

	t.Run("non empty slice", func(t *testing.T) {
		docs := Documents{
			Document{},
		}
		err := docs.Validate()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

}

func TestDocuments_Convert(t *testing.T) {

	t.Run("asd", func(t *testing.T) {
		docs := Documents{
			Document{
				"name": "test1",
				"age":  12,
			},
			Document{
				"name": "test2",
				"age":  24,
			},
		}
		noop := func(i interface{}) (interface{}, error) {
			return i, nil
		}
		clear := func(i interface{}) (interface{}, error) {
			return 0, nil
		}
		converted := docs.Convert(false, map[string]TypeConversionFunc{
			"name": noop,
			"age":  clear,
		})

		if len(converted) != 2 ||
			converted[0].(Document)["name"] != "test1" ||
			converted[0].(Document)["age"] != 0 ||
			converted[1].(Document)["name"] != "test2" ||
			converted[1].(Document)["age"] != 0 {
			t.Errorf("invalid result: %v", converted)
		}

	})

	t.Run("swallows conversion error", func(t *testing.T) {
		docs := Documents{Document{"field": "value"}}
		evil := func(i interface{}) (interface{}, error) {
			return nil, errors.New("evil error")
		}
		converted := docs.Convert(false, map[string]TypeConversionFunc{
			"field": evil,
		})

		if len(converted) != 1 ||
			converted[0].(Document)["field"] != nil {
			t.Errorf("expected nil, got %v", converted[0].(Document)["field"])
		}
	})

	t.Run("continuous flag", func(t *testing.T) {
		docs := Documents{Document{"field": "value"}}
		converted := docs.Convert(true, map[string]TypeConversionFunc{})
		if len(converted) != 1 ||
			!converted[0].(Document)["__pending"].(bool) {
			t.Errorf("missing pending flag")
		}
	})
}

func TestErrValidation_Error(t *testing.T) {
	var e ErrValidation = "test"
	s := e.Error()
	if s != "test" {
		t.Errorf("expected %v, got %v", "test", s)
	}
}

func assertValidity(t *testing.T, d Document, expected bool) {
	err := d.Validate()
	actual := err == nil
	if actual != expected {
		t.Errorf("expected %v, got %v, err: %v", expected, actual, err)
	}
}
