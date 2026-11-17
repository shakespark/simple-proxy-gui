package main

import (
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"testing"

	"github.com/shakespark/simple-proxy-gui/config"
)

func getFreePort(t *testing.T) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("getFreePort failed: %v", err)
	}
	defer ln.Close()
	_, port, _ := net.SplitHostPort(ln.Addr().String())
	return port
}

func TestAppLifecycleAndWhitelistPersistence(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.json")

	app, err := NewApp(cfgPath)
	if err != nil {
		t.Fatalf("NewApp failed: %v", err)
	}
	defer app.Shutdown(nil)

	port := getFreePort(t)
	cfg := app.GetState().Config
	cfg.Bind = "127.0.0.1"
	cfg.Port = port
	res := app.SaveConfig(cfg)
	if !res.OK {
		t.Fatalf("SaveConfig failed: %s", res.Message)
	}

	startRes := app.StartProxy()
	if !startRes.OK || !startRes.Running {
		t.Fatalf("StartProxy failed: %+v", startRes)
	}
	defer app.StopProxy()

	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))
	defer target.Close()

	proxyURL, _ := url.Parse("http://127.0.0.1:" + port)
	client := &http.Client{
		Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)},
	}

	// Whitelist disabled -> allowed
	resp, err := client.Get(target.URL + "/a")
	if err != nil {
		t.Fatalf("GET via proxy failed: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	// Add rule and enable whitelist, then disable 127.0.0.1/32 rule -> blocked
	if r := app.AddWhitelistRule("192.168.50.10", "PC"); !r.OK {
		t.Fatalf("AddWhitelistRule failed: %s", r.Message)
	}
	if r := app.ToggleWhitelistRule("127.0.0.1/32", false); !r.OK {
		t.Fatalf("ToggleWhitelistRule failed: %s", r.Message)
	}
	if r := app.SetWhitelistEnabled(true); !r.OK {
		t.Fatalf("SetWhitelistEnabled failed: %s", r.Message)
	}

	resp2, err := client.Get(target.URL + "/b")
	if err != nil {
		t.Fatalf("blocked GET via proxy failed: %v", err)
	}
	resp2.Body.Close()
	if resp2.StatusCode != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", resp2.StatusCode)
	}

	// Re-enable 127.0.0.1/32 and verify persistence
	if r := app.ToggleWhitelistRule("127.0.0.1/32", true); !r.OK {
		t.Fatalf("re-enable failed: %s", r.Message)
	}
	if r := app.RemoveWhitelistRule("192.168.50.10"); !r.OK {
		t.Fatalf("RemoveWhitelistRule failed: %s", r.Message)
	}

	reloaded, err := config.Load(cfgPath)
	if err != nil {
		t.Fatalf("reload persisted config failed: %v", err)
	}
	if !reloaded.WhitelistEnabled || reloaded.Port != port {
		t.Fatalf("persisted config mismatch: %+v", reloaded)
	}

	live := app.GetLiveSnapshot()
	if !live.Running || live.Monitor.TotalRequests < 2 || live.Monitor.BlockedRequests < 1 {
		t.Fatalf("unexpected live snapshot: %+v", live)
	}

	stopRes := app.StopProxy()
	if !stopRes.OK || stopRes.Running {
		t.Fatalf("StopProxy failed: %+v", stopRes)
	}
}

func TestBeforeCloseBehavior(t *testing.T) {
	dir := t.TempDir()
	app, err := NewApp(filepath.Join(dir, "config.json"))
	if err != nil {
		t.Fatalf("NewApp failed: %v", err)
	}
	defer app.Shutdown(nil)

	// Before frontend is ready, do not block close.
	if app.BeforeClose(nil) {
		t.Fatalf("expected BeforeClose to allow close before frontend ready")
	}

	// Mark frontend ready -> normal close should be intercepted for confirmation.
	_ = app.GetState()
	app.shuttingDownFn = func() bool { return false }
	if !app.BeforeClose(nil) {
		t.Fatalf("expected BeforeClose to intercept normal close")
	}

	// During OS shutdown -> must not prompt / must allow close immediately.
	app.shuttingDownFn = func() bool { return true }
	if app.BeforeClose(nil) {
		t.Fatalf("expected BeforeClose to allow close during OS shutdown")
	}

	// After user confirms quit -> must allow close immediately.
	app.shuttingDownFn = func() bool { return false }
	app.ConfirmQuit()
	if app.BeforeClose(nil) {
		t.Fatalf("expected BeforeClose to allow close after ConfirmQuit")
	}
}

func TestSetLanguageAndLocalization(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.json")

	app, err := NewApp(cfgPath)
	if err != nil {
		t.Fatalf("NewApp failed: %v", err)
	}
	defer app.Shutdown(nil)

	res := app.SetLanguage("en")
	if !res.OK || res.Config.Language != "en" || res.Message != "Language switched to English" {
		t.Fatalf("unexpected en SetLanguage result: %+v", res)
	}
	wlRes := app.SetWhitelistEnabled(true)
	if wlRes.Message != "IP whitelist enabled" {
		t.Fatalf("expected English message, got %q", wlRes.Message)
	}

	resZh := app.SetLanguage("zh-CN")
	if !resZh.OK || resZh.Config.Language != "zh" || resZh.Message != "语言已切换为中文" {
		t.Fatalf("unexpected zh SetLanguage result: %+v", resZh)
	}
	wlResZh := app.SetWhitelistEnabled(false)
	if wlResZh.Message != "已关闭 IP 白名单（允许所有客户端访问）" {
		t.Fatalf("expected Chinese message, got %q", wlResZh.Message)
	}

	reloaded, err := config.Load(cfgPath)
	if err != nil {
		t.Fatalf("reload config failed: %v", err)
	}
	if reloaded.Language != "zh" {
		t.Fatalf("expected persisted language zh, got %q", reloaded.Language)
	}
}
