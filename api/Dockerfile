# ---- build stage ----
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install certs in builder so we can copy them
RUN apk add --no-cache ca-certificates

# Cache deps
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o api

# ---- runtime stage ----
FROM scratch

# Copy CA certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy binary
COPY --from=builder /app/api /api

EXPOSE 8080
ENTRYPOINT ["/api"]
