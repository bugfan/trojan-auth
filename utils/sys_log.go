package utils

import (
	"github.com/bugfan/trojan-auth/env"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	if env.Get("close_log") == "true" {
		logrus.SetOutput(&EmptyLogger{})
	}
}

type EmptyLogger struct{}

func (e *EmptyLogger) Write(data []byte) (int, error) {
	return 0, nil
}
