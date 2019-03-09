package anonmodel

import (
	"fmt"
	"testing"
)

func TestAnonymizationSettings_Validate_K(t *testing.T) {
	tests := []struct {
		k     int
		valid bool
	}{
		{k: -2, valid: false},
		{k: 0, valid: false},
		{k: 1, valid: false},
		{k: 2, valid: true},
		{k: 5, valid: true},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("K = %d", test.k), func(t *testing.T) {
			settings := newAnonymizationSettings()
			settings.K = test.k
			assertValidity(t, settings, test.valid)
		})
	}
}

func TestAnonymizationSettings_Validate_Algorithm(t *testing.T) {
	tests := []struct {
		algorithm string
		valid     bool
	}{
		{algorithm: "mondrian", valid: true},
		{algorithm: "newAnonymizationSettingsWith", valid: false},
		{algorithm: "other", valid: false},
	}
	for _, test := range tests {
		t.Run(test.algorithm, func(t *testing.T) {
			settings := newAnonymizationSettings()
			settings.Algorithm = test.algorithm
			assertValidity(t, settings, test.valid)
		})
	}
}

func TestAnonymizationSettings_Validate_Mode(t *testing.T) {
	tests := []struct {
		mode  string
		valid bool
	}{
		{mode: "single", valid: true},
		{mode: "continuous", valid: true},
		{mode: "test", valid: false},
		{mode: "dummy", valid: false},
	}
	for _, test := range tests {
		t.Run(test.mode, func(t *testing.T) {
			settings := newAnonymizationSettings()
			settings.Mode = test.mode
			assertValidity(t, settings, test.valid)
		})
	}
}

func newAnonymizationSettings() AnonymizationSettings {
	return AnonymizationSettings{
		K:         5,
		Algorithm: "mondrian",
		Mode:      "single",
	}
}

func assertValidity(t *testing.T, s AnonymizationSettings, expected bool) {
	err := s.Validate()
	actual := err == nil
	if actual != expected {
		t.Errorf("expected: %v, got %v - err: %v", expected, actual, err)
	}
}
