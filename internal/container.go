package internal

import (
	"github.com/phramz/go-healthcheck/pkg/contract"
	"github.com/sirupsen/logrus"
)

var _ contract.Ecosystem = (*ecosystem)(nil)

func NewEcosystem() contract.Ecosystem {
	return &ecosystem{}
}

type ecosystem struct {
	config contract.Config
	logger logrus.FieldLogger
}

func (s *ecosystem) clone() *ecosystem {
	i := &ecosystem{
		config: s.config,
		logger: s.logger,
	}

	return i
}

func (s *ecosystem) Config() contract.Config {
	if s.config == nil {
		s.config = NewConfig()
	}

	return s.config
}

func (s *ecosystem) Logger() logrus.FieldLogger {
	if s.logger == nil {
		s.logger = newLogger(s)
	}

	return s.logger
}

func (s *ecosystem) WithConfig(config contract.Config) contract.Ecosystem {
	i := s.clone()
	i.config = config

	return i
}
