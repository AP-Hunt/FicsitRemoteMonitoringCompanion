## IMPORTANT
## This Makefile is intended to be run in a Windows
## environment, and makes use of Windows & Powershell 
## specific things. Don't expect it to work on Linux.

SHELL := powershell.exe
MAKESHELL := powershell.exe
.SHELLFLAGS := -Command

GIT_CHGLOG_VERSION := 0.15.1
SEMVER_CLI_VERSION := 1.1.0
FRMC_VERSION := $(shell Get-Content -Path "./version.txt")

GOBIN := $(shell go env GOBIN)
ifeq ($(GOBIN),)
GOBIN := $(shell go env GOPATH)\bin
endif

.PHONY: release-notes
release-notes: $(GOBIN)\git-chglog.exe
	git-chglog -c "./.chglog/config.yml" -t "./.chglog/CHANGELOG.tpl.md" -o "ReleaseNotes.md" "${FRMC_VERSION}..${FRMC_VERSION}"
	@Get-Content -Path "./InstallationInstructions.md" | Add-Content -Path "./ReleaseNotes.md"

$(GOBIN)\git-chglog.exe:
	go install github.com/git-chglog/git-chglog/cmd/git-chglog@v${GIT_CHGLOG_VERSION}


BUMP=minor
PRERELEASE=
.PHONY: cut-release
cut-release: $(GOBIN)\semver-cli.exe readme
	@if(!@("major", "minor", "patch").Contains("${BUMP}")){ echo "BUMP=major|minor|patch"; exit 1;}
	@echo "Current version is ${FRMC_VERSION}"
	semver-cli inc "${BUMP}" "${FRMC_VERSION}" > version.txt
	if("${PRERELEASE}" -ne "") { semver-cli set prerelease "$$(cat ./version.txt)" "${PRERELEASE}" > version.txt;}
	@echo "New version is $$(cat ./version.txt)"
	@cd map/; npm version "$$(cat ../version.txt)"
	git add "README.md"
	git add "version.txt" "./map/package.json" "./map/package-lock.json"
	git commit -m "Bump version to $$(cat ./version.txt)"
	git tag "$$(cat ./version.txt)"
	@echo "`n`nVersion bump in commit $$(git rev-parse HEAD)"
	@echo "Run the following to push the new version"
	@echo "`t git push origin main"
	@echo "`t git push origin $$(cat ./version.txt)"

$(GOBIN)\semver-cli.exe:
	go install github.com/davidrjonas/semver-cli@${SEMVER_CLI_VERSION}

.PHONY: readme
readme:
	cd Companion/; go run main.go -GenerateReadme > ../README.md
