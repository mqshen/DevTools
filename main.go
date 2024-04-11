package main

import (
	"context"
	"devtools/backend/services"
	"devtools/backend/storage"
	"embed"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
	"runtime"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

var version = "0.0.0"

func main() {
	pref := storage.NewPreferences()
	prefSvc := services.Preferences(pref)
	yamlConvertor := services.YamlConvertor()
	securitySvc := services.Securities()
	hashService := services.HashGenerator()
	encodeSvc := services.EncodeService()
	traceRouteSvc := services.TraceRouter()

	settings := storage.NewSettings()
	ipService := services.IPServices(settings, pref)
	//
	prefSvc.SetAppVersion(version)
	windowWidth, windowHeight, maximised := prefSvc.GetWindowSize()
	windowStartState := options.Normal
	if maximised {
		windowStartState = options.Maximised
	}

	// menu
	appMenu := menu.NewMenu()
	if runtime.GOOS == "darwin" {
		appMenu.Append(menu.AppMenu())
		appMenu.Append(menu.EditMenu())
		appMenu.Append(menu.WindowMenu())
	}

	// Create application with options
	err := wails.Run(&options.App{
		Title:            "DevTools",
		Width:            windowWidth,
		Height:           windowHeight,
		WindowStartState: windowStartState,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		Frameless:        runtime.GOOS != "darwin",
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: func(ctx context.Context) {
			yamlConvertor.Start(ctx)
			ipService.Start(ctx)
			traceRouteSvc.Start(ctx)
		},
		OnDomReady: func(ctx context.Context) {
			x, y := prefSvc.GetWindowPosition(ctx)
			wailsRuntime.WindowSetPosition(ctx, x, y)
			wailsRuntime.WindowShow(ctx)
		},
		OnBeforeClose: func(ctx context.Context) (prevent bool) {
			x, y := wailsRuntime.WindowGetPosition(ctx)
			prefSvc.SaveWindowPosition(x, y)
			return false
		},
		OnShutdown: func(ctx context.Context) {
		},
		Bind: []interface{}{
			prefSvc,
			securitySvc,
			yamlConvertor,
			hashService,
			ipService,
			encodeSvc,
			traceRouteSvc,
		},
		Mac: &mac.Options{
			TitleBar: mac.TitleBarHiddenInset(),
			About: &mac.AboutInfo{
				Title:   "DevTools" + version,
				Message: "A modern lightweight cross-platform Dev Tools.\n\nCopyright © 2024",
				Icon:    icon,
			},
			WebviewIsTransparent: false,
			WindowIsTranslucent:  true,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
