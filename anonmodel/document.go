package anonmodel

// Document represents a data object of any type uploaded by the client
type Document map[string]interface{}

// Documents represents an array of data objects of any type uploaded by the client
type Documents []Document

func (document Document) Validate() error {
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

	return nil
}

// Convert Documents into []interface{}
func (documents Documents) Convert(continuous bool, table map[string]TypeConversionFunc) []interface{} {
	result := make([]interface{}, len(documents))
	for ix, document := range documents {
		if continuous {
			document["__pending"] = true
		}
		for key, value := range document {
			if table[key] != nil {
				document[key], _ = table[key](value)
			}
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
