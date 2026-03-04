package terminal

import (
	"context"
	"fmt"
	"io"
	"log/slog"
)

func SimpleHandler(w io.Writer, level slog.Level) slog.Handler {
	return &simpleHandler{w: w, level: level}
}

type simpleHandler struct {
	w     io.Writer
	level slog.Level
}

func (h *simpleHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *simpleHandler) WithAttrs([]slog.Attr) slog.Handler { return h }
func (h *simpleHandler) WithGroup(string) slog.Handler      { return h }

func (h *simpleHandler) Handle(_ context.Context, r slog.Record) error {
	line := "[" + r.Level.String() + "] " + r.Message
	r.Attrs(func(a slog.Attr) bool {
		line += " " + a.Key + "=" + fmt.Sprint(a.Value.Any())
		return true
	})
	_, err := fmt.Fprintln(h.w, line)
	return err
}
