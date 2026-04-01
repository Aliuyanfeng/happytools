# HappyTools

> 基于 Wails v3 + Vue 3 构建的跨平台桌面工具集

[![Release](https://img.shields.io/github/v/release/Aliuyanfeng/happytools)](https://github.com/Aliuyanfeng/happytools/releases)
[![License](https://img.shields.io/github/license/Aliuyanfeng/happytools)](LICENSE)

## 功能

| 模块 | 说明 |
|------|------|
| 🖥️ 系统仪表盘 | 实时监控 CPU、内存、硬盘、网卡信息 |
| ✅ 智能待办 | 待办事项管理，支持分类与优先级 |
| 🔧 工具盒子 | 单位转换、进制转换、时间戳、颜色转换、加密工具、PNG 注入、批量重命名 |
| 🛡️ VirusTotal | 文件病毒扫描，支持批量上传与任务管理 |
| 🌐 网络调试 | TCP/UDP 调试、FTP 工具、DNS 缓存刷新 |
| 📅 日报管理 | 记录每日工作日报，日历视图展示 |
| 🔀 Git 配置管理 | 多仓库 Git 配置管理与快捷操作 |
| 📄 Makefile 编辑器 | 可视化编辑 Makefile，支持依赖图与原始文本模式 |
| 🔄️ 持续更新中 | ...... |

## 技术栈

- **后端**：Go · Wails v3 · bbolt
- **前端**：Vue 3 · TypeScript · Ant Design Vue · Pinia

## 开发

```bash
# 安装依赖
go mod download
cd frontend && npm install && cd ..

# 开发模式（热重载）
wails3 task dev

# 构建
wails3 task build

# 打包
wails3 task package
```

## 下载

前往 [Releases](https://github.com/Aliuyanfeng/happytools/releases) 下载最新版本。

## License

MIT © 2025 Li6
