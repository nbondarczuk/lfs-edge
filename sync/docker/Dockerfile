FROM docker.github.azc.ext.hp.com/krypton/images/krypton-go-builder as builder

ADD . /go/src/lfs-edge
WORKDIR /go/src/lfs-edge

# build the source
RUN make -C proto build-proto
RUN make -C sync build-binaries

# use a minimal alpine image
FROM docker.github.azc.ext.hp.com/krypton/images/krypton-go-base

# set working directory
WORKDIR /go/bin

COPY --from=builder /go/src/lfs-edge/bin/sync/lfs-edge-sync .
COPY --from=builder /go/src/lfs-edge/bin/sync/config/config.yml .
COPY --from=builder /go/src/lfs-edge/bin/sync/schema /go/bin/schema/

USER 1001
EXPOSE 8383/tcp

# run the binary
CMD ["./lfs-edge-sync"]
