# Container image for polycli. Primarily to run `polycli p2p sensor` on
# Container-Optimized OS (no on-host Go toolchain), but usable for any polycli
# subcommand.
FROM golang:1.26 AS build

WORKDIR /src

# Cache modules first.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# CGO is required (e.g. vectorized-poseidon-gold has no pure-Go path). The
# golang image ships a C toolchain; buildx builds each target platform natively
# under QEMU, so this works multi-arch without cross-compilers.
RUN CGO_ENABLED=1 go build -trimpath -o /polycli main.go

# distroless/base has glibc for the dynamically-linked (CGO) binary.
FROM gcr.io/distroless/base-debian12:nonroot

# The sensor writes a nodes.json discovery cache to its working directory.
WORKDIR /data

COPY --from=build /polycli /usr/local/bin/polycli

ENTRYPOINT ["/usr/local/bin/polycli"]
