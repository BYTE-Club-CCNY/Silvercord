# api/Dockerfile

# ---- build stage ----
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Cache deps
COPY go.mod ./
RUN go mod download

# Copy source
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o api

# ---- runtime stage ----
FROM scratch

WORKDIR /app
COPY --from=builder /app/api /api

EXPOSE 8080
ENTRYPOINT ["/api"]

