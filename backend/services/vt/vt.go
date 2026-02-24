// Package virusTotal
// @author: liuyanfeng
// @date: 2025/9/5 17:48
// @description: VT service for file operations
package virusTotal

import (
	"github.com/wailsapp/wails/v3/pkg/application"
)

// VTService provides file-related operations.
type VTService struct {
	app *application.App
}

// NewVTService creates a new instance of VTService.
func NewVTService(app *application.App) *VTService {
	return &VTService{
		app: app,
	}
}

// OpenFileDialog opens a file dialog and returns the selected file path.
func (v *VTService) OpenFileDialog() string {
	result, err := v.app.Dialog.OpenFile().
		CanChooseFiles(true).
		SetTitle("请选择一个样本文件").
		PromptForSingleSelection()
	if err != nil {
		return ""
	}
	return result
}

// OpenFileDialogs opens a file dialog and returns the selected file paths.
func (v *VTService) OpenFileDialogs() string {
	result, err := v.app.Dialog.OpenFile().
		CanChooseDirectories(true).
		CanChooseFiles(false).
		SetTitle("请选择一个样本目录").
		PromptForSingleSelection()
	if err != nil {
		return ""
	}
	return result
}
