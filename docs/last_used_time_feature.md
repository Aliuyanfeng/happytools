# 上次使用时间功能说明

## 功能概述

本功能实现了在 HappyTools 应用右下角显示真实的"上次使用时间"，替代了之前的硬编码假数据。

## 实现细节

### 后端实现

#### 1. 数据存储层 (backend/store/app_settings.go)

新增了应用设置相关的存储功能：

- `GetLastUsedTime()`: 获取上次使用时间
- `UpdateLastUsedTime()`: 更新上次使用时间为当前时间

数据存储在 `app_settings` bucket 中，使用 `last_used_time` 作为 key。

#### 2. 服务层 (backend/services/appsettings/appsettings.go)

创建了 `AppSettingsService` 服务，提供以下方法：

- `GetLastUsedTime()`: 获取格式化后的上次使用时间字符串
- `UpdateLastUsedTime()`: 更新上次使用时间
- `GetLastUsedTimestamp()`: 获取上次使用时间戳

#### 3. 应用集成 (main.go)

在 `main.go` 中：

1. 注册了 `AppSettingsService` 服务
2. **应用启动时先读取数据库中的时间（上一次启动时间）**
3. **通过 `app.Event.Emit("app:lastUsedTime", time)` 事件发送给前端**
4. **然后更新数据库为当前时间（为下一次启动做准备）**

**关键**: 必须先读取再更新，这样才能显示上一次的时间而不是当前时间。

### 前端实现

#### 1. 状态管理 (frontend/src/stores/app.ts)

创建了 `useAppStore` Pinia store：

- `lastUsedTime`: 存储上次使用时间
- `setupEventListener()`: 设置事件监听器，监听后端发送的 `app:lastUsedTime` 事件

**重要**: 前端不再主动调用后端服务，而是通过事件监听被动接收数据。

#### 2. UI 集成 (frontend/src/App.vue)

在 App.vue 中：

1. 导入并使用 `useAppStore`
2. 在底部状态栏显示真实的上次使用时间
3. 通过 Pinia store 自动更新显示

## 数据流程

```
应用启动
    ↓
main.go 先读取数据库中的时间（上一次启动时间）
    ↓
通过 Events.Emit("app:lastUsedTime") 发送给前端
    ↓
前端监听事件并更新显示
    ↓
然后 main.go 调用 UpdateLastUsedTime()
    ↓
存储当前时间到数据库 (~/.happytools/data.db)
    ↓
为下一次启动做准备
```

## 技术栈

- **后端**: Go + bbolt (嵌入式数据库)
- **前端**: Vue 3 + Pinia + TypeScript
- **通信**: Wails3 事件系统 (Event System)

## 数据存储

- **数据库**: bbolt (嵌入式键值数据库)
- **存储位置**: `~/.happytools/data.db`
- **Bucket**: `app_settings`
- **Key**: `last_used_time`
- **格式**: RFC3339 (ISO 8601)

## 测试

已通过以下测试：

1. ✓ 数据库初始化
2. ✓ 首次读取（无数据）
3. ✓ 更新上次使用时间
4. ✓ 再次读取（有数据）
5. ✓ 时间更新验证

## 使用说明

### 首次使用

第一次启动应用时，右下角会显示"首次使用"。

### 后续使用

每次启动应用时：

1. 应用先读取数据库中的时间（上一次启动时间）
2. 右下角显示上一次打开应用的时间
3. 然后更新数据库为当前时间（为下一次启动做准备）

### 时间格式

显示格式：`YYYY-MM-DD HH:MM:SS`

例如：`2026-02-06 10:50:40`

### 重要说明

- **显示的总是上一次启动的时间**，不是当前启动时间
- 每次启动后，数据库会更新为当前时间
- 下次启动时，会显示这次的启动时间

## 扩展性

此实现具有良好的扩展性，可以轻松添加其他应用设置：

- 用户偏好设置
- 窗口位置和大小
- 主题设置
- 其他持久化配置

只需在 `app_settings` bucket 中添加新的 key-value 对即可。

## 注意事项

1. 数据存储在用户本地，隐私安全
2. 数据库文件位于 `~/.happytools/data.db`
3. 时间使用本地时区
4. 应用关闭时不会自动更新时间（只在启动时更新）

## 未来改进

可以考虑以下改进：

1. 应用关闭时也更新时间
2. 添加使用时长统计
3. 添加使用频率分析
4. 提供时间格式自定义选项
