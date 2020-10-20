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
