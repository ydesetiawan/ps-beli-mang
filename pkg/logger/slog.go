package logger

import (
	"context"
	"fmt"

	"golang.org/x/exp/slog"
)

func NewGlobalSlog(handler slog.Handler) error {
	if handler == nil {
		return fmt.Errorf("invalid handler") //nolint:goerr113
	}

	sl := slog.New(handler)
	slog.SetDefault(sl)

	return nil
}

type (
	// ContextExtractFunc is a type function to extract context values and returns the array of slog.Attr.
	// Using this function gives the client flexibility to determine the way they extract the context values
	// and transform it to the slog.Attr.
	ContextExtractFunc func(ctx context.Context) []slog.Attr

	// SlogContextAdapter will extract the context values using ContextExtractFunc
	// and attach them to the next adapter. The next adapter would be anything that satisfy
	// the slog.Handler (i.e zapslog).
	SlogContextAdapter struct {
		next        slog.Handler
		extractFunc ContextExtractFunc
	}
)

func defaultContextExtractFunc(_ context.Context) []slog.Attr {
	return []slog.Attr{}
}

var _ slog.Handler = (*SlogContextAdapter)(nil)

func NewSlogContextAdapter(next slog.Handler, extractFunc ContextExtractFunc) *SlogContextAdapter {
	if extractFunc == nil {
		extractFunc = defaultContextExtractFunc
	}
	return &SlogContextAdapter{
		next:        next,
		extractFunc: extractFunc,
	}
}

func (s *SlogContextAdapter) Enabled(ctx context.Context, level slog.Level) bool {
	return s.next.Enabled(ctx, level)
}

func (s *SlogContextAdapter) Handle(ctx context.Context, record slog.Record) error {
	record.AddAttrs(s.extractFunc(ctx)...)
	return s.next.Handle(ctx, record)
}

func (s *SlogContextAdapter) WithAttrs(attrs []slog.Attr) slog.Handler {
	return s.next.WithAttrs(attrs)
}

func (s *SlogContextAdapter) WithGroup(name string) slog.Handler {
	return s.next.WithGroup(name)
}
