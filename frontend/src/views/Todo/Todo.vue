<template>
  <div class="todobox">
    <!-- 分类管理 -->
    <CategoryManager ref="categoryManagerRef" />

    <!-- 左右分栏布局 -->
    <div class="layout-container">
      <!-- 左侧：待办列表 -->
      <div class="left-panel">
        <!-- 待办列表头部 -->
        <div class="list-header">
          <span class="title">待办列表</span>
          <div class="header-right">
            <a-button type="link" size="small" @click="clearCompleted" :disabled="completedCount === 0" class="clear-btn">
              清除已完成
            </a-button>
            <span class="count">{{ completedCount }} / {{ todos.length }} 已完成</span>
          </div>
        </div>

        <!-- 待办列表 -->
        <a-list
          size="large"
          bordered
          :data-source="currentPageTodos"
          :loading="loading"
          :locale="{emptyText: '暂时没有待办'	}"
          class="todo-list"
        >
          <template #renderItem="{ item }">
            <TodoItem
              :todo="item"
              :category="getCategoryById(item.category_id)"
              @toggle="toggleTodo"
              @edit="startEdit"
              @delete="deleteTodo"
            />
          </template>
        </a-list>

        <!-- 分页器 -->
        <div class="pagination" v-if="totalPages > 1">
          <a-pagination
            v-model:current="currentPage"
            :total="todos.length"
            :page-size="pageSize"
            :show-size-changer="false"
            :show-total="total => `共 ${total} 条`"
            @change="handlePageChange"
          />
        </div>
      </div>

      <!-- 右侧：新建待办和分类 -->
      <div class="right-panel">
        <!-- 上部：新建待办 -->
        <div class="quick-add-section">
          <div class="section-title">
            <PlusOutlined />
            <span>快速添加待办</span>
          </div>
          <a-form :model="quickAddForm" layout="vertical">
            <a-row :gutter="12">
              <a-col :span="12">
                <a-form-item label="分类">
                  <CategorySelector v-model="quickAddForm.categoryId" />
                </a-form-item>
              </a-col>
              <a-col :span="12">
                <a-form-item label="优先级">
                  <a-select v-model:value="quickAddForm.priority" size="small">
                    <a-select-option :value="0">低</a-select-option>
                    <a-select-option :value="1">中</a-select-option>
                    <a-select-option :value="2">高</a-select-option>
                  </a-select>
                </a-form-item>
              </a-col>
            </a-row>
            <a-form-item label="截止日期">
              <a-date-picker
                v-model:value="quickAddForm.dueDate"
                placeholder="选择截止日期"
                style="width: 100%"
                size="small"
              />
            </a-form-item>
            <a-form-item>
              <a-input
                v-model:value="quickAddForm.title"
                placeholder="输入待办内容，按回车快速添加"
                size="large"
                @pressEnter="quickAdd"
              >
                <template #addonAfter>
                  <a-button type="primary" size="small" @click="quickAdd">
                    添加
                  </a-button>
                </template>
              </a-input>
            </a-form-item>
          </a-form>
        </div>

        <!-- 下部：分类展示 -->
        <div class="categories-section">
          <div class="section-header">
            <div class="section-title">
              <AppstoreOutlined />
              <span>分类列表</span>
            </div>
            <a-button type="link" size="small" @click="showAddCategoryModal" class="add-category-btn">
              <PlusOutlined /> 添加分类
            </a-button>
          </div>
          <div class="categories-cloud">
            <div
              v-for="category in categories"
              :key="category.id"
              class="category-tag"
              :class="{ active: selectedCategory === category.id }"
              :style="getCategoryStyle(category)"
              @click="selectCategory(category.id)"
            >
              {{ category.name }}
              <span class="category-tag-count">{{ getCategoryCount(category.id) }}</span>
              <span
                class="category-tag-delete"
                @click.stop="deleteCategory(category.id)"
                title="删除分类"
              >
                <CloseOutlined />
              </span>
            </div>
            <div
              v-if="categories.length === 0"
              class="category-empty"
            >
              暂无分类
            </div>
          </div>
          <div class="category-actions" v-if="selectedCategory !== null">
            <a-button
              type="link"
              size="small"
              @click="clearCategoryFilter"
            >
              清除筛选
            </a-button>
          </div>
        </div>
      </div>
    </div>

    <!-- 添加分类弹窗 -->
    <a-modal
      v-model:open="addCategoryModalVisible"
      title="添加分类"
      @ok="addCategory"
      @cancel="cancelAddCategory"
    >
      <a-form :model="newCategory" layout="vertical">
        <a-form-item label="分类名称" required>
          <a-input v-model:value="newCategory.name" placeholder="请输入分类名称" />
        </a-form-item>
        <a-form-item label="分类颜色">
          <div class="color-picker-wrapper">
            <div class="color-picker-row">
              <input
                type="color"
                v-model="newCategory.color"
                class="color-input-native"
              />
              <a-input
                v-model:value="newCategory.color"
                placeholder="选择或输入颜色代码"
                class="color-input-text"
              >
                <template #addonBefore>
                  <div
                    class="color-preview"
                    :style="{ backgroundColor: newCategory.color }"
                  ></div>
                </template>
              </a-input>
            </div>
            <div class="quick-colors">
              <div
                v-for="color in quickColors"
                :key="color"
                class="quick-color-item"
                :style="{ backgroundColor: color }"
                @click="newCategory.color = color"
                :title="color"
              ></div>
            </div>
          </div>
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 编辑弹窗 -->
    <a-modal
      v-model:open="editModalVisible"
      title="编辑待办"
      @ok="saveEdit"
      @cancel="cancelEdit"
    >
      <a-form :model="editForm" layout="vertical">
        <a-form-item label="标题">
          <a-input v-model:value="editForm.title" placeholder="待办标题" />
        </a-form-item>
        <a-form-item label="分类">
          <CategorySelector v-model="editForm.categoryId" />
        </a-form-item>
        <a-form-item label="截止日期">
          <a-date-picker
            v-model:value="editForm.dueDate"
            placeholder="选择截止日期"
            style="width: 100%"
          />
        </a-form-item>
        <a-form-item label="优先级">
          <a-select v-model:value="editForm.priority">
            <a-select-option :value="0">低</a-select-option>
            <a-select-option :value="1">中</a-select-option>
            <a-select-option :value="2">高</a-select-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed, onMounted, h } from 'vue'
import { PlusOutlined, AppstoreOutlined, CloseOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import dayjs, { Dayjs } from 'dayjs'
import { TodoService, Todo } from '../../../bindings/github.com/Aliuyanfeng/happytools/backend/services/todo'
import CategorySelector from './components/CategorySelector.vue'
import TodoItem from './components/TodoItem.vue'
import { useCategoryStore } from '../../stores/category'
import { useTodoStatsStore } from '../../stores/todoStats'

const todos = ref<Todo[]>([])
const loading = ref(false)
const categoryStore = useCategoryStore()
const todoStats = useTodoStatsStore()
const categories = computed(() => categoryStore.categories)

// 分页相关
const currentPage = ref(1)
const pageSize = ref(10)

// 分类筛选
const selectedCategory = ref<number | null>(null)

// 快速添加表单
const quickAddForm = ref({
  title: '',
  categoryId: null as number | null,
  dueDate: null as Dayjs | null,
  priority: 0
})

// 添加分类弹窗
const addCategoryModalVisible = ref(false)
const newCategory = ref({
  name: '',
  color: '#1890ff'
})

// 快速选择的常用颜色
const quickColors = [
  '#1890ff', // 蓝色
  '#52c41a', // 绿色
  '#faad14', // 橙色
  '#f5222d', // 红色
  '#722ed1', // 紫色
  '#eb2f96', // 粉色
  '#13c2c2', // 青色
  '#fa8c16', // 深橙色
  '#a0d911', // 柠檬绿
  '#2f54eb', // 深蓝色
  '#fa541c', // 火砖色
  '#607d8b', // 蓝灰色
]

// 编辑相关
const editModalVisible = ref(false)
const editingId = ref<number | null>(null)
const editForm = ref({
  title: '',
  categoryId: null as number | null,
  dueDate: null as Dayjs | null,
  priority: 0,
  completed: false
})

// 计算已完成数量
const completedCount = computed(() => todos.value.filter(t => t.completed).length)

// 计算总页数
const totalPages = computed(() => Math.ceil(todos.value.length / pageSize.value))

// 计算当前页的待办
const currentPageTodos = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return todos.value.slice(start, end)
})

// 根据 ID 获取分类
function getCategoryById(id: number | null) {
  if (id === null) return null
  return categoryStore.getCategoryById(id)
}

// 获取分类的待办数量
function getCategoryCount(categoryId: number) {
  return categoryStore.getCategoryCount(categoryId, todos.value)
}

// 获取分类标签样式（词云效果）
function getCategoryStyle(category: any) {
  const count = getCategoryCount(category.id)
  // 根据待办数量计算字体大小（12px - 20px）
  const fontSize = Math.min(12 + count * 0.8, 20)
  // 根据待办数量计算透明度（0.6 - 1.0）
  const opacity = Math.min(0.6 + count * 0.04, 1.0)

  return {
    backgroundColor: `${category.color}20`, // 20% 透明度背景
    color: category.color,
    fontSize: `${fontSize}px`,
    opacity: opacity,
  }
}

// 选择分类
function selectCategory(categoryId: number) {
  if (selectedCategory.value === categoryId) {
    selectedCategory.value = null
  } else {
    selectedCategory.value = categoryId
  }
  currentPage.value = 1 // 重置到第一页
}

// 清除分类筛选
function clearCategoryFilter() {
  selectedCategory.value = null
  currentPage.value = 1
}



// 加载所有待办
async function loadTodos() {
  loading.value = true
  try {
    const result = await TodoService.GetAll()
    todos.value = result || []
  } catch (e: any) {
    message.error('加载待办失败: ' + e.message)
  } finally {
    loading.value = false
  }
}

// 快速添加待办
async function quickAdd() {
  const title = quickAddForm.value.title.trim()
  if (!title) {
    message.warning('请输入待办内容')
    return
  }

  try {
    const dueDateStr = quickAddForm.value.dueDate ? quickAddForm.value.dueDate.format('YYYY-MM-DD') : null

    const todo = await TodoService.Add(
      title,
      quickAddForm.value.categoryId,
      dueDateStr,
      quickAddForm.value.priority
    )
    if (todo) {
      todos.value.unshift(todo)
      quickAddForm.value = { title: '', categoryId: null, dueDate: null, priority: 0 }
      message.success('添加成功')
      todoStats.refresh()
    }
  } catch (e: any) {
    message.error('添加失败: ' + e.message)
  }
}

// 显示添加分类弹窗
function showAddCategoryModal() {
  newCategory.value = {
    name: '',
    color: '#1890ff'
  }
  addCategoryModalVisible.value = true
}

// 取消添加分类
function cancelAddCategory() {
  addCategoryModalVisible.value = false
  newCategory.value = {
    name: '',
    color: '#1890ff'
  }
}

// 添加分类
async function addCategory() {
  const name = newCategory.value.name.trim()
  if (!name) {
    message.warning('请输入分类名称')
    return
  }

  try {
    const category = await categoryStore.addCategory(name, newCategory.value.color)
    if (category) {
      addCategoryModalVisible.value = false
      message.success('添加成功')
    }
  } catch (e: any) {
    // 错误信息已经在 store 中处理
  }
}

// 删除分类
async function deleteCategory(categoryId: number) {
  const count = getCategoryCount(categoryId)
  if (count > 0) {
    message.warning(`该分类下还有 ${count} 个待办，请先删除或移动这些待办`)
    return
  }

  try {
    await categoryStore.deleteCategory(categoryId)

    // 如果删除的是当前选中的分类，清除筛选
    if (selectedCategory.value === categoryId) {
      selectedCategory.value = null
    }

    message.success('删除成功')
  } catch (e: any) {
    // 错误信息已经在 store 中处理
  }
}

// 切换完成状态
async function toggleTodo(id: number) {
  try {
    await TodoService.Toggle(id)
    const todo = todos.value.find(t => t.id === id)
    if (todo) {
      todo.completed = !todo.completed
      todoStats.refresh()
    }
  } catch (e: any) {
    message.error('更新失败: ' + e.message)
  }
}

// 删除待办
async function deleteTodo(id: number) {
  try {
    await TodoService.Delete(id)
    todos.value = todos.value.filter(t => t.id !== id)
    message.success('删除成功')
    todoStats.refresh()
  } catch (e: any) {
    message.error('删除失败: ' + e.message)
  }
}

// 开始编辑
function startEdit(item: Todo) {
  editingId.value = item.id
  editForm.value = {
    title: item.title,
    categoryId: item.category_id,
    dueDate: item.due_date ? dayjs(item.due_date) : null,
    priority: item.priority ?? 0,
    completed: item.completed
  }
  editModalVisible.value = true
}

// 保存编辑
async function saveEdit() {
  const title = editForm.value.title.trim()
  if (!title) {
    message.warning('标题不能为空')
    return
  }

  if (editingId.value === null) return

  try {
    const dueDateStr = editForm.value.dueDate ? editForm.value.dueDate.format('YYYY-MM-DD') : null
    await TodoService.Update(
      editingId.value,
      title,
      editForm.value.completed,
      editForm.value.categoryId,
      dueDateStr,
      editForm.value.priority
    )

    const todo = todos.value.find(t => t.id === editingId.value)
    if (todo) {
      todo.title = title
      ;(todo as any).category_id = editForm.value.categoryId
      ;(todo as any).due_date = dueDateStr
      todo.priority = editForm.value.priority
      todo.completed = editForm.value.completed
    }

    message.success('更新成功')
    editModalVisible.value = false
    todoStats.refresh()
  } catch (e: any) {
    message.error('更新失败: ' + e.message)
  }
}

// 取消编辑
function cancelEdit() {
  editModalVisible.value = false
  editingId.value = null
  editForm.value = {
    title: '',
    categoryId: null,
    dueDate: null,
    priority: 0,
    completed: false
  }
}

// 清除已完成
async function clearCompleted() {
  const completedIds = todos.value.filter(t => t.completed).map(t => t.id)
  try {
    for (const id of completedIds) {
      await TodoService.Delete(id)
    }
    todos.value = todos.value.filter(t => !t.completed)
    message.success('已清除完成项')
  } catch (e: any) {
    message.error('清除失败: ' + e.message)
  }
}

// 页码变化
function handlePageChange(page: number) {
  currentPage.value = page
}

onMounted(() => {
  loadTodos()
  // 分类数据已经在 store 中自动加载
})
</script>

<style scoped>
.todobox {
  width: 100%;
  max-width: 1400px;
  /* margin: 24px auto; */
  padding: 24px;
  background-color: #ffffff;
  border-radius: 16px;
  box-shadow: 0 12px 30px rgba(15, 23, 42, 0.08);
  box-sizing: border-box;
}

/* 布局容器 */
.layout-container {
  display: flex;
  gap: 24px;
  min-height: 664px;
  width: 100%;
  box-sizing: border-box;
}

/* 左侧面板 */
.left-panel {
  flex: 0 0 60%;
  display: flex;
  flex-direction: column;
  gap: 14px;
  min-width: 0;
  overflow: hidden;
}

.list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 4px;
}

.list-header .title {
  font-size: 18px;
  font-weight: bold;
  color: #1f2937;
}

.list-header .header-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.list-header .clear-btn {
  padding: 0;
  height: auto;
  color: #6b7280;
}

.list-header .clear-btn:hover:not(:disabled) {
  color: #ff4d4f;
}

.list-header .clear-btn:disabled {
  color: #d9d9d9;
}

.list-header .count {
  font-size: 14px;
  color: #6b7280;
}

.todo-list {
  flex: 1;
  min-height: 400px;
  width: 100%;
}

.todo-list :deep(.ant-list) {
  width: 100%;
}

.todo-list :deep(.ant-list-item) {
  width: 100%;
}

.pagination {
  display: flex;
  justify-content: center;
  padding: 16px 0;
}

/* 右侧面板 */
.right-panel {
  flex: 0 0 38%;
  display: flex;
  flex-direction: column;
  gap: 24px;
  min-width: 0;
  overflow: hidden;
}

/* 快速添加区域 */
.quick-add-section {
  background: linear-gradient(360deg, #648bdb 0%, #2f6cd4 100%);
  padding: 20px;
  border-radius: 12px;
  color: white;
}

.quick-add-section .section-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 16px;
}

.quick-add-section :deep(.ant-form-item-label > label) {
  color: rgba(255, 255, 255, 0.9);
}

.quick-add-section :deep(.ant-input),
.quick-add-section :deep(.ant-select-selector) {
  background-color: rgba(255, 255, 255, 0.9);
}

/* 分类展示区域 */
.categories-section {
  flex: 1;
  /* background-color: #f9fafb; */
  padding: 20px;
  border-radius: 12px;
  display: flex;
  flex-direction: column;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.categories-section .section-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
}

.add-category-btn {
  color: #667eea;
  padding: 0;
  height: auto;
}

.add-category-btn:hover {
  color: #764ba2;
}

.categories-cloud {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  padding: 16px;
  flex: 1;
  overflow-y: auto;
  max-height: 400px;
  align-content: flex-start;
}

.category-tag {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  border-radius: 20px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
  border: 2px solid transparent;
  user-select: none;
  white-space: nowrap;
}

.category-tag:hover {
  transform: scale(1.1) translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  border-color: currentColor;
}

.category-tag.active {
  transform: scale(1.15);
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.2);
  border-color: currentColor;
  font-weight: 600;
}

.category-tag-count {
  font-size: 0.85em;
  opacity: 0.8;
  background-color: rgba(255, 255, 255, 0.5);
  padding: 2px 8px;
  border-radius: 10px;
  font-weight: 600;
}

.category-tag-delete {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 18px;
  height: 18px;
  margin-left: 4px;
  border-radius: 50%;
  background-color: rgba(0, 0, 0, 0.1);
  opacity: 0;
  transition: all 0.2s;
  cursor: pointer;
}

.category-tag-delete:hover {
  background-color: rgba(255, 0, 0, 0.2);
  color: #ff4d4f;
}

.category-tag:hover .category-tag-delete {
  opacity: 1;
}

.category-tag.active .category-tag-delete {
  opacity: 1;
}

.category-empty {
  text-align: center;
  color: #9ca3af;
  padding: 40px 20px;
  width: 100%;
}

.category-actions {
  display: flex;
  justify-content: center;
  padding-top: 16px;
  border-top: 1px solid #e5e7eb;
}

/* 颜色选择器样式 */
.color-picker-wrapper {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.color-picker-row {
  display: flex;
  gap: 12px;
  align-items: center;
}

.color-input-native {
  width: 48px;
  height: 32px;
  border: 1px solid #d9d9d9;
  border-radius: 4px;
  cursor: pointer;
  padding: 2px;
  background: #fff;
}

.color-input-native:hover {
  border-color: #40a9ff;
}

.color-input-text {
  flex: 1;
}

.color-preview {
  width: 24px;
  height: 24px;
  border-radius: 4px;
  border: 1px solid #d9d9d9;
  background-color: #ffffff;
}

.quick-colors {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(32px, 1fr));
  gap: 8px;
  padding: 12px;
  background-color: #f9fafb;
  border-radius: 8px;
}

.quick-color-item {
  width: 32px;
  height: 32px;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s;
  border: 2px solid transparent;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.quick-color-item:hover {
  transform: scale(1.1);
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.15);
  border-color: #1890ff;
}

/* 响应式布局 */
@media (max-width: 1200px) {
  /* .layout-container {
    flex-direction: column;
  }

  .left-panel {
    flex: 1;
    width: 100%;
  }

  .right-panel {
    flex: 1;
    width: 100%;
  } */

  /* .categories-grid {
    max-height: 250px;
  } */
}

/* 全局样式覆盖 */
.todobox :deep(*) {
  box-sizing: border-box;
}
</style>
