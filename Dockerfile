# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH environment variable at /go.
FROM golang:1.22 as builder

# Create and change to the app directory.
WORKDIR /app

# Copy go/templ source files.
COPY go.mod go.sum ./
COPY *.go ./
COPY logging/ logging/
COPY types/ types/
COPY templates/ templates/

# Fetch dependencies.
RUN go mod download

# Build the binary.
RUN go build -v -o server

# Use a Docker multi-stage build to create a lean production image.
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM gcr.io/distroless/base-debian12

# Copy the compiled application from the builder stage.
COPY --from=builder /app/server /

# Your application listens on port 8080.
EXPOSE 8080

# Run your application.
CMD ["/server"]
