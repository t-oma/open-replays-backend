# Open-Replays Backend

[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

Backend API for Open-Replays - a video hosting platform for game replays and clips.

## Features

- **Video Upload**: Support for MP4, MOV, WEBM formats (up to 100MB)
- **Thumbnail Generation**: Automatic thumbnail extraction from videos
- **Structured Errors**: Detailed error responses with codes and validation details
- **Fluent Validation**: Type-safe validation API with conditional rules
- **Structured Logging**: JSON/text logging with configurable levels

## Quick Start

### Prerequisites

- Go 1.21+
- SQLite3
- FFmpeg (for thumbnail generation)

### Installation

```bash
# Clone the repository
git clone https://github.com/t-oma/open-replays-backend.git
cd open-replays-backend

# Install dependencies
go mod download

# Run the server
go run ./cmd/api
```

The server will start on `http://localhost:8080`.

### Configuration

Configuration can be done via:

1. **Config file** (`config.yaml`):
```yaml
server:
  host: "0.0.0.0"
  port: 8080

database:
  path: "db.sqlite3"

storage:
  base-dir: "uploads"
  public-url: "http://localhost:8080/media"

video:
  max-file-size: "100MB"
  allowed-extensions: [".mp4", ".mov", ".webm"]
```

2. **Environment variables**:
```bash
export OPEN_REPLAYS_SERVER_PORT=3000
export OPEN_REPLAYS_DATABASE_PATH=/path/to/db.sqlite3
export OPEN_REPLAYS_VIDEO_MAX_FILE_SIZE=50MB
```

## API Endpoints

### Videos

- `GET /api/v1/videos` - List all videos
- `GET /api/v1/videos/:id` - Get video by ID
- `POST /api/v1/videos/upload` - Upload new video
- `DELETE /api/v1/videos/:id` - Delete video

### Media

- `GET /media/:path` - Serve video/thumbnail files

### Example Request

```bash
# Upload a video
curl -X POST http://localhost:8080/api/v1/videos/upload \
  -F "title=My Gameplay" \
  -F "description=Epic play" \
  -F "video=@gameplay.mp4" \
  -F "thumbnail=@thumbnail.jpg"
```

## Error Handling

The API returns structured error responses:

```json
{
  "code": "VALIDATION_ERROR",
  "message": "validation failed",
  "details": [
    {
      "field": "title",
      "message": "Title is required",
      "tag": "required"
    }
  ]
}
```

## Project Structure

```
.
├── cmd/api/           # Application entry point
├── internal/
│   ├── api/
│   │   ├── config/    # Configuration management
│   │   ├── domain/    # Domain models and errors
│   │   ├── http/      # HTTP handlers and validation
│   │   ├── repository/# Data access layer
│   │   ├── usecase/   # Business logic
│   │   └── validation/# Fluent validation API
│   └── pkg/parse/     # Utility parsers
├── configs/           # Configuration files
└── uploads/           # Storage directory
```

## Development

### Build

```bash
go build -o api ./cmd/api
```

### Run Tests

```bash
go test ./...
```

### Run Linter

```bash
golangci-lint run
```

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Version History

- **v0.1.0** - Initial release with validation API and error handling
