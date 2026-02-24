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
import { computed, onMounted, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { HomeOutlined } from '@ant-design/icons-vue';
import TitleBar from './components/TitleBar.vue';
import { useAppStore } from './stores/app';
import {Events} from "@wailsio/runtime";


const router = useRouter();
const route = useRoute();
const appStore = useAppStore();

const showBackHome = computed(() => route.path !== '/');

const lastUsedTime = ref<string>("")

function goHome() {
  if (route.path !== '/') {
    router.push('/');
  }
}
onMounted(()=>{
  Events.On('app:lastUsedTime', (event) => {
    console.log(event)
    if (event && event.data) {
      console.log(event.data)
      lastUsedTime.value = event.data
      appStore.updateLastUsedTime(event.data)
    }
  });
})



</script>

<template>
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
        <!-- <div class="runtime-info">版本: v0.0.1</div> -->
        <div class="copyright">上次使用: {{ appStore.lastUsedTime }}</div>
      </a-layout-footer>
    </a-layout>
  </a-layout>
</template>

<style lang="scss">
/* 全局样式 - 移除滚动条 */
html, body, #app {
  width: 100%;
  height: 100%;
  margin: 0;
  padding: 0;
  overflow: hidden; /* 禁止全局滚动条 */
}

.layout {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden; /* 禁止布局滚动条 */

  /* 主内容区域 */
  .ant-layout {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: hidden;
    
    .ant-layout-content {
      flex: 1;
      overflow: hidden; /* 禁止内容区域滚动条 */
      display: flex;
      flex-direction: column;
    }
  }

  /* 底部状态栏 */
  .app-footer {
    display: flex;
    justify-content: right;
    padding: 0 16px;
    height: 28px;
    line-height: 28px;
    flex-shrink: 0; /* 防止底部状态栏被压缩 */
    // background-color: #f0f0f0;
    font-size: 14px;
    color: #666;
    font-style: italic;
  }
}

.site-layout .site-layout-background {
  background: #fff;
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