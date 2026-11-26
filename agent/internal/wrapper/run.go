package wrapper

import (
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
	"sync"
	"syscall"
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

// fileSink serialises writes from stdout and stderr (which exec copies on
// separate goroutines) into one file. Write errors are swallowed so a full
// disk can't make the command itself fail.
type fileSink struct {
	mu sync.Mutex
	f  *os.File
}

func (s *fileSink) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, _ = s.f.Write(p) // unbuffered: goes to the OS immediately
	return len(p), nil
}

func Run(ctx context.Context, command string, saveFile *string) (Result, error) {
	cmd := exec.CommandContext(ctx, "/bin/sh", "-c", command)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Cancel = func() error {
		// negative pid = whole process group
		return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	}
	cmd.WaitDelay = 3 * time.Second // safety net if something escapes the group
	var stdout, stderr tailBuffer
	stdout.max, stderr.max = maxOutput, maxOutput

	var outW, errW io.Writer = &stdout, &stderr

	if saveFile != nil {
		f, err := os.OpenFile(*saveFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
		if err != nil {
			return Result{ExitCode: -1}, err
		}
		defer f.Close()

		sink := &fileSink{f: f}
		outW = io.MultiWriter(&stdout, sink)
		errW = io.MultiWriter(&stderr, sink)
	}

	cmd.Stdout = outW
	cmd.Stderr = errW

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
		res.ExitCode = exitErr.ExitCode()
		err = nil
	default:
		res.ExitCode = -1
	}
	return res, err
}
