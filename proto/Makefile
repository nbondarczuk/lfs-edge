PROTOS_DIR=.
PROTOC_PATH=/usr/local/bin
PROTOC_CMD=protoc
PROTOC_BUILD=$(PROTOC_PATH)/$(PROTOC_CMD)
BUILDER_IMAGE=docker.github.azc.ext.hp.com/krypton/images/krypton-go-builder

# run build in a builder container so protoc install details
# are not dependent on local machines
build:
	docker run --rm \
		-v$(CURDIR):/src \
		-w/src \
		$(BUILDER_IMAGE) \
	make build-proto

# Run protoc to generate go code for .proto files.
build-proto:
	$(PROTOC_BUILD) -I $(PROTOS_DIR) \
	--go_out=plugins=grpc,paths=source_relative:$(PROTOS_DIR) \
	$(PROTOS_DIR)/file_proxy_ipc_common.proto \
	$(PROTOS_DIR)/file_proxy_ipc.proto

# Run clean to remove all .pb.go files from source code.
clean:
	-rm -f *.pb.go

.PHONY: protoc clean
