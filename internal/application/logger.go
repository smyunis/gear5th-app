package application

type Logger interface {
	Error(usecase string, err error)
	Info(usecase, message string)
}