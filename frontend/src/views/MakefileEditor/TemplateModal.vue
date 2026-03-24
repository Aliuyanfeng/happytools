<template>
  <a-modal
    :open="open"
    :title="t('makefileEditor.templates')"
    :footer="null"
    width="860px"
    @update:open="emit('update:open', $event)"
  >
    <div class="template-modal">
      <!-- Left: template list -->
      <div class="template-modal__list">
        <a-list
          :data-source="store.templates"
          :locale="{ emptyText: t('makefileEditor.noTemplates') }"
          size="small"
        >
          <template #renderItem="{ item }">
            <a-list-item
              class="template-modal__item"
              :class="{ 'template-modal__item--active': selectedId === item.id }"
              @click="selectTemplate(item)"
            >
              <div class="template-modal__item-content">
                <div class="template-modal__item-header">
                  <span class="template-modal__item-name">{{ item.name }}</span>
                  <a-tag v-if="item.isBuiltin" color="blue" size="small">
                    {{ t('makefileEditor.builtin') }}
                  </a-tag>
                  <a-tag v-else color="green" size="small">
                    {{ t('makefileEditor.custom') }}
                  </a-tag>
                </div>
                <div v-if="item.description" class="template-modal__item-desc">
                  {{ item.description }}
                </div>
              </div>
              <template #actions>
                <a-popconfirm
                  v-if="!item.isBuiltin"
                  :title="t('makefileEditor.deleteTemplateConfirm')"
                  :ok-text="t('common.confirm')"
                  :cancel-text="t('common.cancel')"
                  @confirm="handleDeleteTemplate(item.id)"
                >
                  <a-button type="text" danger size="small" @click.stop>
                    {{ t('common.delete') }}
                  </a-button>
                </a-popconfirm>
              </template>
            </a-list-item>
          </template>
        </a-list>
      </div>

      <!-- Right: preview + actions -->
      <div class="template-modal__detail">
        <template v-if="selectedTemplate">
          <!-- Preview -->
          <div class="template-modal__preview-label">{{ t('makefileEditor.templatePreview') }}</div>
          <pre class="template-modal__preview">{{ previewText }}</pre>

          <!-- Apply actions -->
          <div class="template-modal__actions">
            <a-button type="primary" @click="handleApply">
              {{ t('makefileEditor.applyTemplate') }}
            </a-button>
            <a-button @click="handleMerge">
              {{ t('makefileEditor.mergeTemplate') }}
            </a-button>
          </div>
        </template>
        <div v-else class="template-modal__empty">
          {{ t('makefileEditor.selectTemplateHint') }}
        </div>

        <a-divider />

        <!-- Save current as template -->
        <div class="template-modal__save-section">
          <div class="template-modal__save-label">{{ t('makefileEditor.saveAsTemplate') }}</div>
          <a-input
            v-model:value="saveName"
            :placeholder="t('makefileEditor.templateName')"
            class="template-modal__save-input"
          />
          <a-textarea
            v-model:value="saveDescription"
            :placeholder="t('makefileEditor.templateDescription')"
            :rows="2"
            class="template-modal__save-input"
          />
          <a-button
            type="default"
            :disabled="!saveName.trim() || !store.currentDoc"
            :loading="saving"
            @click="handleSaveAsTemplate"
          >
            {{ t('makefileEditor.saveAsTemplate') }}
          </a-button>
        </div>
      </div>
    </div>
  </a-modal>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { message } from 'ant-design-vue'
import { useMakefileEditorStore } from '../../stores/makefileEditor'
import {
  SaveCustomTemplate,
  DeleteCustomTemplate,
} from '../../../bindings/github.com/Aliuyanfeng/happytools/backend/services/makefile/makefileservice.js'
import type { Template } from '../../../bindings/github.com/Aliuyanfeng/happytools/backend/services/makefile/models.js'
import { serializeDoc } from './makefileSerializer'

const { t } = useI18n()
const store = useMakefileEditorStore()

// ── Props & Emits ──────────────────────────────────────────────────────────
defineProps<{ open: boolean }>()
const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
}>()

// ── State ──────────────────────────────────────────────────────────────────
const selectedId = ref<string>('')
const saveName = ref('')
const saveDescription = ref('')
const saving = ref(false)

// ── Computed ───────────────────────────────────────────────────────────────
const selectedTemplate = computed<Template | undefined>(() =>
  store.templates.find(t => t.id === selectedId.value),
)

const previewText = computed<string>(() => {
  if (!selectedTemplate.value) return ''
  return serializeDoc(selectedTemplate.value.doc)
})

// ── Lifecycle ──────────────────────────────────────────────────────────────
onMounted(() => {
  store.fetchTemplates()
})

// ── Handlers ──────────────────────────────────────────────────────────────
function selectTemplate(tmpl: Template) {
  selectedId.value = tmpl.id
}

function handleApply() {
  if (!selectedTemplate.value) return
  // Replace current doc with template doc (req 6.4)
  store.currentDoc = { ...selectedTemplate.value.doc }
  store.isDirty = true
  emit('update:open', false)
}

function handleMerge() {
  if (!selectedId.value) return
  // Append non-existing Targets and Variables from template (req 6.5)
  store.mergeTemplate(selectedId.value)
  emit('update:open', false)
}

async function handleSaveAsTemplate() {
  if (!saveName.value.trim() || !store.currentDoc) return
  saving.value = true
  try {
    await SaveCustomTemplate(saveName.value.trim(), saveDescription.value.trim(), store.currentDoc)
    await store.fetchTemplates()
    saveName.value = ''
    saveDescription.value = ''
    message.success(t('makefileEditor.saveTemplateSuccess'))
  } catch (err) {
    message.error(t('makefileEditor.saveTemplateFailed'))
  } finally {
    saving.value = false
  }
}

async function handleDeleteTemplate(id: string) {
  try {
    await DeleteCustomTemplate(id)
    await store.fetchTemplates()
    if (selectedId.value === id) {
      selectedId.value = ''
    }
    message.success(t('makefileEditor.deleteTemplateSuccess'))
  } catch (err) {
    message.error(t('makefileEditor.deleteTemplateFailed'))
  }
}
</script>

<style scoped>
.template-modal {
  display: flex;
  gap: 16px;
  min-height: 420px;
}

.template-modal__list {
  width: 260px;
  flex-shrink: 0;
  border: 1px solid var(--ant-color-border, #d9d9d9);
  border-radius: 6px;
  overflow-y: auto;
  max-height: 520px;
}

.template-modal__item {
  cursor: pointer;
  padding: 8px 12px !important;
  transition: background 0.15s;
}

.template-modal__item:hover {
  background: var(--ant-color-bg-text-hover, #f5f5f5);
}

.template-modal__item--active {
  background: var(--ant-color-primary-bg, #e6f4ff);
}

.template-modal__item-content {
  flex: 1;
  min-width: 0;
}

.template-modal__item-header {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-wrap: wrap;
}

.template-modal__item-name {
  font-weight: 500;
  font-size: 13px;
}

.template-modal__item-desc {
  font-size: 12px;
  color: var(--ant-color-text-secondary, #8c8c8c);
  margin-top: 2px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.template-modal__detail {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.template-modal__preview-label {
  font-size: 12px;
  font-weight: 600;
  color: var(--ant-color-text-secondary, #8c8c8c);
  margin-bottom: 6px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.template-modal__preview {
  flex: 1;
  background: #1e1e1e;
  color: #d4d4d4;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 12px;
  line-height: 1.6;
  padding: 12px;
  border-radius: 6px;
  overflow: auto;
  max-height: 240px;
  white-space: pre;
  margin-bottom: 12px;
}

.template-modal__actions {
  display: flex;
  gap: 8px;
  margin-bottom: 4px;
}

.template-modal__empty {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--ant-color-text-quaternary, #bfbfbf);
  font-size: 13px;
}

.template-modal__save-section {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.template-modal__save-label {
  font-size: 13px;
  font-weight: 600;
  color: var(--ant-color-text, #262626);
}

.template-modal__save-input {
  width: 100%;
}
</style>
