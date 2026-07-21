# Container image for polycli. Primarily to run `polycli p2p sensor` on
# Container-Optimized OS (no on-host Go toolchain), but usable for any polycli
# subcommand.
#
# The build stage always runs on the native BUILDPLATFORM and cross-compiles to
# the target arch. CGO is required (vectorized-poseidon-gold has no pure-Go
# path), so cross-compiling with a real C toolchain avoids the very slow QEMU
# emulation a per-platform native build would need.
FROM --platform=$BUILDPLATFORM golang:1.26 AS build

# TARGETARCH is provided automatically by buildx (amd64, arm64, ...).
ARG TARGETARCH

WORKDIR /src

# make reuses the Makefile's version stamping; gcc-aarch64-linux-gnu is the same
# cross compiler the release workflow (make cross) uses for arm64. amd64 uses the
# native gcc already in the golang image. NOTE: because vectorized-poseidon-gold
# compiles the amd64 path with -march=native, the amd64 image must be built on an
# amd64 host (true on the ubuntu-latest CI runner); arm64 cross-compiles cleanly.
RUN apt-get update \
  && apt-get install -y --no-install-recommends make gcc-aarch64-linux-gnu \
  && rm -rf /var/lib/apt/lists/*

# Cache modules first (same cache mount as the build step so they aren't
# downloaded twice on a cold build).
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY . .

# Cross-compile a static, version-stamped binary. git describe reads the tags
# from the build context — the publish workflow checks out with fetch-depth: 0.
# Build cache mounts keep re-tag builds fast.
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    GOARCH="$TARGETARCH" \
    CC="$([ "$TARGETARCH" = arm64 ] && echo aarch64-linux-gnu-gcc || echo gcc)" \
    make docker-build

# Staged empty dir so the final image's working directory is owned by the
# nonroot user (distroless has no shell to chown it there).
RUN mkdir -p /data

# The binary is fully static, so distroless/static (no glibc) is enough.
FROM gcr.io/distroless/static-debian12:nonroot

# The sensor writes a nodes.json discovery cache to its working directory, so it
# must be owned by the nonroot runtime user (uid 65532).
COPY --from=build --chown=65532:65532 /data /data
WORKDIR /data

COPY --from=build /src/out/polycli /usr/local/bin/polycli

ENTRYPOINT ["/usr/local/bin/polycli"]
