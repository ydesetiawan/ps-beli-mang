package http

import (
	"net/http"
	"strings"
)

type (
	// SkipFunc is a type function to determine if it should skip to log based on a request value.
	SkipFunc func(r *http.Request) bool
	Client   struct {
		maskingFunc MaskingFunc
		skipFunc    SkipFunc
	}
)

func New(opts ...Option) *Client {
	defaultOpt := &Client{
		maskingFunc: defaultMaskingFunc,
		skipFunc:    skipOnPing,
	}

	for _, opt := range opts {
		opt(defaultOpt)
	}

	return defaultOpt
}

func (c *Client) SlogLogger() *SlogLogger {
	return newSlogLogger(c)
}

func (c *Client) ZapLogger() *ZapLogger {
	return newZapLogger(c)
}

func skipOnPing(r *http.Request) bool {
	cleanPath := strings.TrimSuffix(r.URL.Path, "/")
	return cleanPath == "/ping"
}
