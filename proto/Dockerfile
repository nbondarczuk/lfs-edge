FROM golang:1.19.4-alpine3.17 as stage

ENV PROTOBUF_URL https://github.com/protocolbuffers/protobuf/releases/download/v3.15.8/protoc-3.15.8-linux-x86_64.zip
RUN apk add curl && \
    curl -sLo /tmp/protobuf.zip ${PROTOBUF_URL} && \
    unzip /tmp/protobuf.zip -d /tmp && \
    go install github.com/golang/protobuf/protoc-gen-go@latest

## final
FROM golang:1.19.4-alpine3.17
RUN apk add --no-cache git alpine-sdk gcompat
COPY --from=stage /go/bin/protoc-gen-go /usr/local/bin
COPY --from=stage /tmp/bin/protoc /usr/local/bin
COPY --from=stage /tmp/include /usr/local/include
