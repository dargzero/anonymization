package anonbll

import (
	"anondb"
	"anonmodel"
)

// UploadDocuments validates the documents and inserts them into the database
func UploadDocuments(sessionID string, documents anonmodel.Documents, last bool) (insertSuccessful bool, finalizeSuccessful bool, err error) {
	err = documents.Validate()
	if err != nil {
		return
	}

	dataset, err := anondb.SetUploadSessionBusy(sessionID)
	if err != nil {
		return
	}

	continuous := dataset.Settings.Mode == "continuous"
	err = uploadDocuments(documents, &dataset, continuous, last)
	if err != nil {
		return
	}
	insertSuccessful = true

	if last {
		err = finalizeUpload(&dataset, continuous)
		if err == nil {
			finalizeSuccessful = true
		}
	}

	return
}

func uploadDocuments(documents anonmodel.Documents, dataset *anonmodel.Dataset, continuous bool, last bool) error {
	if !last {
		defer anondb.SetUploadSessionNotBusy(dataset.Name, dataset.UploadSessionData.SessionID)
	}

	return anondb.InsertDocuments(dataset.Name, documents, continuous)
}

func finalizeUpload(dataset *anonmodel.Dataset, continuous bool) error {
	defer anondb.FinishUploadSession(dataset.Name, dataset.UploadSessionData.SessionID)

	return anonymizeDataset(dataset, continuous)
}
