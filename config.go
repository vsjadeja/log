package log

import (
	"go.uber.org/zap/zapcore"
)

var (
	productionEncoderConfig = zapcore.EncoderConfig{
		TimeKey:        `time`,
		LevelKey:       `level`,
		NameKey:        `logger`,
		CallerKey:      `caller`,
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     `message`,
		StacktraceKey:  `stacktrace`,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	developmentEncoderConfig = zapcore.EncoderConfig{
		// Keys can be anything except the empty string.
		TimeKey:        `T`,
		LevelKey:       `L`,
		NameKey:        `N`,
		CallerKey:      `C`,
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     `M`,
		StacktraceKey:  `S`,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout(`2006-01-02 15:04:05`),
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
)
