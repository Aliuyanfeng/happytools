// Package nuclei
// @description: Nuclei POC 模板解析服务
package nuclei

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v3/pkg/application"
	"gopkg.in/yaml.v3"
)

// ─── 数据结构 ────────────────────────────────────────────────────

type TemplateInfo struct {
	Name        string            `yaml:"name" json:"name"`
	Author      interface{}       `yaml:"author" json:"author"`
	Severity    string            `yaml:"severity" json:"severity"`
	Description string            `yaml:"description" json:"description"`
	Reference   interface{}       `yaml:"reference" json:"reference"`
	Tags        interface{}       `yaml:"tags" json:"tags"`
	Metadata    map[string]string `yaml:"metadata" json:"metadata"`
	Remediation string            `yaml:"remediation" json:"remediation"`
	CVSS        string            `yaml:"cvss-metrics" json:"cvss_metrics"`
	CVSSScore   string            `yaml:"cvss-score" json:"cvss_score"`
	CVEID       string            `yaml:"cve-id" json:"cve_id"`
	CWEID       string            `yaml:"cwe-id" json:"cwe_id"`
}

type Matcher struct {
	Type      string      `yaml:"type" json:"type"`
	Part      string      `yaml:"part" json:"part"`
	Words     []string    `yaml:"words" json:"words"`
	Regex     []string    `yaml:"regex" json:"regex"`
	Status    []int       `yaml:"status" json:"status"`
	DSL       []string    `yaml:"dsl" json:"dsl"`
	Condition string      `yaml:"condition" json:"condition"`
	Negative  bool        `yaml:"negative" json:"negative"`
	Name      string      `yaml:"name" json:"name"`
}

type Extractor struct {
	Type  string   `yaml:"type" json:"type"`
	Name  string   `yaml:"name" json:"name"`
	Part  string   `yaml:"part" json:"part"`
	Regex []string `yaml:"regex" json:"regex"`
	Group int      `yaml:"group" json:"group"`
	JSON  []string `yaml:"json" json:"json"`
	XPath []string `yaml:"xpath" json:"xpath"`
	Words []string `yaml:"words" json:"words"`
}

type HTTPRequest struct {
	Method            interface{}  `yaml:"method" json:"method"`
	Path              []string     `yaml:"path" json:"path"`
	Raw               []string     `yaml:"raw" json:"raw"`
	Headers           interface{}  `yaml:"headers" json:"headers"`
	Body              string       `yaml:"body" json:"body"`
	Redirects         bool         `yaml:"redirects" json:"redirects"`
	MaxRedirects      int          `yaml:"max-redirects" json:"max_redirects"`
	StopAtFirstMatch  bool         `yaml:"stop-at-first-match" json:"stop_at_first_match"`
	Matchers          []Matcher    `yaml:"matchers" json:"matchers"`
	MatchersCondition string       `yaml:"matchers-condition" json:"matchers_condition"`
	Extractors        []Extractor  `yaml:"extractors" json:"extractors"`
	Payloads          interface{}  `yaml:"payloads" json:"payloads"`
	AttackType        string       `yaml:"attack" json:"attack"`
	FuzzingRules      interface{}  `yaml:"fuzzing" json:"fuzzing"`
}

type DNSRequest struct {
	Name       string      `yaml:"name" json:"name"`
	Type       string      `yaml:"type" json:"type"`
	Class      string      `yaml:"class" json:"class"`
	Retries    int         `yaml:"retries" json:"retries"`
	Matchers   []Matcher   `yaml:"matchers" json:"matchers"`
	Extractors []Extractor `yaml:"extractors" json:"extractors"`
}

type TCPRequest struct {
	Host       []string    `yaml:"host" json:"host"`
	Port       string      `yaml:"port" json:"port"`
	Data       string      `yaml:"data" json:"data"`
	Matchers   []Matcher   `yaml:"matchers" json:"matchers"`
	Extractors []Extractor `yaml:"extractors" json:"extractors"`
}

type CodeRequest struct {
	Engine     []string    `yaml:"engine" json:"engine"`
	Source     string      `yaml:"source" json:"source"`
	Args       []string    `yaml:"args" json:"args"`
	Pattern    string      `yaml:"pattern" json:"pattern"`
	Matchers   []Matcher   `yaml:"matchers" json:"matchers"`
	MatchersCondition string `yaml:"matchers-condition" json:"matchers_condition"`
	Extractors []Extractor `yaml:"extractors" json:"extractors"`
}

type Variable struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type NucleiTemplate struct {
	ID        string                 `yaml:"id" json:"id"`
	Info      TemplateInfo           `yaml:"info" json:"info"`
	HTTP      []HTTPRequest          `yaml:"http" json:"http"`
	DNS       []DNSRequest           `yaml:"dns" json:"dns"`
	TCP       []TCPRequest           `yaml:"tcp" json:"tcp"`
	Code      []CodeRequest          `yaml:"code" json:"code"`
	Variables map[string]interface{} `yaml:"variables" json:"variables"`
	SelfContained bool               `yaml:"self-contained" json:"self_contained"`
	// 解析后的辅助字段
	RawYAML   string   `json:"raw_yaml"`
	FilePath  string   `json:"file_path"`
	Protocol  string   `json:"protocol"`
	AuthorStr string   `json:"author_str"`
	TagList   []string `json:"tag_list"`
	RefList   []string `json:"ref_list"`
}

// ─── 服务 ────────────────────────────────────────────────────────

type NucleiService struct {
	app *application.App
}

func NewNucleiService(app *application.App) *NucleiService {
	return &NucleiService{app: app}
}

// OpenFile 打开文件选择对话框
func (s *NucleiService) OpenFile() (string, error) {
	path, err := s.app.Dialog.OpenFile().
		SetTitle("选择 Nuclei 模板文件").
		AddFilter("YAML 文件", "*.yaml").
		PromptForSingleSelection()
	if err != nil {
		return "", err
	}
	return path, nil
}

// ParseFile 解析指定路径的 Nuclei 模板文件
func (s *NucleiService) ParseFile(path string) (*NucleiTemplate, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %w", err)
	}
	return parseTemplate(data, path)
}

// ParseContent 直接解析 YAML 文本内容
func (s *NucleiService) ParseContent(content string) (*NucleiTemplate, error) {
	return parseTemplate([]byte(content), "")
}

// ─── 核心解析 ────────────────────────────────────────────────────

func parseTemplate(data []byte, filePath string) (*NucleiTemplate, error) {
	var tpl NucleiTemplate
	if err := yaml.Unmarshal(data, &tpl); err != nil {
		return nil, fmt.Errorf("YAML 解析失败: %w", err)
	}

	tpl.RawYAML = string(data)
	tpl.FilePath = filePath

	// 推断协议
	switch {
	case len(tpl.HTTP) > 0:
		tpl.Protocol = "http"
	case len(tpl.DNS) > 0:
		tpl.Protocol = "dns"
	case len(tpl.TCP) > 0:
		tpl.Protocol = "tcp"
	case len(tpl.Code) > 0:
		tpl.Protocol = "code"
	default:
		tpl.Protocol = "unknown"
	}

	// 规范化 author
	tpl.AuthorStr = toStringSliceJoined(tpl.Info.Author)

	// 规范化 tags
	tpl.TagList = toStringSlice(tpl.Info.Tags)

	// 规范化 reference
	tpl.RefList = toStringSlice(tpl.Info.Reference)

	// 文件名作为备用 ID
	if tpl.ID == "" && filePath != "" {
		tpl.ID = strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
	}

	return &tpl, nil
}

// toStringSlice 将 interface{} 转为 []string（兼容单值和数组）
func toStringSlice(v interface{}) []string {
	if v == nil {
		return nil
	}
	switch val := v.(type) {
	case string:
		if val == "" {
			return nil
		}
		// 逗号分隔的 tags
		parts := strings.Split(val, ",")
		result := make([]string, 0, len(parts))
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p != "" {
				result = append(result, p)
			}
		}
		return result
	case []interface{}:
		result := make([]string, 0, len(val))
		for _, item := range val {
			if s, ok := item.(string); ok && s != "" {
				result = append(result, s)
			}
		}
		return result
	}
	return nil
}

func toStringSliceJoined(v interface{}) string {
	parts := toStringSlice(v)
	if len(parts) == 0 {
		if s, ok := v.(string); ok {
			return s
		}
		return ""
	}
	return strings.Join(parts, ", ")
}
