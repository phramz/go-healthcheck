package internal

import (
	"fmt"

	"github.com/phramz/go-healthcheck/pkg/contract"
)

var _ contract.Result = (*result)(nil)

func NewResult(err error) contract.Result {
	return &result{
		pass: err == nil,
		err:  err,
	}
}

func NewPass() contract.Result {
	return NewResult(nil)
}

func NewFail(err error) contract.Result {
	return NewResult(err)
}

func NewFailf(format string, args ...interface{}) contract.Result {
	return NewResult(fmt.Errorf(format, args...))
}

type result struct {
	pass bool
	err  error
}

func (r *result) Passed() bool {
	return r.pass
}

func (r *result) Failed() bool {
	return !r.Passed()
}

func (r *result) Reason() string {
	return r.Error().Error()
}

func (r *result) Error() error {
	return r.err
}
