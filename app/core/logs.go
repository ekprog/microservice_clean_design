package core

type Logger interface {
	Debug(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Fatal(msg string, args ...interface{})

	DebugWrap(err error, msg string, args ...interface{})
	WarnWrap(err error, msg string, args ...interface{})
	InfoWrap(err error, msg string, args ...interface{})
	ErrorWrap(err error, msg string, args ...interface{})
	FatalWrap(err error, msg string, args ...interface{})
}
