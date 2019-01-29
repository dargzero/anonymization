package anondb

import (
	"anonmodel"
	"log"

	"github.com/globalsign/mgo"
)

// CreateDataset creates a new dataset
func CreateDataset(dataset *anonmodel.Dataset) error {
	session := globalSession.Copy()
	defer session.Close()

	datasets := session.DB("anondb").C("datasets")
	err := datasets.Insert(dataset)
	if err != nil && mgo.IsDup(err) {
		return ErrDuplicate
	}

	return err
}

// GetDataset tries to find a dataset with a given name
func GetDataset(name string) (dataset anonmodel.Dataset, err error) {
	session := globalSession.Copy()
	defer session.Close()

	datasets := session.DB("anondb").C("datasets")
	if err = datasets.FindId(name).One(&dataset); err == mgo.ErrNotFound {
		err = ErrNotFound
	}
	return
}

// DropDataset deletes a dataset with a given name
func DropDataset(name string) error {
	session := globalSession.Copy()
	defer session.Close()

	datasets := session.DB("anondb").C("datasets")
	if err := datasets.RemoveId(name); err == mgo.ErrNotFound {
		return ErrNotFound
	} else if err != nil {
		return err
	}

	data := session.DB("anondb").C("data_" + name)
	if err := data.DropCollection(); err != nil {
		log.Printf("Error deleting data for dataset '%v': %v", name, err.Error())
	}

	anon := session.DB("anondb").C("anon_" + name)
	if err := anon.DropCollection(); err != nil {
		log.Printf("Error deleting anonymized data for dataset '%v': %v", name, err.Error())
	}

	tempAnon := session.DB("anondb").C("temp_anon_" + name)
	if err := tempAnon.DropCollection(); err != nil {
		log.Printf("Error deleting anonymized data for dataset '%v': %v", name, err.Error())
	}

	return nil
}

// ListDatasets lists all the datasets in the database
func ListDatasets() (datasetList []anonmodel.Dataset, err error) {
	session := globalSession.Copy()
	defer session.Close()

	datasets := session.DB("anondb").C("datasets")
	if err = datasets.Find(nil).All(&datasetList); err != nil {
		return
	}

	if datasetList == nil {
		datasetList = []anonmodel.Dataset{}
	}
	return
}
