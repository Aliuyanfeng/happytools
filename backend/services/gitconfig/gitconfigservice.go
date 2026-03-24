// Package gitconfig
// @description: Git Config Manager 主服务，暴露给 Wails 的所有方法
package gitconfig

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Aliuyanfeng/happytools/backend/store"
	"github.com/google/uuid"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// GitConfigService Git 配置管理服务
type GitConfigService struct {
	app *application.App
}

// NewGitConfigService 创建 GitConfigService 实例
func NewGitConfigService(app *application.App) *GitConfigService {
	return &GitConfigService{app: app}
}

// ========== 仓库管理 ==========

// AddRepository 添加仓库，验证路径有效性并写入 bbolt
func (s *GitConfigService) AddRepository(name, path, platform string) (*store.Repository, error) {
	configPath := filepath.Join(path, ".git", "config")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("该路径不是有效的 Git 仓库：%s", path)
	}

	// 检查重复路径
	repos, err := store.GetAllRepositories()
	if err != nil {
		return nil, err
	}
	for _, r := range repos {
		if r.Path == path {
			return nil, fmt.Errorf("该仓库已存在：%s", path)
		}
	}

	repo := &store.Repository{
		ID:        uuid.New().String(),
		Name:      name,
		Path:      path,
		Platform:  platform,
		CreatedAt: time.Now(),
	}
	if err := store.SaveRepository(repo); err != nil {
		return nil, err
	}
	return repo, nil
}

// ListRepositories 获取所有仓库记录
func (s *GitConfigService) ListRepositories() ([]*store.Repository, error) {
	return store.GetAllRepositories()
}

// DeleteRepository 删除仓库记录（不修改磁盘文件）
func (s *GitConfigService) DeleteRepository(id string) error {
	return store.DeleteRepository(id)
}

// OpenDirectoryDialog 打开系统目录选择对话框
func (s *GitConfigService) OpenDirectoryDialog() string {
	result, err := s.app.Dialog.OpenFile().
		CanChooseDirectories(true).
		CanChooseFiles(false).
		SetTitle("选择 Git 仓库根目录").
		PromptForSingleSelection()
	if err != nil {
		return ""
	}
	return result
}

// ========== 配置读写 ==========

// LoadConfig 读取并解析指定仓库的 .git/config
func (s *GitConfigService) LoadConfig(repoID string) ([]ConfigSection, error) {
	repo, err := store.GetRepository(repoID)
	if err != nil {
		return nil, err
	}
	return s.readConfig(repo.Path)
}

// SaveEntry 新增或修改指定节下的键值对，原子写文件
func (s *GitConfigService) SaveEntry(repoID, section, subKey, key, value string) error {
	if key == "" {
		return fmt.Errorf("键名不能为空")
	}
	repo, err := store.GetRepository(repoID)
	if err != nil {
		return err
	}

	sections, err := s.readConfig(repo.Path)
	if err != nil {
		return err
	}

	// 找到目标节
	secIdx := findSection(sections, section, subKey)
	if secIdx == -1 {
		// 节不存在则新建
		sections = append(sections, ConfigSection{
			Name:    section,
			SubKey:  subKey,
			Entries: []ConfigEntry{{Key: key, Value: value}},
		})
	} else {
		// 找到键则更新，否则追加
		entryIdx := findEntry(sections[secIdx].Entries, key)
		if entryIdx == -1 {
			sections[secIdx].Entries = append(sections[secIdx].Entries, ConfigEntry{Key: key, Value: value})
		} else {
			sections[secIdx].Entries[entryIdx].Value = value
		}
	}

	return s.writeConfig(repo.Path, sections)
}

// DeleteEntry 删除指定节下的键值对
func (s *GitConfigService) DeleteEntry(repoID, section, subKey, key string) error {
	repo, err := store.GetRepository(repoID)
	if err != nil {
		return err
	}

	sections, err := s.readConfig(repo.Path)
	if err != nil {
		return err
	}

	secIdx := findSection(sections, section, subKey)
	if secIdx == -1 {
		return fmt.Errorf("节 [%s] 不存在", section)
	}

	entryIdx := findEntry(sections[secIdx].Entries, key)
	if entryIdx == -1 {
		return fmt.Errorf("键 %s 不存在", key)
	}

	entries := sections[secIdx].Entries
	sections[secIdx].Entries = append(entries[:entryIdx], entries[entryIdx+1:]...)
	return s.writeConfig(repo.Path, sections)
}

// AddSection 新增配置节
func (s *GitConfigService) AddSection(repoID, section, subKey string) error {
	if section == "" {
		return fmt.Errorf("节名称不能为空")
	}
	repo, err := store.GetRepository(repoID)
	if err != nil {
		return err
	}

	sections, err := s.readConfig(repo.Path)
	if err != nil {
		return err
	}

	if findSection(sections, section, subKey) != -1 {
		return fmt.Errorf("节名称已存在")
	}

	sections = append(sections, ConfigSection{Name: section, SubKey: subKey, Entries: []ConfigEntry{}})
	return s.writeConfig(repo.Path, sections)
}

// DeleteSection 删除配置节及其所有键值对
func (s *GitConfigService) DeleteSection(repoID, section, subKey string) error {
	repo, err := store.GetRepository(repoID)
	if err != nil {
		return err
	}

	sections, err := s.readConfig(repo.Path)
	if err != nil {
		return err
	}

	secIdx := findSection(sections, section, subKey)
	if secIdx == -1 {
		return fmt.Errorf("节 [%s] 不存在", section)
	}

	sections = append(sections[:secIdx], sections[secIdx+1:]...)
	return s.writeConfig(repo.Path, sections)
}

// ========== KnownKey ==========

// GetKnownKeys 返回全部内置键定义
func (s *GitConfigService) GetKnownKeys() []KnownKey {
	return GetKnownKeys()
}

// GetKnownKeysForSection 返回指定节的内置键定义
func (s *GitConfigService) GetKnownKeysForSection(section string) []KnownKey {
	return GetKnownKeysForSection(section)
}

// ========== QuickPanel ==========

// GetQuickPanel 获取指定仓库的 QuickPanel 配置
func (s *GitConfigService) GetQuickPanel(repoID string) ([]store.QuickPanelItem, error) {
	items, err := store.GetQuickPanel(repoID)
	if err != nil {
		return nil, err
	}

	// 如果是默认配置且包含 remote "origin".url，尝试替换为实际存在的 remote
	if len(items) == 3 && items[2].Section == "remote" && items[2].SubKey == "origin" {
		sections, err := s.LoadConfig(repoID)
		if err == nil {
			// 查找第一个 remote 节
			for _, sec := range sections {
				if sec.Name == "remote" && sec.SubKey != "" {
					items[2].SubKey = sec.SubKey
					break
				}
			}
		}
	}

	return items, nil
}

// SaveQuickPanel 保存指定仓库的 QuickPanel 配置
func (s *GitConfigService) SaveQuickPanel(repoID string, items []store.QuickPanelItem) error {
	return store.SaveQuickPanel(repoID, items)
}

// GetRemoteNames 获取仓库中所有 remote 名称
func (s *GitConfigService) GetRemoteNames(repoID string) ([]string, error) {
	sections, err := s.LoadConfig(repoID)
	if err != nil {
		return nil, err
	}

	var remotes []string
	for _, sec := range sections {
		if sec.Name == "remote" && sec.SubKey != "" {
			remotes = append(remotes, sec.SubKey)
		}
	}
	return remotes, nil
}

// ========== 内部工具方法 ==========

func (s *GitConfigService) readConfig(repoPath string) ([]ConfigSection, error) {
	configPath := filepath.Join(repoPath, ".git", "config")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败：%v", err)
	}
	return Parse(string(data))
}

func (s *GitConfigService) writeConfig(repoPath string, sections []ConfigSection) error {
	configPath := filepath.Join(repoPath, ".git", "config")
	content := Serialize(sections)

	// 原子写：先写临时文件，再 Rename
	tmpPath := configPath + ".tmp"
	if err := os.WriteFile(tmpPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("写入临时文件失败：%v", err)
	}
	if err := os.Rename(tmpPath, configPath); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("替换配置文件失败：%v", err)
	}
	return nil
}

func findSection(sections []ConfigSection, name, subKey string) int {
	for i, s := range sections {
		if s.Name == name && s.SubKey == subKey {
			return i
		}
	}
	return -1
}

func findEntry(entries []ConfigEntry, key string) int {
	for i, e := range entries {
		if e.Key == key {
			return i
		}
	}
	return -1
}
