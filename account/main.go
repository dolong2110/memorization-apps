package main

import (
	"context"
	"github.com/dolong2110/memorization-apps/account/router"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	// you could insert your favorite logger here for structured or leveled logging
	log.Println("Starting server...")

	// load configs
	config, err := router.GetConfig(".", "dev", "json")
	if err != nil {
		log.Fatalf("Failed to get config: %v\n", err)
	}

	ds, err := router.InitDS(config)
	if err != nil {
		log.Fatalf("Unable to initialize data sources: %v\n", err)
	}

	r := router.NewRouters(config, ds)

	engine, err := r.InitGin()
	if err != nil {
		log.Fatalf("Failed to init gin: %v\n", err)
	}

	srv := &http.Server{
		Addr:    ":" + config.Port,
		Handler: engine,
	}

	// Graceful server shutdown - https://github.com/gin-gonic/examples/blob/master/graceful-shutdown/graceful-shutdown/server.go
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to initialize server: %v\n", err)
		}
	}()

	log.Printf("Listening on port %v\n", srv.Addr)

	// Wait for kill signal of channel
	quit := make(chan os.Signal, 1)

	var signalsToIgnore = []os.Signal{os.Interrupt}
	signal.Notify(quit, signalsToIgnore...)

	// This blocks until a signal is passed into the quit channel
	<-quit

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// shutdown data sources
	if err := ds.Close(); err != nil {
		log.Fatalf("A problem occurred gracefully shutting down data sources: %v\n", err)
	}

	// Shutdown server
	log.Println("Shutting down server...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v\n", err)
	}
}
