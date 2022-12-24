package log

type ILogger interface {
	Info(msg string)
	Error(msg string, err error)
	Fatal(msg string, err error)
}
