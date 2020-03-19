#!/usr/bin/make -f

export CGO_ENABLED=0

PROJECT=previousnext/acquia-cli
VERSION=$(shell git describe --tags --always)
COMMIT=$(shell git rev-list -1 HEAD)

# Builds the project.
build:
	gox -os='linux darwin' \
	    -arch='amd64' \
	    -output='bin/acquia-cli_{{.OS}}_{{.Arch}}' \
	    -ldflags='-extldflags "-static"' \
	    github.com/$(PROJECT)/cmd/acquia-cli

# Run all lint checking with exit codes for CI.
lint:
	golint -set_exit_status `go list ./... | grep -v /vendor/`

# Run tests with coverage reporting.
test:
	go test -cover ./...

IMAGE=${PROJECT}

# Releases the project Docker Hub.
release-docker:
	docker build -t ${IMAGE}:${VERSION} -t ${IMAGE}:latest .
	docker push ${IMAGE}:${VERSION}
	docker push ${IMAGE}:latest

release-github: build
	ghr -u previousnext "${VERSION}" ./bin/

release: release-docker release-github

.PHONY: *
