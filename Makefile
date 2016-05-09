SHORT_NAME ?= deisrel

LDFLAGS := "-s -X main.version=${VERSION}"

REPO_PATH := github.com/deis/${SHORT_NAME}
DEV_ENV_IMAGE := quay.io/deis/go-dev:0.9.1
DEV_ENV_WORK_DIR := /go/src/${REPO_PATH}
DEV_ENV_PREFIX := docker run --rm -e GO15VENDOREXPERIMENT=1 -v ${CURDIR}:${DEV_ENV_WORK_DIR} -w ${DEV_ENV_WORK_DIR}
DEV_ENV_CMD := ${DEV_ENV_PREFIX} ${DEV_ENV_IMAGE}

VERSION ?= $(shell git rev-parse --short HEAD)

DEIS_BINARY_NAME ?= ./deis

bootstrap:
	${DEV_ENV_CMD} glide install

build:
	${DEV_ENV_PREFIX} -e CGO_ENABLED=0 ${DEV_ENV_IMAGE} go build

test:
	${DEV_ENV_CMD} sh -c 'go test -v $$(glide nv)'

build-cli-cross:
	${DEV_ENV_CMD} gox -output="bin/${SHORT_NAME}-{{.OS}}-{{.Arch}}"

build-cli:
	go build -o deisrel ./actions

prep-bintray-json:
	@jq '.version.name |= "$(VERSION)"' _scripts/ci/bintray-template.json \
		> _scripts/ci/bintray-ci.json

dist: build-cli-cross
