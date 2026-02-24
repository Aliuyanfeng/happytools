<template>
  <a-list-item 
    class="todo-item"
    :class="{
      'overdue': todo.status === 2,
      'warning': todo.status === 1
    }"
  >
    <div class="todo-content">
      <!-- 标题和完成状态 -->
      <a-checkbox
        v-model:checked="localCompleted"
        @change="handleToggle"
      >
        <span :class="{ completed: localCompleted }">{{ todo.title }}</span>
      </a-checkbox>

      <!-- 优先级指示器 -->
      <span
        v-if="todo.priority > 0"
        :style="{ color: priorityColor }"
        class="priority-text"
      >
        {{ priorityText }}
      </span>

      <!-- 截止日期和状态 -->
      <span v-if="todo.due_date" class="due-date">
        <ClockCircleOutlined />
        {{ todo.due_date }}
        <span
          v-if="todo.status === 2"
          class="status-text error"
        >
          (已逾期)
        </span>
        <span
          v-else-if="todo.status === 1"
          class="status-text warning"
        >
          (即将到期)
        </span>
      </span>

      <!-- 分类名称 -->
      <span v-if="category" :style="{ color: category.color }" class="category-text">
        {{ category.name }}
      </span>
    </div>
    
    <!-- 操作按钮 -->
    <template #actions>
      <a-button type="text" size="small" @click="$emit('edit', todo)">
        <EditOutlined />
      </a-button>
      <a-popconfirm 
        title="确定删除吗？"
        ok-text="删除"
        cancel-text="取消"
        @confirm="$emit('delete', todo.id)">
        <a-button type="text" size="small" danger>
          <DeleteOutlined />
        </a-button>
      </a-popconfirm>
    </template>
  </a-list-item>
</template>

<script lang="ts" setup>
import { ref, computed, watch } from 'vue'
import { ClockCircleOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons-vue'

const props = defineProps<{
  todo: any
  category?: any
}>()

const emit = defineEmits<{
  (e: 'toggle', id: number): void
  (e: 'edit', todo: any): void
  (e: 'delete', id: number): void
}>()

// 本地状态，避免直接修改 props
const localCompleted = ref(props.todo.completed)

// 监听 props 变化同步到本地状态
watch(() => props.todo.completed, (newVal) => {
  localCompleted.value = newVal
})

const priorityColor = computed(() => {
  return props.todo.priority === 2 ? 'red' : 'orange'
})

const priorityText = computed(() => {
  return props.todo.priority === 2 ? '高' : '中'
})

// 处理复选框变化
function handleToggle() {
  emit('toggle', props.todo.id)
}
</script>

<style scoped>
.todo-item {
  transition: all 0.2s;
}

.todo-item.overdue {
  border-left: 3px solid #ff4d4f;
  background-color: #fff1f0;
}

.todo-item.warning {
  border-left: 3px solid #faad14;
  background-color: #fffbe6;
}

.todo-content {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: 1;
}

.category-text {
  font-weight: 500;
  font-size: 12px;
}

.priority-text {
  font-size: 12px;
  font-weight: 500;
}

.due-date {
  color: #999;
  font-size: 12px;
  display: flex;
  align-items: center;
  gap: 4px;
}

.status-text {
  margin-left: 8px;
  font-weight: 500;
}

.status-text.error {
  color: #ff4d4f;
}

.status-text.warning {
  color: #faad14;
}

.completed {
  text-decoration: line-through;
  color: #999;
}
</style>