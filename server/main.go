package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"hazyio.com/servergram/server/internal"
)

func main() {
	startServer()
}

func startServer() {
	env := internal.LoadEnv()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	initLogger()
	database := InitDatabase()
	polls := &internal.Polls{}
	// ensure all polls are stopped when the server shuts down
	defer polls.StopAll()

	router := NewRouter(ctx, database, polls)
	slog.Info("Starting server at https://localhost:" + env.ServerPort)
	// tls certificate and key are provided by the server configuration, so empty strings are passed here
	go func() {
		if err := router.ListenAndServeTLS("", ""); err != nil {
			slog.Error("HTTP server stopped", "error", err)
			os.Exit(1)
		}
	}()
	// watch and wait for shutdown signal
	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = router.Shutdown(shutdownCtx)

}

func initLogger() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

}
