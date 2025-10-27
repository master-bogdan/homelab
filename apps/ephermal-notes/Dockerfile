# Use the official Golang image matching go.mod
FROM golang:1.24-alpine AS builder

# Set environment variables for static build
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set working directory inside container
WORKDIR /build

# Install swag CLI for Swagger doc generation
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy go.mod and go.sum to cache dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire source code
COPY . .

# Generate Swagger docs automatically
RUN /go/bin/swag init -g cmd/server/main.go -o pkg/docs

# Build statically linked binary
RUN go build -trimpath -ldflags="-s -w" -o /app/server ./cmd/server

# Use a minimal Alpine image for the final stage
FROM alpine:3.21

RUN apk add --no-cache ca-certificates

# Copy the compiled binary from the builder stage
COPY --from=builder /app/server /app/server
COPY --from=builder /build/pkg/docs /pkg/docs

# Expose the port your app listens on (adjust as needed)
EXPOSE 8000

# Run the application
CMD ["/app/server"]
