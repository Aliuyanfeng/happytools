<!--
 * @Author: LiuYanFeng
 * @Date: 2025-07-03 17:16:49
 * @LastEditors: LiuYanFeng
 * @LastEditTime: 2025-07-29 14:40:20
 * @FilePath: \happytools\frontend\src\App.vue
 * @Description: 像珍惜礼物一样珍惜今天
 * 
 * Copyright (c) 2025 by ${git_name_email}, All Rights Reserved. 
-->
<script lang="ts" setup>
import {onUnmounted, ref} from 'vue';
import {GithubOutlined, SettingOutlined} from '@ant-design/icons-vue';
import {useRouter} from 'vue-router';

const router = useRouter();
console.log(router.getRoutes())
const selectedKeys = ref<string[]>(['0']);
const collapsed = ref<boolean>(false);
const currentTime = ref<string>(new Date().toLocaleTimeString());

// Update time every second
const updateTime = () => {
  currentTime.value = new Date().toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  });
};

const timeInterval = setInterval(updateTime, 1000);

// Cleanup interval on unmount
onUnmounted(() => {
  clearInterval(timeInterval);
});

// Button click handlers
const handleSettingsClick = () => {
  // Will implement settings panel opening logic
};

const handleGithubClick = () => {
  window.open('https://github.com/Aliuyanfeng/happytools', '_blank');
};

</script>

<template>
  <a-layout class="layout">
    <a-layout>
      <a-layout-header class="header">
        <div class="header-content">
          <div class="logo flex items-center justify-center">
            <img src="@/assets/images/logo.png" alt="logo">
          </div>
          <a-menu mode="horizontal" :selected-keys="selectedKeys" @click="selectedKeys = [$event.key]">
            <a-menu-item :key="index" v-for="(item,index) in router.getRoutes()">
              <router-link :to="item.path">{{item.meta?.title}}</router-link>
            </a-menu-item>
          </a-menu>
        </div>
        <div class="header-actions">
          <a-button type="text" shape="circle" class="flex justify-center items-center" @click="handleSettingsClick">
            <SettingOutlined class="text-xl leading-4"/>
          </a-button>
          <a-button type="text" shape="circle" class="flex justify-center items-center" @click="handleGithubClick">
            <GithubOutlined class="text-xl leading-4"/>
          </a-button>
        </div>
      </a-layout-header>
      <a-layout-content class="content">
        <router-view/>
      </a-layout-content>
      <a-layout-footer class="app-footer">
        <div class="runtime-info">版本: v0.0.1</div>
        <div class="copyright">© 2025 Aliu - Happy Tools. All Rights Reserved.</div>
      </a-layout-footer>
    </a-layout>
  </a-layout>
</template>

<style lang="scss">
.layout {
  width: 100%;
  height: 100%;

  .header {
    background-color: #fff;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 20px;
    height: 46px;
    line-height: 46px;

    .logo {
      width: 100px;
      margin-right: 24px;
    }

    .header-content {
      display: flex;
      align-items: center;
    }

    .header-actions {
      display: flex;
      align-items: center;
      //gap: 16px;
    }

    .current-time {
      color: #666;
      font-size: 14px;
    }
  }

  /* 底部状态栏 */

  .app-footer {
    display: flex;
    justify-content: space-between;
    padding: 0 16px;
    height: 28px;
    line-height: 28px;
    background-color: #f0f0f0;
    font-size: 14px;
    color: #666;
    font-style: italic;
  }
}

.site-layout .site-layout-background {
  background: #fff;
}
</style>