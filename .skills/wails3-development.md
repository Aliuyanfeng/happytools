# Wails3 开发完整指南 Skill

## Skill 描述
这个 skill 提供了 HappyTools 项目（基于 Wails v3 + Vue 3）的完整开发指南，包括技术栈规范、项目结构、开发流程、最佳实践和常见问题解决方案。

## 适用场景
- 开始 Wails v3 项目开发
- 了解项目架构和技术栈
- 学习前后端通信机制
- 掌握数据存储方案
- 解决开发过程中的问题
- 遵循项目开发规范

## 技术栈概览

### 后端技术栈
- **框架**: Wails v3 (v3.0.0-alpha.66)
- **语言**: Go 1.25.0+
- **数据库**: bbolt (嵌入式键值数据库 v1.4.3)
- **系统监控**: gopsutil (v3.24.5)

### 前端技术栈
- **框架**: Vue 3.2.45+
- **语言**: TypeScript 4.9.3+
- **构建工具**: Vite 5.0.0+
- **UI 组件库**: Ant Design Vue 4.2.6+
- **状态管理**: Pinia 3.0.3+
- **路由**: Vue Router 4.5.1+
- **CSS 框架**: Tailwind CSS 3.4.1+
- **Wails 运行时**: @wailsio/runtime 3.0.0-alpha.79+

## 项目结构

```
happytools/
├── main.go                 # 应用入口点
├── go.mod / go.sum         # Go 模块依赖
├── Taskfile.yml            # 任务配置文件
├── build/                  # 构建配置
│   └── config.yml          # Wails v3 配置
├── backend/                # Go 后端代码
│   ├── services/           # 业务服务层
│   │   ├── appsettings/    # 应用设置服务
│   │   ├── category/       # 分类管理服务
│   │   ├── clipboard/      # 剪贴板服务
│   │   ├── dailyreport/    # 日报管理服务
│   │   ├── encryption/     # 加密工具服务
│   │   ├── greetservice/   # 示例服务
│   │   ├── monitor/        # 系统监控服务
│   │   ├── network/        # 网络服务
│   │   ├── rename/         # 文件重命名服务
│   │   ├── todo/           # 待办事项服务
│   │   ├── unitconverter/  # 单位转换服务
│   │   └── vt/             # VirusTotal 服务
│   ├── store/              # 数据存储层
│   │   ├── store.go        # 数据库初始化
│   │   ├── todo_options.go # 待办事项数据操作
│   │   ├── category_options.go # 分类数据操作
│   │   ├── dailyreport_options.go # 日报数据操作
│   │   ├── app_settings.go # 应用设置数据操作
│   │   └── vt_options.go   # VirusTotal 数据操作
│   └── internal/           # 内部工具
│       └── sqlite/         # SQLite 工具
├── frontend/               # Vue 3 前端代码
│   ├── src/
│   │   ├── main.ts         # 应用入口
│   │   ├── App.vue         # 根组件
│   │   ├── router/         # 路由配置
│   │   ├── stores/         # Pinia 状态管理
│   │   ├── views/          # 页面组件
│   │   ├── components/     # 通用组件
│   │   ├── utils/          # 工具函数
│   │   ├── config/         # 配置文件
│   │   ├── locales/        # 国际化
│   │   └── assets/         # 静态资源
│   ├── bindings/           # Wails 自动生成的 TypeScript 绑定
│   └── package.json        # Node.js 依赖
└── bin/                    # 编译输出目录
```

## 核心机制

### 1. 嵌入式资源
通过 Go 的 `embed` 包将前端构建产物嵌入到二进制文件中：

```go
//go:embed all:frontend/dist
var assets embed.FS
```

### 2. 服务绑定
Go 后端服务通过 `application.Service` 暴露给前端调用：

```go
app.RegisterService(application.NewService(todo.NewTodoService()))
```

### 3. 双向通信
- **前端调用后端**: 通过自动生成的 TypeScript 绑定
- **后端推送前端**: 通过 `app.Event.Emit()` 发送事件

### 4. WebView 渲染
使用操作系统原生 WebView 渲染前端界面

## 后端开发规范

### 1. 服务层开发

#### 服务结构定义
```go
// 服务结构体命名规范: [功能名]Service
type TodoService struct{}

// 构造函数命名规范: New[功能名]Service
func NewTodoService() *TodoService {
    return &TodoService{}
}
```

#### 服务方法规范
```go
// 方法必须以大写字母开头（公开方法）
// 返回值应包含 error 类型作为最后一个返回值
// 使用指针接收器定义方法

// 示例：获取所有数据
func (t *TodoService) GetAll() ([]Todo, error) {
    todos, err := store.GetAllTodos()
    if err != nil {
        return nil, err
    }
    // 数据转换逻辑
    return result, nil
}

// 示例：添加数据
func (t *TodoService) Add(title string, categoryID *int, dueDate *string, priority int) (*Todo, error) {
    // 参数验证
    if title == "" {
        return nil, errors.New("title cannot be empty")
    }

    // 业务逻辑
    todo, err := store.CreateTodoEnhanced(title, categoryID, parsedDue, priority)
    if err != nil {
        return nil, err
    }

    return &Todo{
        ID:         todo.ID,
        Title:      todo.Title,
        Completed:  todo.Completed,
        CategoryID: todo.CategoryID,
        DueDate:    dueDateStr,
        Priority:   todo.Priority,
        Status:     todo.GetStatus(),
    }, nil
}
```

#### 数据模型规范
```go
// 数据模型命名规范: [功能名]
type Todo struct {
    ID         int     `json:"id"`
    Title      string  `json:"title"`
    Completed  bool    `json:"completed"`
    CategoryID *int    `json:"category_id"`    // 可空字段使用指针
    DueDate    *string `json:"due_date"`       // 可空字段使用指针
    Priority   int     `json:"priority"`
    Status     int     `json:"status"`
}
```

### 2. 数据存储层开发

#### 数据库初始化
```go
// 数据库路径: ~/.happytools/data.db
homeDir, err := os.UserHomeDir()
if err != nil {
    log.Fatal("获取用户目录失败:", err)
}
dbPath := filepath.Join(homeDir, ".happytools", "data.db")
if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
    log.Fatal("创建数据目录失败:", err)
}
if err := store.Init(dbPath); err != nil {
    log.Fatal("初始化数据库失败:", err)
}
```

#### Bucket 命名规范
```go
var (
    todoBucket        = []byte("todos")
    categoryBucket    = []byte("categories")
    dailyReportBucket = []byte("daily_reports")
    appSettingsBucket = []byte("app_settings")
    vtTaskBucket      = []byte("vt_tasks")
    vtFileBucket      = []byte("vt_files")
)
```

#### 数据库操作规范
```go
// 使用 DB.Update 进行写操作
DB.Update(func(tx *bbolt.Tx) error {
    bucket := tx.Bucket(bucketName)
    if bucket == nil {
        return errors.New("bucket not found")
    }
    // 数据操作
    return nil
})

// 使用 DB.View 进行读操作
DB.View(func(tx *bbolt.Tx) error {
    bucket := tx.Bucket(bucketName)
    if bucket == nil {
        return nil
    }
    // 数据读取
    return nil
})
```

### 3. 服务注册

在 `main.go` 中注册服务：

```go
import (
    "github.com/Aliuyanfeng/happytools/backend/services/todo"
    "github.com/Aliuyanfeng/happytools/backend/services/category"
    // ... 其他服务导入
)

func main() {
    app := application.New(application.Options{
        Name:        "happytools",
        Description: "A tool that can enhance one's sense of happiness",
        Assets: application.AssetOptions{
            Handler: application.AssetFileServerFS(assets),
        },
    })

    // 注册所有服务
    app.RegisterService(application.NewService(&greetservice.GreetService{}))
    app.RegisterService(application.NewService(monitor.NewSysInfoService()))
    app.RegisterService(application.NewService(todo.NewTodoService()))
    app.RegisterService(application.NewService(category.NewCategoryService()))
    app.RegisterService(application.NewService(dailyreport.NewDailyReportService()))
    app.RegisterService(application.NewService(appsettings.NewAppSettingsService()))
    app.RegisterService(application.NewService(rename.NewRenameService(app)))
    app.RegisterService(application.NewService(virusTotal.NewVTService(app)))
    app.RegisterService(application.NewService(network.NewFileTransferService(app)))
    app.RegisterService(application.NewService(network.NewTCPUDPService(app)))
    app.RegisterService(application.NewService(network.NewDNSService()))
    app.RegisterService(application.NewService(unitconverter.NewUnitConverterService()))
    app.RegisterService(application.NewService(encryption.NewEncryptionService()))
    app.RegisterService(application.NewService(clipboard.NewClipboardService(app)))

    // ... 其他代码
}
```

### 4. 事件通信

#### 后端发送事件
```go
// 发送简单事件
app.Event.Emit("time", time.Now().Format(time.RFC1123))

// 发送结构化数据
sysInfo, err := monitor.NewSysInfoService().GetSysInfo()
if err != nil {
    log.Printf("获取系统信息失败: %v", err)
}
app.Event.Emit("monitor:sysInfo", sysInfo)

// 发送应用事件
app.Event.Emit("app:lastUsedTime", lastUsedTime)
```

#### 事件命名规范
- 使用小写字母和冒号分隔
- 格式: `[模块名]:[事件名]`
- 示例: `monitor:sysInfo`, `app:lastUsedTime`

### 5. 错误处理

```go
// 必须处理所有错误
if err != nil {
    log.Printf("操作描述: %v", err)
    return err
}

// 使用 errors.New 创建错误
if title == "" {
    return errors.New("title cannot be empty")
}

// 参数验证
if dueDate != nil && *dueDate != "" {
    d, err := time.Parse("2006-01-02", *dueDate)
    if err != nil {
        return err
    }
    parsedDue = &d
}
```

## 前端开发规范

### 1. 组件开发

#### 组件文件命名
- 使用 PascalCase 命名: `Todo.vue`, `Dashboard.vue`
- 页面组件放在 `views/` 目录
- 通用组件放在 `components/` 目录

#### 组件结构
```vue
<template>
  <div class="p-6">
    <h1 class="text-2xl font-bold text-white mb-4">{{ meta.title }}</h1>
    <!-- 组件内容 -->
  </div>
</template>

<script setup lang="ts">
// 导入依赖
import { ref, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { TodoService } from '../../bindings/github.com/Aliuyanfeng/happytools/backend/services/todo'

// 组件逻辑
const items = ref([])

const loadData = async () => {
  try {
    const result = await TodoService.GetAll()
    items.value = result || []
  } catch (error) {
    message.error('加载数据失败')
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
/* 组件样式 */
</style>
```

### 2. 路由配置

#### 路由定义
```typescript
// frontend/src/router/routes.ts
import { RouteRecordRaw } from "vue-router"

const Home = () => import('../views/Home.vue')
const Dashboard = () => import('../views/Dashboard.vue')
const Todo = () => import('../views/Todo/Todo.vue')

const routes: Array<RouteRecordRaw> = [
    {
        path: '/',
        name: 'home',
        meta: {
            title: '首页'
        },
        component: Home
    },
    {
        path: '/dashboard',
        name: 'dashboard',
        meta: {
            title: '仪表盘'
        },
        component: Dashboard
    },
    {
        path: '/todo',
        name: 'todo',
        meta: {
            title: 'TODO'
        },
        component: Todo
    }
]

export default routes
```

#### 路由命名规范
- path: 使用小写字母和连字符
- name: 使用小写字母
- meta.title: 使用中文标题

### 3. 状态管理

#### Pinia Store 定义
```typescript
// stores/todo.ts
import { defineStore } from 'pinia'

export const useTodoStore = defineStore('todo', {
  state: () => ({
    todos: [],
    loading: false,
    error: null
  }),

  actions: {
    async fetchTodos() {
      this.loading = true
      try {
        const result = await TodoService.GetAll()
        this.todos = result || []
      } catch (error) {
        this.error = error
      } finally {
        this.loading = false
      }
    },

    async addTodo(title: string, categoryID?: number, dueDate?: string, priority: number = 1) {
      try {
        await TodoService.Add(title, categoryID ? categoryID : null, dueDate ? dueDate : null, priority)
        await this.fetchTodos()
      } catch (error) {
        throw error
      }
    }
  }
})
```

### 4. Wails 绑定调用

#### 导入绑定
```typescript
import { TodoService } from '../../bindings/github.com/Aliuyanfeng/happytools/backend/services/todo'
import { CategoryService } from '../../bindings/github.com/Aliuyanfeng/happytools/backend/services/category'
```

#### 调用后端方法
```typescript
// 获取数据
const todos = await TodoService.GetAll()

// 添加数据
await TodoService.Add("新任务", null, null, 1)

// 更新数据
await TodoService.Update(id, title, completed, categoryID, dueDate, priority)

// 删除数据
await TodoService.Delete(id)

// 切换完成状态
await TodoService.Toggle(id)
```

### 5. 事件监听

#### 监听后端事件
```typescript
import { Events } from "@wailsio/runtime"
import { onMounted, onUnmounted } from 'vue'

// 在组件中监听事件
onMounted(() => {
  // 监听时间事件
  Events.On('time', (event) => {
    console.log('Current time:', event.data)
  })

  // 监听系统信息事件
  Events.On('monitor:sysInfo', (event) => {
    const sysInfo = event.data
    console.log('System info:', sysInfo)
  })

  // 监听应用事件
  Events.On('app:lastUsedTime', (event) => {
    const lastUsedTime = event.data
    console.log('Last used time:', lastUsedTime)
  })
})

// 清理事件监听
onUnmounted(() => {
  Events.Off('time')
  Events.Off('monitor:sysInfo')
  Events.Off('app:lastUsedTime')
})
```

### 6. UI 组件使用

#### Ant Design Vue 组件
```vue
<template>
  <a-form :model="form" @finish="handleSubmit">
    <a-form-item label="标题" name="title" :rules="[{ required: true, message: '请输入标题' }]">
      <a-input v-model:value="form.title" placeholder="请输入标题" />
    </a-form-item>

    <a-form-item label="分类" name="categoryID">
      <a-select v-model:value="form.categoryID" placeholder="选择分类" allowClear>
        <a-select-option v-for="cat in categories" :key="cat.id" :value="cat.id">
          {{ cat.name }}
        </a-select-option>
      </a-select>
    </a-form-item>

    <a-form-item label="截止日期" name="dueDate">
      <a-date-picker v-model:value="form.dueDate" format="YYYY-MM-DD" />
    </a-form-item>

    <a-form-item label="优先级" name="priority">
      <a-rate v-model:value="form.priority" :count="3" />
    </a-form-item>

    <a-form-item>
      <a-button type="primary" html-type="submit">提交</a-button>
      <a-button style="margin-left: 10px" @click="handleCancel">取消</a-button>
    </a-form-item>
  </a-form>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { message } from 'ant-design-vue'

const form = ref({
  title: '',
  categoryID: undefined,
  dueDate: null,
  priority: 1
})

const handleSubmit = async () => {
  try {
    await TodoService.Add(
      form.value.title,
      form.value.categoryID || null,
      form.value.dueDate ? form.value.dueDate.format('YYYY-MM-DD') : null,
      form.value.priority
    )
    message.success('添加成功')
  } catch (error) {
    message.error('添加失败')
  }
}
</script>
```

## 构建和部署

### 1. 开发环境

#### 启动开发服务器
```bash
# 方式1: 使用配置文件
wails3 dev -config ./build/config.yml -port 9250

# 方式2: 使用 Task 命令
wails3 task dev
```

#### 开发模式特性
- 自动安装前端依赖
- 启动 Vite 开发服务器（端口 9250）
- 编译 Go 后端
- 监听文件变化并热重载
- 自动生成 TypeScript 绑定

### 2. 生产构建

#### 构建应用
```bash
# 通用构建
wails3 build

# 使用 Task 命令
wails3 task build

# 平台特定构建
wails3 task windows:build
wails3 task darwin:build
wails3 task linux:build
```

#### 打包应用
```bash
wails3 package
# 或
wails3 task package
```

### 3. 构建配置

#### config.yml 配置
```yaml
version: '3'
info:
  companyName: "AliuStudio"
  productName: "HappyTools"
  productIdentifier: "com.aliustudio.happytools"
  description: "A tool that can enhance one's sense of happiness"
  copyright: "(c) 2025, AliuStudio"
  comments: "HappyTools - 提升幸福感的工具集"
  version: "1.0.0"

dev_mode:
  root_path: .
  log_level: warn
  debounce: 1000
  ignore:
    dir:
      - .git
      - node_modules
      - frontend
      - bin
    file:
      - .DS_Store
      - .gitignore
      - .gitkeep
    watched_extension:
      - "*.go"
    git_ignore: true
  executes:
    - cmd: wails3 task common:install:frontend:deps
      type: once
    - cmd: wails3 task common:dev:frontend
      type: background
    - cmd: go mod tidy
      type: blocking
    - cmd: wails3 task build
      type: blocking
    - cmd: wails3 task run
      type: primary
```

## 窗口配置

### 主窗口配置
```go
mainWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:     "happytools",
    Width:     1024,
    Height:    768,
    MinWidth:  1024,
    MinHeight: 768,
    MaxWidth:  1024,
    MaxHeight: 768,
    Mac: application.MacWindow{
        InvisibleTitleBarHeight: 50,
        Backdrop:                application.MacBackdropTranslucent,
        TitleBar:                application.MacTitleBarHiddenInset,
    },
    Windows: application.WindowsWindow{
        DisableIcon: false,
    },
    BackgroundColour: application.NewRGB(27, 38, 54),
    URL:              "/",
    Frameless:        true,
})

mainWindow.SetAlwaysOnTop(false)
```

### 系统托盘配置
```go
systray := app.SystemTray.New()
systray.SetLabel("HappyTools")
systray.SetTooltip("HappyTools工具")

menu := application.NewMenu()
menu.Add("退出 HappyTools").OnClick(func(ctx *application.Context) {
    app.Quit()
})

systray.SetMenu(menu)
```

## 添加新功能流程

### 1. 添加新服务

1. 在 `backend/services/` 创建新的服务包
2. 实现服务结构体和方法
3. 在 `main.go` 中注册服务：
   ```go
   app.RegisterService(application.NewService(&mynewservice.MyNewService{}))
   ```
4. 运行开发服务器，Wails 会自动生成 TypeScript 绑定
5. 在前端导入并使用生成的绑定

### 2. 添加新页面

1. 在 `frontend/src/views/` 创建新的 Vue 组件
2. 在 `frontend/src/router/routes.ts` 添加路由配置：
   ```typescript
   {
       path: '/newpage',
       name: 'newpage',
       meta: { title: '新页面' },
       component: () => import('../views/NewPage.vue')
   }
   ```
3. 在导航菜单中添加入口

### 3. 添加新数据模型

1. 在 `backend/store/` 中定义数据操作函数
2. 在 `store.go` 中添加新的 bucket
3. 在服务层中使用数据操作函数

## 代码风格规范

### Go 代码风格
- 使用 gofmt 格式化代码
- 遵循 Go 官方代码规范
- 注释使用中文或英文，保持一致性
- 文件头部添加作者和日期注释

### TypeScript/Vue 代码风格
- 使用 ESLint 和 Prettier 格式化代码
- 使用 TypeScript 类型注解
- 组件使用 `<script setup>` 语法
- 样式使用 `scoped` 限定作用域

## 调试规范

### 前端调试
- 开发模式下前端运行在 http://localhost:9250
- 使用浏览器开发者工具调试
- 使用 console.log 输出调试信息

### 后端调试
- 使用 log.Printf 输出调试信息
- 检查错误日志
- 使用 Go 的调试工具

## 常见问题处理

### Q: 前端绑定没有生成？
A: 确保服务已正确注册，重启开发服务器即可自动生成

### Q: 数据库文件在哪里？
A: 数据库文件位于 `~/.happytools/data.db` (Windows: `C:\Users\<用户名>\.happytools\data.db`)

### Q: 如何修改应用名称和图标？
A: 编辑 `build/config.yml` 中的 `info` 部分，然后运行 `wails3 task common:update:build-assets`

### Q: 如何调试前端？
A: 开发模式下，前端运行在 http://localhost:9250，可以使用浏览器开发者工具

### Q: 如何处理跨平台问题？
A: 使用 Wails 提供的平台特定配置，在 `build/` 目录下有各平台的构建脚本

## 最佳实践

1. **服务单一职责**: 每个服务只负责一个功能模块
2. **错误处理**: 所有错误必须处理并记录日志
3. **数据验证**: 在服务层进行参数验证
4. **事件命名**: 使用统一的命名规范
5. **代码复用**: 提取公共逻辑到工具函数
6. **类型安全**: 使用 TypeScript 类型注解
7. **组件化**: 将复杂页面拆分为多个组件
8. **状态管理**: 使用 Pinia 管理全局状态
9. **性能优化**: 使用懒加载和代码分割
10. **安全性**: 验证所有用户输入，处理敏感数据

## 版本控制规范

- 使用 Git 进行版本控制
- 提交信息使用中文或英文，保持一致性
- 每个功能模块独立提交
- 重要修改需要添加详细说明

## 总结

这个 skill 提供了 HappyTools 项目（基于 Wails v3 + Vue 3）的完整开发指南，涵盖了从项目结构、技术栈、开发规范到最佳实践的所有内容。通过遵循这个 skill，开发者可以快速上手项目开发，确保代码质量和一致性。
