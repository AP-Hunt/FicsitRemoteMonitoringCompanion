## IMPORTANT
## This Makefile is intended to be run in a Windows
## environment, and makes use of Windows & Powershell 
## specific things. Don't expect it to work on Linux.

SHELL := powershell.exe
.SHELLFLAGS := -Command
.ONESHELL:

GRAFANA_VERSION := 7.5.7
PROMETHEUS_VERSION := 2.27.1

.PHONY: companion
companion: bin/prometheus/prometheus.exe bin/prometheus/prometheus.yml
	go build -o bin/Companion.exe

bin/prometheus/prometheus.exe:
	@New-Item -Path "../Externals/Prometheus/" -ItemType Directory -Force > $$null
	@Invoke-WebRequest "https://github.com/prometheus/prometheus/releases/download/v${PROMETHEUS_VERSION}/prometheus-${PROMETHEUS_VERSION}.windows-amd64.zip" -OutFile "../Externals/Prometheus/prometheus.zip" > $$null
	@Expand-Archive -Path "../Externals/Prometheus/prometheus.zip" -DestinationPath "./bin/" -Force
	@Rename-Item -Path "./bin/prometheus-${PROMETHEUS_VERSION}.windows-amd64" -NewName "prometheus"

bin/prometheus/prometheus.yml:
	@Copy-Item -Path "prometheus.yml" -Destination "bin/prometheus/prometheus.yml"