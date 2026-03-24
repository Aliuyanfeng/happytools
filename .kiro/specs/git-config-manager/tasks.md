# 实现任务：Git Config Manager

## Task 1: bbolt 存储层 ✅

**文件**: `backend/store/gitconfig_options.go`

- [x] 在 `store.go` 的 `Init()` 中注册 `git_repos` 和 `git_quick_panels` 两个 bucket
- [x] 实现 `SaveRepository / GetAllRepositories / GetRepository / DeleteRepository`
- [x] 实现 `GetQuickPanel / SaveQuickPanel`
- [x] 定义 `Repository`、`QuickPanelItem` 结构体

---

## Task 2: INI 解析器 ✅

**文件**: `backend/services/gitconfig/parser.go`

- [x] 实现 `Parse(content string) ([]ConfigSection, error)`
- [x] 实现 `Serialize(sections []ConfigSection) string`
- [x] 定义 `ConfigSection`、`ConfigEntry` 结构体

---

## Task 3: KnownKey 定义库 ✅

**文件**: `backend/services/gitconfig/knownkeys.go`

- [x] 定义 `KnownKey` 结构体
- [x] 硬编码内置 KnownKey 列表（11 个节，35 个键）
- [x] 实现 `GetKnownKeys / GetKnownKeysForSection`

---

## Task 4: GitConfigService 主服务 ✅

**文件**: `backend/services/gitconfig/gitconfigservice.go`

- [x] 仓库管理方法（AddRepository / ListRepositories / DeleteRepository / OpenDirectoryDialog）
- [x] 配置读写方法（LoadConfig / SaveEntry / DeleteEntry / AddSection / DeleteSection）
- [x] KnownKey 方法（GetKnownKeys / GetKnownKeysForSection）
- [x] QuickPanel 方法（GetQuickPanel / SaveQuickPanel）
- [x] 在 `main.go` 中注册服务

---

## Task 5: 生成 Wails Bindings ✅

- [x] 运行 `wails3 generate bindings`，生成 `gitconfigservice.js`

---

## Task 6: Pinia Store ✅

**文件**: `frontend/src/stores/gitconfig.ts`

- [x] 定义 TypeScript 接口（Repository / ConfigSection / ConfigEntry / QuickPanelItem / KnownKey）
- [x] 定义 store state 和全部 actions
- [x] 实现 getter `filteredSections`

---

## Task 7: KnownKeySelector 组件 ✅

**文件**: `frontend/src/views/GitConfig/KnownKeySelector.vue`

- [x] 搜索框 + 按节分组列表 + 右侧详情面板 + 自定义键名输入

---

## Task 8: SectionCard 组件 ✅

**文件**: `frontend/src/views/GitConfig/SectionCard.vue`

- [x] 节头 + 配置项列表 + 内联编辑（enum 用 select）
- [x] KnownKey 蓝点角标 + Tooltip 释义
- [x] 搜索高亮（mark 标签 + dimmed 透明度）

---

## Task 9: QuickPanel 组件 ✅

**文件**: `frontend/src/views/GitConfig/QuickPanel.vue`

- [x] vuedraggable 拖拽排序
- [x] 内联编辑 + 占位符展示
- [x] 添加/删除项，持久化到 bbolt

---

## Task 10: RepoList 组件 ✅

**文件**: `frontend/src/views/GitConfig/RepoList.vue`

- [x] 平台筛选 Tag 组
- [x] 仓库列表 + hover 删除按钮
- [x] 添加仓库 Modal（目录选择对话框）

---

## Task 11: ConfigEditor 组件 ✅

**文件**: `frontend/src/views/GitConfig/ConfigEditor.vue`

- [x] 工具栏（仓库信息 + 搜索框 + 添加节按钮）
- [x] QuickPanel + SectionCard 列表
- [x] 添加节 Modal

---

## Task 12: GitConfig 主视图 ✅

**文件**: `frontend/src/views/GitConfig/GitConfig.vue`

- [x] 左右分栏布局，onMounted 加载数据

---

## Task 13: 路由 & 模块注册 & i18n ✅

- [x] `frontend/src/router/routes.ts`：新增 `/git-config` 路由
- [x] `frontend/src/config/modules.ts`：新增 `gitConfig` 模块（BranchesOutlined, green）
- [x] `frontend/src/locales/zh-CN.ts`：新增完整翻译键
- [x] `frontend/src/locales/en-US.ts`：新增完整翻译键
