// Package makefile
// @description: Makefile 解析器 — 逐行状态机
package makefile

import (
	"regexp"
	"strings"
)

var (
	reVariable = regexp.MustCompile(`^([A-Za-z_][A-Za-z0-9_]*)\s*(=|:=|\?=|\+=)\s*(.*)`)
	reTarget   = regexp.MustCompile(`^([^\s#:][^:]*):(.*)`)
	rePhony    = regexp.MustCompile(`^\.PHONY\s*:\s*(.*)`)
)

// Parse 将 Makefile 文本内容解析为结构化 MakefileDoc。
// 无法识别的语法不返回错误，而是保留为 RawBlock。
func Parse(content string) (*MakefileDoc, error) {
	doc := &MakefileDoc{}
	phonySet := map[string]bool{}

	lines := strings.Split(content, "\n")

	// 第一遍：收集所有 .PHONY 声明
	for _, raw := range lines {
		line := strings.TrimRight(raw, "\r")
		if m := rePhony.FindStringSubmatch(line); m != nil {
			for _, name := range strings.Fields(m[1]) {
				phonySet[name] = true
			}
		}
	}

	// 第二遍：逐行解析
	var currentTarget *Target
	var rawBuf strings.Builder

	flushRaw := func() {
		s := rawBuf.String()
		if s != "" {
			doc.RawBlocks = append(doc.RawBlocks, RawBlock{Content: s})
			rawBuf.Reset()
		}
	}

	for _, raw := range lines {
		line := strings.TrimRight(raw, "\r")

		// 命令行：以 Tab 开头，属于当前 Target
		if strings.HasPrefix(line, "\t") && currentTarget != nil {
			currentTarget.Commands = append(currentTarget.Commands, line[1:])
			continue
		}

		// 空行：结束当前 target 上下文（但不强制）
		if strings.TrimSpace(line) == "" {
			currentTarget = nil
			rawBuf.WriteString(line + "\n")
			continue
		}

		// .PHONY 行：已在第一遍处理，跳过（不放入 RawBlock）
		if rePhony.MatchString(line) {
			currentTarget = nil
			continue
		}

		// 注释行：放入 RawBlock
		if strings.HasPrefix(strings.TrimSpace(line), "#") {
			currentTarget = nil
			rawBuf.WriteString(line + "\n")
			continue
		}

		// 变量定义行
		if m := reVariable.FindStringSubmatch(line); m != nil {
			flushRaw()
			currentTarget = nil
			doc.Variables = append(doc.Variables, Variable{
				Name:     m[1],
				Operator: m[2],
				Value:    strings.TrimSpace(m[3]),
			})
			continue
		}

		// Target 定义行
		if m := reTarget.FindStringSubmatch(line); m != nil {
			flushRaw()
			name := strings.TrimSpace(m[1])
			depsRaw := strings.Fields(m[2])
			deps := make([]string, 0, len(depsRaw))
			for _, d := range depsRaw {
				d = strings.TrimSpace(d)
				if d != "" {
					deps = append(deps, d)
				}
			}
			t := Target{
				Name:     name,
				Deps:     deps,
				Commands: []string{},
			}
			doc.Targets = append(doc.Targets, t)
			currentTarget = &doc.Targets[len(doc.Targets)-1]
			continue
		}

		// 无法识别的行 → RawBlock
		currentTarget = nil
		rawBuf.WriteString(line + "\n")
	}

	flushRaw()

	// 根据 phony 集合设置 IsPhony
	for i := range doc.Targets {
		if phonySet[doc.Targets[i].Name] {
			doc.Targets[i].IsPhony = true
		}
	}

	return doc, nil
}
