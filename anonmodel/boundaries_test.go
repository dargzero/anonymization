package anonmodel

import (
	"fmt"
	"testing"
)

func TestPartition_Clone(t *testing.T) {
	p := &Partition{
		"test1": &MockBoundary{value: "test1"},
		"test2": &MockBoundary{value: "test2"},
	}

	cloned := p.Clone()

	if cloned["test1"].GetGeneralizedValue() != "cloned test1" ||
		cloned["test2"].GetGeneralizedValue() != "cloned test2" {
		t.Errorf("invalid clone: %v", cloned)
	}
}

func TestNumericBoundary_Clone(t *testing.T) {
	lb := 1.456
	ub := 5.678
	bound := &NumericBoundary{
		LowerBound:          &lb,
		LowerBoundInclusive: true,
		UpperBound:          &ub,
		UpperBoundInclusive: true,
	}

	clone := bound.Clone().(*NumericBoundary)

	if *clone.LowerBound != lb || *clone.UpperBound != ub ||
		clone.UpperBoundInclusive != bound.UpperBoundInclusive ||
		clone.LowerBoundInclusive != bound.LowerBoundInclusive {
		t.Errorf("invalid clone: %v", clone)
	}
}

func TestNumericBoundary_GetGeneralizedValue(t *testing.T) {
	tests := []struct {
		lb, ub   float64
		li, ui   bool
		expected string
	}{
		{3.5, 6.7, true, true, "[3.5, 6.7]"},
		{3.5, 6.7, false, true, "]3.5, 6.7]"},
		{3.5, 6.7, true, false, "[3.5, 6.7["},
		{3.5, 3.5, true, true, "3.5"},
	}
	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			bound := &NumericBoundary{
				LowerBound:          &test.lb,
				LowerBoundInclusive: test.li,
				UpperBound:          &test.ub,
				UpperBoundInclusive: test.ui,
			}

			actual := bound.GetGeneralizedValue()

			if test.expected != actual {
				t.Errorf("expected %v, got %v", test.expected, actual)
			}
		})
	}

	t.Run("nil bounds", func(t *testing.T) {
		bound := &NumericBoundary{
			LowerBound:          nil,
			LowerBoundInclusive: true,
			UpperBound:          nil,
			UpperBoundInclusive: false,
		}
		actual := bound.GetGeneralizedValue()

		expected := "]-inf, inf["
		if expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
	})
}

type MockBoundary struct {
	value string
}

func (b *MockBoundary) Clone() Boundary {
	return &MockBoundary{
		value: fmt.Sprintf("cloned %s", b.value),
	}
}

func (b *MockBoundary) GetGeneralizedValue() string {
	return b.value
}
