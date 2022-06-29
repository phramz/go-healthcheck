package contract

import "github.com/sirupsen/logrus"

type HealthCheck interface {
	WithLogger(logger logrus.FieldLogger) HealthCheck
	WithProbe(probe Probe) HealthCheck
}
