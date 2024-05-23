package sloggorm

import (
	"context"
	"errors"
	"time"

	"golang.org/x/exp/slog"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type Logger struct {
	logger                    *slog.Logger
	logLevel                  gormlogger.LogLevel
	slowThreshold             time.Duration
	skipCallerLookup          bool
	ignoreRecordNotFoundError bool
	parameterizedQueries      bool
}

func New(opt *Option) gormlogger.Interface {
	if opt == nil {
		opt = NewOption()
	}

	return Logger{
		logger:                    opt.logger,
		logLevel:                  opt.logLevel,
		slowThreshold:             opt.slowThreshold,
		skipCallerLookup:          opt.skipCallerLookup,
		ignoreRecordNotFoundError: opt.ignoreRecordNotFoundError,
		parameterizedQueries:      opt.parameterizedQueries,
	}
}

func (l Logger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	return Logger{
		logLevel:                  level,
		logger:                    l.logger,
		slowThreshold:             l.slowThreshold,
		skipCallerLookup:          l.skipCallerLookup,
		ignoreRecordNotFoundError: l.ignoreRecordNotFoundError,
		parameterizedQueries:      l.parameterizedQueries,
	}
}

func (l Logger) Info(ctx context.Context, s string, i ...interface{}) {
	if l.logLevel < gormlogger.Info {
		return
	}
	l.logger.InfoCtx(ctx, s, i...)
}

func (l Logger) Warn(ctx context.Context, s string, i ...interface{}) {
	if l.logLevel < gormlogger.Warn {
		return
	}
	l.logger.WarnCtx(ctx, s, i...)
}

func (l Logger) Error(ctx context.Context, s string, i ...interface{}) {
	if l.logLevel < gormlogger.Error {
		return
	}
	l.logger.ErrorCtx(ctx, s, i...)
}

func (l Logger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.logLevel <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && l.logLevel >= gormlogger.Error && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.ignoreRecordNotFoundError):
		sql, rows := fc()
		l.Error(ctx, err.Error(),
			slog.String("sql", sql),
			slog.String("duration", elapsed.String()),
			slog.Int64("rows", rows),
		)
	case elapsed > l.slowThreshold && l.slowThreshold != 0 && l.logLevel >= gormlogger.Warn:
		sql, rows := fc()
		l.Warn(ctx, "slow query",
			slog.String("sql", sql),
			slog.String("duration", elapsed.String()),
			slog.Int64("rows", rows),
		)
	case l.logLevel == gormlogger.Info:
		sql, rows := fc()
		l.Info(ctx, "gorm_logger",
			slog.String("sql", sql),
			slog.String("duration", elapsed.String()),
			slog.Int64("rows", rows),
		)
	}
}

func (l Logger) ParamsFilter(_ context.Context, sql string, params ...interface{}) (string, []interface{}) {
	if l.parameterizedQueries {
		return sql, nil
	}
	return sql, params
}
