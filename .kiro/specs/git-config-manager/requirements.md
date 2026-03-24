# 需求文档

## 简介

Git Config Manager 是一个可视化 Git 仓库配置管理模块，与"工具盒子"（Toolbox）同级别，隶属于 modules 导航。

用户通常管理多个不同平台（GitHub、GitLab、Gitee 等）的 Git 仓库，每个仓库都有项目级配置文件（`.git/config`）。传统方式（`git config` 命令或手动编辑文件）操作繁琐、容易出错。本模块提供可视化界面，让用户能够方便地浏览、编辑和管理多个 Git 仓库的配置，无需手动操作命令行或文本编辑器。

**范围说明：**
- 仅管理项目级配置（`.git/config`），不涉及全局配置（`~/.gitconfig`）
- 不支持跨仓库批量操作

## 词汇表

- **GitConfigManager**：本功能模块的系统名称，负责所有 Git 配置管理操作
- **Repository**：一个 Git 仓库，包含 `.git/config` 配置文件
- **ConfigFile**：Git 仓库的项目级配置文件，路径为 `<仓库根目录>/.git/config`，采用 INI 格式
- **ConfigSection**：配置文件中的一个节（section），如 `[core]`、`[remote "origin"]`、`[branch "main"]`
- **ConfigEntry**：配置文件中某个节下的一个键值对，如 `repositoryformatversion = 0`
- **KnownKey**：系统内置的 Git 配置键定义，包含键名、所属节、数据类型、默认值和中英文释义说明
- **Parser**：负责将 `.git/config` 文件内容解析为结构化数据的组件
- **Pretty_Printer**：负责将结构化配置数据序列化回合法 `.git/config` 文件格式的组件
- **QuickPanel**：仓库配置视图顶部的常用配置快捷面板，展示用户自定义的固定配置项

---

## 需求

### 需求 1：仓库管理

**用户故事：** 作为开发者，我希望能够添加、查看和删除受管理的 Git 仓库，以便集中管理多个仓库的配置。

#### 验收标准

1. THE GitConfigManager SHALL 在持久化存储（bbolt）中维护一个仓库列表，每条记录包含仓库名称、仓库根目录路径和平台标签（如 GitHub、GitLab、Gitee 或自定义）。
2. WHEN 用户提交添加仓库请求时，THE GitConfigManager SHALL 验证指定路径下存在 `.git/config` 文件，验证通过后将仓库信息写入存储。
3. IF 用户提交的仓库路径不存在 `.git/config` 文件，THEN THE GitConfigManager SHALL 返回明确的错误提示，说明该路径不是有效的 Git 仓库。
4. IF 用户提交的仓库路径已存在于仓库列表中，THEN THE GitConfigManager SHALL 返回重复添加的错误提示，不写入重复记录。
5. WHEN 用户删除一个仓库记录时，THE GitConfigManager SHALL 仅从仓库列表中移除该记录，不修改磁盘上的任何文件。
6. THE GitConfigManager SHALL 支持用户通过系统文件选择对话框选择仓库根目录，以降低手动输入路径出错的概率。

---

### 需求 2：解析 Git 配置文件

**用户故事：** 作为开发者，我希望系统能够正确解析 `.git/config` 文件，以便在界面上展示结构化的配置内容。

#### 验收标准

1. WHEN 用户选择一个已添加的仓库时，THE Parser SHALL 读取该仓库的 `.git/config` 文件并将其解析为结构化的 ConfigSection 和 ConfigEntry 列表。
2. WHEN 解析到带子键的节头（如 `[remote "origin"]`）时，THE Parser SHALL 将节名称和子键分别提取，保留原始引号内的子键值。
3. IF `.git/config` 文件不存在或无法读取，THEN THE Parser SHALL 返回包含具体原因的错误信息。
4. IF `.git/config` 文件包含语法错误，THEN THE Parser SHALL 返回包含行号和错误描述的错误信息。
5. THE Pretty_Printer SHALL 将结构化的 ConfigSection 和 ConfigEntry 列表序列化为符合 Git INI 格式规范的文本内容。
6. FOR ALL 合法的 ConfigFile 结构，THE GitConfigManager SHALL 满足：解析后再序列化再解析，得到的结构与原始解析结果语义等价（往返属性）。

---

### 需求 3：可视化查看配置

**用户故事：** 作为开发者，我希望通过可视化界面查看 Git 仓库的配置内容，以便快速了解仓库的配置状态。

#### 验收标准

1. WHEN 用户选择一个仓库后，THE GitConfigManager SHALL 在界面上以分节（section）形式展示所有 ConfigSection 及其下的 ConfigEntry 键值对。
2. THE GitConfigManager SHALL 在界面上展示仓库的平台标签、仓库路径和最后修改时间。
3. WHEN 配置文件内容发生变化时，THE GitConfigManager SHALL 在用户下次打开该仓库配置时展示最新内容。
4. THE GitConfigManager SHALL 支持在配置视图中按关键词搜索 ConfigEntry 的键或值，搜索结果实时高亮匹配项。
5. WHEN 用户将鼠标悬停在某个 ConfigEntry 的键名上时，THE GitConfigManager SHALL 展示该键对应的 KnownKey 释义说明（中英文），若该键不在 KnownKey 列表中则不展示释义。
6. THE GitConfigManager SHALL 在每个 ConfigEntry 行内以图标或角标形式标识该键是否存在于 KnownKey 列表中，以便用户区分标准键与自定义键。

---

### 需求 4：编辑配置条目

**用户故事：** 作为开发者，我希望通过可视化界面编辑 Git 配置的键值对，以便无需手动操作命令行或文本编辑器。

#### 验收标准

1. WHEN 用户修改某个 ConfigEntry 的值并提交保存时，THE GitConfigManager SHALL 将修改后的内容通过 Pretty_Printer 序列化并写回磁盘上的 `.git/config` 文件。
2. WHEN 用户新增一个 ConfigEntry 时，THE GitConfigManager SHALL 展示一个带搜索功能的 KnownKey 选择器，列出当前节下所有可用的 KnownKey，用户可通过关键词搜索过滤后选择，也可手动输入自定义键名。
3. WHEN 用户在 KnownKey 选择器中选中某个 KnownKey 时，THE GitConfigManager SHALL 自动填充键名，并在输入框旁展示该键的释义说明和默认值（如有）。
4. THE GitConfigManager SHALL 验证新增 ConfigEntry 的键名不为空且不与同节内已有键名重复，验证通过后将新条目写入配置文件。
5. IF 用户提交的 ConfigEntry 键名为空，THEN THE GitConfigManager SHALL 拒绝保存并提示键名不能为空。
6. IF 用户提交的 ConfigEntry 键名与同节内已有键名重复，THEN THE GitConfigManager SHALL 拒绝保存并提示键名已存在。
7. WHEN 用户删除一个 ConfigEntry 时，THE GitConfigManager SHALL 从配置文件中移除该键值对并将更新后的内容写回磁盘。
8. IF 写入 `.git/config` 文件时发生 I/O 错误，THEN THE GitConfigManager SHALL 保留原始文件内容不变，并向用户返回包含错误原因的提示。

---

### 需求 5：管理配置节（Section）

**用户故事：** 作为开发者，我希望能够新增和删除配置节，以便灵活管理仓库的配置结构。

#### 验收标准

1. WHEN 用户新增一个 ConfigSection 时，THE GitConfigManager SHALL 验证节名称不为空且不与已有节名称重复，验证通过后将新节写入配置文件。
2. IF 用户提交的 ConfigSection 名称为空，THEN THE GitConfigManager SHALL 拒绝保存并提示节名称不能为空。
3. IF 用户提交的 ConfigSection 名称与已有节名称重复，THEN THE GitConfigManager SHALL 拒绝保存并提示节名称已存在。
4. WHEN 用户删除一个 ConfigSection 时，THE GitConfigManager SHALL 同时删除该节及其下所有 ConfigEntry，并将更新后的内容写回磁盘。

---

### 需求 6：快速操作面板（QuickPanel）

**用户故事：** 作为开发者，我希望能够自定义一个常用配置快捷面板，固定展示我最关心的配置项，以便快速查看和修改，提高日常操作效率。

#### 验收标准

1. THE GitConfigManager SHALL 在仓库配置视图顶部展示 QuickPanel，默认固定显示 `user.name`、`user.email` 和 `remote "origin".url` 三个配置项的当前值。
2. THE GitConfigManager SHALL 支持用户自定义 QuickPanel 中展示的配置项列表，用户可添加任意 `section.key` 格式的配置项，也可移除已有项，自定义配置持久化存储于 bbolt，按仓库独立保存。
3. WHEN 用户在 QuickPanel 中修改某个配置项的值并提交时，THE GitConfigManager SHALL 将修改同步写入 `.git/config` 文件，并刷新完整配置视图。
4. WHERE 某个 QuickPanel 配置项在当前仓库的 `.git/config` 中不存在，THE GitConfigManager SHALL 以占位符形式展示该项，并允许用户填写后新增到配置文件。
5. WHEN 用户向 QuickPanel 添加配置项时，THE GitConfigManager SHALL 展示带搜索功能的 KnownKey 选择器，用户可搜索选择或手动输入 `section.key` 格式的自定义项。
6. THE GitConfigManager SHALL 支持用户通过拖拽调整 QuickPanel 中配置项的显示顺序。

---

### 需求 7：KnownKey 内置键定义库

**用户故事：** 作为开发者，我希望系统内置常见 Git 配置键的定义和释义，以便在新增配置时获得引导，降低配置出错的概率。

#### 验收标准

1. THE GitConfigManager SHALL 内置一份 KnownKey 定义列表，覆盖以下常用节的标准键：`core`、`user`、`remote`、`branch`、`merge`、`pull`、`push`、`fetch`、`diff`、`merge`、`rebase`、`http`、`https`、`credential`、`alias`。
2. 每个 KnownKey 定义 SHALL 包含以下字段：键名（key）、所属节（section）、数据类型（string/bool/int/enum）、默认值（可选）、中文释义、英文释义。
3. THE GitConfigManager SHALL 在新增 ConfigEntry 的 KnownKey 选择器中，按节分组展示 KnownKey 列表，支持按键名或释义关键词进行模糊搜索过滤。
4. WHEN 用户选中某个 KnownKey 时，THE GitConfigManager SHALL 在选择器内展示该键的完整释义、数据类型和默认值，帮助用户理解该配置项的作用。
5. THE GitConfigManager SHALL 对 enum 类型的 KnownKey 提供可选值下拉列表（如 `core.autocrlf` 的可选值为 `true`、`false`、`input`），而非自由文本输入框。

---

### 需求 8：导航集成

**用户故事：** 作为用户，我希望 Git Config Manager 作为独立模块出现在主导航中，与工具盒子同级，以便快速访问。

#### 验收标准

1. THE GitConfigManager SHALL 在 modules 导航列表中注册为独立模块，路径为 `/git-config`，与 `/toolbox` 同级。
2. THE GitConfigManager SHALL 在 `frontend/src/config/modules.ts` 中添加对应的模块配置项，包含 id、nameKey、path、icon 和 theme 字段。
3. THE GitConfigManager SHALL 在 `frontend/src/locales/zh-CN.ts` 和 `frontend/src/locales/en-US.ts` 中添加对应的 i18n 翻译键。
