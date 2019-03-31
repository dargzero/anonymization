package anonmodel

import (
	"fmt"
	"strings"
)

// FieldAnonymizationInfo stores how each data field should be handled during anonymization
type FieldAnonymizationInfo struct {
	Name string            `json:"name" bson:"name"`
	Mode string            `json:"mode" bson:"mode"`
	Type string            `json:"type" bson:"type"`
	Opts map[string]string `json:"opts" bson:"opts"`
}

func (fieldInfo *FieldAnonymizationInfo) Validate() error {
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

func validateFieldName(field string) error {
	if field == "_id" {
		return ErrValidation("Validation error: the '_id' field is not allowed")
	}

	if strings.HasPrefix(field, "__") {
		return ErrValidation(fmt.Sprintf("Validation error (%v): document fields starting with '__' are reserved by the anonymization server", field))
	}

	if strings.ContainsAny(field, ".$") {
		return ErrValidation(fmt.Sprintf("Validation error (%v): document fields containing either '.' or '$' are not allowed", field))
	}

	return nil
}
