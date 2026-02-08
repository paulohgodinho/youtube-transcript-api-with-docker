# Stage 1: Builder stage with Poetry
FROM python:3.13-slim AS builder

WORKDIR /app

# Install Poetry
RUN pip install --no-cache-dir poetry

# Copy dependency files
COPY pyproject.toml ./

# Configure Poetry to not create virtual env (we're in a container)
# Update lock file and install dependencies
RUN poetry config virtualenvs.create false && \
    poetry lock && \
    poetry install --only main --no-interaction --no-ansi --no-root

# Stage 2: Runtime stage with minimal image
FROM python:3.13-slim

WORKDIR /app

# Copy installed packages from builder
COPY --from=builder /usr/local/lib/python3.13/site-packages /usr/local/lib/python3.13/site-packages
COPY --from=builder /usr/local/bin /usr/local/bin

# Copy application code
COPY youtube_transcript_api ./youtube_transcript_api

# Copy entrypoint script
COPY entrypoint.sh ./
RUN chmod +x entrypoint.sh

# Expose port for server mode
EXPOSE 5000

# Use entrypoint script for dual mode
ENTRYPOINT ["./entrypoint.sh"]
