// Package makefile
// @description: Makefile Visual Editor — 数据结构定义
package makefile

// MakefileDoc 解析后的 Makefile 文档
type MakefileDoc struct {
	Variables []Variable `json:"variables"`
	Targets   []Target   `json:"targets"`
	RawBlocks []RawBlock `json:"rawBlocks"`
}

// Variable Makefile 变量定义
type Variable struct {
	Name     string `json:"name"`     // 变量名，如 "BINARY"
	Operator string `json:"operator"` // "=" | ":=" | "?=" | "+="
	Value    string `json:"value"`    // 变量值
}

// Target Makefile 构建目标
type Target struct {
	Name     string   `json:"name"`     // 目标名称，如 "build"
	Deps     []string `json:"deps"`     // 依赖列表
	Commands []string `json:"commands"` // 命令列表（不含前导 Tab）
	IsPhony  bool     `json:"isPhony"`  // 是否声明为 .PHONY
}

// RawBlock 无法识别的原始文本块（保留内容，不丢失）
type RawBlock struct {
	Content string `json:"content"`
}

// Template Makefile 模板
type Template struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	IsBuiltin   bool        `json:"isBuiltin"`
	Doc         MakefileDoc `json:"doc"`
}
