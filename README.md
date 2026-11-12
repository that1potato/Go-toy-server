# Go-toy-server

A toy go server for practice.

# Production-readiness checklist

## Done
- [x] Method-based routing with Go 1.22+ `ServeMux` patterns (`GET /health`, `POST /items`)
- [x] Per-request timeout middleware (context with deadline)
- [x] Basic health endpoint (`/health`)
- [x] Graceful shutdown on SIGINT/SIGTERM via `server.Shutdown(ctx)`
- [x] Shutdown deadline (30 s cap, returns early when handlers finish)
- [x] `http.ErrServerClosed` treated as a planned shutdown, not a fatal error
- [x] Background tasks tracked with a `sync.WaitGroup` (passed by pointer, `Add` before `go`)
- [x] Waiting on background tasks with a deadline (`waitGroupWithContext`)
- [x] One shared shutdown budget for server + background tasks

## Must
- [ ] Server-level timeouts on `http.Server` (`ReadHeaderTimeout`, `ReadTimeout`, `WriteTimeout`, `IdleTimeout`)
- [ ] Request body size limit (`http.MaxBytesReader`)
- [ ] Tests with `net/http/httptest`
- [ ] Remove debug `A`/`B` log lines
- [ ] Request logging with status code and duration (log after the handler runs)
- [ ] Structured logging with `log/slog`

## Should
- [ ] Config from environment variables or flags (port, timeouts, task count)
- [ ] Readiness check that starts failing once shutdown begins
- [ ] Correct status codes (e.g. replace 408 for server-side timeouts)
- [ ] Consistent JSON error response body
- [ ] Make every handler respect `ctx.Done()` (middleware timeout only works if handlers cooperate)
- [ ] Atomic background tasks (open design question: longer budget, stop starting new tasks, or crash-safe commit)

## Later
- [ ] Metrics (e.g. Prometheus)
- [ ] Tracing (e.g. OpenTelemetry)
