# Simple Proxy GUI

English | [简体中文](./README_zh.md)

A lightweight Windows desktop HTTP/HTTPS proxy tool with a clean GUI, built on [Wails v2](https://wails.io/) and forked from [jthomperoo/simple-proxy](https://github.com/jthomperoo/simple-proxy).

Distributed as a single self-contained `.exe` built directly by GitHub Actions from source, with SHA256 checksums included in every release.

## Features

- **HTTP & HTTPS (`CONNECT`) Proxy**: Ready-to-use local or LAN proxy server (`0.0.0.0:8888` by default) with one-click Start / Stop.
- **IP Whitelist (Hot-Reload)**:
  - Master toggle to enable or disable IP filtering without restarting the proxy.
  - Supports individual IPv4/IPv6 addresses and CIDR subnets (e.g., `192.168.1.0/24`), per-rule enable/disable switches, and one-click whitelisting from access logs.
  - Unlisted clients receive `403 Forbidden` and are recorded in the log table.
- **Real-time Traffic Chart**: Live 1-second sampling of upload/download speeds (`B/s`, `KB/s`, `MB/s`) and cumulative traffic over a 3-minute rolling window.
- **IP Access Logs**:
  - Tracks timestamp, client IP, HTTP method, target host/URL, status (`Active`, `200 OK`, `403 Blocked`, `407 Auth`, `Error`), upload/download bytes, and duration.
  - Filter by IP/URL or status, pause auto-refresh, and export logs to UTF-8 CSV.
- **Bilingual UI (`EN` / `中文`)**: Auto-detects system language on first launch and supports instant switching from the top bar.
- **Persistent Configuration**: Saves settings and whitelist rules to `config.json` alongside the executable.
- **Optional Auth & Upstream SOCKS5**: Supports HTTP Basic Authentication and forwarding `CONNECT` tunnels through an upstream SOCKS5 proxy.
- **Safe Close Prompt**: Prompts for confirmation when closing via the window `X` button, while automatically skipping the prompt during Windows shutdown/restart.

## Download & Verify

Download `simple-proxy-gui.exe` or `simple-proxy-gui_windows_amd64.zip` along with `checksums.txt` from the [Releases](https://github.com/shakespark/simple-proxy-gui/releases) page.

All official release binaries are compiled cleanly in public [GitHub Actions](https://github.com/shakespark/simple-proxy-gui/actions) runners directly from this repository's source code. You can verify the downloaded file hash in PowerShell:

```powershell
Get-FileHash .\simple-proxy-gui.exe -Algorithm SHA256
```

> **Note:** Requires Windows 10/11 with Microsoft Edge WebView2 Runtime (pre-installed on Windows 11 and most Windows 10 systems).

## Quick Test

After launching `simple-proxy-gui.exe`, test the proxy with `curl`:

```powershell
curl.exe --proxy "http://127.0.0.1:8888" "https://www.random.org/integers/?num=1&min=1&max=5&col=1&base=10&format=plain&rnd=new"
```

## Build from Source

Requires Go 1.26+:

```powershell
go test -v ./...
go build -tags desktop,production -ldflags="-H windowsgui -s -w" -o dist/simple-proxy-gui.exe .
```

## License

MIT License — see [LICENSE](./LICENSE). Upstream project copyright 2022 Jamie Thompson.
