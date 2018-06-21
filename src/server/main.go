package main

import (
	"anonbll"
	"anondb"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"swagger"
	"syscall"
	"time"
)

func main() {
	log.Printf("Anonymization webserver starting up...")
	defer log.Printf("Anonymization webserver stopped gracefully")

	if err := anondb.InitConnection(); err != nil {
		log.Fatalf("Error while initializing database connection: %v", err.Error())
	}
	defer anondb.CloseConnection()

	if err := anondb.SetupDatabase(); err != nil {
		log.Printf("Error while initializing database: %v", err.Error())
		return
	}

	anonbll.StartUploadSessionMaintenance()
	defer anonbll.StopUploadSessionMaintenance(time.Second * 5)

	runServer()
}

func runServer() bool {
	var signalStop = make(chan os.Signal)
	signal.Notify(signalStop, syscall.SIGHUP, os.Interrupt, syscall.SIGTERM)

	var errorStop = make(chan bool)

	router := swagger.NewRouter()
	server := &http.Server{Addr: ":9137", Handler: router}

	go listenAndServe(server, errorStop)

	select {
	case <-errorStop:
		return false

	case signal := <-signalStop:
		log.Printf("Received signal '%v', starting graceful shutdown...", signal)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Error: %v", err)
			return false
		}

		return true
	}
}

func listenAndServe(server *http.Server, errorStop chan bool) {
	log.Printf("Listening on http://0.0.0.0%v", server.Addr)

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("Error: %v", err)
		errorStop <- true
	}

	log.Printf("Stopped listening")
}
