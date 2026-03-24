// Package makefile
// @description: Makefile 序列化器 — 将 MakefileDoc 输出为合法 GNU Make 语法
package makefile

import (
	"fmt"
	"strings"
)

// Print 将 MakefileDoc 序列化为合法的 GNU Make 语法文本。
// 序列化顺序：Variables → 空行 → .PHONY → 空行 → Targets → RawBlocks
func Print(doc *MakefileDoc) string {
	var sb strings.Builder

	// 1. 输出所有变量定义
	for _, v := range doc.Variables {
		fmt.Fprintf(&sb, "%s %s %s\n", v.Name, v.Operator, v.Value)
	}

	// 2. 收集所有 IsPhony 的 Target 名称
	var phonyNames []string
	for _, t := range doc.Targets {
		if t.IsPhony {
			phonyNames = append(phonyNames, t.Name)
		}
	}

	// 3. 空行 + .PHONY 声明（仅当有 phony target 时）
	if len(phonyNames) > 0 {
		if len(doc.Variables) > 0 {
			sb.WriteByte('\n')
		}
		fmt.Fprintf(&sb, ".PHONY: %s\n", strings.Join(phonyNames, " "))
	}

	// 4. 空行 + Targets
	if len(doc.Targets) > 0 {
		sb.WriteByte('\n')
	}
	for i, t := range doc.Targets {
		if i > 0 {
			sb.WriteByte('\n')
		}
		// target: dep1 dep2
		if len(t.Deps) > 0 {
			fmt.Fprintf(&sb, "%s: %s\n", t.Name, strings.Join(t.Deps, " "))
		} else {
			fmt.Fprintf(&sb, "%s:\n", t.Name)
		}
		// \tcommand
		for _, cmd := range t.Commands {
			fmt.Fprintf(&sb, "\t%s\n", cmd)
		}
	}

	// 5. RawBlocks 追加到末尾
	for _, rb := range doc.RawBlocks {
		content := rb.Content
		// 确保 RawBlock 内容以换行结尾
		if content != "" && !strings.HasSuffix(content, "\n") {
			content += "\n"
		}
		sb.WriteString(content)
	}

	return sb.String()
}
