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
          :class="{ active: activeTab === 'network' }"
          @click="activeTab = 'network'"
        >
          <WifiOutlined class="sidebar-icon" />
          <span>{{ t('settings.network') }}</span>
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

        <!-- 网络设置 -->
        <div v-if="activeTab === 'network'" class="settings-panel">
          <a-form layout="vertical">
            <a-divider>{{ t('settings.dnsCache') }}</a-divider>

            <a-form-item :label="t('settings.dnsCache')">
              <div class="dns-section">
                <p class="dns-description">{{ t('settings.dnsCacheDesc') }}</p>
                <a-button 
                  type="primary" 
                  :loading="dnsFlushing" 
                  @click="flushDNS"
                >
                  {{ dnsFlushing ? t('settings.dnsFlushing') : t('settings.flushDNS') }}
                </a-button>
                <div v-if="dnsResult" class="dns-result" :class="{ 'dns-success': dnsResult.success, 'dns-error': !dnsResult.success }">
                  <p>{{ dnsResult.message }}</p>
                  <p v-if="dnsResult.adminNeeded" class="dns-warning">{{ t('settings.dnsAdminNeeded') }}</p>
                  <pre v-if="dnsResult.output">{{ dnsResult.output }}</pre>
                </div>
              </div>
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
  EyeInvisibleOutlined,
  WifiOutlined
} from '@ant-design/icons-vue'
import { useSettingsStore } from '../stores/settings'
import { VTService } from '../../bindings/github.com/Aliuyanfeng/happytools/backend/services/vt'
import { DNSService } from '../../bindings/github.com/Aliuyanfeng/happytools/backend/services/network/index'

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

// DNS相关状态
const dnsFlushing = ref(false)
const dnsResult = ref<{
  success: boolean;
  message: string;
  output: string;
  adminNeeded: boolean;
} | null>(null)

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

// 刷新DNS缓存
async function flushDNS() {
  dnsFlushing.value = true
  dnsResult.value = null
  
  try {
    const result = await DNSService.FlushDNS()
    dnsResult.value = result
    
    if (result && result.success) {
      message.success(t('settings.dnsFlushSuccess'))
    } else if (result) {
      message.error(result.message)
    }
  } catch (error) {
    console.error('刷新DNS失败:', error)
    dnsResult.value = {
      success: false,
      message: t('settings.dnsFlushFailed'),
      output: String(error),
      adminNeeded: false
    }
    message.error(t('settings.dnsFlushFailed'))
  } finally {
    dnsFlushing.value = false
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

/* DNS 设置样式 */
.dns-section {
  padding: 16px;
  background: #fafafa;
  border-radius: 8px;
}

.dns-description {
  color: #666;
  margin-bottom: 12px;
  font-size: 13px;
}

.dns-result {
  margin-top: 12px;
  padding: 12px;
  border-radius: 4px;
  font-size: 13px;
}

.dns-success {
  background-color: #f6ffed;
  border: 1px solid #b7eb8f;
}

.dns-error {
  background-color: #fff2f0;
  border: 1px solid #ffccc7;
}

.dns-warning {
  color: #fa8c16;
  font-weight: 500;
}

.dns-result pre {
  margin-top: 8px;
  padding: 8px;
  background-color: rgba(0, 0, 0, 0.05);
  border-radius: 4px;
  white-space: pre-wrap;
  word-break: break-all;
  font-size: 12px;
  max-height: 150px;
  overflow-y: auto;
}

[data-theme="dark"] .dns-section {
  background: #1f1f1f;
}

[data-theme="dark"] .dns-description {
  color: #a6a6a6;
}

[data-theme="dark"] .dns-success {
  background-color: #162312;
  border-color: #274916;
}

[data-theme="dark"] .dns-error {
  background-color: #2a1215;
  border-color: #58181c;
}

[data-theme="dark"] .dns-result pre {
  background-color: rgba(255, 255, 255, 0.05);
}
</style>
