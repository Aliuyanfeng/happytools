import { defineStore } from 'pinia'
import { ref } from 'vue'
import { TodoService } from '../../bindings/github.com/Aliuyanfeng/happytools/backend/services/todo'

export const useTodoStatsStore = defineStore('todoStats', () => {
  const pendingCount = ref(0)

  async function refresh() {
    try {
      const todos = await TodoService.GetAll()
      pendingCount.value = (todos ?? []).filter((t: any) => !t.completed).length
    } catch { /* 静默失败 */ }
  }

  return { pendingCount, refresh }
})
