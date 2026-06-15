package main

import (
	"bytes"
	"encoding/json"
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
	ks := keys.NewKeyStore(map[string]string{})
	mock := &mockSender{}
	h := sendHandler(ks, mock)

	req := httptest.NewRequest(http.MethodPost, "/send", bytes.NewBufferString(`{"message":"hi"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestSendHandler_WrongAPIKey(t *testing.T) {
	ks := keys.NewKeyStore(map[string]string{"valid": "test"})
	mock := &mockSender{}
	h := sendHandler(ks, mock)

	req := httptest.NewRequest(http.MethodPost, "/send", bytes.NewBufferString(`{"message":"hi"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", "wrong")
	w := httptest.NewRecorder()
	h(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestSendHandler_EmptyMessage(t *testing.T) {
	ks := keys.NewKeyStore(map[string]string{"valid": "test"})
	mock := &mockSender{}
	h := sendHandler(ks, mock)

	req := httptest.NewRequest(http.MethodPost, "/send", bytes.NewBufferString(`{"message":""}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", "valid")
	w := httptest.NewRecorder()
	h(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestSendHandler_TelegramError(t *testing.T) {
	ks := keys.NewKeyStore(map[string]string{"valid": "test"})
	mock := &mockSender{err: errors.New("bot blocked")}
	h := sendHandler(ks, mock)

	req := httptest.NewRequest(http.MethodPost, "/send", bytes.NewBufferString(`{"message":"hi"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", "valid")
	w := httptest.NewRecorder()
	h(w, req)

	if w.Code != http.StatusBadGateway {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadGateway)
	}
}

func TestSendHandler_OK(t *testing.T) {
	ks := keys.NewKeyStore(map[string]string{"valid": "test"})
	mock := &mockSender{}
	h := sendHandler(ks, mock)

	body := `{"message":"hello","level":"success","title":"CI"}`
	req := httptest.NewRequest(http.MethodPost, "/send", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", "valid")
	w := httptest.NewRecorder()
	h(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	var resp apiResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if !resp.OK {
		t.Errorf("resp.OK = false, want true")
	}
	if mock.lastReq.Message != "hello" {
		t.Errorf("forwarded message = %q, want %q", mock.lastReq.Message, "hello")
	}
	if mock.lastReq.Level != "success" {
		t.Errorf("forwarded level = %q, want %q", mock.lastReq.Level, "success")
	}
	if mock.lastReq.Title != "CI" {
		t.Errorf("forwarded title = %q, want %q", mock.lastReq.Title, "CI")
	}
}
