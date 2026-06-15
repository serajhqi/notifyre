package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"
	"notifyre/internal/config"
)

// TelegramSender is satisfied by TelegramClient and can be mocked in tests.
type TelegramSender interface {
	Send(req SendRequest) error
}

type apiResponse struct {
	OK    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("writeJSON encode: %v", err)
	}
}

func sendHandler(keys *KeyStore, tg TelegramSender) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")
		if !keys.Valid(apiKey) {
			writeJSON(w, http.StatusUnauthorized, apiResponse{Error: "invalid api key"})
			return
		}
		var req SendRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, apiResponse{Error: "invalid request body"})
			return
		}
		if req.Message == "" {
			writeJSON(w, http.StatusBadRequest, apiResponse{Error: "message is empty"})
			return
		}
		if err := tg.Send(req); err != nil {
			writeJSON(w, http.StatusBadGateway, apiResponse{Error: fmt.Sprintf("telegram error: %s", err)})
			return
		}
		writeJSON(w, http.StatusOK, apiResponse{OK: true})
	}
}

func newServeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Start the HTTP notification server",
		RunE:  runServe,
	}
}

func runServe(_ *cobra.Command, _ []string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}
	keys, err := LoadKeys(cfg.KeysFile)
	if err != nil {
		return err
	}
	tg, err := NewTelegramClient(cfg.BotToken, cfg.Channel, cfg.ProxyAddr)
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /send", sendHandler(keys, tg))

	addr := ":" + cfg.Port
	log.Printf("notifyre listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}
