# Simple Proxy GUI

[English](./README.md) | 简体中文

一个轻量的 Windows 桌面端 HTTP/HTTPS 代理工具，基于 [Wails v2](https://wails.io/) 构建，Fork 自开源项目 [jthomperoo/simple-proxy](https://github.com/jthomperoo/simple-proxy)。

所有 Release 二进制文件均由 GitHub Actions 从公开源码直接构建，并附带 SHA256 校验值。

## 主要功能

- **HTTP / HTTPS (`CONNECT`) 代理**：开箱即用的本机或局域网代理（默认 `0.0.0.0:8888`），支持界面一键启停。
- **IP 白名单（热生效）**：
  - 支持总开关与单条规则开关，修改后无需重启代理即刻生效。
  - 支持单个 IPv4/IPv6 及 CIDR 网段（如 `192.168.1.0/24`），并支持从访问记录中一键加白名单。
  - 不在白名单中的请求返回 `403 Forbidden` 并记录拦截日志。
- **实时流量折线图**：每秒采样实时上行/下行速率与累计流量，展示最近 3 分钟双线速率曲线与悬停详情。
- **IP 访问网址记录**：
  - 记录时间、客户端 IP、请求方法、目标域名/网址、状态（连接中 / 放行 / 拦截 / 认证失败 / 错误）、上下行流量与耗时。
  - 支持按 IP 或网址关键词搜索、按状态筛选、暂停刷新及导出 UTF-8 CSV。
- **中英文双语界面（`中文` / `EN`）**：首次启动自动匹配系统语言，顶部栏支持一键切换。
- **配置文件持久化**：自动在程序同目录读写 `config.json`，保存端口、认证、白名单和语言等设置。
- **可选代理认证与上游 SOCKS5**：支持 Basic Auth 账号密码认证及将 `CONNECT` 隧道转发到上游 SOCKS5 代理。
- **安全退出确认**：点击右上角 `X` 时弹出确认框防止误关；Windows 系统关机或重启时自动跳过弹窗，不阻塞系统关机。

## 下载与校验

前往 [Releases](https://github.com/shakespark/simple-proxy-gui/releases) 页面下载 `simple-proxy-gui.exe` 或 `simple-proxy-gui_windows_amd64.zip` 及对应的 `checksums.txt`。

所有官方安装包均在公开的 [GitHub Actions](https://github.com/shakespark/simple-proxy-gui/actions) 环境中由本仓库源码直接编译生成，可在 PowerShell 中核对哈希：

```powershell
Get-FileHash .\simple-proxy-gui.exe -Algorithm SHA256
```

> **运行环境要求**：Windows 10 / 11（需系统内置或已安装 Microsoft Edge WebView2 Runtime，Windows 11 默认已自带）。

## 从源码编译

需要 Go 1.26+：

```powershell
go test -v ./...
go build -tags desktop,production -ldflags="-H windowsgui -s -w" -o dist/simple-proxy-gui.exe .
```

## 开源协议

基于 MIT 协议开源，详见 [LICENSE](./LICENSE)。
