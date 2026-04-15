# HappyTools

> 基于 Wails v3 + Vue 3 构建的跨平台桌面工具集

[![Release](https://img.shields.io/github/v/release/Aliuyanfeng/happytools)](https://github.com/Aliuyanfeng/happytools/releases)
[![License](https://img.shields.io/github/license/Aliuyanfeng/happytools)](LICENSE)

## 功能

| 模块 | 说明 |
|------|------|
| 🖥️ 系统仪表盘 | 实时监控 CPU、内存、硬盘、网卡信息，支持多核心视图与历史趋势 |
| ✅ 智能待办 | 待办管理，支持分类、优先级、截止日期，系统托盘显示未完成数量 |
| 🔧 工具盒子 | 单位/进制/时间戳/颜色转换、加密工具、PNG Chunk 注入、批量重命名、NCM 转 MP3/FLAC |
| 🛡️ VirusTotal | 文件病毒扫描，支持批量目录扫描、任务队列管理与检测结果详情 |
| 🌐 网络调试 | TCP/UDP 客户端与服务端调试、FTP 文件传输、DNS 查询与缓存刷新 |
| 📅 日报管理 | 周/月视图日历，支持标签、摘要、周期日报，连续打卡天数统计 |
| 🔀 Git 配置管理 | 多仓库 Git 配置可视化编辑，内置常用键说明与快捷面板 |
| 📄 Makefile 编辑器 | 可视化编辑 Target 与变量，依赖关系图，内置模板库 |
| 🔍 POC 模板解析 | 解析 Nuclei YAML 模板，支持 HTTP/DNS/TCP/Code 协议，YAML 高亮编辑器（50+ 主题） |

## 技术栈

- **后端**：Go · Wails v3 · bbolt
- **前端**：Vue 3 · TypeScript · Ant Design Vue · Pinia · CodeMirror 6

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
