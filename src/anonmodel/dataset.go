package anonmodel

import (
	"time"
)

// UploadSessionData stores the information about the upload session in the dataset
type UploadSessionData struct {
	SessionID    string    `bson:"sessionId"`
	Busy         bool      `bson:"busy"`
	LastModified time.Time `bson:"lastModified"`
}

// Dataset represents a dataset in the database
type Dataset struct {
	Name              string                   `bson:"_id"`
	Settings          AnonymizationSettings    `bson:"settings"`
	Fields            []FieldAnonymizationInfo `bson:"fields"`
	UploadSessionData *UploadSessionData       `bson:"uploadSessionData,omitempty"`
	Anonymized        bool                     `bson:"anonymized"`
}

// Validate validates a dataset sent by the client
func (dataset *Dataset) Validate() error {
	if err := dataset.Settings.validate(); err != nil {
		return err
	}

	for _, field := range dataset.Fields {
		if err := field.validate(); err != nil {
			return err
		}
	}

	return nil
}
