package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

// mockSender implements TelegramSender for tests.
type mockSender struct {
	lastReq SendRequest
	err     error
}

func (m *mockSender) Send(req SendRequest) error {
	m.lastReq = req
	return m.err
}

func TestSendHandler_MissingAPIKey(t *testing.T) {
	store := &KeyStore{lookup: map[string]string{"valid": "svc"}}
	mock := &mockSender{}
	h := sendHandler(store, mock)

	req := httptest.NewRequest(http.MethodPost, "/send", bytes.NewBufferString(`{"message":"hi"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestSendHandler_WrongAPIKey(t *testing.T) {
	store := &KeyStore{lookup: map[string]string{"valid": "svc"}}
	mock := &mockSender{}
	h := sendHandler(store, mock)

	req := httptest.NewRequest(http.MethodPost, "/send", bytes.NewBufferString(`{"message":"hi"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", "wrong")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestSendHandler_EmptyMessage(t *testing.T) {
	store := &KeyStore{lookup: map[string]string{"valid": "svc"}}
	mock := &mockSender{}
	h := sendHandler(store, mock)

	req := httptest.NewRequest(http.MethodPost, "/send", bytes.NewBufferString(`{"message":""}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", "valid")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestSendHandler_TelegramError(t *testing.T) {
	store := &KeyStore{lookup: map[string]string{"valid": "svc"}}
	mock := &mockSender{err: errors.New("bot blocked")}
	h := sendHandler(store, mock)

	req := httptest.NewRequest(http.MethodPost, "/send", bytes.NewBufferString(`{"message":"hi"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", "valid")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusBadGateway {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadGateway)
	}
}

func TestSendHandler_OK(t *testing.T) {
	store := &KeyStore{lookup: map[string]string{"valid": "svc"}}
	mock := &mockSender{}
	h := sendHandler(store, mock)

	body := `{"message":"hello","level":"success","title":"CI"}`
	req := httptest.NewRequest(http.MethodPost, "/send", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", "valid")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	var resp apiResponse
	json.NewDecoder(w.Body).Decode(&resp)
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
