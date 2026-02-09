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

# Or using Makefile
make server
```

The server will start on `http://localhost:8080`.

### Configuration

Configuration can be done via:

1. **Config file** (`config.yaml`):

```yaml
server:
  host: 0.0.0.0
  port: 8080
database:
  path: db.sqlite3
  busy-timeout: 5s
  journal-mode: WAL
storage:
  base-dir: uploads
  public-url: http://localhost:8080/media
video:
  worker-count: 2
  max-file-size: 100MB
  allowed-extensions: [.mp4, .webm, .mov]
log:
  level: info
  format: json
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
      "message": "title is required",
      "tag": "required"
    }
  ]
}
```

## Development

### Build

```bash
go build -o api ./cmd/api
```

### Run Tests

```bash
go test ./...

# Or using Makefile
make test
```

### Run Linter

```bash
golangci-lint run
```

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
