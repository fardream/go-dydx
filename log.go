package dydx

import (
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log Logger = SetupLogger("info", "GODYDX_LOGLEVEL")

func SetupLogger(defaultLevel, envName string) Logger {
	level_str := os.Getenv(envName)
	if level_str == "" {
		level_str = defaultLevel
	}
	level, err := zapcore.ParseLevel(strings.ToUpper(level_str))
	if err != nil {
		panic(fmt.Errorf("failed to parse level from environment variable %s=\"%s\": %+v", envName, level_str, err))
	}

	var config zap.Config
	if level > zap.DebugLevel {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	config.Level.SetLevel(level)

	logger, _ := config.Build()
	return logger.Sugar()
}

// Logger defines an interface to log.
// By default this is uber's zap.
type Logger interface {
	Debugf(string, ...any)
	Infof(string, ...any)
	Warnf(string, ...any)
	Errorf(string, ...any)
	Fatalf(string, ...any)
}

var _ Logger = (*zap.SugaredLogger)(nil)
