GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test -v --cover
GOGET=$(GOCMD) get
GOVET=$(GOCMD) vet
GORUN=$(GOCMD) run

GIT_COMMIT := $(shell git rev-list -1 HEAD)
BUILT_ON := $(shell hostname)
BUILD_DATE := $(shell date +%FT%T%z)

REPO=lfs-edge

ifndef DOCKER_REPO
  DOCKER_REPO=docker.github.azc.ext.hp.com/krypton
endif

# Publish the image to Github.
publish: docker-image push

tag:
	docker tag $(DOCKER_IMAGE):latest $(DOCKER_REPO)/$(REPO)/$(DOCKER_IMAGE):latest

push: tag
	docker push $(DOCKER_REPO)/$(REPO)/$(DOCKER_IMAGE):latest

vet:
	$(GOVET) ./...

imports:
	goimports -w .
