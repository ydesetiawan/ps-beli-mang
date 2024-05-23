package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

const customEncoder = "custom_encoder"

var skipKeys = map[string]struct{}{
	"password": {},
	"api_key":  {},
}

func init() {
	_ = zap.RegisterEncoder(customEncoder, func(config zapcore.EncoderConfig) (zapcore.Encoder, error) {
		encoder := zapcore.NewJSONEncoder(config)
		return &Encoder{encoder}, nil
	})
}

type Encoder struct {
	zapcore.Encoder
}

func (e *Encoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	filtered := make([]zapcore.Field, 0, len(fields))
	for _, field := range fields {
		if _, ok := skipKeys[field.Key]; ok {
			continue
		}

		filtered = append(filtered, field)
	}
	return e.Encoder.EncodeEntry(entry, filtered)
}
