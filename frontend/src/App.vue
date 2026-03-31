<!--
 * @Author: LiuYanFeng
 * @Date: 2025-07-03 17:16:49
 * @LastEditors: LiuYanFeng
 * @LastEditTime: 2026-02-06 14:29:25
 * @FilePath: \happytools\frontend\src\App.vue
 * @Description: 像珍惜礼物一样珍惜今天
 * 
 * Copyright (c) 2025 by ${git_name_email}, All Rights Reserved. 
-->
<script lang="ts" setup>
import { computed, onMounted, ref, h } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import { HomeOutlined } from '@ant-design/icons-vue';
import { theme as antTheme, notification, Button } from 'ant-design-vue';
import TitleBar from './components/TitleBar.vue';
import { useAppStore } from './stores/app';
import { useSettingsStore } from './stores/settings';
import {Events} from "@wailsio/runtime";
import { UpdateService } from '../bindings/github.com/Aliuyanfeng/happytools/backend/services/update/index'

const { t } = useI18n();
const router = useRouter();
const route = useRoute();
const appStore = useAppStore();
const settingsStore = useSettingsStore();

const showBackHome = computed(() => route.path !== '/');

const lastUsedTime = ref<string>("")

// 动态 antd 主题配置
const antdTheme = computed(() => {
  const fontSizeMap = { small: 12, medium: 14, large: 16 }
  const fontSize = fontSizeMap[settingsStore.fontSize] ?? 14
  const fontFamily = settingsStore.customFont
    ? `${settingsStore.customFont}, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif`
    : undefined
  return {
    algorithm: settingsStore.isDark ? antTheme.darkAlgorithm : antTheme.defaultAlgorithm,
    token: {
      fontSize,
      ...(fontFamily ? { fontFamily } : {}),
    },
  }
})

function goHome() {
  if (route.path !== '/') {
    router.push('/');
  }
}
onMounted(()=>{
  // 初始化设置
  settingsStore.initSettings()
  
  Events.On('app:lastUsedTime', (event) => {
    console.log(event)
    if (event && event.data) {
      console.log(event.data)
      lastUsedTime.value = event.data
      appStore.updateLastUsedTime(event.data)
    }
  });

  // 静默检查更新（延迟 3 秒，等前端完全加载）
  setTimeout(async () => {
    try {
      const result = await UpdateService.CheckUpdate()
      if (!result?.hasUpdate) return

      const currentVersion = await UpdateService.GetCurrentVersion()

      // 检查是否已忽略该版本（dev 模式不受忽略限制）
      const ignoredVersion = localStorage.getItem('ignoredUpdateVersion')
      if (ignoredVersion === result.latest && currentVersion !== 'dev') return

      notification.info({
        key: 'update-notify',
        message: `🎉 发现新版本 v${result.latest}`,
        description: `当前版本 ${currentVersion === 'dev' ? 'dev' : 'v' + currentVersion}，新版本已发布，建议更新。`,
        btn: () => h('div', { style: 'display:flex;gap:8px' }, [
          h(Button, {
            type: 'primary',
            size: 'small',
            onClick: () => {
              import('@wailsio/runtime').then(({ Browser }) => {
                Browser.OpenURL(result.releaseUrl)
              })
              notification.close('update-notify')
            }
          }, '前往下载'),
          h(Button, {
            size: 'small',
            onClick: () => {
              localStorage.setItem('ignoredUpdateVersion', result.latest)
              notification.close('update-notify')
            }
          }, '忽略此版本'),
        ]),
        duration: 0,
        placement: 'bottomRight',
      })
    } catch {
      // 静默失败
    }
  }, 3000)
})



</script>

<template>
  <a-config-provider :theme="antdTheme">
    <a-layout class="layout">
      <TitleBar />
      <a-layout>
        <a-layout-content class="content">
          <router-view />
          <transition name="fade">
            <a-button
              v-if="showBackHome"
              class="back-home-btn"
              type="primary"
              shape="circle"
              size="large"
              @click="goHome"
            >
              <HomeOutlined />
            </a-button>
          </transition>
        </a-layout-content>
        <a-layout-footer class="app-footer">
          <div class="copyright">{{ t('app.lastUsed') }}: {{ appStore.lastUsedTime }}</div>
        </a-layout-footer>
      </a-layout>
    </a-layout>
  </a-config-provider>
</template>

<style lang="scss">
/* 全局样式 - 移除滚动条 */
html, body, #app {
  width: 100%;
  height: 100%;
  margin: 0;
  padding: 0;
  overflow: hidden;
}

/* CSS 变量定义 - 浅色默认 */
:root {
  --bg-primary: #f5f7fa;
  --bg-secondary: #ffffff;
  --text-primary: #262626;
  --text-secondary: #595959;
  --border-color: #e8e8e8;
}

/* 深色主题 CSS 变量 */
[data-theme="dark"] {
  --bg-primary: #141414;
  --bg-secondary: #1f1f1f;
  --text-primary: #ffffff;
  --text-secondary: #a6a6a6;
  --border-color: #434343;

  /* 修正 body/html 背景，防止白边 */
  background-color: var(--bg-primary);
  color: var(--text-primary);

  /* 底部状态栏 */
  .app-footer {
    color: #a6a6a6;
    border-top-color: var(--border-color);
  }

  /* Toolbox 侧边栏等自定义组件 */
  .toolbox-container {
    background: var(--bg-primary);
    border-color: var(--border-color);
  }

  .toolbox-sidebar {
    background: var(--bg-secondary);
    border-right-color: var(--border-color);
  }

  .toolbox-content {
    background: var(--bg-primary);
  }

  .sidebar-header {
    border-bottom-color: var(--border-color);
    h2 { color: var(--text-primary) !important; }
    p { color: var(--text-secondary) !important; }
  }

  /* 欢迎页 */
  .welcome-page h1 { color: var(--text-primary) !important; }
  .welcome-page p { color: var(--text-secondary) !important; }

  /* 拖拽区域 */
  .drop-zone, .select-area {
    border-color: var(--border-color);
    &:hover { border-color: #1890ff; background: rgba(24, 144, 255, 0.06); }
  }

  /* 字体颜色修正 */
  .text-gray-800 { color: var(--text-primary) !important; }
  .text-gray-700 { color: #d9d9d9 !important; }
  .text-gray-600 { color: #bfbfbf !important; }
  .text-gray-500 { color: #a6a6a6 !important; }
  .text-gray-400 { color: #8c8c8c !important; }
}

.layout {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;

  .ant-layout {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: hidden;

    .ant-layout-content {
      flex: 1;
      overflow: hidden;
      display: flex;
      flex-direction: column;
    }
  }

  .app-footer {
    display: flex;
    justify-content: right;
    padding: 0 16px;
    height: 26px;
    line-height: 26px;
    flex-shrink: 0;
    font-size: 11px;
    color: #94a3b8;
    font-style: italic;
    background: rgba(255,255,255,0.8);
    border-top: 1px solid rgba(99,102,241,0.08);
    backdrop-filter: blur(10px);
  }
}

.back-home-btn {
  position: fixed;
  right: 24px;
  bottom: 48px;
  box-shadow: 0 12px 30px rgba(15, 23, 42, 0.35);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease-out, transform 0.2s ease-out;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(8px);
}
</style>