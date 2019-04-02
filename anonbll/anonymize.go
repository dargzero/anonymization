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
	if len(fieldsToSuppress) > 0 {
		if err := anondb.SuppressFields(anonCollectionName, fieldsToSuppress); err != nil {
			return err
		}
	}

	quasiIdentifierFields := anonmodel.GetQuasiIdentifierFields(dataset.Fields)

	var algorithm anonymizerAlgorithm
	switch dataset.Settings.Algorithm {
	case "mondrian":
		algorithm = &mondrian{}
		break
	case "graph":
		algorithm = &graphAnonymizer{}
		break
	default:
		return fmt.Errorf("%v is not supported (must be one of 'mondrian', 'graph')", dataset.Settings.Algorithm)
	}

	algorithm.initialize(dataset, anonCollectionName, quasiIdentifierFields)
	if err := algorithm.anonymize(); err != nil {
		return err
	}

	return anondb.RenameAnonFields(anonCollectionName, quasiIdentifierFields)
}
