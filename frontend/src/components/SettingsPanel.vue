<!--
 * @Author: LiuYanFeng
 * @Date: 2025-07-04 17:59:58
 * @LastEditors: LiuYanFeng
 * @LastEditTime: 2025-07-07 15:34:10
 * @FilePath: \happytools\frontend\src\components\SettingsPanel.vue
 * @Description: 像珍惜礼物一样珍惜今天
 * 
 * Copyright (c) 2025 by ${git_name_email}, All Rights Reserved. 
-->
<template>
  <a-form layout="vertical" :model="settings" @change="handleSettingsChange">
    <a-form-item label="字体设置">
      <a-select v-model:value="settings.fontFamily" placeholder="选择字体">
        <a-select-option value="system">系统默认</a-select-option>
        <a-select-option value="sans-serif">无衬线字体</a-select-option>
        <a-select-option value="serif">衬线字体</a-select-option>
        <a-select-option value="monospace">等宽字体</a-select-option>
      </a-select>
    </a-form-item>

    <a-form-item label="字体大小">
      <a-slider
        v-model:value="settings.fontSize"
        :min="12"
        :max="20"
        :step="1"
        :marks="{ 12: '12px', 14: '14px', 16: '16px', 18: '18px', 20: '20px' }"
      />
      <div style="text-align: center; margin-top: 10px;">{{ settings.fontSize }}px</div>
    </a-form-item>

    <a-form-item label="文件缓存目录">
      <a-input v-model:value="settings.cacheDir" placeholder="缓存目录路径" />
      <a-button type="primary" style="margin-top: 10px;" @click="selectCacheDir">浏览...</a-button>
    </a-form-item>

    <a-divider>关于我们</a-divider>
    <div class="about-section">
      <p>Happy Tools v1.0.0</p>
      <p>© 2023 Happy Tools. 保留所有权利。</p>
    </div>
  </a-form>
</template>

<script setup lang="ts">
import { watch } from 'vue';
import { message } from 'ant-design-vue';

const props = defineProps<{
  settings: {
    fontFamily: string;
    fontSize: number;
    cacheDir: string;
  };
}>();

const emit = defineEmits<{
  (e: 'update:settings', value: typeof props.settings): void;
}>();

const localSettings = {...props.settings};

const handleSettingsChange = () => {
  emit('update:settings', localSettings);
};

const selectCacheDir = () => {
  // 在实际实现中，这里应该调用文件选择对话框
  message.info('文件选择对话框将在完整实现中打开');
  // 模拟选择目录
  localSettings.cacheDir = 'C:\\happytools\\cache';
  handleSettingsChange();
};

// 监听props变化以更新本地设置
watch(
  () => props.settings,
  (newVal) => {
    Object.assign(localSettings, newVal);
  },
  { deep: true }
);
</script>

<style scoped>
.about-section {
  padding: 16px;
  text-align: center;
  color: #666;
  background-color: #f9f9f9;
  border-radius: 8px;
}

.about-section p {
  margin: 6px 0;
}
</style>