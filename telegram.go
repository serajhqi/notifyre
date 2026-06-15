package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/proxy"
)

// SendRequest is the JSON body accepted by POST /send and used internally.
type SendRequest struct {
	Message             string `json:"message"`
	Title               string `json:"title"`
	Level               string `json:"level"`
	ParseMode           string `json:"parse_mode"`
	DisableNotification bool   `json:"disable_notification"`
}

var levelEmoji = map[string]string{
	"info":    "ℹ️",
	"success": "✅",
	"warning": "⚠️",
	"error":   "❌",
}

// RenderMessage formats a SendRequest into a Telegram message string.
// Rules:
//   - level+title  → "EMOJI Title\nMessage"
//   - level only   → "EMOJI Message"
//   - title only   → "Title\nMessage"
//   - neither      → "Message"
func RenderMessage(req SendRequest) string {
	var b strings.Builder
	emoji := levelEmoji[req.Level]

	if req.Title != "" {
		if emoji != "" {
			b.WriteString(emoji + " ")
		}
		b.WriteString(req.Title + "\n")
		b.WriteString(req.Message)
	} else if emoji != "" {
		b.WriteString(emoji + " " + req.Message)
	} else {
		b.WriteString(req.Message)
	}
	return b.String()
}

// TelegramClient sends messages to a Telegram channel via the Bot API.
type TelegramClient struct {
	token      string
	channel    string
	httpClient *http.Client
}

// NewTelegramClient creates a client. If proxyAddr is non-empty it must be a
// socks5:// URL; the client will dial through that proxy.
func NewTelegramClient(token, channel, proxyAddr string) (*TelegramClient, error) {
	transport := &http.Transport{}
	if proxyAddr != "" {
		u, err := url.Parse(proxyAddr)
		if err != nil {
			return nil, fmt.Errorf("invalid proxy address: %w", err)
		}
		if u.Scheme != "socks5" {
			return nil, fmt.Errorf("proxy address must use socks5:// scheme, got %q", u.Scheme)
		}
		var auth *proxy.Auth
		if u.User != nil {
			pass, _ := u.User.Password()
			auth = &proxy.Auth{User: u.User.Username(), Password: pass}
		}
		dialer, err := proxy.SOCKS5("tcp", u.Host, auth, proxy.Direct)
		if err != nil {
			return nil, fmt.Errorf("create socks5 dialer: %w", err)
		}
		transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			type contextDialer interface {
				DialContext(ctx context.Context, network, addr string) (net.Conn, error)
			}
			if cd, ok := dialer.(contextDialer); ok {
				return cd.DialContext(ctx, network, addr)
			}
			return dialer.Dial(network, addr)
		}
	}
	return &TelegramClient{
		token:      token,
		channel:    channel,
		httpClient: &http.Client{Transport: transport, Timeout: 10 * time.Second},
	}, nil
}

// Send renders req and calls the Telegram sendMessage API.
func (c *TelegramClient) Send(req SendRequest) error {
	text := RenderMessage(req)
	payload := map[string]any{
		"chat_id":              c.channel,
		"text":                 text,
		"disable_notification": req.DisableNotification,
	}
	if req.ParseMode != "" {
		payload["parse_mode"] = req.ParseMode
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", c.token)

	resp, err := c.httpClient.Post(apiURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("http post: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		var result struct {
			Description string `json:"description"`
		}
		if jsonErr := json.Unmarshal(respBody, &result); jsonErr != nil || result.Description == "" {
			return fmt.Errorf("telegram api error (status %d): %s", resp.StatusCode, strings.TrimSpace(string(respBody)))
		}
		return fmt.Errorf("telegram api error: %s", result.Description)
	}
	return nil
}
