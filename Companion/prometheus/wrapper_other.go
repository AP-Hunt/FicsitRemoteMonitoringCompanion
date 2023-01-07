//go:build !windows
// +build !windows

package prometheus

import (
	"os"
	"os/exec"
)

// a no-op prometheus wrapper. Not implemented.
type PrometheusWrapper struct {
	cmd    *exec.Cmd
	stdout os.File
	stderr os.File
}

func NewPrometheusWrapper() (*PrometheusWrapper, error) {
	return nil, nil
}

func (pw *PrometheusWrapper) Start() error {
	return nil
}

func (pw *PrometheusWrapper) Stop() error {
	return nil
}
