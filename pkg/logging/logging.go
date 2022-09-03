package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

// Init 初期化
func Init() {
	level := zap.NewAtomicLevelAt(zapcore.DebugLevel)

	config := zap.Config{
		Level:             level,
		Development:       true,
		DisableStacktrace: true,
		Encoding:          "console",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:       "Msg",
			LevelKey:         "Level",
			TimeKey:          "Time",
			NameKey:          "Name",
			CallerKey:        "Caller",
			FunctionKey:      "Func",
			StacktraceKey:    "Trace",
			EncodeLevel:      zapcore.CapitalColorLevelEncoder,
			EncodeTime:       zapcore.RFC3339NanoTimeEncoder,
			EncodeDuration:   zapcore.StringDurationEncoder,
			EncodeCaller:     zapcore.ShortCallerEncoder,
			ConsoleSeparator: "  ",
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	var err error
	logger, err = config.Build()
	if err != nil {
		panic(err)
	}
}

// Sync ログの終了
func Sync() {
	logger.Sync()
}

// Info 情報
func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

// Error エラー
func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

// Fatal 致命的
func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}
