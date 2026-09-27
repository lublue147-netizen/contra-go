# Multi-stage Dockerfile for Contra Go
FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod ./
COPY cmd/ ./cmd/
COPY web/ ./web/
COPY server.go ./
COPY build.sh ./

# Compile WebAssembly game and server binary
RUN apk add --no-cache bash curl && \
    mkdir -p dist && \
    GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o dist/contra.wasm ./cmd/wasm && \
    cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" dist/ && \
    cp web/index.html dist/ && \
    cp web/style.css dist/ && \
    CGO_ENABLED=0 go build -ldflags="-s -w" -o contra-server server.go

# Final minimal runner image
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/contra-server .
COPY --from=builder /app/dist ./dist

ENV PORT=8080
EXPOSE 8080

CMD ["./contra-server"]
