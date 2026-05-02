# ---- build stage ----
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build dependencies for CGO (sqlite3)
RUN apk add --no-cache ca-certificates gcc musl-dev

# Cache deps
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build with CGO enabled for sqlite3
RUN CGO_ENABLED=1 GOOS=linux go build -o api

# ---- runtime stage ----
FROM alpine:latest

# Install sqlite runtime libs and CA certs
RUN apk add --no-cache ca-certificates sqlite-libs

# Copy binary
COPY --from=builder /app/api /api

EXPOSE 8080
ENTRYPOINT ["/api"]
