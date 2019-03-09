/*
 * Data Anonymization Server
 *
 * This is a data anonymization server. You can set the anonymization requirements for the different datasets individually, and upload data to them. The uploaded data is anonymized on the server and can be then downloaded.
 *
 * API version: 0.1-alpha
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

import (
	"github.com/dargzero/anonymization/anonbll"
	"github.com/dargzero/anonymization/anondb"
	"github.com/dargzero/anonymization/anonmodel"
	"net/http"

	"github.com/gorilla/mux"
)

func uploadPost(w http.ResponseWriter, r *http.Request) {
	var request CreateUploadSessionRequest
	if !tryReadRequestBody(r, &request, w) {
		return
	}

	if request.DatasetName == "" {
		respondWithError(w, http.StatusBadRequest, "The value 'datasetName' must be set")
		return
	}

	if sessionID, err := anondb.CreateUploadSession(request.DatasetName); err != nil {
		handleDBNotFound(
			err,
			w,
			http.StatusBadRequest,
			"The dataset with the specified name was not found or an upload session for the specified dataset is currently in use")
	} else {
		respondWithJSON(w, http.StatusOK, CreateUploadSessionResponse{SessionID: sessionID})
	}
}

func uploadSessionIDPost(w http.ResponseWriter, r *http.Request) {
	last, err := readLastQueryParam(r.URL)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	var documents anonmodel.Documents
	if !tryReadRequestBody(r, &documents, w) {
		return
	}

	vars := mux.Vars(r)
	insertSuccessful, finalizeSuccessful, err := anonbll.UploadDocuments(vars["sessionId"], documents, last)
	if !insertSuccessful {
		switch err.(type) {
		case anonmodel.ErrValidation:
			respondWithError(w, http.StatusBadRequest, err.Error())
		default:
			handleDBNotFound(err, w, http.StatusBadRequest, "The upload session with the specified ID was not found or is currently in use")
		}
	} else {
		response := UploadResponse{
			InsertSuccessful:   insertSuccessful,
			FinalizeSuccessful: finalizeSuccessful,
		}
		if err != nil {
			response.Error = err.Error()
		}
		respondWithJSON(w, http.StatusOK, &response)
	}
}