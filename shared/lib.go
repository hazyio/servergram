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


func Map[T, U any](in []T, f func(T) U) []U {
	out := make([]U, 0, len(in))
	for _, v := range in {
		out = append(out, f(v))
	}
	return out
}
