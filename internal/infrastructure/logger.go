package infrastructure

import (
	"io"
	"log/slog"
	"os"
	"path"
)



type AppLogger struct {
	logger *slog.Logger
}

var logWriter io.Writer

func NewAppLogger(config ConfigurationProvider) AppLogger {
	configLogsDir := config.Get("LOGS_DIR", ".")
	logPath := path.Join(configLogsDir, "gear5th-app.log")

	if logWriter == nil {
		f, err := os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
		if err != nil {
			logWriter = os.Stdout
		} else {
			logWriter = f
		}
	}
	return AppLogger{
		slog.New(slog.NewJSONHandler(logWriter, nil)),
	}
}

func (l AppLogger) Error(usecase string, err error) {
	l.logger.Error(err.Error(), "case", usecase)
}

func (l AppLogger) Info(usecase, message string) {
	l.logger.Info(message, "case", usecase)
}
