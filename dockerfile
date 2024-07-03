# Stage 1: Build the Go application
FROM golang:1.22-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -o /app/main .

# Stage 2: Create a lightweight image to run the application
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main .
RUN chmod +x /app/main
# Expose port 8080 to the outside world
EXPOSE 8001

# Command to run the executable
CMD ["/app/main"]
