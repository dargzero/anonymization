package swagger

import (
	"anondb"
	"log"
	"net/http"
)

func ping(w http.ResponseWriter, r *http.Request) {
	if err := anondb.Ping(); err != nil {
		log.Printf("The database responded with error: %v", err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
