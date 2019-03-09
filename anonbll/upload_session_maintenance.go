package anonbll

import (
	"anondb"
	"anonmodel"
	"log"
	"sync/atomic"
	"time"
)

const minSessionAge time.Duration = time.Minute * 120
const maintenancePeriod time.Duration = time.Minute * 30

var stopChannel chan bool
var stop int32

var finishedChannel chan bool

// StartUploadSessionMaintenance starts the upload session maintenance background process
func StartUploadSessionMaintenance() {
	stopChannel = make(chan bool, 1)
	atomic.StoreInt32(&stop, 0)

	finishedChannel = make(chan bool, 1)

	go uploadSessionMaintenance()
}

func uploadSessionMaintenance() {
	defer func() { finishedChannel <- true }()
	log.Printf("Upload session maintenance background task started...")

	for {
		log.Printf("Running upload session maintenance...")
		if atomic.LoadInt32(&stop) > 0 {
			return
		}

		if datasetList, err := anondb.ListOldUploadSessions(minSessionAge); err != nil {
			log.Printf("Error while listing old upload sessions: %v", err.Error())
		} else {
			doMaintenance(datasetList)
		}

		log.Printf("Upload session maintenance finished")
		select {
		case <-stopChannel:
			return
		case <-time.After(maintenancePeriod):
		}
	}
}

func doMaintenance(datasetList []anonmodel.Dataset) {
	for _, dataset := range datasetList {
		log.Printf("Maintenance: session '%v' for dataset '%v'...", dataset.UploadSessionData.SessionID, dataset.Name)
		if atomic.LoadInt32(&stop) > 0 {
			return
		}

		if err := anondb.MaintenanceSetUploadSessionBusy(dataset.Name); err != nil {
			log.Printf("Error during upload session maintenance: %v", err.Error())
			continue
		}

		// TODO finalize documents

		if err := anondb.FinishUploadSession(dataset.Name, dataset.UploadSessionData.SessionID); err != nil {
			log.Printf("Error during upload session maintenance: %v", err.Error())
			continue
		}

		log.Printf(
			"Maintenance: session '%v' for dataset '%v' successfully finished",
			dataset.UploadSessionData.SessionID,
			dataset.Name)
	}
}

// StopUploadSessionMaintenance stops the upload session maintenance background process
func StopUploadSessionMaintenance(timeout time.Duration) {
	atomic.StoreInt32(&stop, 1)
	stopChannel <- true

	select {
	case <-finishedChannel:
	case <-time.After(timeout):
	}

	log.Printf("Upload session maintenance background task finished successfully")
}
