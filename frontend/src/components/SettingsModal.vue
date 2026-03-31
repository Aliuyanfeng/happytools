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
        <div
          class="sidebar-item"
          :class="{ active: activeTab === 'advanced' }"
          @click="activeTab = 'advanced'"
        >
          <WarningOutlined class="sidebar-icon" />
          <span>{{ t('settings.advanced') }}</span>
        </div>
        <div
          class="sidebar-item"
          :class="{ active: activeTab === 'about' }"
          @click="activeTab = 'about'"
        >
          <InfoCircleOutlined class="sidebar-icon" />
          <span>{{ t('settings.about') }}</span>
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

            <!-- 主题设置 - 已隐藏 -->

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
              v-if="!settingsStore.vtApiKey"
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

            <!-- 账户信息 -->
            <div v-if="userInfoLoading || userInfo || userInfoError" class="user-info-card">
              <div v-if="userInfoLoading" class="user-info-loading">
                <LoadingOutlined class="mr-1" spin />
                {{ t('settings.vtUserInfoLoading') }}
              </div>
              <div v-else-if="userInfoError" class="user-info-error">
                <CloseCircleOutlined class="mr-1" />
                {{ userInfoError }}
              </div>
              <div v-else-if="userInfo">
                <div class="user-info-title">
                  <CheckCircleOutlined class="mr-1 text-green" />
                  {{ t('settings.vtUserInfo') }}
                </div>
                <a-descriptions :column="2" size="small" bordered>
                  <a-descriptions-item :label="t('settings.vtUserId')">
                    <span class="font-mono">{{ userInfo.id }}</span>
                  </a-descriptions-item>
                  <a-descriptions-item :label="t('settings.vtUserType')">
                    <a-tag color="blue">{{ userInfo.type }}</a-tag>
                  </a-descriptions-item>
                  <a-descriptions-item :label="t('settings.vtUserEmail')" :span="2">
                    {{ userInfo.email || '—' }}
                  </a-descriptions-item>
                  <a-descriptions-item :label="t('settings.vtUserLastLogin')" :span="2">
                    {{ userInfo.lastLogin ? formatTimestamp(userInfo.lastLogin) : '—' }}
                  </a-descriptions-item>
                </a-descriptions>
              </div>
            </div>

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

        <!-- 高级设置 -->
        <div v-if="activeTab === 'advanced'" class="settings-panel">

          <!-- 首页模块管理 -->
          <div class="module-manage">
            <div class="mm-title">首页显示模块</div>
            <p class="mm-desc">选择在首页展示的功能入口，取消勾选的模块将被隐藏。</p>
            <div class="mm-list">
              <div v-for="mod in allModules" :key="mod.id" class="mm-item">
                <a-switch
                  :checked="!settingsStore.hiddenModules.includes(mod.id)"
                  size="small"
                  @change="(v: boolean) => toggleModule(mod.id, v)"
                />
                <span class="mm-name">{{ t(mod.nameKey) }}</span>
              </div>
            </div>
          </div>

          <a-divider />

          <div class="danger-zone">
            <div class="danger-zone-header">
              <WarningOutlined class="dz-icon" />
              <span class="dz-title">危险操作区</span>
            </div>
            <p class="dz-desc">以下操作不可撤销，请谨慎操作。</p>

            <div class="danger-card">
              <div class="dc-info">
                <div class="dc-name">清除所有缓存数据</div>
                <div class="dc-desc">
                  将清空以下所有数据：待办事项、分类、日报记录、VirusTotal 任务、Git 仓库配置、Makefile 最近文件及自定义模板。<br />
                  <strong>应用设置（语言、字体等）不会被清除。</strong>
                </div>
              </div>
              <a-button
                danger
                type="primary"
                :loading="clearing"
                @mousedown.prevent="confirmClear"
              >
                <template #icon><DeleteOutlined /></template>
                清除数据
              </a-button>
            </div>
          </div>
        </div>

        <!-- 关于 -->
        <div v-if="activeTab === 'about'" class="settings-panel about-panel">
          <div class="about-logo-wrap">
            <img src="@/assets/images/logo.png" alt="logo" class="about-logo" />
          </div>
          <div class="about-name">HappyTools</div>
          <div class="about-version">
            <span class="about-ver-label">当前版本</span>
            <span class="about-ver-val">{{ currentVersion === 'dev' ? 'dev' : 'v' + currentVersion }}</span>
          </div>

          <a-button
            type="primary"
            :loading="checkingUpdate"
            @click="manualCheckUpdate"
            style="margin-top:4px"
          >
            检查更新
          </a-button>
          <div v-if="manualUpdateResult" class="update-result" :class="manualUpdateResult.hasUpdate ? 'has-update' : 'up-to-date'">
            <template v-if="manualUpdateResult.hasUpdate">
              🎉 发现新版本 v{{ manualUpdateResult.latest }}，
              <a @click="openRelease(manualUpdateResult.releaseUrl)">点击下载</a>
            </template>
            <template v-else>
              ✅ 已是最新版本
            </template>
          </div>

          <a-divider style="margin:16px 0 12px" />

          <div class="about-links">
            <div class="about-link-item" @click="openRelease('https://github.com/Aliuyanfeng/happytools')">
              <GithubOutlined />
              <span>源代码 · GitHub</span>
            </div>
          </div>

          <div class="about-meta">
            <div>开源协议：MIT License</div>
            <div>Built with ❤️ by Li6 &amp; Wails3</div>
            <div>© 2025 Li6. All rights reserved.</div>
          </div>
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
  WifiOutlined,
  LoadingOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
  WarningOutlined,
  DeleteOutlined,
  GithubOutlined,
} from '@ant-design/icons-vue'
import { useSettingsStore } from '../stores/settings'
import { modules as allModules } from '../config/modules'
import { UpdateService } from '../../bindings/github.com/Aliuyanfeng/happytools/backend/services/update/index'
import { VTService } from '../../bindings/github.com/Aliuyanfeng/happytools/backend/services/vt'
import { DNSService } from '../../bindings/github.com/Aliuyanfeng/happytools/backend/services/network/index'
import { AppSettingsService } from '../../bindings/github.com/Aliuyanfeng/happytools/backend/services/appsettings'

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

// VT 用户信息状态
interface VTUserInfo {
  id: string
  email: string
  type: string
  lastLogin: number
}
const userInfo = ref<VTUserInfo | null>(null)
const userInfoLoading = ref(false)
const userInfoError = ref('')

// 获取 VT 用户信息
async function fetchUserInfo(apiKey: string) {
  if (!apiKey) {
    userInfo.value = null
    userInfoError.value = ''
    return
  }
  userInfoLoading.value = true
  userInfoError.value = ''
  userInfo.value = null
  try {
    const info = await VTService.GetUserInfo(apiKey)
    console.log(info)
    userInfo.value = info as VTUserInfo
  } catch (e: any) {
    const msg = e?.message || String(e)
    if (msg.includes('invalid') || msg.includes('401') || msg.includes('403')) {
      userInfoError.value = t('settings.vtUserInfoInvalid')
    } else {
      userInfoError.value = t('settings.vtUserInfoFailed')
    }
  } finally {
    userInfoLoading.value = false
  }
}

// 格式化 Unix 时间戳
function formatTimestamp(ts: number): string {
  return new Date(ts * 1000).toLocaleString()
}

watch(() => props.open, (newVal) => {
  visible.value = newVal
})

watch(visible, (newVal) => {
  emit('update:open', newVal)
})

// 切换到 virustotal tab 时自动加载用户信息
watch(activeTab, (newVal) => {
  if (newVal === 'virustotal' && settingsStore.vtApiKey && !userInfo.value && !userInfoLoading.value) {
    fetchUserInfo(settingsStore.vtApiKey)
  }
})

function handleClose() {
  visible.value = false
}

function openVirusTotalDocs() {
  window.open('https://docs.virustotal.com/reference/overview', '_blank')
}

// API Key 输入框失焦时保存到后端并获取用户信息
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
    // 获取用户信息
    await fetchUserInfo(apiKey)
  } else {
    userInfo.value = null
    userInfoError.value = ''
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

// 清除所有数据
const clearing = ref(false)

function toggleModule(id: string, visible: boolean) {
  if (visible) {
    settingsStore.hiddenModules = settingsStore.hiddenModules.filter(m => m !== id)
  } else {
    if (!settingsStore.hiddenModules.includes(id)) {
      settingsStore.hiddenModules = [...settingsStore.hiddenModules, id]
    }
  }
}

// 关于页
const currentVersion = ref('...')
const checkingUpdate = ref(false)
const manualUpdateResult = ref<{ hasUpdate: boolean; latest: string; releaseUrl: string } | null>(null)

watch(() => props.open, async (v) => {
  if (v && currentVersion.value === '...') {
    currentVersion.value = await UpdateService.GetCurrentVersion()
  }
})

async function manualCheckUpdate() {
  checkingUpdate.value = true
  manualUpdateResult.value = null
  try {
    const result = await UpdateService.CheckUpdate()
    manualUpdateResult.value = result
      ? { hasUpdate: result.hasUpdate, latest: result.latest ?? '', releaseUrl: result.releaseUrl ?? '' }
      : { hasUpdate: false, latest: '', releaseUrl: '' }
  } catch {
    message.error('检查更新失败，请检查网络')
  } finally {
    checkingUpdate.value = false
  }
}

function openRelease(url: string) {
  import('@wailsio/runtime').then(({ Browser }) => Browser.OpenURL(url))
}

async function confirmClear() {
  const modal = await import('ant-design-vue').then(m => m.Modal)
  modal.confirm({
    title: '⚠️ 确认清除所有数据？',
    content: '此操作将永久删除：待办事项、分类、日报、VirusTotal 任务、Git 配置、Makefile 记录及模板。应用设置不受影响。此操作不可撤销！',
    okText: '确认清除',
    okType: 'danger',
    cancelText: '取消',
    onOk: async () => {
      clearing.value = true
      try {
        await AppSettingsService.ClearAllData()
        message.success('数据已清除')
      } catch (e: any) {
        message.error(e?.message ?? '清除失败')
      } finally {
        clearing.value = false
      }
    },
  })
}
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

/* 关于页 */
.about-panel {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
  padding-top: 20px;
}
.about-logo-wrap {
  width: 72px; height: 72px;
  border-radius: 18px;
  background: #f0f4ff;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid #e0e7ff;
}
.about-logo { width: 52px; height: 52px; object-fit: contain; }
.about-name { font-size: 18px; font-weight: 700; color: #1e1b4b; }
.about-version {
  display: flex;
  align-items: center;
  gap: 8px;
}
.about-ver-label { font-size: 12px; color: #94a3b8; }
.about-ver-val { font-size: 13px; font-weight: 600; color: #6366f1; font-family: monospace; }

.update-result {
  font-size: 13px;
  padding: 6px 14px;
  border-radius: 8px;
}
.update-result.has-update { background: #eff6ff; color: #2563eb; }
.update-result.up-to-date { background: #f0fdf4; color: #16a34a; }
.update-result a { cursor: pointer; text-decoration: underline; }

.about-links { width: 100%; display: flex; flex-direction: column; gap: 6px; }
.about-link-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  border-radius: 8px;
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  cursor: pointer;
  font-size: 13px;
  color: #334155;
  transition: background 0.15s;
}
.about-link-item:hover { background: #eef2ff; color: #6366f1; }

.about-meta {
  width: 100%;
  text-align: center;
  font-size: 11px;
  color: #94a3b8;
  line-height: 1.8;
}

/* 高级设置 - 模块管理 */
.module-manage {
  margin-bottom: 4px;
}
.mm-title {
  font-size: 14px;
  font-weight: 600;
  color: #262626;
  margin-bottom: 6px;
}
.mm-desc {
  font-size: 12px;
  color: #8c8c8c;
  margin-bottom: 12px;
}
.mm-list {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 8px;
}
.mm-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 10px;
  border-radius: 8px;
  background: #fafafa;
  border: 1px solid #f0f0f0;
}
.mm-name {
  font-size: 13px;
  color: #262626;
  white-space: nowrap;
}

/* 高级设置 - 危险区 */
.danger-zone {
  border: 1px solid #ffa39e;
  border-radius: 8px;
  padding: 20px;
  background: #fff2f0;
}
.danger-zone-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 6px;
}
.dz-icon { font-size: 18px; color: #cf1322; }
.dz-title { font-size: 15px; font-weight: 700; color: #cf1322; }
.dz-desc { font-size: 12px; color: #8c8c8c; margin-bottom: 16px; }

.danger-card {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
  padding: 16px;
  background: #fff;
  border: 1px solid #ffa39e;
  border-radius: 6px;
}
.dc-info { flex: 1; }
.dc-name { font-size: 14px; font-weight: 600; color: #262626; margin-bottom: 6px; }
.dc-desc { font-size: 12px; color: #595959; line-height: 1.6; }
.dc-desc strong { color: #cf1322; }

/* 侧边栏危险项 */
.danger-icon { color: #cf1322 !important; }
.danger-text { color: #cf1322; }

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

/* 用户信息卡片 */
.user-info-card {
  margin-top: 12px;
  padding: 12px 16px;
  border: 1px solid #e8e8e8;
  border-radius: 6px;
  background: #fafafa;
  font-size: 13px;
}

.user-info-loading {
  color: #1890ff;
  display: flex;
  align-items: center;
}

.user-info-error {
  color: #ff4d4f;
  display: flex;
  align-items: center;
}

.user-info-title {
  font-weight: 600;
  color: #52c41a;
  margin-bottom: 10px;
  display: flex;
  align-items: center;
}

.text-green {
  color: #52c41a;
}

.font-mono {
  font-family: 'Courier New', Courier, monospace;
}

[data-theme="dark"] .user-info-card {
  background: #1f1f1f;
  border-color: #434343;
}

[data-theme="dark"] .user-info-title {
  color: #73d13d;
}
</style>
