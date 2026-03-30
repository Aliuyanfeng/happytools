// Package makefile
// @description: Makefile Visual Editor — 暴露给 Wails 的所有方法
package makefile

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Aliuyanfeng/happytools/backend/store"
	"github.com/google/uuid"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// MakefileService Makefile 可视化编辑器服务
type MakefileService struct {
	app *application.App
}

// NewMakefileService 创建 MakefileService 实例
func NewMakefileService(app *application.App) *MakefileService {
	return &MakefileService{app: app}
}

// ========== 文件管理 ==========

// OpenFile 读取并解析指定路径的 Makefile 文件，成功后记录到最近文件列表。
func (s *MakefileService) OpenFile(path string) (*MakefileDoc, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败：%v", err)
	}

	doc, err := Parse(string(data))
	if err != nil {
		return nil, fmt.Errorf("解析文件失败：%v", err)
	}

	// 记录到最近文件列表（忽略存储错误，不影响主流程）
	_ = store.SaveRecentFile(path)

	return doc, nil
}

// NewFile 在指定目录下新建一个空白 Makefile 文件，返回空文档。
func (s *MakefileService) NewFile(dir string) (*MakefileDoc, error) {
	path := filepath.Join(dir, "Makefile")

	// 写入空白文件
	if err := os.WriteFile(path, []byte(""), 0644); err != nil {
		return nil, fmt.Errorf("创建文件失败：%v", err)
	}

	return &MakefileDoc{
		Variables: []Variable{},
		Targets:   []Target{},
		RawBlocks: []RawBlock{},
	}, nil
}

// NewFromTemplate 根据模板 ID 初始化文档，并在指定目录下写入 Makefile 文件。
func (s *MakefileService) NewFromTemplate(dir string, templateID string) (*MakefileDoc, error) {
	templates, err := s.GetTemplates()
	if err != nil {
		return nil, err
	}

	var found *Template
	for i := range templates {
		if templates[i].ID == templateID {
			found = &templates[i]
			break
		}
	}
	if found == nil {
		return nil, fmt.Errorf("模板不存在：%s", templateID)
	}

	doc := found.Doc

	// 将模板内容写入磁盘
	path := filepath.Join(dir, "Makefile")
	content := Print(&doc)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return nil, fmt.Errorf("写入文件失败：%v", err)
	}

	return &doc, nil
}

// SaveFile 将文档序列化后原子写入指定路径。
// 策略：先写 <path>.tmp，成功后 os.Rename 替换；写临时文件失败则直接返回错误，原始文件不受影响。
func (s *MakefileService) SaveFile(path string, doc *MakefileDoc) error {
	content := Print(doc)
	tmpPath := path + ".tmp"

	if err := os.WriteFile(tmpPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("写入临时文件失败：%v", err)
	}

	if err := os.Rename(tmpPath, path); err != nil {
		// 尽力清理临时文件，但不掩盖原始错误
		_ = os.Remove(tmpPath)
		return fmt.Errorf("替换文件失败：%v", err)
	}

	return nil
}

// GetRecentFiles 返回最近打开的 Makefile 路径列表（最多 10 条）。
func (s *MakefileService) GetRecentFiles() ([]string, error) {
	return store.GetRecentFiles()
}

// SaveRawText 将原始文本内容原子写入指定路径（用于原始文本编辑模式保存）。
func (s *MakefileService) SaveRawText(path string, content string) error {
	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("写入临时文件失败：%v", err)
	}
	if err := os.Rename(tmpPath, path); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("替换文件失败：%v", err)
	}
	return nil
}

// ParseRawText 将原始文本内容解析为结构化 MakefileDoc，不写入磁盘。
func (s *MakefileService) ParseRawText(content string) (*MakefileDoc, error) {
	return Parse(content)
}

// OpenFileDialog 打开系统文件选择对话框，过滤 Makefile/makefile/*.mk 文件。
func (s *MakefileService) OpenFileDialog() string {
	result, err := s.app.Dialog.OpenFile().
		SetTitle("打开 Makefile 文件").
		AddFilter("Makefile 文件", "*.mk").
		PromptForSingleSelection()
	if err != nil {
		return ""
	}
	return result
}

// OpenDirectoryDialog 打开系统目录选择对话框，用于新建文件时选择目录。
func (s *MakefileService) OpenDirectoryDialog() string {
	result, err := s.app.Dialog.OpenFile().
		CanChooseDirectories(true).
		CanChooseFiles(false).
		SetTitle("选择目录").
		PromptForSingleSelection()
	if err != nil {
		return ""
	}
	return result
}

// ========== 模板管理 ==========

// GetTemplates 返回内置模板与自定义模板的合并列表。
func (s *MakefileService) GetTemplates() ([]Template, error) {
	builtin := GetBuiltinTemplates()

	blobs, err := store.GetMakefileTemplates()
	if err != nil {
		return nil, fmt.Errorf("获取自定义模板失败：%v", err)
	}

	result := make([]Template, 0, len(builtin)+len(blobs))
	result = append(result, builtin...)
	for _, b := range blobs {
		var t Template
		if err := json.Unmarshal(b, &t); err != nil {
			return nil, fmt.Errorf("解析自定义模板失败：%v", err)
		}
		result = append(result, t)
	}
	return result, nil
}

// SaveCustomTemplate 将当前文档另存为自定义模板，生成新 UUID 作为模板 ID。
func (s *MakefileService) SaveCustomTemplate(name string, description string, doc *MakefileDoc) error {
	t := Template{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		IsBuiltin:   false,
		Doc:         *doc,
	}
	data, err := json.Marshal(t)
	if err != nil {
		return fmt.Errorf("序列化模板失败：%v", err)
	}
	return store.SaveMakefileTemplate(t.ID, data)
}

// DeleteCustomTemplate 删除指定 ID 的自定义模板。
func (s *MakefileService) DeleteCustomTemplate(id string) error {
	return store.DeleteCustomTemplate(id)
}

// ========== 依赖检测 ==========

// ValidateDependencies 使用 DFS 检测目标依赖图中的所有环路，返回每个环路的目标名称列表。
// 若无环路则返回空切片。
func (s *MakefileService) ValidateDependencies(doc *MakefileDoc) ([][]string, error) {
	// 构建邻接表
	adj := make(map[string][]string, len(doc.Targets))
	for _, t := range doc.Targets {
		adj[t.Name] = t.Deps
	}

	var cycles [][]string

	// DFS 状态：0=未访问, 1=访问中（在当前路径上）, 2=已完成
	state := make(map[string]int, len(doc.Targets))
	path := make([]string, 0)
	inPath := make(map[string]bool)

	var dfs func(node string)
	dfs = func(node string) {
		if state[node] == 2 {
			return
		}
		if state[node] == 1 {
			// 找到环：从 path 中截取环路部分
			cycleStart := -1
			for i, n := range path {
				if n == node {
					cycleStart = i
					break
				}
			}
			if cycleStart != -1 {
				cycle := make([]string, len(path)-cycleStart)
				copy(cycle, path[cycleStart:])
				cycles = append(cycles, cycle)
			}
			return
		}

		state[node] = 1
		path = append(path, node)
		inPath[node] = true

		for _, dep := range adj[node] {
			dfs(dep)
		}

		path = path[:len(path)-1]
		inPath[node] = false
		state[node] = 2
	}

	for _, t := range doc.Targets {
		if state[t.Name] == 0 {
			dfs(t.Name)
		}
	}

	return cycles, nil
}
