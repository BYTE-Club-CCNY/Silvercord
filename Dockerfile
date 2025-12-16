# Dockerfile (discord bot)

FROM node:20-alpine

WORKDIR /app

COPY package*.json ./
RUN npm ci --omit=dev

COPY src ./src
COPY entrypoint.sh ./entrypoint.sh

CMD ["./entrypoint.sh"]

