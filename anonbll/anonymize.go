package anonbll

import (
	"fmt"
	"github.com/dargzero/anonymization/anondb"
	"github.com/dargzero/anonymization/anonmodel"
	"log"
	"time"
)

type anonymizerAlgorithm interface {
	initialize(*anonmodel.Dataset, string, []anonmodel.FieldAnonymizationInfo)
	anonymize() error
}

func anonymizeDataset(dataset *anonmodel.Dataset, continuous bool) error {
	start := time.Now()
	defer func() { log.Printf("Anonymization took %v", time.Since(start)) }()

	if err := doAnonymization(dataset, continuous); err != nil {
		return err
	}

	if !continuous {
		return nil
	}
	return anondb.MoveTempAnonymizedData(dataset.Name)
}

func doAnonymization(dataset *anonmodel.Dataset, continuous bool) error {
	var anonCollectionName string
	if continuous {
		anonCollectionName = "temp_anon_" + dataset.Name
	} else {
		anonCollectionName = "anon_" + dataset.Name
	}

	if err := anondb.CopyData(dataset.Name, continuous, anonCollectionName); err != nil {
		return err
	}

	fieldsToSuppress := anonmodel.GetSuppressedFields(dataset.Fields)
	if err := anondb.SuppressFields(anonCollectionName, fieldsToSuppress); err != nil {
		return err
	}

	var algorithm anonymizerAlgorithm
	quasiIdentifierFields := anonmodel.GetQuasiIdentifierFields(dataset.Fields)

	if dataset.Settings.Algorithm == "mondrian" {
		algorithm = &mondrian{}
	} else {
		return fmt.Errorf("The only currently supported anonymization algorithm is 'mondrian', got '%v'", dataset.Settings.Algorithm)
	}

	algorithm.initialize(dataset, anonCollectionName, quasiIdentifierFields)
	if err := algorithm.anonymize(); err != nil {
		return err
	}

	return anondb.RenameAnonFields(anonCollectionName, quasiIdentifierFields)
}
