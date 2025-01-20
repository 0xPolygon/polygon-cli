# Build stage
FROM --platform=${BUILDPLATFORM} golang:1.22 AS builder

# Set the workspace for the build
WORKDIR /workspace

# Copy only necessary Go module files for caching dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the necessary source code
COPY . ./

# Build the Go binary
RUN go build -o /workspace/polycli main.go

# Final stage: minimal base image
FROM --platform=${BUILDPLATFORM} debian:bookworm-slim

# Copy only the necessary files from the builder image
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /workspace/polycli /usr/bin/polycli

# Default cmd for the container
CMD ["/bin/sh", "-c", "polycli"]
ENTRYPOINT ["polycli"]