package http

import (
	"net/http"
	"time"

	"github.com/phramz/go-healthcheck/pkg/contract"
	"github.com/sirupsen/logrus"
)

var _ contract.HealthCheck = (*httpHealthcheck)(nil)

func NewHttpHealthCheck(probe contract.Probe) contract.HealthCheck {
	return (&httpHealthcheck{}).WithProbe(probe)
}

type httpHealthcheck struct {
	logger logrus.FieldLogger
	probe  contract.Probe
}

func (h *httpHealthcheck) clone() *httpHealthcheck {
	i := &httpHealthcheck{
		logger: h.logger,
		probe:  h.probe,
	}

	return i
}

func (h *httpHealthcheck) WithProbe(probe contract.Probe) contract.HealthCheck {
	i := h.clone()
	i.probe = probe

	return i
}

func (h *httpHealthcheck) WithLogger(logger logrus.FieldLogger) contract.HealthCheck {
	i := h.clone()
	i.logger = logger

	return i
}

func (h *httpHealthcheck) Request(method, url string) http.Response {
	rq, err := http.NewRequest(method, url, nil)
	if err != nil {
		h.probe.Fail(err)
		return http.Response{}
	}

	resp, err := h.Client().Do(rq)
	if err != nil {
		h.probe.Fail(err)
		return http.Response{}
	}

	return *resp
}

func (h *httpHealthcheck) Handler() http.HandlerFunc {
	return func(rw http.ResponseWriter, rq *http.Request) {
		client := h.Client()

		resp, err := client.Do(rq)
		if err != nil {
			h.probe.Fail(err)
			return
		}

		err = resp.Write(rw)
		if err != nil {
			h.probe.Fail(err)
			return
		}
	}
}

func (h *httpHealthcheck) Client() *http.Client {
	client := &http.Client{
		Timeout: time.Duration(30) * time.Second,
	}

	return client
}
