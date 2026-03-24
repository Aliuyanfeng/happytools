// Package gitconfig
// @description: .git/config INI 解析器 & Pretty_Printer
package gitconfig

import (
	"fmt"
	"regexp"
	"strings"
)

// ConfigSection 解析后的配置节
type ConfigSection struct {
	Name    string        `json:"name"`
	SubKey  string        `json:"subKey"`
	Entries []ConfigEntry `json:"entries"`
}

// ConfigEntry 配置节下的键值对
type ConfigEntry struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

var sectionWithSubKey = regexp.MustCompile(`^\[(\w[\w-]*)\s+"([^"]+)"\]\s*$`)
var sectionPlain = regexp.MustCompile(`^\[(\w[\w-]*)\]\s*$`)

// Parse 将 .git/config 文件内容解析为结构化数据
func Parse(content string) ([]ConfigSection, error) {
	var sections []ConfigSection
	var current *ConfigSection

	lines := strings.Split(content, "\n")
	for i, raw := range lines {
		line := strings.TrimRight(raw, "\r")
		trimmed := strings.TrimSpace(line)

		// 空行或注释
		if trimmed == "" || strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, ";") {
			continue
		}

		// 带子键的节头：[remote "origin"]
		if m := sectionWithSubKey.FindStringSubmatch(trimmed); m != nil {
			sections = append(sections, ConfigSection{Name: m[1], SubKey: m[2], Entries: []ConfigEntry{}})
			current = &sections[len(sections)-1]
			continue
		}

		// 普通节头：[core]
		if m := sectionPlain.FindStringSubmatch(trimmed); m != nil {
			sections = append(sections, ConfigSection{Name: m[1], SubKey: "", Entries: []ConfigEntry{}})
			current = &sections[len(sections)-1]
			continue
		}

		// 键值对
		if idx := strings.IndexByte(trimmed, '='); idx > 0 {
			if current == nil {
				return nil, fmt.Errorf("line %d: key-value pair outside of section", i+1)
			}
			key := strings.TrimSpace(trimmed[:idx])
			value := strings.TrimSpace(trimmed[idx+1:])
			current.Entries = append(current.Entries, ConfigEntry{Key: key, Value: value})
			continue
		}

		return nil, fmt.Errorf("line %d: unexpected content: %q", i+1, trimmed)
	}

	return sections, nil
}

// Serialize 将结构化配置数据序列化为 .git/config 格式文本
func Serialize(sections []ConfigSection) string {
	var sb strings.Builder
	for i, sec := range sections {
		if i > 0 {
			sb.WriteByte('\n')
		}
		if sec.SubKey != "" {
			fmt.Fprintf(&sb, "[%s \"%s\"]\n", sec.Name, sec.SubKey)
		} else {
			fmt.Fprintf(&sb, "[%s]\n", sec.Name)
		}
		for _, e := range sec.Entries {
			fmt.Fprintf(&sb, "\t%s = %s\n", e.Key, e.Value)
		}
	}
	return sb.String()
}
