package sloggorm

import (
	"time"

	"golang.org/x/exp/slog"
	gormlogger "gorm.io/gorm/logger"
)

type Option struct {
	logger        *slog.Logger
	logLevel      gormlogger.LogLevel
	slowThreshold time.Duration
	skipCallerLookup          bool
	ignoreRecordNotFoundError bool
	parameterizedQueries      bool
}

func NewOption() *Option {
	return &Option{
		logger:                    slog.Default(),
		logLevel:                  gormlogger.Warn,
		slowThreshold:             100 * time.Millisecond, //nolint:gomnd
		skipCallerLookup:          false,
		ignoreRecordNotFoundError: true,
		parameterizedQueries:      true,
	}
}

func (o *Option) WithLogger(logger *slog.Logger) *Option {
	if logger != nil {
		o.logger = logger
	}

	return o
}

func (o *Option) WithLogLevel(logLevel gormlogger.LogLevel) *Option {
	o.logLevel = logLevel
	return o
}

func (o *Option) WithSlowThreshold(duration time.Duration) *Option {
	o.slowThreshold = duration
	return o
}

func (o *Option) WithSkipCallerLookup(skipCallerLookup bool) *Option {
	o.skipCallerLookup = skipCallerLookup
	return o
}

func (o *Option) WithIgnoreRecordNotFoundError(ignoreRecordNotFoundError bool) *Option {
	o.ignoreRecordNotFoundError = ignoreRecordNotFoundError
	return o
}

func (o *Option) WithParameterizedQueries(parameterizedQueries bool) *Option {
	o.parameterizedQueries = parameterizedQueries
	return o
}
