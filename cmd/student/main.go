package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/smartcraze/student-api/internal/config"
)

func main() {
	// load mustload for the configuration file
	cfg := config.MustLoad()
	// database setup

	//setup the router
	router := http.NewServeMux()
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to student api"))
	})

	// setup server

	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}
	fmt.Printf("server is started:  %s", cfg.Addr)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("failed to start the server")
		}
	}()
	<-done // blocking

	slog.Info("Shutting Down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)

	if err != nil {
		slog.Error("failed to Shutdown server", slog.String("error", err.Error()))
	}

}
