# Docker Setup for YouTube Transcript API

This directory contains Docker configuration for running the YouTube Transcript API server in a container.

## Architecture

The setup uses a **multi-stage build** to create a lean, production-ready image:

### Stage 1: Builder
- **Base Image**: `golang:1.21-alpine`
- **Purpose**: Build the Go server binary
- **Output**: Compiled `youtube-api` executable

### Stage 2: Runtime
- **Base Image**: `python:3.11-alpine`
- **Purpose**: Run the server with Python support
- **Contents**:
  - Go binary from builder stage
  - Python 3.11
  - `youtube-transcript-api` Python package
  - System utilities (curl, ca-certificates)

## Image Size

The final Docker image is optimized for minimal size:
- Base Python Alpine: ~100MB
- Go binary: ~8MB
- Python package: ~2MB
- **Total: ~110MB**

## Quick Start

### Using Docker Compose

```bash
# From the repository root
docker compose -f docker/docker-compose.yml up -d
```

The service will be available at `http://localhost:8080`

### Using Docker directly

```bash
# Build the image
docker build -f docker/Dockerfile -t youtube-api:latest .

# Run the container
docker run -d \
  --name youtube-transcript-api \
  -p 8080:8080 \
  -e PYTHON_BIN=python3 \
  -e SERVER_PORT=8080 \
  -e REQUEST_TIMEOUT=30s \
  youtube-api:latest
```

## Configuration

### Environment Variables

Configure the server using environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `PYTHON_BIN` | `python3` | Path to Python executable |
| `SERVER_PORT` | `8080` | Port to listen on |
| `REQUEST_TIMEOUT` | `30s` | Request timeout duration |

### Docker Compose Example

```yaml
environment:
  - PYTHON_BIN=python3
  - SERVER_PORT=8080
  - REQUEST_TIMEOUT=60s
```

### Docker Run Example

```bash
docker run -p 8080:8080 \
  -e SERVER_PORT=8080 \
  -e REQUEST_TIMEOUT=60s \
  youtube-api:latest
```

## Health Check

The container includes a built-in health check:

```bash
curl http://localhost:8080/health
```

Docker Compose will automatically check health every 30 seconds (after 5 second startup grace period).

## API Endpoints

Once running, the API is available at `http://localhost:8080`:

### Health Check
```bash
curl http://localhost:8080/health
```

### Get Version
```bash
curl http://localhost:8080/version
```

### Fetch Transcripts
```bash
curl -X POST http://localhost:8080/transcripts \
  -H "Content-Type: application/json" \
  -d '{
    "videoIds": ["dQw4w9WgXcQ"],
    "languages": ["en"],
    "format": "json"
  }'
```

### List Transcripts
```bash
curl -X POST http://localhost:8080/list \
  -H "Content-Type: application/json" \
  -d '{
    "videoIds": ["dQw4w9WgXcQ"]
  }'
```

See the main README in `server/` for full API documentation.

## Container Logs

### Docker Compose
```bash
docker compose -f docker/docker-compose.yml logs -f youtube-api
```

### Docker
```bash
docker logs -f youtube-transcript-api
```

## Stopping the Container

### Docker Compose
```bash
docker compose -f docker/docker-compose.yml down
```

### Docker
```bash
docker stop youtube-transcript-api
docker rm youtube-transcript-api
```

## Build Details

### Multi-Stage Build Benefits

1. **Smaller Final Image**: Builder dependencies not included
2. **Security**: No Go build tools or git in runtime image
3. **Clean Separation**: Clear distinction between build and runtime
4. **Faster Iteration**: Change Go code without rebuilding Python layer

### Build Arguments

If you want to customize the build:

```dockerfile
# To use different Go version
docker build --build-arg GO_VERSION=1.22 -f docker/Dockerfile -t youtube-api .

# To use different Python version
docker build --build-arg PYTHON_VERSION=3.12 -f docker/Dockerfile -t youtube-api .
```

(Not currently parameterized - can be added if needed)

## Troubleshooting

### Python Package Not Found

If you see "youtube_transcript_api must be installed" error:

```bash
# Rebuild the image without cache
docker compose -f docker/docker-compose.yml build --no-cache
docker compose -f docker/docker-compose.yml up -d
```

### Port Already in Use

If port 8080 is already in use:

```bash
# Run on a different port
docker run -p 9000:8080 youtube-api:latest
```

Or in docker-compose.yml:
```yaml
ports:
  - "9000:8080"
```

### Health Check Failing

Check the container logs:

```bash
docker logs youtube-transcript-api
```

The health check requires curl inside the container. If disabled in your setup:

```bash
# Remove HEALTHCHECK from Dockerfile or
docker inspect youtube-transcript-api | grep -A 10 Health
```

## Files

- `Dockerfile` - Multi-stage Docker build configuration
- `docker-compose.yml` - Docker Compose service definition
- `README.md` - This documentation

## See Also

- `../server/README.md` - Server documentation and API reference
- `../README.md` - Main project documentation
