package contract

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type Probe interface {
	assert.TestingT

	HasFailed() bool
	HasPassed() bool
	IsDone() bool
	Run() (string, error)
	Fail(err error)
	Failf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	WithConfig(config Config) Probe
	WithLogger(logger logrus.FieldLogger) Probe
}
