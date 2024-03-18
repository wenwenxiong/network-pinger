SHELL = /bin/bash

GOLDFLAGS = "-w -s -extldflags '-z now' -X github.com/kubeovn/kube-ovn/versions.COMMIT=$(COMMIT) -X github.com/kubeovn/kube-ovn/versions.VERSION=$(RELEASE_TAG) -X github.com/kubeovn/kube-ovn/versions.BUILDDATE=$(DATE)"

REGISTRY = kubesphere
RELEASE_TAG = $(shell cat VERSION)
# ARCH could be amd64,arm64
ARCH = amd64

.PHONY: build-go
build-go:
	go mod tidy
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -buildmode=pie -o $(CURDIR)/dist/images/network-pinger -ldflags $(GOLDFLAGS) -v ./cmd

.PHONY: image-network-pinger
image-network-pinger: build-go
	docker buildx build --platform linux/amd64 -t $(REGISTRY)/network-pinger:$(RELEASE_TAG) --build-arg VERSION=$(RELEASE_TAG) -o type=docker -f dist/images/Dockerfile dist/images/
