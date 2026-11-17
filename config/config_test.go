package config

import (
	"path/filepath"
	"testing"
)

func TestSaveAndLoadRoundTrip(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.json")

	cfg, err := Load(cfgPath)
	if err != nil {
		t.Fatalf("Load initial config failed: %v", err)
	}
	if cfg.Port != "8888" {
		t.Fatalf("expected default port 8888, got %s", cfg.Port)
	}

	cfg.Port = "9090"
	cfg.WhitelistEnabled = true
	cfg.Whitelist = append(cfg.Whitelist, WhitelistRule{
		CIDR:    "192.168.10.0/24",
		Note:    "LAN",
		Enabled: true,
	})
	if err := Save(cfgPath, cfg); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	reloaded, err := Load(cfgPath)
	if err != nil {
		t.Fatalf("Reload failed: %v", err)
	}
	if reloaded.Port != "9090" || !reloaded.WhitelistEnabled {
		t.Fatalf("unexpected reloaded config: %+v", reloaded)
	}
	if len(reloaded.Whitelist) != 3 || reloaded.Whitelist[2].CIDR != "192.168.10.0/24" {
		t.Fatalf("unexpected whitelist: %+v", reloaded.Whitelist)
	}
}

func TestValidateWhitelistInput(t *testing.T) {
	cases := []struct {
		in      string
		want    string
		wantErr bool
	}{
		{"192.168.1.5", "192.168.1.5", false},
		{" 192.168.1.55/24 ", "192.168.1.0/24", false},
		{"::1", "::1", false},
		{"", "", true},
		{"not-an-ip", "", true},
		{"192.168.1.1/99", "", true},
	}

	for _, tc := range cases {
		got, err := ValidateWhitelistInput(tc.in)
		if tc.wantErr {
			if err == nil {
				t.Errorf("ValidateWhitelistInput(%q) expected error, got %q", tc.in, got)
			}
			continue
		}
		if err != nil {
			t.Errorf("ValidateWhitelistInput(%q) unexpected error: %v", tc.in, err)
			continue
		}
		if got != tc.want {
			t.Errorf("ValidateWhitelistInput(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
