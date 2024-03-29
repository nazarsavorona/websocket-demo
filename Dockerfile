# Stage 1: Build the Go executable
FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o websocket-server

# Stage 2: Create a minimal image to run the executable
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/websocket-server .

EXPOSE 8081

ENV PORT 8081

ENV NODE_CONNECTOR_URL node-connector-production.up.railway.app
ENV NODE_URL 46.211.11.198:8081

CMD ["./websocket-server"]
