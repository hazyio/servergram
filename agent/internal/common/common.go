package common

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

type Result[T any] struct {
	Val T
	Err error
}

func RunAll[T any](fns ...func() (T, error)) []Result[T] {
	out := make([]Result[T], len(fns))
	var wg sync.WaitGroup
	wg.Add(len(fns))

	for i, fn := range fns {
		go func() {
			defer wg.Done()
			v, err := fn()
			out[i] = Result[T]{v, err}
		}()
	}
	wg.Wait()
	return out
}

func FileExt(name string) string {
	ext := filepath.Ext(name)
	if ext == name { // hidden file with no real extension
		return ""
	}
	return strings.ToLower(strings.TrimPrefix(ext, "."))
}

const sniffLen = 8192

// IsReadableText reports whether f is a regular file whose start looks like UTF-8 text.
// The error is non-nil if the file couldn't be read at all.
func IsReadableText(f *os.File) (bool, error) {
	info, err := f.Stat()
	if err != nil {
		return false, err
	}
	if !info.Mode().IsRegular() { // skip dirs, devices, pipes, sockets
		return false, nil
	}
	if info.Size() == 0 {
		return true, nil
	}

	buf := make([]byte, sniffLen)
	n, err := io.ReadFull(f, buf)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return false, err // permission or I/O failure
	}
	return looksLikeText(buf[:n], n == sniffLen), nil
}

func looksLikeText(b []byte, truncated bool) bool {
	if bytes.IndexByte(b, 0) >= 0 { // NUL byte: almost always binary
		return false
	}
	if truncated {
		// The sniff window may have cut a multi-byte rune in half; drop it.
		for i := 1; i <= utf8.UTFMax && i <= len(b); i++ {
			if utf8.RuneStart(b[len(b)-i]) {
				if !utf8.FullRune(b[len(b)-i:]) {
					b = b[:len(b)-i]
				}
				break
			}
		}
	}
	return utf8.Valid(b)
}

func LogError(e error) error {
	if e != nil {
		slog.Error("error", "err", e)
	}
	return e
}

func UnWrapLogErr[T any](data T, e error) *T {
	if e != nil {
		slog.Error("error", "err", e)
		return nil
	}
	return &data
}
func GetSavePath() string {
	if p := os.Getenv("SERVERGRAM_BIN_SAVE_PATH"); p != "" {
		return p
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".servergram")
}

func GetAgentPath() string {
	if p := os.Getenv("SERVERGRAM_AGENT_PATH"); p != "" {
		return p
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".servergram", "agent")
}

func HasAgentPrefix(s string) bool {
	return strings.HasPrefix(s, fmt.Sprintf("%s run-job", GetAgentPath()))
}

func Map[T, U any](in []T, f func(T) U) []U {
	out := make([]U, 0, len(in))
	for _, v := range in {
		out = append(out, f(v))
	}
	return out
}
func GetLogSavePath() string {
	return filepath.Join(GetSavePath(), "logs")
}

func MakeLogFileName(id uint64, t time.Time) (*string, error) {
	saveDir := filepath.Join(GetLogSavePath(), fmt.Sprintf("%d", id))
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		slog.Error("failed to create log directory", "dir", saveDir, "err", err)
		return nil, err
	}
	path := filepath.Join(saveDir, fmt.Sprintf("log-%s.log", t.Format(time.RFC3339)))
	return &path, nil
}

func DeleteDirectory(path string) error {
	err := os.RemoveAll(path)
	if err != nil {
		slog.Error("failed to delete directory", "path", path, "err", err)
	}
	return err
}
func RemoveCronRunLog(data []struct {
	StartedAt time.Time
	ID        uint64
}) {
	for _, d := range data {
		path, err := MakeLogFileName(d.ID, d.StartedAt)
		if err != nil {
			continue
		}
		if err := os.Remove(*path); err != nil {
			slog.Error("failed to remove cron run log", "path", path, "err", err)
		}
	}
}
