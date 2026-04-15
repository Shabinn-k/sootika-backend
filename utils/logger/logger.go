package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func InitLogger() *logrus.Logger {
	Log = logrus.New()

	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	if err := os.MkdirAll("logs", 0755); err != nil {
		panic(err)
	}

	file, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	Log.SetOutput(io.MultiWriter(os.Stdout, file))
	Log.SetLevel(logrus.InfoLevel)
	Log.Info("Logger initialised")

	return Log
}
