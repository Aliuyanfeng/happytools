<template>
  <div class="p-6 encryption-tool">
    <a-card :title="t('toolbox.encryption.title')" class="mb-6">
      <a-tabs v-model:activeKey="activeTab">
        <!-- MD5 -->
        <a-tab-pane key="md5" :tab="t('toolbox.encryption.tabMd5')">
          <a-form layout="vertical">
            <a-form-item :label="t('toolbox.encryption.inputLabel')">
              <a-textarea v-model:value="md5Input" :rows="4" :placeholder="t('toolbox.encryption.inputPlaceholder')" @keydown="handleBeforeInput" />
            </a-form-item>
            <a-form-item>
              <a-button type="primary" @click="encodeMD5" :loading="loading">{{ t('toolbox.encryption.encodeMd5') }}</a-button>
            </a-form-item>
            <a-form-item :label="t('toolbox.encryption.resultLabel')">
              <a-input-group compact>
                <a-input v-model:value="md5Result" readonly style="width: calc(100% - 80px)" />
                <a-button type="primary" @click="copyToClipboard(md5Result)">{{ t('toolbox.encryption.copy') }}</a-button>
              </a-input-group>
            </a-form-item>
          </a-form>
        </a-tab-pane>

        <!-- SHA1 -->
        <a-tab-pane key="sha1" :tab="t('toolbox.encryption.tabSha1')">
          <a-form layout="vertical">
            <a-form-item :label="t('toolbox.encryption.inputLabel')">
              <a-textarea v-model:value="sha1Input" :rows="4" :placeholder="t('toolbox.encryption.inputPlaceholder')" />
            </a-form-item>
            <a-form-item>
              <a-button type="primary" @click="encodeSHA1" :loading="loading">{{ t('toolbox.encryption.encodeSha1') }}</a-button>
            </a-form-item>
            <a-form-item :label="t('toolbox.encryption.resultLabel')">
              <a-input-group compact>
                <a-input v-model:value="sha1Result" readonly style="width: calc(100% - 80px)" />
                <a-button type="primary" @click="copyToClipboard(sha1Result)">{{ t('toolbox.encryption.copy') }}</a-button>
              </a-input-group>
            </a-form-item>
          </a-form>
        </a-tab-pane>

        <!-- SHA256 -->
        <a-tab-pane key="sha256" :tab="t('toolbox.encryption.tabSha256')">
          <a-form layout="vertical">
            <a-form-item :label="t('toolbox.encryption.inputLabel')">
              <a-textarea v-model:value="sha256Input" :rows="4" :placeholder="t('toolbox.encryption.inputPlaceholder')" />
            </a-form-item>
            <a-form-item>
              <a-button type="primary" @click="encodeSHA256" :loading="loading">{{ t('toolbox.encryption.encodeSha256') }}</a-button>
            </a-form-item>
            <a-form-item :label="t('toolbox.encryption.resultLabel')">
              <a-input-group compact>
                <a-input v-model:value="sha256Result" readonly style="width: calc(100% - 80px)" />
                <a-button type="primary" @click="copyToClipboard(sha256Result)">{{ t('toolbox.encryption.copy') }}</a-button>
              </a-input-group>
            </a-form-item>
          </a-form>
        </a-tab-pane>

        <!-- SHA512 -->
        <a-tab-pane key="sha512" :tab="t('toolbox.encryption.tabSha512')">
          <a-form layout="vertical">
            <a-form-item :label="t('toolbox.encryption.inputLabel')">
              <a-textarea v-model:value="sha512Input" :rows="4" :placeholder="t('toolbox.encryption.inputPlaceholder')" />
            </a-form-item>
            <a-form-item>
              <a-button type="primary" @click="encodeSHA512" :loading="loading">{{ t('toolbox.encryption.encodeSha512') }}</a-button>
            </a-form-item>
            <a-form-item :label="t('toolbox.encryption.resultLabel')">
              <a-input-group compact>
                <a-input v-model:value="sha512Result" readonly style="width: calc(100% - 80px)" />
                <a-button type="primary" @click="copyToClipboard(sha512Result)">{{ t('toolbox.encryption.copy') }}</a-button>
              </a-input-group>
            </a-form-item>
          </a-form>
        </a-tab-pane>

        <!-- Base64 -->
        <a-tab-pane key="base64" :tab="t('toolbox.encryption.tabBase64')">
          <a-form layout="vertical">
            <a-form-item :label="t('toolbox.encryption.inputLabel')">
              <a-textarea v-model:value="base64Input" :rows="4" :placeholder="t('toolbox.encryption.inputPlaceholderCodec')" />
            </a-form-item>
            <a-form-item>
              <a-space>
                <a-button type="primary" @click="encodeBase64" :loading="loading">{{ t('toolbox.encryption.encodeBase64') }}</a-button>
                <a-button @click="decodeBase64" :loading="loading">{{ t('toolbox.encryption.decodeBase64') }}</a-button>
              </a-space>
            </a-form-item>
            <a-form-item :label="t('toolbox.encryption.resultLabelCodec')">
              <a-input-group compact>
                <a-input v-model:value="base64Result" readonly style="width: calc(100% - 80px)" />
                <a-button type="primary" @click="copyToClipboard(base64Result)">{{ t('toolbox.encryption.copy') }}</a-button>
              </a-input-group>
            </a-form-item>
          </a-form>
        </a-tab-pane>

        <!-- UTF-8 -->
        <a-tab-pane key="utf8" :tab="t('toolbox.encryption.tabUtf8')">
          <a-form layout="vertical">
            <a-form-item :label="t('toolbox.encryption.inputLabel')">
              <a-textarea v-model:value="utf8Input" :rows="4" :placeholder="t('toolbox.encryption.inputPlaceholderCodec')" />
            </a-form-item>
            <a-form-item>
              <a-space>
                <a-button type="primary" @click="encodeUTF8ToHex" :loading="loading">{{ t('toolbox.encryption.toHex') }}</a-button>
                <a-button @click="decodeHexToUTF8" :loading="loading">{{ t('toolbox.encryption.fromHex') }}</a-button>
              </a-space>
            </a-form-item>
            <a-form-item :label="t('toolbox.encryption.resultLabelCodec')">
              <a-input-group compact>
                <a-input v-model:value="utf8Result" readonly style="width: calc(100% - 80px)" />
                <a-button type="primary" @click="copyToClipboard(utf8Result)">{{ t('toolbox.encryption.copy') }}</a-button>
              </a-input-group>
            </a-form-item>
          </a-form>
        </a-tab-pane>

        <!-- 批量编码 -->
        <a-tab-pane key="batch" :tab="t('toolbox.encryption.tabBatch')">
          <a-form layout="vertical">
            <a-form-item :label="t('toolbox.encryption.inputLabel')">
              <a-textarea v-model:value="batchInput" :rows="4" :placeholder="t('toolbox.encryption.inputPlaceholderBatch')" />
            </a-form-item>
            <a-form-item>
              <a-button type="primary" @click="batchEncode" :loading="loading">{{ t('toolbox.encryption.batchEncode') }}</a-button>
            </a-form-item>
            <div v-if="batchResult" class="result-section">
              <a-divider>{{ t('toolbox.encryption.batchResultTitle') }}</a-divider>
              <a-row :gutter="[16, 16]">
                <a-col :span="24" v-for="(value, key) in batchResult" :key="key">
                  <a-card size="small">
                    <template #title>
                      <span class="text-sm font-bold">{{ getAlgorithmLabel(key) }}</span>
                    </template>
                    <a-input-group compact>
                      <a-input :value="value" readonly style="width: calc(100% - 80px)" />
                      <a-button type="primary" size="small" @click="copyToClipboard(value)">{{ t('toolbox.encryption.copy') }}</a-button>
                    </a-input-group>
                  </a-card>
                </a-col>
              </a-row>
            </div>
          </a-form>
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { message } from 'ant-design-vue'
import { useI18n } from 'vue-i18n'
import * as EncryptionService from '../../../bindings/github.com/Aliuyanfeng/happytools/backend/services/encryption/encryptionservice'
import * as ClipboardService from '../../../bindings/github.com/Aliuyanfeng/happytools/backend/services/clipboard/clipboardservice'

const { t } = useI18n()
const activeTab = ref('md5')
const loading = ref(false)

// MD5
const md5Input = ref('')
const md5Result = ref('')

// SHA1
const sha1Input = ref('')
const sha1Result = ref('')

// SHA256
const sha256Input = ref('')
const sha256Result = ref('')

// SHA512
const sha512Input = ref('')
const sha512Result = ref('')

// Base64
const base64Input = ref('')
const base64Result = ref('')

// UTF-8
const utf8Input = ref('')
const utf8Result = ref('')

// 批量编码
const batchInput = ref('')
const batchResult = ref<any>(null)

// 复制锁，防止重复复制
let copyLock = false

// MD5 编码
const encodeMD5 = async () => {
  if (!md5Input.value) {
    message.warning(t('toolbox.encryption.pleaseInput'))
    return
  }
  loading.value = true
  try {
    const result = await EncryptionService.MD5Encode(md5Input.value)
    md5Result.value = result
  } catch (error: any) {
    message.error(error.message || t('toolbox.encryption.encodeFailed'))
  } finally {
    loading.value = false
  }
}

const handleBeforeInput = async (e) => {
  if (e.ctrlKey && e.key.toLowerCase() === 'v') {
    e.preventDefault()
    e.stopPropagation()

    const text = await navigator.clipboard.readText()

    const el = e.target
    const start = el.selectionStart
    const end = el.selectionEnd

    const value = el.value
    const newValue =
      value.slice(0, start) + text + value.slice(end)

    el.value = newValue

    el.dispatchEvent(new Event('input', { bubbles: true }))
  }
}
const encodeSHA1 = async () => {
  if (!sha1Input.value) { message.warning(t('toolbox.encryption.pleaseInput')); return }
  loading.value = true
  try {
    sha1Result.value = await EncryptionService.SHA1Encode(sha1Input.value)
  } catch (error: any) { message.error(error.message || t('toolbox.encryption.encodeFailed')) }
  finally { loading.value = false }
}

const encodeSHA256 = async () => {
  if (!sha256Input.value) { message.warning(t('toolbox.encryption.pleaseInput')); return }
  loading.value = true
  try {
    sha256Result.value = await EncryptionService.SHA256Encode(sha256Input.value)
  } catch (error: any) { message.error(error.message || t('toolbox.encryption.encodeFailed')) }
  finally { loading.value = false }
}

const encodeSHA512 = async () => {
  if (!sha512Input.value) { message.warning(t('toolbox.encryption.pleaseInput')); return }
  loading.value = true
  try {
    sha512Result.value = await EncryptionService.SHA512Encode(sha512Input.value)
  } catch (error: any) { message.error(error.message || t('toolbox.encryption.encodeFailed')) }
  finally { loading.value = false }
}

const encodeBase64 = async () => {
  if (!base64Input.value) { message.warning(t('toolbox.encryption.pleaseInput')); return }
  loading.value = true
  try {
    base64Result.value = await EncryptionService.Base64Encode(base64Input.value)
  } catch (error: any) { message.error(error.message || t('toolbox.encryption.encodeFailed')) }
  finally { loading.value = false }
}

const decodeBase64 = async () => {
  if (!base64Input.value) { message.warning(t('toolbox.encryption.pleaseInput')); return }
  loading.value = true
  try {
    base64Result.value = await EncryptionService.Base64Decode(base64Input.value)
  } catch (error: any) { message.error(error.message || t('toolbox.encryption.decodeFailed')) }
  finally { loading.value = false }
}

const encodeUTF8ToHex = async () => {
  if (!utf8Input.value) { message.warning(t('toolbox.encryption.pleaseInput')); return }
  loading.value = true
  try {
    utf8Result.value = await EncryptionService.UTF8ToHex(utf8Input.value)
  } catch (error: any) { message.error(error.message || t('toolbox.encryption.encodeFailed')) }
  finally { loading.value = false }
}

const decodeHexToUTF8 = async () => {
  if (!utf8Input.value) { message.warning(t('toolbox.encryption.pleaseInputHex')); return }
  loading.value = true
  try {
    utf8Result.value = await EncryptionService.HexToUTF8(utf8Input.value)
  } catch (error: any) { message.error(error.message || t('toolbox.encryption.decodeHexFailed')) }
  finally { loading.value = false }
}

const batchEncode = async () => {
  if (!batchInput.value) { message.warning(t('toolbox.encryption.pleaseInput')); return }
  loading.value = true
  try {
    batchResult.value = await EncryptionService.BatchEncode(batchInput.value)
  } catch (error: any) { message.error(error.message || t('toolbox.encryption.encodeFailed')) }
  finally { loading.value = false }
}

// 复制到剪贴板 - 使用后端剪贴板服务
const copyToClipboard = async (text: string) => {
  if (!text) {
    message.warning(t('toolbox.encryption.noCopyContent'))
    return
  }
  if (copyLock) return
  copyLock = true
  try {
    await ClipboardService.SetText(text)
    message.success(t('toolbox.encryption.copied'))
    setTimeout(() => { copyLock = false }, 500)
  } catch (error) {
    message.error(t('toolbox.encryption.copyFailed'))
    copyLock = false
  }
}

// 获取算法标签
const getAlgorithmLabel = (key: string | number) => {
  const keyStr = String(key)
  const labels: Record<string, string> = {
    md5: 'MD5',
    sha1: 'SHA1',
    sha256: 'SHA256',
    sha512: 'SHA512',
    base64: 'Base64'
  }
  return labels[keyStr] || keyStr.toUpperCase()
}
</script>

<style scoped>
.encryption-tool {
  max-width: 1200px;
  margin: 0 auto;
}

.result-section {
  margin-top: 24px;
  padding: 16px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 8px;
}
</style>
