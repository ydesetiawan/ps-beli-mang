package http

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type ZapLogger struct {
	maskingFunc MaskingFunc
	isSkip      SkipFunc
}

func newZapLogger(c *Client) *ZapLogger {
	return &ZapLogger{
		maskingFunc: c.maskingFunc,
		isSkip:      c.skipFunc,
	}
}

func (z *ZapLogger) Middleware(next http.Handler) http.Handler {
	return http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if z.isSkip(r) {
			next.ServeHTTP(w, r)
			return
		}

		startTime := time.Now()

		reqData := captureRequestData(r, z.maskingFunc)

		customRespWriter := &customResponseWriter{
			ResponseWriter: w,
		}
		next.ServeHTTP(customRespWriter, r)

		respData := captureResponseFromWriter(customRespWriter, z.maskingFunc)

		zap.L().Info("http_logger",
			zap.String("duration", time.Since(startTime).String()),
			zap.Object("request", reqData),
			zap.Object("response", respData),
		)
	}))
}
