package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewGlobalZap(cfg Config, opts ...zap.Option) (*zap.Logger, error) {
	zapOpts := []zap.Option{
		zap.AddStacktrace(zap.ErrorLevel),
		zap.AddCaller(),
	}

	if len(opts) > 0 {
		zapOpts = append(zapOpts, opts...)
	}

	level, err := zap.ParseAtomicLevel(cfg.Level)
	if err != nil {
		return nil, err
	}

	config := zap.NewProductionConfig()
	config.Encoding = customEncoder
	config.Level = level
	config.EncoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.MessageKey = "message"

	l, err := config.Build(zapOpts...)
	if err != nil {
		return nil, err
	}
	_ = zap.ReplaceGlobals(l)

	return l, nil
}
