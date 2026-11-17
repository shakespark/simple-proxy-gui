package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/netip"
	"os"
	"path/filepath"
	"strings"
)

// WhitelistRule represents a single IP or CIDR entry in the IP whitelist.
type WhitelistRule struct {
	CIDR    string `json:"cidr"`
	Note    string `json:"note"`
	Enabled bool   `json:"enabled"`
}

// Config holds persisted proxy and monitoring settings.
type Config struct {
	Language          string          `json:"language"`
	Bind              string          `json:"bind"`
	Port              string          `json:"port"`
	TimeoutSeconds    int             `json:"timeout_seconds"`
	AutoStartProxy    bool            `json:"auto_start_proxy"`
	BasicAuthEnabled  bool            `json:"basic_auth_enabled"`
	BasicAuthUsername string          `json:"basic_auth_username"`
	BasicAuthPassword string          `json:"basic_auth_password"`
	Socks5Address     string          `json:"socks5_address"`
	Socks5Username    string          `json:"socks5_username"`
	Socks5Password    string          `json:"socks5_password"`
	WhitelistEnabled  bool            `json:"whitelist_enabled"`
	Whitelist         []WhitelistRule `json:"whitelist"`
	MaxLogEntries     int             `json:"max_log_entries"`
}

// DefaultConfig returns sensible default settings for first-time startup.
func DefaultConfig() *Config {
	return &Config{
		Bind:           "0.0.0.0",
		Port:           "8888",
		TimeoutSeconds: 10,
		AutoStartProxy: true,
		Whitelist: []WhitelistRule{
			{CIDR: "127.0.0.1/32", Note: "Loopback (IPv4)", Enabled: true},
			{CIDR: "::1/128", Note: "Loopback (IPv6)", Enabled: true},
		},
		MaxLogEntries: 2000,
	}
}

// ResolveConfigPath returns config.json next to the executable, or in the current
// working directory when running via `go run` (where the binary lives in TempDir).
func ResolveConfigPath() string {
	exe, err := os.Executable()
	if err != nil {
		return "config.json"
	}
	exeDir := filepath.Dir(exe)
	tempDir := filepath.Clean(os.TempDir())
	if strings.HasPrefix(strings.ToLower(filepath.Clean(exeDir)), strings.ToLower(tempDir)) {
		return "config.json"
	}
	return filepath.Join(exeDir, "config.json")
}

// Load reads configuration from disk, creating a default file if it does not exist.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			cfg := DefaultConfig()
			_ = Save(path, cfg)
			return cfg, nil
		}
		return nil, fmt.Errorf("read config %s: %w", path, err)
	}

	cfg := DefaultConfig()
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parse config %s: %w", path, err)
	}
	normalizeConfig(cfg)
	return cfg, nil
}

// Save writes configuration to disk as formatted JSON.
func Save(path string, cfg *Config) error {
	if cfg == nil {
		return errors.New("nil config")
	}
	normalizeConfig(cfg)
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}
	data = append(data, '\n')

	dir := filepath.Dir(path)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("create config dir: %w", err)
		}
	}

	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0o644); err != nil {
		return fmt.Errorf("write temp config: %w", err)
	}
	if err := os.Rename(tmpPath, path); err != nil {
		// Fallback for Windows if target exists and Rename fails across states.
		if writeErr := os.WriteFile(path, data, 0o644); writeErr != nil {
			_ = os.Remove(tmpPath)
			return fmt.Errorf("write config: %w", writeErr)
		}
		_ = os.Remove(tmpPath)
	}
	return nil
}

func normalizeConfig(cfg *Config) {
	if strings.TrimSpace(cfg.Bind) == "" {
		cfg.Bind = "0.0.0.0"
	}
	if strings.TrimSpace(cfg.Port) == "" {
		cfg.Port = "8888"
	}
	if cfg.TimeoutSeconds <= 0 {
		cfg.TimeoutSeconds = 10
	}
	if cfg.MaxLogEntries <= 0 {
		cfg.MaxLogEntries = 2000
	}
	if cfg.Whitelist == nil {
		cfg.Whitelist = []WhitelistRule{}
	}
	if cfg.Language != "" {
		cfg.Language = NormalizeLanguage(cfg.Language)
	}
}

// NormalizeLanguage maps a language code to "en" or "zh".
func NormalizeLanguage(lang string) string {
	l := strings.ToLower(strings.TrimSpace(lang))
	if strings.HasPrefix(l, "zh") {
		return "zh"
	}
	return "en"
}

// ValidateWhitelistInput parses an IP address or CIDR prefix and returns a normalized representation.
func ValidateWhitelistInput(raw string) (string, error) {
	return ValidateWhitelistInputLang(raw, "zh")
}

// ValidateWhitelistInputLang parses an IP address or CIDR prefix with localized errors.
func ValidateWhitelistInputLang(raw, lang string) (string, error) {
	en := NormalizeLanguage(lang) == "en"
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		if en {
			return "", errors.New("IP or CIDR cannot be empty")
		}
		return "", errors.New("IP 或网段不能为空")
	}
	if strings.Contains(trimmed, "/") {
		prefix, err := netip.ParsePrefix(trimmed)
		if err != nil {
			if en {
				return "", fmt.Errorf("invalid CIDR prefix: %s", trimmed)
			}
			return "", fmt.Errorf("无效的 CIDR 网段: %s", trimmed)
		}
		return prefix.Masked().String(), nil
	}
	addr, err := netip.ParseAddr(trimmed)
	if err != nil {
		if en {
			return "", fmt.Errorf("invalid IP address: %s", trimmed)
		}
		return "", fmt.Errorf("无效的 IP 地址: %s", trimmed)
	}
	return addr.Unmap().String(), nil
}
