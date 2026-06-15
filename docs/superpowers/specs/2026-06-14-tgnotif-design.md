# notifyre — Design Spec

**Date:** 2026-06-14  
**Status:** Approved

## Context

Long-running processes (scripts, CI jobs, build tools) need to report their status to the developer over Telegram. `notifyre` is a small Go service that accepts authenticated HTTP requests from any caller and forwards the message to a Telegram channel via a bot. It also ships as a CLI so bash scripts and colleagues can send notifications without writing HTTP code.

## Architecture

A single Go binary with three subcommands. The `serve` subcommand starts an HTTP server; `send` and `snippet` are thin client-side helpers that do not require the server to be on the same machine.

```
[caller script / service]
        |
  POST /send  +  X-API-Key
        |
 [notifyre HTTP server]
        |
  [key validation]   <-- keys.yaml
        |
  [telegram client]
        |
  [SOCKS5 dialer]   (optional — only if PROXY_ADDR is set)
        |
  api.telegram.org/bot{token}/sendMessage
        |
  [Telegram channel]
```

## Configuration

All server configuration via environment variables.

| Variable | Required | Description |
|---|---|---|
| `BOT_TOKEN` | yes | Telegram bot token from BotFather |
| `CHANNEL` | yes | Target channel: `@channelusername` or numeric ID `-100…` |
| `PROXY_ADDR` | no | SOCKS5 proxy URL: `socks5://user:pass@host:1080`. Omit for direct connection. |
| `KEYS_FILE` | yes | Path to the YAML API keys file |
| `PORT` | no | HTTP listen port. Default: `8080` |

**Keys file** (`keys.yaml`):
```yaml
keys:
  - key: "abc123secret"
    label: "service-a"
  - key: "xyz789secret"
    label: "service-b"
```

Keys are loaded once at startup. A restart is required to reload the file.

## HTTP API

**Endpoint:** `POST /send`

**Authentication:** `X-API-Key: <key>` header.

**Request body (JSON):**
```json
{
  "message": "Build finished in 42s",
  "title": "CI Pipeline",
  "level": "success",
  "parse_mode": "HTML",
  "disable_notification": false
}
```

| Field | Required | Values / Notes |
|---|---|---|
| `message` | yes | The notification body text |
| `title` | no | Short label prepended to the message |
| `level` | no | `info` (ℹ️), `success` (✅), `warning` (⚠️), `error` (❌). Default: no prefix |
| `parse_mode` | no | `HTML` or `MarkdownV2`. Default: plain text |
| `disable_notification` | no | `true` to send silently (no buzz). Default: `false` |

**Message rendering rules:**

| `level` | `title` | Rendered output |
|---|---|---|
| set | set | `✅ CI Pipeline\nBuild finished in 42s` |
| set | not set | `✅ Build finished in 42s` (emoji + space + message, single line) |
| not set | set | `CI Pipeline\nBuild finished in 42s` |
| not set | not set | `Build finished in 42s` (plain message) |

**Responses:**

| Status | Body |
|---|---|
| `200 OK` | `{"ok": true}` |
| `400 Bad Request` | `{"ok": false, "error": "message is empty"}` |
| `401 Unauthorized` | `{"ok": false, "error": "invalid api key"}` |
| `502 Bad Gateway` | `{"ok": false, "error": "telegram error: ..."}` |

## CLI Subcommands

### `notifyre serve`
Starts the HTTP server. Reads config from environment.

```
notifyre serve [--port 8080]
```

### `notifyre send`
Sends a single notification from the command line. Calls the running server over HTTP.

```
notifyre send \
  --url http://localhost:8080 \
  --key abc123secret \
  --message "Backup complete" \
  --title "Backup Job" \
  --level success \
  --disable-notification
```

All flags except `--url`, `--key`, and `--message` are optional.

### `notifyre snippet`
Prints a ready-to-paste integration snippet for a given language. Callers copy it into their repo.

```
notifyre snippet --lang go --url http://host:8080 --key abc123secret
```

Supported languages: `go`, `node`, `python`, `bash`, `curl`

## Project Structure

```
notifyre/
├── main.go        # cobra root command + subcommand wiring
├── serve.go       # serve subcommand: HTTP server + /send handler
├── send.go        # send subcommand: CLI → HTTP call
├── snippet.go     # snippet subcommand: prints language snippets
├── config.go      # Config struct + env loading
├── keys.go        # Keys YAML loading + O(1) lookup map
├── telegram.go    # Telegram client: builds message, calls Bot API, optional SOCKS5
├── go.mod
└── go.sum
```

All files are in `package main`. No internal packages needed at this scale.

## Dependencies

| Module | Purpose |
|---|---|
| `github.com/spf13/cobra` | Subcommand routing + auto `--help` |
| `golang.org/x/net/proxy` | SOCKS5 dialer |
| `gopkg.in/yaml.v3` | Keys file parsing |

## Error Handling

- Missing required env vars → fatal log + exit at startup, before binding the port
- Keys file not found or malformed → fatal log + exit at startup
- Invalid JSON request body → `400`
- Unknown or missing API key → `401`
- Telegram API returns non-200 → `502` with Telegram's error message forwarded
- SOCKS5 proxy unreachable → `502` with connection error

## Verification

1. **Unit-testable pieces:** key lookup (`keys.go`), message rendering (`telegram.go` formatter), config parsing (`config.go`)
2. **Integration test:** spin up `notifyre serve` against a real bot + channel (no proxy), call `POST /send`, confirm message arrives in Telegram
3. **CLI smoke test:** `notifyre send` → verify HTTP request is constructed correctly
4. **Snippet test:** `notifyre snippet --lang bash` → output is valid bash
5. **Proxy test:** set `PROXY_ADDR` to a local SOCKS5 proxy, verify traffic routes through it
