# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Build
go build -o notifyre ./cmd/notifyre/

# Test all
go test ./...

# Test single package/file
go test -run TestRenderMessage ./internal/telegram
go test -run TestSendHandler ./cmd/notifyre

# Lint (uses go vet)
go vet ./...
```

## Architecture

Single Go binary (`package main` in `cmd/notifyre/`) with three Cobra subcommands. Library logic in `internal/config`, `internal/keys`, and `internal/telegram` packages.

- **`cmd/notifyre/`** — commands: `main.go` (cobra root), `serve.go` (HTTP server), `send.go` (CLI client), `snippet.go` (code snippets)
- **`internal/config/`** — `Config` struct, `LoadConfig()` validates required env vars
- **`internal/keys/`** — `KeyStore` type, YAML keys file parsing
- **`internal/telegram/`** — `TelegramClient`, `SendRequest`, `RenderMessage`

### Request flow

```
POST /send + X-API-Key
  → sendHandler (serve.go)        # auth via KeyStore
  → TelegramClient.Send (telegram.go) # renders message, calls Telegram API
  → api.telegram.org/bot{token}/sendMessage
```

### Key types

| Type | File | Role |
|---|---|---|
| `Config` | `config.go` | Env var config; `LoadConfig()` validates required vars |
| `KeyStore` | `keys.go` | YAML keys file; `Valid(key)` / `Label(key)` |
| `TelegramClient` | `telegram.go` | HTTP client with optional SOCKS5 proxy |
| `TelegramSender` | `serve.go` | Interface satisfied by `TelegramClient`; used for mocking in tests |
| `SendRequest` | `telegram.go` | Shared JSON struct for the HTTP API and internal Send calls |

### Message rendering (`RenderMessage` in `telegram.go`)

| `level` | `title` | Output |
|---|---|---|
| set | set | `EMOJI Title\nMessage` |
| set | unset | `EMOJI Message` |
| unset | set | `Title\nMessage` |
| unset | unset | `Message` |

### Configuration (env vars for `serve`)

| Variable | Required | Default |
|---|---|---|
| `BOT_TOKEN` | yes | — |
| `CHANNEL` | yes | — |
| `KEYS_FILE` | yes | — |
| `PROXY_ADDR` | no | direct connection |
| `PORT` | no | `8080` |

`PROXY_ADDR` must be a `socks5://` URL. Keys are loaded once at startup; restart to reload.

## graphify

This project has a knowledge graph at graphify-out/ with god nodes, community structure, and cross-file relationships.

Rules:
- For codebase questions, first run `graphify query "<question>"` when graphify-out/graph.json exists. Use `graphify path "<A>" "<B>"` for relationships and `graphify explain "<concept>"` for focused concepts. These return a scoped subgraph, usually much smaller than GRAPH_REPORT.md or raw grep output.
- If graphify-out/wiki/index.md exists, use it for broad navigation instead of raw source browsing.
- Read graphify-out/GRAPH_REPORT.md only for broad architecture review or when query/path/explain do not surface enough context.
- After modifying code, run `graphify update .` to keep the graph current (AST-only, no API cost).
