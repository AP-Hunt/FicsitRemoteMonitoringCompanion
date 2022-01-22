## IMPORTANT
## This Makefile is intended to be run in a Windows
## environment, and makes use of Windows & Powershell 
## specific things. Don't expect it to work on Linux.

SHELL := powershell.exe
MAKESHELL := powershell.exe

.SHELLFLAGS := -Command
.DEFAULT_GOAL := bin/Companion.exe

GRAFANA_VERSION := 7.5.7
PROMETHEUS_VERSION := 2.27.1
FRMC_VERSION := $(shell Get-Content -Path "../version.txt")

GO_SRC := $(shell Get-ChildItem -Recurse -Filter "*.go" -Exclude "*test.go" -Name | where {$$_ -notlike "vendor*" })
GO_BUILD_ENV := $$Env:GOOS = "windows"; $$Env:GOARCH = "amd64";


.PHONY: companion
bin/Companion.exe: $(GO_SRC) bin/prometheus/prometheus.exe bin/prometheus/prometheus.yml bin/map/ vendor/
	${GO_BUILD_ENV} go build -ldflags="-X 'main.Version=${FRMC_VERSION}'" -o bin/Companion.exe

bin/prometheus/prometheus.exe:
	@New-Item -Path "../Externals/Prometheus/" -ItemType Directory -Force > $$null
	@Invoke-WebRequest "https://github.com/prometheus/prometheus/releases/download/v${PROMETHEUS_VERSION}/prometheus-${PROMETHEUS_VERSION}.windows-amd64.zip" -OutFile "../Externals/Prometheus/prometheus.zip" > $$null
	@Expand-Archive -Path "../Externals/Prometheus/prometheus.zip" -DestinationPath "./bin/" -Force
	@Rename-Item -Path "./bin/prometheus-${PROMETHEUS_VERSION}.windows-amd64" -NewName "prometheus"

bin/prometheus/prometheus.yml: FORCE
	@Copy-Item -Path "./prometheus.yml" -Destination "bin/prometheus/prometheus.yml" -Force
FORCE:

bin/map/:
	@New-Item -Path "./bin/map" -ItemType Directory -Force > $$null
	@cd ../map/; npm install; npm run compile
	@cd ../map/; Copy-Item -Path "index.html","map-16k.png","vendor","img","js" -Destination "../Companion/bin/map/" -Recurse -Force

vendor/: go.mod go.sum
	@go mod vendor

.PHONY: test
test: 
	go test ./...

.PHONY: clean
clean: 
	Remove-Item -Recurse -Force "./bin"

.PHONY: release
release: bin/Companion.exe
	Compress-Archive -Path "./bin/*" -DestinationPath "./FicsitRemoteMonitoringCompanion-v${FRMC_VERSION}.zip"