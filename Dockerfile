# Use an official Go runtime as a parent image
FROM golang:1.21 as builder

# Set the working directory inside the container
WORKDIR /go/src/app

# Copy the Go source code and .git directory into the container
COPY . .

# Build your Go app using the 'build' target in your Makefile
RUN CGO_ENABLED=0 make build

# Use a smaller base image to create a minimal final image
FROM scratch

# Set working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /go/src/app/out/polycli .

# Command to run the binary
CMD ["./polycli"]\