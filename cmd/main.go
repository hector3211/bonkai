package main

import (
	"context"
	"log"
	"my-chi/pkg"
	"my-chi/pkg/handlers"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	db, err := pkg.NewDB()
	if err != nil || db == nil {
		log.Fatal("failed initializing DB")
	}
	ctx, cancle := context.WithCancel(context.Background())
	defer cancle()

	userRouter := handlers.NewUserRouter(ctx, db)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	// User Router
	r.Get("/users", func(w http.ResponseWriter, r *http.Request) {
		userRouter.GetAllusers(w, r)
	})
	r.Post("/new-user", func(w http.ResponseWriter, r *http.Request) {
		userRouter.CreateUser(w, r)
	})

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
