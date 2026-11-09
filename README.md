# URSH

A URL shortener written in Go, built stage by stage as a learning project.

## What it does

The finished service will accept a long URL, return a short code for it, and redirect visitors from that short code to the original destination.

The repository is currently at Stage 1. It runs an HTTP server on port 8080 with sensible timeouts. There are no endpoints, URL validation, short-code generation, or storage yet. Each stage is implemented, tested, and reviewed before the next one starts.

## Why it exists

This is an educational backend project, not a product competing with existing link shorteners. The goals are a working service and a real understanding of the Go code that implements it: HTTP handling, validation, random code generation, persistence, concurrency, and eventually performance measurement.

Alternatives such as Bitly, Shlink, or YOURLS are better choices if you just need a shortener today. This repository is worth reading if you want to see how one is built from the standard library up.

## Features

Working now:

- A Go module named `ursh`.
- An HTTP server using `net/http`, listening on `:8080`.
- Server timeouts for reading, writing, and idle connections.
- A startup log line and a fatal error if the listener fails.

Planned, in stage order: health endpoint, URL creation API, in-memory redirects, persistent storage, refactoring, operational readiness, and measurement-driven performance work. Nothing from that list exists yet.

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
curl.exe -i http://localhost:8080/
```

The server is running, but no routes are registered yet, so every path currently returns `404 Not Found` from `net/http`'s default mux. Stop the server with Ctrl+C.

## How it works

`main.go` builds an `http.Server` value rather than calling the package-level `http.ListenAndServe`. That matters because the `http.Server` struct is where timeouts live; the convenience function has none.

- `Addr` is the listening address, here the constant `:8080`, meaning all interfaces on port 8080.
- `ReadHeaderTimeout` bounds how long a client may take to send request headers, which limits slow-header attacks.
- `ReadTimeout` and `WriteTimeout` bound the whole request read and response write.
- `IdleTimeout` bounds how long an idle keep-alive connection is kept.

`ListenAndServe` blocks until the server stops. It returns `http.ErrServerClosed` on a clean shutdown, so that error is ignored and every other error is fatal via `log.Fatalf`.

## Project layout

| Path | Purpose |
| --- | --- |
| `go.mod` | Module path and Go version. No dependencies yet. |
| `main.go` | Server setup and entry point. |
| `.gitignore` | Build artifacts and local notes that must stay out of history. |

## Tests

There are no tests yet. Testable behaviour starts with the health endpoint in the next commit.

```powershell
go test ./...
```

reports `no test files` at this point.

## Notes and limitations

- The port is a constant in `main.go`. Configuration for environments is deferred until the operational readiness stage.
- There is no graceful shutdown yet. Add it when the process has resources worth closing, not before.
- The module path is deliberately local (`ursh`). If the project is published, the `go.mod` module line and its import paths must be renamed to the real repository path.
- Nothing is deployed, containerized, or hardened for public traffic. Public URL shorteners can be abused for phishing, so abuse prevention and rate limiting are required before exposing it publicly.
- Markdown files are ignored by default in `.gitignore`. Only explicitly allowlisted files such as this README are tracked; local notes and the project specification stay out of history.
