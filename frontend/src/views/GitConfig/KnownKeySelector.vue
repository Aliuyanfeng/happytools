<template>
  <a-modal
    :open="visible"
    :title="t('gitConfig.knownKeySelector')"
    width="640px"
    :footer="null"
    @cancel="emit('update:visible', false)"
  >
    <a-input
      v-model:value="keyword"
      :placeholder="t('gitConfig.knownKeySearch')"
      allow-clear
      class="mb-3"
    />

    <div class="selector-body">
      <!-- 左侧列表 -->
      <div class="key-list">
        <!-- 按节分组 -->
        <template v-for="(group, sec) in groupedKeys" :key="sec">
          <div class="group-label">{{ sec }}</div>
          <div
            v-for="k in group"
            :key="k.section + '.' + k.key"
            class="key-item"
            :class="{ active: selected?.key === k.key && selected?.section === k.section }"
            @click="selected = k"
          >
            <span class="key-name">{{ k.key }}</span>
            <span class="key-desc">{{ k.descZh }}</span>
          </div>
        </template>

        <!-- 无结果时显示自定义输入 -->
        <div v-if="Object.keys(groupedKeys).length === 0" class="no-result">
          {{ t('gitConfig.knownKeySearch') }}
        </div>
      </div>

      <!-- 右侧详情 -->
      <div class="key-detail">
        <template v-if="selected">
          <div class="detail-title">{{ selected.section }}.{{ selected.key }}</div>
          <a-descriptions :column="1" size="small" class="mb-3">
            <a-descriptions-item :label="t('gitConfig.knownKeyType')">
              <a-tag>{{ selected.type }}</a-tag>
            </a-descriptions-item>
            <a-descriptions-item v-if="selected.default" :label="t('gitConfig.knownKeyDefault')">
              {{ selected.default }}
            </a-descriptions-item>
            <a-descriptions-item v-if="selected.enumValues?.length" label="可选值">
              <a-tag v-for="v in selected.enumValues" :key="v" class="mr-1">{{ v }}</a-tag>
            </a-descriptions-item>
          </a-descriptions>
          <div class="desc-zh">{{ selected.descZh }}</div>
          <div class="desc-en">{{ selected.descEn }}</div>
        </template>
        <template v-else>
          <div class="detail-empty">选择左侧配置键查看详情</div>
        </template>
      </div>
    </div>

    <!-- 自定义键名 -->
    <a-divider class="my-3" />
    <div class="custom-key-row">
      <span class="mr-2 text-sm text-gray-500">{{ t('gitConfig.customKey') }}：</span>
      <a-input
        v-model:value="customKey"
        placeholder="例如: mykey"
        style="width: 200px"
        class="mr-2"
      />
      <a-button type="primary" :disabled="!confirmKey" @click="handleConfirm">
        确认
      </a-button>
    </div>
  </a-modal>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useGitConfigStore, type KnownKey } from '@/stores/gitconfig'

const props = defineProps<{
  visible: boolean
  section?: string // 预过滤节
}>()

const emit = defineEmits<{
  (e: 'update:visible', v: boolean): void
  (e: 'select', payload: { key: string; knownKey?: KnownKey }): void
}>()

const { t } = useI18n()
const store = useGitConfigStore()

const keyword = ref('')
const selected = ref<KnownKey | null>(null)
const customKey = ref('')

// 打开时加载 KnownKeys
watch(() => props.visible, async (v) => {
  if (v) {
    await store.loadKnownKeys()
    keyword.value = ''
    selected.value = null
    customKey.value = ''
  }
})

const filteredKeys = computed(() => {
  const kw = keyword.value.trim().toLowerCase()
  return store.knownKeys.filter(k => {
    const matchSection = !props.section || k.section === props.section
    if (!kw) return matchSection
    return matchSection && (
      k.key.toLowerCase().includes(kw) ||
      k.descZh.toLowerCase().includes(kw) ||
      k.descEn.toLowerCase().includes(kw)
    )
  })
})

const groupedKeys = computed(() => {
  const groups: Record<string, KnownKey[]> = {}
  for (const k of filteredKeys.value) {
    if (!groups[k.section]) groups[k.section] = []
    groups[k.section].push(k)
  }
  return groups
})

// 确认键：优先用选中的 KnownKey，其次用自定义输入
const confirmKey = computed(() => customKey.value.trim() || selected.value?.key || '')

function handleConfirm() {
  const key = customKey.value.trim() || selected.value?.key || ''
  if (!key) return
  emit('select', { key, knownKey: selected.value ?? undefined })
  emit('update:visible', false)
}
</script>

<style scoped>
.selector-body {
  display: flex;
  gap: 12px;
  height: 320px;
}

.key-list {
  flex: 1;
  overflow-y: auto;
  border: 1px solid #f0f0f0;
  border-radius: 6px;
  padding: 4px 0;
}

.group-label {
  padding: 4px 12px;
  font-size: 11px;
  font-weight: 600;
  color: #8c8c8c;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  background: #fafafa;
  border-bottom: 1px solid #f0f0f0;
}

.key-item {
  padding: 6px 12px;
  cursor: pointer;
  display: flex;
  flex-direction: column;
  gap: 2px;
  transition: background 0.15s;
}

.key-item:hover { background: #e6f7ff; }
.key-item.active { background: #1890ff; }
.key-item.active .key-name,
.key-item.active .key-desc { color: white; }

.key-name { font-size: 13px; font-weight: 500; color: #262626; }
.key-desc { font-size: 11px; color: #8c8c8c; }

.key-detail {
  width: 220px;
  border: 1px solid #f0f0f0;
  border-radius: 6px;
  padding: 12px;
  overflow-y: auto;
}

.detail-title {
  font-size: 14px;
  font-weight: 600;
  color: #1890ff;
  margin-bottom: 10px;
  word-break: break-all;
}

.desc-zh { font-size: 13px; color: #262626; margin-bottom: 6px; }
.desc-en { font-size: 12px; color: #8c8c8c; }
.detail-empty { color: #bfbfbf; font-size: 13px; text-align: center; margin-top: 40px; }
.no-result { padding: 12px; color: #bfbfbf; font-size: 13px; text-align: center; }

.custom-key-row { display: flex; align-items: center; }
</style>
