# HappyTools 发布指南

本文档说明如何使用 git-chglog 和发布脚本进行版本发布。

## 📋 前置要求

### 1. 安装 git-chglog

```bash
# macOS
brew install git-chglog

# Windows (使用 Scoop)
scoop install git-chglog

# 或使用 Go 安装
go install github.com/git-chglog/git-chglog/cmd/git-chglog@latest
```

### 2. 验证安装

```bash
git-chglog --version
```

---

## 🚀 发布流程

### 方式 1：使用发布脚本（推荐）

#### Windows (PowerShell)
```powershell
# 发布正式版本
.\release.ps1 v1.0.0

# 发布预发布版本
.\release.ps1 v1.0.0-beta.1

# 预览发布流程（不执行实际操作）
.\release.ps1 v1.0.0 -DryRun
```

#### Linux/macOS (Bash)
```bash
# 添加执行权限
chmod +x release.sh

# 发布正式版本
./release.sh v1.0.0

# 发布预发布版本
./release.sh v1.0.0-beta.1

# 预览发布流程
./release.sh v1.0.0 --dry-run
```

---

### 方式 2：手动发布

```bash
# 1. 确保所有更改已提交
git add .
git commit -m "feat: 添加新功能"

# 2. 生成 CHANGELOG（使用 --next-tag 预览新标签）
git-chglog --next-tag v1.0.0 -o CHANGELOG.md

# 3. 提交 CHANGELOG
git add CHANGELOG.md
git commit -m "build: update v1.0.0"

# 4. 创建标签
git tag v1.0.0

# 5. 推送代码和标签
git push origin main
git push origin v1.0.0
```

---

## 🎯 --next-tag 参数说明

### **什么是 --next-tag？**

`--next-tag` 参数允许你在创建标签之前预览 CHANGELOG 的内容。

### **优势**

1. ✅ **只需生成一次 CHANGELOG** - 不需要两次生成
2. ✅ **流程更简单** - 减少步骤
3. ✅ **避免使用 --amend** - 提交历史更清晰

### **对比**

#### **传统方式（需要两次生成）**
```bash
# 第一次生成（不包含新标签）
git-chglog -o CHANGELOG.md
git commit -m "docs: update CHANGELOG"

# 创建标签
git tag v1.0.0

# 第二次生成（包含新标签）
git-chglog -o CHANGELOG.md
git commit --amend --no-edit
```

#### **使用 --next-tag（只需一次）**
```bash
# 生成 CHANGELOG（预览新标签）
git-chglog --next-tag v1.0.0 -o CHANGELOG.md
git commit -m "docs: update CHANGELOG for v1.0.0"

# 创建标签
git tag v1.0.0
```

---

## 📝 提交信息规范

使用规范的提交信息，git-chglog 会自动分类：

### 提交类型

| 类型 | 说明 | CHANGELOG 分类 |
|------|------|---------------|
| `feat` | 新功能 | 新功能 |
| `fix` | 修复 bug | 修复 |
| `perf` | 性能优化 | 性能优化 |
| `refactor` | 重构代码 | 重构 |
| `docs` | 文档更新 | 文档 |
| `test` | 添加测试 | 测试 |
| `chore` | 构建/工具链更新 | 杂项 |

### 示例

```bash
# 新功能
git commit -m "feat: 添加用户认证功能"

# 修复
git commit -m "fix: 修复登录页面样式问题"

# 性能优化
git commit -m "perf: 优化启动速度"

# 文档
git commit -m "docs: 更新 README"

# 重构
git commit -m "refactor: 重构用户模块"

# 测试
git commit -m "test: 添加单元测试"

# 杂项
git commit -m "chore: 更新依赖版本"
```

---

## 🎯 版本号规则

遵循语义化版本（Semantic Versioning）：

```
MAJOR.MINOR.PATCH

MAJOR: 不兼容的 API 变更
MINOR: 向后兼容的功能新增
PATCH: 向后兼容的问题修复
```

### 示例

```bash
# 补丁版本（修复 bug）
v1.0.1

# 次版本（新功能）
v1.1.0

# 主版本（重大变更）
v2.0.0

# 预发布版本
v1.0.0-beta.1
v1.0.0-rc.1
```

---

## 📦 发布后的流程

### 1. GitHub Actions 自动执行

推送标签后，GitHub Actions 会自动：

- ✅ 构建 Windows 和 macOS 应用
- ✅ 从 CHANGELOG.md 提取内容
- ✅ 生成发布说明
- ✅ 创建 GitHub Release
- ✅ 上传安装程序

### 2. 查看发布进度

```
https://github.com/Aliuyanfeng/happytools/actions
```

### 3. 查看发布结果

```
https://github.com/Aliuyanfeng/happytools/releases
```

---

## 🔧 配置文件说明

### `.chglog/config.yml`

git-chglog 的配置文件，定义了：
- 提交类型过滤
- 分类标题映射
- 提交信息解析规则

### `.chglog/CHANGELOG.tpl.md`

CHANGELOG 的模板文件，定义了：
- 输出格式
- 版本分组
- 提交分组

---

## 🎨 CHANGELOG 示例

```markdown
## [v1.0.0] - 2026-03-09

### 新功能
- 添加用户认证功能
- 添加数据导出功能

### 修复
- 修复登录页面样式问题
- 修复数据保存失败的问题

### 性能优化
- 优化启动速度
- 优化内存占用

### 重构
- 重构用户模块

### 文档
- 更新 README

### 测试
- 添加单元测试
```

---

## 🐛 常见问题

### Q: 如何预览 CHANGELOG 而不生成文件？

```bash
git-chglog
```

### Q: 如何生成特定版本的 CHANGELOG？

```bash
# 生成 v1.0.0 到 v1.1.0 之间的变更
git-chglog v1.0.0..v1.1.0

# 生成 v1.0.0 之后的所有变更
git-chglog v1.0.0..
```

### Q: 如何只生成未发布的变更？

```bash
# 使用 --next-tag 预览下一个版本
git-chglog --next-tag v1.1.0
```

### Q: 如何自定义 CHANGELOG 格式？

编辑 `.chglog/CHANGELOG.tpl.md` 模板文件。

### Q: --next-tag 和直接生成有什么区别？

- **直接生成**：只包含已存在的标签
- **--next-tag**：可以预览即将创建的标签内容

---

## 📚 参考资料

- [git-chglog 官方文档](https://github.com/git-chglog/git-chglog)
- [语义化版本](https://semver.org/lang/zh-CN/)
- [Conventional Commits](https://www.conventionalcommits.org/zh-hans/)

---

**最后更新**: 2026-03-09
