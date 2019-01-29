package anonmodel

import "fmt"

// FieldAnonymizationInfo stores how each data field should be handled during anonymization
type FieldAnonymizationInfo struct {
	Name string `json:"name" bson:"name"`
	Mode string `json:"mode" bson:"mode"`
	Type string `json:"type" bson:"type"`
}

func (fieldInfo *FieldAnonymizationInfo) validate() error {
	if err := validateFieldName(fieldInfo.Name); err != nil {
		return err
	}

	if fieldInfo.Mode != "id" && fieldInfo.Mode != "qid" && fieldInfo.Mode != "keep" && fieldInfo.Mode != "drop" {
		return fmt.Errorf("Field 'mode' should be one of 'id', 'qid', 'keep' or 'drop', got '%v'", fieldInfo.Mode)
	}

	if fieldInfo.Mode == "qid" && fieldInfo.Type != "numeric" && fieldInfo.Type != "prefix" && fieldInfo.Type != "coords" {
		return fmt.Errorf("Field 'type' should be one of 'numeric' or 'prefix' or 'coords', got '%v'", fieldInfo.Type)
	}

	return nil
}

// GetSuppressedFields gets which fields should be suppressed in the anonymized data
func GetSuppressedFields(fieldInfos []FieldAnonymizationInfo) []string {
	var result []string

	for _, fieldInfo := range fieldInfos {
		if fieldInfo.Mode == "id" || fieldInfo.Mode == "drop" {
			result = append(result, fieldInfo.Name)
		}
	}

	return result
}

// GetQuasiIdentifierFields gets the fields that are specified as quasi identifiers
func GetQuasiIdentifierFields(fieldInfos []FieldAnonymizationInfo) []FieldAnonymizationInfo {
	var result []FieldAnonymizationInfo

	for _, fieldInfo := range fieldInfos {
		if fieldInfo.Mode == "qid" {
			result = append(result, fieldInfo)
		}
	}

	return result
}
