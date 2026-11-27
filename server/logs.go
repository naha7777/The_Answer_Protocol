package server

import (
	"bytes"
	"fmt"
	"golang.org/x/term"
	"io"
	"log/slog"
	"os"
)

const (
	reset   = "\033[0m"
	red     = "\033[31m"
	yellow  = "\033[33m"
	RED     = "\033[41m"
	gray    = "\033[90m"
	green   = "\033[32m"
	cyan    = "\033[36m"
	magenta = "\033[35m"
)

type colorWriter struct {
	w io.Writer
}

func (c colorWriter) Write(p []byte) (int, error) {
	color := ""
	switch {
	case bytes.Contains(p, []byte(`client disconnected`)),
		bytes.Contains(p, []byte(`client connected`)):
		color = yellow
	case bytes.Contains(p, []byte(`server shutting down`)),
		bytes.Contains(p, []byte(`server started`)),
		bytes.Contains(p, []byte(`server stopped`)),
		bytes.Contains(p, []byte(`shutdown requested`)):
		color = magenta
	case !bytes.Contains(p, []byte(`,"source":"SERVER"`)):
		color = cyan
	case bytes.Contains(p, []byte(`,"level":"ERROR"`)):
		color = RED
	case bytes.Contains(p, []byte(`,"level":"WARN"`)):
		color = red
	case bytes.Contains(p, []byte(`EVT`)):
		color = green
	case bytes.Contains(p, []byte(`,"level":"INFO"`)):
		color = reset
	}
	if color == "" {
		return c.w.Write(p)
	}
	line := bytes.TrimRight(p, "\n")
	if _, err := fmt.Fprintf(c.w, "%s%s%s\n", color, line, reset); err != nil {
		return 0, err
	}
	return len(p), nil
}

func stdoutWriter() io.Writer {
	if term.IsTerminal(int(os.Stdout.Fd())) {
		return colorWriter{os.Stdout}
	}
	return os.Stdout
}

func newLogger(w io.Writer) *slog.Logger {
	var level slog.Level
	switch os.Getenv("LOG_LEVEL") {
	case "debug":
		level = slog.LevelDebug
	case "off":
		level = slog.Level(100)
	default:
		level = slog.LevelInfo
	}
	opts := &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey && len(groups) == 0 {
				a.Value = slog.StringValue(
					a.Value.Time().Format("2006/01/02 15:04:05.000"))
			}
			return a
		},
	}
	return slog.New(slog.NewJSONHandler(w, opts))
}

var Log = newLogger(stdoutWriter())

func Srv() *slog.Logger {
	return Log.With("source", "SERVER")
}

func (c *Client) log() *slog.Logger {
	return Log.With("source", c.label(), "ip", c.conn.RemoteAddr().String())
}

func InitLogger(path string) (func(), error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, err
	}
	Log = newLogger(io.MultiWriter(stdoutWriter(), f))
	return func() { f.Close() }, nil
}
