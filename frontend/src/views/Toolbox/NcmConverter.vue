<template>
  <div class="ncm-page">
    <!-- 输出目录 -->
    <div class="section">
      <div class="section-label">
        输出目录
        <span v-if="isDefault" class="default-badge">默认</span>
      </div>
      <div class="dir-row">
        <div class="dir-display" :class="{ placeholder: !outputDir }">
          {{ outputDir || '默认输出到原文件所在目录' }}
        </div>
        <a-button size="small" @click="selectOutputDir">选择目录</a-button>
        <a-tooltip :title="isDefault ? '取消默认目录' : '设为默认目录'" v-if="outputDir">
          <a-button
            size="small"
            :type="isDefault ? 'primary' : 'default'"
            @click="toggleDefault"
          >
            {{ isDefault ? '取消默认' : '设为默认' }}
          </a-button>
        </a-tooltip>
        <a-button size="small" type="text" v-if="outputDir" @click="clearOutputDir" danger>清除</a-button>
      </div>
    </div>

    <!-- 输入模式切换 -->
    <div class="section">
      <div class="section-label">输入模式</div>
      <a-radio-group v-model:value="inputMode" button-style="solid" size="small">
        <a-radio-button value="files">选择文件</a-radio-button>
        <a-radio-button value="dir">选择目录</a-radio-button>
      </a-radio-group>
    </div>

    <!-- 拖拽 / 选择区域 -->
    <div
      class="drop-zone"
      :class="{ dragging: isDragging, 'has-files': files.length > 0 }"
      id="ncm-drop-zone"
      @dragover.prevent="isDragging = true"
      @dragleave="isDragging = false"
      @drop.prevent="onDrop"
      @click="inputMode === 'dir' ? selectDir() : selectFiles()"
    >
      <template v-if="files.length === 0">
        <div class="drop-icon">🎵</div>
        <template v-if="inputMode === 'files'">
          <div class="drop-text">拖拽 NCM 文件到此处</div>
          <div class="drop-hint">或点击选择文件，支持批量</div>
        </template>
        <template v-else>
          <div class="drop-text">拖拽目录到此处</div>
          <div class="drop-hint">或点击选择目录，自动扫描所有 NCM 文件</div>
        </template>
      </template>
      <template v-else>
        <div class="drop-icon small">🎵</div>
        <div class="drop-text small">已选 {{ files.length }} 个文件，可继续拖入</div>
      </template>
    </div>

    <!-- 文件列表 -->
    <div class="file-list" v-if="files.length > 0">
      <div class="list-header">
        <span class="list-title">待转换文件</span>
        <a-button type="text" size="small" @click="clearFiles" :disabled="converting" class="clear-all">清空</a-button>
      </div>
      <div class="file-items">
        <div v-for="(f, i) in files" :key="i" class="file-item">
          <div class="file-info">
            <span class="file-name">{{ basename(f.path) }}</span>
            <span class="file-dir">{{ dirname(f.path) }}</span>
          </div>
          <div class="file-status">
            <span v-if="f.hasLrc" class="lrc-badge" title="存在歌词文件">🎵 LRC</span>
            <span v-if="f.status === 'pending'" class="status-pending">等待中</span>
            <span v-else-if="f.status === 'converting'" class="status-converting">
              <LoadingOutlined /> 转换中
            </span>
            <span v-else-if="f.status === 'done'" class="status-done">
              <CheckCircleOutlined /> 完成
            </span>
            <span v-else-if="f.status === 'error'" class="status-error" :title="f.error">
              <CloseCircleOutlined /> 失败
            </span>
          </div>
          <a-button
            v-if="!converting"
            type="text"
            size="small"
            class="remove-btn"
            @click.stop="removeFile(i)"
          >
            <CloseOutlined />
          </a-button>
        </div>
      </div>
    </div>

    <!-- 进度 & 操作 -->
    <div class="action-bar" v-if="files.length > 0">
      <div class="progress-wrap" v-if="converting || doneCount > 0">
        <a-progress
          :percent="progressPercent"
          :status="progressStatus"
          size="small"
        />
        <span class="progress-text">{{ doneCount + errorCount }} / {{ files.length }}</span>
      </div>
      <a-button
        type="primary"
        :loading="converting"
        :disabled="files.length === 0"
        @click="startConvert"
        class="convert-btn"
        size="large"
      >
        {{ converting ? '转换中...' : '开始转换' }}
      </a-button>
      <a-button
        v-if="doneCount > 0 && !converting"
        size="large"
        @click="openOutputDir"
        class="open-dir-btn"
      >
        <FolderOpenOutlined /> 打开输出目录
      </a-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { message } from 'ant-design-vue'
import { LoadingOutlined, CheckCircleOutlined, CloseCircleOutlined, CloseOutlined, FolderOpenOutlined } from '@ant-design/icons-vue'
import * as NCMService from '../../../bindings/github.com/Aliuyanfeng/happytools/backend/services/ncm/ncmservice'
import { Events } from '@wailsio/runtime'

interface FileItem {
  path: string
  status: 'pending' | 'converting' | 'done' | 'error'
  error?: string
  hasLrc?: boolean
}

const outputDir = ref('')
const files = ref<FileItem[]>([])
const converting = ref(false)
const isDragging = ref(false)
const inputMode = ref<'files' | 'dir'>('files')
const isDefault = ref(false)

const DEFAULT_DIR_KEY = 'ncm-default-output-dir'

function loadDefaultDir() {
  const saved = localStorage.getItem(DEFAULT_DIR_KEY)
  if (saved) {
    outputDir.value = saved
    isDefault.value = true
  }
}

function toggleDefault() {
  if (isDefault.value) {
    localStorage.removeItem(DEFAULT_DIR_KEY)
    isDefault.value = false
    message.success('已取消默认输出目录')
  } else {
    localStorage.setItem(DEFAULT_DIR_KEY, outputDir.value)
    isDefault.value = true
    message.success('已设为默认输出目录')
  }
}

function clearOutputDir() {
  // 清除目录时如果是默认目录，同时清除默认设置
  if (isDefault.value) {
    localStorage.removeItem(DEFAULT_DIR_KEY)
    isDefault.value = false
  }
  outputDir.value = ''
}

const doneCount = computed(() => files.value.filter(f => f.status === 'done').length)
const errorCount = computed(() => files.value.filter(f => f.status === 'error').length)
const progressPercent = computed(() => {
  if (files.value.length === 0) return 0
  return Math.round((doneCount.value + errorCount.value) / files.value.length * 100)
})
const progressStatus = computed(() => {
  if (converting.value) return 'active'
  if (errorCount.value > 0 && doneCount.value === 0) return 'exception'
  return 'normal'
})

function basename(p: string) {
  return p.replace(/\\/g, '/').split('/').pop() || p
}

function dirname(p: string) {
  const parts = p.replace(/\\/g, '/').split('/')
  parts.pop()
  return parts.join('/') || p
}

async function addPaths(paths: string[]) {
  const ncmPaths = paths.filter(p => p.toLowerCase().endsWith('.ncm'))
  if (ncmPaths.length === 0) {
    message.warning('未找到 .ncm 格式文件')
    return
  }
  const existing = new Set(files.value.map(f => f.path))
  const newPaths = ncmPaths.filter(p => !existing.has(p))
  if (newPaths.length === 0) {
    message.info('所选文件已全部添加')
    return
  }

  // 批量检查 lrc
  let lrcMap: Record<string, string> = {}
  try {
    lrcMap = await NCMService.CheckLrcBatch(newPaths) ?? {}
  } catch {}

  const newItems: FileItem[] = newPaths.map(p => ({
    path: p,
    status: 'pending' as const,
    hasLrc: !!(lrcMap[p])
  }))
  files.value.push(...newItems)

  if (newPaths.length < ncmPaths.length) {
    message.info(`已过滤 ${ncmPaths.length - newPaths.length} 个重复文件`)
  }
  message.success(`已添加 ${newItems.length} 个文件`)
}

function onDrop(e: DragEvent) {
  isDragging.value = false
  // Wails 桌面环境下 dataTransfer.files 没有真实路径，路径由 wails:file-drop 事件提供
}

async function selectFiles() {
  try {
    const paths = await NCMService.SelectFiles()
    if (paths?.length) await addPaths(paths)
  } catch {}
}

async function selectDir() {
  try {
    const paths = await NCMService.SelectDir()
    if (paths?.length) await addPaths(paths)
    else message.warning('该目录下未找到 NCM 文件')
  } catch {}
}

async function selectOutputDir() {
  try {
    const dir = await NCMService.SelectOutputDir()
    if (dir) {
      outputDir.value = dir
      // 检查是否和已保存的默认目录一致
      const saved = localStorage.getItem(DEFAULT_DIR_KEY)
      isDefault.value = saved === dir
    }
  } catch {}
}

function removeFile(i: number) {
  files.value.splice(i, 1)
}

function clearFiles() {
  files.value = []
}

async function openOutputDir() {
  // 优先用设置的输出目录，否则取第一个成功文件的目录
  let dir = outputDir.value
  if (!dir) {
    const done = files.value.find(f => f.status === 'done')
    if (done) {
      dir = done.path.replace(/\\/g, '/').split('/').slice(0, -1).join('/')
    }
  }
  if (!dir) { message.warning('没有可打开的目录'); return }
  try {
    await NCMService.OpenOutputDir(dir)
  } catch (e: any) {
    message.error('打开目录失败: ' + e.message)
  }
}

async function startConvert() {
  if (converting.value || files.value.length === 0) return
  converting.value = true

  // 重置所有状态为 pending
  files.value.forEach(f => { f.status = 'pending'; f.error = undefined })

  let doneN = 0
  let failN = 0

  for (const f of files.value) {
    f.status = 'converting'
    try {
      const r = await NCMService.ConvertOne(f.path, outputDir.value)
      if (r?.success) {
        f.status = 'done'
        doneN++
      } else {
        f.status = 'error'
        f.error = r?.error ?? '未知错误'
        failN++
      }
    } catch (e: any) {
      f.status = 'error'
      f.error = e?.message ?? '未知错误'
      failN++
    }
  }

  converting.value = false

  if (failN === 0) {
    message.success(`全部转换完成，共 ${doneN} 个文件`)
  } else {
    message.warning(`完成 ${doneN} 个，失败 ${failN} 个`)
  }
}

// 监听 Wails 文件拖拽事件（参考 PNGInjector 实现）
let offFileDrop: (() => void) | null = null

onMounted(() => {
  loadDefaultDir()
  offFileDrop = Events.On('wails:file-drop', async (event: any) => {
    const data = event?.data
    const dropped: string[] = Array.isArray(data?.files) ? data.files : []
    if (!dropped.length) return

    isDragging.value = false

    if (inputMode.value === 'dir') {
      for (const p of dropped) {
        try {
          const paths = await NCMService.ScanDir(p)
          if (paths?.length) await addPaths(paths)
          else message.warning(`目录 "${basename(p)}" 下未找到 NCM 文件`)
        } catch {
          await addPaths([p])
        }
      }
    } else {
      await addPaths(dropped)
    }
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
.ncm-page {
  padding: 28px 32px;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

/* 输出目录 */
.section-label {
  font-size: 13px;
  font-weight: 600;
  color: #555;
  margin-bottom: 8px;
  display: flex;
  align-items: center;
  gap: 6px;
}

.default-badge {
  font-size: 11px;
  font-weight: 500;
  color: #1677ff;
  background: #e6f4ff;
  padding: 1px 6px;
  border-radius: 4px;
}

.dir-row {
  display: flex;
  align-items: center;
  gap: 8px;
}

.dir-display {
  flex: 1;
  padding: 6px 12px;
  background: #f7f8fa;
  border: 1px solid #e8e8e8;
  border-radius: 8px;
  font-size: 13px;
  color: #333;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.dir-display.placeholder {
  color: #bbb;
}

/* 拖拽区 */
.drop-zone {
  border: 2px dashed #d9d9d9;
  border-radius: 14px;
  padding: 40px 24px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  transition: all 0.2s;
  background: #fafafa;
  user-select: none;
}

.drop-zone:hover,
.drop-zone.dragging {
  border-color: #1677ff;
  background: #e6f4ff;
}

.drop-zone.has-files {
  padding: 20px 24px;
}

.drop-icon {
  font-size: 40px;
  line-height: 1;
}

.drop-icon.small {
  font-size: 24px;
}

.drop-text {
  font-size: 15px;
  font-weight: 500;
  color: #333;
}

.drop-text.small {
  font-size: 13px;
  color: #666;
}

.drop-hint {
  font-size: 12px;
  color: #bbb;
}

/* 文件列表 */
.file-list {
  background: #fff;
  border: 1px solid #f0f0f0;
  border-radius: 12px;
  overflow: hidden;
}

.list-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 16px;
  border-bottom: 1px solid #f5f5f5;
}

.list-title {
  font-size: 13px;
  font-weight: 600;
  color: #555;
}

.clear-all {
  color: #bbb;
  font-size: 12px;
}
.clear-all:hover { color: #ff4d4f !important; }

.file-items {
  max-height: 280px;
  overflow-y: auto;
}

.file-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 16px;
  border-bottom: 1px solid #f9f9f9;
  transition: background 0.12s;
}

.file-item:last-child { border-bottom: none; }
.file-item:hover { background: #fafafa; }

.file-info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.file-name {
  font-size: 13px;
  color: #1a1a1a;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.file-dir {
  font-size: 11px;
  color: #bbb;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.file-status {
  font-size: 12px;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  gap: 6px;
}

.lrc-badge {
  font-size: 11px;
  color: #1677ff;
  background: #e6f4ff;
  padding: 1px 6px;
  border-radius: 4px;
  white-space: nowrap;
}

.status-pending  { color: #bbb; }
.status-converting { color: #1677ff; }
.status-done     { color: #52c41a; }
.status-error    { color: #ff4d4f; cursor: help; }

.remove-btn {
  color: #ccc;
  padding: 0;
  width: 24px;
  height: 24px;
  flex-shrink: 0;
}
.remove-btn:hover { color: #ff4d4f !important; }

/* 操作栏 */
.action-bar {
  display: flex;
  align-items: center;
  gap: 16px;
}

.progress-wrap {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 10px;
}

.progress-text {
  font-size: 12px;
  color: #999;
  white-space: nowrap;
}

.convert-btn {
  border-radius: 10px;
  padding: 0 32px;
  font-weight: 500;
  flex-shrink: 0;
}

.open-dir-btn {
  border-radius: 10px;
  flex-shrink: 0;
}
</style>
