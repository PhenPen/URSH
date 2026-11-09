# URSH

A URL shortener written in Go, built stage by stage as a learning project.

## What it does

The finished service will accept a long URL, return a short code for it, and redirect visitors from that short code to the original destination.

The repository is currently at Stage 0. It is a Go module with one minimal executable. There is no HTTP server, URL validation, short-code generation, or storage yet. Each stage is implemented, tested, and reviewed before the next one starts.

## Why it exists

This is an educational backend project, not a product competing with existing link shorteners. The goals are a working service and a real understanding of the Go code that implements it: HTTP handling, validation, random code generation, persistence, concurrency, and eventually performance measurement.

Alternatives such as Bitly, Shlink, or YOURLS are better choices if you just need a shortener today. This repository is worth reading if you want to see how one is built from the standard library up.

## Features

Working now:

- A Go module named `ursh`.
- A minimal `main` executable.

Planned, in stage order: health endpoint, URL creation API, in-memory redirects, persistent storage, refactoring, operational readiness, and measurement-driven performance work. Nothing from that list exists yet.

## Requirements

- Go 1.27 or newer, check with `go version`.

## Install

Nothing to install beyond Go. The module path `ursh` is local-only, so `go get` from another machine is not possible; build from a checkout instead:

```powershell
go build .
.\URSH.exe
```

`go build` compiles the current directory's `main` package into an executable named after the directory. Running it prints:

```text
Hello
```

## Quick start

```powershell
go run .
```

`go run` compiles the package and runs the binary in one step, without leaving an executable behind. Output is `Hello`, because that is all the program does at this stage.

## How it works

`main.go` is a `package main` with a single `func main()`, the entry point Go calls when the program starts. It writes `Hello` to standard output with `fmt.Print`. `go.mod` declares the module path `ursh` and the Go version the module targets; it also records dependencies, of which there are none.

When the HTTP server arrives, this section will describe the actual request flow.

## Project layout

| Path | Purpose |
| --- | --- |
| `go.mod` | Module path and Go version. No dependencies yet. |
| `main.go` | Entry point. Prints `Hello`. |
| `.gitignore` | Build artifacts and local notes that must stay out of history. |

## Tests

There are no tests yet. Nothing beyond compilation and one print call exists to test, and testable behaviour starts at Stage 1 with the health endpoint.

```powershell
go test ./...
```

reports `no test files` at this stage.

## Notes and limitations

- The module path is deliberately local (`ursh`). If the project is published, the `go.mod` module line and its import paths must be renamed to the real repository path.
- Nothing is deployed, containerized, or hardened for public traffic. Public URL shorteners can be abused for phishing, so abuse prevention and rate limiting are required before exposing it publicly.
- Markdown files are ignored by default in `.gitignore`. Only explicitly allowlisted files such as this README are tracked; local notes and the project specification stay out of history.
