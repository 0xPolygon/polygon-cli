FROM golang:1.21 AS builder
WORKDIR /workspace
COPY go.mod go.sum ./
RUN go mod download

COPY abi/ abi/
COPY bindings/ bindings/
COPY cmd/ cmd/
COPY dashboard/ dashboard/
COPY gethkeystore/ gethkeystore/
COPY hdwallet/ hdwallet/
COPY metrics/ metrics/
COPY p2p/ p2p/
COPY proto/ proto/
COPY rpctypes/ rpctypes/
COPY util/ util/
COPY main.go ./
RUN CGO_ENABLED=0 go build -o polycli main.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /workspace/polycli /usr/bin/polycli
USER 65532:65532
ENTRYPOINT ["polycli"]
CMD ["--help"]

# How to test this image?
# https://github.com/maticnetwork/polygon-cli/pull/189#discussion_r1464486344
