package logger

import (
	"context"
	"os"

	"github.com/willie-lin/cloud-terminal/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var global *zap.Logger

// Init initializes the global logger from config.
func Init(cfg config.LoggerConfig) error {
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(cfg.Level)); err != nil {
		level = zapcore.InfoLevel
	}

	encCfg := zap.NewProductionEncoderConfig()
	encCfg.TimeKey = "ts"
	encCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encCfg.LevelKey = "level"
	encCfg.EncodeLevel = zapcore.CapitalLevelEncoder
	encCfg.CallerKey = "caller"
	encCfg.EncodeCaller = zapcore.ShortCallerEncoder
	encCfg.MessageKey = "msg"

	var encoder zapcore.Encoder
	if cfg.Format == "console" {
		devEnc := zap.NewDevelopmentEncoderConfig()
		devEnc.EncodeTime = encCfg.EncodeTime
		encoder = zapcore.NewConsoleEncoder(devEnc)
	} else {
		encoder = zapcore.NewJSONEncoder(encCfg)
	}

	ws := zapcore.AddSync(os.Stdout)
	if cfg.ConsoleToStd == "stderr" {
		ws = zapcore.AddSync(os.Stderr)
	}

	core := zapcore.NewCore(encoder, ws, level)

	opts := []zap.Option{zap.AddCaller()}
	if cfg.Development {
		opts = append(opts, zap.Development())
	}

	global = zap.New(core, opts...)
	return nil
}

// New creates and returns a new logger instance.
func New(cfg config.LoggerConfig) (*zap.Logger, error) {
	if err := Init(cfg); err != nil {
		return nil, err
	}
	return global, nil
}

// Shutdown flushes the global logger.
func Shutdown(_ context.Context) {
	if global != nil {
		_ = global.Sync()
	}
}

// L returns the global logger.
func L() *zap.Logger {
	if global == nil {
		panic("logger not initialized")
	}
	return global
}

// FromContext returns a logger with trace fields (no-op without otel).
func FromContext(_ context.Context) *zap.Logger {
	return L()
}

// Info logs at info level.
func Info(msg string, fields ...zap.Field) {
	L().Info(msg, fields...)
}

// Error logs at error level.
func Error(msg string, fields ...zap.Field) {
	L().Error(msg, fields...)
}

// Debug logs at debug level.
func Debug(msg string, fields ...zap.Field) {
	L().Debug(msg, fields...)
}

// Warn logs at warn level.
func Warn(msg string, fields ...zap.Field) {
	L().Warn(msg, fields...)
}

// Fatal logs at fatal level.
func Fatal(msg string, fields ...zap.Field) {
	L().Fatal(msg, fields...)
}
