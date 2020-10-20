/*
 *    Copyright 2020 Yury Makarov
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 *
 */

package logger

import (
	"github.com/TheZeroSlave/zapsentry"
	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	once      sync.Once
	logger    *zap.Logger
)

func InitLogger(logPath string, logLevel string, sentryDSN string, version string) {
	level := zapcore.InfoLevel
	logLevel = strings.TrimSpace(logLevel)
	logLevel = strings.ToLower(logLevel)
	switch logLevel {
	case "debug":
		level = zapcore.DebugLevel
	case "warn", "warning":
		level = zapcore.WarnLevel
	case "err", "error":
		level = zapcore.ErrorLevel
	case "dpanic":
		level = zapcore.DPanicLevel
	case "panic":
		level = zapcore.PanicLevel
	case "fatal":
		level = zapcore.FatalLevel
	}

	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		if level > zapcore.ErrorLevel {
			return lvl >= level
		}
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		if level >= zapcore.ErrorLevel {
			return false
		}
		return lvl < zapcore.ErrorLevel && lvl >= level
	})
	wInfo := zapcore.AddSync(&lumberjack.Logger{
		Filename:   filepath.Join(logPath, "info.log"),
		MaxSize:    5, // megabytes
		MaxBackups: 30,
		MaxAge:     0,
		Compress:   true,
	})
	wErr := zapcore.AddSync(&lumberjack.Logger{
		Filename:   filepath.Join(logPath, "error.log"),
		MaxSize:    5, // megabytes
		MaxBackups: 30,
		MaxAge:     0,
		Compress:   true,
	})
	consoleInfo := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)

	fileEncoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, wErr, highPriority),
		zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
		zapcore.NewCore(fileEncoder, wInfo, lowPriority),
		zapcore.NewCore(consoleEncoder, consoleInfo, lowPriority),
	)

	logger = zap.New(core)
	if sentryDSN != "" {
		sentryClient, err := sentry.NewClient(sentry.ClientOptions{
			Release:          version,
			Dsn:              sentryDSN,
			DebugWriter:      os.Stderr,
			AttachStacktrace: true,
		})
		if err != nil {
			logger.Error("error init sentry client", zap.Error(err))
		}
		sentryCore, err := zapsentry.NewCore(zapsentry.Configuration{
			Tags:              nil,
			DisableStacktrace: false,
			Level:             zapcore.ErrorLevel,
			FlushTimeout:      0,
			Hub:               nil,
		}, zapsentry.NewSentryClientFromClient(sentryClient))
		if err != nil {
			logger.Error("error init sentry core", zap.Error(err))
		}
		logger.Debug("sentry attach")
		logger = zapsentry.AttachCoreToLogger(sentryCore, logger)
	}
}

func GetLogger() *zap.Logger {
	return logger
}

func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}
func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}
func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}
func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}
func DPanic(msg string, fields ...zap.Field) {
	logger.DPanic(msg, fields...)
}
func Panic(msg string, fields ...zap.Field) {
	logger.Panic(msg, fields...)
}
func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}
func Sync() error {
	return logger.Sync()
}
