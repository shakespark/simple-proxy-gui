package main

import (
	"context"
	"crypto/tls"
	"encoding/csv"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/shakespark/simple-proxy-gui/config"
	"github.com/shakespark/simple-proxy-gui/proxy"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// ActionResult is a uniform response returned to the frontend for user actions.
type ActionResult struct {
	OK      bool          `json:"ok"`
	Message string        `json:"message,omitempty"`
	Config  config.Config `json:"config"`
	Running bool          `json:"running"`
	Addr    string        `json:"addr"`
}

// DashboardState represents the full initial UI payload.
type DashboardState struct {
	Running    bool                  `json:"running"`
	ListenAddr string                `json:"listen_addr"`
	LastError  string                `json:"last_error,omitempty"`
	ConfigPath string                `json:"config_path"`
	LocalIPs   []string              `json:"local_ips"`
	Config     config.Config         `json:"config"`
	Monitor    proxy.MonitorSnapshot `json:"monitor"`
}

// LiveSnapshot is polled periodically by the frontend to refresh charts and access logs.
type LiveSnapshot struct {
	Running    bool                  `json:"running"`
	ListenAddr string                `json:"listen_addr"`
	LastError  string                `json:"last_error,omitempty"`
	Monitor    proxy.MonitorSnapshot `json:"monitor"`
}

// App manages the proxy lifecycle, configuration persistence, and Wails UI bindings.
type App struct {
	ctx        context.Context
	configPath string

	mu        sync.Mutex
	cfg       *config.Config
	whitelist *proxy.Whitelist
	monitor   *proxy.Monitor

	server   *http.Server
	listener net.Listener
	running  bool
	lastErr  string

	allowForceClose atomic.Bool
	frontendReady   atomic.Bool
	shuttingDownFn  func() bool
}

// NewApp creates an App instance with loaded configuration and initialized monitoring.
func NewApp(cfgPath string) (*App, error) {
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return nil, err
	}
	wl := proxy.NewWhitelist(cfg.WhitelistEnabled, cfg.Whitelist)
	mon := proxy.NewMonitor(cfg.MaxLogEntries)

	return &App{
		configPath:     cfgPath,
		cfg:            cfg,
		whitelist:      wl,
		monitor:        mon,
		shuttingDownFn: isSystemShuttingDown,
	}, nil
}

// Startup is invoked by Wails when the application window initializes.
func (a *App) Startup(ctx context.Context) {
	a.mu.Lock()
	a.ctx = ctx
	autoStart := a.cfg.AutoStartProxy
	a.mu.Unlock()

	if autoStart {
		_ = a.StartProxy()
	}
}

// Shutdown is invoked by Wails when the application exits.
func (a *App) Shutdown(ctx context.Context) {
	a.StopProxy()
	if a.monitor != nil {
		a.monitor.Close()
	}
}

// BeforeClose intercepts window close requests to prompt for confirmation,
// except when Windows is shutting down or the user already confirmed exit.
func (a *App) BeforeClose(ctx context.Context) bool {
	if a.allowForceClose.Load() {
		return false
	}
	if a.shuttingDownFn != nil && a.shuttingDownFn() {
		return false
	}
	if !a.frontendReady.Load() {
		return false
	}

	a.mu.Lock()
	running := a.running
	a.mu.Unlock()

	if ctx != nil {
		go wailsRuntime.EventsEmit(ctx, "confirm-exit", running)
	}
	return true
}

// ConfirmQuit allows the application to close without prompting again.
func (a *App) ConfirmQuit() {
	a.allowForceClose.Store(true)
	a.mu.Lock()
	ctx := a.ctx
	a.mu.Unlock()
	if ctx != nil {
		wailsRuntime.Quit(ctx)
	}
}

// GetState returns the full initial state including configuration, local IPs, and monitor data.
func (a *App) GetState() DashboardState {
	a.frontendReady.Store(true)
	a.mu.Lock()
	cfgCopy := cloneConfig(a.cfg)
	running := a.running
	listenAddr := fmt.Sprintf("%s:%s", a.cfg.Bind, a.cfg.Port)
	lastErr := a.lastErr
	cfgPath := a.configPath
	a.mu.Unlock()

	return DashboardState{
		Running:    running,
		ListenAddr: listenAddr,
		LastError:  lastErr,
		ConfigPath: cfgPath,
		LocalIPs:   detectLocalIPs(),
		Config:     cfgCopy,
		Monitor:    a.monitor.Snapshot(500),
	}
}

// GetLiveSnapshot returns current runtime status, bandwidth history, and recent access logs.
func (a *App) GetLiveSnapshot() LiveSnapshot {
	a.mu.Lock()
	running := a.running
	listenAddr := fmt.Sprintf("%s:%s", a.cfg.Bind, a.cfg.Port)
	lastErr := a.lastErr
	a.mu.Unlock()

	return LiveSnapshot{
		Running:    running,
		ListenAddr: listenAddr,
		LastError:  lastErr,
		Monitor:    a.monitor.Snapshot(500),
	}
}

// SetLanguage updates and persists the UI language ("zh" or "en").
func (a *App) SetLanguage(lang string) ActionResult {
	normalized := config.NormalizeLanguage(lang)

	a.mu.Lock()
	defer a.mu.Unlock()

	a.cfg.Language = normalized
	if err := config.Save(a.configPath, a.cfg); err != nil {
		return a.resultLocked(false, trLang(normalized, fmt.Sprintf("保存配置失败: %v", err), fmt.Sprintf("Failed to save config: %v", err)))
	}
	return a.resultLocked(true, trLang(normalized, "语言已切换为中文", "Language switched to English"))
}

// StartProxy starts the HTTP/HTTPS proxy server using the current configuration.
func (a *App) StartProxy() ActionResult {
	a.mu.Lock()
	defer a.mu.Unlock()

	lang := a.cfg.Language
	if a.running {
		return a.resultLocked(true, trLang(lang, "代理已在运行中", "Proxy is already running"))
	}

	addr := net.JoinHostPort(strings.TrimSpace(a.cfg.Bind), strings.TrimSpace(a.cfg.Port))
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		a.lastErr = trLang(lang, fmt.Sprintf("监听 %s 失败: %v", addr, err), fmt.Sprintf("Failed to listen on %s: %v", addr, err))
		return a.resultLocked(false, a.lastErr)
	}

	handler := a.buildHandlerLocked()
	srv := &http.Server{
		Addr:         addr,
		Handler:      handler,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}

	a.server = srv
	a.listener = ln
	a.running = true
	a.lastErr = ""

	go func(s *http.Server, l net.Listener, curLang string) {
		serveErr := s.Serve(l)
		if serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
			a.mu.Lock()
			if a.server == s {
				a.running = false
				a.lastErr = trLang(curLang, fmt.Sprintf("代理异常退出: %v", serveErr), fmt.Sprintf("Proxy exited unexpectedly: %v", serveErr))
			}
			a.mu.Unlock()
		}
	}(srv, ln, lang)

	return a.resultLocked(true, trLang(lang, fmt.Sprintf("代理已启动 (%s)", addr), fmt.Sprintf("Proxy started (%s)", addr)))
}

// StopProxy gracefully stops the running proxy server.
func (a *App) StopProxy() ActionResult {
	a.mu.Lock()
	srv := a.server
	ln := a.listener
	wasRunning := a.running
	lang := a.cfg.Language
	a.server = nil
	a.listener = nil
	a.running = false
	a.mu.Unlock()

	if wasRunning && srv != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		_ = srv.Shutdown(ctx)
		cancel()
		if ln != nil {
			_ = ln.Close()
		}
	}

	a.mu.Lock()
	defer a.mu.Unlock()
	return a.resultLocked(true, trLang(lang, "代理已停止", "Proxy stopped"))
}

// SaveConfig validates, persists, and applies updated proxy settings.
func (a *App) SaveConfig(newCfg config.Config) ActionResult {
	a.mu.Lock()
	lang := newCfg.Language
	if lang == "" {
		lang = a.cfg.Language
	}
	a.mu.Unlock()

	if err := validateSettings(&newCfg, lang); err != nil {
		a.mu.Lock()
		defer a.mu.Unlock()
		return a.resultLocked(false, err.Error())
	}

	a.mu.Lock()
	if newCfg.Whitelist == nil {
		newCfg.Whitelist = append([]config.WhitelistRule(nil), a.cfg.Whitelist...)
	}
	if newCfg.Language == "" {
		newCfg.Language = a.cfg.Language
	}
	needRestart := a.running && serverConfigChanged(a.cfg, &newCfg)
	a.cfg = cloneConfigPtr(&newCfg)
	a.whitelist.Update(a.cfg.WhitelistEnabled, a.cfg.Whitelist)
	a.monitor.SetMaxLogs(a.cfg.MaxLogEntries)

	if err := config.Save(a.configPath, a.cfg); err != nil {
		a.mu.Unlock()
		return ActionResult{OK: false, Message: trLang(lang, fmt.Sprintf("保存配置失败: %v", err), fmt.Sprintf("Failed to save config: %v", err))}
	}
	a.mu.Unlock()

	if needRestart {
		a.StopProxy()
		res := a.StartProxy()
		if !res.OK {
			return res
		}
		res.Message = trLang(lang, "配置已保存并重启代理", "Settings saved and proxy restarted")
		return res
	}

	a.mu.Lock()
	defer a.mu.Unlock()
	return a.resultLocked(true, trLang(lang, "配置已保存", "Settings saved"))
}

// SetWhitelistEnabled toggles the IP whitelist without restarting the proxy.
func (a *App) SetWhitelistEnabled(enabled bool) ActionResult {
	a.mu.Lock()
	defer a.mu.Unlock()

	lang := a.cfg.Language
	a.cfg.WhitelistEnabled = enabled
	a.whitelist.Update(a.cfg.WhitelistEnabled, a.cfg.Whitelist)
	if err := config.Save(a.configPath, a.cfg); err != nil {
		return a.resultLocked(false, trLang(lang, fmt.Sprintf("保存配置失败: %v", err), fmt.Sprintf("Failed to save config: %v", err)))
	}
	msg := trLang(lang, "已关闭 IP 白名单（允许所有客户端访问）", "IP whitelist disabled (all clients allowed)")
	if enabled {
		msg = trLang(lang, "已启用 IP 白名单", "IP whitelist enabled")
	}
	return a.resultLocked(true, msg)
}

// AddWhitelistRule validates and appends a new IP/CIDR entry to the whitelist.
func (a *App) AddWhitelistRule(cidr string, note string) ActionResult {
	a.mu.Lock()
	lang := a.cfg.Language
	a.mu.Unlock()

	normalized, err := config.ValidateWhitelistInputLang(cidr, lang)
	if err != nil {
		a.mu.Lock()
		defer a.mu.Unlock()
		return a.resultLocked(false, err.Error())
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	for _, existing := range a.cfg.Whitelist {
		if existing.CIDR == normalized {
			return a.resultLocked(false, trLang(lang, fmt.Sprintf("规则 %s 已存在", normalized), fmt.Sprintf("Rule %s already exists", normalized)))
		}
	}

	a.cfg.Whitelist = append(a.cfg.Whitelist, config.WhitelistRule{
		CIDR:    normalized,
		Note:    strings.TrimSpace(note),
		Enabled: true,
	})
	a.whitelist.Update(a.cfg.WhitelistEnabled, a.cfg.Whitelist)
	if err := config.Save(a.configPath, a.cfg); err != nil {
		return a.resultLocked(false, trLang(lang, fmt.Sprintf("保存配置失败: %v", err), fmt.Sprintf("Failed to save config: %v", err)))
	}
	return a.resultLocked(true, trLang(lang, fmt.Sprintf("已添加白名单规则: %s", normalized), fmt.Sprintf("Added whitelist rule: %s", normalized)))
}

// ToggleWhitelistRule enables or disables a single whitelist entry.
func (a *App) ToggleWhitelistRule(cidr string, enabled bool) ActionResult {
	a.mu.Lock()
	defer a.mu.Unlock()

	lang := a.cfg.Language
	found := false
	for i := range a.cfg.Whitelist {
		if a.cfg.Whitelist[i].CIDR == cidr {
			a.cfg.Whitelist[i].Enabled = enabled
			found = true
			break
		}
	}
	if !found {
		return a.resultLocked(false, trLang(lang, "未找到对应规则", "Rule not found"))
	}

	a.whitelist.Update(a.cfg.WhitelistEnabled, a.cfg.Whitelist)
	if err := config.Save(a.configPath, a.cfg); err != nil {
		return a.resultLocked(false, trLang(lang, fmt.Sprintf("保存配置失败: %v", err), fmt.Sprintf("Failed to save config: %v", err)))
	}
	return a.resultLocked(true, trLang(lang, "规则状态已更新", "Rule updated"))
}

// RemoveWhitelistRule deletes a whitelist entry by CIDR string.
func (a *App) RemoveWhitelistRule(cidr string) ActionResult {
	a.mu.Lock()
	defer a.mu.Unlock()

	lang := a.cfg.Language
	filtered := make([]config.WhitelistRule, 0, len(a.cfg.Whitelist))
	removed := false
	for _, r := range a.cfg.Whitelist {
		if r.CIDR == cidr {
			removed = true
			continue
		}
		filtered = append(filtered, r)
	}
	if !removed {
		return a.resultLocked(false, trLang(lang, "未找到对应规则", "Rule not found"))
	}

	a.cfg.Whitelist = filtered
	a.whitelist.Update(a.cfg.WhitelistEnabled, a.cfg.Whitelist)
	if err := config.Save(a.configPath, a.cfg); err != nil {
		return a.resultLocked(false, trLang(lang, fmt.Sprintf("保存配置失败: %v", err), fmt.Sprintf("Failed to save config: %v", err)))
	}
	return a.resultLocked(true, trLang(lang, fmt.Sprintf("已删除规则: %s", cidr), fmt.Sprintf("Removed rule: %s", cidr)))
}

// ClearLogs clears all in-memory access logs.
func (a *App) ClearLogs() ActionResult {
	a.monitor.ClearLogs()
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.resultLocked(true, trLang(a.cfg.Language, "访问记录已清空", "Access logs cleared"))
}

// ExportLogsCSV opens a native save dialog and writes current access logs to a CSV file.
func (a *App) ExportLogsCSV() ActionResult {
	a.mu.Lock()
	ctx := a.ctx
	lang := a.cfg.Language
	a.mu.Unlock()

	if ctx == nil {
		return ActionResult{OK: false, Message: trLang(lang, "窗口上下文未就绪", "Window context not ready")}
	}

	defaultName := fmt.Sprintf("proxy-logs-%s.csv", time.Now().Format("20060102-150405"))
	savePath, err := wailsRuntime.SaveFileDialog(ctx, wailsRuntime.SaveDialogOptions{
		Title:           trLang(lang, "导出访问记录", "Export Access Logs"),
		DefaultFilename: defaultName,
		Filters: []wailsRuntime.FileFilter{
			{DisplayName: trLang(lang, "CSV 文件 (*.csv)", "CSV Files (*.csv)"), Pattern: "*.csv"},
		},
	})
	if err != nil {
		return ActionResult{OK: false, Message: trLang(lang, fmt.Sprintf("打开保存对话框失败: %v", err), fmt.Sprintf("Failed to open save dialog: %v", err))}
	}
	if savePath == "" {
		return ActionResult{OK: false, Message: trLang(lang, "已取消导出", "Export cancelled")}
	}

	snap := a.monitor.Snapshot(0)
	f, err := os.Create(savePath)
	if err != nil {
		return ActionResult{OK: false, Message: trLang(lang, fmt.Sprintf("创建文件失败: %v", err), fmt.Sprintf("Failed to create file: %v", err))}
	}
	defer f.Close()

	// Write UTF-8 BOM so Windows Excel renders headers properly.
	_, _ = f.Write([]byte{0xEF, 0xBB, 0xBF})

	writer := csv.NewWriter(f)
	if config.NormalizeLanguage(lang) == "en" {
		_ = writer.Write([]string{"Time", "Client IP", "Method", "Target URL", "Status", "Status Code", "Upload Bytes", "Download Bytes", "Duration (ms)", "Error"})
	} else {
		_ = writer.Write([]string{"时间", "客户端IP", "方法", "目标网址", "状态", "状态码", "上行字节", "下行字节", "耗时(ms)", "错误信息"})
	}
	for _, row := range snap.Logs {
		_ = writer.Write([]string{
			row.Time,
			row.ClientIP,
			row.Method,
			row.Target,
			row.Status,
			strconv.Itoa(row.StatusCode),
			strconv.FormatInt(row.UpBytes, 10),
			strconv.FormatInt(row.DownBytes, 10),
			strconv.FormatInt(row.DurationMs, 10),
			row.Error,
		})
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return ActionResult{OK: false, Message: trLang(lang, fmt.Sprintf("写入 CSV 失败: %v", err), fmt.Sprintf("Failed to write CSV: %v", err))}
	}

	a.mu.Lock()
	defer a.mu.Unlock()
	return a.resultLocked(true, trLang(lang, fmt.Sprintf("已导出 %d 条记录到 %s", len(snap.Logs), savePath), fmt.Sprintf("Exported %d records to %s", len(snap.Logs), savePath)))
}

func (a *App) buildHandlerLocked() http.Handler {
	var socks5Forward *proxy.Socks5Forward
	if addr := strings.TrimSpace(a.cfg.Socks5Address); addr != "" {
		socks5Forward = &proxy.Socks5Forward{
			Address: addr,
		}
		if a.cfg.Socks5Username != "" || a.cfg.Socks5Password != "" {
			u := a.cfg.Socks5Username
			p := a.cfg.Socks5Password
			socks5Forward.Username = &u
			socks5Forward.Password = &p
		}
	}

	h := &proxy.ProxyHandler{
		Timeout:       time.Duration(a.cfg.TimeoutSeconds) * time.Second,
		Socks5Forward: socks5Forward,
		Whitelist:     a.whitelist,
		Monitor:       a.monitor,
	}
	if a.cfg.BasicAuthEnabled && a.cfg.BasicAuthUsername != "" {
		u := a.cfg.BasicAuthUsername
		p := a.cfg.BasicAuthPassword
		h.Username = &u
		h.Password = &p
	}
	return h
}

func (a *App) resultLocked(ok bool, message string) ActionResult {
	return ActionResult{
		OK:      ok,
		Message: message,
		Config:  cloneConfig(a.cfg),
		Running: a.running,
		Addr:    fmt.Sprintf("%s:%s", a.cfg.Bind, a.cfg.Port),
	}
}

func trLang(lang, zh, en string) string {
	if config.NormalizeLanguage(lang) == "en" {
		return en
	}
	return zh
}

func validateSettings(cfg *config.Config, lang string) error {
	cfg.Bind = strings.TrimSpace(cfg.Bind)
	if cfg.Bind == "" {
		return errors.New(trLang(lang, "监听地址不能为空", "Bind address cannot be empty"))
	}
	if net.ParseIP(cfg.Bind) == nil {
		return errors.New(trLang(lang, fmt.Sprintf("无效的监听 IP 地址: %s", cfg.Bind), fmt.Sprintf("Invalid bind IP address: %s", cfg.Bind)))
	}

	cfg.Port = strings.TrimSpace(cfg.Port)
	portNum, err := strconv.Atoi(cfg.Port)
	if err != nil || portNum < 1 || portNum > 65535 {
		return errors.New(trLang(lang, fmt.Sprintf("端口号必须在 1 到 65535 之间: %s", cfg.Port), fmt.Sprintf("Port must be between 1 and 65535: %s", cfg.Port)))
	}

	if cfg.TimeoutSeconds < 1 || cfg.TimeoutSeconds > 300 {
		return errors.New(trLang(lang, "超时时间必须在 1 到 300 秒之间", "Timeout must be between 1 and 300 seconds"))
	}
	if cfg.MaxLogEntries < 100 || cfg.MaxLogEntries > 50000 {
		return errors.New(trLang(lang, "最大日志条数必须在 100 到 50000 之间", "Max log entries must be between 100 and 50000"))
	}
	if cfg.BasicAuthEnabled && strings.TrimSpace(cfg.BasicAuthUsername) == "" {
		return errors.New(trLang(lang, "启用代理认证时用户名不能为空", "Username is required when Basic Auth is enabled"))
	}
	return nil
}

func serverConfigChanged(a, b *config.Config) bool {
	return a.Bind != b.Bind ||
		a.Port != b.Port ||
		a.TimeoutSeconds != b.TimeoutSeconds ||
		a.BasicAuthEnabled != b.BasicAuthEnabled ||
		a.BasicAuthUsername != b.BasicAuthUsername ||
		a.BasicAuthPassword != b.BasicAuthPassword ||
		a.Socks5Address != b.Socks5Address ||
		a.Socks5Username != b.Socks5Username ||
		a.Socks5Password != b.Socks5Password
}

func cloneConfig(cfg *config.Config) config.Config {
	out := *cfg
	out.Whitelist = append([]config.WhitelistRule(nil), cfg.Whitelist...)
	return out
}

func cloneConfigPtr(cfg *config.Config) *config.Config {
	c := cloneConfig(cfg)
	return &c
}

func detectLocalIPs() []string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil
	}
	var ips []string
	seen := make(map[string]struct{})
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok || ipNet.IP == nil {
				continue
			}
			v4 := ipNet.IP.To4()
			if v4 == nil {
				continue
			}
			s := v4.String()
			if _, exists := seen[s]; !exists {
				seen[s] = struct{}{}
				ips = append(ips, s)
			}
		}
	}
	return ips
}
