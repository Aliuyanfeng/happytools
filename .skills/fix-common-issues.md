# 修复常见问题 Skill

## Skill 描述
这个 skill 用于快速诊断和修复 HappyTools 项目中的常见问题，包括编译错误、运行时错误、前后端通信问题等。

## 适用场景
- 前端绑定未生成
- 后端服务无法调用
- 数据库操作失败
- 事件通信异常
- 构建打包问题

## 常见问题及解决方案

### 1. 前端绑定未生成

#### 问题现象
- TypeScript 绑定文件未生成到 `frontend/bindings/` 目录
- 前端无法导入后端服务
- 类型检查失败

#### 诊断步骤
```bash
# 1. 检查服务是否正确注册
grep "RegisterService" main.go

# 2. 检查服务方法是否公开（首字母大写）
grep -r "func.*Service.*GetAll" backend/services/

# 3. 检查开发服务器是否正常运行
wails3 dev -config ./build/config.yml -port 9250
```

#### 解决方案
1. **确保服务已注册**：
   ```go
   // main.go
   app.RegisterService(application.NewService(todo.NewTodoService()))
   ```

2. **确保方法公开**：
   ```go
   // 正确：首字母大写
   func (t *TodoService) GetAll() ([]Todo, error) { }
   
   // 错误：首字母小写
   func (t *TodoService) getAll() ([]Todo, error) { }
   ```

3. **重启开发服务器**：
   ```bash
   # 停止当前服务器
   Ctrl+C
   
   # 重新启动
   wails3 dev -config ./build/config.yml -port 9250
   ```

### 2. 数据库操作失败

#### 问题现象
- 数据无法保存或读取
- 报错 "bucket not found"
- 数据库文件损坏

#### 诊断步骤
```bash
# 1. 检查数据库文件是否存在
ls -la ~/.happytools/data.db

# 2. 检查数据库路径
echo $HOME
# Windows: echo %USERPROFILE%

# 3. 检查 bucket 是否初始化
grep "CreateBucketIfNotExists" backend/store/store.go
```

#### 解决方案
1. **检查数据库初始化**：
   ```go
   // backend/store/store.go
   func InitDB() error {
       // 确保所有 bucket 都被创建
       err = db.Update(func(tx *bbolt.Tx) error {
           _, err := tx.CreateBucketIfNotExists(todoBucket)
           if err != nil {
               return err
           }
           // ... 其他 bucket
           return nil
       })
       return err
   }
   ```

2. **检查数据库路径**：
   ```go
   homeDir, err := os.UserHomeDir()
   if err != nil {
       return err
   }
   dbPath := filepath.Join(homeDir, ".happytools", "data.db")
   ```

3. **重建数据库**（谨慎操作）：
   ```bash
   # 备份现有数据
   cp ~/.happytools/data.db ~/.happytools/data.db.backup
   
   # 删除数据库
   rm ~/.happytools/data.db
   
   # 重启应用，会自动创建新数据库
   ```

### 3. 前后端通信失败

#### 问题现象
- 前端调用后端方法无响应
- 报错 "service not found"
- 数据传输异常

#### 诊断步骤
```typescript
// 前端添加调试日志
import { TodoService } from '../../bindings/github.com/Aliuyanfeng/happytools/backend/services/todo'

console.log('TodoService:', TodoService)
console.log('Available methods:', Object.keys(TodoService))

try {
  const result = await TodoService.GetAll()
  console.log('Result:', result)
} catch (error) {
  console.error('Error:', error)
}
```

#### 解决方案
1. **检查服务注册**：
   ```go
   // main.go
   // 确保服务在应用初始化时注册
   app := application.New(application.Options{
       Services: []application.Service{
           application.NewService(&todo.TodoService{}),
       },
   })
   ```

2. **检查绑定路径**：
   ```typescript
   // 确保导入路径正确
   import { TodoService } from '../../bindings/github.com/Aliuyanfeng/happytools/backend/services/todo'
   ```

3. **检查方法签名**：
   ```go
   // 后端方法签名必须匹配前端调用
   func (t *TodoService) GetAll() ([]Todo, error) { }
   ```

### 4. 事件通信异常

#### 问题现象
- 前端无法接收后端事件
- 事件数据格式错误
- 事件监听器未触发

#### 诊断步骤
```go
// 后端添加调试日志
log.Printf("Emitting event: %s, data: %v", "monitor:sysInfo", sysInfo)
app.Event.Emit("monitor:sysInfo", sysInfo)
```

```typescript
// 前端添加调试日志
import { Events } from "@wailsio/runtime"

Events.On('monitor:sysInfo', (event) => {
    console.log('Event received:', event)
    console.log('Event data:', event.data)
})
```

#### 解决方案
1. **检查事件名称**：
   ```go
   // 使用统一的事件命名规范
   app.Event.Emit("monitor:sysInfo", sysInfo)  // 正确
   app.Event.Emit("Monitor:SysInfo", sysInfo)  // 错误：大小写不一致
   ```

2. **检查事件监听**：
   ```typescript
   // 确保事件名称完全匹配
   Events.On('monitor:sysInfo', (event) => { })  // 正确
   Events.On('Monitor:SysInfo', (event) => { })  // 错误
   ```

3. **检查数据序列化**：
   ```go
   // 确保数据可以被 JSON 序列化
   type SysInfo struct {
       CPUUsage    float64 `json:"cpu_usage"`
       MemoryUsage float64 `json:"memory_usage"`
   }
   ```

### 5. 构建打包失败

#### 问题现象
- 编译错误
- 打包失败
- 运行时缺少依赖

#### 诊断步骤
```bash
# 1. 检查 Go 版本
go version

# 2. 检查依赖
go mod tidy
go mod verify

# 3. 检查前端依赖
cd frontend
npm install
npm run build
cd ..

# 4. 尝试构建
wails3 build
```

#### 解决方案
1. **更新依赖**：
   ```bash
   # 更新 Go 依赖
   go get -u
   go mod tidy
   
   # 更新前端依赖
   cd frontend
   npm update
   cd ..
   ```

2. **清理构建缓存**：
   ```bash
   # 清理 Go 缓存
   go clean -cache
   
   # 清理前端构建产物
   rm -rf frontend/dist
   rm -rf bin
   
   # 重新构建
   wails3 build
   ```

3. **检查平台特定问题**：
   ```bash
   # Windows: 检查 WebView2 是否安装
   # macOS: 检查 Xcode 命令行工具
   xcode-select --install
   
   # Linux: 检查依赖
   sudo apt-get install libgtk-3-dev libwebkit2gtk-4.0-dev
   ```

### 6. 窗口显示异常

#### 问题现象
- 窗口无法显示
- 窗口大小异常
- 窗口样式错误

#### 诊断步骤
```go
// 检查窗口配置
mainWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:     "happytools",
    Width:     1024,
    Height:    768,
    MinWidth:  1024,
    MinHeight: 768,
    MaxWidth:  1024,
    MaxHeight: 768,
    BackgroundColour: application.NewRGB(27, 38, 54),
    URL:              "/",
    Frameless:        true,
})
```

#### 解决方案
1. **检查 URL 配置**：
   ```go
   // 开发模式
   URL: "http://localhost:9250"
   
   // 生产模式
   URL: "/"
   ```

2. **检查窗口显示**：
   ```go
   // 确保调用了 Show()
   mainWindow.Show()
   ```

3. **检查无边框窗口**：
   ```go
   // 如果使用无边框窗口，需要自定义标题栏
   Frameless: true,
   ```

## 调试技巧

### 1. 后端调试
```go
// 使用 log.Printf 输出调试信息
log.Printf("Debug: %+v", data)

// 检查错误
if err != nil {
    log.Printf("Error: %v", err)
    return err
}
```

### 2. 前端调试
```typescript
// 使用 console.log 调试
console.log('Debug:', data)

// 使用浏览器开发者工具
// F12 打开开发者工具
// 查看 Console 和 Network 标签
```

### 3. 数据库调试
```go
// 添加数据库操作日志
DB.Update(func(tx *bbolt.Tx) error {
    log.Printf("Updating bucket: %s", bucketName)
    bucket := tx.Bucket(bucketName)
    // ... 操作
    return nil
})
```

## 预防措施

1. **代码规范**：
   - 遵循项目开发规范
   - 使用一致的命名规范
   - 添加必要的错误处理

2. **测试验证**：
   - 开发过程中频繁测试
   - 使用版本控制
   - 定期提交代码

3. **文档维护**：
   - 更新 README.md
   - 记录重要变更
   - 维护 CHANGELOG.md

## 总结

这个 skill 提供了 HappyTools 项目常见问题的诊断和解决方案，帮助开发者快速定位和修复问题，提高开发效率。
