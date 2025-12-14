#!/bin/sh
set -e

echo "Deploying Discord slash commands..."
node src/deploy_commands.js

echo "Starting Discord bot..."
exec node src/index.js

