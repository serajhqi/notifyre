package main

import (
	"os"
	"testing"
)

func writeTempKeys(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "keys-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Remove(f.Name()) })
	f.WriteString(content)
	f.Close()
	return f.Name()
}

func TestLoadKeys_Valid(t *testing.T) {
	path := writeTempKeys(t, `
keys:
  - key: "abc123"
    label: "service-a"
  - key: "xyz789"
    label: "service-b"
`)
	store, err := LoadKeys(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !store.Valid("abc123") {
		t.Error("expected abc123 to be valid")
	}
	if store.Valid("unknown") {
		t.Error("expected unknown to be invalid")
	}
	if store.Label("abc123") != "service-a" {
		t.Errorf("Label = %q, want %q", store.Label("abc123"), "service-a")
	}
}

func TestLoadKeys_EmptyFile(t *testing.T) {
	path := writeTempKeys(t, "keys: []\n")
	_, err := LoadKeys(path)
	if err == nil {
		t.Error("expected error for empty keys file")
	}
}

func TestLoadKeys_EmptyKey(t *testing.T) {
	path := writeTempKeys(t, `
keys:
  - key: ""
    label: "bad"
`)
	_, err := LoadKeys(path)
	if err == nil {
		t.Error("expected error for entry with empty key")
	}
}

func TestLoadKeys_FileNotFound(t *testing.T) {
	_, err := LoadKeys("/nonexistent/path/keys.yaml")
	if err == nil {
		t.Error("expected error for missing file")
	}
}
