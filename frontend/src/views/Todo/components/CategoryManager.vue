<template>
  <div class="category-manager">
    <!-- 分类列表 -->
    <div class="category-list">
      <span 
        v-for="category in categories" 
        :key="category.id"
        class="category-item"
        :style="{ color: category.color }"
      >
        {{ category.name }}
        <a-button type="text" size="small" @click="deleteCategory(category.id)" style="margin-left: 4px;">
          <DeleteOutlined />
        </a-button>
      </span>
    </div>
    
    <!-- 添加分类按钮 -->
    <a-button type="dashed" @click="showAddModal" style="margin-top: 8px;">
      <PlusOutlined /> 添加分类
    </a-button>
    
    <!-- 添加/编辑弹窗 -->
    <a-modal
      v-model:open="modalVisible"
      :title="editingCategory ? '编辑分类' : '添加分类'"
      @ok="saveCategory"
      @cancel="cancelEdit"
    >
      <a-form :model="form" layout="vertical">
        <a-form-item label="分类名称" :rules="[{ required: true, message: '请输入分类名称' }]">
          <a-input v-model:value="form.name" placeholder="输入分类名称" />
        </a-form-item>
        <a-form-item label="颜色">
          <a-color-picker v-model:value="form.color" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed } from 'vue'
import { PlusOutlined, DeleteOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { useCategoryStore } from '../../../stores/category'

// 使用 Pinia store
const categoryStore = useCategoryStore()

const categories = computed(() => categoryStore.categories)
const modalVisible = ref(false)
const editingCategory = ref<any | null>(null)

const form = ref({
  name: '',
  color: '#1890ff'
})

// 显示添加弹窗
function showAddModal() {
  editingCategory.value = null
  form.value = { name: '', color: '#1890ff' }
  modalVisible.value = true
}

// 显示编辑弹窗
function showEditModal(category: any) {
  editingCategory.value = category
  form.value = { name: category.name, color: category.color }
  modalVisible.value = true
}

// 保存分类
async function saveCategory() {
  if (!form.value.name.trim()) {
    message.warning('请输入分类名称')
    return
  }

  try {
    if (editingCategory.value) {
      // 更新分类
      await categoryStore.updateCategory(editingCategory.value.id, form.value.name, form.value.color)
      message.success('更新成功')
    } else {
      // 添加分类
      await categoryStore.addCategory(form.value.name, form.value.color)
      message.success('添加成功')
    }
    
    modalVisible.value = false
  } catch (e: any) {
    // 错误信息已经在 store 中处理
  }
}

// 取消编辑
function cancelEdit() {
  modalVisible.value = false
  editingCategory.value = null
  form.value = { name: '', color: '#1890ff' }
}

// 删除分类
async function deleteCategory(id: number) {
  try {
    await categoryStore.deleteCategory(id)
    message.success('删除成功')
  } catch (e: any) {
    // 错误信息已经在 store 中处理
  }
}

// 暴露给父组件的方法
defineExpose({
  loadCategories: categoryStore.loadCategories
})
</script>

<style scoped>
.category-manager {
  margin-bottom: 16px;
}
.category-list {
  min-height: 40px;
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}
.category-item {
  display: inline-flex;
  align-items: center;
  font-weight: 500;
}
</style>