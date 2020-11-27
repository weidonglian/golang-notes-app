package logging

import (
	"github.com/sirupsen/logrus"
)

func NewLogger() *logrus.Logger {
	// Setup the logger backend using sirupsen/logrus and configure
	// it to use a custom JSONFormatter. See the logrus docs for how to
	// configure the backend at github.com/sirupsen/logrus
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		// FullTimestamp: true,
		DisableTimestamp: true,
	})
	/* logger.Formatter = &logrus.JSONFormatter{
		// disable, as we set our own
		DisableTimestamp: true,
	}*/
	return logger
}
