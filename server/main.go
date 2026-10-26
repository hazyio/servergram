package main

import (
	"log/slog"
	"net/http"
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

	router := NewRouter(database)
	slog.Info("Starting server at http://0.0.0.0:" + env.ServerPort)
	if err := http.ListenAndServe(":"+env.ServerPort, router); err != nil {
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
