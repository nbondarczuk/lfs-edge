BINARY_DIR=../bin/proxy
BINARY_NAME=lfs-edge-proxy
DOCKER_IMAGE=lfs-edge-proxy

all: build docker-image

build: build-proto build-binaries

# Build the binaries for the service.
build-binaries: 
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) \
	-ldflags "-X main.gitCommitHash=$(GIT_COMMIT) -X main.builtAt=$(BUILD_DATE) \
	-X main.builtBy=$(USER) -X main.builtOn=$(BUILT_ON)" \
	-o $(BINARY_DIR)/$(BINARY_NAME) -v
	cp -r config/*.yml $(BINARY_DIR)

# Create a docker image for the service.
docker-image:
	make -C docker

# Make proto files if never done
build-proto:
	make -C ../proto build

# Start the service and its dependencies within a local Docker network.
start:

run:
	$(GORUN) ./... \
		--config_file=config/config.yml

unit-test:
	$(GOTEST) ./...

clean:
	$(GOCLEAN)
	make -C docker clean
	@rm -rf $(BINARY_DIR)

.SILENT:
.PHONY: all build build-binaries docker-image start vet clean

include ../common/common.mk
