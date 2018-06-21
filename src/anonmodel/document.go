package anonmodel

import (
	"fmt"
	"strings"
)

// Document represents a data object of any type uploaded by the client
type Document map[string]interface{}

// Documents represents an array of data objects of any type uploaded by the client
type Documents []Document

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

func (document Document) validate() error {
	for key := range document {
		if err := validateFieldName(key); err != nil {
			return err
		}
	}

	return nil
}

// Validate validates the set of documents
func (documents Documents) Validate() error {
	if len(documents) == 0 {
		return ErrValidation("No documents sent to upload")
	}

	for _, doc := range documents {
		if err := doc.validate(); err != nil {
			return err
		}
	}

	return nil
}

// Convert convert the array of Documents into an array if interface{}s
func (documents Documents) Convert(continuous bool) []interface{} {
	result := make([]interface{}, len(documents))
	for ix, document := range documents {
		if continuous {
			document["__pending"] = true
		}
		result[ix] = document
	}

	return result
}

// ErrValidation signals that some of the documents had problems with them
type ErrValidation string

func (err ErrValidation) Error() string {
	return string(err)
}
