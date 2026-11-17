//go:build windows

package main

import "syscall"

const smShuttingDown = 0x2000

var (
	user32DLL        = syscall.NewLazyDLL("user32.dll")
	getSystemMetrics = user32DLL.NewProc("GetSystemMetrics")
)

func isSystemShuttingDown() bool {
	ret, _, _ := getSystemMetrics.Call(uintptr(smShuttingDown))
	return ret != 0
}
