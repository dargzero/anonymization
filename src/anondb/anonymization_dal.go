package anondb

import (
	"anonmodel"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// CopyData copies the data so that it can be anonymized
func CopyData(datasetName string, continuous bool, anonCollectionName string) error {
	session := globalSession.Copy()
	defer session.Close()
	data := session.DB("anondb").C("data_" + datasetName)

	var pipeline []bson.M
	if continuous {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"__pending": true}})
		pipeline = append(pipeline, bson.M{"$project": bson.M{"__pending": 0}})
	}
	pipeline = append(pipeline, bson.M{"$addFields": bson.M{"__anonymized": false}})
	pipeline = append(pipeline, bson.M{"$out": anonCollectionName})

	if err := data.Pipe(pipeline).Iter().Close(); err != nil {
		return err
	}

	anon := session.DB("anondb").C(anonCollectionName)
	index := mgo.Index{
		Key:        []string{"__anonymized"},
		Unique:     false,
		Background: false,
		Sparse:     false,
		Name:       "isAnonymized",
	}
	return anon.EnsureIndex(index)
}

// SuppressFields suppresses the specified fields in the anonymized data
func SuppressFields(anonCollectionName string, fields []string) error {
	session := globalSession.Copy()
	defer session.Close()
	anon := session.DB("anondb").C(anonCollectionName)

	unset := bson.M{}
	for _, field := range fields {
		unset[field] = ""
	}

	_, err := anon.UpdateAll(bson.M{}, bson.M{"$unset": unset})
	return err
}

// CreateIndices creates indices of the quasi identifiers for the anon database
func CreateIndices(anonCollectionName string, fields []anonmodel.FieldAnonymizationInfo) error {
	session := globalSession.Copy()
	defer session.Close()
	anon := session.DB("anondb").C(anonCollectionName)

	for _, field := range fields {
		sparse := true
		if field.Type == "numeric" {
			sparse = false
		}

		index := mgo.Index{
			Key:        []string{field.Name},
			Unique:     false,
			Background: false,
			Sparse:     sparse,
			Name:       "anonIndex_" + field.Name,
		}

		if err := anon.EnsureIndex(index); err != nil {
			return err
		}
	}

	return nil
}

// DropIndices drops the indices of the quasi identifiers
func DropIndices(anonCollectionName string, fields []anonmodel.FieldAnonymizationInfo) {
	session := globalSession.Copy()
	defer session.Close()
	anon := session.DB("anondb").C(anonCollectionName)

	for _, field := range fields {
		anon.DropIndexName("anonIndex_" + field.Name)
	}
}

// RenameAnonFields moves the anonymized fields to their original places
func RenameAnonFields(anonCollectionName string, fields []anonmodel.FieldAnonymizationInfo) error {
	session := globalSession.Copy()
	defer session.Close()
	anon := session.DB("anondb").C(anonCollectionName)

	replace := bson.M{}
	for _, field := range fields {
		replace["__anon_"+field.Name] = field.Name
	}

	_, err := anon.UpdateAll(bson.M{}, bson.M{"$rename": replace})
	return err
}

// MoveTempAnonymizedData moves the data from the temporary anonimized collection to the main one
func MoveTempAnonymizedData(datasetName string) error {
	session := globalSession.Copy()
	defer session.Close()

	tempAnon := session.DB("anondb").C("temp_anon_" + datasetName)
	anon := session.DB("anondb").C("anon_" + datasetName)
	data := session.DB("anondb").C("data_" + datasetName)

	if err := moveTempAnonymizedData(tempAnon, anon, data); err != nil {
		return err
	}

	return tempAnon.DropCollection()
}

const moveBatchSize int = 1000

func moveTempAnonymizedData(tempAnon *mgo.Collection, anon *mgo.Collection, data *mgo.Collection) error {
	iter := tempAnon.Find(bson.M{}).Batch(moveBatchSize).Iter()
	defer iter.Close()

	for {
		anonBulk := anon.Bulk()
		anonBulk.Unordered()

		dataBulk := data.Bulk()
		dataBulk.Unordered()

		i := 0
		for ; i < moveBatchSize && !iter.Done(); i++ {
			var document bson.M
			if success := iter.Next(&document); !success {
				if err := iter.Err(); err != nil {
					return err
				}

				// this means we have no error, there are no more documents in the iteration
				// to be honest, this should not happen, as we are checking iter.Done, however, it does
				break
			}

			anonBulk.Insert(document)
			dataBulk.Update(bson.M{"_id": document["_id"]}, bson.M{"$set": bson.M{"__pending": false}})
		}

		if i > 0 {
			if _, err := anonBulk.Run(); err != nil {
				return err
			}
			if _, err := dataBulk.Run(); err != nil {
				return err
			}
		}

		// premature batch termination signals no more documents available
		if i < moveBatchSize {
			return nil
		}
	}
}
