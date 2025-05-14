package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/adityaadpandey/students-api/internal/config"
	"github.com/adityaadpandey/students-api/internal/http/handlers/student"
)

func main() {
	cfg := config.MustLoad()

	router := http.NewServeMux()
	router.HandleFunc("POST /api/students", student.New())

	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}
	slog.Info("server startded", "addr", cfg.Addr)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	<-done

	slog.Info("shutting down server...")
	cts, cancel := context.WithTimeout(context.Background(), 5-time.Second)
	defer cancel()
	if err := server.Shutdown(cts); err != nil {
		slog.Error("failed to shutdown server", "error", err)
	} else {
		slog.Info("server shut down gracefully")
	}
}
