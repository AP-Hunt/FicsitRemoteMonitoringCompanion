//go:build !windows
// +build !windows

package prometheus

import (
	"errors"
	"os"
	"os/exec"
)

type PrometheusWrapper struct {
	cmd    *exec.Cmd
	stdout os.File
	stderr os.File
}

func NewPrometheusWrapper() (*PrometheusWrapper, error) {
	return nil, errors.New("Not Implemented")
}

func (pw *PrometheusWrapper) Start() error {
	return errors.New("Not Implemented")
}

func (pw *PrometheusWrapper) Stop() error {
	return errors.New("Not Implemented")
}
