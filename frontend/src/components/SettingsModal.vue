<template>
  <a-modal
    v-model:open="visible"
    :title="t('settings.title')"
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
          <span>{{ t('settings.basic') }}</span>
        </div>
        <div
          class="sidebar-item"
          :class="{ active: activeTab === 'virustotal' }"
          @click="activeTab = 'virustotal'"
        >
          <SafetyOutlined class="sidebar-icon" />
          <span>VirusTotal</span>
        </div>
      </div>
  
      <!-- 右侧内容区 -->
      <div class="settings-content">
        <!-- 基本设置 -->
        <div v-if="activeTab === 'basic'" class="settings-panel">
          <a-form layout="vertical">
            <!-- 界面语言 -->
            <a-form-item :label="t('settings.language')">
              <a-radio-group v-model:value="settingsStore.language" button-style="solid">
                <a-radio-button value="auto">
                  <GlobalOutlined class="mr-1" />
                  {{ t('settings.languageFollowSystem') }}
                </a-radio-button>
                <a-radio-button value="zh-CN">
                  中文
                </a-radio-button>
                <a-radio-button value="en-US">
                  English
                </a-radio-button>
              </a-radio-group>
            </a-form-item>

            <!-- 关闭按钮行为 -->
            <a-form-item :label="t('settings.closeBehavior')">
              <a-radio-group v-model:value="settingsStore.closeBehavior" button-style="solid">
                <a-radio-button value="exit">
                  <LogoutOutlined class="mr-1" />
                  {{ t('settings.closeBehaviorExit') }}
                </a-radio-button>
                <a-radio-button value="hide">
                  <EyeInvisibleOutlined class="mr-1" />
                  {{ t('settings.closeBehaviorHide') }}
                </a-radio-button>
              </a-radio-group>
            </a-form-item>

            <!-- 主题设置 -->
            <a-form-item :label="t('settings.theme')">
              <a-radio-group v-model:value="settingsStore.themeMode" button-style="solid">
                <a-radio-button value="light">
                  <BulbOutlined class="mr-1" />
                  {{ t('settings.themeLight') }}
                </a-radio-button>
                <a-radio-button value="dark">
                  <BulbOutlined class="mr-1" />
                  {{ t('settings.themeDark') }}
                </a-radio-button>
                <a-radio-button value="auto">
                  <SyncOutlined class="mr-1" />
                  {{ t('settings.themeAuto') }}
                </a-radio-button>
              </a-radio-group>
              <div class="setting-hint">
                <InfoCircleOutlined class="mr-1" />
                {{ t('settings.themeAutoHint') }}
              </div>
            </a-form-item>
  
            <!-- 字体大小 -->
            <a-form-item :label="t('settings.fontSize')">
              <a-radio-group v-model:value="settingsStore.fontSize" button-style="solid">
                <a-radio-button value="small">{{ t('settings.fontSizeSmall') }}</a-radio-button>
                <a-radio-button value="medium">{{ t('settings.fontSizeMedium') }}</a-radio-button>
                <a-radio-button value="large">{{ t('settings.fontSizeLarge') }}</a-radio-button>
              </a-radio-group>
            </a-form-item>
  
            <!-- 自定义字体 -->
            <a-form-item :label="t('settings.customFont')">
              <a-select
                v-model:value="settingsStore.customFont"
                :placeholder="t('settings.customFont')"
                allowClear
                style="width: 100%"
              >
                <a-select-option value="">{{ t('settings.defaultFont') }}</a-select-option>
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
                {{ t('settings.customFontHint') }}
              </div>
            </a-form-item>
          </a-form>
        </div>
  
        <!-- VirusTotal 设置 -->
        <div v-if="activeTab === 'virustotal'" class="settings-panel">
          <a-form layout="vertical">
            <a-alert
              :message="t('settings.apiKeyNotConfigured')"
              :description="t('settings.apiKeyNotConfiguredDesc')"
              type="info"
              show-icon
              class="mb-4"
            />

            <a-form-item :label="t('settings.apiKey')">
              <a-input-password
                v-model:value="settingsStore.vtApiKey"
                :placeholder="t('settings.apiKeyPlaceholder')"
                style="width: 100%"
                @blur="handleApiKeyBlur"
              />
              <div class="setting-hint">
                <InfoCircleOutlined class="mr-1" />
                {{ t('settings.apiKeyHint') }}
              </div>
            </a-form-item>

            <a-divider>{{ t('settings.concurrency') }}</a-divider>

            <a-form-item :label="t('settings.concurrency')">
              <a-slider
                v-model:value="settingsStore.vtConcurrency"
                :min="1"
                :max="10"
                :marks="{ 1: '1', 3: '3', 5: '5', 7: '7', 10: '10' }"
                @afterChange="handleConcurrencyChange"
              />
              <div class="setting-hint">
                <InfoCircleOutlined class="mr-1" />
                {{ t('settings.concurrencyHint') }}
              </div>
            </a-form-item>

            <a-form-item>
              <a-button type="link" @click="openVirusTotalDocs">
                <LinkOutlined class="mr-1" />
                {{ t('settings.viewApiDocs') }}
              </a-button>
            </a-form-item>
          </a-form>
        </div>
      </div>
    </div>
  </a-modal>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { message } from 'ant-design-vue'
import { useI18n } from 'vue-i18n'
import {
  SettingOutlined,
  BulbOutlined,
  SyncOutlined,
  InfoCircleOutlined,
  SafetyOutlined,
  LinkOutlined,
  GlobalOutlined,
  LogoutOutlined,
  EyeInvisibleOutlined
} from '@ant-design/icons-vue'
import { useSettingsStore } from '../stores/settings'
import { VTService } from '../../bindings/github.com/Aliuyanfeng/happytools/backend/services/vt'

const { t } = useI18n()

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

function openVirusTotalDocs() {
  window.open('https://docs.virustotal.com/reference/overview', '_blank')
}

// API Key 输入框失焦时保存到后端
async function handleApiKeyBlur() {
  const apiKey = settingsStore.vtApiKey
  if (apiKey) {
    try {
      await VTService.SetAPIKey(apiKey)
      message.success(t('settings.apiKeySaved'))
    } catch (error) {
      console.error('Failed to save API Key:', error)
      message.error(t('settings.apiKeySaveFailed'))
    }
  }
}

// 并发数变更时保存到后端
async function handleConcurrencyChange(value: number) {
  try {
    await VTService.SetConcurrency(value)
    message.success(`${t('settings.concurrency')}: ${value}`)
  } catch (error) {
    console.error('Failed to save concurrency:', error)
    message.error(t('common.failed'))
  }
}
</script>

<style scoped>
.settings-container {
  display: flex;
  height: 500px;
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

.mb-4 {
  margin-bottom: 16px;
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
