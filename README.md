# notifyre

A lightweight Go service that forwards authenticated HTTP requests to Telegram, with optional SOCKS5 proxy support and CLI integration helpers.

## Features

- **HTTP API** — POST `/send` with JSON payload to send notifications
- **API Key Authentication** — YAML-based key store with labels
- **Message Formatting** — Optional emojis and titles
- **Parse Modes** — HTML and MarkdownV2 support
- **SOCKS5 Proxy** — Optional proxy routing for Telegram API calls
- **CLI Client** — `notifyre send` command for shell scripts
- **Code Snippets** — Generate ready-to-paste integration code (Go, Node, Python, Bash, cURL)

## Quick Start

### Prerequisites

- Go 1.26.3+
- Telegram bot token from [@BotFather](https://t.me/BotFather)
- Target Telegram channel (username like `@mychannel` or numeric ID like `-100123456789`)

### Setup

1. **Clone and build:**

```bash
git clone https://github.com/serajhqi/notifyre.git
cd notifyre
go build -o notifyre ./cmd/notifyre/
```

2. **Create keys file** (copy from template):

```bash
cp keys.yaml.example keys.yaml
```

Edit `keys.yaml` with your API keys:

```yaml
keys:
  - key: "your-secret-key-1"
    label: "ci-pipeline"
  - key: "your-secret-key-2"
    label: "backup-job"
```

3. **Start the server:**

```bash
export BOT_TOKEN="your-telegram-bot-token"
export CHANNEL="@your-channel"
export KEYS_FILE="./keys.yaml"
./notifyre serve
```

Server listens on `http://localhost:8080` by default.

## Usage

### HTTP API

Send a notification:

```bash
curl -X POST http://localhost:8080/send \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-key-1" \
  -d '{
    "message": "Build succeeded",
    "title": "CI Pipeline",
    "level": "success",
    "parse_mode": "HTML"
  }'
```

**Request body:**

| Field | Required | Values |
|-------|----------|--------|
| `message` | yes | Notification text |
| `title` | no | Short label (prepended with emoji if level set) |
| `level` | no | `info` (ℹ️), `success` (✅), `warning` (⚠️), `error` (❌) |
| `parse_mode` | no | `HTML` or `MarkdownV2` |
| `disable_notification` | no | `true` to send silently |

**Response:**

- `200 OK`: `{"ok": true}`
- `400 Bad Request`: `{"ok": false, "error": "..."}`
- `401 Unauthorized`: `{"ok": false, "error": "invalid api key"}`
- `502 Bad Gateway`: `{"ok": false, "error": "telegram error: ..."}`

### CLI Client

Send notifications from shell:

```bash
./notifyre send \
  --url http://localhost:8080 \
  --key your-secret-key-1 \
  --message "Backup complete" \
  --title "Backup Job" \
  --level success
```

All flags except `--url`, `--key`, and `--message` are optional.

### Code Snippets

Generate integration code for your language:

```bash
./notifyre snippet --lang go --url http://localhost:8080 --key your-secret-key-1
./notifyre snippet --lang python --url http://localhost:8080 --key your-secret-key-1
./notifyre snippet --lang bash --url http://localhost:8080 --key your-secret-key-1
./notifyre snippet --lang node --url http://localhost:8080 --key your-secret-key-1
./notifyre snippet --lang curl --url http://localhost:8080 --key your-secret-key-1
```

Copy the output directly into your project.

## Configuration

All configuration via environment variables:

| Variable | Required | Default | Notes |
|----------|----------|---------|-------|
| `BOT_TOKEN` | yes | — | Telegram bot token from @BotFather |
| `CHANNEL` | yes | — | Target channel: `@username` or `-100...` (numeric ID) |
| `KEYS_FILE` | yes | — | Path to YAML keys file |
| `PORT` | no | `8080` | HTTP server listen port |
| `PROXY_ADDR` | no | — | SOCKS5 proxy (e.g., `socks5://user:pass@host:1080`) |

### Docker

Build and run with Docker:

```bash
docker build -t notifyre .
docker run -e BOT_TOKEN=... -e CHANNEL=... -e KEYS_FILE=/keys/keys.yaml \
  -v ./keys.yaml:/keys/keys.yaml \
  -p 8080:8080 \
  notifyre serve
```

Or use `docker-compose`:

```bash
docker-compose up
```

(Edit `.env` for BOT_TOKEN and CHANNEL before running.)

## Project Structure

```
notifyre/
├── cmd/notifyre/              # CLI commands (package main)
│   ├── main.go                # Cobra root command
│   ├── serve.go               # HTTP server subcommand
│   ├── send.go                # CLI send subcommand
│   ├── snippet.go             # Snippet generator
│   └── serve_test.go          # Server handler tests
├── internal/
│   ├── config/                # Configuration loading
│   │   ├── config.go
│   │   └── config_test.go
│   ├── keys/                  # API key store
│   │   ├── keys.go
│   │   └── keys_test.go
│   └── telegram/              # Telegram client
│       ├── telegram.go        # Client, SendRequest, RenderMessage
│       └── telegram_test.go
├── go.mod / go.sum            # Go dependencies
├── CLAUDE.md                  # Development guide
├── Dockerfile / docker-compose.yml
├── keys.yaml.example          # Keys template
└── README.md
```

## Development

### Build

```bash
go build -o notifyre ./cmd/notifyre/
```

### Test

```bash
# All tests
go test ./...

# Specific package
go test ./internal/telegram -v
go test ./cmd/notifyre -v

# Specific test
go test -run TestRenderMessage ./internal/telegram
```

### Lint

```bash
go vet ./...
```

## Message Rendering

Output format depends on `level` and `title`:

| `level` | `title` | Output |
|---------|---------|--------|
| ✅ | ✅ | `✅ Title\nMessage` |
| ✅ | ❌ | `✅ Message` |
| ❌ | ✅ | `Title\nMessage` |
| ❌ | ❌ | `Message` |

Examples:

```
// Both set
POST /send with level="success", title="CI", message="Build done"
→ Telegram receives: "✅ CI\nBuild done"

// Level only
POST /send with level="error", message="Connection timeout"
→ Telegram receives: "❌ Connection timeout"

// Title only
POST /send with title="Backup", message="3 GB uploaded"
→ Telegram receives: "Backup\n3 GB uploaded"

// Neither
POST /send with message="Job complete"
→ Telegram receives: "Job complete"
```

## Dependencies

- `github.com/spf13/cobra` — CLI framework
- `golang.org/x/net/proxy` — SOCKS5 dialer
- `gopkg.in/yaml.v3` — YAML parsing
- stdlib `net/http` — HTTP server and client

## License

MIT

## Contributing

Improvements welcome! Please ensure:
- All tests pass: `go test ./...`
- Code is formatted: `gofmt -w .`
- Linting passes: `go vet ./...`
