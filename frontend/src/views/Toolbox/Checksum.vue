<template>
  <div class="p-6">
    <a-card title="校验和计算">
      <!-- 拖拽/点击区域 -->
      <div
        id="checksum-drop-zone"
        data-file-drop-target
        class="drop-zone"
        :class="{ 'drop-zone--active': isDragging }"
        @click="handleSelectFiles"
      >
        <SafetyCertificateOutlined class="drop-icon" />
        <p class="mt-3 text-base text-gray-600">拖拽文件到此处，或点击选择文件</p>
        <p class="text-xs text-gray-400">支持任意类型文件，可多选</p>
      </div>

      <!-- 操作栏 -->
      <div v-if="fileList.length > 0" class="mt-4 mb-3 flex items-center justify-between">
        <span class="text-sm text-gray-500">共 {{ fileList.length }} 个文件</span>
        <a-space>
          <a-button size="small" danger @click="clearList">清空列表</a-button>
        </a-space>
      </div>

      <!-- 文件列表 -->
      <a-table
        v-if="fileList.length > 0"
        :columns="columns"
        :data-source="fileList"
        :loading="loading"
        :pagination="false"
        :scroll="{ y: 420 }"
        row-key="path"
        size="small"
        bordered
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'name'">
            <a-tooltip :title="record.path">
              <span class="font-medium">{{ record.name }}</span>
            </a-tooltip>
          </template>
          <template v-if="column.key === 'size'">
            {{ record.error ? '-' : formatBytes(record.size) }}
          </template>
          <template v-if="column.key === 'status'">
            <a-tag v-if="record.error" color="error">失败</a-tag>
            <a-tag v-else color="success">完成</a-tag>
          </template>
          <template v-if="column.key === 'hashes'">
            <div v-if="record.error" class="text-red-400 text-xs">{{ record.error }}</div>
            <div v-else class="hash-stack">
              <div v-for="item in hashItems(record)" :key="item.label" class="hash-row">
                <span class="hash-label">{{ item.label }}</span>
                <span class="hash-val font-mono">{{ item.value }}</span>
                <span class="copy-btn" @click="copyText(item.value)" title="复制">
                  <CopyOutlined />
                </span>
              </div>
            </div>
          </template>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { message } from 'ant-design-vue'
import { SafetyCertificateOutlined, CopyOutlined } from '@ant-design/icons-vue'
import { Events, Clipboard } from '@wailsio/runtime'
import * as ChecksumService from '../../../bindings/github.com/Aliuyanfeng/happytools/backend/services/checksum/checksumservice'

interface FileChecksum {
  path: string
  name: string
  size: number
  md5: string
  sha1: string
  sha256: string
  crc32: string
  error: string
}

const fileList = ref<FileChecksum[]>([])
const loading = ref(false)
const isDragging = ref(false)

const columns = [
  { title: '文件名', dataIndex: 'name', key: 'name', width: 200, ellipsis: true },
  { title: '大小', dataIndex: 'size', key: 'size', width: 90 },
  { title: '状态', key: 'status', width: 70 },
  { title: '校验和', key: 'hashes' },
]

function hashItems(record: FileChecksum) {
  return [
    { label: 'MD5', value: record.md5 },
    { label: 'SHA1', value: record.sha1 },
    { label: 'SHA256', value: record.sha256 },
    { label: 'CRC32', value: record.crc32 },
  ]
}

async function copyText(text: string) {
  if (!text) return
  try {
    await Clipboard.SetText(text)
    message.success('已复制')
  } catch {
    navigator.clipboard.writeText(text).then(() => message.success('已复制'))
  }
}

const formatBytes = (n: number): string => {
  if (n === 0) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(n) / Math.log(1024))
  return (n / Math.pow(1024, i)).toFixed(i === 0 ? 0 : 2) + ' ' + units[i]
}

// 去重：已存在的路径不重复添加
const addFiles = async (paths: string[]) => {
  const existing = new Set(fileList.value.map((f) => f.path))
  const newPaths = paths.filter((p) => !existing.has(p))
  if (newPaths.length === 0) {
    message.info('所选文件已在列表中')
    return
  }
  loading.value = true
  try {
    const results: FileChecksum[] = await ChecksumService.Calculate(newPaths)
    fileList.value.push(...results)
  } catch (e: any) {
    message.error('计算失败：' + (e?.message || String(e)))
  } finally {
    loading.value = false
  }
}

const handleSelectFiles = async () => {
  try {
    const paths: string[] = await ChecksumService.SelectFiles()
    if (!paths || paths.length === 0) return
    await addFiles(paths)
  } catch (e: any) {
    message.error('选择文件失败：' + (e?.message || String(e)))
  }
}

const clearList = () => {
  fileList.value = []
}

let offFileDrop: (() => void) | null = null

onMounted(() => {
  offFileDrop = Events.On('wails:file-drop', (event: any) => {
    const data = event?.data
    const files: string[] = Array.isArray(data?.files) ? data.files : []
    const targetId: string = data?.target ?? ''
    if (targetId && targetId !== 'checksum-drop-zone') return
    if (files.length === 0) return
    addFiles(files)
  })
})

onUnmounted(() => {
  if (offFileDrop) {
    offFileDrop()
    offFileDrop = null
  }
})
</script>

<style scoped>
.drop-zone {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 160px;
  border: 2px dashed #d9d9d9;
  border-radius: 8px;
  cursor: pointer;
  transition: border-color 0.2s, background 0.2s;
  user-select: none;
}

.drop-zone:hover,
.drop-zone--active {
  border-color: #1890ff;
  background: #f0f8ff;
}

.drop-icon {
  font-size: 48px;
  color: #bfbfbf;
  transition: color 0.2s;
}

.drop-zone:hover .drop-icon,
.drop-zone--active .drop-icon {
  color: #1890ff;
}

.hash-stack {
  display: flex;
  flex-direction: column;
  gap: 5px;
  padding: 2px 0;
}

.hash-row {
  display: flex;
  align-items: center;
  gap: 8px;
  min-height: 22px;
}

.hash-label {
  font-size: 11px;
  font-weight: 600;
  color: #999;
  width: 48px;
  flex-shrink: 0;
  text-align: right;
}

.hash-val {
  flex: 1;
  font-size: 12px;
  color: #1a1a1a;
  word-break: break-all;
  line-height: 1.4;
}

.copy-btn {
  cursor: pointer;
  color: #bbb;
  flex-shrink: 0;
  font-size: 12px;
  transition: color 0.15s;
  padding: 2px;
  border-radius: 3px;
}

.copy-btn:hover {
  color: #1677ff;
  background: #e6f4ff;
}

:deep(.hash-cell) {
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

:deep(.hash-text) {
  word-break: break-all;
}

:deep(.copy-btn) {
  cursor: pointer;
  color: #8c8c8c;
  flex-shrink: 0;
  font-size: 12px;
  transition: color 0.2s;
}

:deep(.copy-btn:hover) {
  color: #1890ff;
}</style>
