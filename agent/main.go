package main

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"hazyio.com/servergram/agent/internal"
	"hazyio.com/servergram/agent/internal/db"
	"hazyio.com/servergram/agent/internal/wrapper"
	"hazyio.com/servergram/shared"
)

var logLevel slog.LevelVar
var transportOut = os.Stdout

func setupLogging() {
	if os.Getenv("LOG_LEVEL") == "debug" {
		logLevel.Set(slog.LevelDebug)
	}
	h := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: &logLevel})
	slog.SetDefault(slog.New(h))
}

func main() {
	os.Stdout = os.Stderr

	setupLogging()
	if len(os.Args) < 2 {
		shared.Fatal("usage: agent <serve|run-job> ...")
	}
	database, err := db.InitDatabase()
	if err != nil {
		shared.Fatal("failed to initialize database: %v", err)
	}

	switch os.Args[1] {
	case "run-job":
		if len(os.Args) != 4 && len(os.Args) != 5 {
			fmt.Fprintln(os.Stderr, `usage: agent run-job "command" id=<id>`)
			os.Exit(2)
		}
		command := os.Args[2]
		if command == "" {
			fmt.Fprintln(os.Stderr, "command must not be empty")
			os.Exit(2)
		}
		idStr, ok := strings.CutPrefix(os.Args[3], "id=")
		if !ok {
			fmt.Fprintln(os.Stderr, "expected id=<id> as the last argument")
			os.Exit(2)
		}
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			fmt.Fprintln(os.Stderr, "invalid job id")
			os.Exit(2)
		}
		trigger := db.CronRunTriggerScheduled
		if len(os.Args) > 4 && os.Args[4] == "--manual" {
			trigger = db.CronRunTriggerManual
		}
		if err := wrapper.RunJob(database, id, command, trigger); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		os.Exit(0)
	case "serve":
		internal.Serve(database, transportOut)

	default:
		shared.Fatal("unknown command: %s", os.Args[1])
	}

}
