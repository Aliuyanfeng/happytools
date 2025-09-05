<!--
 * @Author: LiuYanFeng
 * @Date: 2025-07-04 17:59:06
 * @LastEditors: LiuYanFeng
 * @LastEditTime: 2025-07-07 14:55:23
 * @FilePath: \happytools\frontend\src\components\ToolCard.vue
 * @Description: 像珍惜礼物一样珍惜今天
 * 
 * Copyright (c) 2025 by ${git_name_email}, All Rights Reserved. 
-->
<template>
  <div class="tool-card" @click="handleClick">
    <div class="tool-icon">
      <a-icon :type="tool.icon" :style="{ fontSize: '24px' }" />
    </div>
    <div class="tool-name">{{ tool.name }}</div>
    <div class="tool-metrics">
      <span v-if="tool.recent" class="recent-badge">最近使用</span>
      <span class="usage-count">使用: {{ tool.usage }}次</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { defineProps, defineEmits } from 'vue';
import { FileSearchOutlined } from '@ant-design/icons-vue';

const props = defineProps<{
  tool: {
    id: number;
    name: string;
    icon: string;
    path: string;
    usage: number;
    recent: boolean;
  };
}>();

const emit = defineEmits<{
  (e: 'click'): void;
}>();

const handleClick = () => {
  emit('click');
};
</script>

<style scoped>
.tool-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 16px;
  background-color: #ffffff;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  cursor: pointer;
  transition: transform 0.2s, box-shadow 0.2s;
}

.tool-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.tool-icon {
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #3498db;
  margin-bottom: 12px;
}

.tool-name {
  font-size: 14px;
  font-weight: 500;
  margin-bottom: 8px;
  text-align: center;
  color: #333;
}

.tool-metrics {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
  font-size: 12px;
  color: #666;
}

.recent-badge {
  background-color: #e3f2fd;
  color: #1976d2;
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 11px;
}

.usage-count {
  color: #888;
}
</style>