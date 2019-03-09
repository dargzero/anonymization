package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type errorResponse struct {
	Error string `json:"error,omitempty"`
}

func main() {
	os.Exit(healthCheck())
}

func healthCheck() int {
	var netClient = &http.Client{
		Timeout: time.Second * 40,
	}

	response, err := netClient.Get("http://localhost:9137/v1/ping")
	if err != nil {
		log.Printf("Error while accessing webserver: %v", err.Error())
		return 1
	}
	defer response.Body.Close()

	if response.StatusCode == 200 {
		log.Printf("Webserver is running OK")
		return 0
	}

	errorMessage := getErrorMessage(response.Body)
	log.Printf("Webserver returned error: %v (%v)", errorMessage, response.Status)
	return 1
}

func getErrorMessage(body io.ReadCloser) string {
	var response errorResponse

	decoder := json.NewDecoder(body)
	if err := decoder.Decode(&response); err != nil {
		return "Unknown error"
	}

	return response.Error
}
