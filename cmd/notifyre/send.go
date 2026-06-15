package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
	"notifyre/internal/telegram"
)

func newSendCmd() *cobra.Command {
	var (
		serverURL           string
		apiKey              string
		message             string
		title               string
		level               string
		parseMode           string
		disableNotification bool
	)

	cmd := &cobra.Command{
		Use:   "send",
		Short: "Send a notification via the running notifyre server",
		RunE: func(_ *cobra.Command, _ []string) error {
			req := telegram.SendRequest{
				Message:             message,
				Title:               title,
				Level:               level,
				ParseMode:           parseMode,
				DisableNotification: disableNotification,
			}
			body, err := json.Marshal(req)
			if err != nil {
				return fmt.Errorf("marshal request: %w", err)
			}

			httpReq, err := http.NewRequest(http.MethodPost, serverURL+"/send", bytes.NewReader(body))
			if err != nil {
				return fmt.Errorf("create request: %w", err)
			}
			httpReq.Header.Set("Content-Type", "application/json")
			httpReq.Header.Set("X-API-Key", apiKey)

			resp, err := http.DefaultClient.Do(httpReq)
			if err != nil {
				return fmt.Errorf("connect to server: %w", err)
			}
			defer resp.Body.Close()

			respBody, err := io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("read response: %w", err)
			}
			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("server error %d: %s", resp.StatusCode, respBody)
			}
			fmt.Println("ok")
			return nil
		},
	}

	cmd.Flags().StringVar(&serverURL, "url", "http://localhost:8080", "notifyre server URL")
	cmd.Flags().StringVar(&apiKey, "key", "", "API key")
	cmd.Flags().StringVar(&message, "message", "", "Notification message")
	cmd.Flags().StringVar(&title, "title", "", "Optional title")
	cmd.Flags().StringVar(&level, "level", "", "Level: info, success, warning, error")
	cmd.Flags().StringVar(&parseMode, "parse-mode", "", "Telegram parse mode: HTML or MarkdownV2")
	cmd.Flags().BoolVar(&disableNotification, "disable-notification", false, "Send silently (no buzz)")

	cmd.MarkFlagRequired("key")
	cmd.MarkFlagRequired("message")

	return cmd
}
