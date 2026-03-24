/*
 * @Author: LiuYanFeng
 * @Date: 2025-07-03 17:16:49
 * @LastEditors: LiuYanFeng
 * @LastEditTime: 2026-03-06 16:51:00
 * @FilePath: \happytools\main.go
 * @Description: 像珍惜礼物一样珍惜今天
 *
 * Copyright (c) 2025 by ${git_name_email}, All Rights Reserved.
 */
package main

import (
	"embed"
	_ "embed"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Aliuyanfeng/happytools/backend/services/appsettings"
	"github.com/Aliuyanfeng/happytools/backend/services/category"
	"github.com/Aliuyanfeng/happytools/backend/services/clipboard"
	"github.com/Aliuyanfeng/happytools/backend/services/dailyreport"
	"github.com/Aliuyanfeng/happytools/backend/services/encryption"
	gitconfig "github.com/Aliuyanfeng/happytools/backend/services/gitconfig"
	"github.com/Aliuyanfeng/happytools/backend/services/makefile"
	"github.com/Aliuyanfeng/happytools/backend/services/pnginjector"
	"github.com/Aliuyanfeng/happytools/backend/services/greetservice"
	"github.com/Aliuyanfeng/happytools/backend/services/monitor"
	"github.com/Aliuyanfeng/happytools/backend/services/network"
	"github.com/Aliuyanfeng/happytools/backend/services/rename"
	"github.com/Aliuyanfeng/happytools/backend/services/todo"
	"github.com/Aliuyanfeng/happytools/backend/services/unitconverter"
	virusTotal "github.com/Aliuyanfeng/happytools/backend/services/vt"
	"github.com/Aliuyanfeng/happytools/backend/store"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
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

	// 初始化 bbolt 数据库
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("获取用户目录失败:", err)
	}
	dbPath := filepath.Join(homeDir, ".happytools", "data.db")
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		log.Fatal("创建数据目录失败:", err)
	}
	if err := store.Init(dbPath); err != nil {
		log.Fatal("初始化数据库失败:", err)
	}
	defer store.Close()

	// Create a new Wails application by providing the necessary options.
	// Variables 'Name' and 'Description' are for application metadata.
	// 'Assets' configures the asset server with the 'FS' variable pointing to the frontend files.
	// 'Bind' is a list of Go struct instances. The frontend has access to the methods of these instances.
	// 'Mac' options tailor the application when running an macOS.
	app := application.New(application.Options{
		Name:        "happytools",
		Description: "A tool that can enhance one's sense of happiness",
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})
	// Register the backend services with the application.
	// These services are accessible from the frontend and can be called using the 'backend' object.
	app.RegisterService(application.NewService(&greetservice.GreetService{}))
	app.RegisterService(application.NewService(monitor.NewSysInfoService()))
	app.RegisterService(application.NewService(todo.NewTodoService()))
	app.RegisterService(application.NewService(category.NewCategoryService()))
	app.RegisterService(application.NewService(dailyreport.NewDailyReportService()))
	app.RegisterService(application.NewService(appsettings.NewAppSettingsService()))
	app.RegisterService(application.NewService(rename.NewRenameService(app)))
	app.RegisterService(application.NewService(virusTotal.NewVTService(app)))
	app.RegisterService(application.NewService(network.NewFileTransferService(app)))
	app.RegisterService(application.NewService(network.NewTCPUDPService(app)))
	app.RegisterService(application.NewService(network.NewDNSService()))
	app.RegisterService(application.NewService(unitconverter.NewUnitConverterService()))
	app.RegisterService(application.NewService(encryption.NewEncryptionService()))
	app.RegisterService(application.NewService(clipboard.NewClipboardService(app)))
	app.RegisterService(application.NewService(pnginjector.NewPNGInjectorService(app)))
	app.RegisterService(application.NewService(gitconfig.NewGitConfigService(app)))
	app.RegisterService(application.NewService(makefile.NewMakefileService(app)))

	// Create a new window with the necessary options.
	// 'Title' is the title of the window.
	// 'Mac' options tailor the window when running on macOS.
	// 'BackgroundColour' is the background colour of the window.
	// 'URL' is the URL that will be loaded into the webview.
	mainWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
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
			DisableIcon: false,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/",
		Frameless:        true,
		EnableFileDrop: true,
	})

	mainWindow.SetAlwaysOnTop(false)

	// 监听文件拖拽事件，将文件路径推送给前端
	mainWindow.OnWindowEvent(events.Common.WindowFilesDropped, func(event *application.WindowEvent) {
		files := event.Context().DroppedFiles()
		log.Printf("[FileDrop] 收到拖拽文件: %+v", files)
		targetID := ""
		if details := event.Context().DropTargetDetails(); details != nil {
			targetID = details.ElementID
			log.Printf("[FileDrop] 目标元素 ID: %s", targetID)
		}
		app.Event.Emit("wails:file-drop", map[string]any{
			"files":  files,
			"target": targetID,
		})
	})

	systray := app.SystemTray.New()
	systray.SetLabel("HappyTools")
	systray.SetTooltip("HappyTools工具")

	systray.OnClick(func() {
		if mainWindow.IsVisible() {
			mainWindow.Hide()
		} else {
			mainWindow.Show()
		}
	})

	// Optional: Set debounce time for window show/hide
	systray.WindowDebounce(200 * time.Millisecond)

	// Create tray menu with only Quit option
	menu := application.NewMenu()
	menu.Add("退出 HappyTools").OnClick(func(ctx *application.Context) {
		app.Quit()
	})

	systray.SetMenu(menu)

	// 初始化应用设置服务
	appSettingsService := appsettings.NewAppSettingsService()

	// 延迟发送上次使用时间，确保前端已经准备好监听事件
	go func() {
		time.Sleep(2 * time.Second) // 等待2秒，让前端完全加载

		// 应用启动时先获取上次使用时间（这是上一次启动的时间）
		lastUsedTime := appSettingsService.GetLastUsedTime()
		if lastUsedTime != "" {
			// 通过事件将上次使用时间发送给前端
			app.Event.Emit("app:lastUsedTime", lastUsedTime)
			log.Printf("发送上次使用时间事件: %s", lastUsedTime)
		} else {
			// 如果没有上次使用时间，发送空字符串
			app.Event.Emit("app:lastUsedTime", "")
			log.Printf("没有上次使用时间记录")
		}

		// 然后更新数据库为当前时间（为下一次启动做准备）
		if err := appSettingsService.UpdateLastUsedTime(); err != nil {
			log.Printf("更新上次使用时间失败: %v", err)
		} else {
			log.Printf("已更新上次使用时间为: %s", time.Now().Format(time.RFC3339))
		}
	}()

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
	err = app.Run()

	// If an error occurred while running the application, log it and exit.
	if err != nil {
		log.Fatal(err)
	}
}

func getFileType(path string) string {
	if path == "" {
		return "unknown"
	}
	ext := strings.ToLower(filepath.Ext(path))
	if ext == "" {
		if fi, err := os.Stat(path); err == nil && fi.IsDir() {
			return "directory"
		}
		return "no extension"
	}
	return ext
}
