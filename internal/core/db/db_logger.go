package db

import (
	"context"
	"errors"
	"github.com/sendya/pkg/log"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	gormlogger "gorm.io/gorm/logger"
)

var (
	gormPackage    = filepath.Join("gorm.io", "gorm")
	zapgormPackage = filepath.Join("internal/core/db", "zapgorm")
)

type GormLogger struct {
	dbLog                 *zap.Logger
	SlowThreshold         time.Duration
	SourceField           string
	SkipCallerLookup      bool
	SkipErrRecordNotFound bool
}

func NewLog(dbLog *zap.Logger) *GormLogger {
	newLog := dbLog.WithOptions(zap.IncreaseLevel(zap.DebugLevel))
	return &GormLogger{
		dbLog:                 newLog,
		SlowThreshold:         200 * time.Millisecond,
		SkipCallerLookup:      false,
		SkipErrRecordNotFound: true,
	}
}

func (l *GormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	return l
}

func (l *GormLogger) Info(ctx context.Context, s string, args ...interface{}) {
	l.logger().Info(s, log.Any("args", args))
}

func (l *GormLogger) Warn(ctx context.Context, s string, args ...interface{}) {
	l.logger().Warn(s, log.Any("args", args))
}

func (l *GormLogger) Error(ctx context.Context, s string, args ...interface{}) {
	l.logger().Error(s, log.Any("args", args))
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if !l.dbLog.Core().Enabled(log.DebugLevel) {
		return
	}
	elapsed := time.Since(begin)
	fields := make([]log.Field, 0)
	timeUsed := float64(elapsed.Nanoseconds()) / 1e6
	fields = append(fields, log.Float64("timeUsed(ms)", timeUsed))

	switch {
	case err != nil && !(errors.Is(err, gorm.ErrRecordNotFound) && l.SkipErrRecordNotFound):
		sql, rows := fc()
		fields = append(fields, log.String("error", err.Error()))
		fields = append(fields, log.String("sql", sql))
		fields = append(fields, log.Int64("rows", rows))
		l.logger().Error("", fields...)
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0:
		sql, rows := fc()
		fields = append(fields, log.String("sql", sql))
		fields = append(fields, log.Duration("SLOW SQL", l.SlowThreshold))
		fields = append(fields, log.Int64("rows", rows))
		l.logger().Warn("trace", fields...)
	default:
		sql, rows := fc()
		fields = append(fields, log.Int64("rows", rows))
		fields = append(fields, log.String("sql", sql))
		l.logger().Debug("trace", fields...)
	}
}

func (l *GormLogger) logger() *zap.Logger {
	for i := 2; i < 15; i++ {
		_, file, _, ok := runtime.Caller(i)
		switch {
		case !ok:
		case strings.HasSuffix(file, "_test.go"):
		case strings.Contains(file, gormPackage):
		case strings.Contains(file, zapgormPackage):
		default:
			return l.dbLog.WithOptions(zap.AddCallerSkip(i))
		}
	}
	return l.dbLog
}
