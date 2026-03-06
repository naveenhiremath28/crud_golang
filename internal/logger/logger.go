package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New creates a SugaredLogger. Uses colored console output at Debug level
// by default. Set APP_ENV=production for JSON output at Info level.
func New() (*zap.SugaredLogger, error) {
	var cfg zap.Config

	if os.Getenv("APP_ENV") == "production" {
		cfg = zap.NewProductionConfig()
		cfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	} else {
		cfg = zap.NewDevelopmentConfig()
		cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	zapLogger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return zapLogger.Sugar(), nil
}
