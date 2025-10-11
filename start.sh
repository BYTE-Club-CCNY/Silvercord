#!/bin/bash

# Silvercord Startup Script
set -e

echo "Starting Silvercord services..."
echo ""

# Cleanup function
cleanup() {
    echo ""
    echo "Stopping services..."
    kill $(jobs -p) 2>/dev/null || true
    exit 0
}

trap cleanup SIGINT SIGTERM

# Start Python gRPC Server
echo "[1/3] Starting Python gRPC server (port 50051)..."
cd silvercord_agent
uv run python grpc_server.py &
cd ..
sleep 2

# Start Go API Server
echo "[2/3] Starting Go API server (port 8080)..."
cd api
go run main.go &
cd ..
sleep 2

# Start Discord Bot
echo "[3/3] Starting Discord bot..."
node src/index.js &

echo ""
echo "All services started!"
echo "Press Ctrl+C to stop all services"
echo ""

# Wait for all background processes
wait
