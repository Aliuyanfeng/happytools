<template>
  <a-modal
    v-model:open="visible"
    title="全局设置"
    width="700px"
    :footer="null"
    :bodyStyle="{ padding: 0 }"
    @cancel="handleClose"
  >
    <div class="settings-container">
      <!-- 左侧 Tab 栏 -->
      <div class="settings-sidebar">
        <div
          class="sidebar-item"
          :class="{ active: activeTab === 'basic' }"
          @click="activeTab = 'basic'"
        >
          <SettingOutlined class="sidebar-icon" />
          <span>基本设置</span>
        </div>
        <!-- 未来可以添加更多设置项 -->
        <!-- <div
          class="sidebar-item"
          :class="{ active: activeTab === 'advanced' }"
          @click="activeTab = 'advanced'"
        >
          <ToolOutlined class="sidebar-icon" />
          <span>高级设置</span>
        </div> -->
      </div>

      <!-- 右侧内容区 -->
      <div class="settings-content">
        <!-- 基本设置 -->
        <div v-if="activeTab === 'basic'" class="settings-panel">
          <a-form layout="vertical">
            <!-- 主题设置 -->
            <a-form-item label="主题模式">
              <a-radio-group v-model:value="settingsStore.themeMode" button-style="solid">
                <a-radio-button value="light">
                  <BulbOutlined class="mr-1" />
                  浅色
                </a-radio-button>
                <a-radio-button value="dark">
                  <BulbOutlined class="mr-1" />
                  深色
                </a-radio-button>
                <a-radio-button value="auto">
                  <SyncOutlined class="mr-1" />
                  自动
                </a-radio-button>
              </a-radio-group>
              <div class="setting-hint">
                <InfoCircleOutlined class="mr-1" />
                自动模式将跟随系统主题设置
              </div>
            </a-form-item>

            <!-- 字体大小 -->
            <a-form-item label="字体大小">
              <a-radio-group v-model:value="settingsStore.fontSize" button-style="solid">
                <a-radio-button value="small">小</a-radio-button>
                <a-radio-button value="medium">中</a-radio-button>
                <a-radio-button value="large">大</a-radio-button>
              </a-radio-group>
            </a-form-item>

            <!-- 自定义字体 -->
            <a-form-item label="自定义字体">
              <a-select
                v-model:value="settingsStore.customFont"
                placeholder="选择字体"
                allowClear
                style="width: 100%"
              >
                <a-select-option value="">默认字体</a-select-option>
                <a-select-option value="Microsoft YaHei">微软雅黑</a-select-option>
                <a-select-option value="SimSun">宋体</a-select-option>
                <a-select-option value="SimHei">黑体</a-select-option>
                <a-select-option value="KaiTi">楷体</a-select-option>
                <a-select-option value="FangSong">仿宋</a-select-option>
                <a-select-option value="Arial">Arial</a-select-option>
                <a-select-option value="Times New Roman">Times New Roman</a-select-option>
              </a-select>
              <div class="setting-hint">
                <InfoCircleOutlined class="mr-1" />
                选择系统已安装的字体
              </div>
            </a-form-item>
          </a-form>
        </div>

        <!-- 高级设置（预留） -->
        <!-- <div v-if="activeTab === 'advanced'" class="settings-panel">
          <a-empty description="高级设置开发中..." />
        </div> -->
      </div>
    </div>
  </a-modal>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import {
  SettingOutlined,
  BulbOutlined,
  SyncOutlined,
  InfoCircleOutlined
} from '@ant-design/icons-vue'
import { useSettingsStore } from '../stores/settings'

const props = defineProps<{
  open: boolean
}>()

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
}>()

const settingsStore = useSettingsStore()
const activeTab = ref('basic')

const visible = ref(props.open)

watch(() => props.open, (newVal) => {
  visible.value = newVal
})

watch(visible, (newVal) => {
  emit('update:open', newVal)
})

function handleClose() {
  visible.value = false
}
</script>

<style scoped>
.settings-container {
  display: flex;
  height: 450px;
}

.settings-sidebar {
  width: 180px;
  background: #fafafa;
  border-right: 1px solid #e8e8e8;
  padding: 16px 0;
  flex-shrink: 0;
}

.sidebar-item {
  display: flex;
  align-items: center;
  padding: 12px 20px;
  cursor: pointer;
  transition: all 0.3s;
  color: #595959;
  font-size: 14px;
}

.sidebar-item:hover {
  background: #e6f7ff;
  color: #1890ff;
}

.sidebar-item.active {
  background: #e6f7ff;
  color: #1890ff;
  border-right: 3px solid #1890ff;
  font-weight: 600;
}

.sidebar-icon {
  margin-right: 8px;
  font-size: 16px;
}

.settings-content {
  flex: 1;
  padding: 24px;
  overflow-y: auto;
}

.settings-panel {
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.setting-hint {
  margin-top: 8px;
  padding: 8px 12px;
  background: #f0f5ff;
  border: 1px solid #d6e4ff;
  border-radius: 4px;
  font-size: 12px;
  color: #2f54eb;
  display: flex;
  align-items: center;
}

.mr-1 {
  margin-right: 4px;
}

/* 深色主题样式 */
[data-theme="dark"] .settings-sidebar {
  background: #1f1f1f;
  border-right-color: #434343;
}

[data-theme="dark"] .sidebar-item {
  color: #a6a6a6;
}

[data-theme="dark"] .sidebar-item:hover {
  background: #2a2a2a;
  color: #40a9ff;
}

[data-theme="dark"] .sidebar-item.active {
  background: #2a2a2a;
  color: #40a9ff;
  border-right-color: #40a9ff;
}

[data-theme="dark"] .setting-hint {
  background: #1a2a3a;
  border-color: #1e3a5f;
  color: #69b1ff;
}
</style>
