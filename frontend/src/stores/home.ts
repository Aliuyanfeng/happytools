/*
 * @Author: LiuYanFeng
 * @Date: 2026-07-21
 * @Description: 首页状态管理 - 收藏、排序、分类筛选
 */
import { defineStore } from 'pinia'
import { ref, computed, watch } from 'vue'
import { modules, type ModuleConfig } from '@/config/modules'

export type CategoryId = 'all' | 'security' | 'dev' | 'efficiency' | 'network' | 'doc'

export const useHomeStore = defineStore('home', () => {
  // ── 收藏 ──
  const favorites = ref<string[]>(
    JSON.parse(localStorage.getItem('homeFavorites') || '[]')
  )

  function toggleFavorite(id: string) {
    const idx = favorites.value.indexOf(id)
    if (idx === -1) {
      favorites.value.push(id)
    } else {
      favorites.value.splice(idx, 1)
    }
  }

  function isFavorite(id: string): boolean {
    return favorites.value.includes(id)
  }

  // ── 排序 ──
  const sortOrder = ref<string[]>(
    JSON.parse(localStorage.getItem('homeSortOrder') || '[]')
  )

  /** 获取排序后的模块列表 */
  const sortedModules = computed(() => {
    const order = sortOrder.value.length ? sortOrder.value : modules.map(m => m.id)
    const indexed = new Map(order.map((id, i) => [id, i]))
    return [...modules].sort((a, b) => {
      const ai = indexed.get(a.id) ?? modules.length
      const bi = indexed.get(b.id) ?? modules.length
      return ai - bi
    })
  })

  /** 拖拽排序后持久化 */
  function updateSort(ids: string[]) {
    sortOrder.value = ids
  }

  // ── 分类筛选 ──
  const activeCategory = ref<CategoryId>(
    (localStorage.getItem('homeCategory') as CategoryId) || 'all'
  )

  function setCategory(cat: CategoryId) {
    activeCategory.value = cat
  }

  /** 按分类过滤模块 */
  const filteredModules = computed(() => {
    let list = sortedModules.value
    if (activeCategory.value !== 'all') {
      list = list.filter(m => m.category === activeCategory.value)
    }
    return list
  })

  // ── 搜索 ──
  const searchQuery = ref('')

  /** 搜索 + 分类过滤 */
  const displayedModules = computed(() => {
    let list = filteredModules.value
    if (searchQuery.value.trim()) {
      const q = searchQuery.value.trim().toLowerCase()
      list = list.filter(m =>
        m.nameKey.toLowerCase().includes(q) ||
        m.description.toLowerCase().includes(q) ||
        m.id.toLowerCase().includes(q)
      )
    }
    // 排除隐藏模块
    return list
  })

  // ── 持久化 ──
  watch(favorites, (v) => {
    localStorage.setItem('homeFavorites', JSON.stringify(v))
  }, { deep: true })

  watch(sortOrder, (v) => {
    localStorage.setItem('homeSortOrder', JSON.stringify(v))
  }, { deep: true })

  watch(activeCategory, (v) => {
    localStorage.setItem('homeCategory', v)
  })

  return {
    favorites,
    toggleFavorite,
    isFavorite,
    sortOrder,
    sortedModules,
    updateSort,
    activeCategory,
    setCategory,
    filteredModules,
    searchQuery,
    displayedModules,
  }
})