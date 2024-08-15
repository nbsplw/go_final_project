package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"main/core/config"
	"os"
)

var Logger *zap.Logger

func Init() {
	var zapConfig zap.Config
	var loggerCore zapcore.Core
	logLevel := zapcore.DebugLevel

	switch config.Get().Fiber.Mode {
	case "debug":
		zapConfig = zap.NewDevelopmentConfig()
		zapConfig.EncoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
		break
	default:
		zapConfig = zap.NewProductionConfig()
		logLevel = zapcore.ErrorLevel
		break
	}
	var encoder zapcore.Encoder
	zapConfig.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	if config.Get().Logger.Output == "json" {
		encoder = zapcore.NewJSONEncoder(zapConfig.EncoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(zapConfig.EncoderConfig)
	}

	loggerCore = zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), logLevel),
	)

	Logger = zap.New(loggerCore, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	Logger.Info("logger: init")
}

func Get() *zap.Logger {
	if Logger == nil {
		Init()
	}
	return Logger
}
