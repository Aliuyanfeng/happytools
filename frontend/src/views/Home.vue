<template>
  <div class="home">
    <!-- ═══ 顶部标题栏 ═══ -->
    <header class="home-header">
      <div class="header-left">
        <!-- 主机信息 -->
        <div class="host-info">
          <LaptopOutlined class="host-icon" />
          <span class="host-name">{{ hostName }}</span>
          <span class="host-sep">/</span>
          <span class="host-os">{{ hostOs }}</span>
        </div>
        <!-- 系统资源极简预览 -->
        <div class="sys-stats">
          <span class="stat-item" :class="cpuClass">
            <DashboardOutlined class="stat-icon" />
            <span class="stat-val">{{ cpuPercent }}%</span>
          </span>
          <span class="stat-item" :class="memClass">
            <DatabaseOutlined class="stat-icon" />
            <span class="stat-val">{{ memPercent }}%</span>
          </span>
        </div>
      </div>

      <div class="header-right">
        <!-- 全局搜索 -->
        <div class="search-box">
          <SearchOutlined class="search-icon" />
          <input
            v-model="homeStore.searchQuery"
            :placeholder="t('home.search')"
            class="search-input"
          />
          <span v-if="homeStore.searchQuery" class="search-clear" @click="homeStore.searchQuery = ''">×</span>
        </div>

        <!-- 待办气泡 -->
        <a-popover trigger="click" placement="bottomRight" :overlay-style="{ minWidth: '300px', maxWidth: '340px' }">
          <template #content>
            <div class="todo-preview-panel">
              <div class="todo-preview-header">
                <span class="todo-preview-title">{{ t('home.todoPreviewTitle') }}</span>
                <span class="todo-preview-count">{{ todoStats.pendingCount }}</span>
              </div>
              <div v-if="allTodos.length === 0" class="todo-empty">{{ t('home.noPendingTodos') }}</div>
              <div v-else class="todo-preview-list">
                <div
                  v-for="todo in allTodos"
                  :key="todo.id"
                  class="todo-preview-item"
                  :class="{ 'is-done': todo.completed }"
                >
                  <span
                    class="todo-check"
                    :class="{ checked: todo.completed }"
                    @click.stop="toggleTodoInPreview(todo)"
                  >
                    <CheckOutlined v-if="todo.completed" />
                  </span>
                  <span class="todo-preview-title-text">{{ todo.title }}</span>
                  <span v-if="todo.due_date" class="todo-due" :class="{ overdue: todo.status === 2, warning: todo.status === 1 }">
                    {{ todo.due_date }}
                  </span>
                </div>
              </div>
              <div v-if="allTodos.length > 0" class="todo-preview-footer">
                <span class="todo-view-all" @click="go('/todo')">{{ t('home.viewAllTodos') }} →</span>
              </div>
            </div>
          </template>
          <div class="todo-bubble">
            <CheckCircleOutlined />
            <span v-if="todoStats.pendingCount > 0" class="todo-badge">{{ todoStats.pendingCount }}</span>
          </div>
        </a-popover>

        <!-- 时钟 -->
        <span class="header-clock">{{ currentTime }}</span>
      </div>
    </header>

    <!-- ═══ 分类导航 ═══ -->
    <nav class="category-nav">
      <button
        v-for="cat in categories"
        :key="cat.id"
        class="cat-btn"
        :class="{ active: homeStore.activeCategory === cat.id }"
        @click="homeStore.setCategory(cat.id)"
      >
        <component :is="cat.icon" class="cat-icon" />
        <span>{{ t(`home.category.${cat.id}`) }}</span>
      </button>
    </nav>

    <!-- ═══ 模块卡片网格 ═══ -->
    <main class="home-main">
      <draggable
        v-model="draggableModules"
        item-key="id"
        class="module-grid"
        :animation="200"
        ghost-class="card-ghost"
        @end="onDragEnd"
      >
        <template #item="{ element: module }">
          <div
            class="module-card"
            :class="`cat-${module.category}`"
            @click="go(module.path)"
            @mouseenter="onCardEnter($event)"
            @mouseleave="onCardLeave($event)"
          >
            <!-- 收藏按钮 -->
            <span
              class="card-fav"
              :class="{ faved: homeStore.isFavorite(module.id) }"
              @click.stop="homeStore.toggleFavorite(module.id)"
            >
              <StarFilled v-if="homeStore.isFavorite(module.id)" />
              <StarOutlined v-else />
            </span>

            <!-- 图标 + 名称 -->
            <div class="card-top">
              <div class="card-icon-wrap" :class="`icon-${module.theme}`">
                <component :is="getIconComponent(module.icon)" />
              </div>
              <span class="card-name">{{ t(module.nameKey) }}</span>
            </div>

            <!-- 描述 -->
            <p class="card-desc">{{ module.description }}</p>

            <!-- 消息角标 -->
            <span v-if="module.badge && module.badge > 0" class="card-badge" :class="`badge-${module.category}`">
              {{ module.badge > 99 ? '99+' : module.badge }}
            </span>
          </div>
        </template>
      </draggable>

      <!-- 空状态 -->
      <div v-if="homeStore.displayedModules.length === 0" class="empty-state">
        <SearchOutlined class="empty-icon" />
        <p>{{ t('home.noResults') }}</p>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import draggable from 'vuedraggable'
import {
  DashboardOutlined, CheckCircleOutlined, ToolOutlined,
  SafetyOutlined, ApartmentOutlined, CalendarOutlined,
  BranchesOutlined, FileTextOutlined, BugOutlined,
  FilePdfOutlined, SearchOutlined, StarOutlined, StarFilled,
  DatabaseOutlined, SecurityScanOutlined,
  CodeOutlined, ThunderboltOutlined, GlobalOutlined,
  FileProtectOutlined, AppstoreOutlined, LaptopOutlined,
  CheckOutlined,
} from '@ant-design/icons-vue'
import { useHomeStore, type CategoryId } from '@/stores/home'
import { useTodoStatsStore } from '@/stores/todoStats'
import { GetCPUInfo, GetMemoryInfo, GetHostInfo } from '../../bindings/github.com/Aliuyanfeng/happytools/backend/services/monitor/sysinfoservice'
import { TodoService } from '../../bindings/github.com/Aliuyanfeng/happytools/backend/services/todo'
import type { Todo } from '../../bindings/github.com/Aliuyanfeng/happytools/backend/services/todo/models'

const { t } = useI18n()
const router = useRouter()
const homeStore = useHomeStore()
const todoStats = useTodoStatsStore()

// ── 图标映射 ──
const iconMap: Record<string, any> = {
  DashboardOutlined, CheckCircleOutlined, ToolOutlined,
  SafetyOutlined, ApartmentOutlined, CalendarOutlined,
  BranchesOutlined, FileTextOutlined, BugOutlined,
  FilePdfOutlined,
}
function getIconComponent(n: string) { return iconMap[n] || DashboardOutlined }
function go(path: string) { router.push(path) }

// ── 分类导航 ──
const categories: { id: CategoryId; icon: any }[] = [
  { id: 'all', icon: AppstoreOutlined },
  { id: 'security', icon: SecurityScanOutlined },
  { id: 'dev', icon: CodeOutlined },
  { id: 'efficiency', icon: ThunderboltOutlined },
  { id: 'network', icon: GlobalOutlined },
  { id: 'doc', icon: FileProtectOutlined },
]

// ── 拖拽排序 ──
const draggableModules = computed({
  get() {
    return homeStore.displayedModules
  },
  set(val: any[]) {
    homeStore.updateSort(val.map(m => m.id))
  },
})
function onDragEnd() {
  // 排序已通过 computed setter 持久化
}

// ── 主机信息 ──
const hostName = ref('--')
const hostOs = ref('--')

async function refreshHostInfo() {
  try {
    const info = await GetHostInfo()
    if (info) {
      hostName.value = info.hostname || '--'
      hostOs.value = info.platform || '--'
    }
  } catch { /* 静默 */ }
}

// ── 系统资源 ──
const cpuPercent = ref(0)
const memPercent = ref(0)
let statsTimer: ReturnType<typeof setInterval> | null = null

async function refreshStats() {
  try {
    const [cpu, mem] = await Promise.all([GetCPUInfo(), GetMemoryInfo()])
    if (cpu && cpu.core_usages?.length) {
      cpuPercent.value = Math.round(cpu.core_usages.reduce((a, b) => a + b, 0) / cpu.core_usages.length)
    }
    if (mem) {
      memPercent.value = Math.round(mem.used_percent)
    }
  } catch { /* 静默 */ }
}

const cpuClass = computed(() => cpuPercent.value > 80 ? 'stat-warn' : '')
const memClass = computed(() => memPercent.value > 80 ? 'stat-warn' : '')

// ── 待办气泡 ──
const allTodos = ref<Todo[]>([])

async function refreshTodos() {
  try {
    const all = await TodoService.GetAll()
    allTodos.value = all ?? []
    todoStats.pendingCount = allTodos.value.filter((t: any) => !t.completed).length
  } catch { /* 静默 */ }
}

async function toggleTodoInPreview(todo: Todo) {
  try {
    await TodoService.Toggle(todo.id)
    todo.completed = !todo.completed
    todoStats.pendingCount = allTodos.value.filter((t: any) => !t.completed).length
  } catch { /* 静默 */ }
}

// ── 时钟 ──
const now = ref(new Date())
let clockTimer: ReturnType<typeof setInterval>
onMounted(() => {
  clockTimer = setInterval(() => { now.value = new Date() }, 1000)
  refreshHostInfo()
  refreshStats()
  statsTimer = setInterval(refreshStats, 5000)
  refreshTodos()
})
onUnmounted(() => {
  clearInterval(clockTimer)
  if (statsTimer) clearInterval(statsTimer)
})

const pad = (n: number) => n.toString().padStart(2, '0')
const currentTime = computed(() =>
  `${pad(now.value.getHours())}:${pad(now.value.getMinutes())}:${pad(now.value.getSeconds())}`
)

// ── 卡片 hover 扫光动画 ──
function onCardEnter(e: MouseEvent) {
  const card = (e.currentTarget as HTMLElement)
  card.classList.remove('sweep-done')
  card.classList.add('sweep-active')
}
function onCardLeave(e: MouseEvent) {
  const card = (e.currentTarget as HTMLElement)
  card.classList.remove('sweep-active')
  card.classList.add('sweep-done')
}
</script>

<style scoped>
/* ── 根容器 ── */
.home {
  width: 100%;
  height: 100%;
  background: #f8fafc;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

/* ══════════════════════════════════════
   顶部标题栏
══════════════════════════════════════ */
.home-header {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 24px;
  background: #ffffff;
  border-bottom: 1px solid #e2e8f0;
  gap: 16px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

/* 主机信息 */
.host-info {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: #334155;
  background: #f1f5f9;
  padding: 4px 12px;
  border-radius: 8px;
  white-space: nowrap;
}

.host-icon {
  font-size: 14px;
  color: #6366f1;
  opacity: 0.8;
}

.host-name {
  font-weight: 600;
  color: #1e293b;
}

.host-sep {
  color: #cbd5e1;
  margin: 0 1px;
}

.host-os {
  color: #64748b;
  font-size: 12px;
}

/* 系统资源极简预览 */
.sys-stats {
  display: flex;
  align-items: center;
  gap: 12px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #64748b;
  padding: 2px 8px;
  border-radius: 6px;
  background: #f1f5f9;
  transition: color 0.2s, background 0.2s;
}

.stat-item.stat-warn {
  color: #ef4444;
  background: #fef2f2;
}

.stat-icon {
  font-size: 12px;
  opacity: 0.7;
}

.stat-val {
  font-weight: 600;
  font-variant-numeric: tabular-nums;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

/* 全局搜索 */
.search-box {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 4px 10px;
  border-radius: 8px;
  border: 1px solid #e2e8f0;
  background: #f8fafc;
  transition: border-color 0.2s, box-shadow 0.2s;
}

.search-box:focus-within {
  border-color: #6366f1;
  box-shadow: 0 0 0 2px rgba(99, 102, 241, 0.1);
}

.search-icon {
  font-size: 13px;
  color: #94a3b8;
}

.search-input {
  border: none;
  outline: none;
  background: transparent;
  font-size: 13px;
  color: #1e293b;
  width: 140px;
}

.search-input::placeholder {
  color: #cbd5e1;
}

.search-clear {
  font-size: 14px;
  color: #94a3b8;
  cursor: pointer;
  line-height: 1;
  padding: 0 2px;
}

.search-clear:hover {
  color: #64748b;
}

/* 待办气泡 */
.todo-bubble {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: 8px;
  background: #f1f5f9;
  cursor: pointer;
  font-size: 16px;
  color: #64748b;
  transition: background 0.2s, color 0.2s;
}

.todo-bubble:hover {
  background: #e2e8f0;
  color: #1e293b;
}

.todo-badge {
  position: absolute;
  top: -4px;
  right: -4px;
  min-width: 16px;
  height: 16px;
  padding: 0 4px;
  border-radius: 8px;
  background: #ef4444;
  color: #fff;
  font-size: 10px;
  font-weight: 600;
  display: flex;
  align-items: center;
  justify-content: center;
  line-height: 1;
}

/* 待办预览面板 */
.todo-preview-panel {
  display: flex;
  flex-direction: column;
}

.todo-preview-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-bottom: 8px;
  border-bottom: 1px solid #f1f5f9;
  margin-bottom: 4px;
}

.todo-preview-title {
  font-size: 13px;
  font-weight: 600;
  color: #1e293b;
}

.todo-preview-count {
  font-size: 11px;
  color: #fff;
  background: #6366f1;
  padding: 1px 7px;
  border-radius: 8px;
  font-weight: 600;
}

.todo-preview-list {
  max-height: 260px;
  overflow-y: auto;
  scrollbar-width: thin;
  scrollbar-color: #e2e8f0 transparent;
}

.todo-preview-list::-webkit-scrollbar {
  width: 4px;
}

.todo-preview-list::-webkit-scrollbar-thumb {
  background: #e2e8f0;
  border-radius: 2px;
}

.todo-preview-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 7px 4px;
  border-bottom: 1px solid #f8fafc;
  font-size: 13px;
  color: #334155;
  transition: background 0.15s;
  border-radius: 4px;
}

.todo-preview-item:hover {
  background: #f8fafc;
}

.todo-preview-item.is-done {
  opacity: 0.5;
}

.todo-preview-item.is-done .todo-preview-title-text {
  text-decoration: line-through;
}

.todo-check {
  width: 16px;
  height: 16px;
  border-radius: 50%;
  border: 1.5px solid #d1d5db;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  flex-shrink: 0;
  transition: all 0.2s;
  font-size: 10px;
  color: transparent;
}

.todo-check:hover {
  border-color: #6366f1;
}

.todo-check.checked {
  background: #6366f1;
  border-color: #6366f1;
  color: #fff;
}

.todo-preview-title-text {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  min-width: 0;
}

.todo-due {
  font-size: 11px;
  color: #94a3b8;
  flex-shrink: 0;
  white-space: nowrap;
}

.todo-due.overdue {
  color: #ef4444;
}

.todo-due.warning {
  color: #f59e0b;
}

.todo-preview-footer {
  padding-top: 8px;
  border-top: 1px solid #f1f5f9;
  margin-top: 4px;
  text-align: center;
}

.todo-view-all {
  font-size: 12px;
  color: #6366f1;
  cursor: pointer;
  transition: color 0.2s;
}

.todo-view-all:hover {
  color: #4f46e5;
  text-decoration: underline;
}

.todo-empty {
  font-size: 13px;
  color: #94a3b8;
  text-align: center;
  padding: 16px 0;
}

.header-clock {
  font-size: 13px;
  font-weight: 600;
  color: #64748b;
  font-variant-numeric: tabular-nums;
  letter-spacing: 1px;
}

/* ══════════════════════════════════════
   分类导航
══════════════════════════════════════ */
.category-nav {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 8px 24px;
  background: #ffffff;
  border-bottom: 1px solid #f1f5f9;
  overflow-x: auto;
}

.cat-btn {
  display: flex;
  align-items: center;
  gap: 5px;
  padding: 5px 14px;
  border-radius: 8px;
  border: 1px solid transparent;
  background: transparent;
  font-size: 13px;
  color: #64748b;
  cursor: pointer;
  transition: all 0.2s;
  white-space: nowrap;
}

.cat-btn:hover {
  background: #f1f5f9;
  color: #334155;
}

.cat-btn.active {
  background: #eef2ff;
  color: #4f46e5;
  border-color: #c7d2fe;
  font-weight: 500;
}

.cat-icon {
  font-size: 14px;
}

/* ══════════════════════════════════════
   模块卡片网格
══════════════════════════════════════ */
.home-main {
  flex: 1;
  overflow-y: auto;
  padding: 20px 24px;
}

.module-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 12px;
}

/* 卡片入场 */
@keyframes card-in {
  from { opacity: 0; transform: translateY(8px); }
  to   { opacity: 1; transform: none; }
}

/* 模块卡片 */
.module-card {
  position: relative;
  display: flex;
  flex-direction: column;
  background: #ffffff;
  border-radius: 12px;
  border: 1px solid #e2e8f0;
  padding: 16px;
  cursor: pointer;
  overflow: hidden;
  transition: transform 0.25s cubic-bezier(0.23, 1, 0.32, 1),
              box-shadow 0.25s ease,
              border-color 0.25s ease;
  animation: card-in 0.3s ease both;
}

.module-card:hover {
  transform: translateY(-3px) scale(1.01);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.06);
  border-color: #cbd5e1;
}

/* 卡片顶部装饰线 */
.module-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 12px;
  right: 12px;
  height: 2px;
  border-radius: 0 0 2px 2px;
  opacity: 0.5;
  transition: opacity 0.25s;
}

.module-card:hover::before {
  opacity: 1;
}

/* 分类主色 - VT/POC 红色、Git/工具盒紫色 */
.cat-security::before { background: #ef4444; }
.cat-dev::before      { background: #6366f1; }
.cat-efficiency::before { background: #22c55e; }
.cat-network::before  { background: #f97316; }
.cat-doc::before      { background: #06b6d4; }

/* ── 扫光动画 ── */
@keyframes sweep-light {
  0%   { left: -60%; opacity: 0; }
  10%  { opacity: 0.7; }
  90%  { opacity: 0.7; }
  100% { left: 110%; opacity: 0; }
}

.module-card::after {
  content: '';
  position: absolute;
  top: 0;
  left: -60%;
  width: 40%;
  height: 2px;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.7), transparent);
  border-radius: 0 0 2px 2px;
  pointer-events: none;
  opacity: 0;
}

.module-card.sweep-active::after {
  animation: sweep-light 0.6s ease-out forwards;
}

.module-card.sweep-done::after {
  opacity: 0;
  animation: none;
}

/* 收藏按钮 */
.card-fav {
  position: absolute;
  top: 8px;
  right: 8px;
  font-size: 14px;
  color: #cbd5e1;
  cursor: pointer;
  transition: color 0.2s, transform 0.2s;
  z-index: 2;
}

.card-fav:hover {
  transform: scale(1.2);
}

.card-fav.faved {
  color: #f59e0b;
}

/* 卡片内容 */
.card-top {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 8px;
}

.card-icon-wrap {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 17px;
  flex-shrink: 0;
}

.icon-blue   { background: rgba(99, 102, 241, 0.1);  color: #4f46e5; }
.icon-green  { background: rgba(34, 197, 94, 0.1);   color: #16a34a; }
.icon-purple { background: rgba(168, 85, 247, 0.1);  color: #9333ea; }
.icon-orange { background: rgba(249, 115, 22, 0.1);  color: #ea580c; }
.icon-red    { background: rgba(239, 68, 68, 0.1);   color: #dc2626; }
.icon-cyan   { background: rgba(6, 182, 212, 0.1);   color: #0891b2; }

.card-name {
  font-size: 14px;
  font-weight: 600;
  color: #1e293b;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.card-desc {
  margin: 0;
  font-size: 12px;
  color: #94a3b8;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

/* 消息角标 - 按业务域统一色彩 */
.card-badge {
  position: absolute;
  top: 8px;
  left: 8px;
  min-width: 18px;
  height: 18px;
  padding: 0 5px;
  border-radius: 9px;
  font-size: 10px;
  font-weight: 600;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  line-height: 1;
}

.badge-security   { background: #ef4444; }
.badge-dev        { background: #6366f1; }
.badge-efficiency { background: #94a3b8; }
.badge-network    { background: #f97316; }
.badge-doc        { background: #06b6d4; }

/* 拖拽占位 */
.card-ghost {
  opacity: 0.4;
  border-style: dashed;
}

/* 空状态 */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 0;
  color: #94a3b8;
}

.empty-icon {
  font-size: 32px;
  margin-bottom: 8px;
  opacity: 0.5;
}

.empty-state p {
  font-size: 14px;
  margin: 0;
}
</style>