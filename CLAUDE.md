# Flow Limit Proxy - Technical Documentation

## Project Overview
Flow Limit Proxy is a Go-based HTTP proxy that provides concurrent connection limiting and retry functionality.

## Architecture

### Current Structure
```
main.go (51 lines)
├── CLI argument parsing
├── Proxy server startup
└── Graceful shutdown

proxy.go (116 lines)
├── ListenProxy function
├── customTransport with semaphore
└── Exponential backoff retry logic
```

## Key Components

### Connection Limiting
- Uses `golang.org/x/sync/semaphore` for concurrent connection control
- Configurable limit via CLI parameter
- Applied at HTTP transport level

### Retry Logic
- Exponential backoff with `github.com/cenkalti/backoff/v4`
- Configured for network resilience
- Applied to failed proxy requests

### Graceful Shutdown
- Handles OS signals for clean termination
- 10-second timeout for ongoing requests

## Known Technical Debt

### Code Organization
- Single package structure limits reusability
- CLI logic mixed with application logic in `main()` function
- No separation of concerns between components

### Error Handling
- Context not properly used for semaphore acquisition (proxy.go:84-88)
- Limited error context for debugging
- Basic error logging without structured information

### Configuration
- Hardcoded backoff settings (proxy.go:111-116)
- No configuration file support
- Limited CLI parameter validation

### Testing
- No unit tests or integration tests
- No benchmarks for performance validation
- No test coverage reporting

## Implementation Notes

### Current Architecture
```
main.go (51 lines)
├── CLI argument parsing
├── Proxy server startup
└── Graceful shutdown

proxy.go (116 lines)
├── ListenProxy function
├── customTransport with semaphore
└── Exponential backoff retry logic
```

### Suggested Package Structure
```
cmd/
├── main.go (CLI only)
pkg/proxy/
├── proxy.go (core logic)
├── config.go (configuration)
└── transport.go (HTTP transport)
internal/
├── health/ (health checks)
└── metrics/ (monitoring)
```

## Test Commands
- Build: `go build`
- Test: `go test ./...` (after adding tests)
- Lint: `golangci-lint run` (if available)

## Notes for Implementation
- Keep backward compatibility
- Follow Go conventions and idioms
- Add proper error messages with context
- Use structured logging if needed
- Consider using cobra for CLI if complexity grows