<template>
  <div class="file-panel">
    <!-- Action buttons -->
    <div class="panel-actions">
      <a-button block @click="handleOpenFile">
        <template #icon><FolderOpenOutlined /></template>
        {{ t('makefileEditor.openFile') }}
      </a-button>
      <a-button block @click="handleNewFile">
        <template #icon><FileAddOutlined /></template>
        {{ t('makefileEditor.newFile') }}
      </a-button>
      <a-button block @click="emit('open-template-modal')">
        <template #icon><AppstoreOutlined /></template>
        {{ t('makefileEditor.newFromTemplate') }}
      </a-button>
    </div>

    <!-- Current file -->
    <div class="section-title">{{ t('makefileEditor.currentFile') }}</div>
    <div class="current-file">
      <template v-if="store.currentPath">
        <FileTextOutlined class="file-icon" />
        <a-tooltip :title="store.currentPath">
          <span class="file-path">{{ shortPath(store.currentPath) }}</span>
        </a-tooltip>
        <a-button
          size="small"
          type="primary"
          :disabled="!store.isDirty"
          @click="handleSave"
          class="save-btn"
        >
          {{ t('makefileEditor.save') }}
        </a-button>
      </template>
      <span v-else class="no-file">—</span>
    </div>

    <!-- Recent files -->
    <div class="section-title">{{ t('makefileEditor.recentFiles') }}</div>
    <div class="recent-list">
      <div
        v-for="path in store.recentFiles"
        :key="path"
        class="recent-item"
        :class="{ active: path === store.currentPath }"
        @click="handleOpenRecent(path)"
      >
        <FileTextOutlined class="file-icon" />
        <a-tooltip :title="path">
          <span class="recent-path">{{ shortPath(path) }}</span>
        </a-tooltip>
      </div>
      <div v-if="store.recentFiles.length === 0" class="empty-hint">
        {{ t('makefileEditor.noRecentFiles') }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { FolderOpenOutlined, FileAddOutlined, AppstoreOutlined, FileTextOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { useMakefileEditorStore } from '../../stores/makefileEditor'
import {
  OpenFileDialog,
  OpenDirectoryDialog,
} from '../../../bindings/github.com/Aliuyanfeng/happytools/backend/services/makefile/makefileservice.js'

const { t } = useI18n()
const store = useMakefileEditorStore()

const emit = defineEmits<{
  (e: 'open-template-modal'): void
}>()

onMounted(() => {
  store.fetchRecentFiles()
})

function shortPath(path: string): string {
  const parts = path.replace(/\\/g, '/').split('/')
  if (parts.length <= 3) return path
  return '.../' + parts.slice(-2).join('/')
}

async function handleOpenFile() {
  try {
    const path = await OpenFileDialog()
    if (!path) return
    await store.loadFile(path)
  } catch (e: any) {
    message.error(e?.message ?? t('makefileEditor.openFile'))
  }
}

async function handleNewFile() {
  try {
    const dir = await OpenDirectoryDialog()
    if (!dir) return
    await store.newFile(dir)
  } catch (e: any) {
    message.error(e?.message ?? t('makefileEditor.newFile'))
  }
}

async function handleOpenRecent(path: string) {
  try {
    await store.loadFile(path)
  } catch (e: any) {
    message.error(e?.message ?? path)
  }
}

async function handleSave() {
  try {
    await store.saveFile()
  } catch (e: any) {
    message.error(e?.message ?? t('makefileEditor.save'))
  }
}
</script>

<style scoped>
.file-panel {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 12px;
  height: 100%;
  overflow-y: auto;
}

.panel-actions {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.section-title {
  font-size: 11px;
  font-weight: 600;
  color: #8c8c8c;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin-top: 8px;
}

.current-file {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 8px;
  background: #f5f5f5;
  border-radius: 6px;
  min-height: 32px;
}

.file-icon {
  color: #8c8c8c;
  flex-shrink: 0;
  font-size: 13px;
}

.file-path {
  flex: 1;
  font-size: 12px;
  color: #262626;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  font-family: monospace;
  cursor: default;
}

.save-btn {
  flex-shrink: 0;
}

.no-file {
  font-size: 12px;
  color: #bfbfbf;
}

.recent-list {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.recent-item {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 5px 8px;
  border-radius: 5px;
  cursor: pointer;
  transition: background 0.15s;
}

.recent-item:hover {
  background: #f0f0f0;
}

.recent-item.active {
  background: #e6f7ff;
}

.recent-path {
  font-size: 12px;
  color: #595959;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  font-family: monospace;
}

.empty-hint {
  font-size: 12px;
  color: #bfbfbf;
  padding: 8px;
  text-align: center;
}
</style>
