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

	"github.com/joho/godotenv"
	"github.com/smartcraze/student-api/internal/config"
	httphandler "github.com/smartcraze/student-api/internal/http"
	"github.com/smartcraze/student-api/internal/storage"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	cfg := config.MustLoad()
	// database setup
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)
	db, err := storage.NewPostgresStorage(connStr)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer db.Close()

	slog.Info("Database connected successfully")

	//setup the router
	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to student api"))
	})

	router.HandleFunc("POST /api/student/create", httphandler.CreateStudentHandler(db))

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

	err = server.Shutdown(ctx)

	if err != nil {
		slog.Error("failed to Shutdown server", slog.String("error", err.Error()))
	}

}
