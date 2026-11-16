package wrapper

import (
	"context"
	"errors"
	"os/exec"
	"time"
)

type Result struct {
	ExitCode int
	Stdout   string
	Stderr   string
	Duration time.Duration
	TimedOut bool
}

const maxOutput = 64 * 1024

// tailBuffer keeps only the last `max` bytes written to it.
type tailBuffer struct {
	buf       []byte
	max       int
	Truncated bool
}

func (t *tailBuffer) Write(p []byte) (int, error) {
	n := len(p)
	t.buf = append(t.buf, p...)
	if len(t.buf) > t.max {
		t.buf = t.buf[len(t.buf)-t.max:]
		t.Truncated = true
	}
	return n, nil
}

func (t *tailBuffer) String() string { return string(t.buf) }

func Run(ctx context.Context, command string) (Result, error) {

	cmd := exec.CommandContext(ctx, "/bin/sh", "-c", command)
	var stdout, stderr tailBuffer
	stdout.max, stderr.max = maxOutput, maxOutput
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	start := time.Now()
	err := cmd.Run()
	res := Result{
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		Duration: time.Since(start),
		TimedOut: errors.Is(ctx.Err(), context.DeadlineExceeded),
	}

	var exitErr *exec.ExitError
	switch {
	case err == nil:
		res.ExitCode = 0
	case errors.As(err, &exitErr):
		res.ExitCode = exitErr.ExitCode() // non-zero exit is a result, not a failure to run
		err = nil
	default:
		res.ExitCode = -1 // couldn't start (e.g. /bin/sh missing)
	}
	return res, err
}
