package main

import (
	"context"
	"fmt"
	"github.com/asyauqi15/go-flag"
	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rdb := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})

	flagClient := flag.New(rdb)

	routes := chi.NewRouter()

	flagClient.InitiateRoutes(routes)

	srv := &http.Server{
		Addr:              "0.0.0.0:8000",
		Handler:           routes,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	errCh := make(chan error, 1)
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)

	go func() {
		log.Println("http server is running")
		if err := srv.ListenAndServe(); err != nil {
			errCh <- fmt.Errorf("failed to run http server: %w", err)
		}
	}()

	go func() {
		<-signalCh
		signal.Reset(os.Interrupt)
		errCh <- fmt.Errorf("interrupted")
	}()

	<-errCh

	if err := srv.Shutdown(ctx); err != nil {
		log.Println(err)
	}
}
