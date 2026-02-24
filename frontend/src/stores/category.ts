/*
 * @Author: LiuYanFeng
 * @Date: 2026-02-12
 * @Description: 分类状态管理
 */
import { defineStore } from 'pinia'
import { ref } from 'vue'
import { message } from 'ant-design-vue'
import { CategoryService, Category } from '../../bindings/github.com/Aliuyanfeng/happytools/backend/services/category'

export const useCategoryStore = defineStore('category', () => {
  // 分类列表
  const categories = ref<Category[]>([])
  
  // 加载状态
  const loading = ref(false)

  // 加载所有分类
  async function loadCategories() {
    loading.value = true
    try {
      const result = await CategoryService.GetAll()
      categories.value = result || []
    } catch (e: any) {
      message.error('加载分类失败: ' + e.message)
      categories.value = []
    } finally {
      loading.value = false
    }
  }

  // 添加分类
  async function addCategory(name: string, color: string) {
    try {
      const category = await CategoryService.Add(name, color)
      if (category) {
        categories.value.push(category)
        return category
      }
    } catch (e: any) {
      message.error('添加分类失败: ' + e.message)
      throw e
    }
  }

  // 更新分类
  async function updateCategory(id: number, name: string, color: string) {
    try {
      await CategoryService.Update(id, name, color)
      const index = categories.value.findIndex(c => c.id === id)
      if (index !== -1) {
        categories.value[index] = {
          ...categories.value[index],
          name,
          color
        }
      }
    } catch (e: any) {
      message.error('更新分类失败: ' + e.message)
      throw e
    }
  }

  // 删除分类
  async function deleteCategory(id: number) {
    try {
      await CategoryService.Delete(id)
      categories.value = categories.value.filter(c => c.id !== id)
    } catch (e: any) {
      message.error('删除分类失败: ' + e.message)
      throw e
    }
  }

  // 根据ID获取分类
  function getCategoryById(id: number | null) {
    if (id === null) return null
    return categories.value.find(c => c.id === id)
  }

  // 获取分类的待办数量（需要传入待办列表）
  function getCategoryCount(categoryId: number, todos: any[]) {
    return todos.filter(t => t.category_id === categoryId).length
  }

  // 初始化时加载分类
  loadCategories()

  return {
    categories,
    loading,
    loadCategories,
    addCategory,
    updateCategory,
    deleteCategory,
    getCategoryById,
    getCategoryCount
  }
})