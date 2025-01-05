# Stage 1: Build
FROM golang:1.23.4 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY epg-approved/go.mod epg-approved/go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application source code
COPY epg-approved/ ./

# Build the application (statically linked to avoid glibc issues)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o epg-approved ./cmd/epg-approved/main.go

# Stage 2: Final Image
FROM alpine:latest

# Set up a working directory
WORKDIR /app

# Install necessary tools
RUN apk add --no-cache bash curl

# Copy the binary from the builder stage
COPY --from=builder /app/epg-approved .

# Copy the config file
COPY epg-approved/config.toml .

# Copy the wait-for-it.sh script
COPY wait-for-it.sh /app/wait-for-it.sh

# Make wait-for-it.sh executable
RUN chmod +x /app/wait-for-it.sh

# Expose necessary port
EXPOSE 8080

# Command to run the application with wait-for-it.sh
ENTRYPOINT ["/app/wait-for-it.sh", "iptv-search:8000", "--"]
CMD ["./epg-approved", "update", "vod"]

