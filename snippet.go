package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newSnippetCmd() *cobra.Command {
	var (
		lang      string
		serverURL string
		apiKey    string
	)

	cmd := &cobra.Command{
		Use:   "snippet",
		Short: "Print a copy-paste integration snippet",
		Long:  "Prints ready-to-use code for integrating notifyre into a project. Pipe or copy into your repo.",
		RunE: func(_ *cobra.Command, _ []string) error {
			snippet, err := getSnippet(lang, serverURL, apiKey)
			if err != nil {
				return err
			}
			fmt.Println(snippet)
			return nil
		},
	}

	cmd.Flags().StringVar(&lang, "lang", "", "Language: curl, bash, go, node, python")
	cmd.Flags().StringVar(&serverURL, "url", "http://localhost:8080", "notifyre server URL")
	cmd.Flags().StringVar(&apiKey, "key", "", "API key")
	cmd.MarkFlagRequired("lang")
	cmd.MarkFlagRequired("key")

	return cmd
}

func getSnippet(lang, serverURL, apiKey string) (string, error) {
	switch lang {
	case "curl":
		return fmt.Sprintf(`curl -s -X POST %s/send \
  -H "Content-Type: application/json" \
  -H "X-API-Key: %s" \
  -d '{"message": "Hello", "level": "info"}'`, serverURL, apiKey), nil

	case "bash":
		return fmt.Sprintf(`#!/usr/bin/env bash
# notifyre helper — paste this function into your script

notify() {
  local message="$1"
  local level="${2:-info}"
  local title="${3:-}"
  curl -s -X POST %s/send \
    -H "Content-Type: application/json" \
    -H "X-API-Key: %s" \
    -d "{\"message\": \"$message\", \"level\": \"$level\", \"title\": \"$title\"}" \
    > /dev/null
}

# Usage:
# notify "Build started"
# notify "Build succeeded" success "CI Pipeline"
# notify "Build failed" error "CI Pipeline"`, serverURL, apiKey), nil

	case "go":
		return fmt.Sprintf(`package main

import (
	"bytes"
	"encoding/json"
	"net/http"
)

const notifyreURL = "%s/send"
const notifyreKey = "%s"

func notify(message, level, title string) error {
	payload, _ := json.Marshal(map[string]string{
		"message": message,
		"level":   level,
		"title":   title,
	})
	req, _ := http.NewRequest("POST", notifyreURL, bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", notifyreKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

// Usage:
// notify("Build started", "info", "My App")
// notify("Build succeeded", "success", "My App")
// notify("Build failed", "error", "My App")`, serverURL, apiKey), nil

	case "node":
		return fmt.Sprintf(`// notifyre integration
const NOTIFYRE_URL = '%s/send';
const NOTIFYRE_KEY = '%s';

async function notify(message, level = 'info', title = '') {
  await fetch(NOTIFYRE_URL, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'X-API-Key': NOTIFYRE_KEY,
    },
    body: JSON.stringify({ message, level, title }),
  });
}

// Usage:
// await notify('Build started')
// await notify('Build succeeded', 'success', 'My App')
// await notify('Build failed', 'error', 'My App')`, serverURL, apiKey), nil

	case "python":
		return fmt.Sprintf(`import requests

NOTIFYRE_URL = "%s/send"
NOTIFYRE_KEY = "%s"

def notify(message: str, level: str = "info", title: str = "") -> None:
    requests.post(
        NOTIFYRE_URL,
        headers={"X-API-Key": NOTIFYRE_KEY},
        json={"message": message, "level": level, "title": title},
    )

# Usage:
# notify("Build started")
# notify("Build succeeded", level="success", title="My App")
# notify("Build failed", level="error", title="My App")`, serverURL, apiKey), nil

	default:
		return "", fmt.Errorf("unsupported language %q — choose: go, node, python, bash, curl", lang)
	}
}
