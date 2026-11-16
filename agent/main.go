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

func main() {
	h := slog.NewJSONHandler(os.Stderr, nil)
	slog.SetDefault(slog.New(h))
	if len(os.Args) < 2 {
		shared.Fatal("usage: agent <serve|run-job> ...")
	}
	database, err := db.InitDatabase()
	if err != nil {
		shared.Fatal("failed to initialize database: %v", err)
	}

	switch os.Args[1] {
	case "run-job":
		if len(os.Args) != 4 {
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
		if err := wrapper.RunJob(database, id, command, db.CronRunTriggerScheduled); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	case "serve":
		internal.Serve(database)

	default:
		shared.Fatal("unknown command: %s", os.Args[1])
	}

}
