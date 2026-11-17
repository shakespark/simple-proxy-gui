package proxy

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/shakespark/simple-proxy-gui/config"
)

func TestWhitelistAllow(t *testing.T) {
	wl := NewWhitelist(true, []config.WhitelistRule{
		{CIDR: "127.0.0.1/32", Enabled: true},
		{CIDR: "192.168.1.0/24", Enabled: true},
		{CIDR: "10.0.0.5", Enabled: false},
	})

	if ip, ok := wl.Allow("127.0.0.1:54321"); !ok || ip != "127.0.0.1" {
		t.Fatalf("expected 127.0.0.1 allowed, got ip=%s ok=%v", ip, ok)
	}
	if ip, ok := wl.Allow("[::ffff:127.0.0.1]:54321"); !ok || ip != "127.0.0.1" {
		t.Fatalf("expected IPv4-mapped 127.0.0.1 allowed, got ip=%s ok=%v", ip, ok)
	}
	if ip, ok := wl.Allow("192.168.1.42:1234"); !ok || ip != "192.168.1.42" {
		t.Fatalf("expected 192.168.1.42 allowed, got ip=%s ok=%v", ip, ok)
	}
	if _, ok := wl.Allow("10.0.0.5:1234"); ok {
		t.Fatalf("expected disabled rule 10.0.0.5 to be blocked")
	}
	if _, ok := wl.Allow("8.8.8.8:1234"); ok {
		t.Fatalf("expected 8.8.8.8 to be blocked")
	}

	wl.Update(false, nil)
	if _, ok := wl.Allow("8.8.8.8:1234"); !ok {
		t.Fatalf("expected disabled whitelist to allow all IPs")
	}
}

func TestHTTPProxyWhitelistAndMonitoring(t *testing.T) {
	targetSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("hello-proxy-world"))
	}))
	defer targetSrv.Close()

	mon := NewMonitor(100)
	defer mon.Close()
	wl := NewWhitelist(true, []config.WhitelistRule{
		{CIDR: "127.0.0.1/32", Enabled: true},
		{CIDR: "::1/128", Enabled: true},
	})

	handler := &ProxyHandler{
		Timeout:   5 * time.Second,
		Whitelist: wl,
		Monitor:   mon,
	}
	proxySrv := httptest.NewServer(handler)
	defer proxySrv.Close()

	proxyURL, _ := url.Parse(proxySrv.URL)
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	}

	resp, err := client.Get(targetSrv.URL + "/test-path")
	if err != nil {
		t.Fatalf("proxied HTTP GET failed: %v", err)
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK || string(body) != "hello-proxy-world" {
		t.Fatalf("unexpected response: status=%d body=%q", resp.StatusCode, string(body))
	}

	// Now block localhost in whitelist and verify 403 + blocked log
	wl.Update(true, []config.WhitelistRule{
		{CIDR: "10.20.30.40", Enabled: true},
	})
	resp2, err := client.Get(targetSrv.URL + "/blocked-path")
	if err != nil {
		t.Fatalf("blocked HTTP GET request error: %v", err)
	}
	resp2.Body.Close()
	if resp2.StatusCode != http.StatusForbidden {
		t.Fatalf("expected 403 Forbidden, got %d", resp2.StatusCode)
	}

	snap := mon.Snapshot(10)
	if snap.TotalRequests != 2 || snap.BlockedRequests != 1 {
		t.Fatalf("unexpected request counts: total=%d blocked=%d", snap.TotalRequests, snap.BlockedRequests)
	}
	if snap.TotalDownBytes <= 0 || snap.TotalUpBytes <= 0 {
		t.Fatalf("expected non-zero traffic bytes, up=%d down=%d", snap.TotalUpBytes, snap.TotalDownBytes)
	}
	if len(snap.Logs) != 2 {
		t.Fatalf("expected 2 logs, got %d", len(snap.Logs))
	}
	if snap.Logs[0].Status != "blocked" || !strings.Contains(snap.Logs[0].Target, "/blocked-path") {
		t.Fatalf("unexpected newest log: %+v", snap.Logs[0])
	}
	if snap.Logs[1].Status != "allowed" || !strings.Contains(snap.Logs[1].Target, "/test-path") {
		t.Fatalf("unexpected previous log: %+v", snap.Logs[1])
	}
}

func TestConnectTunnelMonitoring(t *testing.T) {
	echoLn, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen echo failed: %v", err)
	}
	defer echoLn.Close()

	go func() {
		conn, err := echoLn.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		_, _ = io.Copy(conn, conn)
	}()

	mon := NewMonitor(100)
	defer mon.Close()
	wl := NewWhitelist(false, nil)

	proxySrv := httptest.NewServer(&ProxyHandler{
		Timeout:   5 * time.Second,
		Whitelist: wl,
		Monitor:   mon,
	})
	defer proxySrv.Close()

	proxyAddr := strings.TrimPrefix(proxySrv.URL, "http://")
	conn, err := net.DialTimeout("tcp", proxyAddr, 3*time.Second)
	if err != nil {
		t.Fatalf("dial proxy failed: %v", err)
	}

	req := fmt.Sprintf("CONNECT %s HTTP/1.1\r\nHost: %s\r\n\r\n", echoLn.Addr().String(), echoLn.Addr().String())
	if _, err := conn.Write([]byte(req)); err != nil {
		t.Fatalf("write CONNECT failed: %v", err)
	}

	br := bufio.NewReader(conn)
	statusLine, err := br.ReadString('\n')
	if err != nil || !strings.Contains(statusLine, "200") {
		t.Fatalf("unexpected CONNECT response %q, err=%v", statusLine, err)
	}
	for {
		line, err := br.ReadString('\n')
		if err != nil {
			t.Fatalf("read header failed: %v", err)
		}
		if line == "\r\n" {
			break
		}
	}

	payload := "tunnel-ping-12345"
	if _, err := conn.Write([]byte(payload)); err != nil {
		t.Fatalf("write payload failed: %v", err)
	}
	buf := make([]byte, len(payload))
	if _, err := io.ReadFull(br, buf); err != nil {
		t.Fatalf("read echo failed: %v", err)
	}
	if string(buf) != payload {
		t.Fatalf("got %q, want %q", string(buf), payload)
	}
	conn.Close()

	// Give tunnel goroutine a brief moment to mark entry finished.
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		snap := mon.Snapshot(10)
		if len(snap.Logs) == 1 && snap.Logs[0].Status == "allowed" {
			if snap.Logs[0].UpBytes < int64(len(payload)) || snap.Logs[0].DownBytes < int64(len(payload)) {
				t.Fatalf("expected tunnel byte counts >= %d, got up=%d down=%d", len(payload), snap.Logs[0].UpBytes, snap.Logs[0].DownBytes)
			}
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("tunnel entry did not finish in time: %+v", mon.Snapshot(10).Logs)
}
