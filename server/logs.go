package server

import (
	"io"
	"log/slog"
	"os"
)

func newLogger(w io.Writer) *slog.Logger {
	var level slog.Level
	if os.Getenv("LOG_LEVEL") == "debug" {
		level = slog.LevelDebug
	} else {
		level = slog.LevelInfo
	}
	opts := &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey && len(groups) == 0 {
				a.Value = slog.StringValue(
					a.Value.Time().Format("2006/01/02 15:04:05"))
			}
			return a
		},
	}
	return slog.New(slog.NewJSONHandler(w, opts))
}

var Log = newLogger(os.Stdout)

func Srv() *slog.Logger {
	return Log.With("source", "SERVER")
}

func (c *Client) log() *slog.Logger {
	return Log.With("source", c.label(), "ip", c.conn.RemoteAddr().String())
}
