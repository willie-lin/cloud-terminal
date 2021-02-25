package utils

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Log zapLog对象
var Log *zap.Logger

// 日志切割设置
func getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "api.log", // 日志文件位置
		MaxSize:    10,        // 日志文件最大大小(MB)
		MaxBackups: 5,         // 保留旧文件最大数量
		MaxAge:     30,        // 保留旧文件最长天数
		Compress:   false,     // 是否压缩旧文件
	}
	return zapcore.AddSync(lumberJackLogger)
}

// 编码器
func getEncoder() zapcore.Encoder {
	// 使用默认的JSON编码
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

// InitLogger 初始化Logger
func InitLogger() {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	Log = zap.New(core, zap.AddCaller())
}

func createLogger(encoding string) (*zap.Logger, error) {
	if encoding == "json" {
		return zap.NewProduction()
	}
	return zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "console",
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}.Build()
}
