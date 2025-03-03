package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	myDb "my-chi/db"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	_ "modernc.org/sqlite"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	ctx := context.Background()

	wd, _ := os.Getwd()

	schema, err := os.ReadFile(filepath.Join(fmt.Sprintf("%s/schema.sql", wd)))
	if err != nil {
		slog.Error("failed reading schema.sql", "error", err)
		os.Exit(1)
	}

	// db, err := sql.Open("sqlite", ":memory:")
	db, err := sql.Open("sqlite", "../db.sqlite")
	if err != nil {
		slog.Error("failed initializing db")
		os.Exit(1)
	}

	// create tables
	if _, err := db.ExecContext(ctx, string(schema)); err != nil {
		slog.Error("failed executing seed")
		os.Exit(1)
	}
	queries := myDb.New(db)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})
	r.Get("/users", func(w http.ResponseWriter, r *http.Request) {
		data, err := queries.ListUsers(ctx)
		if err != nil {
			http.Error(w, "no data found", http.StatusNoContent)
			return
		}
		j, _ := json.Marshal(data)
		w.Write([]byte(j))
		w.WriteHeader(200)
	})

	r.Post("/users/new", func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("name")
		email := r.FormValue("email")

		if name == "" || email == "" {
			slog.Error("failed recieving name and email from from", "error", err)
			w.WriteHeader(500)
			return
		}

		insertedUser, err := queries.CreateUser(ctx, myDb.CreateUserParams{
			ID:    uuid.New(),
			Name:  name,
			Email: email,
		})
		if err != nil {
			slog.Error("failed inserting new user", "error", err)
			http.Error(w, "failed inserting user", http.StatusInternalServerError)
			return

		}

		data, _ := json.Marshal(insertedUser)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(data)
	})

	// User Router
	// r.Get("/users", userRouter.GetAllusers)
	// r.Post("/new-user", userRouter.CreateUser)

	server := &http.Server{
		Addr:    ":3000",
		Handler: r,
	}

	go func() {
		log.Println("Server is running on port 3000...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Block until we receive an interrupt signal
	<-sigChan
	log.Println("Shutting down server...")

	// Gracefully shutdown the server
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("server shutdown failed: %v", err)
	}
}
