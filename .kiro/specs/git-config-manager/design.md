# 设计文档：Git Config Manager

## 概述

Git Config Manager 是一个与 Toolbox 同级的独立模块，采用左右分栏布局：左侧仓库列表，右侧配置编辑区。后端使用 Go 实现 INI 解析、bbolt 持久化和文件读写；前端使用 Vue 3 + Ant Design Vue 实现可视化编辑界面。

---

## 架构总览

```
frontend/src/views/GitConfig/
  GitConfig.vue          # 主视图（左右分栏容器）
  RepoList.vue           # 左侧仓库列表面板
  ConfigEditor.vue       # 右侧配置编辑区（含 QuickPanel + Section 列表）
  QuickPanel.vue         # 快速操作面板组件
  SectionCard.vue        # 单个 Section 折叠卡片
  KnownKeySelector.vue   # 带搜索的 KnownKey 选择器弹窗

backend/services/gitconfig/
  gitconfigservice.go    # 主 Service，暴露给 Wails 的所有方法
  parser.go              # .git/config INI 解析器 & Pretty_Printer
  knownkeys.go           # KnownKey 内置定义库

backend/store/
  gitconfig_options.go   # bbolt CRUD：仓库记录 + QuickPanel 配置
```

---

## 数据模型

### 后端 Go 结构体

```go
// Repository 受管理的 Git 仓库记录（持久化到 bbolt）
type Repository struct {
    ID        string    `json:"id"`        // UUID
    Name      string    `json:"name"`      // 显示名称
    Path      string    `json:"path"`      // 仓库根目录绝对路径
    Platform  string    `json:"platform"`  // "github" | "gitlab" | "gitee" | "custom"
    CreatedAt time.Time `json:"createdAt"`
}

// ConfigSection 解析后的配置节
type ConfigSection struct {
    Name    string        `json:"name"`    // 节名，如 "remote"
    SubKey  string        `json:"subKey"`  // 子键，如 "origin"（无则为空）
    Entries []ConfigEntry `json:"entries"`
}

// ConfigEntry 配置节下的键值对
type ConfigEntry struct {
    Key   string `json:"key"`
    Value string `json:"value"`
}

// KnownKey 内置键定义
type KnownKey struct {
    Section     string   `json:"section"`
    Key         string   `json:"key"`
    Type        string   `json:"type"`        // "string" | "bool" | "int" | "enum"
    Default     string   `json:"default"`     // 可选
    EnumValues  []string `json:"enumValues"`  // type=enum 时有效
    DescZh      string   `json:"descZh"`
    DescEn      string   `json:"descEn"`
}

// QuickPanelItem QuickPanel 中的一个配置项
type QuickPanelItem struct {
    Section string `json:"section"` // 如 "remote"
    SubKey  string `json:"subKey"`  // 如 "origin"（无则为空）
    Key     string `json:"key"`     // 如 "url"
    Order   int    `json:"order"`
}
```

### bbolt Bucket 设计

| Bucket | Key | Value |
|--------|-----|-------|
| `git_repos` | `repo.ID` | JSON(Repository) |
| `git_quick_panels` | `repoID` | JSON([]QuickPanelItem) |

---

## 后端设计

### `parser.go` — INI 解析器

**Parse(content string) ([]ConfigSection, error)**
- 逐行扫描，识别 `[section]`、`[section "subkey"]`、`key = value`、注释（`#` / `;`）
- 带子键的节头：正则 `^\[(\w+)\s+"([^"]+)"\]$` 分别提取 Name 和 SubKey
- 语法错误时返回行号和描述

**Serialize(sections []ConfigSection) string**
- 按节顺序输出，节头格式：有 SubKey 则 `[name "subkey"]`，否则 `[name]`
- 每个 entry 缩进一个 tab：`\tkey = value`
- 节之间空一行

### `knownkeys.go` — KnownKey 定义库

以 Go 切片硬编码，覆盖以下节的常用键：

| 节 | 代表键 |
|----|--------|
| `core` | `repositoryformatversion`, `filemode`, `bare`, `logallrefupdates`, `autocrlf`, `eol`, `ignorecase`, `editor` |
| `user` | `name`, `email`, `signingkey` |
| `remote` | `url`, `fetch`, `pushurl`, `mirror` |
| `branch` | `remote`, `merge`, `rebase` |
| `pull` | `rebase`, `ff` |
| `push` | `default`, `followtags` |
| `merge` | `tool`, `conflictstyle` |
| `rebase` | `autosquash`, `autostash` |
| `diff` | `tool`, `algorithm` |
| `http` | `proxy`, `sslverify`, `sslcainfo` |
| `credential` | `helper` |

**GetKnownKeys() []KnownKey** — 返回全部定义
**GetKnownKeysForSection(section string) []KnownKey** — 按节过滤

### `gitconfigservice.go` — 暴露给 Wails 的方法

```go
// 仓库管理
func (s *GitConfigService) AddRepository(name, path, platform string) (*Repository, error)
func (s *GitConfigService) ListRepositories() ([]*Repository, error)
func (s *GitConfigService) DeleteRepository(id string) error
func (s *GitConfigService) OpenDirectoryDialog() string

// 配置读写
func (s *GitConfigService) LoadConfig(repoID string) ([]ConfigSection, error)
func (s *GitConfigService) SaveEntry(repoID, section, subKey, key, value string) error
func (s *GitConfigService) DeleteEntry(repoID, section, subKey, key string) error
func (s *GitConfigService) AddSection(repoID, section, subKey string) error
func (s *GitConfigService) DeleteSection(repoID, section, subKey string) error

// KnownKey
func (s *GitConfigService) GetKnownKeys() []KnownKey
func (s *GitConfigService) GetKnownKeysForSection(section string) []KnownKey

// QuickPanel
func (s *GitConfigService) GetQuickPanel(repoID string) ([]QuickPanelItem, error)
func (s *GitConfigService) SaveQuickPanel(repoID string, items []QuickPanelItem) error
```

**写文件安全策略**：先写临时文件，成功后原子替换（`os.Rename`），失败时保留原文件不变。

---

## 前端设计

### 路由

在 `frontend/src/router/routes.ts` 新增：

```ts
{
  path: '/git-config',
  name: 'git-config',
  component: () => import('../views/GitConfig/GitConfig.vue')
}
```

### 组件结构

```
GitConfig.vue
├── RepoList.vue          左侧 200px 固定宽度面板
│   ├── 平台筛选 Tag 组
│   ├── a-list（仓库列表）
│   └── 添加仓库 Modal
└── ConfigEditor.vue      右侧弹性宽度
    ├── QuickPanel.vue    顶部快捷面板（可折叠）
    │   ├── 拖拽排序（vuedraggable）
    │   └── 每项：label + 内联编辑 Input
    └── Section 列表
        └── SectionCard.vue × N
            ├── a-collapse-panel（节头 + 操作按钮）
            └── ConfigEntry 行 × N
                ├── 键名（带 KnownKey 角标 + Tooltip 释义）
                ├── 值（内联编辑）
                └── 删除按钮
```

### 状态管理（Pinia store）

```ts
// stores/gitconfig.ts
interface GitConfigStore {
  repos: Repository[]
  activeRepoID: string | null
  sections: ConfigSection[]
  quickPanel: QuickPanelItem[]
  searchKeyword: string
}
```

actions：`loadRepos`, `selectRepo`, `loadConfig`, `saveEntry`, `deleteEntry`, `addSection`, `deleteSection`, `loadQuickPanel`, `saveQuickPanel`

### KnownKeySelector 弹窗

- `a-modal` + `a-input` 搜索框
- `a-list` 按节分组，支持模糊匹配键名/释义
- 选中后展示：类型、默认值、中英文释义
- enum 类型：选中后值输入框替换为 `a-select`

### 搜索高亮

`searchKeyword` 非空时，对每个 ConfigEntry 的 key 和 value 做字符串匹配，命中则用 `<mark>` 包裹高亮文本，不命中的 entry 降低透明度。

---

## 模块注册

### `frontend/src/config/modules.ts`

```ts
{
  id: 'gitConfig',
  nameKey: 'home.modules.gitConfig',
  description: '',
  path: '/git-config',
  icon: 'BranchesOutlined',
  theme: 'green'
}
```

### i18n 键（zh-CN / en-US）

```ts
// zh-CN
gitConfig: {
  title: 'Git 配置管理',
  // 仓库管理
  addRepo: '添加仓库',
  repoName: '仓库名称',
  repoPath: '仓库路径',
  platform: '平台',
  platformGithub: 'GitHub',
  platformGitlab: 'GitLab',
  platformGitee: 'Gitee',
  platformCustom: '自定义',
  selectDirectory: '选择目录',
  deleteRepoConfirm: '确定要移除此仓库吗？（不会删除本地文件）',
  // 配置编辑
  searchPlaceholder: '搜索键或值...',
  addSection: '添加节',
  addEntry: '添加配置项',
  sectionName: '节名称',
  subKey: '子键（可选）',
  key: '键名',
  value: '值',
  deleteSectionConfirm: '确定要删除此节及其所有配置项吗？',
  deleteEntryConfirm: '确定要删除此配置项吗？',
  saveSuccess: '保存成功',
  saveFailed: '保存失败',
  loadFailed: '加载配置失败',
  invalidRepo: '该路径不是有效的 Git 仓库',
  duplicateRepo: '该仓库已存在',
  keyEmpty: '键名不能为空',
  keyDuplicate: '键名已存在',
  sectionEmpty: '节名称不能为空',
  sectionDuplicate: '节名称已存在',
  // QuickPanel
  quickPanel: '快速面板',
  quickPanelHint: '常用配置项快捷访问',
  addQuickItem: '添加快捷项',
  notConfigured: '未配置',
  // KnownKey
  knownKeySelector: '选择配置键',
  knownKeySearch: '搜索键名或说明...',
  knownKeyType: '类型',
  knownKeyDefault: '默认值',
  customKey: '自定义键名',
}
```

---

## 关键交互流程

### 添加仓库

1. 用户点击「添加仓库」→ 弹出 Modal
2. 点击「选择目录」→ 调用 `OpenDirectoryDialog()` 获取路径
3. 填写名称、选择平台 → 提交
4. 后端验证 `<path>/.git/config` 存在 → 写入 bbolt → 返回 Repository
5. 前端刷新仓库列表，自动选中新仓库

### 编辑配置值

1. 用户点击某 ConfigEntry 的值 → 切换为内联 Input
2. 失焦或按 Enter → 调用 `SaveEntry()`
3. 后端：读取文件 → 修改对应 entry → Serialize → 写临时文件 → Rename
4. 前端：更新本地 sections 状态，无需重新加载整个文件

### 新增 ConfigEntry

1. 用户点击节内「+ 添加配置项」→ 打开 KnownKeySelector
2. 搜索/选择 KnownKey 或手动输入键名
3. enum 类型自动切换为下拉选择值
4. 确认 → 调用 `SaveEntry()` → 刷新该节 entries

### QuickPanel 拖拽排序

1. 使用 `vuedraggable` 包裹 QuickPanel 项列表
2. 拖拽结束后更新本地 `order` 字段
3. 调用 `SaveQuickPanel(repoID, items)` 持久化到 bbolt
