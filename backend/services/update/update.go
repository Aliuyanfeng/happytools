// Package update 提供检查 GitHub Releases 更新的功能
package update

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// UpdateService 更新检查服务
type UpdateService struct {
	CurrentVersion string // 当前版本，如 "1.0.0"
	Owner          string // GitHub 用户名
	Repo           string // 仓库名
}

// ReleaseInfo GitHub Release 信息
type ReleaseInfo struct {
	TagName string `json:"tag_name"` // "v1.1.0"
	Name    string `json:"name"`
	Body    string `json:"body"`     // 更新说明
	HTMLURL string `json:"html_url"` // release 页面链接
}

// UpdateResult 检查结果
type UpdateResult struct {
	HasUpdate bool   `json:"hasUpdate"`
	Latest    string `json:"latest"`
	Current   string `json:"current"`
	ReleaseURL string `json:"releaseUrl"`
	ReleaseNote string `json:"releaseNote"`
}

func NewUpdateService(version, owner, repo string) *UpdateService {
	return &UpdateService{
		CurrentVersion: version,
		Owner:          owner,
		Repo:           repo,
	}
}

// CheckUpdate 检查是否有新版本，无更新返回 hasUpdate=false
func (s *UpdateService) CheckUpdate() (*UpdateResult, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", s.Owner, s.Repo)

	client := &http.Client{Timeout: 8 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "happytools-updater")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("网络请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("GitHub API 返回 %d", resp.StatusCode)
	}

	var release ReleaseInfo
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	latest := strings.TrimPrefix(release.TagName, "v")
	current := strings.TrimPrefix(s.CurrentVersion, "v")

	result := &UpdateResult{
		HasUpdate:   isNewer(latest, current),
		Latest:      latest,
		Current:     current,
		ReleaseURL:  release.HTMLURL,
		ReleaseNote: release.Body,
	}
	return result, nil
}

// GetCurrentVersion 返回当前版本号
func (s *UpdateService) GetCurrentVersion() string {
	return s.CurrentVersion
}

// isNewer 比较两个语义化版本号，latest > current 返回 true
func isNewer(latest, current string) bool {
	if latest == current || latest == "" {
		return false
	}
	// dev 版本始终认为有更新，方便本地调试
	if current == "dev" {
		return latest != ""
	}
	lp := parseSemver(latest)
	cp := parseSemver(current)
	for i := 0; i < 3; i++ {
		if lp[i] > cp[i] {
			return true
		}
		if lp[i] < cp[i] {
			return false
		}
	}
	return false
}

func parseSemver(v string) [3]int {
	parts := strings.Split(v, ".")
	var result [3]int
	for i := 0; i < 3 && i < len(parts); i++ {
		n, _ := strconv.Atoi(parts[i])
		result[i] = n
	}
	return result
}
