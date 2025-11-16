# Build stage
FROM golang:1.24.4-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/server .

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./server"]