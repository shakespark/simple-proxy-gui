package main

import (
	"embed"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/shakespark/simple-proxy-gui/config"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

var (
	Version = "development"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	var version bool
	flag.BoolVar(&version, "version", false, "prints current simple-proxy version")
	flag.Parse()

	if version {
		fmt.Println(Version)
		os.Exit(0)
	}

	cfgPath := config.ResolveConfigPath()
	app, err := NewApp(cfgPath)
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}

	err = wails.Run(&options.App{
		Title:     "Simple Proxy",
		Width:     1160,
		Height:    760,
		MinWidth:  920,
		MinHeight: 620,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 15, G: 23, B: 42, A: 1},
		OnStartup:        app.Startup,
		OnBeforeClose:    app.BeforeClose,
		OnShutdown:       app.Shutdown,
		Bind: []interface{}{
			app,
		},
		Windows: &windows.Options{
			Theme: windows.Dark,
			CustomTheme: &windows.ThemeSettings{
				DarkModeTitleBar:          windows.RGB(22, 32, 50),
				DarkModeTitleBarInactive:  windows.RGB(15, 23, 42),
				DarkModeTitleText:         windows.RGB(241, 245, 249),
				DarkModeTitleTextInactive: windows.RGB(148, 163, 184),
				DarkModeBorder:            windows.RGB(51, 65, 85),
				DarkModeBorderInactive:    windows.RGB(30, 41, 59),
			},
		},
	})
	if err != nil {
		log.Fatalf("wails error: %v", err)
	}
}
