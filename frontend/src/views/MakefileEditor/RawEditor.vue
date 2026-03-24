<template>
  <div class="raw-editor">
    <div class="raw-editor__toolbar">
      <span class="raw-editor__label">{{ t('makefileEditor.rawMode') }}</span>
      <a-button type="primary" size="small" @click="handleApply">
        {{ t('makefileEditor.applyToVisual') }}
      </a-button>
    </div>
    <div class="raw-editor__body">
      <textarea
        ref="textareaRef"
        class="raw-editor__textarea"
        :value="localText"
        spellcheck="false"
        autocomplete="off"
        autocorrect="off"
        autocapitalize="off"
        @input="handleInput"
        @keydown.tab.prevent="handleTab"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import type { MakefileDoc } from '../../../bindings/github.com/Aliuyanfeng/happytools/backend/services/makefile/models.js'

const { t } = useI18n()

// ── Props & Emits ──────────────────────────────────────────────────────────
const props = defineProps<{
  modelValue: string
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
  (e: 'apply', text: string): void
}>()

// ── Local state ────────────────────────────────────────────────────────────
const localText = ref(props.modelValue)
const textareaRef = ref<HTMLTextAreaElement | null>(null)

// Sync external changes (visual mode edits) into the textarea — req 7.4
watch(
  () => props.modelValue,
  (newVal) => {
    if (newVal !== localText.value) {
      localText.value = newVal
    }
  },
  { immediate: true },
)

// ── Handlers ───────────────────────────────────────────────────────────────
function handleInput(event: Event) {
  const value = (event.target as HTMLTextAreaElement).value
  localText.value = value
  emit('update:modelValue', value)
}

// Insert a real Tab character instead of moving focus — req 4.7 / usability
function handleTab(event: KeyboardEvent) {
  const el = event.target as HTMLTextAreaElement
  const start = el.selectionStart
  const end = el.selectionEnd
  const newValue = localText.value.substring(0, start) + '\t' + localText.value.substring(end)
  localText.value = newValue
  emit('update:modelValue', newValue)
  // Restore cursor position after Vue updates the DOM
  requestAnimationFrame(() => {
    el.selectionStart = start + 1
    el.selectionEnd = start + 1
  })
}

// Emit apply event so the parent can re-parse and refresh the store — req 7.2
function handleApply() {
  emit('apply', localText.value)
}


</script>

<style scoped>
.raw-editor {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: #1e1e1e;
  border-radius: 6px;
  overflow: hidden;
}

.raw-editor__toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 6px 12px;
  background: #2d2d2d;
  border-bottom: 1px solid #3c3c3c;
  flex-shrink: 0;
}

.raw-editor__label {
  font-size: 12px;
  color: #9cdcfe;
  font-family: monospace;
  font-weight: 600;
}

.raw-editor__body {
  flex: 1;
  overflow: hidden;
  position: relative;
}

.raw-editor__textarea {
  width: 100%;
  height: 100%;
  padding: 12px 16px;
  background: transparent;
  color: #d4d4d4;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 13px;
  line-height: 1.6;
  border: none;
  outline: none;
  resize: none;
  box-sizing: border-box;
  tab-size: 4;
  white-space: pre;
  overflow-wrap: normal;
  overflow-x: auto;
  overflow-y: auto;
}

.raw-editor__textarea::selection {
  background: #264f78;
}

.raw-editor__textarea::placeholder {
  color: #6a6a6a;
}
</style>
