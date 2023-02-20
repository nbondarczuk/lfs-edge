FROM docker.github.azc.ext.hp.com/krypton/images/krypton-go-builder as builder

ADD . /go/src/lfs-edge
WORKDIR /go/src/lfs-edge

# build the source
RUN make -C proto build-proto
RUN make -C proxy build-binaries

# use a minimal alpine image
FROM docker.github.azc.ext.hp.com/krypton/images/krypton-go-base

# set working directory
WORKDIR /go/bin

COPY --from=builder /go/src/lfs-edge/bin/proxy/lfs-edge-proxy .
COPY --from=builder /go/src/lfs-edge/bin/proxy/config.yml .

USER 1001
EXPOSE 8181/tcp

# run the binary
CMD ["./lfs-edge-proxy"]
