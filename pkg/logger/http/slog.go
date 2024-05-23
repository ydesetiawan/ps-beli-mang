package http

import (
	"net/http"
	"time"

	"golang.org/x/exp/slog"
)

type SlogLogger struct {
	maskingFunc MaskingFunc
	isSkip      SkipFunc
}

func newSlogLogger(c *Client) *SlogLogger {
	return &SlogLogger{
		maskingFunc: c.maskingFunc,
		isSkip:      c.skipFunc,
	}
}

func (s *SlogLogger) Middleware(next http.Handler) http.Handler {
	return http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.isSkip(r) {
			next.ServeHTTP(w, r)
			return
		}

		startTime := time.Now()

		reqData := captureRequestData(r, s.maskingFunc)

		customRespWriter := &customResponseWriter{
			ResponseWriter: w,
		}
		next.ServeHTTP(customRespWriter, r)

		respData := captureResponseFromWriter(customRespWriter, s.maskingFunc)

		slog.InfoCtx(r.Context(), "http_logger",
			slog.String("duration", time.Since(startTime).String()),
			slog.Any("request", reqData),
			slog.Any("response", respData),
		)
	}))
}

func (s *SlogLogger) ClientTransport(name string, rt http.RoundTripper) http.RoundTripper {
	return &slogRoundTripper{rt: rt, maskingFunc: s.maskingFunc, name: name}
}

type slogRoundTripper struct {
	name        string
	rt          http.RoundTripper
	maskingFunc MaskingFunc
}

// RoundTrip is a method to satisfy the http.RoundTripper interface.
// It will log the request and the response data from the another http.RoundTripper.
func (s *slogRoundTripper) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	startTime := time.Now()
	reqData := captureRequestData(req, s.maskingFunc)

	resp, err = s.rt.RoundTrip(req)
	duration := time.Since(startTime)

	if err != nil {
		slog.WarnCtx(req.Context(), s.name,
			slog.String("duration", duration.String()),
			slog.String("error", err.Error()),
			slog.Any("request", reqData),
		)
		return
	}

	respData := captureResponseFromHTTP(resp, s.maskingFunc)

	slog.InfoCtx(req.Context(), s.name,
		slog.String("duration", duration.String()),
		slog.Any("request", reqData),
		slog.Any("response", respData),
	)

	return
}
