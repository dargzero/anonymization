package anondb

import (
	"anonmodel"
	"log"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// InsertDocuments inserts the given documents into the database
func InsertDocuments(datasetName string, documents anonmodel.Documents, continuous bool) error {
	session := globalSession.Copy()
	defer session.Close()

	data := session.DB("anondb").C("data_" + datasetName)
	if err := ensureDataCollectionExists(data, continuous); err != nil {
		return err
	}

	table, err := MakeTypeConversionTable(datasetName)
	if err != nil {
		return err
	}
	log.Println("Converting documents")
	convertedDocs := documents.Convert(continuous, table)

	bulk := data.Bulk()
	bulk.Unordered()
	bulk.Insert(convertedDocs...)
	_, err = bulk.Run()
	return err
}

func ensureDataCollectionExists(data *mgo.Collection, continuous bool) error {
	if !continuous {
		return nil
	}

	index := mgo.Index{
		Key:        []string{"__pending"},
		Unique:     false,
		Background: false,
		Sparse:     false,
		Name:       "isPending",
	}
	return data.EnsureIndex(index)
}

// IsValidID tells whether an ID value is valid or not
func IsValidID(from string) bool {
	return bson.IsObjectIdHex(from)
}

// ListDocuments lists the stored documents for the specified dataset
func ListDocuments(datasetName string, size int, from string, documents *anonmodel.Documents) (string, error) {
	session := globalSession.Copy()
	defer session.Close()
	data := session.DB("anondb").C("data_" + datasetName)

	query := bson.M{}
	return listDocuments(size, from, query, data, documents)
}

// ListAnonDocuments lists the anonymized documents for the specified dataset
func ListAnonDocuments(datasetName string, size int, from string, documents *anonmodel.Documents) (string, error) {
	session := globalSession.Copy()
	defer session.Close()
	data := session.DB("anondb").C("anon_" + datasetName)

	query := bson.M{
		"__anonymized": true,
	}
	return listDocuments(size, from, query, data, documents)
}

func listDocuments(size int, from string, query bson.M, collection *mgo.Collection, documents *anonmodel.Documents) (string, error) {
	if from != "" {
		query["_id"] = bson.M{"$gt": bson.ObjectIdHex(from)}
	}

	if err := collection.Find(query).Sort("_id").Limit(size + 1).All(documents); err != nil {
		return "", err
	}

	if documents == nil || len(*documents) == 0 {
		return "", ErrNotFound
	}

	if len(*documents) < size+1 {
		return "", nil
	}

	*documents = (*documents)[:size]
	return (*documents)[size-1]["_id"].(bson.ObjectId).Hex(), nil
}

// GetDocument gets a document by its id from the specified dataset
func GetDocument(datasetName string, documentID string, document *anonmodel.Document) error {
	session := globalSession.Copy()
	defer session.Close()
	data := session.DB("anondb").C("data_" + datasetName)

	query := bson.M{}
	return getDocument(documentID, query, data, document)
}

// GetAnonDocument gets an anonymized document by its id from the specified dataset
func GetAnonDocument(datasetName string, documentID string, document *anonmodel.Document) error {
	session := globalSession.Copy()
	defer session.Close()
	data := session.DB("anondb").C("anon_" + datasetName)

	query := bson.M{
		"__anonymized": true,
	}
	return getDocument(documentID, query, data, document)
}

func getDocument(documentID string, query bson.M, collection *mgo.Collection, document *anonmodel.Document) error {
	query["_id"] = bson.ObjectIdHex(documentID)
	err := collection.Find(query).One(document)
	if err == mgo.ErrNotFound {
		return ErrNotFound
	}

	return err
}
