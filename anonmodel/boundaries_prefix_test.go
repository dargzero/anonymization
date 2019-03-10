package anonmodel

import (
	"reflect"
	"testing"
)

func TestPrefixBoundary_Clone(t *testing.T) {
	b := newPrefixBoundary("tests are first class citizens")
	clone := b.Clone().(*PrefixBoundary)
	if !prefixEquals(clone, b) {
		t.Errorf("invalid clone: %v, original: %v", clone, b)
	}
}

func TestPrefixBoundary_AddFilter(t *testing.T) {
	t.Run("add simple filter", func(t *testing.T) {
		b := newPrefixBoundary("test 123")
		b.AddFilter("new")
		_, contains := b.Filters["new"]
		if !contains {
			t.Errorf("filter not found")
		}
	})

	t.Run("add existing filter", func(t *testing.T) {
		b := newPrefixBoundary("test 123")
		b.AddFilter("cats are wild")
		_, contains := b.Filters["cats are wild"]
		if !contains {
			t.Errorf("filter not found")
		}
	})

	t.Run("overwrite filter", func(t *testing.T) {
		b := newPrefixBoundary("test 123")
		b.AddFilter("cats")
		_, kept := b.Filters["cats are wild"]
		_, replaced := b.Filters["cats"]
		if kept || !replaced {
			t.Errorf("invalid filters: %v", b.Filters)
		}
	})
}

func TestPrefixBoundary_SetPrefix(t *testing.T) {
	prefix := "cats"
	b := newPrefixBoundary("")
	b.SetPrefix(prefix)
	if b.Prefix != prefix {
		t.Errorf("prefix mismatch: %v", b)
	}
	_, cat := b.Filters["cats are wild"]
	_, dog := b.Filters["dogs are playful"]
	if !cat || dog {
		t.Errorf("incorrect filters: %v", b)
	}
}

func TestPrefixBoundary_GetGeneralizedValue(t *testing.T) {

	t.Run("non empty prefix", func(t *testing.T) {
		prefix := "cats are agile"
		b := newPrefixBoundary(prefix)
		actual := b.GetGeneralizedValue()
		expected := prefix
		if expected != actual {
			t.Errorf("expected %s, got %s", expected, actual)
		}
	})

	t.Run("empty prefix", func(t *testing.T) {
		b := newPrefixBoundary("")
		actual := b.GetGeneralizedValue()
		expected := "-"
		if expected != actual {
			t.Errorf("expected %s, got %s", expected, actual)
		}
	})
}

func prefixEquals(p1, p2 *PrefixBoundary) bool {
	return p1.Prefix == p2.Prefix &&
		reflect.DeepEqual(p1.Filters, p2.Filters)
}

func newPrefixBoundary(prefix string) *PrefixBoundary {
	return &PrefixBoundary{
		Prefix: prefix,
		Filters: map[string]struct{}{
			"cats are wild":    {},
			"dogs are playful": {},
		},
	}
}
