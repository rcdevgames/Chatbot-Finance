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

# Final stage - use alpine nonroot
FROM alpine:latest

RUN addgroup -S nonroot && adduser -S nonroot -G nonroot

WORKDIR /app

# Install ca-certificates and timezone data
RUN apk add --no-cache ca-certificates tzdata

# Copy the binary from builder stage
COPY --from=builder --chown=nonroot:nonroot /app/chatbot /app/chatbot

# Copy timezone data from builder stage
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Expose port
EXPOSE 8080

ENTRYPOINT ["/app/chatbot"]
