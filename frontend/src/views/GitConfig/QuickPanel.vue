<template>
  <div class="quick-panel">
    <div class="panel-header">
      <span class="panel-title">
        <ThunderboltOutlined class="mr-1" />
        {{ t('gitConfig.quickPanel') }}
      </span>
      <div class="panel-actions">
        <a-button size="small" type="text" @click="selectorVisible = true">
          <template #icon><PlusOutlined /></template>
          {{ t('gitConfig.addQuickItem') }}
        </a-button>
        <a-button size="small" type="text" @click="editMode = !editMode">
          <template #icon><EditOutlined /></template>
        </a-button>
      </div>
    </div>

    <div class="panel-body">
      <draggable
        v-model="localItems"
        item-key="order"
        handle=".drag-handle"
        @end="onDragEnd"
        class="items-grid"
      >
        <template #item="{ element: item, index }">
          <div class="quick-item">
            <DragOutlined v-if="editMode" class="drag-handle" />
            <div class="item-label">
              <span class="item-section">{{ itemLabel(item) }}</span>
            </div>
            <div class="item-value">
              <template v-if="editingIndex === index">
                <a-input
                  v-model:value="editingValue"
                  size="small"
                  style="width: 160px"
                  @blur="commitQuickEdit(item)"
                  @keyup.enter="commitQuickEdit(item)"
                  @keyup.esc="editingIndex = -1"
                  ref="quickInputRef"
                />
              </template>
              <span
                v-else
                class="value-text"
                :class="{ placeholder: !currentValue(item) }"
                @click="startQuickEdit(item, index)"
              >
                {{ currentValue(item) || t('gitConfig.notConfigured') }}
              </span>
            </div>
            <a-button
              v-if="editMode"
              size="small"
              type="text"
              danger
              @click="removeItem(index)"
            >
              <template #icon><CloseOutlined /></template>
            </a-button>
          </div>
        </template>
      </draggable>

      <div v-if="localItems.length === 0" class="empty-hint">
        {{ t('gitConfig.quickPanelHint') }}
      </div>
    </div>

    <KnownKeySelector
      v-model:visible="selectorVisible"
      @select="handleAddItem"
    />

    <!-- 添加 remote 快捷项时选择 remote 名称 -->
    <a-modal
      v-model:open="remotePickerVisible"
      title="选择 Remote"
      @ok="confirmRemoteItem"
      @cancel="remotePickerVisible = false"
    >
      <p class="mb-3 text-sm text-gray-500">该配置项属于 remote 节，请选择要关联的 remote 名称：</p>
      <a-select v-model:value="pendingRemoteName" style="width: 100%">
        <a-select-option v-for="r in store.remoteNames" :key="r" :value="r">{{ r }}</a-select-option>
      </a-select>
      <div v-if="store.remoteNames.length === 0" class="mt-2 text-xs text-gray-400">
        该仓库暂无 remote 配置
      </div>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import { PlusOutlined, EditOutlined, CloseOutlined, DragOutlined, ThunderboltOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import draggable from 'vuedraggable'
import { useGitConfigStore, type QuickPanelItem, type KnownKey } from '@/stores/gitconfig'
import KnownKeySelector from './KnownKeySelector.vue'

const { t } = useI18n()
const store = useGitConfigStore()

const editMode = ref(false)
const selectorVisible = ref(false)
const editingIndex = ref(-1)
const editingValue = ref('')
const quickInputRef = ref()

// remote 名称选择器状态
const remotePickerVisible = ref(false)
const pendingRemoteName = ref('')
const pendingRemoteKey = ref('')
const pendingRemoteKnownKey = ref<KnownKey | null>(null)

// 本地副本用于拖拽
const localItems = ref<QuickPanelItem[]>([])

watch(() => store.quickPanel, (v) => {
  localItems.value = [...v].sort((a, b) => a.order - b.order)
}, { immediate: true, deep: true })

function itemLabel(item: QuickPanelItem): string {
  if (item.subKey) return `${item.section}["${item.subKey}"].${item.key}`
  return `${item.section}.${item.key}`
}

function currentValue(item: QuickPanelItem): string {
  const sec = store.sections.find(
    s => s.name === item.section && s.subKey === item.subKey
  )
  return sec?.entries.find(e => e.key === item.key)?.value ?? ''
}

function startQuickEdit(item: QuickPanelItem, index: number) {
  editingIndex.value = index
  editingValue.value = currentValue(item)
  nextTick(() => quickInputRef.value?.focus?.())
}

async function commitQuickEdit(item: QuickPanelItem) {
  const original = currentValue(item)
  // 值未变化，直接取消
  if (editingValue.value === original) {
    editingIndex.value = -1
    return
  }
  // 新值为空且原值也为空（误点击），不写入
  if (editingValue.value.trim() === '' && original === '') {
    editingIndex.value = -1
    return
  }
  try {
    await store.saveEntry(item.section, item.subKey, item.key, editingValue.value)
  } catch (e: any) {
    message.error(e?.message ?? t('gitConfig.saveFailed'))
  }
  editingIndex.value = -1
}

function handleAddItem(payload: { key: string; knownKey?: KnownKey }) {
  const key = payload.key
  let section = '', entryKey = key

  if (payload.knownKey) {
    section = payload.knownKey.section
    entryKey = payload.knownKey.key
  } else {
    const parts = key.split('.')
    if (parts.length >= 2) {
      section = parts[0]
      entryKey = parts.slice(1).join('.')
    }
  }

  if (!section) return

  // remote 节需要先选择 remote 名称
  if (section === 'remote') {
    pendingRemoteKey.value = entryKey
    pendingRemoteKnownKey.value = payload.knownKey ?? null
    pendingRemoteName.value = store.remoteNames[0] ?? 'origin'
    remotePickerVisible.value = true
    return
  }

  const newItem: QuickPanelItem = {
    section,
    subKey: '',
    key: entryKey,
    order: localItems.value.length,
  }
  localItems.value.push(newItem)
  store.saveQuickPanel(localItems.value)
}

function confirmRemoteItem() {
  const name = pendingRemoteName.value.trim()
  if (!name) return
  const newItem: QuickPanelItem = {
    section: 'remote',
    subKey: name,
    key: pendingRemoteKey.value,
    order: localItems.value.length,
  }
  localItems.value.push(newItem)
  store.saveQuickPanel(localItems.value)
  remotePickerVisible.value = false
}

function removeItem(index: number) {
  localItems.value.splice(index, 1)
  localItems.value.forEach((item, i) => { item.order = i })
  store.saveQuickPanel(localItems.value)
}

function onDragEnd() {
  localItems.value.forEach((item, i) => { item.order = i })
  store.saveQuickPanel(localItems.value)
}
</script>

<style scoped>
.quick-panel {
  background: #fff;
  border: 1px solid #e8e8e8;
  border-radius: 8px;
  margin-bottom: 16px;
  overflow: hidden;
}

.panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 16px;
  background: linear-gradient(135deg, #e6f7ff 0%, #f0f5ff 100%);
  border-bottom: 1px solid #e8e8e8;
}

.panel-title {
  font-weight: 600;
  font-size: 14px;
  color: #1890ff;
}

.panel-actions { display: flex; gap: 4px; }

.panel-body { padding: 12px 16px; }

.items-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.quick-item {
  display: flex;
  align-items: center;
  gap: 8px;
  background: #f5f5f5;
  border: 1px solid #e8e8e8;
  border-radius: 6px;
  padding: 6px 10px;
  min-width: 200px;
}

.drag-handle {
  cursor: grab;
  color: #bfbfbf;
  font-size: 14px;
}

.item-label { flex-shrink: 0; }
.item-section {
  font-size: 12px;
  font-family: monospace;
  color: #8c8c8c;
}

.item-value { flex: 1; }

.value-text {
  font-size: 13px;
  color: #262626;
  cursor: pointer;
  padding: 2px 6px;
  border-radius: 4px;
  border: 1px solid transparent;
  transition: border-color 0.15s;
  display: inline-block;
}

.value-text:hover { border-color: #d9d9d9; background: #fff; }
.value-text.placeholder { color: #bfbfbf; font-style: italic; }

.empty-hint {
  color: #bfbfbf;
  font-size: 13px;
  text-align: center;
  padding: 12px 0;
}
</style>
