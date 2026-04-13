<template>
  <div
    class="todo-item"
    :class="{
      'is-overdue': todo.status === 2,
      'is-warning': todo.status === 1,
      'is-done': localCompleted
    }"
  >
    <!-- 优先级色条 -->
    <div v-if="todo.priority > 0 && !localCompleted" class="priority-bar" :class="todo.priority === 2 ? 'high' : 'mid'" />

    <!-- 勾选框 -->
    <div class="check-wrap" @click="handleToggle">
      <div class="check-circle" :class="{ checked: localCompleted }">
        <CheckOutlined v-if="localCompleted" class="check-icon" />
      </div>
    </div>

    <!-- 内容 -->
    <div class="item-body">
      <span class="item-title">{{ todo.title }}</span>
      <div class="item-meta" v-if="todo.due_date || category">
        <span v-if="category" class="meta-cat" :style="{ color: category.color, background: category.color + '18' }">
          {{ category.name }}
        </span>
        <span v-if="todo.due_date" class="meta-date" :class="{ overdue: todo.status === 2, warning: todo.status === 1 }">
          <ClockCircleOutlined />
          {{ todo.due_date }}
          <span v-if="todo.status === 2"> · 逾期</span>
          <span v-else-if="todo.status === 1"> · 即将到期</span>
        </span>
      </div>
    </div>

    <!-- 操作 -->
    <div class="item-actions">
      <button class="act-btn" @click="$emit('edit', todo)" title="编辑">
        <EditOutlined />
      </button>
      <a-popconfirm title="确定删除吗？" ok-text="删除" cancel-text="取消" @confirm="$emit('delete', todo.id)">
        <button class="act-btn danger" title="删除">
          <DeleteOutlined />
        </button>
      </a-popconfirm>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, watch } from 'vue'
import { CheckOutlined, ClockCircleOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons-vue'

const props = defineProps<{ todo: any; category?: any }>()
const emit = defineEmits<{
  (e: 'toggle', id: number): void
  (e: 'edit', todo: any): void
  (e: 'delete', id: number): void
}>()

const localCompleted = ref(props.todo.completed)
watch(() => props.todo.completed, v => { localCompleted.value = v })

function handleToggle() { emit('toggle', props.todo.id) }
</script>

<style scoped>
.todo-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 11px 18px;
  border-bottom: 1px solid #f7f7f7;
  transition: background 0.12s;
  position: relative;
  overflow: hidden;
}

.todo-item:last-child { border-bottom: none; }
.todo-item:hover { background: #fafafa; }

.todo-item.is-overdue { background: #fff9f9; }
.todo-item.is-warning { background: #fffdf5; }

/* 优先级色条 */
.priority-bar {
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 3px;
}
.priority-bar.high { background: #ff4d4f; }
.priority-bar.mid  { background: #faad14; }

/* 自定义勾选 */
.check-wrap {
  flex-shrink: 0;
  cursor: pointer;
}

.check-circle {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  border: 2px solid #d9d9d9;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.check-circle:hover {
  border-color: #1677ff;
}

.check-circle.checked {
  background: #52c41a;
  border-color: #52c41a;
}

.check-icon {
  font-size: 11px;
  color: #fff;
}

/* 内容 */
.item-body {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.item-title {
  font-size: 14px;
  color: #1a1a1a;
  line-height: 1.5;
  transition: color 0.2s;
}

.is-done .item-title {
  text-decoration: line-through;
  color: #c0c0c0;
}

.item-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.meta-cat {
  font-size: 11px;
  font-weight: 500;
  padding: 1px 7px;
  border-radius: 4px;
}

.meta-date {
  font-size: 11px;
  color: #bbb;
  display: flex;
  align-items: center;
  gap: 3px;
}

.meta-date.overdue { color: #ff4d4f; }
.meta-date.warning { color: #faad14; }

/* 操作按钮 */
.item-actions {
  display: flex;
  gap: 2px;
  opacity: 0;
  transition: opacity 0.15s;
  flex-shrink: 0;
}

.todo-item:hover .item-actions { opacity: 1; }

.act-btn {
  width: 28px;
  height: 28px;
  border: none;
  background: transparent;
  border-radius: 6px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #bbb;
  font-size: 13px;
  transition: all 0.15s;
}

.act-btn:hover { background: #f0f0f0; color: #555; }
.act-btn.danger:hover { background: #fff1f0; color: #ff4d4f; }
</style>
