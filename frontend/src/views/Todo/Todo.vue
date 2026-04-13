<template>
  <div class="todo-page">

    <!-- 顶部输入区 -->
    <div class="input-zone">
      <div class="input-row">
        <div class="input-wrap">
          <PlusCircleOutlined class="input-icon" />
          <input
            v-model="quickAddForm.title"
            class="main-input"
            placeholder="添加一项待办..."
            @keydown.enter="quickAdd"
          />
        </div>
        <div class="input-meta">
          <CategorySelector v-model="quickAddForm.categoryId" class="meta-select" size="small" />
          <a-select v-model:value="quickAddForm.priority" size="small" class="meta-select priority-select">
            <a-select-option :value="0"><span class="dot dot-low" />低</a-select-option>
            <a-select-option :value="1"><span class="dot dot-mid" />中</a-select-option>
            <a-select-option :value="2"><span class="dot dot-high" />高</a-select-option>
          </a-select>
          <a-date-picker
            v-model:value="quickAddForm.dueDate"
            placeholder="截止日期"
            size="small"
            class="meta-select"
            style="width: 120px"
          />
          <a-button type="primary" size="small" @click="quickAdd" class="add-btn">添加</a-button>
        </div>
      </div>
    </div>

    <!-- 分类筛选栏 -->
    <div class="filter-bar">
      <div class="filter-tabs">
        <span
          class="filter-tab"
          :class="{ active: selectedCategory === null }"
          @click="clearCategoryFilter"
        >
          全部
          <span class="tab-count">{{ todos.length }}</span>
        </span>
        <span
          v-for="cat in categories"
          :key="cat.id"
          class="filter-tab"
          :class="{ active: selectedCategory === cat.id }"
          @click="selectCategory(cat.id)"
        >
          <span class="tab-dot" :style="{ background: cat.color }" />
          {{ cat.name }}
          <span class="tab-count">{{ getCategoryCount(cat.id) }}</span>
          <span class="tab-del" @click.stop="deleteCategory(cat.id)" title="删除分类">
            <CloseOutlined />
          </span>
        </span>
      </div>
      <div class="filter-actions">
        <span class="progress-text">{{ completedCount }} / {{ todos.length }} 已完成</span>
        <div class="progress-bar">
          <div
            class="progress-fill"
            :style="{ width: todos.length ? (completedCount / todos.length * 100) + '%' : '0%' }"
          />
        </div>
        <a-button
          type="text"
          size="small"
          :disabled="completedCount === 0"
          @click="clearCompleted"
          class="clear-btn"
        >清除已完成</a-button>
        <a-button type="text" size="small" @click="showAddCategoryModal" class="manage-btn">
          <PlusOutlined /> 分类
        </a-button>
      </div>
    </div>

    <!-- 任务列表 -->
    <div class="list-area">
      <div v-if="loading" class="state-view">
        <a-spin size="large" />
      </div>
      <div v-else-if="filteredTodos.length === 0" class="state-view">
        <div class="empty-illus">
          <CheckCircleOutlined class="empty-icon" />
          <p>{{ selectedCategory ? '该分类下暂无待办' : '暂无待办，享受空闲时光吧 ☕' }}</p>
        </div>
      </div>
      <div v-else class="task-groups">
        <!-- 待处理 -->
        <div v-if="pendingTodos.length" class="task-group">
          <div class="group-label">待处理 · {{ pendingTodos.length }}</div>
          <TodoItem
            v-for="item in pendingTodos"
            :key="item.id"
            :todo="item"
            :category="getCategoryById(item.category_id)"
            @toggle="toggleTodo"
            @edit="startEdit"
            @delete="deleteTodo"
          />
        </div>
        <!-- 已完成 -->
        <div v-if="doneTodos.length" class="task-group">
          <div class="group-label done-label">已完成 · {{ doneTodos.length }}</div>
          <TodoItem
            v-for="item in doneTodos"
            :key="item.id"
            :todo="item"
            :category="getCategoryById(item.category_id)"
            @toggle="toggleTodo"
            @edit="startEdit"
            @delete="deleteTodo"
          />
        </div>
      </div>
    </div>

    <!-- 分页 -->
    <div class="pagination" v-if="totalPages > 1">
      <a-pagination
        v-model:current="currentPage"
        :total="filteredTodos.length"
        :page-size="pageSize"
        :show-size-changer="false"
        size="small"
        @change="handlePageChange"
      />
    </div>

    <!-- 新建分类弹窗 -->
    <a-modal v-model:open="addCategoryModalVisible" title="新建分类" @ok="addCategory" @cancel="cancelAddCategory" width="380px">
      <a-form :model="newCategory" layout="vertical">
        <a-form-item label="名称" required>
          <a-input v-model:value="newCategory.name" placeholder="分类名称" />
        </a-form-item>
        <a-form-item label="颜色">
          <div class="color-row">
            <input type="color" v-model="newCategory.color" class="color-native" />
            <a-input v-model:value="newCategory.color" class="color-text">
              <template #addonBefore>
                <div class="color-preview" :style="{ backgroundColor: newCategory.color }" />
              </template>
            </a-input>
          </div>
          <div class="quick-colors">
            <div v-for="c in quickColors" :key="c" class="qc" :style="{ background: c }" @click="newCategory.color = c" />
          </div>
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 编辑弹窗 -->
    <a-modal v-model:open="editModalVisible" title="编辑待办" @ok="saveEdit" @cancel="cancelEdit" width="400px">
      <a-form :model="editForm" layout="vertical">
        <a-form-item label="标题">
          <a-input v-model:value="editForm.title" />
        </a-form-item>
        <a-form-item label="分类">
          <CategorySelector v-model="editForm.categoryId" />
        </a-form-item>
        <a-form-item label="截止日期">
          <a-date-picker v-model:value="editForm.dueDate" style="width:100%" />
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
import { ref, computed, onMounted } from 'vue'
import { PlusOutlined, PlusCircleOutlined, CheckCircleOutlined, CloseOutlined } from '@ant-design/icons-vue'
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

const currentPage = ref(1)
const pageSize = ref(20)
const selectedCategory = ref<number | null>(null)

const quickAddForm = ref({ title: '', categoryId: null as number | null, dueDate: null as Dayjs | null, priority: 0 })
const addCategoryModalVisible = ref(false)
const newCategory = ref({ name: '', color: '#1890ff' })
const quickColors = ['#1890ff','#52c41a','#faad14','#f5222d','#722ed1','#eb2f96','#13c2c2','#fa8c16','#a0d911','#2f54eb','#fa541c','#607d8b']

const editModalVisible = ref(false)
const editingId = ref<number | null>(null)
const editForm = ref({ title: '', categoryId: null as number | null, dueDate: null as Dayjs | null, priority: 0, completed: false })

const completedCount = computed(() => todos.value.filter(t => t.completed).length)

const filteredTodos = computed(() => {
  if (selectedCategory.value === null) return todos.value
  return todos.value.filter(t => t.category_id === selectedCategory.value)
})

const pendingTodos = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  return filteredTodos.value.filter(t => !t.completed).slice(start, start + pageSize.value)
})

const doneTodos = computed(() => filteredTodos.value.filter(t => t.completed))
const totalPages = computed(() => Math.ceil(filteredTodos.value.filter(t => !t.completed).length / pageSize.value))

function getCategoryById(id: number | null) {
  return id === null ? null : categoryStore.getCategoryById(id)
}

function getCategoryCount(categoryId: number) {
  return categoryStore.getCategoryCount(categoryId, todos.value)
}

function selectCategory(id: number) {
  selectedCategory.value = selectedCategory.value === id ? null : id
  currentPage.value = 1
}

function clearCategoryFilter() {
  selectedCategory.value = null
  currentPage.value = 1
}

async function loadTodos() {
  loading.value = true
  try {
    todos.value = (await TodoService.GetAll()) || []
  } catch (e: any) {
    message.error('加载失败: ' + e.message)
  } finally {
    loading.value = false
  }
}

async function quickAdd() {
  const title = quickAddForm.value.title.trim()
  if (!title) { message.warning('请输入待办内容'); return }
  try {
    const dueDateStr = quickAddForm.value.dueDate?.format('YYYY-MM-DD') ?? null
    const todo = await TodoService.Add(title, quickAddForm.value.categoryId, dueDateStr, quickAddForm.value.priority)
    if (todo) {
      todos.value.unshift(todo)
      quickAddForm.value = { title: '', categoryId: null, dueDate: null, priority: 0 }
      todoStats.refresh()
    }
  } catch (e: any) { message.error('添加失败: ' + e.message) }
}

function showAddCategoryModal() { newCategory.value = { name: '', color: '#1890ff' }; addCategoryModalVisible.value = true }
function cancelAddCategory() { addCategoryModalVisible.value = false }

async function addCategory() {
  const name = newCategory.value.name.trim()
  if (!name) { message.warning('请输入分类名称'); return }
  try {
    const cat = await categoryStore.addCategory(name, newCategory.value.color)
    if (cat) { addCategoryModalVisible.value = false; message.success('添加成功') }
  } catch {}
}

async function deleteCategory(categoryId: number) {
  const count = getCategoryCount(categoryId)
  if (count > 0) { message.warning(`该分类下还有 ${count} 个待办`); return }
  try {
    await categoryStore.deleteCategory(categoryId)
    if (selectedCategory.value === categoryId) selectedCategory.value = null
  } catch {}
}

async function toggleTodo(id: number) {
  try {
    await TodoService.Toggle(id)
    const t = todos.value.find(t => t.id === id)
    if (t) { t.completed = !t.completed; todoStats.refresh() }
  } catch (e: any) { message.error('更新失败') }
}

async function deleteTodo(id: number) {
  try {
    await TodoService.Delete(id)
    todos.value = todos.value.filter(t => t.id !== id)
    todoStats.refresh()
  } catch (e: any) { message.error('删除失败') }
}

function startEdit(item: Todo) {
  editingId.value = item.id
  editForm.value = { title: item.title, categoryId: item.category_id, dueDate: item.due_date ? dayjs(item.due_date) : null, priority: item.priority ?? 0, completed: item.completed }
  editModalVisible.value = true
}

async function saveEdit() {
  const title = editForm.value.title.trim()
  if (!title) { message.warning('标题不能为空'); return }
  if (editingId.value === null) return
  try {
    const dueDateStr = editForm.value.dueDate?.format('YYYY-MM-DD') ?? null
    await TodoService.Update(editingId.value, title, editForm.value.completed, editForm.value.categoryId, dueDateStr, editForm.value.priority)
    const t = todos.value.find(t => t.id === editingId.value)
    if (t) { t.title = title; (t as any).category_id = editForm.value.categoryId; (t as any).due_date = dueDateStr; t.priority = editForm.value.priority; t.completed = editForm.value.completed }
    message.success('更新成功'); editModalVisible.value = false; todoStats.refresh()
  } catch (e: any) { message.error('更新失败') }
}

function cancelEdit() {
  editModalVisible.value = false; editingId.value = null
  editForm.value = { title: '', categoryId: null, dueDate: null, priority: 0, completed: false }
}

async function clearCompleted() {
  const ids = todos.value.filter(t => t.completed).map(t => t.id)
  try {
    for (const id of ids) await TodoService.Delete(id)
    todos.value = todos.value.filter(t => !t.completed)
    message.success('已清除'); todoStats.refresh()
  } catch (e: any) { message.error('清除失败') }
}

function handlePageChange(page: number) { currentPage.value = page }

onMounted(loadTodos)
</script>

<style scoped>
.todo-page {
  height: calc(100vh - 120px);
  display: flex;
  flex-direction: column;
  background: #f7f8fa;
  overflow: hidden;
}

/* ── 输入区 ── */
.input-zone {
  padding: 20px 28px 0;
  flex-shrink: 0;
}

.input-row {
  background: #fff;
  border-radius: 14px;
  box-shadow: 0 2px 12px rgba(0,0,0,0.06);
  padding: 14px 18px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.input-wrap {
  display: flex;
  align-items: center;
  gap: 10px;
}

.input-icon {
  font-size: 20px;
  color: #1677ff;
  flex-shrink: 0;
}

.main-input {
  flex: 1;
  border: none;
  outline: none;
  font-size: 15px;
  color: #1a1a1a;
  background: transparent;
  font-family: inherit;
}

.main-input::placeholder {
  color: #bbb;
}

.input-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  padding-left: 30px;
}

.meta-select {
  min-width: 90px;
}

.add-btn {
  margin-left: auto;
  border-radius: 8px;
  padding: 0 18px;
}

/* ── 筛选栏 ── */
.filter-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 28px 0;
  flex-shrink: 0;
  gap: 16px;
}

.filter-tabs {
  display: flex;
  align-items: center;
  gap: 4px;
  flex-wrap: wrap;
}

.filter-tab {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  padding: 5px 12px;
  border-radius: 20px;
  font-size: 13px;
  color: #666;
  cursor: pointer;
  transition: all 0.15s;
  user-select: none;
  background: transparent;
}

.filter-tab:hover {
  background: #fff;
  color: #1a1a1a;
}

.filter-tab.active {
  background: #fff;
  color: #1677ff;
  font-weight: 600;
  box-shadow: 0 1px 6px rgba(0,0,0,0.08);
}

.tab-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  flex-shrink: 0;
}

.tab-count {
  font-size: 11px;
  color: #999;
  background: #f0f0f0;
  padding: 1px 5px;
  border-radius: 8px;
  font-weight: 400;
}

.filter-tab.active .tab-count {
  background: #e6f4ff;
  color: #1677ff;
}

.tab-del {
  display: none;
  align-items: center;
  justify-content: center;
  width: 14px;
  height: 14px;
  border-radius: 3px;
  font-size: 9px;
  color: #999;
  margin-left: 1px;
  transition: all 0.15s;
}

.filter-tab:hover .tab-del {
  display: inline-flex;
}

.tab-del:hover {
  background: #fff1f0;
  color: #ff4d4f;
}

.filter-actions {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-shrink: 0;
}

.progress-text {
  font-size: 12px;
  color: #999;
  white-space: nowrap;
}

.progress-bar {
  width: 80px;
  height: 4px;
  background: #e8e8e8;
  border-radius: 2px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: #52c41a;
  border-radius: 2px;
  transition: width 0.3s ease;
}

.clear-btn {
  font-size: 12px;
  color: #bbb;
  padding: 0 6px;
}
.clear-btn:not(:disabled):hover { color: #ff4d4f !important; }

.manage-btn {
  font-size: 12px;
  color: #999;
  padding: 0 6px;
}
.manage-btn:hover { color: #1677ff !important; }

/* ── 列表区 ── */
.list-area {
  flex: 1;
  overflow-y: auto;
  padding: 14px 28px 0;
  min-height: 0;
}

.state-view {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  min-height: 200px;
}

.empty-illus {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  color: #ccc;
}

.empty-icon {
  font-size: 48px;
  color: #e0e0e0;
}

.empty-illus p {
  font-size: 14px;
  color: #bbb;
  margin: 0;
}

.task-groups {
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding-bottom: 16px;
}

.task-group {
  background: #fff;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 1px 4px rgba(0,0,0,0.04);
}

.group-label {
  font-size: 12px;
  font-weight: 600;
  color: #999;
  padding: 10px 18px 8px;
  letter-spacing: 0.5px;
  text-transform: uppercase;
  border-bottom: 1px solid #f5f5f5;
}

.done-label {
  color: #bbb;
}

/* ── 分页 ── */
.pagination {
  padding: 12px 28px;
  display: flex;
  justify-content: center;
  flex-shrink: 0;
}

/* ── 颜色选择器 ── */
.color-row {
  display: flex;
  gap: 10px;
  align-items: center;
  margin-bottom: 12px;
}

.color-native {
  width: 40px;
  height: 32px;
  border: 1px solid #d9d9d9;
  border-radius: 6px;
  cursor: pointer;
  padding: 2px;
  background: #fff;
  flex-shrink: 0;
}

.color-text { flex: 1; }

.color-preview {
  width: 20px;
  height: 20px;
  border-radius: 4px;
  border: 1px solid #d9d9d9;
}

.quick-colors {
  display: grid;
  grid-template-columns: repeat(6, 1fr);
  gap: 8px;
}

.qc {
  width: 100%;
  aspect-ratio: 1;
  border-radius: 6px;
  cursor: pointer;
  transition: transform 0.15s;
  border: 2px solid transparent;
}

.qc:hover {
  transform: scale(1.15);
  border-color: rgba(0,0,0,0.15);
}
</style>
