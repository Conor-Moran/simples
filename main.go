package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	addrPort := ":6969"
	directory := "."
	delegate := http.FileServer(http.Dir(directory))

	// create a new serve mux and register the handlers
	sm := http.NewServeMux()
	sm.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		log.Printf("Handling request for: %s\n", req.URL)
		delegate.ServeHTTP(w, req)
	})

	// create a new server
	s := http.Server{
		Addr:         addrPort,
		Handler:      sm,
		ErrorLog:     log.Default(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// start the server
	go func() {
		log.Printf("Serving HTTP on port: %s\n", addrPort)

		err := s.ListenAndServe()
		if err != nil {
			log.Printf("Server shutting down: %s\n", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	// Block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
}
