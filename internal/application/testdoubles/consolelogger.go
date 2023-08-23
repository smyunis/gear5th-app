package testdoubles

import "log/slog"


type ConsoleLogger struct {
}


func (l ConsoleLogger) Error(usecase string, err error) {
	slog.Error(err.Error(), "case", usecase)
}

func (l ConsoleLogger) Info(usecase, message string) {
	slog.Info(message, "case", usecase)
}