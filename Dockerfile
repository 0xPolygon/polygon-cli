# Container image for polycli. Primarily to run `polycli p2p sensor` on
# Container-Optimized OS (no on-host Go toolchain), but usable for any polycli
# subcommand.
FROM golang:1.26 AS build

WORKDIR /src

# Cache modules first.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Reuse the Makefile's build target so the binary is version-stamped exactly like
# a normal `make build` (needs make + git; git describe reads the tags from the
# build context — the publish workflow checks out with fetch-depth: 0). CGO is on
# by default (vectorized-poseidon-gold has no pure-Go path); buildx builds each
# target platform natively under QEMU, so no cross-compilers are needed.
RUN apt-get update \
  && apt-get install -y --no-install-recommends make \
  && rm -rf /var/lib/apt/lists/* \
  && make build

# distroless/base has glibc for the dynamically-linked (CGO) binary.
FROM gcr.io/distroless/base-debian12:nonroot

# The sensor writes a nodes.json discovery cache to its working directory.
WORKDIR /data

COPY --from=build /src/out/polycli /usr/local/bin/polycli

ENTRYPOINT ["/usr/local/bin/polycli"]
