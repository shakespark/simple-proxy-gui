(() => {
  const I18N = {
    zh: {
      loading: "加载中...",
      initializing: "初始化...",
      lanHintTitle: "其他设备可连接的局域网代理地址",
      configLabel: "配置文件",
      runningStatus: "运行中",
      stoppedStatus: "已停止",
      startProxy: "启动代理",
      stopProxy: "停止代理",
      lanProxyPrefix: "局域网代理",
      proxyAddrPrefix: "代理地址",
      startupError: "启动异常",
      backendConnectError: "无法连接后端服务",
      backendTimeoutError: "后端通信未就绪",

      tabMonitor: "实时监控",
      tabWhitelist: "IP 白名单",
      tabSettings: "代理设置",

      kpiUpLabel: "实时上行速率",
      kpiTotalUpLabel: "累计上行",
      kpiDownLabel: "实时下行速率",
      kpiTotalDownLabel: "累计下行",
      kpiConnsLabel: "活跃连接 / 总请求",
      kpiConnsSub: "实时并发隧道与请求",
      kpiWlLabel: "IP 白名单",
      kpiBlockedLabel: "已拦截请求",
      wlEnabledCount: "已启用 ({n}条)",
      wlDisabledState: "未启用",

      chartTitle: "实时流量监控（最近 3 分钟）",
      legendUp: "上行速率",
      legendDown: "下行速率",
      tooltipUp: "上行",
      tooltipDown: "下行",

      logsTitle: "IP 访问网址记录",
      logSearchPlaceholder: "搜索客户端 IP 或网址...",
      filterAll: "全部状态",
      filterAllowed: "放行 / 活跃",
      filterBlocked: "白名单拦截",
      filterAuthFailed: "认证失败",
      filterError: "连接错误",
      pauseRefresh: "暂停刷新",
      resumeRefresh: "恢复刷新",
      exportCsv: "导出 CSV",
      clearLogs: "清空",

      colTime: "时间",
      colClientIp: "客户端 IP",
      colMethod: "方法",
      colTarget: "访问网址 / 目标",
      colStatus: "状态",
      colTraffic: "上行 / 下行",
      colDuration: "耗时",
      logsEmpty: "暂无匹配的访问记录",

      quickAllowBtn: "+白名单",
      quickAllowTitle: "加入 IP 白名单",
      quickAllowNote: "从访问记录添加",
      durationActive: "进行中",
      badgeActive: "连接中",
      badgeAllowed: "{code} 放行",
      badgeBlocked: "403 拦截",
      badgeAuth: "407 认证",
      badgeError: "错误",

      wlMasterTitle: "IP 白名单控制",
      wlMasterDesc: "开启后，仅允许下方已启用列表中的客户端 IP 或 CIDR 网段使用本代理；其余来源将返回 403 并记录到拦截日志中。修改实时生效。",
      wlAddTitle: "添加允许的 IP / 网段",
      wlCidrPlaceholder: "IP 或 CIDR，例如 192.168.1.100 或 192.168.1.0/24",
      wlNotePlaceholder: "备注（可选），例如 办公电脑",
      wlAddBtn: "添加规则",
      wlListTitle: "白名单规则列表",
      wlColCidr: "IP / 网段",
      wlColNote: "备注",
      wlColEnabled: "启用",
      wlColAction: "操作",
      wlEmpty: "暂无白名单规则",
      wlDelete: "删除",
      quickFillLabel: "快捷填入:",
      chipLoopback: "本机回环",
      chipLocalIp: "本机网卡 IP",
      chipLocalSubnet: "本机所在局域网段",

      setBasicTitle: "基础代理配置",
      setBindLabel: "监听地址 (Bind)",
      bindPresetTitle: "快速填入网卡地址",
      presetSelect: "快捷选择...",
      presetAllInterfaces: "0.0.0.0 (全部网卡)",
      presetLocalOnly: "127.0.0.1 (仅本机)",
      presetInterfaceSuffix: "本机网卡",
      setPortLabel: "监听端口 (Port)",
      setTimeoutLabel: "连接超时（秒）",
      setMaxLogsLabel: "内存日志保留条数",
      setAutoStartLabel: "打开 Simple Proxy 时自动启动代理服务",
      setAuthTitle: "代理认证 (Basic Auth)",
      setAuthUserLabel: "用户名",
      setAuthPassLabel: "密码",
      setSocksTitle: "上游 SOCKS5 转发（可选）",
      setSocksDesc: "留空则直接连接目标服务器；填写后 HTTPS (CONNECT) 隧道流量将经由此 SOCKS5 代理转发。",
      setSocksAddrLabel: "SOCKS5 地址",
      setSocksAddrPlaceholder: "例如 127.0.0.1:1080",
      setSocksUserLabel: "SOCKS5 用户名",
      setSocksPassLabel: "SOCKS5 密码",
      optionalPlaceholder: "可选",
      saveSettingsBtn: "保存并应用配置",

      exitTitle: "确认退出 Simple Proxy",
      exitMsgRunning: "当前代理服务正在运行中，退出后将中断所有代理连接。确定要退出吗？",
      exitMsgStopped: "确定要退出 Simple Proxy 吗？",
      cancel: "取消",
      confirmExit: "确认退出"
    },
    en: {
      loading: "Loading...",
      initializing: "Initializing...",
      lanHintTitle: "LAN proxy address for other devices",
      configLabel: "Config",
      runningStatus: "Running",
      stoppedStatus: "Stopped",
      startProxy: "Start Proxy",
      stopProxy: "Stop Proxy",
      lanProxyPrefix: "LAN Proxy",
      proxyAddrPrefix: "Proxy",
      startupError: "Startup Error",
      backendConnectError: "Failed to connect to backend",
      backendTimeoutError: "Backend IPC not ready",

      tabMonitor: "Monitor",
      tabWhitelist: "IP Whitelist",
      tabSettings: "Settings",

      kpiUpLabel: "Upload Speed",
      kpiTotalUpLabel: "Total Up",
      kpiDownLabel: "Download Speed",
      kpiTotalDownLabel: "Total Down",
      kpiConnsLabel: "Active / Total Requests",
      kpiConnsSub: "Concurrent tunnels & requests",
      kpiWlLabel: "IP Whitelist",
      kpiBlockedLabel: "Blocked",
      wlEnabledCount: "Enabled ({n})",
      wlDisabledState: "Disabled",

      chartTitle: "Real-time Traffic (Last 3 Minutes)",
      legendUp: "Upload",
      legendDown: "Download",
      tooltipUp: "Up",
      tooltipDown: "Down",

      logsTitle: "Access Logs",
      logSearchPlaceholder: "Filter by client IP or URL...",
      filterAll: "All Statuses",
      filterAllowed: "Allowed / Active",
      filterBlocked: "Blocked",
      filterAuthFailed: "Auth Failed",
      filterError: "Error",
      pauseRefresh: "Pause",
      resumeRefresh: "Resume",
      exportCsv: "Export CSV",
      clearLogs: "Clear",

      colTime: "Time",
      colClientIp: "Client IP",
      colMethod: "Method",
      colTarget: "Target / URL",
      colStatus: "Status",
      colTraffic: "Up / Down",
      colDuration: "Duration",
      logsEmpty: "No matching access logs",

      quickAllowBtn: "+Allow",
      quickAllowTitle: "Add IP to whitelist",
      quickAllowNote: "Added from access logs",
      durationActive: "Active",
      badgeActive: "Active",
      badgeAllowed: "{code} OK",
      badgeBlocked: "403 Blocked",
      badgeAuth: "407 Auth",
      badgeError: "Error",

      wlMasterTitle: "IP Whitelist Control",
      wlMasterDesc: "When enabled, only clients matching enabled IP or CIDR rules below can use this proxy. Others receive 403 Forbidden and are logged. Changes take effect immediately.",
      wlAddTitle: "Add Allowed IP / CIDR",
      wlCidrPlaceholder: "IP or CIDR, e.g. 192.168.1.100 or 192.168.1.0/24",
      wlNotePlaceholder: "Note (optional), e.g. Office PC",
      wlAddBtn: "Add Rule",
      wlListTitle: "Whitelist Rules",
      wlColCidr: "IP / CIDR",
      wlColNote: "Note",
      wlColEnabled: "Enabled",
      wlColAction: "Action",
      wlEmpty: "No whitelist rules",
      wlDelete: "Delete",
      quickFillLabel: "Quick Fill:",
      chipLoopback: "Loopback",
      chipLocalIp: "Local NIC IP",
      chipLocalSubnet: "Local LAN Subnet",

      setBasicTitle: "Basic Proxy Settings",
      setBindLabel: "Bind Address",
      bindPresetTitle: "Quick fill network interface address",
      presetSelect: "Quick select...",
      presetAllInterfaces: "0.0.0.0 (All Interfaces)",
      presetLocalOnly: "127.0.0.1 (Localhost Only)",
      presetInterfaceSuffix: "Local NIC",
      setPortLabel: "Listen Port",
      setTimeoutLabel: "Connect Timeout (sec)",
      setMaxLogsLabel: "Max In-Memory Logs",
      setAutoStartLabel: "Automatically start proxy when Simple Proxy opens",
      setAuthTitle: "Proxy Authentication (Basic Auth)",
      setAuthUserLabel: "Username",
      setAuthPassLabel: "Password",
      setSocksTitle: "Upstream SOCKS5 Forwarding (Optional)",
      setSocksDesc: "Leave empty to connect directly. When configured, HTTPS (CONNECT) tunnels forward through this SOCKS5 proxy.",
      setSocksAddrLabel: "SOCKS5 Address",
      setSocksAddrPlaceholder: "e.g. 127.0.0.1:1080",
      setSocksUserLabel: "SOCKS5 Username",
      setSocksPassLabel: "SOCKS5 Password",
      optionalPlaceholder: "Optional",
      saveSettingsBtn: "Save & Apply Settings",

      exitTitle: "Confirm Exit",
      exitMsgRunning: "Proxy server is currently running. Exiting will close all active connections. Are you sure you want to exit?",
      exitMsgStopped: "Are you sure you want to exit Simple Proxy?",
      cancel: "Cancel",
      confirmExit: "Exit"
    }
  };

  const state = {
    lang: (navigator.language || "").toLowerCase().startsWith("zh") ? "zh" : "en",
    configPath: "",
    config: null,
    localIPs: [],
    running: false,
    listenAddr: "",
    lastError: "",
    monitor: {
      total_up_bytes: 0,
      total_down_bytes: 0,
      current_up_bps: 0,
      current_down_bps: 0,
      active_conns: 0,
      total_requests: 0,
      blocked_requests: 0,
      history: [],
      logs: []
    },
    logsPaused: false,
    hoverIndex: -1
  };

  function t(key, vars) {
    const dict = I18N[state.lang] || I18N.en;
    let text = dict[key] ?? I18N.en[key] ?? key;
    if (vars) {
      for (const [k, v] of Object.entries(vars)) {
        text = text.replace(`{${k}}`, String(v));
      }
    }
    return text;
  }

  const els = {
    configPathHint: document.getElementById("configPathHint"),
    statusPill: document.getElementById("statusPill"),
    statusText: document.getElementById("statusText"),
    lanHint: document.getElementById("lanHint"),
    langSwitch: document.getElementById("langSwitch"),
    toggleProxyBtn: document.getElementById("toggleProxyBtn"),
    errorBanner: document.getElementById("errorBanner"),

    kpiUpBps: document.getElementById("kpiUpBps"),
    kpiDownBps: document.getElementById("kpiDownBps"),
    kpiTotalUp: document.getElementById("kpiTotalUp"),
    kpiTotalDown: document.getElementById("kpiTotalDown"),
    kpiActiveConns: document.getElementById("kpiActiveConns"),
    kpiTotalReqs: document.getElementById("kpiTotalReqs"),
    kpiWhitelistState: document.getElementById("kpiWhitelistState"),
    kpiBlockedReqs: document.getElementById("kpiBlockedReqs"),

    trafficChart: document.getElementById("trafficChart"),
    chartTooltip: document.getElementById("chartTooltip"),

    logSearch: document.getElementById("logSearch"),
    logStatusFilter: document.getElementById("logStatusFilter"),
    pauseLogsBtn: document.getElementById("pauseLogsBtn"),
    exportLogsBtn: document.getElementById("exportLogsBtn"),
    clearLogsBtn: document.getElementById("clearLogsBtn"),
    logsTableBody: document.getElementById("logsTableBody"),

    whitelistMasterToggle: document.getElementById("whitelistMasterToggle"),
    whitelistAddForm: document.getElementById("whitelistAddForm"),
    wlCidrInput: document.getElementById("wlCidrInput"),
    wlNoteInput: document.getElementById("wlNoteInput"),
    wlQuickChips: document.getElementById("wlQuickChips"),
    wlCount: document.getElementById("wlCount"),
    whitelistTableBody: document.getElementById("whitelistTableBody"),

    settingsForm: document.getElementById("settingsForm"),
    cfgBind: document.getElementById("cfgBind"),
    bindPresetSelect: document.getElementById("bindPresetSelect"),
    cfgPort: document.getElementById("cfgPort"),
    cfgTimeout: document.getElementById("cfgTimeout"),
    cfgMaxLogs: document.getElementById("cfgMaxLogs"),
    cfgAutoStart: document.getElementById("cfgAutoStart"),
    cfgAuthEnabled: document.getElementById("cfgAuthEnabled"),
    cfgAuthUser: document.getElementById("cfgAuthUser"),
    cfgAuthPass: document.getElementById("cfgAuthPass"),
    cfgSocksAddr: document.getElementById("cfgSocksAddr"),
    cfgSocksUser: document.getElementById("cfgSocksUser"),
    cfgSocksPass: document.getElementById("cfgSocksPass"),

    exitModal: document.getElementById("exitModal"),
    exitModalMessage: document.getElementById("exitModalMessage"),
    exitCancelBtn: document.getElementById("exitCancelBtn"),
    exitConfirmBtn: document.getElementById("exitConfirmBtn"),
    toast: document.getElementById("toast")
  };

  let toastTimer = null;

  function api() {
    return window.go && window.go.main && window.go.main.App;
  }

  function waitForBackend(timeoutMs = 5000) {
    return new Promise((resolve, reject) => {
      const start = Date.now();
      const check = () => {
        if (api()) {
          resolve(api());
          return;
        }
        if (Date.now() - start > timeoutMs) {
          reject(new Error(t("backendTimeoutError")));
          return;
        }
        setTimeout(check, 40);
      };
      check();
    });
  }

  function applyTranslations() {
    document.documentElement.lang = state.lang === "zh" ? "zh-CN" : "en";

    document.querySelectorAll("[data-i18n]").forEach((el) => {
      const key = el.getAttribute("data-i18n");
      el.textContent = t(key);
    });
    document.querySelectorAll("[data-i18n-placeholder]").forEach((el) => {
      const key = el.getAttribute("data-i18n-placeholder");
      el.setAttribute("placeholder", t(key));
    });
    document.querySelectorAll("[data-i18n-title]").forEach((el) => {
      const key = el.getAttribute("data-i18n-title");
      el.setAttribute("title", t(key));
    });

    els.langSwitch.querySelectorAll(".lang-btn").forEach((btn) => {
      btn.classList.toggle("active", btn.getAttribute("data-lang") === state.lang);
    });

    if (state.configPath) {
      els.configPathHint.textContent = `${t("configLabel")}: ${state.configPath}`;
    }
    els.pauseLogsBtn.textContent = state.logsPaused ? t("resumeRefresh") : t("pauseRefresh");

    renderBindPresets();
    renderStatus();
    renderKPIs();
    renderWhitelist();
    renderLogs();
  }

  function showToast(msg, ok = true) {
    if (!msg) return;
    els.toast.textContent = msg;
    els.toast.className = `toast ${ok ? "ok" : "err"}`;
    clearTimeout(toastTimer);
    toastTimer = setTimeout(() => {
      els.toast.classList.add("hidden");
    }, 2800);
  }

  function formatBytes(bytes) {
    const n = Number(bytes) || 0;
    if (n < 1024) return `${n} B`;
    if (n < 1024 * 1024) return `${(n / 1024).toFixed(1)} KB`;
    if (n < 1024 * 1024 * 1024) return `${(n / (1024 * 1024)).toFixed(2)} MB`;
    return `${(n / (1024 * 1024 * 1024)).toFixed(2)} GB`;
  }

  function formatSpeed(bps) {
    return `${formatBytes(bps)}/s`;
  }

  function escapeHtml(str) {
    return String(str ?? "")
      .replace(/&/g, "&amp;")
      .replace(/</g, "&lt;")
      .replace(/>/g, "&gt;")
      .replace(/"/g, "&quot;")
      .replace(/'/g, "&#39;");
  }

  // Language switcher
  els.langSwitch.addEventListener("click", async (e) => {
    const btn = e.target.closest("[data-lang]");
    if (!btn) return;
    const nextLang = btn.getAttribute("data-lang");
    if (nextLang === state.lang) return;
    state.lang = nextLang;
    applyTranslations();
    if (api()) {
      const res = await api().SetLanguage(nextLang);
      if (res && res.config) {
        state.config = res.config;
      }
      if (res && res.message) {
        showToast(res.message, res.ok);
      }
    }
  });

  // Tabs
  document.querySelectorAll(".tab-btn").forEach((btn) => {
    btn.addEventListener("click", () => {
      const target = btn.dataset.tab;
      document.querySelectorAll(".tab-btn").forEach((b) => b.classList.toggle("active", b === btn));
      document.querySelectorAll(".tab-panel").forEach((panel) => {
        panel.classList.toggle("active", panel.id === `tab-${target}`);
      });
      if (target === "monitor") {
        requestAnimationFrame(drawChart);
      }
    });
  });

  // Topbar status & actions
  function renderStatus() {
    els.statusPill.classList.toggle("running", state.running);
    if (state.running) {
      els.statusText.textContent = `${t("runningStatus")} · ${state.listenAddr}`;
      els.toggleProxyBtn.textContent = t("stopProxy");
      els.toggleProxyBtn.className = "btn btn-danger";
    } else {
      els.statusText.textContent = `${t("stoppedStatus")} · ${state.listenAddr}`;
      els.toggleProxyBtn.textContent = t("startProxy");
      els.toggleProxyBtn.className = "btn btn-primary";
    }

    if (state.lastError) {
      els.errorBanner.textContent = state.lastError;
      els.errorBanner.classList.remove("hidden");
    } else {
      els.errorBanner.classList.add("hidden");
    }

    const port = state.config ? state.config.port : "8888";
    if (state.localIPs.length > 0 && state.config && state.config.bind === "0.0.0.0") {
      els.lanHint.textContent = `${t("lanProxyPrefix")}: ${state.localIPs[0]}:${port}`;
    } else if (state.listenAddr) {
      els.lanHint.textContent = `${t("proxyAddrPrefix")}: ${state.listenAddr}`;
    } else {
      els.lanHint.textContent = "";
    }
  }

  function renderKPIs() {
    const m = state.monitor;
    els.kpiUpBps.textContent = formatSpeed(m.current_up_bps);
    els.kpiDownBps.textContent = formatSpeed(m.current_down_bps);
    els.kpiTotalUp.textContent = formatBytes(m.total_up_bytes);
    els.kpiTotalDown.textContent = formatBytes(m.total_down_bytes);
    els.kpiActiveConns.textContent = String(m.active_conns || 0);
    els.kpiTotalReqs.textContent = String(m.total_requests || 0);
    els.kpiBlockedReqs.textContent = String(m.blocked_requests || 0);

    if (state.config && state.config.whitelist_enabled) {
      const enabledCount = (state.config.whitelist || []).filter((r) => r.enabled).length;
      els.kpiWhitelistState.textContent = t("wlEnabledCount", { n: enabledCount });
      els.kpiWhitelistState.className = "kpi-value text-emerald";
    } else {
      els.kpiWhitelistState.textContent = t("wlDisabledState");
      els.kpiWhitelistState.className = "kpi-value";
    }
  }

  // Canvas Traffic Chart
  const chartLayout = { padLeft: 66, padRight: 16, padTop: 14, padBottom: 24 };

  function drawChart() {
    const canvas = els.trafficChart;
    if (!canvas) return;
    const rect = canvas.parentElement.getBoundingClientRect();
    if (rect.width <= 0 || rect.height <= 0) return;

    const dpr = window.devicePixelRatio || 1;
    const width = Math.floor(rect.width);
    const height = Math.floor(rect.height);
    if (canvas.width !== width * dpr || canvas.height !== height * dpr) {
      canvas.width = width * dpr;
      canvas.height = height * dpr;
    }

    const ctx = canvas.getContext("2d");
    ctx.setTransform(dpr, 0, 0, dpr, 0, 0);
    ctx.clearRect(0, 0, width, height);

    const history = state.monitor.history || [];
    const plotW = width - chartLayout.padLeft - chartLayout.padRight;
    const plotH = height - chartLayout.padTop - chartLayout.padBottom;
    if (plotW <= 10 || plotH <= 10) return;

    let maxBps = 1024;
    for (const pt of history) {
      if (pt.up_bps > maxBps) maxBps = pt.up_bps;
      if (pt.down_bps > maxBps) maxBps = pt.down_bps;
    }
    maxBps = niceCeiling(maxBps * 1.15);

    const ticks = 4;
    ctx.font = '11px "Segoe UI", sans-serif';
    ctx.textAlign = "right";
    ctx.textBaseline = "middle";

    for (let i = 0; i <= ticks; i++) {
      const ratio = i / ticks;
      const y = chartLayout.padTop + plotH * (1 - ratio);
      ctx.strokeStyle = i === 0 ? "#475569" : "rgba(51, 65, 85, 0.55)";
      ctx.lineWidth = 1;
      ctx.beginPath();
      ctx.moveTo(chartLayout.padLeft, y);
      ctx.lineTo(chartLayout.padLeft + plotW, y);
      ctx.stroke();

      ctx.fillStyle = "#94a3b8";
      ctx.fillText(formatSpeed(Math.round(maxBps * ratio)), chartLayout.padLeft - 8, y);
    }

    if (history.length < 2) return;

    const xAt = (idx) => chartLayout.padLeft + (idx / (history.length - 1)) * plotW;
    const yAt = (val) => chartLayout.padTop + plotH - (Math.min(val, maxBps) / maxBps) * plotH;

    ctx.textAlign = "center";
    ctx.textBaseline = "top";
    ctx.fillStyle = "#64748b";
    const xSteps = 5;
    for (let i = 0; i <= xSteps; i++) {
      const idx = Math.round((i / xSteps) * (history.length - 1));
      const x = xAt(idx);
      ctx.fillText(history[idx].label || "", x, chartLayout.padTop + plotH + 6);
    }

    drawSeries(ctx, history, (pt) => pt.down_bps, xAt, yAt, chartLayout.padTop + plotH, "#10b981", "rgba(16, 185, 129, 0.22)", "rgba(16, 185, 129, 0.0)");
    drawSeries(ctx, history, (pt) => pt.up_bps, xAt, yAt, chartLayout.padTop + plotH, "#38bdf8", "rgba(56, 189, 248, 0.20)", "rgba(56, 189, 248, 0.0)");

    if (state.hoverIndex >= 0 && state.hoverIndex < history.length) {
      const hx = xAt(state.hoverIndex);
      const pt = history[state.hoverIndex];
      ctx.strokeStyle = "rgba(226, 232, 240, 0.45)";
      ctx.lineWidth = 1;
      ctx.setLineDash([4, 4]);
      ctx.beginPath();
      ctx.moveTo(hx, chartLayout.padTop);
      ctx.lineTo(hx, chartLayout.padTop + plotH);
      ctx.stroke();
      ctx.setLineDash([]);

      drawDot(ctx, hx, yAt(pt.down_bps), "#10b981");
      drawDot(ctx, hx, yAt(pt.up_bps), "#38bdf8");
    }
  }

  function drawSeries(ctx, history, getVal, xAt, yAt, baselineY, strokeColor, fillTop, fillBottom) {
    const grad = ctx.createLinearGradient(0, chartLayout.padTop, 0, baselineY);
    grad.addColorStop(0, fillTop);
    grad.addColorStop(1, fillBottom);

    ctx.beginPath();
    for (let i = 0; i < history.length; i++) {
      const x = xAt(i);
      const y = yAt(getVal(history[i]));
      if (i === 0) ctx.moveTo(x, y);
      else ctx.lineTo(x, y);
    }
    ctx.lineTo(xAt(history.length - 1), baselineY);
    ctx.lineTo(xAt(0), baselineY);
    ctx.closePath();
    ctx.fillStyle = grad;
    ctx.fill();

    ctx.beginPath();
    for (let i = 0; i < history.length; i++) {
      const x = xAt(i);
      const y = yAt(getVal(history[i]));
      if (i === 0) ctx.moveTo(x, y);
      else ctx.lineTo(x, y);
    }
    ctx.strokeStyle = strokeColor;
    ctx.lineWidth = 2;
    ctx.stroke();
  }

  function drawDot(ctx, x, y, color) {
    ctx.beginPath();
    ctx.arc(x, y, 4, 0, Math.PI * 2);
    ctx.fillStyle = color;
    ctx.fill();
    ctx.strokeStyle = "#0f172a";
    ctx.lineWidth = 1.5;
    ctx.stroke();
  }

  function niceCeiling(val) {
    if (val <= 1024) return 1024;
    const units = [1024, 1024 * 1024, 1024 * 1024 * 1024];
    let unit = 1;
    for (const u of units) {
      if (val >= u) unit = u;
    }
    const scaled = val / unit;
    const steps = [1, 2, 5, 10, 20, 50, 100, 200, 500, 1024];
    for (const s of steps) {
      if (scaled <= s) return s * unit;
    }
    return Math.ceil(scaled) * unit;
  }

  els.trafficChart.addEventListener("mousemove", (e) => {
    const history = state.monitor.history || [];
    if (history.length < 2) return;
    const rect = els.trafficChart.getBoundingClientRect();
    const x = e.clientX - rect.left;
    const plotW = rect.width - chartLayout.padLeft - chartLayout.padRight;
    if (x < chartLayout.padLeft || x > chartLayout.padLeft + plotW) {
      state.hoverIndex = -1;
      els.chartTooltip.classList.add("hidden");
      drawChart();
      return;
    }
    const ratio = (x - chartLayout.padLeft) / plotW;
    const idx = Math.max(0, Math.min(history.length - 1, Math.round(ratio * (history.length - 1))));
    state.hoverIndex = idx;
    const pt = history[idx];

    els.chartTooltip.innerHTML = `
      <div><strong>${escapeHtml(pt.label)}</strong></div>
      <div style="color:#38bdf8">${escapeHtml(t("tooltipUp"))}: ${formatSpeed(pt.up_bps)}</div>
      <div style="color:#10b981">${escapeHtml(t("tooltipDown"))}: ${formatSpeed(pt.down_bps)}</div>
    `;
    els.chartTooltip.classList.remove("hidden");

    const tipWidth = 140;
    const left = Math.min(rect.width - tipWidth, Math.max(8, x + 12));
    els.chartTooltip.style.left = `${left}px`;
    els.chartTooltip.style.top = "12px";
    drawChart();
  });

  els.trafficChart.addEventListener("mouseleave", () => {
    state.hoverIndex = -1;
    els.chartTooltip.classList.add("hidden");
    drawChart();
  });

  window.addEventListener("resize", () => drawChart());

  // Access Logs
  function renderLogs() {
    const logs = state.monitor.logs || [];
    const query = els.logSearch.value.trim().toLowerCase();
    const statusFilter = els.logStatusFilter.value;

    const filtered = logs.filter((item) => {
      if (statusFilter !== "all") {
        if (statusFilter === "allowed") {
          if (item.status !== "allowed" && item.status !== "active") return false;
        } else if (item.status !== statusFilter) {
          return false;
        }
      }
      if (query) {
        const hay = `${item.client_ip} ${item.target} ${item.method}`.toLowerCase();
        if (!hay.includes(query)) return false;
      }
      return true;
    });

    if (filtered.length === 0) {
      els.logsTableBody.innerHTML = `<tr><td colspan="7" class="empty-cell">${escapeHtml(t("logsEmpty"))}</td></tr>`;
      return;
    }

    const whitelistedSet = new Set(
      (state.config && state.config.whitelist ? state.config.whitelist : []).map((r) => r.cidr)
    );

    const rowsHtml = filtered.slice(0, 300).map((item) => {
      const shortTime = item.time && item.time.length >= 19 ? item.time.slice(11, 19) : item.time;
      const statusBadge = renderStatusBadge(item);
      const canQuickAdd = item.client_ip && !whitelistedSet.has(item.client_ip) && !whitelistedSet.has(`${item.client_ip}/32`);
      const quickBtn = canQuickAdd
        ? `<button class="ip-quick-allow" data-allow-ip="${escapeHtml(item.client_ip)}" title="${escapeHtml(t("quickAllowTitle"))}">${escapeHtml(t("quickAllowBtn"))}</button>`
        : "";
      const trafficText = `${formatBytes(item.up_bytes)} / ${formatBytes(item.down_bytes)}`;
      const durText = item.status === "active" ? t("durationActive") : `${item.duration_ms} ms`;

      return `
        <tr>
          <td title="${escapeHtml(item.time)}">${escapeHtml(shortTime)}</td>
          <td>
            <div class="ip-cell">
              <span class="selectable">${escapeHtml(item.client_ip)}</span>
              ${quickBtn}
            </div>
          </td>
          <td><span class="badge badge-method">${escapeHtml(item.method)}</span></td>
          <td class="selectable" title="${escapeHtml(item.target)}">${escapeHtml(item.target)}</td>
          <td title="${escapeHtml(item.error || "")}">${statusBadge}</td>
          <td>${escapeHtml(trafficText)}</td>
          <td>${escapeHtml(durText)}</td>
        </tr>
      `;
    }).join("");

    els.logsTableBody.innerHTML = rowsHtml;
  }

  function renderStatusBadge(item) {
    switch (item.status) {
      case "active":
        return `<span class="badge badge-active">${escapeHtml(t("badgeActive"))}</span>`;
      case "allowed":
        return `<span class="badge badge-allowed">${escapeHtml(t("badgeAllowed", { code: item.status_code || 200 }))}</span>`;
      case "blocked":
        return `<span class="badge badge-blocked">${escapeHtml(t("badgeBlocked"))}</span>`;
      case "auth_failed":
        return `<span class="badge badge-auth">${escapeHtml(t("badgeAuth"))}</span>`;
      default:
        return `<span class="badge badge-error">${escapeHtml(item.status_code ? String(item.status_code) : t("badgeError"))}</span>`;
    }
  }

  els.logsTableBody.addEventListener("click", async (e) => {
    const btn = e.target.closest("[data-allow-ip]");
    if (!btn) return;
    const ip = btn.getAttribute("data-allow-ip");
    const res = await api().AddWhitelistRule(ip, t("quickAllowNote"));
    handleActionResult(res);
  });

  els.logSearch.addEventListener("input", renderLogs);
  els.logStatusFilter.addEventListener("change", renderLogs);

  els.pauseLogsBtn.addEventListener("click", () => {
    state.logsPaused = !state.logsPaused;
    els.pauseLogsBtn.textContent = state.logsPaused ? t("resumeRefresh") : t("pauseRefresh");
    if (!state.logsPaused) {
      renderLogs();
    }
  });

  els.clearLogsBtn.addEventListener("click", async () => {
    const res = await api().ClearLogs();
    if (res.ok) {
      state.monitor.logs = [];
      renderLogs();
    }
    showToast(res.message, res.ok);
  });

  els.exportLogsBtn.addEventListener("click", async () => {
    const res = await api().ExportLogsCSV();
    showToast(res.message, res.ok);
  });

  // Whitelist UI
  function renderWhitelist() {
    if (!state.config) return;
    els.whitelistMasterToggle.checked = Boolean(state.config.whitelist_enabled);
    const rules = state.config.whitelist || [];
    els.wlCount.textContent = String(rules.length);

    if (rules.length === 0) {
      els.whitelistTableBody.innerHTML = `<tr><td colspan="4" class="empty-cell">${escapeHtml(t("wlEmpty"))}</td></tr>`;
    } else {
      els.whitelistTableBody.innerHTML = rules.map((r) => `
        <tr>
          <td class="selectable"><code>${escapeHtml(r.cidr)}</code></td>
          <td>${escapeHtml(r.note || "-")}</td>
          <td>
            <label class="switch">
              <input type="checkbox" data-wl-toggle="${escapeHtml(r.cidr)}" ${r.enabled ? "checked" : ""} />
              <span class="slider"></span>
            </label>
          </td>
          <td>
            <button class="btn btn-danger-ghost btn-sm" data-wl-remove="${escapeHtml(r.cidr)}" type="button">${escapeHtml(t("wlDelete"))}</button>
          </td>
        </tr>
      `).join("");
    }

    const chips = [
      { cidr: "127.0.0.1", note: t("chipLoopback") }
    ];
    for (const ip of state.localIPs) {
      chips.push({ cidr: ip, note: t("chipLocalIp") });
      const parts = ip.split(".");
      if (parts.length === 4) {
        chips.push({ cidr: `${parts[0]}.${parts[1]}.${parts[2]}.0/24`, note: t("chipLocalSubnet") });
      }
    }
    els.wlQuickChips.innerHTML = `<span class="chip-label">${escapeHtml(t("quickFillLabel"))}</span>` + chips.map((c) => `
      <button class="chip-btn" type="button" data-chip-cidr="${escapeHtml(c.cidr)}" data-chip-note="${escapeHtml(c.note)}">
        + ${escapeHtml(c.cidr)} (${escapeHtml(c.note)})
      </button>
    `).join("");
  }

  els.whitelistMasterToggle.addEventListener("change", async () => {
    const res = await api().SetWhitelistEnabled(els.whitelistMasterToggle.checked);
    handleActionResult(res);
  });

  els.whitelistAddForm.addEventListener("submit", async (e) => {
    e.preventDefault();
    const cidr = els.wlCidrInput.value.trim();
    const note = els.wlNoteInput.value.trim();
    if (!cidr) return;
    const res = await api().AddWhitelistRule(cidr, note);
    handleActionResult(res);
    if (res.ok) {
      els.wlCidrInput.value = "";
      els.wlNoteInput.value = "";
    }
  });

  els.wlQuickChips.addEventListener("click", (e) => {
    const btn = e.target.closest("[data-chip-cidr]");
    if (!btn) return;
    els.wlCidrInput.value = btn.getAttribute("data-chip-cidr");
    els.wlNoteInput.value = btn.getAttribute("data-chip-note");
    els.wlCidrInput.focus();
  });

  els.whitelistTableBody.addEventListener("change", async (e) => {
    const input = e.target.closest("[data-wl-toggle]");
    if (!input) return;
    const cidr = input.getAttribute("data-wl-toggle");
    const res = await api().ToggleWhitelistRule(cidr, input.checked);
    handleActionResult(res);
  });

  els.whitelistTableBody.addEventListener("click", async (e) => {
    const btn = e.target.closest("[data-wl-remove]");
    if (!btn) return;
    const cidr = btn.getAttribute("data-wl-remove");
    const res = await api().RemoveWhitelistRule(cidr);
    handleActionResult(res);
  });

  // Settings UI
  function renderBindPresets() {
    const baseOpts = [
      `<option value="">${escapeHtml(t("presetSelect"))}</option>`,
      `<option value="0.0.0.0">${escapeHtml(t("presetAllInterfaces"))}</option>`,
      `<option value="127.0.0.1">${escapeHtml(t("presetLocalOnly"))}</option>`
    ];
    for (const ip of state.localIPs) {
      baseOpts.push(`<option value="${escapeHtml(ip)}">${escapeHtml(ip)} (${escapeHtml(t("presetInterfaceSuffix"))})</option>`);
    }
    els.bindPresetSelect.innerHTML = baseOpts.join("");
  }

  function populateSettings() {
    if (!state.config) return;
    const c = state.config;
    els.cfgBind.value = c.bind || "0.0.0.0";
    els.cfgPort.value = c.port || "8888";
    els.cfgTimeout.value = String(c.timeout_seconds || 10);
    els.cfgMaxLogs.value = String(c.max_log_entries || 2000);
    els.cfgAutoStart.checked = Boolean(c.auto_start_proxy);
    els.cfgAuthEnabled.checked = Boolean(c.basic_auth_enabled);
    els.cfgAuthUser.value = c.basic_auth_username || "";
    els.cfgAuthPass.value = c.basic_auth_password || "";
    els.cfgSocksAddr.value = c.socks5_address || "";
    els.cfgSocksUser.value = c.socks5_username || "";
    els.cfgSocksPass.value = c.socks5_password || "";
    renderBindPresets();
  }

  els.bindPresetSelect.addEventListener("change", () => {
    if (els.bindPresetSelect.value) {
      els.cfgBind.value = els.bindPresetSelect.value;
      els.bindPresetSelect.value = "";
    }
  });

  els.settingsForm.addEventListener("submit", async (e) => {
    e.preventDefault();
    if (!state.config) return;
    const updated = {
      ...state.config,
      language: state.lang,
      bind: els.cfgBind.value.trim(),
      port: String(els.cfgPort.value).trim(),
      timeout_seconds: parseInt(els.cfgTimeout.value, 10) || 10,
      max_log_entries: parseInt(els.cfgMaxLogs.value, 10) || 2000,
      auto_start_proxy: els.cfgAutoStart.checked,
      basic_auth_enabled: els.cfgAuthEnabled.checked,
      basic_auth_username: els.cfgAuthUser.value.trim(),
      basic_auth_password: els.cfgAuthPass.value,
      socks5_address: els.cfgSocksAddr.value.trim(),
      socks5_username: els.cfgSocksUser.value.trim(),
      socks5_password: els.cfgSocksPass.value
    };
    const res = await api().SaveConfig(updated);
    handleActionResult(res);
  });

  // Start / Stop Proxy
  els.toggleProxyBtn.addEventListener("click", async () => {
    els.toggleProxyBtn.disabled = true;
    try {
      const res = state.running ? await api().StopProxy() : await api().StartProxy();
      handleActionResult(res);
    } finally {
      els.toggleProxyBtn.disabled = false;
    }
  });

  function handleActionResult(res) {
    if (!res) return;
    if (res.config) {
      state.config = res.config;
    }
    state.running = Boolean(res.running);
    if (res.addr) {
      state.listenAddr = res.addr;
    }
    if (!res.ok && res.message) {
      state.lastError = res.message;
    } else if (res.ok) {
      state.lastError = "";
    }
    renderStatus();
    renderKPIs();
    renderWhitelist();
    if (!state.logsPaused) {
      renderLogs();
    }
    showToast(res.message, res.ok);
  }

  async function refreshLive() {
    try {
      const snap = await api().GetLiveSnapshot();
      state.running = Boolean(snap.running);
      state.listenAddr = snap.listen_addr || state.listenAddr;
      state.lastError = snap.last_error || "";
      state.monitor = snap.monitor || state.monitor;

      renderStatus();
      renderKPIs();
      drawChart();
      if (!state.logsPaused) {
        renderLogs();
      }
    } catch (err) {
      console.error("refreshLive failed:", err);
    }
  }

  async function init() {
    applyTranslations();
    try {
      await waitForBackend();
      const initial = await api().GetState();
      state.running = Boolean(initial.running);
      state.listenAddr = initial.listen_addr || "";
      state.lastError = initial.last_error || "";
      state.localIPs = initial.local_ips || [];
      state.configPath = initial.config_path || "";
      state.config = initial.config;
      state.monitor = initial.monitor || state.monitor;

      if (state.config && state.config.language) {
        state.lang = state.config.language;
      } else {
        const detected = (navigator.language || "").toLowerCase().startsWith("zh") ? "zh" : "en";
        state.lang = detected;
        const res = await api().SetLanguage(detected);
        if (res && res.config) {
          state.config = res.config;
        }
      }

      populateSettings();
      applyTranslations();
      drawChart();

      if (window.runtime && typeof window.runtime.EventsOn === "function") {
        window.runtime.EventsOn("confirm-exit", (proxyRunning) => {
          showExitModal(Boolean(proxyRunning));
        });
      }

      setInterval(refreshLive, 1000);
    } catch (err) {
      els.statusText.textContent = t("startupError");
      els.errorBanner.textContent = err.message || t("backendConnectError");
      els.errorBanner.classList.remove("hidden");
    }
  }

  function showExitModal(proxyRunning) {
    els.exitModalMessage.textContent = proxyRunning ? t("exitMsgRunning") : t("exitMsgStopped");
    els.exitModal.classList.remove("hidden");
    els.exitCancelBtn.focus();
  }

  function hideExitModal() {
    els.exitModal.classList.add("hidden");
  }

  els.exitCancelBtn.addEventListener("click", hideExitModal);
  els.exitConfirmBtn.addEventListener("click", async () => {
    els.exitConfirmBtn.disabled = true;
    await api().ConfirmQuit();
  });
  els.exitModal.addEventListener("click", (e) => {
    if (e.target === els.exitModal) hideExitModal();
  });
  document.addEventListener("keydown", (e) => {
    if (e.key === "Escape" && !els.exitModal.classList.contains("hidden")) {
      hideExitModal();
    }
  });

  init();
})();
