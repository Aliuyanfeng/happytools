<template>
  <div class="section-card">
    <!-- 节头 -->
    <div class="section-header">
      <div class="section-title">
        <span class="section-name">{{ sectionLabel }}</span>
        <a-tag v-if="section.entries.length" size="small" class="ml-2">{{ section.entries.length }}</a-tag>
      </div>
      <div class="section-actions">
        <a-button size="small" type="text" @click="selectorVisible = true">
          <template #icon><PlusOutlined /></template>
          {{ t('gitConfig.addEntry') }}
        </a-button>
        <a-popconfirm
          :title="t('gitConfig.deleteSectionConfirm')"
          @confirm="handleDeleteSection"
        >
          <a-button size="small" type="text" danger>
            <template #icon><DeleteOutlined /></template>
          </a-button>
        </a-popconfirm>
      </div>
    </div>

    <!-- 配置项列表 -->
    <div class="entry-list">
      <div
        v-for="entry in section.entries"
        :key="entry.key"
        class="entry-row"
        :class="{ dimmed: isFiltered && !isMatch(entry) }"
      >
        <!-- 键名 -->
        <div class="entry-key">
          <a-tooltip v-if="getKnownKey(entry.key)" :title="knownKeyTooltip(entry.key)" placement="right">
            <span class="key-label known">
              <span class="known-dot" />
              {{ entry.key }}
            </span>
          </a-tooltip>
          <span v-else class="key-label">{{ entry.key }}</span>
        </div>

        <!-- 值（内联编辑） -->
        <div class="entry-value">
          <template v-if="editingKey === entry.key">
            <!-- enum 类型用 select -->
            <a-select
              v-if="getKnownKey(entry.key)?.type === 'enum'"
              v-model:value="editingValue"
              size="small"
              style="width: 200px"
              :options="enumOptions(entry.key)"
              @blur="commitEdit(entry)"
              @change="commitEdit(entry)"
              ref="editInputRef"
            />
            <a-input
              v-else
              v-model:value="editingValue"
              size="small"
              style="width: 200px"
              @blur="commitEdit(entry)"
              @keyup.enter="commitEdit(entry)"
              @keyup.esc="cancelEdit"
              ref="editInputRef"
            />
          </template>
          <span
            v-else
            class="value-text"
            :class="{ highlight: isFiltered && isMatchValue(entry) }"
            @click="startEdit(entry)"
          >
            <mark v-if="isFiltered && isMatchValue(entry)">{{ entry.value }}</mark>
            <template v-else>{{ entry.value || '—' }}</template>
          </span>
        </div>

        <!-- 删除 -->
        <div class="entry-action">
          <a-popconfirm :title="t('gitConfig.deleteEntryConfirm')" @confirm="handleDeleteEntry(entry.key)">
            <a-button size="small" type="text" danger>
              <template #icon><DeleteOutlined /></template>
            </a-button>
          </a-popconfirm>
        </div>
      </div>

      <div v-if="section.entries.length === 0" class="empty-entries">
        {{ t('common.loading') === '加载中...' ? '暂无配置项' : 'No entries' }}
      </div>
    </div>

    <!-- KnownKey 选择器 -->
    <KnownKeySelector
      v-model:visible="selectorVisible"
      :section="section.name"
      @select="handleAddEntry"
    />

    <!-- 新增 entry 值输入 Modal -->
    <a-modal
      v-model:open="addEntryModal"
      :title="`添加 ${pendingKey}`"
      @ok="confirmAddEntry"
      @cancel="addEntryModal = false"
    >
      <a-select
        v-if="pendingKnownKey?.type === 'enum'"
        v-model:value="pendingValue"
        style="width: 100%"
        :options="pendingKnownKey.enumValues.map(v => ({ label: v, value: v }))"
      />
      <a-input v-else v-model:value="pendingValue" :placeholder="`请输入 ${pendingKey} 的值`" />
      <div v-if="pendingKnownKey" class="mt-2 text-xs text-gray-400">{{ pendingKnownKey.descZh }}</div>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import { PlusOutlined, DeleteOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { useGitConfigStore, type ConfigSection, type ConfigEntry, type KnownKey } from '@/stores/gitconfig'
import KnownKeySelector from './KnownKeySelector.vue'

const props = defineProps<{ section: ConfigSection }>()
const { t } = useI18n()
const store = useGitConfigStore()

const selectorVisible = ref(false)
const editingKey = ref<string | null>(null)
const editingValue = ref('')
const editInputRef = ref()

// 新增 entry 状态
const addEntryModal = ref(false)
const pendingKey = ref('')
const pendingValue = ref('')
const pendingKnownKey = ref<KnownKey | null>(null)

const sectionLabel = computed(() =>
  props.section.subKey ? `${props.section.name} "${props.section.subKey}"` : props.section.name
)

const isFiltered = computed(() => store.searchKeyword.trim() !== '')

function isMatch(entry: ConfigEntry) {
  const kw = store.searchKeyword.trim().toLowerCase()
  return entry.key.toLowerCase().includes(kw) || entry.value.toLowerCase().includes(kw)
}

function isMatchValue(entry: ConfigEntry) {
  return entry.value.toLowerCase().includes(store.searchKeyword.trim().toLowerCase())
}

function getKnownKey(key: string): KnownKey | undefined {
  return store.knownKeys.find(k => k.section === props.section.name && k.key === key)
}

function knownKeyTooltip(key: string): string {
  const k = getKnownKey(key)
  if (!k) return ''
  return `${k.descZh}\n${k.descEn}`
}

function enumOptions(key: string) {
  const k = getKnownKey(key)
  return (k?.enumValues ?? []).map(v => ({ label: v, value: v }))
}

function startEdit(entry: ConfigEntry) {
  editingKey.value = entry.key
  editingValue.value = entry.value
  nextTick(() => editInputRef.value?.focus?.())
}

function cancelEdit() {
  editingKey.value = null
  editingValue.value = ''
}

async function commitEdit(entry: ConfigEntry) {
  // 值未变化，直接取消
  if (editingValue.value === entry.value) {
    cancelEdit()
    return
  }
  // 新值为空且原值也为空（误点击），不写入
  if (editingValue.value.trim() === '' && entry.value === '') {
    cancelEdit()
    return
  }
  try {
    await store.saveEntry(props.section.name, props.section.subKey, entry.key, editingValue.value)
  } catch (e: any) {
    message.error(e?.message ?? t('gitConfig.saveFailed'))
  }
  cancelEdit()
}

function handleAddEntry(payload: { key: string; knownKey?: KnownKey }) {
  pendingKey.value = payload.key
  pendingKnownKey.value = payload.knownKey ?? null
  pendingValue.value = payload.knownKey?.default ?? ''
  addEntryModal.value = true
}

async function confirmAddEntry() {
  try {
    await store.saveEntry(props.section.name, props.section.subKey, pendingKey.value, pendingValue.value)
    message.success(t('gitConfig.saveSuccess'))
  } catch (e: any) {
    message.error(e?.message ?? t('gitConfig.saveFailed'))
  }
  addEntryModal.value = false
}

async function handleDeleteEntry(key: string) {
  try {
    await store.deleteEntry(props.section.name, props.section.subKey, key)
  } catch (e: any) {
    message.error(e?.message ?? t('gitConfig.saveFailed'))
  }
}

async function handleDeleteSection() {
  try {
    await store.deleteSection(props.section.name, props.section.subKey)
  } catch (e: any) {
    message.error(e?.message ?? t('gitConfig.saveFailed'))
  }
}
</script>

<style scoped>
.section-card {
  background: #fff;
  border: 1px solid #f0f0f0;
  border-radius: 8px;
  margin-bottom: 12px;
  overflow: hidden;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 16px;
  background: #fafafa;
  border-bottom: 1px solid #f0f0f0;
}

.section-title { display: flex; align-items: center; }
.section-name { font-weight: 600; font-size: 14px; color: #262626; font-family: monospace; }
.section-actions { display: flex; gap: 4px; }

.entry-list { padding: 4px 0; }

.entry-row {
  display: flex;
  align-items: center;
  padding: 6px 16px;
  gap: 12px;
  transition: opacity 0.2s, background 0.15s;
}

.entry-row:hover { background: #f5f5f5; }
.entry-row.dimmed { opacity: 0.3; }

.entry-key { width: 180px; flex-shrink: 0; }

.key-label {
  font-family: monospace;
  font-size: 13px;
  color: #595959;
  cursor: default;
}

.key-label.known { color: #1890ff; }

.known-dot {
  display: inline-block;
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #1890ff;
  margin-right: 5px;
  vertical-align: middle;
}

.entry-value { flex: 1; }

.value-text {
  font-size: 13px;
  color: #262626;
  cursor: pointer;
  padding: 2px 6px;
  border-radius: 4px;
  border: 1px solid transparent;
  transition: border-color 0.15s;
  display: inline-block;
  min-width: 40px;
}

.value-text:hover { border-color: #d9d9d9; background: #fff; }

.entry-action { width: 32px; flex-shrink: 0; }

.empty-entries {
  padding: 12px 16px;
  color: #bfbfbf;
  font-size: 13px;
}
</style>
