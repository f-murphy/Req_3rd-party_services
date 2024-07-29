package logger

import (
	"github.com/sirupsen/logrus"
	"os"
	"io"
)

func InitLogger() (*os.File, error) {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	logFile, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)
	logrus.SetOutput(multiWriter)

	return logFile, err
}
