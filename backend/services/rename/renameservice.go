// Package rename
// @author: liuyanfeng
// @date: 2026/2/12
// @description: 批量文件重命名服务
package rename

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// RenameRule 重命名规则
type RenameRule struct {
	Prefix              string `json:"prefix"`              // 前缀
	Suffix              string `json:"suffix"`              // 后缀
	StartNumber         int    `json:"startNumber"`         // 起始编号
	NumberDigits        int    `json:"numberDigits"`        // 编号位数
	NumberStep          int    `json:"numberStep"`          // 编号步长
	ReplaceFrom         string `json:"replaceFrom"`         // 要替换的文本
	ReplaceTo           string `json:"replaceTo"`           // 替换为的文本
	CaseSensitive       bool   `json:"caseSensitive"`       // 是否区分大小写
	KeepExtension       bool   `json:"keepExtension"`       // 是否保留扩展名
	PreviewBeforeRename bool   `json:"previewBeforeRename"` // 重命名前预览
}

// FileInfo 文件信息
type FileInfo struct {
	OriginalPath string `json:"originalPath"` // 原始路径
	OriginalName string `json:"originalName"` // 原始文件名
	NewName      string `json:"newName"`      // 新文件名
	NewPath      string `json:"newPath"`      // 新文件路径
	Size         int64  `json:"size"`         // 文件大小
	IsDir        bool   `json:"isDir"`        // 是否是目录
	Error        string `json:"error"`        // 错误信息
}

// RenameResult 重命名结果
type RenameResult struct {
	SuccessCount int        `json:"successCount"` // 成功数量
	FailedCount  int        `json:"failedCount"`  // 失败数量
	TotalCount   int        `json:"totalCount"`   // 总数量
	Results      []FileInfo `json:"results"`      // 重命名结果
}

// RenameService 重命名服务
type RenameService struct {
	app *application.App
}

// NewRenameService 创建重命名服务实例
func NewRenameService(app *application.App) *RenameService {
	return &RenameService{
		app: app,
	}
}

// OpenFileDialog 打开文件对话框并返回选中的文件路径
func (r *RenameService) OpenFileDialog() string {
	result, err := r.app.Dialog.OpenFile().
		SetTitle("选择要重命名的文件").
		PromptForMultipleSelection()
	if err != nil {
		return ""
	}
	// 返回第一个文件路径（保持与VTService一致）
	if len(result) > 0 {
		return result[0]
	}
	return ""
}

// OpenFileDialogs 打开文件对话框并返回选中的多个文件路径
func (r *RenameService) OpenFileDialogs() []string {
	result, err := r.app.Dialog.OpenFile().
		SetTitle("选择要重命名的文件").
		PromptForMultipleSelection()
	if err != nil {
		return []string{}
	}
	return result
}

// PreviewRename 预览重命名结果
func (r *RenameService) PreviewRename(files []FileInfo, rule RenameRule) ([]FileInfo, error) {
	if len(files) == 0 {
		return nil, errors.New("没有文件需要重命名")
	}

	results := make([]FileInfo, len(files))

	for i, file := range files {
		newName, err := r.generateNewFileName(file.OriginalName, i, rule)
		if err != nil {
			results[i] = FileInfo{
				OriginalPath: file.OriginalPath,
				OriginalName: file.OriginalName,
				NewName:      "",
				Size:         file.Size,
				IsDir:        file.IsDir,
				Error:        err.Error(),
			}
			continue
		}

		results[i] = FileInfo{
			OriginalPath: file.OriginalPath,
			OriginalName: file.OriginalName,
			NewName:      newName,
			Size:         file.Size,
			IsDir:        file.IsDir,
			Error:        "",
		}
	}

	return results, nil
}

// ExecuteRename 执行批量重命名
func (r *RenameService) ExecuteRename(files []FileInfo, rule RenameRule) (*RenameResult, error) {
	if len(files) == 0 {
		return nil, errors.New("没有文件需要重命名")
	}

	result := &RenameResult{
		SuccessCount: 0,
		FailedCount:  0,
		TotalCount:   len(files),
		Results:      make([]FileInfo, len(files)),
	}

	for i, file := range files {
		// 首先检查源文件是否存在
		if _, err := os.Stat(file.OriginalPath); err != nil {
			errorMsg := fmt.Sprintf("源文件不存在或无法访问: %s (错误: %v)", file.OriginalPath, err)
			result.Results[i] = FileInfo{
				OriginalPath: file.OriginalPath,
				OriginalName: file.OriginalName,
				NewName:      "",
				Size:         file.Size,
				IsDir:        file.IsDir,
				Error:        errorMsg,
			}
			result.FailedCount++
			continue
		}

		// 生成新文件名
		newName, err := r.generateNewFileName(file.OriginalName, i, rule)
		if err != nil {
			result.Results[i] = FileInfo{
				OriginalPath: file.OriginalPath,
				OriginalName: file.OriginalName,
				NewName:      "",
				Size:         file.Size,
				IsDir:        file.IsDir,
				Error:        err.Error(),
			}
			result.FailedCount++
			continue
		}

		// 构建完整的新文件路径
		dir := filepath.Dir(file.OriginalPath)
		newPath := filepath.Join(dir, newName)

		// 检查新文件是否已存在
		if _, err := os.Stat(newPath); err == nil {
			errorMsg := fmt.Sprintf("文件已存在: %s", newName)
			result.Results[i] = FileInfo{
				OriginalPath: file.OriginalPath,
				OriginalName: file.OriginalName,
				NewName:      newName,
				Size:         file.Size,
				IsDir:        file.IsDir,
				Error:        errorMsg,
			}
			result.FailedCount++
			continue
		}

		// 执行重命名
		err = os.Rename(file.OriginalPath, newPath)
		if err != nil {
			errorMsg := fmt.Sprintf("重命名失败: %v", err)
			result.Results[i] = FileInfo{
				OriginalPath: file.OriginalPath,
				OriginalName: file.OriginalName,
				NewName:      newName,
				Size:         file.Size,
				IsDir:        file.IsDir,
				Error:        errorMsg,
			}
			result.FailedCount++
			continue
		}

		// 重命名成功
		result.Results[i] = FileInfo{
			OriginalPath: file.OriginalPath,
			OriginalName: file.OriginalName,
			NewName:      newName,
			NewPath:      newPath,
			Size:         file.Size,
			IsDir:        file.IsDir,
			Error:        "",
		}
		result.SuccessCount++
	}

	return result, nil
}

// generateNewFileName 生成新文件名
func (r *RenameService) generateNewFileName(originalName string, index int, rule RenameRule) (string, error) {
	// 提取文件名和扩展名
	var nameWithoutExt string
	var extension string

	if rule.KeepExtension {
		ext := filepath.Ext(originalName)
		if ext != "" {
			nameWithoutExt = strings.TrimSuffix(originalName, ext)
			extension = ext
		} else {
			nameWithoutExt = originalName
			extension = ""
		}
	} else {
		nameWithoutExt = originalName
		extension = ""
	}

	// 应用文本替换
	if rule.ReplaceFrom != "" {
		if rule.CaseSensitive {
			nameWithoutExt = strings.ReplaceAll(nameWithoutExt, rule.ReplaceFrom, rule.ReplaceTo)
		} else {
			nameWithoutExt = strings.ReplaceAll(strings.ToLower(nameWithoutExt), strings.ToLower(rule.ReplaceFrom), rule.ReplaceTo)
		}
	}

	// 添加前缀
	if rule.Prefix != "" {
		nameWithoutExt = rule.Prefix + nameWithoutExt
	}

	// 添加编号
	number := rule.StartNumber + (index * rule.NumberStep)
	formattedNumber := fmt.Sprintf("%0*d", rule.NumberDigits, number)
	nameWithoutExt = nameWithoutExt + formattedNumber

	// 添加后缀
	if rule.Suffix != "" {
		nameWithoutExt = nameWithoutExt + rule.Suffix
	}

	// 组合新文件名
	newName := nameWithoutExt + extension

	// 验证新文件名
	if newName == "" {
		return "", errors.New("生成的文件名为空")
	}

	// 检查文件名是否包含非法字符
	if strings.ContainsAny(newName, `<>:"/\|?*`) {
		return "", errors.New("文件名包含非法字符")
	}

	return newName, nil
}

// ValidateRule 验证重命名规则
func (r *RenameService) ValidateRule(rule RenameRule) error {
	if rule.StartNumber < 0 {
		return errors.New("起始编号不能为负数")
	}

	if rule.NumberDigits < 1 || rule.NumberDigits > 10 {
		return errors.New("编号位数必须在1-10之间")
	}

	if rule.NumberStep < 1 {
		return errors.New("编号步长必须大于0")
	}

	return nil
}

// GetFileInfo 获取文件信息
func (r *RenameService) GetFileInfo(filePath string) (*FileInfo, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	return &FileInfo{
		OriginalPath: filePath,
		OriginalName: filepath.Base(filePath),
		NewName:      "",
		Size:         fileInfo.Size(),
		IsDir:        fileInfo.IsDir(),
		Error:        "",
	}, nil
}

// BatchGetFileInfo 批量获取文件信息
func (r *RenameService) BatchGetFileInfo(filePaths []string) ([]FileInfo, error) {
	results := make([]FileInfo, len(filePaths))

	for i, filePath := range filePaths {
		fileInfo, err := r.GetFileInfo(filePath)
		if err != nil {
			results[i] = FileInfo{
				OriginalPath: filePath,
				OriginalName: filepath.Base(filePath),
				NewName:      "",
				Size:         0,
				IsDir:        false,
				Error:        err.Error(),
			}
			continue
		}

		results[i] = *fileInfo
	}

	return results, nil
}
