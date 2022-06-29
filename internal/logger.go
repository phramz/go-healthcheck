package internal

import (
	"github.com/phramz/go-healthcheck/pkg/contract"
	"github.com/sirupsen/logrus"
)

const (
	LogFieldStatus    = "status"
	LogFieldComponent = "component"
	LogFieldProbe     = "probe"
	LogFieldProbeKind = "probe"

	LogLevelFatal   = "fatal"
	LogLevelError   = "error"
	LogLevelWarning = "warning"
	LogLevelInfo    = "info"
	LogLevelDebug   = "debug"
)

func newLogger(container contract.Ecosystem) logrus.FieldLogger {
	log := logrus.New()

	log.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "15:04:05",
		FullTimestamp:   true,
		PadLevelText:    true,
	})

	logLevel := container.Config().String(contract.ConfigKeyLogLevel, LogLevelInfo)
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		log.Fatalf("invalid log level %q: %v", logLevel, err)
	}

	log.Debugf("log level set to %q", logLevel)
	log.SetLevel(level)

	return log
}
