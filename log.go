package dydx

import (
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.SugaredLogger

// set up the logger
func init() {
	level_str := os.Getenv("GODYDX_LOG_LEVEL")
	if level_str == "" {
		level_str = "info"
	}
	level, err := zapcore.ParseLevel(strings.ToUpper(level_str))
	if err != nil {
		panic(fmt.Errorf("failed to parse level from environment variable GODYDX_LOG_LEVEL=\"%s\": %+v", level_str, err))
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
	log = logger.Sugar()
}
