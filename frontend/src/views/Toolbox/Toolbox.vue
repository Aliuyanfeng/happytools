<template>
  <div class="toolbox-container">
    <!-- 左侧导航栏 -->
    <div class="toolbox-sidebar" :class="{ collapsed: isCollapsed }">
      <div class="sidebar-header">
        <div v-if="!isCollapsed" class="header-content">
          <h2 class="mb-1 text-lg font-bold text-gray-800">工具盒子</h2>
          <p class="text-xs text-gray-500">实用工具集合</p>
        </div>
        <div v-else class="header-collapsed">
          <ToolOutlined class="text-xl" />
        </div>
      </div>

      <a-menu
        v-model:selectedKeys="selectedKeys"
        mode="inline"
        class="toolbox-menu"
        :inline-collapsed="isCollapsed"
        @click="handleMenuClick"
      >
        <!-- 数据转换分类 -->
        <a-menu-item-group key="data-convert" title="数据转换">
          <a-menu-item key="unit-converter">
            <template #icon>
              <RetweetOutlined />
            </template>
            <span v-if="!isCollapsed">单位转换</span>
          </a-menu-item>
        </a-menu-item-group>

        <!-- 辅助工具分类 -->
        <a-menu-item-group key="tools" title="辅助工具">
          <a-menu-item key="encryption">
            <template #icon>
              <SafetyOutlined />
            </template>
            <span v-if="!isCollapsed">加密工具</span>
          </a-menu-item>
          <a-menu-item key="png-injector">
            <template #icon>
              <FileImageOutlined />
            </template>
            <span v-if="!isCollapsed">PNG 注入</span>
          </a-menu-item>
        </a-menu-item-group>
      </a-menu>

      <!-- 收缩按钮 -->
      <div class="collapse-btn" @click="toggleCollapse">
        <MenuFoldOutlined v-if="!isCollapsed" />
        <MenuUnfoldOutlined v-else />
      </div>
    </div>

    <!-- 右侧内容区 -->
    <div class="toolbox-content">
      <!-- 欢迎页面 -->
      <div v-if="!currentTool" class="welcome-page">
        <div class="welcome-content">
          <ToolOutlined class="welcome-icon" />
          <h1 class="mb-4 text-3xl font-bold text-gray-800">欢迎使用工具盒子</h1>
          <p class="mb-8 text-lg text-gray-600">请从左侧选择一个工具开始使用</p>

          <!-- <div class="quick-access">
            <h3 class="mb-4 text-lg font-semibold text-gray-800">快速访问</h3>
            <a-row :gutter="[16, 16]">
              <a-col :span="12">
                <a-card hoverable class="quick-card" @click="navigateTo('unit-converter')">
                  <template #cover>
                    <div class="card-cover blue">
                      <RetweetOutlined class="cover-icon" />
                    </div>
                  </template>
                  <a-card-meta title="单位转换" description="字节、长度、时间转换" />
                </a-card>
              </a-col>
              <a-col :span="12">
                <a-card hoverable class="quick-card" @click="navigateTo('encryption')">
                  <template #cover>
                    <div class="card-cover green">
                      <SafetyOutlined class="cover-icon" />
                    </div>
                  </template>
                  <a-card-meta title="加密工具" description="MD5、SHA、Base64编解码" />
                </a-card>
              </a-col>
            </a-row>
          </div> -->
        </div>
      </div>

      <!-- 子路由内容 -->
      <router-view v-else />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import {
  RetweetOutlined,
  SafetyOutlined,
  ToolOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  FileImageOutlined
} from '@ant-design/icons-vue'

const router = useRouter()
const route = useRoute()

const selectedKeys = ref<string[]>([])
const isCollapsed = ref(false)

// 当前选中的工具
const currentTool = computed(() => {
  const path = route.path
  if (path.includes('unit-converter')) return 'unit-converter'
  if (path.includes('encryption')) return 'encryption'
  if (path.includes('png-injector')) return 'png-injector'
  return ''
})

// 监听路由变化，更新选中状态
watch(currentTool, (newVal) => {
  if (newVal) {
    selectedKeys.value = [newVal]
  } else {
    selectedKeys.value = []
  }
}, { immediate: true })

// 菜单点击处理
const handleMenuClick = ({ key }: { key: string }) => {
  navigateTo(key)
}

// 导航到指定工具
const navigateTo = (tool: string) => {
  router.push(`/toolbox/${tool}`)
}

// 切换收缩状态
const toggleCollapse = () => {
  isCollapsed.value = !isCollapsed.value
}
</script>

<style scoped>
.toolbox-container {
  display: flex;
  height: 100%;
  background: #f5f7fa;
  border-radius: 8px;
  overflow: hidden;
  border: 1px solid #e8e8e8;
}

/* 左侧导航栏 */
.toolbox-sidebar {
  width: 180px;
  background: #ffffff;
  border-right: 1px solid #e8e8e8;
  display: flex;
  flex-direction: column;
  transition: width 0.3s ease;
  position: relative;
}

.toolbox-sidebar.collapsed {
  width: 64px;
}

.sidebar-header {
  padding: 20px 16px;
  border-bottom: 1px solid #f0f0f0;
  min-height: 72px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.header-content {
  text-align: center;
}

.header-collapsed {
  display: flex;
  align-items: center;
  justify-content: center;
  color: #1890ff;
}

.toolbox-menu {
  flex: 1;
  border: none;
  background: transparent;
  padding: 8px;
}

.toolbox-menu :deep(.ant-menu-item-group-title) {
  color: #8c8c8c;
  font-size: 12px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  padding: 12px 8px 8px;
}

.toolbox-menu :deep(.ant-menu-item) {
  margin: 4px 0;
  border-radius: 6px;
  color: #595959;
  height: 36px;
  line-height: 36px;
}

.toolbox-menu :deep(.ant-menu-item:hover) {
  background: #e6f7ff;
  color: #1890ff;
}

.toolbox-menu :deep(.ant-menu-item-selected) {
  background: #1890ff;
  color: white;
}

/* 收缩按钮 */
.collapse-btn {
  position: absolute;
  bottom: 16px;
  right: -12px;
  width: 24px;
  height: 24px;
  background: #1890ff;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  cursor: pointer;
  box-shadow: 0 2px 8px rgba(24, 144, 255, 0.3);
  transition: all 0.3s;
  z-index: 10;
}

.collapse-btn:hover {
  background: #40a9ff;
  transform: scale(1.1);
}

/* 右侧内容区 */
.toolbox-content {
  flex: 1;
  overflow-y: auto;
  background: #f5f7fa;
}

/* 欢迎页面 */
.welcome-page {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px;
}

.welcome-content {
  text-align: center;
  max-width: 800px;
}

.welcome-icon {
  font-size: 80px;
  color: #1890ff;
  margin-bottom: 24px;
}

.quick-access {
  margin-top: 48px;
  text-align: left;
}

.quick-card {
  background: #ffffff;
  border: 1px solid #e8e8e8;
  border-radius: 8px;
  transition: all 0.3s;
  overflow: hidden;
}

.quick-card:hover {
  border-color: #1890ff;
  transform: translateY(-4px);
  box-shadow: 0 8px 16px rgba(24, 144, 255, 0.15);
}

.card-cover {
  height: 120px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.card-cover.blue {
  background: linear-gradient(135deg, #1890ff 0%, #096dd9 100%);
}

.card-cover.green {
  background: linear-gradient(135deg, #52c41a 0%, #389e0d 100%);
}

.cover-icon {
  font-size: 48px;
  color: white;
}

.quick-card :deep(.ant-card-meta-title) {
  color: #262626;
  font-weight: 600;
}

.quick-card :deep(.ant-card-meta-description) {
  color: #8c8c8c;
}

/* 滚动条样式 */
.toolbox-content::-webkit-scrollbar {
  width: 6px;
}

.toolbox-content::-webkit-scrollbar-track {
  background: #f5f7fa;
}

.toolbox-content::-webkit-scrollbar-thumb {
  background: #d9d9d9;
  border-radius: 3px;
}

.toolbox-content::-webkit-scrollbar-thumb:hover {
  background: #bfbfbf;
}
</style>
