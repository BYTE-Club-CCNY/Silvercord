#!/bin/bash

set -e

PROJECT_DIR="/root/Silvercord"

echo "Changing into Silvercord directory"
cd "${PROJECT_DIR}"

echo "Fetching latest changes from GitHub"
git fetch && git reset origin/main --hard

echo "Stopping all running containers"
docker-compose down

echo "Rebuilding and starting services"
docker-compose up -d --build

echo "Waiting for services to be healthy..."
sleep 5

echo "Deploying Discord commands"
docker-compose run --rm discord-deploy

echo "Checking service status"
docker-compose ps

echo "Deployment successful!"
echo ""
echo "View logs with: docker-compose logs -f"
echo "Check agent health: curl http://localhost:5000/health"
