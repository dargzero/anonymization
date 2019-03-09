package swagger

import (
	"encoding/json"
	"github.com/dargzero/anonymization/anondb"
	"log"
	"net/http"
)

type errorResponse struct {
	Error string `json:"error,omitempty"`
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, errorResponse{Error: message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)
	w.Write(response)
}

func tryReadRequestBody(r *http.Request, payload interface{}, w http.ResponseWriter) bool {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(payload); err != nil {
		log.Printf("Invalid request payload: %v", err.Error())
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")

		return false
	}
	defer r.Body.Close()

	return true
}

func handleDBNotFound(err error, w http.ResponseWriter, code int, message string) {
	logDBError(err)
	if err == anondb.ErrNotFound {
		respondWithError(w, code, message)
	} else {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
}

func handleDBDuplicate(err error, w http.ResponseWriter, code int, message string) {
	logDBError(err)
	if err == anondb.ErrDuplicate {
		respondWithError(w, code, message)
	} else {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
}

func logDBError(err error) {
	log.Printf("The database responded with error: %v", err.Error())
}
