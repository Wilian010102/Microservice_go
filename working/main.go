package main

import (
	"context"
	"log"
	"microservicesgo/working/handlers"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	l := log.New(os.Stdout, "product-api ", log.LstdFlags)

	ph := handlers.NewProducts(l)

	sm := http.NewServeMux()
	sm.Handle("/", ph)

	server := &http.Server{
		Addr:         ":9090",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			l.Fatalf("Error starting the server: %s", err)
		}
	}()

	l.Println("Server running on port 9090")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	l.Println("Received termination signal, graceful shutdown", sig)

	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	server.Shutdown(tc)
}
