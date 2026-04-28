package common

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger
var SugaredLogger *zap.SugaredLogger

// SetupLogger initializes the global logger based on the current environment.
// In production mode, structured JSON logging is used; otherwise, console logging is used.
func SetupLogger() {
	var cfg zap.Config

	if gin.Mode() == gin.ReleaseMode {
		cfg = zap.NewProductionConfig()
		cfg.EncoderConfig.TimeKey = "timestamp"
		cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	} else {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// Allow overriding log level via environment variable
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel != "" {
		var level zapcore.Level
		if err := level.UnmarshalText([]byte(logLevel)); err == nil {
			cfg.Level = zap.NewAtomicLevelAt(level)
		}
	}

	var err error
	Logger, err = cfg.Build(zap.AddCallerSkip(0))
	if err != nil {
		// Fallback to a no-op logger if setup fails
		Logger = zap.NewNop()
	}
	SugaredLogger = Logger.Sugar()
}

// SyncLogger flushes any buffered log entries. Should be called before program exit.
func SyncLogger() {
	if Logger != nil {
		_ = Logger.Sync()
	}
}

// Info logs a message at info level.
func Info(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Info(msg, fields...)
	}
}

// Warn logs a message at warn level.
func Warn(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Warn(msg, fields...)
	}
}

// Error logs a message at error level.
func Error(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Error(msg, fields...)
	}
}

// Fatal logs a message at fatal level and then calls os.Exit(1).
func Fatal(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Fatal(msg, fields...)
	}
}

// Infof logs a formatted message at info level.
func Infof(format string, args ...interface{}) {
	if SugaredLogger != nil {
		SugaredLogger.Infof(format, args...)
	}
}

// Warnf logs a formatted message at warn level.
func Warnf(format string, args ...interface{}) {
	if SugaredLogger != nil {
		SugaredLogger.Warnf(format, args...)
	}
}

// Errorf logs a formatted message at error level.
func Errorf(format string, args ...interface{}) {
	if SugaredLogger != nil {
		SugaredLogger.Errorf(format, args...)
	}
}
