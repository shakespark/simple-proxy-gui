package proxy

import (
	"net"
	"net/netip"
	"strings"
	"sync"

	"github.com/shakespark/simple-proxy-gui/config"
)

// Whitelist evaluates whether client IP addresses are allowed to use the proxy.
type Whitelist struct {
	mu       sync.RWMutex
	enabled  bool
	prefixes []netip.Prefix
}

// NewWhitelist builds a Whitelist from config settings.
func NewWhitelist(enabled bool, rules []config.WhitelistRule) *Whitelist {
	w := &Whitelist{}
	w.Update(enabled, rules)
	return w
}

// Update replaces the active whitelist state and compiled CIDR rules atomically.
func (w *Whitelist) Update(enabled bool, rules []config.WhitelistRule) {
	compiled := make([]netip.Prefix, 0, len(rules))
	for _, rule := range rules {
		if !rule.Enabled {
			continue
		}
		if prefix, ok := parseRulePrefix(rule.CIDR); ok {
			compiled = append(compiled, prefix)
		}
	}

	w.mu.Lock()
	w.enabled = enabled
	w.prefixes = compiled
	w.mu.Unlock()
}

// Allow extracts the normalized client IP from remoteAddr and checks whether it is permitted.
func (w *Whitelist) Allow(remoteAddr string) (string, bool) {
	ipStr, addr, ok := ExtractClientIP(remoteAddr)
	if w == nil {
		return ipStr, true
	}

	w.mu.RLock()
	enabled := w.enabled
	prefixes := w.prefixes
	w.mu.RUnlock()

	if !enabled {
		return ipStr, true
	}
	if !ok {
		return ipStr, false
	}
	for _, p := range prefixes {
		if p.Contains(addr) {
			return ipStr, true
		}
	}
	return ipStr, false
}

// ExtractClientIP parses the IP portion of a network address (host:port or raw IP)
// and unmaps IPv4-mapped IPv6 addresses.
func ExtractClientIP(remoteAddr string) (string, netip.Addr, bool) {
	host := strings.TrimSpace(remoteAddr)
	if h, _, err := net.SplitHostPort(host); err == nil {
		host = h
	}
	host = strings.TrimPrefix(host, "[")
	host = strings.TrimSuffix(host, "]")

	addr, err := netip.ParseAddr(host)
	if err != nil {
		return host, netip.Addr{}, false
	}
	unmapped := addr.Unmap()
	return unmapped.String(), unmapped, true
}

func parseRulePrefix(raw string) (netip.Prefix, bool) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return netip.Prefix{}, false
	}
	if strings.Contains(trimmed, "/") {
		prefix, err := netip.ParsePrefix(trimmed)
		if err != nil {
			return netip.Prefix{}, false
		}
		return prefix.Masked(), true
	}
	addr, err := netip.ParseAddr(trimmed)
	if err != nil {
		return netip.Prefix{}, false
	}
	unmapped := addr.Unmap()
	return netip.PrefixFrom(unmapped, unmapped.BitLen()), true
}
