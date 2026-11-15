package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"hazyio.com/servergram/agent/internal"
	"hazyio.com/servergram/agent/internal/db"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: agent <serve|run-job> ...")
	}
	database, err := db.InitDatabase()
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	switch os.Args[1] {
	case "run-job":
		if len(os.Args) != 3 {
			fmt.Fprintln(os.Stderr, "usage: agent run-job <id>")
			os.Exit(2)
		}
		id, err := strconv.ParseUint(os.Args[2], 10, 64) // "7; rm -rf" fails here
		if err != nil {
			fmt.Fprintln(os.Stderr, "invalid job id")
			os.Exit(2)
		}
		if err := internal.RunJob(id, 0); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	default:
		internal.Serve(database)
	}

}
