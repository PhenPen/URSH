# URSH

A URL shortener written in Go, built stage by stage as a learning project.

## What it does

The finished service will accept a long URL, return a short code for it, and redirect visitors from that short code to the original destination.

The repository is currently at Stage 1. It runs an HTTP server on port 8080 with sensible timeouts and exposes a health endpoint. There is no URL creation, validation, short-code generation, or storage yet. Each stage is implemented, tested, and reviewed before the next one starts.

## Why it exists

This is an educational backend project, not a product competing with existing link shorteners. The goals are a working service and a real understanding of the Go code that implements it: HTTP handling, validation, random code generation, persistence, concurrency, and eventually performance measurement.

Alternatives such as Bitly, Shlink, or YOURLS are better choices if you just need a shortener today. This repository is worth reading if you want to see how one is built from the standard library up.

## Features

Working now:

- A Go module named `ursh`.
- An HTTP server using `net/http`, listening on `:8080`.
- Server timeouts for reading, writing, and idle connections.
- `GET /health`, which returns `200 OK` with the body `ok`.
- A startup log line and a fatal error if the listener fails.

Planned, in stage order: URL creation API, in-memory redirects, persistent storage, refactoring, operational readiness, and measurement-driven performance work. Nothing from that list exists yet.

## Requirements

- Go 1.27 or newer, check with `go version`.

## Install

Nothing to install beyond Go. The module path `ursh` is local-only, so `go get` from another machine is not possible; build from a checkout instead:

```powershell
go build .
.\URSH.exe
```

`go build` compiles the current directory's `main` package into an executable named after the directory. Running it starts the server:

```text
2026/11/09 12:00:00 listening on http://localhost:8080
```

## Quick start

```powershell
go run .
```

`go run` compiles the package and runs the binary in one step, without leaving an executable behind. Then in another terminal:

```powershell
curl.exe -i http://localhost:8080/health
```

The response is a `200 OK` with the plain-text body `ok`:

```http
HTTP/1.1 200 OK
Content-Type: text/plain; charset=utf-8

ok
```

Stop the server with Ctrl+C.

## API

| Method | Path | Response |
| --- | --- | --- |
| `GET` | `/health` | `200 OK`, body `ok`, `Content-Type: text/plain; charset=utf-8` |

Requests to any other path return `404 Not Found`. A non-GET method on `/health`, such as `POST`, returns `405 Method Not Allowed`.

## How it works

`main.go` builds an `http.Server` value rather than calling the package-level `http.ListenAndServe`. That matters because the `http.Server` struct is where timeouts live; the convenience function has none.

- `Addr` is the listening address, here the constant `:8080`, meaning all interfaces on port 8080.
- `Handler` is `newMux()`, the router that maps patterns to handler functions.
- `ReadHeaderTimeout` bounds how long a client may take to send request headers, which limits slow-header attacks.
- `ReadTimeout` and `WriteTimeout` bound the whole request read and response write.
- `IdleTimeout` bounds how long an idle keep-alive connection is kept.

`newMux` returns an `*http.ServeMux` with one route registered as `"GET /health"`. Since Go 1.22 a pattern may include the method, so the mux itself rejects other methods with `405 Method Not Allowed` and answers `GET` and `HEAD`.

`healthHandler` has the standard handler signature, `func(http.ResponseWriter, *http.Request)`. The `http.ResponseWriter` is the write side of the response: the server sends headers the first time the handler writes a body or calls `WriteHeader`. The `*http.Request` is the parsed request, unused here because health checks need no input. The handler sets `Content-Type`, writes the status, then writes `ok\n`. A write error means the client has gone away, so it is logged rather than reported to a client that cannot receive it.

`ListenAndServe` blocks until the server stops. It returns `http.ErrServerClosed` on a clean shutdown, so that error is ignored and every other error is fatal via `log.Fatalf`.

## Project layout

| Path | Purpose |
| --- | --- |
| `go.mod` | Module path and Go version. No dependencies yet. |
| `main.go` | Server setup, routing, and handler functions. |
| `main_test.go` | Handler and routing tests using `net/http/httptest`. |
| `.gitignore` | Build artifacts and local notes that must stay out of history. |

## Tests

Handler tests use `net/http/httptest`, so no live server or real port is needed. They cover the health handler's status, headers and body, plus the router's handling of `GET`, `HEAD`, a rejected `POST`, and an unknown path:

```powershell
go test ./...
```

The race detector requires cgo and a C compiler; the current machine has no `gcc`, so `go test -race` is not runnable here yet. See Notes.

## Notes and limitations

- The port is a constant in `main.go`. Configuration for environments is deferred until the operational readiness stage.
- There is no graceful shutdown yet. Add it when the process has resources worth closing, not before.
- The health endpoint reports only that the process is answering. It does not check storage, because no storage exists yet.
- `go test -race` cannot run on the current development machine: the race detector requires cgo and a C compiler, and `gcc` is not installed. Install MinGW-w64 and run with `CGO_ENABLED=1` before Stage 3 introduces concurrency.
- The module path is deliberately local (`ursh`). If the project is published, the `go.mod` module line and its import paths must be renamed to the real repository path.
- Nothing is deployed, containerized, or hardened for public traffic. Public URL shorteners can be abused for phishing, so abuse prevention and rate limiting are required before exposing it publicly.
- Markdown files are ignored by default in `.gitignore`. Only explicitly allowlisted files such as this README are tracked; local notes and the project specification stay out of history.
