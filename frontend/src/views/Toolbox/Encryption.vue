<template>
  <div class="p-6 encryption-tool">
    <a-card title="加密工具" class="mb-6">
      <a-tabs v-model:activeKey="activeTab">
        <!-- MD5 编码 -->
        <a-tab-pane key="md5" tab="MD5">
          <a-form layout="vertical">
            <a-form-item label="输入文本">
              <a-textarea
                v-model:value="md5Input"
                :rows="4"
                placeholder="请输入要编码的文本"
                @keydown="handleBeforeInput"
              />
            </a-form-item>
            <a-form-item>
              <a-button type="primary" @click="encodeMD5" :loading="loading">
                MD5 编码
              </a-button>
            </a-form-item>
            <a-form-item label="编码结果">
              <a-input-group compact>
                <a-input v-model:value="md5Result" readonly style="width: calc(100% - 80px)" />
                <a-button type="primary" @click="copyToClipboard(md5Result)">复制</a-button>
              </a-input-group>
            </a-form-item>
          </a-form>
        </a-tab-pane>

        <!-- SHA1 编码 -->
        <a-tab-pane key="sha1" tab="SHA1">
          <a-form layout="vertical">
            <a-form-item label="输入文本">
              <a-textarea
                v-model:value="sha1Input"
                :rows="4"
                placeholder="请输入要编码的文本"
              />
            </a-form-item>
            <a-form-item>
              <a-button type="primary" @click="encodeSHA1" :loading="loading">
                SHA1 编码
              </a-button>
            </a-form-item>
            <a-form-item label="编码结果">
              <a-input-group compact>
                <a-input v-model:value="sha1Result" readonly style="width: calc(100% - 80px)" />
                <a-button type="primary" @click="copyToClipboard(sha1Result)">复制</a-button>
              </a-input-group>
            </a-form-item>
          </a-form>
        </a-tab-pane>

        <!-- SHA256 编码 -->
        <a-tab-pane key="sha256" tab="SHA256">
          <a-form layout="vertical">
            <a-form-item label="输入文本">
              <a-textarea
                v-model:value="sha256Input"
                :rows="4"
                placeholder="请输入要编码的文本"
              />
            </a-form-item>
            <a-form-item>
              <a-button type="primary" @click="encodeSHA256" :loading="loading">
                SHA256 编码
              </a-button>
            </a-form-item>
            <a-form-item label="编码结果">
              <a-input-group compact>
                <a-input v-model:value="sha256Result" readonly style="width: calc(100% - 80px)" />
                <a-button type="primary" @click="copyToClipboard(sha256Result)">复制</a-button>
              </a-input-group>
            </a-form-item>
          </a-form>
        </a-tab-pane>

        <!-- SHA512 编码 -->
        <a-tab-pane key="sha512" tab="SHA512">
          <a-form layout="vertical">
            <a-form-item label="输入文本">
              <a-textarea
                v-model:value="sha512Input"
                :rows="4"
                placeholder="请输入要编码的文本"
              />
            </a-form-item>
            <a-form-item>
              <a-button type="primary" @click="encodeSHA512" :loading="loading">
                SHA512 编码
              </a-button>
            </a-form-item>
            <a-form-item label="编码结果">
              <a-input-group compact>
                <a-input v-model:value="sha512Result" readonly style="width: calc(100% - 80px)" />
                <a-button type="primary" @click="copyToClipboard(sha512Result)">复制</a-button>
              </a-input-group>
            </a-form-item>
          </a-form>
        </a-tab-pane>

        <!-- Base64 编解码 -->
        <a-tab-pane key="base64" tab="Base64">
          <a-form layout="vertical">
            <a-form-item label="输入文本">
              <a-textarea
                v-model:value="base64Input"
                :rows="4"
                placeholder="请输入要编解码的文本"
              />
            </a-form-item>
            <a-form-item>
              <a-space>
                <a-button type="primary" @click="encodeBase64" :loading="loading">
                  Base64 编码
                </a-button>
                <a-button @click="decodeBase64" :loading="loading">
                  Base64 解码
                </a-button>
              </a-space>
            </a-form-item>
            <a-form-item label="结果">
              <a-input-group compact>
                <a-input v-model:value="base64Result" readonly style="width: calc(100% - 80px)" />
                <a-button type="primary" @click="copyToClipboard(base64Result)">复制</a-button>
              </a-input-group>
            </a-form-item>
          </a-form>
        </a-tab-pane>

        <!-- UTF-8 编解码 -->
        <a-tab-pane key="utf8" tab="UTF-8">
          <a-form layout="vertical">
            <a-form-item label="输入文本">
              <a-textarea
                v-model:value="utf8Input"
                :rows="4"
                placeholder="请输入要编解码的文本"
              />
            </a-form-item>
            <a-form-item>
              <a-space>
                <a-button type="primary" @click="encodeUTF8ToHex" :loading="loading">
                  转16进制
                </a-button>
                <a-button @click="decodeHexToUTF8" :loading="loading">
                  16进制转文本
                </a-button>
              </a-space>
            </a-form-item>
            <a-form-item label="结果">
              <a-input-group compact>
                <a-input v-model:value="utf8Result" readonly style="width: calc(100% - 80px)" />
                <a-button type="primary" @click="copyToClipboard(utf8Result)">复制</a-button>
              </a-input-group>
            </a-form-item>
          </a-form>
        </a-tab-pane>

        <!-- 批量编码 -->
        <a-tab-pane key="batch" tab="批量编码">
          <a-form layout="vertical">
            <a-form-item label="输入文本">
              <a-textarea
                v-model:value="batchInput"
                :rows="4"
                placeholder="请输入要编码的文本，将一次性生成所有哈希值"
              />
            </a-form-item>
            <a-form-item>
              <a-button type="primary" @click="batchEncode" :loading="loading">
                批量编码
              </a-button>
            </a-form-item>
            <div v-if="batchResult" class="result-section">
              <a-divider>编码结果</a-divider>
              <a-row :gutter="[16, 16]">
                <a-col :span="24" v-for="(value, key) in batchResult" :key="key">
                  <a-card size="small">
                    <template #title>
                      <span class="text-sm font-bold">{{ getAlgorithmLabel(key) }}</span>
                    </template>
                    <a-input-group compact>
                      <a-input :value="value" readonly style="width: calc(100% - 80px)" />
                      <a-button type="primary" size="small" @click="copyToClipboard(value)">
                        复制
                      </a-button>
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
import * as EncryptionService from '../../../bindings/github.com/Aliuyanfeng/happytools/backend/services/encryption/encryptionservice'
import * as ClipboardService from '../../../bindings/github.com/Aliuyanfeng/happytools/backend/services/clipboard/clipboardservice'

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
    message.warning('请输入文本')
    return
  }

  loading.value = true
  try {
    const result = await EncryptionService.MD5Encode(md5Input.value)
    md5Result.value = result
  } catch (error: any) {
    message.error(error.message || '编码失败')
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
// SHA1 编码
const encodeSHA1 = async () => {
  if (!sha1Input.value) {
    message.warning('请输入文本')
    return
  }

  loading.value = true
  try {
    const result = await EncryptionService.SHA1Encode(sha1Input.value)
    sha1Result.value = result
  } catch (error: any) {
    message.error(error.message || '编码失败')
  } finally {
    loading.value = false
  }
}

// SHA256 编码
const encodeSHA256 = async () => {
  if (!sha256Input.value) {
    message.warning('请输入文本')
    return
  }

  loading.value = true
  try {
    const result = await EncryptionService.SHA256Encode(sha256Input.value)
    sha256Result.value = result
  } catch (error: any) {
    message.error(error.message || '编码失败')
  } finally {
    loading.value = false
  }
}

// SHA512 编码
const encodeSHA512 = async () => {
  if (!sha512Input.value) {
    message.warning('请输入文本')
    return
  }

  loading.value = true
  try {
    const result = await EncryptionService.SHA512Encode(sha512Input.value)
    sha512Result.value = result
  } catch (error: any) {
    message.error(error.message || '编码失败')
  } finally {
    loading.value = false
  }
}

// Base64 编码
const encodeBase64 = async () => {
  if (!base64Input.value) {
    message.warning('请输入文本')
    return
  }

  loading.value = true
  try {
    const result = await EncryptionService.Base64Encode(base64Input.value)
    base64Result.value = result
  } catch (error: any) {
    message.error(error.message || '编码失败')
  } finally {
    loading.value = false
  }
}

// Base64 解码
const decodeBase64 = async () => {
  if (!base64Input.value) {
    message.warning('请输入文本')
    return
  }

  loading.value = true
  try {
    const result = await EncryptionService.Base64Decode(base64Input.value)
    base64Result.value = result
  } catch (error: any) {
    message.error(error.message || '解码失败，请检查输入是否为有效的 Base64 字符串')
  } finally {
    loading.value = false
  }
}

// UTF-8 转16进制
const encodeUTF8ToHex = async () => {
  if (!utf8Input.value) {
    message.warning('请输入文本')
    return
  }

  loading.value = true
  try {
    const result = await EncryptionService.UTF8ToHex(utf8Input.value)
    utf8Result.value = result
  } catch (error: any) {
    message.error(error.message || '编码失败')
  } finally {
    loading.value = false
  }
}

// 16进制转 UTF-8
const decodeHexToUTF8 = async () => {
  if (!utf8Input.value) {
    message.warning('请输入16进制字符串')
    return
  }

  loading.value = true
  try {
    const result = await EncryptionService.HexToUTF8(utf8Input.value)
    utf8Result.value = result
  } catch (error: any) {
    message.error(error.message || '解码失败，请检查输入是否为有效的16进制字符串')
  } finally {
    loading.value = false
  }
}

// 批量编码
const batchEncode = async () => {
  if (!batchInput.value) {
    message.warning('请输入文本')
    return
  }

  loading.value = true
  try {
    const result = await EncryptionService.BatchEncode(batchInput.value)
    batchResult.value = result
  } catch (error: any) {
    message.error(error.message || '编码失败')
  } finally {
    loading.value = false
  }
}

// 复制到剪贴板 - 使用后端剪贴板服务
const copyToClipboard = async (text: string) => {
  if (!text) {
    message.warning('没有可复制的内容')
    return
  }

  // 防止重复调用
  if (copyLock) {
    return
  }

  copyLock = true

  try {
    // 使用后端剪贴板服务
    await ClipboardService.SetText(text)
    message.success('已复制到剪贴板')
    console.log('复制成功:', text)
    // 延迟解锁，防止快速重复点击
    setTimeout(() => {
      copyLock = false
    }, 500)
  } catch (error) {
    message.error('复制失败')
    console.error('复制失败:', error)
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
