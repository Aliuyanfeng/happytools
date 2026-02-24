<template>
  <a-select
    v-model:value="selectedCategoryId"
    placeholder="选择分类"
    allow-clear
    style="width: 120px"
    @change="handleChange"
  >
    <a-select-option :value="null">无分类</a-select-option>
    <a-select-option 
      v-for="category in categories" 
      :key="category.id" 
      :value="category.id"
    >
      <span :style="{ color: category.color }">{{ category.name }}</span>
    </a-select-option>
  </a-select>
</template>

<script lang="ts" setup>
import { ref, watch, computed } from 'vue'
import { useCategoryStore } from '../../../stores/category'

const selectedCategoryId = ref<number | null>(null)

// 使用 Pinia store
const categoryStore = useCategoryStore()

// 计算属性获取分类列表
const categories = computed(() => categoryStore.categories)

// 监听选择变化
function handleChange(value: number | null) {
  selectedCategoryId.value = value
  emit('update:modelValue', value)
}

// 定义 props 和 emits
const props = defineProps<{
  modelValue: number | null
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: number | null): void
}>()

// 同步外部值
watch(() => props.modelValue, (newValue) => {
  selectedCategoryId.value = newValue
}, { immediate: true })
</script>