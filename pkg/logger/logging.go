package logger

import (
	"os"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	log  *zap.Logger
	once sync.Once
)

func Init() error {
	var err error
	once.Do(func() {
		cfg := zap.Config{
			OutputPaths: []string{getOutputLogs(), "stdout"},
			Level:       zap.NewAtomicLevelAt(getLevelLogs()),
			Encoding:    "json",
			EncoderConfig: zapcore.EncoderConfig{
				MessageKey:   "message",
				LevelKey:     "level",
				TimeKey:      "timestamp",
				EncodeTime:   zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
				EncodeLevel:  zapcore.LowercaseColorLevelEncoder,
				EncodeCaller: zapcore.ShortCallerEncoder,
			},
		}
		log, err = cfg.Build()
	})
	return err

}

// TODO: logging saindo em um json na raiz da pasta mas será trocado para algo mais viavel.
func getOutputLogs() string {
	output := strings.ToLower(strings.TrimSpace(os.Getenv("LOG_OUTPUT")))

	if output == "" {
		return "logs.json"
	}

	return output
}

func getLevelLogs() zapcore.Level {
	switch strings.ToLower(strings.TrimSpace(os.Getenv("LOG_LEVEL"))) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

func ensureLogger() {
	if log == nil {
		panic("logger.Init() must be called before any logging")
	}
}

func Info(msg string, fields ...zap.Field) {
	ensureLogger()
	log.Info(msg, fields...)
}

func Error(msg string, err error, fields ...zap.Field) {
	ensureLogger()
	tags := append(fields, zap.NamedError("error", err))
	log.Error(msg, tags...)
}

func Debug(msg string, fields ...zap.Field) {
	ensureLogger()
	log.Debug(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	ensureLogger()
	log.Warn(msg, fields...)
}

func Sync() error {
	if log == nil {
		return nil
	}

	return log.Sync()
}
