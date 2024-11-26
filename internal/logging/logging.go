/*
 *    Copyright 2024 okdp.io
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
 */

package log

import (
	"fmt"
	"sync"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/okdp/okdp-server/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	instance *zap.SugaredLogger
	once     sync.Once
)

func SetupGlobalLogger(loggingConf config.Logging) {
	once.Do(func() {

		if (loggingConf.Level == ""){
			loggingConf.Level = "info"
		}
		if (loggingConf.Format == ""){
			loggingConf.Format = "console"
		}
		
		fmt.Println("Initializing logger with: logging Level: ", loggingConf.Level, "Logging format: ", loggingConf.Format)

		logLevel, err := zapcore.ParseLevel(loggingConf.Level)

		if err != nil {
			panic(err)
		}

		config := zap.NewProductionConfig()
		config.Level = zap.NewAtomicLevelAt(logLevel)
		config.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
		config.Encoding = loggingConf.Format

		logger, err := config.Build()
		if err != nil {
			panic(err)
		}

		zap.ReplaceGlobals(zap.Must(logger, err))
		instance = logger.Sugar()
	})
}

// Logs in 'DEBUG' level according to a format specifier and writes to standard output.
func Debug(args ...interface{}) {
	instance.Debugf(args[0].(string), args[1:]...)
}

// Logs in 'INFO' level according to a format specifier and writes to standard output.
func Info(args ...interface{}) {
	instance.Infof(args[0].(string), args[1:]...)
}

// Logs in 'WARN' level according to a format specifier and writes to standard output.
func Warn(args ...interface{}) {
	instance.Warnf(args[0].(string), args[1:]...)
}

// Logs in 'ERROR' level according to a format specifier and writes to standard output.
func Error(args ...interface{}) {
	instance.Errorf(args[0].(string), args[1:]...)
}

// Logs in 'FATAL' level according to a format specifier and writes to standard output.
func Fatal(args ...interface{}) {
	instance.Fatalf(args[0].(string), args[1:]...)
}

// Logs in 'PANIC' level according to a format specifier and panics.
func Panic(args ...interface{}) {
	instance.Panicf(args[0].(string), args[1:]...)
}

// Logger returns a middleware that will write the logs to gin.DefaultWriter.
// By default, gin.DefaultWriter = os.Stdout.
func Logger() []gin.HandlerFunc {
	return []gin.HandlerFunc {ginzap.Ginzap(instance.Desugar(), time.RFC3339, true),
		ginzap.RecoveryWithZap(instance.Desugar(), true)}
}
