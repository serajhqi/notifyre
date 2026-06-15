package config

import (
	"testing"
)

func TestLoadConfig_AllVars(t *testing.T) {
	t.Setenv("BOT_TOKEN", "test-token")
	t.Setenv("CHANNEL", "@testchan")
	t.Setenv("PROXY_ADDR", "socks5://localhost:1080")
	t.Setenv("KEYS_FILE", "/tmp/keys.yaml")
	t.Setenv("PORT", "9090")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.BotToken != "test-token" {
		t.Errorf("BotToken = %q, want %q", cfg.BotToken, "test-token")
	}
	if cfg.ProxyAddr != "socks5://localhost:1080" {
		t.Errorf("ProxyAddr = %q, want %q", cfg.ProxyAddr, "socks5://localhost:1080")
	}
	if cfg.Port != "9090" {
		t.Errorf("Port = %q, want %q", cfg.Port, "9090")
	}
	if cfg.Channel != "@testchan" {
		t.Errorf("Channel = %q, want %q", cfg.Channel, "@testchan")
	}
	if cfg.KeysFile != "/tmp/keys.yaml" {
		t.Errorf("KeysFile = %q, want %q", cfg.KeysFile, "/tmp/keys.yaml")
	}
}

func TestLoadConfig_DefaultPort(t *testing.T) {
	t.Setenv("BOT_TOKEN", "tok")
	t.Setenv("CHANNEL", "@chan")
	t.Setenv("KEYS_FILE", "/tmp/keys.yaml")
	t.Setenv("PORT", "")
	t.Setenv("PROXY_ADDR", "")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Port != "8080" {
		t.Errorf("Port = %q, want default %q", cfg.Port, "8080")
	}
	if cfg.ProxyAddr != "" {
		t.Errorf("ProxyAddr = %q, want empty", cfg.ProxyAddr)
	}
}

func TestLoadConfig_MissingRequired(t *testing.T) {
	cases := []struct {
		name  string
		unset string
	}{
		{"missing BOT_TOKEN", "BOT_TOKEN"},
		{"missing CHANNEL", "CHANNEL"},
		{"missing KEYS_FILE", "KEYS_FILE"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("BOT_TOKEN", "tok")
			t.Setenv("CHANNEL", "@chan")
			t.Setenv("KEYS_FILE", "/tmp/keys.yaml")
			t.Setenv(tc.unset, "")

			_, err := LoadConfig()
			if err == nil {
				t.Error("expected error, got nil")
			}
		})
	}
}
