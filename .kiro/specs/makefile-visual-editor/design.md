# 设计文档：Makefile Visual Editor

## 概述

Makefile Visual Editor 是一个与 Toolbox 同级的独立模块，路径为 `/makefile-editor`。用户可通过可视化界面打开、新建、编辑本地 Makefile 文件，无需手动编写 GNU Make 语法。后端使用 Go 实现 Makefile 解析器（Parser）、序列化器（Pretty_Printer）、循环依赖检测和 bbolt 持久化；前端使用 Vue 3 + TypeScript + Ant Design Vue 实现依赖关系图（DependencyGraph）、Target/Variable 编辑器、模板库和原始文本编辑模式。

---

## 架构总览

```
frontend/src/views/MakefileEditor/
  MakefileEditor.vue      # 主视图（整体布局容器）
  FilePanel.vue           # 左侧文件管理面板（打开/新建/最近文件/保存）
  VariableList.vue        # 变量列表（展示 + 内联编辑/删除）
  VariableForm.vue        # 新增/编辑变量表单
  TargetList.vue          # Target 列表（展示 + 编辑/删除入口）
  TargetForm.vue          # 新增/编辑 Target 表单（含依赖多选 + 命令编辑区）
  DependencyGraph.vue     # 依赖关系有向图（节点 + 连线 + 缩放/平移）
  RawEditor.vue           # 原始文本编辑模式（代码编辑器 + 语法高亮）
  TemplateModal.vue       # 模板库弹窗（预览 + 应用/合并 + 另存为模板）

frontend/src/stores/
  makefileEditor.ts       # Pinia store（编辑器全局状态 + actions）

backend/services/makefile/
  types.go                # 数据结构定义（MakefileDoc、Target、Variable、RawBlock、Template）
  parser.go               # Parse(content string) (*MakefileDoc, error)
  printer.go              # Print(doc *MakefileDoc) string
  templates.go            # 5 个内置模板定义
  makefileservice.go      # 暴露给 Wails 的所有方法

backend/store/
  makefile_options.go     # bbolt CRUD：最近文件列表 + 自定义模板
```

---

## 数据模型

### 后端 Go 结构体

```go
// MakefileDoc 解析后的 Makefile 文档
type MakefileDoc struct {
    Variables []Variable  `json:"variables"`
    Targets   []Target    `json:"targets"`
    RawBlocks []RawBlock  `json:"rawBlocks"` // 无法识别的原始文本块
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
    Deps     []string `json:"deps"`     // 依赖列表，如 ["clean", "fmt"]
    Commands []string `json:"commands"` // 命令列表（不含前导 Tab）
    IsPhony  bool     `json:"isPhony"`  // 是否声明为 .PHONY
}

// RawBlock 无法识别的原始文本块（保留内容，不丢失）
type RawBlock struct {
    Content string `json:"content"` // 原始文本内容
}

// Template Makefile 模板
type Template struct {
    ID          string     `json:"id"`          // UUID（内置模板使用固定 ID）
    Name        string     `json:"name"`        // 模板名称
    Description string     `json:"description"` // 适用场景描述
    IsBuiltin   bool       `json:"isBuiltin"`   // 是否为内置模板
    Doc         MakefileDoc `json:"doc"`        // 模板内容
}
```

### 前端 TypeScript 接口

```ts
// 与后端 Go 结构体一一对应，由 Wails 自动生成绑定
export interface Variable {
  name: string
  operator: '=' | ':=' | '?=' | '+='
  value: string
}

export interface Target {
  name: string
  deps: string[]
  commands: string[]
  isPhony: boolean
}

export interface RawBlock {
  content: string
}

export interface MakefileDoc {
  variables: Variable[]
  targets: Target[]
  rawBlocks: RawBlock[]
}

export interface Template {
  id: string
  name: string
  description: string
  isBuiltin: boolean
  doc: MakefileDoc
}
```

### bbolt Bucket 设计

| Bucket | Key | Value |
|--------|-----|-------|
| `makefile_recent` | `"list"` | JSON([]string)，最多 10 条路径 |
| `makefile_templates` | `template.ID` | JSON(Template) |

---

## 后端设计

### `types.go` — 数据结构

定义上述所有 Go 结构体，附带 JSON 标签，供 Wails 绑定生成 TypeScript 接口。

### `parser.go` — Makefile 解析器

**Parse(content string) (\*MakefileDoc, error)**

解析策略（逐行扫描状态机）：

1. 收集 `.PHONY` 声明行，提取所有 phony target 名称集合
2. 识别变量定义行：正则 `^([A-Za-z_][A-Za-z0-9_]*)\s*(=|:=|\?=|\+=)\s*(.*)$`
3. 识别 Target 定义行：正则 `^([^\s#:][^:]*):(.*)$`（非空格、非注释、非冒号开头）
4. Target 下的命令行：以 Tab 字符开头，去除前导 Tab 后存入 Commands
5. 无法归类的行（include、条件指令等）：追加到当前 RawBlock
6. 解析完成后，根据 phony 集合为每个 Target 设置 `IsPhony` 字段

**错误处理**：文件不存在或无法读取时返回包含具体原因的 error；无法识别的语法不返回错误，而是保留为 RawBlock。

### `printer.go` — Makefile 序列化器

**Print(doc \*MakefileDoc) string**

序列化顺序：

1. 输出所有 Variable 定义：`NAME OPERATOR VALUE\n`
2. 空一行
3. 收集所有 `IsPhony == true` 的 Target 名称，输出 `.PHONY: name1 name2 ...\n`
4. 空一行
5. 逐个输出 Target：
   ```
   name: dep1 dep2
   \tcommand1
   \tcommand2
   ```
6. Target 之间空一行
7. 在末尾输出所有 RawBlock 内容（保留原始文本）

**关键约束**：每条 Command 行必须以 Tab 字符（`\t`）开头，不得使用空格替代。

### `makefileservice.go` — 暴露给 Wails 的方法

```go
type MakefileService struct {
    app *application.App
}

// 文件管理
func (s *MakefileService) OpenFile(path string) (*MakefileDoc, error)
func (s *MakefileService) OpenFileDialog() string
func (s *MakefileService) NewFile(dir string) (*MakefileDoc, error)
func (s *MakefileService) NewFromTemplate(dir string, templateID string) (*MakefileDoc, error)
func (s *MakefileService) SaveFile(path string, doc *MakefileDoc) error
func (s *MakefileService) GetRecentFiles() ([]string, error)

// 模板管理
func (s *MakefileService) GetTemplates() ([]Template, error)
func (s *MakefileService) SaveCustomTemplate(name string, description string, doc *MakefileDoc) error
func (s *MakefileService) DeleteCustomTemplate(id string) error

// 依赖检测
func (s *MakefileService) ValidateDependencies(doc *MakefileDoc) ([][]string, error)
```

**SaveFile 原子写策略**：先将序列化内容写入 `<path>.tmp` 临时文件，成功后调用 `os.Rename` 原子替换；若写临时文件失败，直接返回错误，原始文件不受影响。

**OpenFile 最近文件记录**：每次成功打开文件后，将路径插入最近文件列表头部，去重后截断至最多 10 条，写回 bbolt。

**ValidateDependencies 循环检测**：使用深度优先搜索（DFS）检测有向图中的环，返回所有环路（每个环路为一个 Target 名称列表）。

### `templates.go` — 内置模板定义

