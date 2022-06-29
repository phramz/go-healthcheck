package contract

import (
	"github.com/sirupsen/logrus"
)

type Ecosystem interface {
	Config() Config
	Logger() logrus.FieldLogger
	WithConfig(config Config) Ecosystem
}
