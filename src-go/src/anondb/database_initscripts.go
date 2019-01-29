package anondb

import "github.com/globalsign/mgo"

func createDatasetsCollection(db *mgo.Database) error {
	datasets := db.C("datasets")

	index := mgo.Index{
		Key:        []string{"uploadSessionData.sessionId"},
		Unique:     true,
		Background: false,
		Sparse:     true,
		Name:       "unique_uploadSessionId",
	}
	if err := datasets.EnsureIndex(index); err != nil {
		return err
	}

	index = mgo.Index{
		Key:        []string{"uploadSessionData.lastModified"},
		Unique:     false,
		Background: false,
		Sparse:     true,
		Name:       "uploadSessionLastModified",
	}
	return datasets.EnsureIndex(index)
}
