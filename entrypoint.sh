#!/bin/bash
set -e

if [ $# -eq 0 ]; then
    # No arguments → Server mode
    echo "Starting Flask server on port ${PORT:-5000}..."
    exec python -m youtube_transcript_api.server
else
    # Arguments provided → CLI mode
    exec python -m youtube_transcript_api "$@"
fi
