# HappyTools

一个基于 Wails v3 + Vue 3 开发的跨平台桌面应用程序，旨在提升用户的幸福感和工作效率。

## 项目简介

HappyTools 是一款功能丰富的桌面工具应用，集成了待办事项管理、分类管理、日报管理、系统监控、应用设置等多个实用功能。通过现代化的界面设计和流畅的用户体验，帮助用户更好地管理日常任务和生活。

## 技术栈

### 后端 (Go)
- **Wails v3**: 跨平台桌面应用框架 (v3.0.0-alpha.66)
- **Go 1.25.0**: 编程语言
- **bbolt**: 嵌入式键值数据库 (etcd.io/bbolt v1.4.3)
- **gopsutil**: 系统信息监控 (github.com/shirou/gopsutil/v3 v3.24.5)

### 前端 (Vue 3)
- **Vue 3**: 渐进式 JavaScript 框架
- **TypeScript**: 类型安全的 JavaScript 超集
- **Vite**: 现代化的前端构建工具
- **Ant Design Vue**: 企业级 UI 组件库
- **Pinia**: Vue 3 官方状态管理库
- **Vue Router**: 官方路由管理器
- **Tailwind CSS**: 实用优先的 CSS 框架
- **@wailsio/runtime**: Wails v3 运行时库

## 运行原理

本项目基于 Wails v3 框架，采用 Go 后端 + Web 前端的架构模式。参考官方文档：https://v3alpha.wails.io/quick-start/first-app/#how-it-works

### 核心机制

1. **嵌入式资源**: 通过 Go 的 `embed` 包将前端构建产物（`frontend/dist`）嵌入到二进制文件中
2. **服务绑定**: Go 后端服务通过 `application.Service` 暴露给前端调用
3. **双向通信**: 
   - 前端通过自动生成的 TypeScript 绑定调用 Go 后端方法
   - 后端通过 `app.Event.Emit()` 向前端推送实时事件
4. **WebView 渲染**: 使用操作系统原生 WebView 渲染前端界面

### 项目结构

```
happytools/
├── main.go                 # 应用入口点，初始化 Wails 应用
├── go.mod                  # Go 模块依赖
├── go.sum                  # Go 依赖校验
├── Taskfile.yml            # 任务配置文件
├── build/                  # 构建配置
│   ├── config.yml          # Wails v3 配置
│   └── ...                 # 平台特定构建脚本
├── backend/                # Go 后端代码
│   ├── services/           # 业务服务层
│   │   ├── appsettings/    # 应用设置服务
│   │   ├── category/       # 分类管理服务
│   │   ├── dailyreport/    # 日报管理服务
│   │   ├── greetservice/   # 示例服务
│   │   ├── monitor/        # 系统监控服务
│   │   ├── rename/         # 文件重命名服务
│   │   ├── todo/           # 待办事项服务
│   │   └── vt/             # VirusTotal 服务
│   ├── store/              # 数据存储层 (bbolt)
│   │   ├── store.go        # 数据库初始化
│   │   ├── todo_options.go # 待办事项数据操作
│   │   ├── category_options.go # 分类数据操作
│   │   ├── dailyreport_options.go # 日报数据操作
│   │   └── app_settings.go # 应用设置数据操作
│   └── internal/           # 内部工具
├── frontend/               # Vue 3 前端代码
│   ├── src/
│   │   ├── main.ts         # 应用入口
│   │   ├── App.vue         # 根组件
│   │   ├── router/         # 路由配置
│   │   ├── stores/         # Pinia 状态管理
│   │   ├── views/          # 页面组件
│   │   │   ├── Dashboard.vue    # 仪表盘
│   │   │   ├── Home.vue         # 首页
│   │   │   ├── Todo/            # 待办事项页面
│   │   │   ├── Toolbox/         # 工具箱页面
│   │   │   └── VirusTotal/      # VirusTotal 页面
│   │   ├── components/     # 通用组件
│   │   ├── utils/          # 工具函数
│   │   ├── config/         # 配置文件
│   │   └── assets/         # 静态资源
│   ├── bindings/           # Wails 自动生成的 TypeScript 绑定
│   │   └── github.com/Aliuyanfeng/happytools/backend/services/
│   ├── dist/               # 构建输出目录
│   ├── package.json        # Node.js 依赖
│   ├── vite.config.ts      # Vite 配置
│   └── tailwind.config.js  # Tailwind CSS 配置
├── bin/                    # 编译输出目录
├── docs/                   # 文档目录
└── example/                # 示例代码
```

## 开发流程

### 环境要求

- Go 1.25.0 或更高版本
- Node.js 18+ (推荐使用 LTS 版本)
- Wails v3 CLI: `go install github.com/wailsapp/wails/v3/cmd/wails3@latest`

### 安装依赖

#### 后端依赖
```bash
go mod download
```

#### 前端依赖
```bash
cd frontend
npm install
cd ..
```

### 开发模式

使用 Task 运行开发服务器（支持热重载）：

```bash
wails3 dev -config ./build/config.yml -port 9250
```

或者使用 Wails Task 命令：

```bash
wails3 task dev
```

开发模式会自动：
1. 安装前端依赖
2. 启动 Vite 开发服务器（端口 9250）
3. 编译 Go 后端
4. 启动应用并监听文件变化

### 构建生产版本

#### 构建应用
```bash
wails3 build
```

或使用 Wails Task 命令：
```bash
wails3 task build
```

或使用平台特定命令：
```bash
# Windows
wails3 task windows:build

# macOS
wails3 task darwin:build

# Linux
wails3 task linux:build
```

#### 打包应用
```bash
wails3 package
```

或使用 Wails Task 命令：
```bash
wails3 task package
```

构建产物将输出到 `bin/` 目录。

## 功能模块

### 1. 待办事项管理 (Todo)
- 创建、编辑、删除待办事项
- 支持分类和优先级
- 支持截止日期
- 标记完成状态

### 2. 分类管理 (Category)
- 创建和管理待办事项分类
- 分类统计和查询

### 3. 日报管理 (DailyReport)
- 记录日常工作日报
- 支持富文本内容
- 日期归档

### 4. 系统监控 (Monitor)
- 实时监控系统资源
- CPU 使用率
- 内存使用情况
- 磁盘空间

### 5. 应用设置 (AppSettings)
- 记录上次使用时间
- 应用偏好设置

### 6. 文件重命名 (Rename)
- 批量文件重命名工具

### 7. VirusTotal 集成
- 文件病毒扫描
- 威胁情报查询

## 前后端通信

### 服务绑定

Go 服务通过 `application.NewService()` 绑定到 Wails 应用：

```go
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(&greetservice.GreetService{}),
        application.NewService(monitor.NewSysInfoService()),
        application.NewService(todo.NewTodoService()),
        // ... 更多服务
    },
})
```

### 前端调用

Wails 会自动生成 TypeScript 绑定文件到 `frontend/bindings/` 目录。前端可以直接导入并调用：

```typescript
import { TodoService } from '../../bindings/github.com/Aliuyanfeng/happytools/backend/services/todo'

// 调用后端方法
const todos = await TodoService.GetAll()
await TodoService.Add("新任务", null, null, 1)
```

### 事件通信

后端通过事件向前端推送数据：

```go
// Go 后端发送事件
app.Event.Emit("time", time.Now().Format(time.RFC1123))
app.Event.Emit("monitor:sysInfo", sysInfo)
```

前端监听事件：

```typescript
import { Events } from "@wailsio/runtime"

Events.On('time', (event) => {
  console.log(event.data)
})

Events.On('monitor:sysInfo', (event) => {
  // 处理系统信息
})
```

## 数据存储

项目使用 **bbolt** 嵌入式数据库，数据存储在用户主目录下的 `.happytools/data.db` 文件中。

数据库结构：
- `todos`: 待办事项数据
- `categories`: 分类数据
- `daily_reports`: 日报数据
- `app_settings`: 应用设置数据

## 系统托盘

应用支持系统托盘功能：
- 托盘图标显示应用名称
- 右键菜单提供 "Open" 和 "Quit" 选项
- 点击托盘图标可显示/隐藏窗口

## 窗口配置

主窗口配置（main.go:127）：
- 固定大小：1024 x 768
- 无边框窗口 (Frameless: true)
- 背景色：#1B2636
- 支持置顶功能

## 开发指南

### 添加新服务

1. 在 `backend/services/` 创建新的服务包
2. 实现服务结构体和方法
3. 在 `main.go` 中注册服务：

```go
application.NewService(&mynewservice.MyNewService{})
```

4. 运行开发服务器，Wails 会自动生成 TypeScript 绑定
5. 在前端导入并使用生成的绑定

### 添加新页面

1. 在 `frontend/src/views/` 创建新的 Vue 组件
2. 在 `frontend/src/router/router.ts` 添加路由配置
3. 在导航菜单中添加入口

### 数据库操作

参考 `backend/store/` 中的实现，使用 bbolt 进行数据存储：

```go
DB.Update(func(tx *bbolt.Tx) error {
    bucket := tx.Bucket(bucketName)
    // 数据操作
    return nil
})
```

## 常见问题

### Q: 如何修改应用名称和图标？
A: 编辑 `build/config.yml` 中的 `info` 部分，然后运行 `wails3 task common:update:build-assets`

### Q: 前端绑定没有生成？
A: 确保服务已正确注册，重启开发服务器即可自动生成

### Q: 数据库文件在哪里？
A: 数据库文件位于 `~/.happytools/data.db` (Windows: `C:\Users\<用户名>\.happytools\data.db`)

### Q: 如何调试前端？
A: 开发模式下，前端运行在 http://localhost:9250，可以使用浏览器开发者工具

## 许可证

本项目采用开源许可证，详见 [LICENSE](LICENSE) 文件。

## 贡献

欢迎提交 Issue 和 Pull Request！

## 联系方式

如有问题或建议，请通过 GitHub Issues 联系。

---

**HappyTools** - 让工作更高效，让生活更快乐！

