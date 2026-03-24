# 需求文档

## 简介

Makefile Visual Editor 是一个可视化 Makefile 编辑器模块，与"工具盒子"（Toolbox）同级别，隶属于 modules 导航。

用户的每个项目通常都需要 Makefile 来管理构建、部署、重启等自动化任务，但 Makefile 语法较为特殊（缩进必须使用 Tab、变量引用、依赖关系等），不熟悉的用户每次都需要查阅文档或复制模板。本模块提供可视化界面，让用户能够通过图形化方式新增、编辑 Makefile 的 Target（构建目标）及其依赖关系，并支持常用场景（多平台构建、部署、重启等）的预置模板，无需手动编写 Makefile 语法。

**范围说明：**
- 支持打开本地已有 Makefile 文件进行可视化编辑
- 支持从空白或预置模板新建 Makefile 文件
- 不支持 Makefile 的完整语法调试或执行

## 词汇表

- **MakefileEditor**：本功能模块的系统名称，负责所有 Makefile 可视化编辑操作
- **Makefile**：项目根目录下的构建配置文件，包含若干 Target 定义，遵循 GNU Make 语法规范
- **Target**：Makefile 中的一个构建目标，由目标名称、依赖列表和命令列表组成，格式为 `target: deps\n\tcommand`
- **Dependency**：一个 Target 所依赖的其他 Target 名称列表，构成有向无环图（DAG）
- **Command**：Target 下执行的 Shell 命令，每行必须以 Tab 字符开头
- **Variable**：Makefile 中定义的变量，格式为 `NAME = value` 或 `NAME := value`
- **Parser**：负责将 Makefile 文件文本内容解析为结构化数据（Target 列表、Variable 列表）的组件
- **Pretty_Printer**：负责将结构化数据序列化回合法 Makefile 文本内容的组件
- **DependencyGraph**：以有向图形式展示各 Target 之间依赖关系的可视化组件
- **Template**：预置的 Makefile 模板，包含常用场景（多平台构建、部署、重启等）的 Target 定义集合
- **PhonyTarget**：声明为 `.PHONY` 的 Target，表示该目标不对应实际文件，如 `clean`、`build`、`deploy`

---

## 需求

### 需求 1：文件管理

**用户故事：** 作为开发者，我希望能够打开本地已有的 Makefile 文件或新建 Makefile 文件，以便对其进行可视化编辑。

#### 验收标准

1. THE MakefileEditor SHALL 支持用户通过系统文件选择对话框选择本地 Makefile 文件（文件名为 `Makefile` 或 `makefile` 或 `*.mk`）进行加载。
2. THE MakefileEditor SHALL 支持用户指定目录路径，在该目录下新建一个空白 Makefile 文件。
3. WHERE 用户选择从模板新建，THE MakefileEditor SHALL 展示预置 Template 列表供用户选择，选择后以该模板内容初始化新文件。
4. THE MakefileEditor SHALL 在界面上展示当前编辑文件的完整路径和文件名。
5. WHEN 用户保存时，THE MakefileEditor SHALL 将当前编辑内容通过 Pretty_Printer 序列化后写入对应的磁盘文件。
6. IF 写入文件时发生 I/O 错误，THEN THE MakefileEditor SHALL 保留原始文件内容不变，并向用户返回包含错误原因的提示。
7. THE MakefileEditor SHALL 在持久化存储（bbolt）中记录最近打开的文件路径列表（最多 10 条），供用户快速重新打开。

---

### 需求 2：解析 Makefile 文件

**用户故事：** 作为开发者，我希望系统能够正确解析 Makefile 文件，以便在界面上展示结构化的 Target 和 Variable 信息。

#### 验收标准

1. WHEN 用户打开一个 Makefile 文件时，THE Parser SHALL 将文件内容解析为结构化的 Target 列表和 Variable 列表。
2. WHEN 解析到 Target 定义时，THE Parser SHALL 提取目标名称、Dependency 列表和 Command 列表，并识别该 Target 是否被声明为 PhonyTarget。
3. WHEN 解析到变量定义时，THE Parser SHALL 提取变量名称、赋值操作符（`=`、`:=`、`?=`、`+=`）和变量值。
4. IF Makefile 文件不存在或无法读取，THEN THE Parser SHALL 返回包含具体原因的错误信息。
5. IF Makefile 文件包含 Parser 无法识别的语法结构，THEN THE Parser SHALL 将该部分作为原始文本块保留，不丢失内容，并在界面上标注为"原始内容"。
6. THE Pretty_Printer SHALL 将结构化的 Target 列表和 Variable 列表序列化为符合 GNU Make 语法规范的 Makefile 文本，Command 行必须以 Tab 字符开头。
7. FOR ALL 合法的 Makefile 结构，THE MakefileEditor SHALL 满足：解析后再序列化再解析，得到的结构与原始解析结果语义等价（往返属性）。

---

### 需求 3：可视化查看依赖关系

**用户故事：** 作为开发者，我希望通过可视化图形直观地查看各 Target 之间的依赖关系，以便快速理解构建流程。

#### 验收标准

1. WHEN 用户打开 Makefile 文件后，THE DependencyGraph SHALL 以有向图形式展示所有 Target 节点及其依赖关系连线。
2. THE DependencyGraph SHALL 以不同颜色或图标区分 PhonyTarget 和普通 Target。
3. WHEN 用户点击 DependencyGraph 中的某个 Target 节点时，THE MakefileEditor SHALL 在右侧面板中展示该 Target 的详细信息（名称、依赖列表、命令列表）。
4. THE DependencyGraph SHALL 支持缩放和平移操作，以便用户在 Target 数量较多时仍能清晰查看全图。
5. IF Makefile 中存在循环依赖（Dependency 形成环），THEN THE MakefileEditor SHALL 在图中以红色高亮标注循环依赖的节点和连线，并展示警告提示。

---

### 需求 4：编辑 Target

**用户故事：** 作为开发者，我希望通过可视化界面新增、编辑和删除 Makefile 的 Target，以便无需手动编写 Makefile 语法。

#### 验收标准

1. THE MakefileEditor SHALL 提供新增 Target 的表单界面，包含目标名称输入框、是否声明为 PhonyTarget 的开关、Dependency 多选列表（从已有 Target 中选择）和 Command 多行编辑区。
2. WHEN 用户提交新增 Target 时，THE MakefileEditor SHALL 验证目标名称不为空且不与已有 Target 名称重复，验证通过后将新 Target 写入文件。
3. IF 用户提交的 Target 名称为空，THEN THE MakefileEditor SHALL 拒绝保存并提示目标名称不能为空。
4. IF 用户提交的 Target 名称与已有 Target 名称重复，THEN THE MakefileEditor SHALL 拒绝保存并提示目标名称已存在。
5. WHEN 用户编辑已有 Target 并提交保存时，THE MakefileEditor SHALL 将修改后的内容序列化并写回磁盘文件。
6. WHEN 用户删除一个 Target 时，THE MakefileEditor SHALL 同时从 `.PHONY` 声明中移除该 Target（如有），并将更新后的内容写回磁盘。
7. THE MakefileEditor SHALL 在 Command 编辑区提供语法提示，自动处理 Tab 缩进，确保用户输入的命令行在序列化时以 Tab 字符开头。
8. WHEN 用户在 Dependency 列表中选择依赖时，THE MakefileEditor SHALL 实时检测是否会产生循环依赖，IF 检测到循环依赖，THEN THE MakefileEditor SHALL 禁止选择并提示循环依赖路径。

---

### 需求 5：编辑变量（Variable）

**用户故事：** 作为开发者，我希望通过可视化界面管理 Makefile 中的变量定义，以便统一维护构建参数。

#### 验收标准

1. THE MakefileEditor SHALL 提供变量列表视图，展示所有已定义的 Variable，包含变量名、赋值操作符和变量值。
2. THE MakefileEditor SHALL 支持新增 Variable，用户需填写变量名、选择赋值操作符（`=`、`:=`、`?=`、`+=`）和变量值。
3. IF 用户提交的 Variable 名称为空，THEN THE MakefileEditor SHALL 拒绝保存并提示变量名不能为空。
4. IF 用户提交的 Variable 名称与已有变量名重复，THEN THE MakefileEditor SHALL 拒绝保存并提示变量名已存在。
5. WHEN 用户修改或删除 Variable 时，THE MakefileEditor SHALL 将更新后的内容写回磁盘文件。

---

### 需求 6：预置模板库

**用户故事：** 作为开发者，我希望系统内置常用场景的 Makefile 模板，以便快速生成符合需求的 Makefile，无需从零编写。

#### 验收标准

1. THE MakefileEditor SHALL 内置以下场景的预置 Template：Go 项目多平台构建（linux/windows/darwin）、Docker 镜像构建与推送、服务部署与重启、前端项目构建（npm/yarn）、通用清理（clean）。
2. 每个 Template SHALL 包含以下字段：模板名称、适用场景描述、预置 Variable 列表和预置 Target 列表（含 PhonyTarget 声明）。
3. THE MakefileEditor SHALL 支持用户将当前编辑的 Makefile 内容另存为自定义 Template，自定义 Template 持久化存储于 bbolt。
4. WHEN 用户选择一个 Template 时，THE MakefileEditor SHALL 展示该模板的预览内容（序列化后的 Makefile 文本），用户确认后再应用。
5. WHERE 用户选择将模板内容合并到当前文件，THE MakefileEditor SHALL 仅追加模板中当前文件不存在的 Target 和 Variable，不覆盖已有内容。

---

### 需求 7：原始文本编辑模式

**用户故事：** 作为开发者，我希望在可视化编辑之外，还能直接查看和编辑 Makefile 的原始文本内容，以便处理可视化模式无法覆盖的复杂语法。

#### 验收标准

1. THE MakefileEditor SHALL 提供原始文本编辑模式，以代码编辑器形式展示当前 Makefile 的完整文本内容，支持语法高亮。
2. WHEN 用户在原始文本编辑模式下修改内容并切换回可视化模式时，THE MakefileEditor SHALL 重新解析文本内容并刷新可视化视图。
3. IF 原始文本内容存在 Parser 无法解析的语法，THEN THE MakefileEditor SHALL 在可视化模式中保留原始文本块，不丢失内容。
4. THE MakefileEditor SHALL 保持可视化模式与原始文本模式的内容实时同步，WHEN 用户在可视化模式下修改 Target 或 Variable 时，原始文本视图 SHALL 同步更新。

---

### 需求 8：导航集成

**用户故事：** 作为用户，我希望 Makefile Visual Editor 作为独立模块出现在主导航中，与工具盒子同级，以便快速访问。

#### 验收标准

1. THE MakefileEditor SHALL 在 modules 导航列表中注册为独立模块，路径为 `/makefile-editor`，与 `/toolbox` 同级。
2. THE MakefileEditor SHALL 在 `frontend/src/config/modules.ts` 中添加对应的模块配置项，包含 id、nameKey、path、icon 和 theme 字段。
3. THE MakefileEditor SHALL 在 `frontend/src/locales/zh-CN.ts` 和 `frontend/src/locales/en-US.ts` 中添加对应的 i18n 翻译键。
