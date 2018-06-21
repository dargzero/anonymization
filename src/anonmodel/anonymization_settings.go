package anonmodel

import "fmt"

// AnonymizationSettings stores the settings about the dataset
type AnonymizationSettings struct {
	K         int    `json:"k" bson:"k"`
	Algorithm string `json:"algorithm" bson:"algorithm"`
	Mode      string `json:"mode" bson:"mode"`
}

func (settings *AnonymizationSettings) validate() error {
	if settings.K < 2 {
		return fmt.Errorf("The 'k' value should be at least 2, got: %v", settings.K)
	}

	if settings.Algorithm != "mondrian" {
		return fmt.Errorf("The only currently supported anonymization is 'mondrian', got '%v'", settings.Algorithm)
	}

	if settings.Mode != "single" && settings.Mode != "continuous" {
		return fmt.Errorf("Anonymization mode should be 'single' or 'continuous', got '%v'", settings.Mode)
	}

	return nil
}
