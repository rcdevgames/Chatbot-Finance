# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata gcc musl-dev

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with CGO enabled for PostgreSQL
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o chatbot cmd/bot/main.go

# Final stage - use distroless debian nonroot with C library support
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/chatbot .

# Copy timezone data from builder stage
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy ca-certificates from builder stage
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# Expose port
EXPOSE 8080

CMD ["/app/chatbot"]
