package anonmodel

import "fmt"

// AnonymizationSettings stores the settings about the dataset
type AnonymizationSettings struct {
	K         int    `json:"k" bson:"k"`
	Algorithm string `json:"algorithm" bson:"algorithm"`
	Mode      string `json:"mode" bson:"mode"`
}

func (settings *AnonymizationSettings) Validate() error {
	if settings.K < 2 {
		return fmt.Errorf("The 'k' value should be at least 2, got: %v", settings.K)
	}

	if settings.Algorithm != "mondrian" && settings.Algorithm != "graph" {
		return fmt.Errorf("Expected mondrian or graph. Algorithm '%v' not supported.", settings.Algorithm)
	}

	if settings.Mode != "single" && settings.Mode != "continuous" {
		return fmt.Errorf("Anonymization mode should be 'single' or 'continuous', got '%v'", settings.Mode)
	}

	return nil
}
