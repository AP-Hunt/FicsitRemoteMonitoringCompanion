//go:build windows
// +build windows

package prometheus

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"syscall"
)

type PrometheusWrapper struct {
	cmd    *exec.Cmd
	stdout os.File
	stderr os.File
}

func NewPrometheusWrapper() (*PrometheusWrapper, error) {

	curExePath, err := os.Executable()
	if err != nil {
		return nil, err
	}

	curExeDir := filepath.Dir(curExePath)

	if err != nil {
		return nil, err
	}

	promPath := path.Join(curExeDir, "prometheus", "prometheus.exe")
	cfgPath := path.Join(curExeDir, "prometheus", "prometheus.yml")

	stdout, err := os.Create(path.Join(curExeDir, "prometheus.stdout.log"))
	if err != nil {
		return nil, err
	}

	stderr, err := os.Create(path.Join(curExeDir, "prometheus.stderr.log"))
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(promPath, "--config.file", cfgPath)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x08000000} // CREATE_NO_WINDOW

	return &PrometheusWrapper{
		cmd:    cmd,
		stdout: *stdout,
		stderr: *stderr,
	}, nil
}

func (pw *PrometheusWrapper) Start() error {
	return pw.cmd.Start()
}

func (pw *PrometheusWrapper) Stop() error {
	pw.stdout.Close()
	pw.stderr.Close()
	err := pw.cmd.Process.Kill()
	pw.cmd.Wait()
	return err
}
