package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Log   *zap.Logger
	Sugar *zap.SugaredLogger
)

func Init(level, format string) error {
	var cfg zap.Config

	if format == "json" {
		cfg = zap.NewProductionConfig()
	} else {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	lvl, err := zapcore.ParseLevel(level)
	if err != nil {
		lvl = zapcore.InfoLevel
	}
	cfg.Level = zap.NewAtomicLevelAt(lvl)

	Log, err = cfg.Build()
	if err != nil {
		return err
	}

	Sugar = Log.Sugar()
	return nil
}

func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
}

func InitDefault() {
	env := os.Getenv("APP_ENV")
	logLevel := os.Getenv("LOG_LEVEL")

	if logLevel == "" {
		if env == "production" {
			logLevel = "info"
		} else {
			logLevel = "debug"
		}
	}

	format := "console"
	if env == "production" {
		format = "json"
	}

	if err := Init(logLevel, format); err != nil {
		panic(err)
	}
}
