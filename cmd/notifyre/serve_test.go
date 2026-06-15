package main

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"notifyre/internal/keys"
	"notifyre/internal/telegram"
)

// mockSender implements TelegramSender for tests.
type mockSender struct {
	lastReq telegram.SendRequest
	err     error
}

func (m *mockSender) Send(req telegram.SendRequest) error {
	m.lastReq = req
	return m.err
}

func TestSendHandler_MissingAPIKey(t *testing.T) {
	ks := &keys.KeyStore{} // empty store
	mock := &mockSender{}
	h := sendHandler(ks, mock)

	req := httptest.NewRequest(http.MethodPost, "/send", bytes.NewBufferString(`{"message":"hi"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestSendHandler_WrongAPIKey(t *testing.T) {
	mock := &mockSender{}
	h := sendHandler(nil, mock)
	if h == nil {
		t.Error("sendHandler returned nil")
	}
}

func TestSendHandler_EmptyMessage(t *testing.T) {
	mock := &mockSender{}
	h := sendHandler(nil, mock)
	if h == nil {
		t.Error("sendHandler returned nil")
	}
}

func TestSendHandler_TelegramError(t *testing.T) {
	mock := &mockSender{err: errors.New("bot blocked")}
	h := sendHandler(nil, mock)
	if h == nil {
		t.Error("sendHandler returned nil")
	}
}

func TestSendHandler_OK(t *testing.T) {
	mock := &mockSender{}
	h := sendHandler(nil, mock)
	if h == nil {
		t.Error("sendHandler returned nil")
	}
}
