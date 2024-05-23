package sloggorm_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"bitbucket.org/jurnal/scm/pkg/logger"
	sloggorm "bitbucket.org/jurnal/scm/pkg/logger/gorm"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slog"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func TestNew(t *testing.T) {
	slogger, _ := logger.SlogOption{
		Leveler: slog.LevelError,
	}.NewSlog()
	l := sloggorm.New(sloggorm.NewOption().WithLogger(slogger))
	assert.NotNil(t, l)

	l2 := sloggorm.New(nil)
	assert.NotNil(t, l2)
}

func TestEveryLogMode(t *testing.T) {
	slogger, _ := logger.SlogOption{
		Leveler: slog.LevelDebug,
	}.NewSlog()
	l := sloggorm.New(sloggorm.NewOption().
		WithLogger(slogger))
	ctx := context.TODO()
	count := 0
	fc := func() (sql string, rowsAffected int64) {
		count++
		return "SELECT * FROM users WHERE id = ?", 0
	}

	l.LogMode(gormlogger.Silent)
	l.LogMode(gormlogger.Silent).Info(ctx, "logs should not be written")
	l.LogMode(gormlogger.Silent).Warn(ctx, "logs should not be written")
	l.LogMode(gormlogger.Silent).Error(ctx, "logs should not be written")
	l.LogMode(gormlogger.Silent).Trace(ctx, time.Now().Add(-100*time.Millisecond), fc, errors.New("logs should not be written")) //nolint:goerr113
	l.LogMode(gormlogger.Silent).Trace(ctx, time.Now().Add(-500*time.Millisecond), fc, nil)                                      // logs should not be written
	l.LogMode(gormlogger.Silent).Trace(ctx, time.Now(), fc, nil)                                                                 // logs should not be written
	l.LogMode(gormlogger.Silent).Trace(ctx, time.Now(), fc, gorm.ErrRecordNotFound)                                              // logs should not be written

	l.LogMode(gormlogger.Error)
	l.LogMode(gormlogger.Error).Info(ctx, "logs should not be written")
	l.LogMode(gormlogger.Error).Warn(ctx, "logs should not be written")
	l.LogMode(gormlogger.Error).Error(ctx, "log error gorm")
	l.LogMode(gormlogger.Error).Trace(ctx, time.Now().Add(-100*time.Millisecond), fc, errors.New("database error")) //nolint:goerr113
	l.LogMode(gormlogger.Error).Trace(ctx, time.Now().Add(-500*time.Millisecond), fc, nil)                          // logs should not be written
	l.LogMode(gormlogger.Error).Trace(ctx, time.Now(), fc, nil)                                                     // logs should not be written
	l.LogMode(gormlogger.Error).Trace(ctx, time.Now(), fc, gorm.ErrRecordNotFound)                                  // logs should not be written

	l.LogMode(gormlogger.Warn)
	l.LogMode(gormlogger.Warn).Info(ctx, "logs should not be written")
	l.LogMode(gormlogger.Warn).Warn(ctx, "log warn gorm")
	l.LogMode(gormlogger.Warn).Error(ctx, "log error gorm")
	l.LogMode(gormlogger.Warn).Trace(ctx, time.Now().Add(-100*time.Millisecond), fc, errors.New("database error")) //nolint:goerr113
	l.LogMode(gormlogger.Warn).Trace(ctx, time.Now().Add(-500*time.Millisecond), fc, nil)
	l.LogMode(gormlogger.Warn).Trace(ctx, time.Now(), fc, nil)                    // logs should not be written
	l.LogMode(gormlogger.Warn).Trace(ctx, time.Now(), fc, gorm.ErrRecordNotFound) // logs should not be written

	l = sloggorm.New(sloggorm.NewOption().
		WithLogger(slogger).
		WithIgnoreRecordNotFoundError(false))
	l.LogMode(gormlogger.Info)
	l.LogMode(gormlogger.Info).Info(ctx, "log info gorm")
	l.LogMode(gormlogger.Info).Warn(ctx, "log warn gorm")
	l.LogMode(gormlogger.Info).Error(ctx, "log error gorm")
	l.LogMode(gormlogger.Info).Trace(ctx, time.Now().Add(-100*time.Millisecond), fc, errors.New("database error")) //nolint:goerr113
	l.LogMode(gormlogger.Info).Trace(ctx, time.Now().Add(-500*time.Millisecond), fc, nil)
	l.LogMode(gormlogger.Info).Trace(ctx, time.Now(), fc, nil)
	l.LogMode(gormlogger.Info).Trace(ctx, time.Now(), fc, gorm.ErrRecordNotFound)

	assert.Equal(t, count, 7)
}
