package anondb

import (
	"anonmodel"
	"log"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// CreateUploadSession creates an upload session for the given dataset
func CreateUploadSession(datasetName string) (string, error) {
	session := globalSession.Copy()
	defer session.Close()

	datasets := session.DB("anondb").C("datasets")

	query := bson.M{
		"_id":               datasetName,
		"uploadSessionData": bson.M{"$exists": false},
		"$or": []bson.M{
			bson.M{"settings.mode": "continuous"},
			bson.M{"anonymized": false},
		},
	}

	// in the VERY unlikely case when there is a UUID collision, try again
	for {
		sessionID := generateUUID()
		if err := createUploadSession(datasets, query, sessionID); !mgo.IsDup(err) {
			return sessionID, err
		}
	}
}

func createUploadSession(datasets *mgo.Collection, query bson.M, sessionID string) error {
	update := bson.M{
		"$set": bson.M{
			"uploadSessionData": anonmodel.UploadSessionData{
				SessionID:    sessionID,
				Busy:         false,
				LastModified: time.Now().UTC(),
			},
		},
	}

	err := datasets.Update(query, update)
	if err == mgo.ErrNotFound {
		return ErrNotFound
	}

	return err
}

// SetUploadSessionBusy sets an upload session busy if it is available
func SetUploadSessionBusy(sessionID string) (dataset anonmodel.Dataset, err error) {
	session := globalSession.Copy()
	defer session.Close()

	datasets := session.DB("anondb").C("datasets")
	query := bson.M{
		"uploadSessionData.sessionId": sessionID,
		"uploadSessionData.busy":      false,
	}
	change := mgo.Change{
		Update: bson.M{
			"$set": bson.M{
				"uploadSessionData.busy":         true,
				"uploadSessionData.lastModified": time.Now().UTC(),
			},
		},
		Upsert:    false,
		ReturnNew: true,
	}

	if _, err := datasets.Find(query).Apply(change, &dataset); err == mgo.ErrNotFound {
		return dataset, ErrNotFound
	}

	return
}

// SetUploadSessionNotBusy sets an upload session no longer busy
func SetUploadSessionNotBusy(datasetName string, sessionID string) {
	session := globalSession.Copy()
	defer session.Close()

	datasets := session.DB("anondb").C("datasets")
	err := datasets.Update(
		bson.M{"_id": datasetName, "uploadSessionData.sessionId": sessionID},
		bson.M{
			"$set": bson.M{
				"uploadSessionData.busy":         false,
				"uploadSessionData.lastModified": time.Now().UTC(),
			},
		})
	if err != nil {
		log.Printf("Error while setting upload session not busy for dataset '%v' (%v)", datasetName, sessionID)
	}
}

// MaintenanceSetUploadSessionBusy signals that an upload session is busy - used from maintenance
func MaintenanceSetUploadSessionBusy(datasetName string) error {
	session := globalSession.Copy()
	defer session.Close()

	datasets := session.DB("anondb").C("datasets")
	return datasets.UpdateId(datasetName, bson.M{"$set": bson.M{"uploadSessionData.busy": true}})
}

// FinishUploadSession signals that an upload session is done
func FinishUploadSession(datasetName string, sessionID string) error {
	session := globalSession.Copy()
	defer session.Close()

	datasets := session.DB("anondb").C("datasets")
	return datasets.Update(
		bson.M{"_id": datasetName, "uploadSessionData.sessionId": sessionID},
		bson.M{"$set": bson.M{"anonymized": true}, "$unset": bson.M{"uploadSessionData": ""}})
}

// ListOldUploadSessions lists all upload sessions that are older than a given age
func ListOldUploadSessions(minAge time.Duration) (datasetList []anonmodel.Dataset, err error) {
	session := globalSession.Copy()
	defer session.Close()

	datasets := session.DB("anondb").C("datasets")
	query := bson.M{
		"uploadSessionData.lastModified": bson.M{"$lt": time.Now().UTC().Add(-minAge)},
	}

	err = datasets.Find(query).All(&datasetList)
	return
}
