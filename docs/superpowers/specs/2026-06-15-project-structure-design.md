# notifyre — Project Structure Redesign

**Date:** 2026-06-15
**Status:** Approved

## Context

All Go source files currently live flat in the repo root as `package main`. This works but doesn't follow standard Go project layout. The goal is to reorganise into `cmd/` (command wiring) and `internal/` (library logic) without changing any behaviour.

## Target Layout

```
cmd/notifyre/
├── main.go          # cobra root, Execute()
├── serve.go         # TelegramSender interface, sendHandler, newServeCmd
├── send.go          # newSendCmd
├── snippet.go       # newSnippetCmd, getSnippet
└── serve_test.go

internal/
├── config/
│   ├── config.go        # Config, LoadConfig
│   └── config_test.go
├── keys/
│   ├── keys.go          # KeyStore, LoadKeys, Valid, Label
│   └── keys_test.go
└── telegram/
    ├── telegram.go      # SendRequest, TelegramClient, RenderMessage
    └── telegram_test.go

go.mod
go.sum
Dockerfile
docker-compose.yml
CLAUDE.md
```

## Package Assignments

| File | Old package | New package |
|---|---|---|
| `main.go` → `cmd/notifyre/main.go` | `main` | `main` |
| `serve.go` → `cmd/notifyre/serve.go` | `main` | `main` |
| `send.go` → `cmd/notifyre/send.go` | `main` | `main` |
| `snippet.go` → `cmd/notifyre/snippet.go` | `main` | `main` |
| `config.go` → `internal/config/config.go` | `main` | `config` |
| `keys.go` → `internal/keys/keys.go` | `main` | `keys` |
| `telegram.go` → `internal/telegram/telegram.go` | `main` | `telegram` |

## Import Changes

`cmd/notifyre/*.go` gains three imports:

```go
"notifyre/internal/config"
"notifyre/internal/keys"
"notifyre/internal/telegram"
```

`SendRequest` lives in `internal/telegram` — `serve.go` references it as `telegram.SendRequest`.

## Build Command

```bash
go build -o notifyre ./cmd/notifyre/
```

Tests are unchanged in behaviour; run with:

```bash
go test ./...
```

## CLAUDE.md Updates

The Commands section build line changes to `go build -o notifyre ./cmd/notifyre/`. The project structure table is updated to reflect the new paths.

## What Does Not Change

- All type names, method signatures, and behaviour are identical
- HTTP API contract is unchanged
- Environment variable config is unchanged
- No new dependencies
