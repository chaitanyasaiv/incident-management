package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/ChaitanyaSaiV/Incident-Management/internal/handlers"
	"github.com/ChaitanyaSaiV/Incident-Management/internal/logging"
	"github.com/ChaitanyaSaiV/Incident-Management/internal/router"
	"github.com/ChaitanyaSaiV/Incident-Management/internal/storage"
)

var healthCheck atomic.Bool

func main() {

	storageName := flag.String("storage", "inMemory", "flag to decide which storage to use during compile time")
	flag.Parse()
	var store handlers.IncidentStore
	switch *storageName {
	case "inMemory":
		store = storage.NewInMemoryStore()
	case "file":
		fileStore, err := storage.NewFileStorage("localFile.json")
		if err != nil {
			log.Fatal("Error creating new file storage")
		}
		store = fileStore
	default:
		log.Fatalf("unknown store type: %s", *storageName)
	}

	logging := logging.NewLoggingIncidentStore(store)

	handler := handlers.NewHandler(logging)

	server := router.Routes(":8080", handler)

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM) // This sends the message to the quit channel
	<-quit
	log.Println("Shutdown Signal Received")
	healthCheck.Store(false)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Forced shutdown after timeout: %v", err)
	}
	log.Println("Server stopped cleanly")
}
