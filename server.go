package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	"go-server-counters/config"
	handlers "go-server-counters/handers"
	"go-server-counters/service"
	"go-server-counters/transport"
)

func runServer(cfg *config.Config) {
	redisCfg := transport.NewRedisConfig(cfg.RedisAddress, cfg.RedisPassword, cfg.RedisDB)
	redisTransport := transport.NewRedisTransport(redisCfg)

	svc := service.NewMessageCounter(redisTransport)
	handlers := handlers.NewHandler(svc)

	// Create a new serve mux
	mux := mux.NewRouter()

	// Add basic health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	mux.HandleFunc("/increment", handlers.SendMessage).Methods("POST")
	mux.HandleFunc("/read", handlers.ReadMessages).Methods("PATCH")
	mux.HandleFunc("/get", handlers.GetMessageCount).Methods("GET")

	// Create server with timeouts
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// Channel to listen for errors coming from the listener.
	serverErrors := make(chan error, 1)

	// Start the server
	go func() {
		log.Printf("Server is starting on port %s", cfg.Port)
		serverErrors <- srv.ListenAndServe()
	}()

	// Channel to listen for an interrupt or terminate signal from the OS.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Blocking select waiting for either a server error or a signal.
	select {
	case err := <-serverErrors:
		log.Fatalf("Error starting server: %v", err)

	case sig := <-shutdown:
		log.Printf("Start shutdown... Signal: %v", sig)

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		// Shutdown the server
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("Could not stop server gracefully: %v", err)
			os.Exit(1)
		}
	}
}
