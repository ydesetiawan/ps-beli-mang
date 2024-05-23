package sloggorm_test

import (
	"testing"
	"time"

	"bitbucket.org/jurnal/scm/pkg/logger"
	sloggorm "bitbucket.org/jurnal/scm/pkg/logger/gorm"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slog"
	gormLogger "gorm.io/gorm/logger"
)

func TestNewOption(t *testing.T) {
	o := sloggorm.NewOption()
	assert.NotNil(t, o)
}

func TestOption_WithLogger(t *testing.T) {
	tests := []struct {
		name   string
		logger *slog.Logger
	}{
		{
			name:   "default slogger",
			logger: slog.Default(),
		},
		{
			name: "custom slogger",
			logger: func() *slog.Logger {
				slogger, _ := logger.SlogOption{
					Leveler: slog.LevelError,
				}.NewSlog()
				return slogger
			}(),
		},
		{
			name:   "send nil",
			logger: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := sloggorm.NewOption().WithLogger(tt.logger)
			assert.NotNil(t, o)
		})
	}
}

func TestOption_WithLogLevel(t *testing.T) {
	tests := []struct {
		name     string
		logLevel gormLogger.LogLevel
	}{
		{
			name:     "log level silent",
			logLevel: gormLogger.Silent,
		},
		{
			name:     "log level error",
			logLevel: gormLogger.Error,
		},
		{
			name:     "log level warn",
			logLevel: gormLogger.Warn,
		},
		{
			name:     "log level info",
			logLevel: gormLogger.Info,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := sloggorm.NewOption().WithLogLevel(tt.logLevel)
			assert.NotNil(t, o)
		})
	}
}

func TestOption_WithSlowThreshold(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
	}{
		{
			name:     "500 ms",
			duration: 500 * time.Millisecond,
		},
		{
			name:     "1 second",
			duration: 1 * time.Second,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := sloggorm.NewOption().WithSlowThreshold(tt.duration)
			assert.NotNil(t, o)
		})
	}
}

func TestOption_WithSkipCallerLookup(t *testing.T) {
	tests := []struct {
		name             string
		skipCallerLookup bool
	}{
		{
			name:             "true",
			skipCallerLookup: true,
		},
		{
			name:             "false",
			skipCallerLookup: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := sloggorm.NewOption().WithSkipCallerLookup(tt.skipCallerLookup)
			assert.NotNil(t, o)
		})
	}
}

func TestOption_WithIgnoreRecordNotFoundError(t *testing.T) {
	tests := []struct {
		name                      string
		ignoreRecordNotFoundError bool
	}{
		{
			name:                      "true",
			ignoreRecordNotFoundError: true,
		},
		{
			name:                      "false",
			ignoreRecordNotFoundError: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := sloggorm.NewOption().WithIgnoreRecordNotFoundError(tt.ignoreRecordNotFoundError)
			assert.NotNil(t, o)
		})
	}
}

func TestOption_WithParameterizedQueries(t *testing.T) {
	tests := []struct {
		name                 string
		parameterizedQueries bool
	}{
		{
			name:                 "true",
			parameterizedQueries: true,
		},
		{
			name:                 "false",
			parameterizedQueries: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := sloggorm.NewOption().WithParameterizedQueries(tt.parameterizedQueries)
			assert.NotNil(t, o)
		})
	}
}
