FROM golang:1.21 as builder
WORKDIR /go/src/app
COPY go.mod go.sum ./
RUN go mod download
COPY abi bindings cmd dashboard gethkeystore hdwallet metrics p2p proto rpctypes util ./
RUN CGO_ENABLED=0 make build

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /go/src/app/out/polycli /usr/bin/polycli
USER 65532:65532
ENTRYPOINT ["polycli"]
CMD ["--help"]
