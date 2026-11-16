package shared

import (
	"fmt"
	"log/slog"
	"os"
)

func Fatal(msg string, args ...any) {
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	slog.Error(msg)
	os.Exit(1)
}
