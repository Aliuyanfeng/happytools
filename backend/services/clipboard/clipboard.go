/*
 * @Author: LiuYanFeng
 * @Date: 2026-02-25 15:29:00
 * @LastEditors: LiuYanFeng
 * @LastEditTime: 2026-02-26 11:06:45
 * @FilePath: \happytools\backend\services\clipboard\clipboard.go
 * @Description: 剪贴板服务
 *
 * Copyright (c) 2026 by ${git_name_email}, All Rights Reserved.
 */
package clipboard

import (
	"github.com/wailsapp/wails/v3/pkg/application"
)

// ClipboardService 剪贴板服务
type ClipboardService struct {
	app *application.App
}

// NewClipboardService 创建剪贴板服务实例
func NewClipboardService(app *application.App) *ClipboardService {
	return &ClipboardService{app: app}
}

// SetText 设置剪贴板文本
func (c *ClipboardService) SetText(text string) bool {
	return c.app.Clipboard.SetText(text)
}

// GetText 获取剪贴板文本
func (c *ClipboardService) GetText() string {
	text, ok := c.app.Clipboard.Text()
	if !ok {
		return ""
	}
	return text
}
