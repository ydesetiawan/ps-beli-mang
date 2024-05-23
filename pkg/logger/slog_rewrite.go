package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
	"golang.org/x/exp/slog"
)

const (
	SDKVersionKey   = "logger.sdk.version"
	SDKVersionValue = "v1.0.0"
)

type Resource map[string]string

type AttributeFormatterFunc func(groups []string, a slog.Attr) slog.Attr

type SlogOption struct {
	Resource           Resource
	ContextExtractor   ContextExtractFunc
	AttributeFormatter AttributeFormatterFunc
	Writer             io.Writer
	Leveler            slog.Leveler
}

type slogHandler struct {
	opt           SlogOption
	groups        []string
	attributesEnc zapcore.Encoder
}

var _ slog.Handler = (*slogHandler)(nil)

func NewSlogWithDefault() (*slog.Logger, error) {
	return SlogOption{
		Resource: Resource{
			SDKVersionKey: SDKVersionValue,
		},
	}.NewSlog()
}

func (opt SlogOption) NewSlog() (*slog.Logger, error) {
	// Prepare encoder to be use again.
	zapEnc := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		// SkipLineEnding ensure that the last byte is not contains new line,
		// because we need these bytes to be appended later.
		SkipLineEnding: true,
		EncodeTime:     zapcore.RFC3339NanoTimeEncoder, // used by zap.Time("k", time.Time{})
		EncodeDuration: zapcore.MillisDurationEncoder,  // used by zap.Duration("k", time.Duration())
	})

	// Always override sdk.version key
	if opt.Resource == nil {
		opt.Resource = make(Resource)
	}

	opt.Resource[SDKVersionKey] = SDKVersionValue

	h := &slogHandler{
		opt:           opt,
		groups:        make([]string, 0),
		attributesEnc: zapEnc,
	}
	return slog.New(h), nil
}

func (s *slogHandler) Enabled(_ context.Context, l slog.Level) bool {
	minLevel := slog.LevelDebug
	if s.opt.Leveler != nil {
		minLevel = s.opt.Leveler.Level()
	}

	return l >= minLevel
}

// Handle will handler the log format (adding attribute, etc.) before log is written.
// The same attribute key will be printed (not discarded) similar like in Uber Zap logger.
func (s *slogHandler) Handle(ctx context.Context, record slog.Record) error {
	// Clone, to ensure that when s.attributesEnc called again, the previous attributes not stick.
	attributesEncoder := s.attributesEnc.Clone()

	// This is the left-over attributes that defined by user in the:
	// slog.DebugCtx(context.Background(), "test debug", slog.String("last-pushed-attr", "value"))
	// So, this need to be formatted before printed.
	record.Attrs(func(attr slog.Attr) bool {
		if s.opt.AttributeFormatter != nil {
			attr = s.opt.AttributeFormatter(s.groups, attr)
		}

		convertAttrToField(attr).AddTo(attributesEncoder)
		return true
	})

	var attributesBytes string
	{
		// This will create valid JSON output made by zapcore.NewJSONEncoder.
		// Use empty zapcore.Entry and nil zapcore.Field to ensure it only print
		// the object we append in the field using AddTo function.
		attributesBuf, err := attributesEncoder.EncodeEntry(zapcore.Entry{}, nil)
		if err != nil {
			err = fmt.Errorf("cannot encode attributes entry: %w", err)
			return err
		}

		attributesBytes = attributesBuf.String() // when empty, it will return: {}
		attributesBuf.Free()
	}

	var logLineBytes string
	{
		finalLogEnc := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			MessageKey: "message",
			LevelKey:   "level",
			TimeKey:    "timestamp",
			NameKey:    "logger",
			CallerKey:  "caller",
			LineEnding: zapcore.DefaultLineEnding,
			EncodeLevel: func(level zapcore.Level, encoder zapcore.PrimitiveArrayEncoder) {
				encoder.AppendString(level.CapitalString())
			},
			EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
			EncodeDuration: zapcore.MillisDurationEncoder,
			EncodeCaller:   zapcore.FullCallerEncoder,
		})

		// From [OTEP0114](https://github.com/open-telemetry/oteps/pull/114)
		// https://github.com/open-telemetry/opentelemetry-specification/blob/v1.18.0/specification/logs/README.md?plain=1#L474-L526
		spanCtx := trace.SpanContextFromContext(ctx)
		zap.String("trace_id", spanCtx.TraceID().String()).AddTo(finalLogEnc)
		zap.String("span_id", spanCtx.SpanID().String()).AddTo(finalLogEnc)

		// 00: no flags
		zap.String("trace_flags", spanCtx.TraceFlags().String()).AddTo(finalLogEnc)

		// Add resource.
		if s.opt.Resource != nil {
			zap.Any("resource", s.opt.Resource).AddTo(finalLogEnc)
		} else {
			zap.Any("resource", Resource{SDKVersionKey: SDKVersionValue}).AddTo(finalLogEnc)
		}

		// extractCtx does not need to be re-formatted using AttributeFormatterFunc
		// Because, user already have the control of the returned value when creating
		// the extractor function. So, it MUST be already correctly formatted.
		if s.opt.ContextExtractor != nil {
			// OpenNamespace will make any AddTo call after this line become under "context".
			// Beware of this behavior! But, you SHOULD NOT worry about this as long as you're not adding
			// new top-level field name. Top-level field name is field with the same level as "timestamp", "message", etc.
			finalLogEnc.OpenNamespace("context")
			for _, attr := range s.opt.ContextExtractor(ctx) {
				convertAttrToField(attr).AddTo(finalLogEnc)
			}
		}

		logEntry := zapcore.Entry{
			Level:      convertSlogLevel(record.Level),
			Time:       record.Time,
			LoggerName: "",
			Message:    record.Message,
			Caller:     zapcore.EntryCaller{},
			Stack:      "",
		}

		// Add caller.
		if record.PC > 0 {
			frames := runtime.CallersFrames([]uintptr{record.PC})
			if frames != nil {
				frame, _ := frames.Next()
				logEntry.Caller = zapcore.EntryCaller{
					Defined:  true,
					PC:       frame.PC,
					File:     frame.File,
					Line:     frame.Line,
					Function: frame.Function,
				}
			}
		}

		logLineBuf, err := finalLogEnc.EncodeEntry(logEntry, nil)
		if err != nil {
			err = fmt.Errorf("cannot encode final log entry: %w", err)
			return err
		}

		logLineBuf.TrimNewline()           // trim new line to ensure we can merge two bytes.
		logLineBytes = logLineBuf.String() // when empty this will return: {}
		logLineBuf.Free()
	}

	// We use append attributes (in spite of open new namespace) to make sure that
	// attributes will be added on the last element in the log output.
	// This makes the log output easier to read because if it truncated in your screen, we will still can see
	// the trace id in the first occurrence.
	finalLogBytes := AppendAttributes(logLineBytes, attributesBytes)

	if s.opt.Writer != nil {
		_, _errWrite := s.opt.Writer.Write(finalLogBytes)
		return _errWrite
	}

	// By default, using stdout as target output.
	_, _errWrite := os.Stdout.Write(finalLogBytes)
	return _errWrite
}

// WithAttrs is called when user call slog.With(attrs...)
// or when slog.With(attrs..).WithGroup(name).With(attrs...)
//
// If user already call WithGroup, the appended valued will append to this group name child's.
// We need to call AttributeFormatterFunc here before appended to the Zap logger buffer.
func (s *slogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// clone to ensure the next call will not contain the same attributes from previous calls
	zapEncCloned := s.attributesEnc.Clone()

	for _, attr := range attrs {
		if s.opt.AttributeFormatter != nil {
			attr = s.opt.AttributeFormatter(s.groups, attr)
		}

		convertAttrToField(attr).AddTo(zapEncCloned)
	}

	return &slogHandler{
		opt:           s.opt,
		groups:        s.groups,
		attributesEnc: zapEncCloned,
	}
}

// WithGroup open new namespace in the zapcore.Encoder so the following WithAttrs
// will append to this group name.
func (s *slogHandler) WithGroup(name string) slog.Handler {
	zapEncoder := s.attributesEnc.Clone()
	zapEncoder.OpenNamespace(name)

	return &slogHandler{
		opt:           s.opt,
		groups:        append(s.groups, name),
		attributesEnc: zapEncoder,
	}
}

// groupObject holds all the slog.Attr saved in a slog.GroupValue.
type groupObject []slog.Attr

func (gs groupObject) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	for _, attr := range gs {
		convertAttrToField(attr).AddTo(enc)
	}
	return nil
}

// convertAttrToField convert slog.Attr to zapcore.Field
// based on slog.Kind.
func convertAttrToField(attr slog.Attr) zapcore.Field {
	switch attr.Value.Kind() {
	case slog.KindBool:
		return zap.Bool(attr.Key, attr.Value.Bool())
	case slog.KindDuration:
		return zap.Duration(attr.Key, attr.Value.Duration())
	case slog.KindFloat64:
		return zap.Float64(attr.Key, attr.Value.Float64())
	case slog.KindInt64:
		return zap.Int64(attr.Key, attr.Value.Int64())
	case slog.KindString:
		return zap.String(attr.Key, attr.Value.String())
	case slog.KindTime:
		return zap.Time(attr.Key, attr.Value.Time())
	case slog.KindUint64:
		return zap.Uint64(attr.Key, attr.Value.Uint64())
	case slog.KindGroup:
		return zap.Object(attr.Key, groupObject(attr.Value.Group()))
	case slog.KindLogValuer:
		return convertAttrToField(slog.Attr{
			Key:   attr.Key,
			Value: attr.Value.Resolve(),
		})

	default:
		return zap.Any(attr.Key, attr.Value.Any())
	}
}

// convertSlogLevel maps slog Levels to zap Levels.
// Note that there is some room between slog levels while zap levels are continuous, so we can't 1:1 map them.
// See also https://go.googlesource.com/proposal/+/master/design/56345-structured-logging.md?pli=1#levels
func convertSlogLevel(l slog.Level) zapcore.Level {
	switch {
	case l >= slog.LevelError:
		return zapcore.ErrorLevel
	case l >= slog.LevelWarn:
		return zapcore.WarnLevel
	case l >= slog.LevelInfo:
		return zapcore.InfoLevel
	default:
		return zapcore.DebugLevel
	}
}

// AppendAttributes will append s2 into s1 string, with key "attributes".
// s1 and s2 is assumed to have valid JSON structure since this SHOULD come from
// Zap buffer output. If this is not valid JSON, then the output of this function will also not valid JSON.
//
// As the last mitigation, we check whether the s1 is surrounded by `{` and `}`.
// This mean that the s1 is JSON object, and we can safely put any key-value in the top level JSON.
// For example, s1= {"foo": "bar"} and s2= {"a": "b"}
// This output function will: {"foo": "bar", "attributes": {"a": "b"}}
//
// Key in JSON must always string, so this payload is not valid JSON: `{ 0: { "foo": "bar" }}` because 0 is not string.
// And the output function will become: { 0: { "foo": "bar" }, "attributes": {"a": "b"}}
// which also not valid JSON.
//
// See the JSON Spec here https://www.json.org/json-en.html
func AppendAttributes(s1, s2 string) []byte {
	const minLength = 2
	if len(s1) < minLength {
		return []byte(s1)
	}

	lenByte := len(s1)
	firstChar := s1[:1]
	lastChar := s1[lenByte-1:]

	if firstChar != "{" || lastChar != "}" {
		return []byte(s1)
	}

	buf := buffer.NewPool().Get()

	// To be valid JSON, it must start with { and end with }
	buf.AppendString(s1[:lenByte-1])

	// Append string must only exist if previous bytes is not "{".
	// In other word, this is mean that the length of s1 is more than 2 and not only contains "{}".
	if lenByte > minLength {
		buf.AppendString(",")
	}

	if s2 != "" {
		buf.AppendString(`"attributes":` + s2)
	}

	buf.AppendString("}")
	buf.AppendString("\n") // add new line after append attributes

	defer buf.Free()
	return buf.Bytes()
}
