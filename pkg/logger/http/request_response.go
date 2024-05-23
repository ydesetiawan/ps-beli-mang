package http

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"go.uber.org/zap/zapcore"
	"golang.org/x/exp/slog"
)

type requestData struct {
	URL    string          `json:"url"`
	Method string          `json:"method"`
	Body   json.RawMessage `json:"body"`
}

func captureRequestData(r *http.Request, fn MaskingFunc) *requestData {
	requestLog := &requestData{
		URL:    r.URL.String(),
		Method: r.Method,
	}

	if r.Body == nil {
		return requestLog
	}

	buf, err := io.ReadAll(r.Body)
	if err != nil {
		return requestLog
	}

	//_ = r.Body.Close()
	r.Body = io.NopCloser(bytes.NewBuffer(buf))

	if len(buf) == 0 {
		return requestLog
	}

	if fn == nil {
		requestLog.Body = buf
		return requestLog
	}

	bodyStr := fn(string(buf))
	requestLog.Body = json.RawMessage(bodyStr)
	return requestLog
}

func (r *requestData) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("url", r.URL)
	enc.AddString("method", r.Method)
	_ = enc.AddReflected("body", r.Body)

	return nil
}

// LogValue method is to satisfy the slog.LogValuer interface.
// https://pkg.go.dev/golang.org/x/exp/slog#hdr-Customizing_a_type_s_logging_behavior
func (r *requestData) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("url", r.URL),
		slog.String("method", r.Method),
		slog.Any("body", r.Body),
	)
}

type responseData struct {
	StatusCode int             `json:"status_code"`
	Body       json.RawMessage `json:"body"`
}

func captureResponseFromHTTP(resp *http.Response, fn MaskingFunc) *responseData {
	if resp == nil {
		return &responseData{}
	}

	responseLog := &responseData{
		StatusCode: resp.StatusCode,
	}

	if resp.Body == nil {
		return responseLog
	}

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return responseLog
	}

	//_ = resp.Body.Close()
	resp.Body = io.NopCloser(bytes.NewBuffer(buf))

	if len(buf) == 0 {
		return responseLog
	}

	if fn == nil {
		responseLog.Body = buf
		return responseLog
	}

	bodyStr := fn(string(buf))
	responseLog.Body = json.RawMessage(bodyStr)

	return responseLog
}

func captureResponseFromWriter(w *customResponseWriter, fn MaskingFunc) *responseData {
	responseLog := &responseData{
		StatusCode: w.statusCode,
	}

	if w.body == nil || len(w.body) == 0 {
		return responseLog
	}

	if fn == nil {
		responseLog.Body = w.body
		return responseLog
	}

	bodyStr := fn(string(w.body))
	responseLog.Body = json.RawMessage(bodyStr)

	return responseLog
}

func (r *responseData) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt("status_code", r.StatusCode)
	_ = enc.AddReflected("body", r.Body)

	return nil
}

// LogValue method is to satisfy the slog.LogValuer interface.
// https://pkg.go.dev/golang.org/x/exp/slog#hdr-Customizing_a_type_s_logging_behavior
func (r *responseData) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Int("status_code", r.StatusCode),
		slog.Any("body", r.Body),
	)
}

// customRepsonseWriter is an adapter for the http.ResponseWriter.
// we can't directly get the response values from the http.ResponseWriter,
// so we need to make an adapter that capture all the values,
// but also write to the real http.ResponseWriter.
type customResponseWriter struct {
	http.ResponseWriter
	body       []byte
	statusCode int
}

func (w *customResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w *customResponseWriter) Write(b []byte) (int, error) {
	w.body = b
	return w.ResponseWriter.Write(b)
}

func (w *customResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
