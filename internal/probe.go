package internal

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/Masterminds/sprig/v3"
	"github.com/phramz/go-healthcheck/pkg/contract"
	"github.com/phramz/go-healthcheck/pkg/http"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var _ contract.Probe = (*probe)(nil)
var _ assert.TestingT = (*probe)(nil)

func NewProbe(name, tpl string, logger logrus.FieldLogger, config contract.Config) (contract.Probe, error) {
	functions := sprig.TxtFuncMap()
	functions["regexp"] = regexp.MustCompile

	t, err := template.New(name).Funcs(functions).Parse(tpl)
	if err != nil {
		return nil, err
	}

	return &probe{
		name:    name,
		tpl:     t,
		config:  config,
		logger:  logger,
		results: make([]contract.Result, 0),
	}, nil
}

type probe struct {
	config   contract.Config
	logger   logrus.FieldLogger
	tpl      *template.Template
	results  []contract.Result
	name     string
	started  *time.Time
	finished *time.Time
}

func (p *probe) HasFailed() bool {
	for _, r := range p.results {
		if r.Failed() {
			return true
		}
	}

	return false
}

func (p *probe) HasPassed() bool {
	return !p.HasFailed()
}

func (p *probe) IsDone() bool {
	return p.finished != nil
}

func (p *probe) clone() *probe {
	i := &probe{
		tpl:     p.tpl,
		config:  p.config,
		logger:  p.logger,
		results: p.results,
	}

	return i
}

func (p *probe) Run() (string, error) {
	if p.started != nil {
		panic("already started")
	}

	ts := time.Now()
	p.started = &ts

	logger := p.logger.WithField(LogFieldStatus, "STARTED")
	logger.Info(`starting probe`)

	defer func() {
		tf := time.Now()
		p.finished = &tf

		if p.HasFailed() {
			logger.WithField(LogFieldStatus, "FAIL").Warning(`finished probe`)
			return
		}

		logger.WithField(LogFieldStatus, "PASS").Info(`finished probe`)
	}()

	param := struct {
		Assert *assert.Assertions
		HTTP   contract.HealthCheck
	}{
		Assert: assert.New(p),
		HTTP:   http.NewHttpHealthCheck(p).WithLogger(logger.WithField(LogFieldComponent, "http")),
	}

	var b bytes.Buffer
	out := bufio.NewWriter(&b)

	err := p.tpl.Execute(out, param)

	if err != nil {
		p.logger.Error(err)
		return "", err
	}

	rendered := b.String()

	return rendered, err
}

func (p *probe) Fail(err error) {
	p.logger.Warning(err.Error())
	p.results = append(p.results, NewFail(err))
}

func (p *probe) Failf(format string, args ...interface{}) {
	p.Errorf(format, args...)
}

func (p *probe) Infof(format string, args ...interface{}) {
	p.logger.Infof(format, args...)
}

func (p *probe) Errorf(format string, args ...interface{}) {
	// remove stack
	if len(args) == 1 {
		if str, ok := args[0].(string); ok && strings.HasPrefix(str, "\tError Trace:") {
			start := false
			args[0] = ""
			for _, l := range strings.Split(str, "\n") {
				if strings.HasPrefix(strings.TrimSpace(l), "Error:") {
					start = true
				}

				if !start {
					continue
				}

				args[0] = fmt.Sprintf("%s\n%s", args[0], l)
			}
		}
	}

	p.logger.Warningf(format, args...)
	p.results = append(p.results, NewFailf(format, args...))
}

func (p *probe) WithConfig(config contract.Config) contract.Probe {
	i := p.clone()
	i.config = config

	return i
}

func (p *probe) WithLogger(logger logrus.FieldLogger) contract.Probe {
	i := p.clone()
	i.logger = logger

	return i
}
