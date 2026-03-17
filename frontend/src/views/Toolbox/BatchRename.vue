<template>
  <div class="p-6">
    <a-card title="批量文件重命名">
      <!-- 拖拽区域 -->
      <div
        id="batch-rename-drop"
        data-file-drop-target
        class="drop-zone"
        @click="selectFiles"
      >
        <InboxOutlined class="drop-icon" />
        <p class="mt-2 text-sm text-gray-600">拖拽文件到此处，或点击选择文件</p>
        <p class="text-xs text-gray-400">支持多文件，拖拽可追加</p>
      </div>

      <!-- 文件列表 -->
      <div v-if="files.length > 0" class="mt-4">
        <div class="mb-2 flex items-center justify-between">
          <span class="text-sm text-gray-500">已添加 <strong>{{ files.length }}</strong> 个文件</span>
          <a-space>
            <a-button size="small" danger @click="clearFiles">清空</a-button>
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
              <a-button type="link" size="small" danger @click="removeFile(record.originalPath)">移除</a-button>
            </template>
          </template>
        </a-table>
      </div>
    </a-card>

    <!-- 重命名规则配置 -->
    <a-card title="重命名规则" class="mt-4">
      <a-radio-group v-model:value="mode" button-style="solid" class="mb-4">
        <a-radio-button value="custom">自定义重命名</a-radio-button>
        <a-radio-button value="hash">哈希重命名</a-radio-button>
      </a-radio-group>

      <!-- 自定义模式 -->
      <div v-if="mode === 'custom'">
        <a-form layout="inline" class="mb-2">
          <a-form-item label="前缀文字">
            <a-input v-model:value="customRule.prefix" placeholder="例如: file" style="width:140px" />
          </a-form-item>
          <a-form-item label="起始序号">
            <a-input-number v-model:value="customRule.startNumber" :min="0" style="width:90px" />
          </a-form-item>
          <a-form-item label="序号位数">
            <a-input-number v-model:value="customRule.numberDigits" :min="1" :max="10" style="width:80px" />
          </a-form-item>
          <a-form-item label="步长">
            <a-input-number v-model:value="customRule.numberStep" :min="1" style="width:80px" />
          </a-form-item>
          <a-form-item label="保留扩展名">
            <a-switch v-model:checked="customRule.keepExtension" />
          </a-form-item>
        </a-form>
        <div class="text-xs text-gray-400 mb-3">
          预览格式：<span class="font-mono text-gray-600">{{ previewCustomName }}</span>
        </div>
      </div>

      <!-- 哈希模式 -->
      <div v-if="mode === 'hash'">
        <a-form layout="inline" class="mb-2">
          <a-form-item label="哈希算法">
            <a-radio-group v-model:value="hashRule.algorithm" button-style="outline">
              <a-radio-button value="md5">MD5</a-radio-button>
              <a-radio-button value="sha1">SHA1</a-radio-button>
              <a-radio-button value="sha256">SHA256</a-radio-button>
            </a-radio-group>
          </a-form-item>
          <a-form-item label="保留扩展名">
            <a-switch v-model:checked="hashRule.keepExtension" />
          </a-form-item>
        </a-form>
        <div class="text-xs text-gray-400">文件名将替换为文件内容的哈希值</div>
      </div>

      <a-divider class="my-3" />

      <a-space>
        <a-button @click="doPreview" :loading="loading" :disabled="files.length === 0">
          预览结果
        </a-button>
        <a-button
          type="primary"
          @click="doRename"
          :loading="loading"
          :disabled="files.length === 0"
        >
          执行重命名
        </a-button>
      </a-space>
    </a-card>

    <!-- 执行结果 -->
    <a-card v-if="result" title="执行结果" class="mt-4">
      <a-space class="mb-3">
        <a-tag color="green">成功 {{ result.successCount }}</a-tag>
        <a-tag color="red" v-if="result.failedCount > 0">失败 {{ result.failedCount }}</a-tag>
        <a-tag>共 {{ result.totalCount }}</a-tag>
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
              {{ record.error ? '失败' : '成功' }}
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

const fileCols = [
  { title: '原文件名', dataIndex: 'originalName', key: 'originalName', ellipsis: true },
  { title: '新文件名（预览）', dataIndex: 'newName', key: 'newName', ellipsis: true },
  { title: '大小', dataIndex: 'size', key: 'size', width: 90 },
  { title: '操作', key: 'action', width: 70 }
]

const resultCols = [
  { title: '原文件名', dataIndex: 'originalName', key: 'originalName', ellipsis: true },
  { title: '新文件名', dataIndex: 'newName', key: 'newName', ellipsis: true },
  { title: '状态', key: 'status', width: 80 }
]

// 点击选择文件
const selectFiles = async () => {
  try {
    const paths = await RenameService.OpenFileDialogs()
    if (paths?.length) await addPaths(paths)
  } catch (e: any) {
    message.error('选择文件失败：' + (e?.message || e))
  }
}

// 添加路径（去重）
const addPaths = async (paths: string[]) => {
  const existing = new Set(files.value.map(f => f.originalPath))
  const newPaths = paths.filter(p => !existing.has(p))
  if (!newPaths.length) {
    message.info('所选文件已全部添加')
    return
  }
  try {
    const infos = await RenameService.BatchGetFileInfo(newPaths)
    files.value.push(...infos)
    previewList.value = []
    result.value = null
  } catch (e: any) {
    message.error('获取文件信息失败：' + (e?.message || e))
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
    message.success('预览完成')
  } catch (e: any) {
    message.error('预览失败：' + (e?.message || e))
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
        message.success(`全部重命名成功，共 ${res.successCount} 个文件`)
        clearFiles()
      } else {
        message.warning(`完成：${res.successCount} 成功，${res.failedCount} 失败`)
      }
    }
  } catch (e: any) {
    message.error('重命名失败：' + (e?.message || e))
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
