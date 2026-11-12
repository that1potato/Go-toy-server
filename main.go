package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const DEFAULT_TIMEOUT time.Duration = time.Millisecond * 3000 // 3 second default timeout
const PORT = 8080
const NUMBER_OF_BACKGROUND_TASKS = 10

func registerHandlers(serveMux *http.ServeMux) {
	serveMux.Handle(
		"GET /health",
		middlewareHandler(http.HandlerFunc(healthHandler), DEFAULT_TIMEOUT),
	)
	serveMux.Handle(
		"POST /items",
		middlewareHandler(http.HandlerFunc(itemsHandler), DEFAULT_TIMEOUT),
	)
	serveMux.Handle(
		"POST /tenSec",
		middlewareHandler(http.HandlerFunc(tenSecHandler), DEFAULT_TIMEOUT),
	)

	e := experimentHandler{}
	serveMux.Handle("/exp", e)
}

func waitGroupWithContext(ctx context.Context, wg *sync.WaitGroup) error {
	done := make(chan struct{})
	go func() {
		defer close(done)
		wg.Wait()
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return nil
	}
}

func main() {
	// ---- starting the http server ----
	address := fmt.Sprintf(":%d", PORT)

	serveMux := http.NewServeMux()
	server := &http.Server{
		Addr:    address,
		Handler: serveMux,
	}

	registerHandlers(serveMux)

	go func() {
		log.Printf("Starting server at %s", address)
		if err := server.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal("An error occurred when the server is running: ", err)
			}
		}
	}()

	// ---- starting some background task ----
	taskWaitGroup := &sync.WaitGroup{}
	startBackgroundTask(taskWaitGroup, NUMBER_OF_BACKGROUND_TASKS)

	// ---- handle quit signal (ctl+c) ----
	quitSig := make(chan os.Signal, 1)
	signal.Notify(
		quitSig,
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	<-quitSig
	log.Print("Quit signal received, shutting down gracefully...")
	ctx, cancel := context.WithTimeout(
		context.Background(),
		30*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Shutdown error: ", err)
	}
	log.Println("Server exited cleanly, waiting for background tasks to exit...")

	if err := waitGroupWithContext(ctx, taskWaitGroup); err != nil {
		log.Println("Background tasks timed out.")
	} else {
		log.Println("Background tasks exited cleanly.")
	}
}
