<template>
  <div class="pdf-parser-container">
    <!-- 左侧导航栏 -->
    <div class="pdf-parser-sidebar" :class="{ collapsed: isCollapsed }">
      <div class="sidebar-header">
        <div v-if="!isCollapsed" class="header-content">
          <h2 class="mb-1 text-lg font-bold text-gray-800">{{ t('pdfParser.title') }}</h2>
          <p class="text-xs text-gray-500">{{ t('pdfParser.subtitle') }}</p>
        </div>
        <div v-else class="header-collapsed">
          <FilePdfOutlined class="text-xl" />
        </div>
      </div>

      <a-menu
        v-model:selectedKeys="selectedKeys"
        mode="inline"
        class="pdf-parser-menu"
        :inline-collapsed="isCollapsed"
        @click="handleMenuClick"
      >
        <!-- 元数据分类 -->
        <a-menu-item-group key="metadata-group" :title="t('pdfParser.categoryMetadata')">
          <a-menu-item key="metadata">
            <template #icon><InfoCircleOutlined /></template>
            <span v-if="!isCollapsed">{{ t('pdfParser.menuMetadata') }}</span>
          </a-menu-item>
        </a-menu-item-group>
      </a-menu>

      <div class="collapse-btn" @click="toggleCollapse">
        <MenuFoldOutlined v-if="!isCollapsed" />
        <MenuUnfoldOutlined v-else />
      </div>
    </div>

    <!-- 右侧内容区 -->
    <div class="pdf-parser-content">
      <div v-if="!currentTool" class="welcome-page">
        <div class="welcome-icon">
          <FilePdfOutlined class="text-6xl text-blue-400" />
        </div>
        <h2 class="mt-4 text-xl font-semibold text-gray-700">{{ t('pdfParser.welcome') }}</h2>
        <p class="mt-2 text-sm text-gray-400">{{ t('pdfParser.welcomeHint') }}</p>
      </div>
      <router-view v-else />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import {
  FilePdfOutlined,
  InfoCircleOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
} from '@ant-design/icons-vue'

const { t } = useI18n()
const router = useRouter()
const route = useRoute()
const selectedKeys = ref<string[]>([])
const isCollapsed = ref(false)

const currentTool = computed(() => {
  const path = route.path
  if (path.includes('metadata')) return 'metadata'
  return ''
})

watch(currentTool, (newVal) => {
  if (newVal) {
    selectedKeys.value = [newVal]
  }
}, { immediate: true })

const handleMenuClick = ({ key }: { key: string }) => {
  router.push(`/pdf-parser/${key}`)
}

const toggleCollapse = () => {
  isCollapsed.value = !isCollapsed.value
}
</script>

<style scoped>
.pdf-parser-container {
  display: flex;
  height: 100%;
  background: #f5f7fa;
}

.pdf-parser-sidebar {
  width: 220px;
  min-width: 220px;
  background: #fff;
  border-right: 1px solid #e8e8e8;
  display: flex;
  flex-direction: column;
  transition: all 0.3s;
  overflow: hidden;
}

.pdf-parser-sidebar.collapsed {
  width: 64px;
  min-width: 64px;
}

.sidebar-header {
  padding: 16px;
  border-bottom: 1px solid #f0f0f0;
  min-height: 64px;
  display: flex;
  align-items: center;
}

.header-content {
  white-space: nowrap;
  overflow: hidden;
}

.header-collapsed {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
}

.pdf-parser-menu {
  flex: 1;
  border-right: none;
}

.collapse-btn {
  padding: 12px 16px;
  border-top: 1px solid #f0f0f0;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #666;
  transition: color 0.3s;
}

.collapse-btn:hover {
  color: #1890ff;
}

.pdf-parser-content {
  flex: 1;
  overflow: auto;
  padding: 20px;
}

.welcome-page {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  min-height: 400px;
}

.welcome-icon {
  opacity: 0.6;
}
</style>