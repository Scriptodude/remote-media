package log

import "github.com/sirupsen/logrus"

func GetLoggerForHandler(handler string) *logrus.Entry {
	log := logrus.New()

	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	return log.WithFields(logrus.Fields{"handler": handler})
}
