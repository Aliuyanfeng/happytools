<template>
  <div class="home">
    <div class="home-container">
      <!-- 左侧：功能入口 -->
      <div class="left-panel">
        <div class="panel-header">
          <div class="panel-icon">🚀</div>
          <h2 class="panel-title">{{ t('home.featureEntry') }}</h2>
        </div>

        <div class="menu-list">
          <div
            v-for="module in modules"
            :key="module.id"
            :class="['menu-item', `menu-${module.theme}`]"
            @click="go(module.path)"
          >
            <div class="menu-icon" :style="getIconStyle(module.theme)">
              <component :is="getIconComponent(module.icon)" />
            </div>
            <div class="menu-content">
              <h3 class="menu-title">{{ t(module.nameKey) }}</h3>
            </div>
            <div class="menu-arrow">
              <ArrowRightOutlined />
            </div>
          </div>
        </div>
      </div>

      <!-- 右侧：欢迎信息 -->
      <div class="right-panel">
        <div class="welcome-card">
          <div class="welcome-emoji">👋</div>
          <h1 class="welcome-title">{{ t('home.welcome') }}</h1>
          <h2 class="welcome-subtitle">{{ t('home.subtitle') }}</h2>

          <div class="feature-list">
            <div class="feature-item">
              <div class="feature-icon">⚡</div>
              <div class="feature-content">
                <h4>{{ t('home.lightFast') }}</h4>
              </div>
            </div>
            <div class="feature-item">
              <div class="feature-icon">🔒</div>
              <div class="feature-content">
                <h4>{{ t('home.secure') }}</h4>
              </div>
            </div>
            <div class="feature-item">
              <div class="feature-icon">🎯</div>
              <div class="feature-content">
                <h4>{{ t('home.simple') }}</h4>
              </div>
            </div>
          </div>

          <div class="version-info">
            <span class="version-text">v1.0.0</span>
            <span class="copyright-text">© 2025 Li6 - Happy Tools. Built with Wails3</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import {
  DashboardOutlined,
  CheckCircleOutlined,
  ToolOutlined,
  SafetyOutlined,
  ArrowRightOutlined,
  ApartmentOutlined
} from '@ant-design/icons-vue'
import { modules, themeColors } from '@/config/modules'

const { t } = useI18n()
const router = useRouter()

// 图标组件映射
const iconComponents: Record<string, any> = {
  DashboardOutlined,
  CheckCircleOutlined,
  ToolOutlined,
  SafetyOutlined,
  ApartmentOutlined,
}

// 获取图标组件
function getIconComponent(iconName: string) {
  return iconComponents[iconName] || DashboardOutlined
}

// 获取图标样式
function getIconStyle(theme: string) {
  const colors = themeColors[theme as keyof typeof themeColors] || themeColors.blue
  return {
    background: `linear-gradient(135deg, ${colors.primary}, ${colors.secondary})`
  }
}

function go(path: string) {
  router.push(path)
}
</script>

<style scoped>
.home {
  min-height: calc(100vh - 36px);
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 28px;
}

.home-container {
  width: 100%;
  max-width: 1200px;
  display: grid;
  grid-template-columns: 1.3fr 1fr;
  gap: 40px;
  align-items: stretch;
}

/* 左侧面板 */
.left-panel {
  background: #ffffff;
  border-radius: 24px;
  padding: 32px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.1);
  display: flex;
  flex-direction: column;
}

.panel-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 24px;
  padding-bottom: 20px;
  border-bottom: 2px solid #f0f0f0;
}

.panel-icon {
  font-size: 32px;
}

.panel-title {
  font-size: 20px;
  font-weight: 700;
  color: #1f2937;
  margin: 0;
}

.menu-list {
  flex: 1;
  display: grid;
  grid-template-rows: repeat(5, 1fr);
  grid-auto-flow: column;
  gap: 12px;
  overflow: hidden;
  padding: 0 4px;
}

.menu-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px;
  border-radius: 16px;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
  border: 2px solid transparent;
  background: #f9fafb;
  height: 60px;
  box-sizing: border-box;
}

.menu-item:hover {
  transform: translateY(-8px);
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.1);
}

.menu-icon {
  width: 40px;
  height: 40px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  color: #ffffff;
  flex-shrink: 0;
  transition: all 0.3s ease;
}

.menu-item:hover .menu-icon {
  transform: scale(1.1) rotate(5deg);
}

.menu-content {
  flex: 1;
  min-width: 0;
}

.menu-title {
  font-size: 15px;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 2px 0;
}



.menu-arrow {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background: #e5e7eb;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #6b7280;
  font-size: 12px;
  flex-shrink: 0;
  transition: all 0.3s ease;
}

.menu-item:hover .menu-arrow {
  background: #1f2937;
  color: #ffffff;
  transform: translateX(4px);
}

/* 右侧面板 */
.right-panel {
  display: flex;
  /* align-items: center; */
  justify-content: center;
}

.welcome-card {
  background: #ffffff;
  border-radius: 24px;
  padding: 32px 40px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.1);
  text-align: center;
  max-width: 480px;
  width: 100%;
}

.welcome-emoji {
  font-size: 64px;
  margin-bottom: 24px;
  animation: wave 2s ease-in-out infinite;
}

@keyframes wave {
  0%, 100% {
    transform: rotate(0deg);
  }
  10%, 30% {
    transform: rotate(14deg);
  }
  20%, 40% {
    transform: rotate(-8deg);
  }
  50% {
    transform: rotate(0deg);
  }
}

.welcome-title {
  font-size: 24px;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 8px 0;
}

.welcome-subtitle {
  font-size: 32px;
  font-weight: 700;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin: 0 0 32px 0;
}



.feature-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 24px;
}

.feature-item {
  display: flex;
  align-items: end;
  gap: 12px;
  text-align: left;
}

.feature-icon {
  font-size: 24px;
  flex-shrink: 0;
}

.feature-content h4 {
  font-size: 15px;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 4px 0;
}



.version-info {
  padding-top: 24px;
  border-top: 1px solid #e5e7eb;
  display: flex;
  flex-direction: column;
  gap: 8px;
  align-items: center;
}

.version-text {
  font-size: 13px;
  color: #9ca3af;
  font-weight: 500;
}

.copyright-text {
  font-size: 12px;
  color: #a8a8a8;
  font-weight: 400;
}
</style>
