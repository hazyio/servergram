package main

import (
	"log/slog"
	"os"

	"hazyio.com/servergram/internal"
)

func main() {
	startServer()
}

func startServer() {
	env := internal.LoadEnv()

	initLogger()
	database := InitDatabase()
	polls := &internal.Polls{}
	// ensure all polls are stopped when the server shuts down
	defer polls.StopAll()

	router := NewRouter(database, polls)
	slog.Info("Starting server at https://localhost:" + env.ServerPort)
	// tls certificate and key are provided by the server configuration, so empty strings are passed here
	if err := router.ListenAndServeTLS("", ""); err != nil {
		slog.Error("HTTP server stopped", "error", err)
		os.Exit(1)
	}
}

func initLogger() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

}
