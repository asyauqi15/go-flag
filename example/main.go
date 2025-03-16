package main

import (
	"context"
	"encoding/json"
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

	flagClient, err := flag.New(rdb)
	if err != nil {
		log.Fatalf("failed to initiate flag client: %v", err)
	}

	routes := chi.NewRouter()
	routes.Get("/example/{feature_name}", func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "feature_name")
		isActive, err := flagClient.IsActive(r.Context(), name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp := map[string]any{
			"feature_name": name,
			"is_active":    isActive,
		}

		respJson, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(respJson)
	})

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
