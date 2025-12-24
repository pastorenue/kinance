package logger

import (
	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
}

func New(level string) *Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})

	switch level {
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	case "info":
		log.SetLevel(logrus.InfoLevel)
	case "warn":
		log.SetLevel(logrus.WarnLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	default:
		log.SetLevel(logrus.InfoLevel)
	}

	return &Logger{log}
}

func (l *Logger) Info(msg string, fields ...interface{}) {
	if len(fields) > 0 {
		l.WithFields(convertFields(fields)).Info(msg)
	} else {
		l.Logger.Info(msg)
	}
}

func (l *Logger) Error(msg string, fields ...interface{}) {
	if len(fields) > 0 {
		l.WithFields(convertFields(fields)).Error(msg)
	} else {
		l.Logger.Error(msg)
	}
}

func (l *Logger) Debug(msg string, fields ...interface{}) {
	if len(fields) > 0 {
		l.WithFields(convertFields(fields)).Debug(msg)
	} else {
		l.Logger.Debug(msg)
	}
}

func convertFields(fields []interface{}) logrus.Fields {
	f := make(logrus.Fields)
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			key := fields[i].(string)
			value := fields[i+1]
			f[key] = value
		}
	}
	return f
}
