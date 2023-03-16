package app

import (
	"fmt"
	logrus "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"microservice/app/core"
	"microservice/app/logs_hooks"
)

var (
	log core.Logger
)

func InitLogs(rootDir ...string) (core.Logger, error) {

	basePath := "."
	if len(rootDir) != 0 {
		basePath = rootDir[0]
	}

	//log.SetFormatter(&easy.Formatter{
	//	TimestampFormat: "2006-01-02 15:04:05",
	//	LogFormat:       "[%lvl%]: %time% - %msg%\n",
	//})

	logrusLogger := logrus.New()
	logrusLogger.SetFormatter(&logrus.TextFormatter{
		ForceColors:      true,
		DisableTimestamp: true,
	})
	logrusLogger.SetReportCaller(true)
	if viper.GetString("app.debug") == "true" {
		logrusLogger.SetLevel(logrus.TraceLevel)
	} else {
		logrusLogger.SetLevel(logrus.InfoLevel)
	}

	// hooks
	logrusLogger.AddHook(logs_hooks.NewToFileHook(basePath))
	logrusLogger.AddHook(logs_hooks.NewToFileErrorHook(basePath))

	log = NewDefaultLogger(logrusLogger)
	return log, nil
}

type DefaultLogger struct {
	logger *logrus.Logger
}

func NewDefaultLogger(logger *logrus.Logger) *DefaultLogger {
	return &DefaultLogger{logger: logger}
}

func (l *DefaultLogger) Debug(msg string, args ...interface{}) {
	l.logger.Debugf(msg, args...)
}

func (l *DefaultLogger) Warn(msg string, args ...interface{}) {
	l.logger.Warnf(msg, args...)
}

func (l *DefaultLogger) Info(msg string, args ...interface{}) {
	l.logger.Infof(msg, args...)
}

func (l *DefaultLogger) Error(msg string, args ...interface{}) {
	l.logger.Errorf(msg, args...)
}

func (l *DefaultLogger) Fatal(msg string, args ...interface{}) {
	l.logger.Fatalf(msg, args...)
}

func (l *DefaultLogger) DebugWrap(err error, msg string, args ...interface{}) {
	msg = fmt.Sprintf(msg, args...)
	l.logger.Debugf("%s: %s", msg, err.Error())
}

func (l *DefaultLogger) WarnWrap(err error, msg string, args ...interface{}) {
	msg = fmt.Sprintf(msg, args...)
	l.logger.Warnf("%s: %s", msg, err.Error())
}

func (l *DefaultLogger) InfoWrap(err error, msg string, args ...interface{}) {
	msg = fmt.Sprintf(msg, args...)
	l.logger.Infof("%s: %s", msg, err.Error())
}

func (l *DefaultLogger) ErrorWrap(err error, msg string, args ...interface{}) {
	msg = fmt.Sprintf(msg, args...)
	l.logger.Errorf("%s: %s", msg, err.Error())
}

func (l *DefaultLogger) FatalWrap(err error, msg string, args ...interface{}) {
	msg = fmt.Sprintf(msg, args...)
	l.logger.Fatalf("%s: %s", msg, err.Error())
}
