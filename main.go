/*
 * @Author: LiuYanFeng
 * @Date: 2025-07-03 17:16:49
 * @LastEditors: LiuYanFeng
 * @LastEditTime: 2025-07-29 09:08:00
 * @FilePath: \happytools\main.go
 * @Description: 像珍惜礼物一样珍惜今天
 *
 * Copyright (c) 2025 by ${git_name_email}, All Rights Reserved.
 */
package main

import (
	"embed"
	_ "embed"
	"github.com/Aliuyanfeng/happytools/backend/services/greetservice"
	"github.com/Aliuyanfeng/happytools/backend/services/monitor"
	"log"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// Wails uses Go's `embed` package to embed the frontend files into the binary.
// Any files in the frontend/dist folder will be embedded into the binary and
// made available to the frontend.
// See https://pkg.go.dev/embed for more information.

//go:embed all:frontend/dist
var assets embed.FS

// main function serves as the application's entry point. It initializes the application, creates a window,
// and starts a goroutine that emits a time-based event every second. It subsequently runs the application and
// logs any error that might occur.
func main() {

	// Create a new Wails application by providing the necessary options.
	// Variables 'Name' and 'Description' are for application metadata.
	// 'Assets' configures the asset server with the 'FS' variable pointing to the frontend files.
	// 'Bind' is a list of Go struct instances. The frontend has access to the methods of these instances.
	// 'Mac' options tailor the application when running an macOS.
	app := application.New(application.Options{
		Name:        "happytools",
		Description: "A tool that can enhance one's sense of happiness",
		Services: []application.Service{
			application.NewService(&greetservice.GreetService{}),
			application.NewService(monitor.NewSysInfoService()),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	systray := app.SystemTray.New()
	systray.SetLabel("HappyTools")
	//systray.SetIcon()
	systray.SetTooltip("HappyTools工具")

	// Create a window
	window := app.Window.New()

	// Attach the window to the system tray
	systray.AttachWindow(window)

	// Optional: Set window offset from tray icon
	systray.WindowOffset(10)

	// Optional: Set debounce time for window show/hide
	systray.WindowDebounce(200 * time.Millisecond)

	menu := application.NewMenu()
	menu.Add("Open").OnClick(func(ctx *application.Context) {
		app.Show()
	})
	menu.Add("Quit").OnClick(func(ctx *application.Context) {
		app.Quit()
	})

	systray.SetMenu(menu)

	// Create a new window with the necessary options.
	// 'Title' is the title of the window.
	// 'Mac' options tailor the window when running on macOS.
	// 'BackgroundColour' is the background colour of the window.
	// 'URL' is the URL that will be loaded into the webview.
	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:     "happytools",
		Width:     1024,
		Height:    768,
		MinWidth:  1024,
		MinHeight: 768,
		MaxWidth:  1024,
		MaxHeight: 768,
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		Windows: application.WindowsWindow{
			//ExStyle: w32.WS_EX_TOOLWINDOW | w32.WS_EX_NOREDIRECTIONBITMAP | w32.WS_EX_TOPMOST,
			DisableIcon: true,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/",
	})

	// Create a goroutine that emits an event containing the current time every second.
	// The frontend can listen to this event and update the UI accordingly.
	go func() {
		for {
			now := time.Now().Format(time.RFC1123)
			app.Event.Emit("time", now)
			time.Sleep(time.Second)
		}
	}()

	// 可随时终止和启动
	go func() {
		// 默认执行一次
		// 获取系统信息
		sysInfo, err := monitor.NewSysInfoService().GetSysInfo()
		if err != nil {
			log.Printf("获取系统信息失败: %v", err)
		}
		// 通过 Events 推送到前端
		app.Event.Emit("monitor:sysInfo", sysInfo)
		// 5秒执行一次
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			// 获取系统信息
			sysInfo, err := monitor.NewSysInfoService().GetSysInfo()
			if err != nil {
				log.Printf("获取系统信息失败: %v", err)
				continue
			}
			// 通过 Events 推送到前端
			app.Event.Emit("monitor:sysInfo", sysInfo)
		}
	}()

	// Run the application. This blocks until the application has been exited.
	err := app.Run()

	// If an error occurred while running the application, log it and exit.
	if err != nil {
		log.Fatal(err)
	}
}
