1. Testing (High Priority)

- [ ] Unit tests for use case layer
- [ ] Integration tests for API endpoints
- [ ] Tests for repository with test containers or SQLite in-memory

2. Graceful Shutdown (High Priority)

- [ ] Correct server shutdown
- [ ] Closing database connections
- [ ] Termination of background jobs (video processor)

3. API Documentation (Medium Priority)

- [ ] Swagger/OpenAPI specification
- [ ] Documentation of endpoints with examples

4. Docker & CI/CD (Medium Priority)

- [ ] Dockerfile for the application
- [ ] Docker Compose with SQLite volume
- [ ] GitHub Actions for tests and linter

5. Monitoring & Health Checks (Medium Priority)

- [ ] /health endpoint
- [ ] Prometheus metrics
- [ ] Ready/Liveness probes

6. Rate Limiting & Security (Low Priority)

- [ ] Rate limiting middleware
- [ ] Request size limits
- [ ] CORS configuration

7. Advanced Features (Future)

- [ ] Pagination for video list
- [ ] Video search/filtering
- [ ] Users and authorization
- [ ] S3 storage instead of local
