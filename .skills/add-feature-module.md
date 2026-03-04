# 添加新功能模块 Skill

## Skill 描述
这个 skill 用于在 HappyTools 项目中快速添加新的功能模块，包括后端服务、数据存储、前端页面和路由配置。

## 适用场景
- 需要添加新的业务功能模块
- 需要创建完整的后端服务 + 前端页面
- 需要实现数据持久化存储

## 执行步骤

### 1. 需求确认
询问用户以下信息：
- 功能模块名称（英文，如：note, bookmark, reminder）
- 功能模块中文名称（如：笔记管理、书签管理、提醒管理）
- 功能描述和主要功能点
- 是否需要数据持久化存储
- 是否需要分类管理
- 是否需要时间相关字段（创建时间、更新时间）

### 2. 后端服务实现

#### 2.1 创建服务目录
在 `backend/services/` 下创建新的服务包：
```
backend/services/[模块名]/
├── [模块名].go          # 服务实现
└── models.go            # 数据模型定义
```

#### 2.2 实现数据模型 (models.go)
```go
package [模块名]

// [模块名首字母大写] 数据模型
type [模块名首字母大写] struct {
    ID          int     `json:"id"`
    Title       string  `json:"title"`
    Content     string  `json:"content"`
    CategoryID  *int    `json:"category_id"`    // 可选：分类ID
    CreatedAt   string  `json:"created_at"`
    UpdatedAt   string  `json:"updated_at"`
    Status      int     `json:"status"`         // 状态：0-禁用，1-启用
}
```

#### 2.3 实现服务层 ([模块名].go)
```go
package [模块名]

import (
    "encoding/json"
    "errors"
    "time"
    
    "github.com/Aliuyanfeng/happytools/backend/store"
    "go.etcd.io/bbolt"
)

type [模块名首字母大写]Service struct{}

func New[模块名首字母大写]Service() *[模块名首字母大写]Service {
    return &[模块名首字母大写]Service{}
}

// GetAll 获取所有数据
func (s *[模块名首字母大写]Service) GetAll() ([][模块名首字母大写], error) {
    var items [][模块名首字母大写]
    
    err := store.DB.View(func(tx *bbolt.Tx) error {
        bucket := tx.Bucket(store.[模块名]Bucket)
        if bucket == nil {
            return nil
        }
        
        return bucket.ForEach(func(k, v []byte) error {
            var item [模块名首字母大写]
            if err := json.Unmarshal(v, &item); err != nil {
                return err
            }
            items = append(items, item)
            return nil
        })
    })
    
    return items, err
}

// Add 添加新数据
func (s *[模块名首字母大写]Service) Add(title, content string, categoryID *int) (*[模块名首字母大写], error) {
    if title == "" {
        return nil, errors.New("title cannot be empty")
    }
    
    item := [模块名首字母大写]{
        Title:      title,
        Content:    content,
        CategoryID: categoryID,
        CreatedAt:  time.Now().Format("2006-01-02 15:04:05"),
        UpdatedAt:  time.Now().Format("2006-01-02 15:04:05"),
        Status:     1,
    }
    
    err := store.DB.Update(func(tx *bbolt.Tx) error {
        bucket := tx.Bucket(store.[模块名]Bucket)
        if bucket == nil {
            return errors.New("bucket not found")
        }
        
        // 生成自增ID
        id, _ := bucket.NextSequence()
        item.ID = int(id)
        
        data, err := json.Marshal(item)
        if err != nil {
            return err
        }
        
        return bucket.Put(itob(item.ID), data)
    })
    
    if err != nil {
        return nil, err
    }
    
    return &item, nil
}

// Update 更新数据
func (s *[模块名首字母大写]Service) Update(id int, title, content string, categoryID *int, status int) error {
    if title == "" {
        return errors.New("title cannot be empty")
    }
    
    return store.DB.Update(func(tx *bbolt.Tx) error {
        bucket := tx.Bucket(store.[模块名]Bucket)
        if bucket == nil {
            return errors.New("bucket not found")
        }
        
        data := bucket.Get(itob(id))
        if data == nil {
            return errors.New("item not found")
        }
        
        var item [模块名首字母大写]
        if err := json.Unmarshal(data, &item); err != nil {
            return err
        }
        
        // 更新字段
        item.Title = title
        item.Content = content
        item.CategoryID = categoryID
        item.Status = status
        item.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
        
        updatedData, err := json.Marshal(item)
        if err != nil {
            return err
        }
        
        return bucket.Put(itob(id), updatedData)
    })
}

// Delete 删除数据
func (s *[模块名首字母大写]Service) Delete(id int) error {
    return store.DB.Update(func(tx *bbolt.Tx) error {
        bucket := tx.Bucket(store.[模块名]Bucket)
        if bucket == nil {
            return errors.New("bucket not found")
        }
        
        return bucket.Delete(itob(id))
    })
}

// itob 将 int 转换为 []byte
func itob(v int) []byte {
    b := make([]byte, 8)
    b[0] = byte(v >> 56)
    b[1] = byte(v >> 48)
    b[2] = byte(v >> 40)
    b[3] = byte(v >> 32)
    b[4] = byte(v >> 24)
    b[5] = byte(v >> 16)
    b[6] = byte(v >> 8)
    b[7] = byte(v)
    return b
}
```

### 3. 数据存储层实现

#### 3.1 添加 Bucket 定义
在 `backend/store/store.go` 中添加新的 bucket：
```go
var (
    todoBucket        = []byte("todos")
    categoryBucket    = []byte("categories")
    dailyReportBucket = []byte("daily_reports")
    appSettingsBucket = []byte("app_settings")
    [模块名]Bucket     = []byte("[模块名]s")  // 添加这行
)
```

#### 3.2 初始化 Bucket
在 `InitDB()` 函数中添加 bucket 创建：
```go
func InitDB() error {
    // ... 现有代码
    
    // 创建新模块的 bucket
    err = db.Update(func(tx *bbolt.Tx) error {
        _, err := tx.CreateBucketIfNotExists([模块名]Bucket)
        return err
    })
    if err != nil {
        return err
    }
    
    // ... 其他代码
}
```

### 4. 服务注册

在 `main.go` 中注册新服务：
```go
import (
    // ... 现有导入
    "[模块名] "github.com/Aliuyanfeng/happytools/backend/services/[模块名]"
)

func main() {
    // ... 现有代码
    
    app.RegisterService(application.NewService([模块名].New[模块名首字母大写]Service()))
    
    // ... 其他代码
}
```

### 5. 前端页面实现

#### 5.1 创建页面组件
在 `frontend/src/views/` 下创建新目录和组件：
```
frontend/src/views/[模块名首字母大写]/
├── Index.vue           # 主页面
├── [模块名首字母大写]Form.vue  # 表单组件
└── [模块名首字母大写]List.vue  # 列表组件
```

#### 5.2 实现主页面 (Index.vue)
```vue
<template>
  <div class="p-6">
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-white mb-4">{{ meta.title }}</h1>
      <a-button type="primary" @click="showAddModal">
        <template #icon><PlusOutlined /></template>
        添加[中文名称]
      </a-button>
    </div>

    <!-- 列表 -->
    <[模块名首字母大写]List 
      :data="items" 
      @edit="handleEdit"
      @delete="handleDelete"
    />

    <!-- 添加/编辑弹窗 -->
    <[模块名首字母大写]Form
      v-model:visible="formVisible"
      :item="currentItem"
      @success="loadData"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import { [模块名首字母大写]Service } from '../../bindings/github.com/Aliuyanfeng/happytools/backend/services/[模块名]'
import [模块名首字母大写]List from './[模块名首字母大写]List.vue'
import [模块名首字母大写]Form from './[模块名首字母大写]Form.vue'

const items = ref([])
const formVisible = ref(false)
const currentItem = ref(null)

const loadData = async () => {
  try {
    const result = await [模块名首字母大写]Service.GetAll()
    items.value = result || []
  } catch (error) {
    message.error('加载数据失败')
  }
}

const showAddModal = () => {
  currentItem.value = null
  formVisible.value = true
}

const handleEdit = (item: any) => {
  currentItem.value = item
  formVisible.value = true
}

const handleDelete = async (id: number) => {
  try {
    await [模块名首字母大写]Service.Delete(id)
    message.success('删除成功')
    loadData()
  } catch (error) {
    message.error('删除失败')
  }
}

onMounted(() => {
  loadData()
})
</script>
```

### 6. 路由配置

在 `frontend/src/router/routes.ts` 中添加路由：
```typescript
{
  path: '/[模块名]',
  name: '[模块名]',
  meta: { title: '[中文名称]' },
  component: () => import('../views/[模块名首字母大写]/Index.vue')
}
```

### 7. 导航菜单更新

在导航菜单组件中添加新菜单项：
```vue
<a-menu-item key="/[模块名]">
  <template #icon>
    <[图标名] />
  </template>
  <span>[中文名称]</span>
</a-menu-item>
```

### 8. 测试验证

#### 8.1 启动开发服务器
```bash
wails3 dev -config ./build/config.yml -port 9250
```

#### 8.2 验证功能
- 检查 TypeScript 绑定是否自动生成
- 测试添加、编辑、删除功能
- 验证数据持久化
- 检查页面路由和导航

## 注意事项

1. **命名规范**：
   - 模块名使用小写字母，如：note, bookmark
   - 服务名首字母大写，如：NoteService, BookmarkService
   - 数据库 bucket 名使用复数形式，如：notes, bookmarks

2. **错误处理**：
   - 所有错误必须处理并返回
   - 使用 errors.New 创建错误信息
   - 前端使用 message 提示用户

3. **数据验证**：
   - 后端必须验证必填字段
   - 前端表单添加验证规则

4. **类型安全**：
   - 使用 TypeScript 类型注解
   - 可空字段使用指针类型

5. **代码风格**：
   - Go 代码使用 gofmt 格式化
   - Vue 组件使用 `<script setup>` 语法
   - 样式使用 `scoped` 限定作用域

## 示例：添加笔记管理模块

### 需求
- 模块名：note
- 中文名称：笔记管理
- 功能：创建、编辑、删除笔记，支持分类

### 执行结果
- 后端服务：`backend/services/note/note.go`
- 数据模型：`backend/services/note/models.go`
- 数据存储：`backend/store/store.go` 添加 `noteBucket`
- 服务注册：`main.go` 注册 `NoteService`
- 前端页面：`frontend/src/views/Note/Index.vue`
- 路由配置：`/note` 路由
- 导航菜单：添加"笔记管理"菜单项

## 总结

这个 skill 提供了完整的端到端功能模块开发流程，遵循 Wails v3 + Vue 3 的最佳实践，确保代码质量和一致性。通过这个 skill，可以快速、规范地添加新功能模块到 HappyTools 项目中。
