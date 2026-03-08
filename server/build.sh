#!/bin/bash

set -e

export CGO_ENABLED=0
export LDFLAGS="-s -w"

# Generate OpenAPI documentation from Go comments
echo "Generating OpenAPI documentation..."
swag init -q

# Ensure openapi_generated.yaml is created from generated swagger.yaml
if [ -f docs/swagger.yaml ]; then
	cp docs/swagger.yaml openapi_generated.yaml
	echo "✓ OpenAPI spec generated: openapi_generated.yaml"
else
	echo "✗ Warning: Failed to generate OpenAPI spec"
fi

rm -rf ./bin
mkdir -p ./bin

GOOS=darwin GOARCH=arm64 go build -trimpath -ldflags "$LDFLAGS" -o ./bin/youtube-api-macos-arm64 .

GOOS=darwin GOARCH=amd64 go build -trimpath -ldflags "$LDFLAGS" -o ./bin/youtube-api-macos-amd64 .

GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "$LDFLAGS" -o ./bin/youtube-api-linux-amd64 .

GOOS=windows GOARCH=amd64 go build -trimpath -ldflags "$LDFLAGS" -o ./bin/youtube-api-windows-amd64.exe .

chmod +x ./bin/youtube-api-*

echo "✓ Build complete! Binaries in ./bin/"
ls -lh ./bin/
