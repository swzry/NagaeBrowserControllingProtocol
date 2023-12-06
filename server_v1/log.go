package server_v1

type IBaseLogger interface {
	Error(d ...interface{})
	ErrorF(format string, d ...interface{})
	ErrorC(data map[string]string)

	Warn(d ...interface{})
	WarnF(format string, d ...interface{})
	WarnC(data map[string]string)

	Info(d ...interface{})
	InfoF(format string, d ...interface{})
	InfoC(data map[string]string)

	Verbose(d ...interface{})
	VerboseF(format string, d ...interface{})
	VerboseC(data map[string]string)

	Debug(d ...interface{})
	DebugF(format string, d ...interface{})
	DebugC(data map[string]string)
}

type IServerLogger interface {
	IBaseLogger
	GetSessionLogger(sessionID int64) ISessionLogger
	DiscardSessionLogger(sessionID int64)
}

type ISessionLogger interface {
	IBaseLogger
}

var _ IServerLogger = (*EmptyLogger)(nil)
var _ ISessionLogger = (*EmptyLogger)(nil)

type EmptyLogger struct {
}

func (e EmptyLogger) DiscardSessionLogger(sessionID int64) {
}

func (e EmptyLogger) Error(d ...interface{}) {
}

func (e EmptyLogger) ErrorF(format string, d ...interface{}) {
}

func (e EmptyLogger) ErrorC(data map[string]string) {
}

func (e EmptyLogger) Warn(d ...interface{}) {
}

func (e EmptyLogger) WarnF(format string, d ...interface{}) {
}

func (e EmptyLogger) WarnC(data map[string]string) {
}

func (e EmptyLogger) Info(d ...interface{}) {
}

func (e EmptyLogger) InfoF(format string, d ...interface{}) {
}

func (e EmptyLogger) InfoC(data map[string]string) {
}

func (e EmptyLogger) Verbose(d ...interface{}) {
}

func (e EmptyLogger) VerboseF(format string, d ...interface{}) {
}

func (e EmptyLogger) VerboseC(data map[string]string) {
}

func (e EmptyLogger) Debug(d ...interface{}) {
}

func (e EmptyLogger) DebugF(format string, d ...interface{}) {
}

func (e EmptyLogger) DebugC(data map[string]string) {
}

func (e EmptyLogger) GetSessionLogger(sessionID int64) ISessionLogger {
	return &EmptyLogger{}
}
