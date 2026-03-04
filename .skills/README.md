# HappyTools Skills 使用指南

## 概述

本目录包含 HappyTools 项目的 Skills 配置，用于提供智能化的开发辅助功能。Skills 是预定义的任务模板，可以帮助开发者快速完成常见的开发任务。

## 目录结构

```
.qoder/skills/
├── skill-metadata.json          # Skills 元数据配置
├── add-feature-module.md        # 添加新功能模块 Skill
├── fix-common-issues.md         # 修复常见问题 Skill
└── README.md                    # 本文档
```

## 已配置的 Skills

### 1. 添加新功能模块 (add-feature-module)

**用途**: 快速添加新的功能模块到 HappyTools 项目

**功能**:
- 创建后端服务（Go）
- 实现数据存储层（bbolt）
- 创建前端页面（Vue 3）
- 配置路由和导航菜单
- 自动生成 TypeScript 绑定

**触发关键词**:
- "添加新功能"
- "创建新模块"
- "新增功能模块"
- "添加模块"

**使用示例**:
```
用户: 我想添加一个笔记管理功能
AI: [自动触发 add-feature-module skill，引导完成模块创建]
```

### 2. 修复常见问题 (fix-common-issues)

**用途**: 快速诊断和修复项目中的常见问题

**功能**:
- 前端绑定未生成
- 数据库操作失败
- 前后端通信异常
- 事件通信问题
- 构建打包失败
- 窗口显示异常

**触发关键词**:
- "修复问题"
- "解决错误"
- "调试问题"
- "排查错误"
- "绑定未生成"
- "数据库错误"

**使用示例**:
```
用户: 前端绑定没有生成怎么办？
AI: [自动触发 fix-common-issues skill，提供解决方案]
```

## Skill 配置格式

每个 skill 在 `skill-metadata.json` 中包含以下字段：

```json
{
  "name": "skill-name",
  "displayName": "Skill 显示名称",
  "description": "Skill 功能描述",
  "version": "1.0.0",
  "author": "作者",
  "tags": ["标签1", "标签2"],
  "triggers": ["触发词1", "触发词2"],
  "filePath": "skill-file.md",
  "enabled": true,
  "priority": 1,
  "createdAt": "2026-02-24",
  "updatedAt": "2026-02-24"
}
```

## 如何添加新的 Skill

### 1. 创建 Skill 文档

在 `.qoder/skills/` 目录下创建新的 Markdown 文件：

```bash
touch .qoder/skills/my-new-skill.md
```

### 2. 编写 Skill 内容

参考现有 skill 的格式，编写详细的执行步骤：

```markdown
# Skill 名称

## Skill 描述
简要描述 skill 的用途和功能

## 适用场景
列出适用的场景和情况

## 执行步骤
详细的执行步骤和代码示例

## 注意事项
需要注意的问题和最佳实践
```

### 3. 注册 Skill

在 `skill-metadata.json` 中添加新的 skill 配置：

```json
{
  "name": "my-new-skill",
  "displayName": "我的新 Skill",
  "description": "Skill 功能描述",
  "version": "1.0.0",
  "author": "Your Name",
  "tags": ["tag1", "tag2"],
  "triggers": ["触发词1", "触发词2"],
  "filePath": "my-new-skill.md",
  "enabled": true,
  "priority": 3,
  "createdAt": "2026-02-24",
  "updatedAt": "2026-02-24"
}
```

## Skill 开发最佳实践

### 1. 命名规范
- 使用小写字母和连字符：`add-feature-module`
- 文件名与 skill name 保持一致
- 显示名称使用中文，简洁明了

### 2. 内容组织
- 提供清晰的步骤说明
- 包含完整的代码示例
- 添加必要的注释和说明
- 提供错误处理和调试建议

### 3. 触发词设计
- 使用常见的问题描述
- 包含中英文关键词
- 避免过于宽泛的触发词
- 保持触发词的准确性

### 4. 代码示例
- 遵循项目开发规范
- 使用实际的代码片段
- 包含错误处理逻辑
- 添加必要的类型注解

## Skill 使用场景

### 开发阶段
- 快速添加新功能模块
- 解决开发过程中的问题
- 学习项目架构和规范

### 维护阶段
- 修复生产环境问题
- 优化现有功能
- 重构代码结构

### 团队协作
- 统一开发流程
- 降低学习成本
- 提高代码质量

## 常见问题

### Q: Skill 如何被触发？
A: Skill 通过关键词触发，当用户输入包含触发词时，AI 会自动识别并执行相应的 skill。

### Q: 可以禁用某个 skill 吗？
A: 可以，在 `skill-metadata.json` 中将对应 skill 的 `enabled` 字段设置为 `false`。

### Q: 如何更新 skill？
A: 直接编辑对应的 Markdown 文件，并更新 `skill-metadata.json` 中的 `updatedAt` 字段。

### Q: Skill 文件可以放在其他目录吗？
A: 目前 skill 文件必须放在 `.qoder/skills/` 目录下，并在 `skill-metadata.json` 中正确配置路径。

## 技术支持

如有问题或建议，请通过以下方式联系：
- 提交 GitHub Issue
- 查阅项目文档
- 联系项目维护者

---

**HappyTools** - 让开发更高效，让协作更简单！
