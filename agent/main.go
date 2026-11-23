package main

import (
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
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
		fs := flag.NewFlagSet("run-job", flag.ContinueOnError)
		fs.SetOutput(os.Stderr)
		fs.Usage = func() {
			fmt.Fprintln(os.Stderr, `usage: agent run-job -id=<id> [-manual] "command"`)
			fs.PrintDefaults()
		}
		id := fs.Uint64("id", 0, "job id (required)")
		manual := fs.Bool("manual", false, "mark run as manually triggered")

		if err := fs.Parse(os.Args[2:]); err != nil {
			if errors.Is(err, flag.ErrHelp) {
				os.Exit(0)
			}
			os.Exit(2)
		}

		command := strings.Join(fs.Args(), " ")
		if command == "" {
			fmt.Fprintln(os.Stderr, "command must not be empty")
			fs.Usage()
			os.Exit(2)
		}
		if *id == 0 {
			fmt.Fprintln(os.Stderr, "-id is required")
			fs.Usage()
			os.Exit(2)
		}

		trigger := db.CronRunTriggerScheduled
		if *manual {
			trigger = db.CronRunTriggerManual
		}
		if err := wrapper.RunJob(database, *id, command, trigger); err != nil {
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