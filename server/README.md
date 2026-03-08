# YouTube Transcript API Go Server

A minimal Go HTTP server that wraps the Python `youtube-transcript-api` CLI tool, providing a REST API with JSON request/response bodies.

## Requirements

- Go 1.16+ (for building)
- Python 3.8+ 
- `youtube-transcript-api` Python package (`pip install youtube-transcript-api`)

## Installation

### 1. Install the Python package

```bash
pip install youtube-transcript-api
```

### 2. Build the Go server

```bash
cd server
go build -o youtube-api .
```

### 3. Run the server

```bash
./youtube-api
```

The server will start on `http://localhost:8080` by default.

## Configuration

Configure the server using command-line flags or environment variables:

### Command-line Flags

```bash
./youtube-api -python python3 -port 8080 -timeout 30s
```

- `-python` - Path to Python executable (default: `python3`)
- `-port` - Port to listen on (default: `8080`)
- `-timeout` - Request timeout duration (default: `30s`)

### Environment Variables

- `PYTHON_BIN` - Path to Python executable
- `SERVER_PORT` - Port to listen on
- `REQUEST_TIMEOUT` - Request timeout duration

Example:

```bash
PYTHON_BIN=python3 SERVER_PORT=3000 REQUEST_TIMEOUT=60s ./youtube-api
```

## API Endpoints

### GET /health

Health check endpoint.

**Response:**
```json
{
  "status": "ok"
}
```

### GET /version

Get the installed version of `youtube-transcript-api`.

**Response:**
```json
{
  "version": "1.2.4"
}
```

### POST /transcripts

Fetch transcripts for one or more YouTube videos.

**Request Body:**
```json
{
  "videoIds": ["dQw4w9WgXcQ", "jNQXAC9IVRw"],
  "languages": ["en", "de"],
  "format": "pretty",
  "excludeGenerated": false,
  "excludeManuallyCreated": false,
  "translate": ""
}
```

**Request Fields:**

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `videoIds` | array[string] | Yes | - | List of YouTube video IDs to fetch transcripts for |
| `languages` | array[string] | No | `["en"]` | Language codes in descending priority (e.g., `["en", "de"]`) |
| `format` | string | No | `"pretty"` | Output format: `pretty`, `json`, `txt`, `vtt` |
| `excludeGenerated` | boolean | No | `false` | Exclude auto-generated transcripts |
| `excludeManuallyCreated` | boolean | No | `false` | Exclude manually created transcripts |
| `translate` | string | No | `""` | Translate transcript to this language code |

**Response:**
```json
{
  "success": true,
  "transcripts": [
    {
      "videoId": "dQw4w9WgXcQ",
      "data": "Hey there\nhow are you\n..."
    }
  ]
}
```

### POST /list

List available transcripts for one or more YouTube videos.

**Request Body:**
```json
{
  "videoIds": ["dQw4w9WgXcQ"]
}
```

**Response:**
```json
{
  "success": true,
  "transcripts": [
    {
      "videoId": "dQw4w9WgXcQ",
      "data": "Available transcripts:\n..."
    }
  ]
}
```

## Example Usage

### Fetch a transcript

```bash
curl -X POST http://localhost:8080/transcripts \
  -H "Content-Type: application/json" \
  -d '{
    "videoIds": ["dQw4w9WgXcQ"],
    "languages": ["en"],
    "format": "pretty"
  }'
```

### Fetch and translate a transcript

```bash
curl -X POST http://localhost:8080/transcripts \
  -H "Content-Type: application/json" \
  -d '{
    "videoIds": ["dQw4w9WgXcQ"],
    "languages": ["en"],
    "format": "json",
    "translate": "de"
  }'
```

### List available transcripts

```bash
curl -X POST http://localhost:8080/list \
  -H "Content-Type: application/json" \
  -d '{
    "videoIds": ["dQw4w9WgXcQ"]
  }'
```

### Check server health

```bash
curl http://localhost:8080/health
```

### Get API version

```bash
curl http://localhost:8080/version
```

## Error Handling

### Python package not installed

If `youtube-transcript-api` is not installed, the server will exit on startup with an error message:

```
Error: youtube_transcript_api must be installed. Run: pip install youtube-transcript-api
```

### Invalid request

```json
{
  "success": false,
  "error": "videoIds is required and must not be empty"
}
```

### Command failed

If a Python command fails (e.g., invalid video ID):

```json
{
  "success": true,
  "transcripts": [
    {
      "videoId": "invalid-id",
      "data": {
        "error": "command failed: unable to get video id_invalid-id"
      }
    }
  ]
}
```

## OpenAPI / Swagger Documentation

The API includes comprehensive OpenAPI/Swagger documentation that is automatically generated from source code comments using [swag](https://github.com/swaggo/swag).

### Accessing the Specification

The OpenAPI specification is available in two formats:

- **YAML (auto-generated)**: `server/openapi_generated.yaml` - The authoritative spec that stays in sync with code
- **JSON (auto-generated)**: `server/docs/swagger.json` - Same spec in JSON format
- **YAML (manual reference)**: `server/openapi.yaml` - Manual reference documentation for learning

### Regenerating the Specification

If you modify the API (add/remove endpoints, change request/response types), regenerate the specification:

```bash
cd server
swag init -q
```

This command:
1. Parses Swagger/OpenAPI comments from Go source code
2. Generates `docs/swagger.yaml` and `docs/swagger.json`
3. The build script (`./build.sh`) automatically regenerates the spec

Alternatively, just run the build script which includes spec generation:

```bash
cd server
./build.sh
```

### Documentation Comments

The API documentation is maintained as comments in the source code:

- **`main.go`** - API-level metadata (title, version, description, contact, license)
- **`handlers.go`** - Endpoint documentation (@Summary, @Description, @Param, @Success, @Failure)
- **`types.go`** - Data model documentation with field descriptions and examples

When you modify endpoints or data structures, update the corresponding documentation comments and regenerate the spec.

## Development

### Project Structure

```
server/
├── main.go              # Server entrypoint + API metadata
├── handlers.go          # HTTP endpoint handlers + endpoint docs
├── types.go             # Request/response types + type docs
├── cli.go               # Python CLI wrapper
├── build.sh             # Build script with OpenAPI generation
├── docs/                # Auto-generated OpenAPI specifications (generated by swag)
│   ├── docs.go          # Embedded spec in Go code
│   ├── swagger.json     # OpenAPI 2.0 spec in JSON format
│   └── swagger.yaml     # OpenAPI 2.0 spec in YAML format
├── openapi_generated.yaml   # Copy of swagger.yaml (main spec file)
├── openapi.yaml         # Manual reference documentation
├── go.mod               # Go module definition
├── .gitignore
└── README.md
```

### Building

```bash
cd server
./build.sh
```

This builds binaries for multiple platforms and automatically generates the OpenAPI specification.

For a quick development build:

```bash
go build -o youtube-api .
```

## Architecture

The server is single-threaded and non-concurrent by design:

1. **HTTP Handler** - Receives JSON requests and validates input
2. **Python CLI Wrapper** - Builds Python CLI commands and executes them via subprocess
3. **Output Parsing** - Returns Python CLI output as JSON responses

Each request is processed sequentially. There is no request queueing, connection pooling, or concurrent execution.

## Notes

- The server performs a startup check to ensure `youtube-transcript-api` is installed
- Each HTTP request executes a new Python subprocess
- There is a configurable request timeout (default: 30 seconds)
- The server uses only the Go standard library (no external dependencies)

## License

MIT (inherited from youtube-transcript-api)
