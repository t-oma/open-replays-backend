# AGENTS.md - Open-Replays Backend

## Build/Lint/Test Commands
- **Build**: `go build ./cmd/api`
- **Run**: `go run ./cmd/api`
- **Lint**: `golangci-lint run` (uses .golangci.yaml config)
- **Format**: `goimports -w .` (also `gofmt -w .`, `gofumpt -w .`)
- **Generate SQL**: `sqlc generate` (uses sqlc.yaml config)
- **Test**: No tests configured yet

## Code Style Guidelines

### Imports
- Standard library first, then third-party packages, then local packages
- Use `goimports` for automatic import organization and formatting

### Formatting
- Use `gofmt`, `gofumpt`, and `goimports` for consistent formatting
- Replace `interface{}` with `any`

### Types & Naming
- **Exported types/functions**: PascalCase (e.g., `VideosHandler`, `NewVideosHandler`)
- **Unexported**: camelCase (e.g., `videosRepo`)
- **Struct methods**: Receiver name should be 1-2 characters (e.g., `func (h *VideosHandler)`)
- **JSON tags**: Use camelCase for API fields (e.g., `json:"filename"`)

### Error Handling
- Return errors from functions, don't panic
- Use `errors.Is()` for error type checking
- Define custom domain errors in `domain/errors.go`
- Handle errors at HTTP layer with appropriate status codes

### Architecture
- **Clean Architecture**: domain → usecase → repository → http layers
- **Domain models**: Pure Go structs in `internal/api/domain/`
- **DTOs**: API request/response structs in `internal/api/http/v1/`
- **Database**: Use sqlc for type-safe queries, migrations in `internal/api/db/`

### HTTP API
- Use Gin framework with structured responses: `{success: bool, data: any, error: string, code: int}`
- Handle validation errors (400), not found (404), server errors (500)
- Use `ShouldBind()` for request parsing, `ShouldBindUri()` for path params