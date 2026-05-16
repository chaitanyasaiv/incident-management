package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ChaitanyaSaiV/Incident-Management/internal/handlers"
	"github.com/ChaitanyaSaiV/Incident-Management/internal/router"
	"github.com/ChaitanyaSaiV/Incident-Management/internal/storage"
)

func main() {

	store := storage.NewInMemoryStore()
	handler := handlers.NewHandler(store)

	server := router.Routes("8080", handler)

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM) // This sends the message to the quit channel
	<-quit
	log.Println("Shutdown Signal Received")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Forced shutdown after timeout: %v", err)
	}
	log.Println("Server stopped cleanly")
}
