FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git

WORKDIR /build

COPY go.mod go.sum* ./

RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o youtube-api .

FROM python:3.11-alpine

WORKDIR /app

RUN apk add --no-cache \
    ca-certificates \
    curl

RUN pip install --no-cache-dir youtube-transcript-api==1.2.4

COPY --from=builder /build/youtube-api /app/youtube-api

RUN chmod +x /app/youtube-api

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

ENV PYTHON_BIN=python3 \
    SERVER_PORT=8080 \
    REQUEST_TIMEOUT=30s

ENTRYPOINT ["/app/youtube-api"]
CMD ["-port", "8080"]
