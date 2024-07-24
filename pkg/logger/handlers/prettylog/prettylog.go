package prettylog

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"strconv"
	"sync"
)

type Handler struct {
	h slog.Handler
	b *bytes.Buffer
	m *sync.Mutex
}

const (
	timeFormat = "[15:04:05.000]"

	reset = "\033[0m"

	black        = 30
	Red          = 31
	green        = 32
	yellow       = 33
	blue         = 34
	magenta      = 35
	cyan         = 36
	lightGray    = 37
	darkGray     = 90
	lightRed     = 91
	lightGreen   = 92
	lightYellow  = 93
	lightBlue    = 94
	lightMagenta = 95
	lightCyan    = 96
	white        = 97
)

func NewHandler(opts *slog.HandlerOptions) *Handler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}

	b := &bytes.Buffer{}

	return &Handler{
		b: b,
		h: slog.NewJSONHandler(b, &slog.HandlerOptions{
			Level:       opts.Level,
			AddSource:   opts.AddSource,
			ReplaceAttr: SuppressDefaults(opts.ReplaceAttr),
		}),
		m: &sync.Mutex{},
	}
}

func Colorize(colorCode int, v string) string {
	return fmt.Sprintf("\033[%sm%s%s", strconv.Itoa(colorCode), v, reset)
}

func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.h.Enabled(ctx, level)
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &Handler{h: h.h.WithAttrs(attrs), b: h.b, m: h.m}
}

func (h *Handler) WithGroup(name string) slog.Handler {
	return &Handler{h: h.h.WithGroup(name), b: h.b, m: h.m}
}

func (h *Handler) Handle(ctx context.Context, rec slog.Record) error {
	level := rec.Level.String() + ":"

	switch rec.Level {
	case slog.LevelDebug:
		level = Colorize(darkGray, level)
	case slog.LevelInfo:
		level = Colorize(cyan, level)
	case slog.LevelWarn:
		level = Colorize(lightYellow, level)
	case slog.LevelError:
		level = Colorize(lightRed, level)
	}

	attrs, err := h.ComputeAttrs(ctx, rec)
	if err != nil {
		return err
	}

	attrBytes, err := json.MarshalIndent(attrs, "", "  ")
	if err != nil {
		return fmt.Errorf("error when marshaling attrs: %w", err)
	}

	log.Println(
		Colorize(lightGray, rec.Time.Format(timeFormat)),
		level,
		Colorize(white, rec.Message),
		Colorize(darkGray, string(attrBytes)),
	)

	return nil
}

func (h *Handler) ComputeAttrs(
	ctx context.Context,
	rec slog.Record,
) (map[string]any, error) {
	h.m.Lock()
	defer func() {
		h.b.Reset()
		h.m.Unlock()
	}()

	if err := h.h.Handle(ctx, rec); err != nil {
		return nil, fmt.Errorf("error when calling inner handler's Handle: %w", err)
	}

	var attrs map[string]any

	err := json.Unmarshal(h.b.Bytes(), &attrs)
	if err != nil {
		return nil, fmt.Errorf("error when unmarshaling inner handler's Handle result: %w", err)
	}

	return attrs, nil
}

func SuppressDefaults(
	next func([]string, slog.Attr) slog.Attr,
) func([]string, slog.Attr) slog.Attr {
	return func(groups []string, attr slog.Attr) slog.Attr {
		if attr.Key == slog.TimeKey ||
			attr.Key == slog.LevelKey ||
			attr.Key == slog.MessageKey {
			return slog.Attr{}
		}

		if next == nil {
			return attr
		}

		return next(groups, attr)
	}
}
