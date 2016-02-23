# Short name: Short name, following [a-zA-Z_], used all over the place.
# Some uses for short name:
# - Docker image name
# - Kubernetes service, rc, pod, secret, volume names
SHORT_NAME := pkg

# Enable vendor/ directory support.
export GO15VENDOREXPERIMENT=1

# SemVer with build information is defined in the SemVer 2 spec, but Docker
# doesn't allow +, so we use -.
VERSION := 0.0.1-$(shell date "+%Y%m%d%H%M%S")

# Common flags passed into Go's linker.
LDFLAGS := "-s -X main.version=${VERSION}"

NV_PKGS := $(shell glide nv)
GO_PKGS := $(shell glide nv -x)

all: build test

# This builds .a files, which will be placed in $GOPATH/pkg
build:
	go build ${NV_PKGS}

test: test-style
	go test ${NV_PKGS}

test-style:
	@if [ $(shell gofmt -e -l -s ${GO_PKGS}) ]; then \
		echo "gofmt check failed:"; gofmt -e -d -s ${GO_PKGS}; exit 1; \
	fi
	@for i in . ${GO_PKGS}; do \
		golint $$i; \
	done
	@for i in . ${GO_PKGS}; do \
		go vet github.com/deis/pkg/$$i; \
	done

.PHONY: all build test test-style
