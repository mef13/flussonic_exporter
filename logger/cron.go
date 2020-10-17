package logger

import (
	"go.uber.org/zap"
)

type cronLogger struct {
	logger  *zap.Logger
	jobName string
}

// GetLoggerForCron returns a struct that implements cron.Logger interface.
func GetLoggerForCron(jobName string) cronLogger {
	return cronLogger{logger: GetLogger(), jobName: jobName}
}

// Info implements the cron.Logger interface.
func (cl cronLogger) Info(msg string, keysAndValues ...interface{}) {
	cl.logger.Info(msg, zap.String("JobName", cl.jobName), zap.Any("params", keysAndValues))
}

// Error implements the cron.Logger interface.
func (cl cronLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	cl.logger.Error(msg, zap.String("JobName", cl.jobName), zap.Any("params", keysAndValues))
}
