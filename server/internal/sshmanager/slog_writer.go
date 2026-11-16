package sshmanager

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"sync"
)

type slogWriter struct {
	mu    sync.Mutex
	buf   []byte
	log   *slog.Logger
	level slog.Level
}

func newSlogWriter(l *slog.Logger, level slog.Level) *slogWriter {
	return &slogWriter{log: l, level: level}
}

func (w *slogWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.buf = append(w.buf, p...)
	for {
		i := bytes.IndexByte(w.buf, '\n')
		if i < 0 {
			break
		}
		line := strings.TrimRight(string(w.buf[:i]), "\r")
		w.buf = w.buf[i+1:]
		if line != "" {
			w.log.Log(context.Background(), w.level, line)
		}
	}
	if len(w.buf) > 64*1024 { // no newline in 64 KB: flush instead of growing forever
		w.log.Log(context.Background(), w.level, string(w.buf))
		w.buf = w.buf[:0]
	}
	return len(p), nil
}

// Flush logs any trailing partial line. Call after session.Wait().
func (w *slogWriter) Flush() {
	w.mu.Lock()
	defer w.mu.Unlock()
	if len(w.buf) > 0 {
		w.log.Log(context.Background(), w.level, string(w.buf))
		w.buf = nil
	}
}