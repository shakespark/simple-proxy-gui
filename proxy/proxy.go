package proxy

import (
	"encoding/base64"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/golang/glog"
	netProxy "golang.org/x/net/proxy"
)

func NewProxyHandler(timeoutSeconds int) *ProxyHandler {
	return &ProxyHandler{
		Timeout: time.Duration(timeoutSeconds) * time.Second,
	}
}

type ProxyHandler struct {
	Timeout       time.Duration
	Username      *string
	Password      *string
	LogAuth       bool
	LogHeaders    bool
	Socks5Forward *Socks5Forward
	Whitelist     *Whitelist
	Monitor       *Monitor
}

type Socks5Forward struct {
	Address  string
	Username *string
	Password *string
}

func (p *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	target := formatTarget(r)
	clientIP, allowed := p.Whitelist.Allow(r.RemoteAddr)

	glog.V(1).Infof("Serving '%s' request from '%s' to '%s'\n", r.Method, r.RemoteAddr, r.Host)
	if p.LogHeaders {
		for name, values := range r.Header {
			for i, value := range values {
				glog.V(1).Infof("'%s': [%d] %s", name, i, value)
			}
		}
	}

	if !allowed {
		glog.V(1).Infof("Blocked by whitelist: %s -> %s\n", clientIP, target)
		p.Monitor.RecordBlocked(clientIP, r.Method, target, "blocked", http.StatusForbidden, "IP not in whitelist")
		http.Error(w, "Forbidden: IP not in whitelist", http.StatusForbidden)
		return
	}

	if p.Username != nil && p.Password != nil {
		username, password, ok := proxyBasicAuth(r)
		if !ok || username != *p.Username || password != *p.Password {
			if p.LogAuth {
				glog.Errorf("Unauthorized, username: %s, password: %s\n", username, password)
			} else {
				glog.Errorln("Unauthorized")
			}
			p.Monitor.RecordBlocked(clientIP, r.Method, target, "auth_failed", http.StatusProxyAuthRequired, "Proxy authentication failed")
			w.Header().Set("Proxy-Authenticate", "Basic")
			http.Error(w, "Unauthorized", http.StatusProxyAuthRequired)
			return
		}
	}

	entry := p.Monitor.StartRequest(clientIP, r.Method, target)
	if r.Method == http.MethodConnect {
		handleTunneling(w, r, p.Timeout, p.Socks5Forward, p.Monitor, entry)
	} else {
		handleHTTP(w, r, p.Monitor, entry)
	}
}

func handleTunneling(w http.ResponseWriter, r *http.Request, timeout time.Duration, socks5Forward *Socks5Forward, monitor *Monitor, entry *AccessEntry) {
	var destConn net.Conn
	var err error

	if socks5Forward == nil {
		destConn, err = net.DialTimeout("tcp", r.Host, timeout)
	} else {
		var socks5Auth *netProxy.Auth
		if socks5Forward.Username != nil && socks5Forward.Password != nil {
			socks5Auth = &netProxy.Auth{
				User:     *socks5Forward.Username,
				Password: *socks5Forward.Password,
			}
		}

		var socks5Dialer netProxy.Dialer
		socks5Dialer, err = netProxy.SOCKS5("tcp", socks5Forward.Address, socks5Auth, &net.Dialer{
			Timeout:   timeout,
			KeepAlive: 30 * time.Second,
		})

		if err != nil {
			glog.Errorf("Failed to dial socks5 proxy %s, %s\n", socks5Forward.Address, err.Error())
			entry.Finish("error", http.StatusServiceUnavailable, err.Error())
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}

		destConn, err = socks5Dialer.Dial("tcp", r.Host)
	}

	if err != nil {
		glog.Errorf("Failed to dial host, %s\n", err.Error())
		entry.Finish("error", http.StatusServiceUnavailable, err.Error())
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	hijacker, ok := w.(http.Hijacker)
	if !ok {
		destConn.Close()
		glog.Errorln("Attempted to hijack connection that does not support it")
		entry.Finish("error", http.StatusInternalServerError, "Hijacking not supported")
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		destConn.Close()
		glog.Errorf("Failed to hijack connection, %s\n", err.Error())
		entry.Finish("error", http.StatusServiceUnavailable, err.Error())
		return
	}

	monitor.IncActiveConns(1)
	var wg sync.WaitGroup
	wg.Add(2)
	go transfer(destConn, clientConn, entry.AddUp, &wg)
	go transfer(clientConn, destConn, entry.AddDown, &wg)

	go func() {
		wg.Wait()
		monitor.IncActiveConns(-1)
		entry.Finish("allowed", http.StatusOK, "")
	}()
}

func transfer(destination io.WriteCloser, source io.ReadCloser, onBytes func(int64), wg *sync.WaitGroup) {
	defer wg.Done()
	defer destination.Close()
	defer source.Close()
	buf := make([]byte, 32*1024)
	io.CopyBuffer(destination, &countingReader{Reader: source, onRead: onBytes}, buf)
}

func handleHTTP(w http.ResponseWriter, req *http.Request, monitor *Monitor, entry *AccessEntry) {
	monitor.IncActiveConns(1)
	defer monitor.IncActiveConns(-1)

	entry.AddUp(estimateRequestBytes(req))
	if req.Body != nil {
		req.Body = &countingReadCloser{
			ReadCloser: req.Body,
			onRead:     entry.AddUp,
		}
	}

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		glog.Errorf("Failed to proxy request, %s\n", err.Error())
		entry.Finish("error", http.StatusServiceUnavailable, err.Error())
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	entry.AddDown(estimateResponseHeaderBytes(resp))
	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	buf := make([]byte, 32*1024)
	io.CopyBuffer(w, &countingReader{Reader: resp.Body, onRead: entry.AddDown}, buf)
	entry.Finish("allowed", resp.StatusCode, "")
}

type countingReader struct {
	io.Reader
	onRead func(int64)
}

func (c *countingReader) Read(p []byte) (int, error) {
	n, err := c.Reader.Read(p)
	if n > 0 && c.onRead != nil {
		c.onRead(int64(n))
	}
	return n, err
}

type countingReadCloser struct {
	io.ReadCloser
	onRead func(int64)
}

func (c *countingReadCloser) Read(p []byte) (int, error) {
	n, err := c.ReadCloser.Read(p)
	if n > 0 && c.onRead != nil {
		c.onRead(int64(n))
	}
	return n, err
}

func formatTarget(r *http.Request) string {
	if r.Method == http.MethodConnect {
		return r.Host
	}
	if r.URL != nil && r.URL.String() != "" {
		if r.URL.IsAbs() {
			return r.URL.String()
		}
		if r.Host != "" {
			return "http://" + r.Host + r.URL.RequestURI()
		}
		return r.URL.String()
	}
	return r.Host
}

func estimateRequestBytes(r *http.Request) int64 {
	total := int64(len(r.Method) + len(r.RequestURI) + len(r.Proto) + 4)
	for k, vv := range r.Header {
		for _, v := range vv {
			total += int64(len(k) + len(v) + 4)
		}
	}
	return total
}

func estimateResponseHeaderBytes(resp *http.Response) int64 {
	total := int64(len(resp.Proto) + len(resp.Status) + 4)
	for k, vv := range resp.Header {
		for _, v := range vv {
			total += int64(len(k) + len(v) + 4)
		}
	}
	return total
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func proxyBasicAuth(r *http.Request) (username, password string, ok bool) {
	auth := r.Header.Get("Proxy-Authorization")
	if auth == "" {
		return
	}
	return parseBasicAuth(auth)
}

// parseBasicAuth parses an HTTP Basic Authentication string.
// "Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==" returns ("Aladdin", "open sesame", true).
func parseBasicAuth(auth string) (username, password string, ok bool) {
	const prefix = "Basic "
	if len(auth) < len(prefix) || !equalFold(auth[:len(prefix)], prefix) {
		return
	}
	c, err := base64.StdEncoding.DecodeString(auth[len(prefix):])
	if err != nil {
		return
	}
	cs := string(c)
	s := strings.IndexByte(cs, ':')
	if s < 0 {
		return
	}
	return cs[:s], cs[s+1:], true
}

// EqualFold is strings.EqualFold, ASCII only. It reports whether s and t
// are equal, ASCII-case-insensitively.
func equalFold(s, t string) bool {
	if len(s) != len(t) {
		return false
	}
	for i := 0; i < len(s); i++ {
		if lower(s[i]) != lower(t[i]) {
			return false
		}
	}
	return true
}

// lower returns the ASCII lowercase version of b.
func lower(b byte) byte {
	if 'A' <= b && b <= 'Z' {
		return b + ('a' - 'A')
	}
	return b
}
