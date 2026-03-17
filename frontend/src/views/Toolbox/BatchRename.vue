<template>
  <div class="p-6">
    <a-card :title="t('toolbox.batchRename.title')">
      <!-- 拖拽区域 -->
      <div
        id="batch-rename-drop"
        data-file-drop-target
        class="drop-zone"
        @click="selectFiles"
      >
        <InboxOutlined class="drop-icon" />
        <p class="mt-2 text-sm text-gray-600">{{ t('toolbox.batchRename.dropHint') }}</p>
        <p class="text-xs text-gray-400">{{ t('toolbox.batchRename.dropHintSub') }}</p>
      </div>

      <!-- 文件列表 -->
      <div v-if="files.length > 0" class="mt-4">
        <div class="mb-2 flex items-center justify-between">
          <span class="text-sm text-gray-500">{{ t('toolbox.batchRename.filesAdded', { count: files.length }) }}</span>
          <a-space>
            <a-button size="small" danger @click="clearFiles">{{ t('toolbox.batchRename.clear') }}</a-button>
          </a-space>
        </div>
        <a-table
          :columns="fileCols"
          :data-source="previewList.length ? previewList : files"
          :pagination="{ pageSize: 8, size: 'small' }"
          row-key="originalPath"
          size="small"
          bordered
        >
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'originalName'">
              <span class="font-mono text-xs">{{ record.originalName }}</span>
            </template>
            <template v-if="column.key === 'newName'">
              <span v-if="record.error" class="text-red-500 text-xs">{{ record.error }}</span>
              <span v-else-if="record.newName" class="font-mono text-xs text-green-600">{{ record.newName }}</span>
              <span v-else class="text-gray-300 text-xs">—</span>
            </template>
            <template v-if="column.key === 'size'">
              <span class="text-xs">{{ formatBytes(record.size) }}</span>
            </template>
            <template v-if="column.key === 'action'">
              <a-button type="link" size="small" danger @click="removeFile(record.originalPath)">{{ t('toolbox.batchRename.remove') }}</a-button>
            </template>
          </template>
        </a-table>
      </div>
    </a-card>

    <!-- 重命名规则配置 -->
    <a-card :title="t('toolbox.batchRename.rulesTitle')" class="mt-4">
      <a-radio-group v-model:value="mode" button-style="solid" class="mb-4">
        <a-radio-button value="custom">{{ t('toolbox.batchRename.modeCustom') }}</a-radio-button>
        <a-radio-button value="hash">{{ t('toolbox.batchRename.modeHash') }}</a-radio-button>
      </a-radio-group>

      <!-- 自定义模式 -->
      <div v-if="mode === 'custom'">
        <a-form layout="inline" class="mb-2">
          <a-form-item :label="t('toolbox.batchRename.prefix')">
            <a-input v-model:value="customRule.prefix" :placeholder="t('toolbox.batchRename.prefixPlaceholder')" style="width:140px" />
          </a-form-item>
          <a-form-item :label="t('toolbox.batchRename.startNumber')">
            <a-input-number v-model:value="customRule.startNumber" :min="0" style="width:90px" />
          </a-form-item>
          <a-form-item :label="t('toolbox.batchRename.numberDigits')">
            <a-input-number v-model:value="customRule.numberDigits" :min="1" :max="10" style="width:80px" />
          </a-form-item>
          <a-form-item :label="t('toolbox.batchRename.numberStep')">
            <a-input-number v-model:value="customRule.numberStep" :min="1" style="width:80px" />
          </a-form-item>
          <a-form-item :label="t('toolbox.batchRename.keepExtension')">
            <a-switch v-model:checked="customRule.keepExtension" />
          </a-form-item>
        </a-form>
        <div class="text-xs text-gray-400 mb-3">
          {{ t('toolbox.batchRename.previewFormat') }}：<span class="font-mono text-gray-600">{{ previewCustomName }}</span>
        </div>
      </div>

      <!-- 哈希模式 -->
      <div v-if="mode === 'hash'">
        <a-form layout="inline" class="mb-2">
          <a-form-item :label="t('toolbox.batchRename.hashAlgorithm')">
            <a-radio-group v-model:value="hashRule.algorithm" button-style="outline">
              <a-radio-button value="md5">MD5</a-radio-button>
              <a-radio-button value="sha1">SHA1</a-radio-button>
              <a-radio-button value="sha256">SHA256</a-radio-button>
            </a-radio-group>
          </a-form-item>
          <a-form-item :label="t('toolbox.batchRename.keepExtension')">
            <a-switch v-model:checked="hashRule.keepExtension" />
          </a-form-item>
        </a-form>
        <div class="text-xs text-gray-400">{{ t('toolbox.batchRename.hashHint') }}</div>
      </div>

      <a-divider class="my-3" />

      <a-space>
        <a-button @click="doPreview" :loading="loading" :disabled="files.length === 0">
          {{ t('toolbox.batchRename.preview') }}
        </a-button>
        <a-button
          type="primary"
          @click="doRename"
          :loading="loading"
          :disabled="files.length === 0"
        >
          {{ t('toolbox.batchRename.execute') }}
        </a-button>
      </a-space>
    </a-card>

    <!-- 执行结果 -->
    <a-card v-if="result" :title="t('toolbox.batchRename.resultTitle')" class="mt-4">
      <a-space class="mb-3">
        <a-tag color="green">{{ t('toolbox.batchRename.resultSuccess', { count: result.successCount }) }}</a-tag>
        <a-tag color="red" v-if="result.failedCount > 0">{{ t('toolbox.batchRename.resultFailed', { count: result.failedCount }) }}</a-tag>
        <a-tag>{{ t('toolbox.batchRename.resultTotal', { count: result.totalCount }) }}</a-tag>
      </a-space>
      <a-table
        :columns="resultCols"
        :data-source="result.results"
        :pagination="{ pageSize: 8, size: 'small' }"
        row-key="originalPath"
        size="small"
        bordered
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'status'">
            <a-tag :color="record.error ? 'red' : 'green'">
              {{ record.error ? t('toolbox.batchRename.statusFailed') : t('toolbox.batchRename.statusSuccess') }}
            </a-tag>
          </template>
          <template v-if="column.key === 'newName'">
            <span v-if="record.error" class="text-red-500 text-xs">{{ record.error }}</span>
            <span v-else class="font-mono text-xs text-green-600">{{ record.newName }}</span>
          </template>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, onUnmounted } from 'vue'
import { message } from 'ant-design-vue'
import { InboxOutlined } from '@ant-design/icons-vue'
import { useI18n } from 'vue-i18n'
import { Events } from '@wailsio/runtime'
import * as RenameService from '../../../bindings/github.com/Aliuyanfeng/happytools/backend/services/rename/renameservice'

interface FileInfo {
  originalPath: string
  originalName: string
  newName: string
  newPath: string
  size: number
  isDir: boolean
  error: string
}

interface RenameResult {
  successCount: number
  failedCount: number
  totalCount: number
  results: FileInfo[]
}

const { t } = useI18n()
const isDragOver = ref(false)
const loading = ref(false)
const mode = ref<'custom' | 'hash'>('custom')
const files = ref<FileInfo[]>([])
const previewList = ref<FileInfo[]>([])
const result = ref<RenameResult | null>(null)

// 监听 Wails 文件拖拽事件
let offFileDrop: (() => void) | null = null

onMounted(() => {
  // 调试：检查 enableFileDrop flag
  const flags = (window as any)._wails?.flags
  console.log('[FileDrop] _wails.flags:', flags)

  offFileDrop = Events.On('wails:file-drop', (event: any) => {
    console.log('[FileDrop] 收到事件:', event)
    const data = event?.data
    const droppedFiles: string[] = Array.isArray(data?.files) ? data.files : []
    const targetId: string = data?.target ?? ''

    console.log('[FileDrop] targetId:', targetId, 'files:', droppedFiles)

    // 只处理拖到 batch-rename-drop 区域的文件（targetId 为空时也接受，兼容 body drop）
    if (targetId && targetId !== 'batch-rename-drop') return
    if (droppedFiles.length === 0) return

    addPaths(droppedFiles)
  })
})

onUnmounted(() => {
  if (offFileDrop) {
    offFileDrop()
    offFileDrop = null
  }
})

const customRule = reactive({
  prefix: 'file',
  startNumber: 1,
  numberDigits: 3,
  numberStep: 1,
  keepExtension: true,
  // 以下字段保持默认（不使用替换功能）
  suffix: '',
  replaceFrom: '',
  replaceTo: '',
  caseSensitive: false,
  previewBeforeRename: false
})

const hashRule = reactive({
  algorithm: 'md5',
  keepExtension: true
})

// 自定义模式预览名称示例
const previewCustomName = computed(() => {
  const num = String(customRule.startNumber).padStart(customRule.numberDigits || 1, '0')
  const ext = customRule.keepExtension ? '.ext' : ''
  const prefix = customRule.prefix ? customRule.prefix + '-' : ''
  const suffix = customRule.suffix ? '-' + customRule.suffix : ''
  return `${prefix}${num}${suffix}${ext}`
})

const fileCols = computed(() => [
  { title: t('toolbox.batchRename.colOriginalName'), dataIndex: 'originalName', key: 'originalName', ellipsis: true },
  { title: t('toolbox.batchRename.colNewName'), dataIndex: 'newName', key: 'newName', ellipsis: true },
  { title: t('toolbox.batchRename.colSize'), dataIndex: 'size', key: 'size', width: 90 },
  { title: t('toolbox.batchRename.colAction'), key: 'action', width: 70 }
])

const resultCols = computed(() => [
  { title: t('toolbox.batchRename.colOriginalName'), dataIndex: 'originalName', key: 'originalName', ellipsis: true },
  { title: t('toolbox.batchRename.colNewNameResult'), dataIndex: 'newName', key: 'newName', ellipsis: true },
  { title: t('toolbox.batchRename.colStatus'), key: 'status', width: 80 }
])

// 点击选择文件
const selectFiles = async () => {
  try {
    const paths = await RenameService.OpenFileDialogs()
    if (paths?.length) await addPaths(paths)
  } catch (e: any) {
    message.error(t('toolbox.batchRename.selectFailed') + '：' + (e?.message || e))
  }
}

// 添加路径（去重）
const addPaths = async (paths: string[]) => {
  const existing = new Set(files.value.map(f => f.originalPath))
  const newPaths = paths.filter(p => !existing.has(p))
  if (!newPaths.length) {
    message.info(t('toolbox.batchRename.allAdded'))
    return
  }
  try {
    const infos = await RenameService.BatchGetFileInfo(newPaths)
    files.value.push(...infos)
    previewList.value = []
    result.value = null
  } catch (e: any) {
    message.error(t('toolbox.batchRename.getInfoFailed') + '：' + (e?.message || e))
  }
}

// 移除单个文件
const removeFile = (path: string) => {
  files.value = files.value.filter(f => f.originalPath !== path)
  previewList.value = previewList.value.filter(f => f.originalPath !== path)
}

// 清空
const clearFiles = () => {
  files.value = []
  previewList.value = []
  result.value = null
}

// 预览
const doPreview = async () => {
  loading.value = true
  try {
    if (mode.value === 'custom') {
      previewList.value = await RenameService.PreviewRename(files.value, customRule)
    } else {
      previewList.value = await RenameService.PreviewHashRename(files.value, hashRule)
    }
    message.success(t('toolbox.batchRename.previewDone'))
  } catch (e: any) {
    message.error(t('toolbox.batchRename.previewFailed') + '：' + (e?.message || e))
  } finally {
    loading.value = false
  }
}

// 执行重命名
const doRename = async () => {
  loading.value = true
  result.value = null
  try {
    let res: RenameResult | null
    if (mode.value === 'custom') {
      res = await RenameService.ExecuteRename(files.value, customRule)
    } else {
      res = await RenameService.ExecuteHashRename(files.value, hashRule)
    }
    result.value = res
    if (res) {
      if (res.failedCount === 0) {
        message.success(t('toolbox.batchRename.allSuccess', { count: res.successCount }))
        clearFiles()
      } else {
        message.warning(t('toolbox.batchRename.partialSuccess', { success: res.successCount, failed: res.failedCount }))
      }
    }
  } catch (e: any) {
    message.error(t('toolbox.batchRename.renameFailed') + '：' + (e?.message || e))
  } finally {
    loading.value = false
  }
}

const formatBytes = (n: number) => {
  if (!n) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(n) / Math.log(1024))
  return (n / Math.pow(1024, i)).toFixed(i === 0 ? 0 : 1) + ' ' + units[i]
}
</script>

<style scoped>
.drop-zone {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 140px;
  border: 2px dashed #d9d9d9;
  border-radius: 8px;
  cursor: pointer;
  transition: border-color 0.2s, background 0.2s;
}
.drop-zone:hover,
.drop-zone.file-drop-target-active {
  border-color: #1890ff;
  background: #f0f8ff;
}
.drop-icon {
  font-size: 40px;
  color: #bfbfbf;
}
.drop-zone:hover .drop-icon,
.drop-zone.file-drop-target-active .drop-icon {
  color: #1890ff;
}
.font-mono {
  font-family: 'Courier New', Courier, monospace;
}
</style>
