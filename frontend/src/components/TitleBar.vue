<template>
  <div class="titlebar">
    <div class="titlebar-left" @dblclick="toggleMaximize">
      <!-- <div class="app-icon">🔧</div> -->
      <!-- <div class="app-title">happytools</div> -->
      <img class=app-logo src="@/assets/images/logo.png" alt="">
    </div>

    <div class="titlebar-right">
      <button class="title-btn settings" :title="t('settings.title')" @click="openSettings">
        <SettingOutlined class="icon" />
      </button>
      <button class="title-btn pin" :class="{ active: isAlwaysOnTop }" :title="isAlwaysOnTop ? t('app.unpin') : t('app.pin')" @click="toggleAlwaysOnTop">
        <PushpinOutlined class="icon" />
      </button>
      <button class="title-btn minimize" :title="t('app.minimize')" @click="minimize"><span class="icon">−</span></button>
      <!-- <button class="title-btn maximize" :title="t('app.maximize')" @click="toggleMaximize"><span class="icon">▢</span></button> -->
      <button class="title-btn close" :title="t('app.close')" @click="close"><span class="icon">✕</span></button>
    </div>

    <!-- 全局设置弹窗 -->
    <SettingsModal v-model:open="settingsVisible" />
  </div>
</template>

<script lang="ts" setup>
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { Window } from '@wailsio/runtime'
import { SettingOutlined, PushpinOutlined } from '@ant-design/icons-vue'
import SettingsModal from './SettingsModal.vue'
import { useSettingsStore } from '../stores/settings'

const { t } = useI18n()
const settingsStore = useSettingsStore()
const settingsVisible = ref(false)
const isAlwaysOnTop = ref(false)

function minimize() {
  Window.Minimise()
}

function toggleMaximize() {
  Window.ToggleMaximise()
}

function close() {
  // 根据设置决定关闭行为
  if (settingsStore.closeBehavior === 'hide') {
    Window.Hide()
  } else {
    Window.Close()
  }
}

function openSettings() {
  settingsVisible.value = true
}

async function toggleAlwaysOnTop() {
  isAlwaysOnTop.value = !isAlwaysOnTop.value
  await Window.SetAlwaysOnTop(isAlwaysOnTop.value)
}

// 为了兼容你给出的用法（直接使用 document.querySelector），同时在组件挂载时绑定选择器事件
onMounted(() => {
 
})
</script>

<style scoped>
/* Orbitronio - 科技感字体 */
@font-face {
  font-family: 'Orbitronio';
  src: url('@/assets/fonts/Orbitronio-1.ttf') format('truetype');
  font-weight: normal;
  font-style: normal;
  font-display: swap;
}

/* Ethnocentric RG - 常规字体 */
@font-face {
  font-family: 'EthnocentricRG';
  src: url('@/assets/fonts/ethnocentric-rg-2.ttf') format('truetype');
  font-weight: normal;
  font-style: normal;
  font-display: swap;
}

/* Ethnocentric RG Italic - 斜体字体 */
@font-face {
  font-family: 'EthnocentricRGItalic';
  src: url('@/assets/fonts/ethnocentric-rg-it-1.ttf') format('truetype');
  font-weight: normal;
  font-style: italic;
  font-display: swap;
}


.titlebar {
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 8px;
  background: linear-gradient(90deg, #eff1f3, #4285dc);
  color: #323131;
  -webkit-app-region: drag; /* 允许拖拽窗口 */
  user-select: none;
  --wails-draggable: drag;
}

.titlebar-left {
  display: flex;
  align-items: center;
  gap: 8px;
}

.app-icon {
  font-size: 20px;
  
}
.app-logo {
  height: 28px;
}

.app-title {
  /* font-family: 'Orbitronio', sans-serif; */
  /* font-family: 'EthnocentricRG', sans-serif; */
  font-family: 'EthnocentricRGItalic', sans-serif;
  font-weight: 600;
  font-size: 24px;
  letter-spacing: 1px;
}

.titlebar-right {
  display: flex;
  gap: 6px;
}

.title-btn {
  height: 30px;
  min-width: 40px;
  padding: 0;
  border: none;
  border-radius: 6px;
  background: rgba(255, 255, 255, 0.02);
  color: #ffffff;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 700;
  cursor: pointer;
  -webkit-app-region: no-drag; /* 按钮区域不触发拖拽 */
  transition: background 0.12s ease, transform 0.08s ease;
  box-shadow: none;
}

/* 更明显的 hover 效果 */
.title-btn:hover {
  background: rgba(255, 255, 255, 0.08);
  transform: translateY(-1px);
}

/* 设置按钮特殊样式 */
.title-btn.settings:hover {
  background: rgba(24, 144, 255, 0.8);
  color: #ffffff;
}

/* 置顶按钮特殊样式 */
.title-btn.pin:hover {
  background: rgba(250, 173, 20, 0.8);
  color: #ffffff;
}

.title-btn.pin.active {
  background: rgba(250, 173, 20, 0.9);
  color: #ffffff;
}

/* 关闭按钮特殊样式 */
.title-btn.close:hover {
  background: #e74c3c;
  color: #ffffff;
}

/* 更清晰的图标 */
.title-btn .icon {
  font-size: 15px;
  line-height: 1;
  display: inline-block;
}

/* 小屏或触控优化 */
@media (hover: none) {
  .title-btn:hover {
    transform: none;
  }
}
</style>
