<template>
  <div class="makefile-editor" ref="editorRef">
    <!-- Left sidebar -->
    <aside class="makefile-editor__sidebar" :style="{ width: leftWidth + 'px' }">
      <FilePanel @open-template-modal="templateModalVisible = true" />
      <a-divider style="margin: 8px 0" />
      <VariableList />
    </aside>

    <!-- Left resizer -->
    <div class="resizer" @mousedown="startResize('left', $event)">
      <div class="resizer-handle" />
    </div>

    <!-- Center panel -->
    <main class="makefile-editor__center">
      <div class="center-header">
        <span class="center-title">
          {{ isRawMode ? t('makefileEditor.rawMode') : t('makefileEditor.visualMode') }}
        </span>
        <a-button
          size="small"
          :type="isRawMode ? 'default' : 'primary'"
          @click="toggleMode"
        >
          {{ isRawMode ? t('makefileEditor.switchToVisual') : t('makefileEditor.switchToRaw') }}
        </a-button>
      </div>

      <div class="center-body">
        <DependencyGraph
          v-if="!isRawMode"
          @select-target="handleSelectTarget"
        />
        <RawEditor
          v-else
          v-model="rawText"
          @apply="handleRawApply"
        />
      </div>
    </main>

    <!-- Right resizer -->
    <div class="resizer" @mousedown="startResize('right', $event)">
      <div class="resizer-handle" />
    </div>

    <!-- Right panel -->
    <aside class="makefile-editor__right" :style="{ width: rightWidth + 'px' }">
      <template v-if="store.selectedTargetName || showAddTarget">
        <div class="right-header">
          <span class="section-title">
            {{ store.selectedTargetName
              ? t('makefileEditor.editTarget')
              : t('makefileEditor.addTarget') }}
          </span>
          <a-button size="small" type="text" @click="closeTargetForm">✕</a-button>
        </div>
        <TargetForm
          :initial-value="selectedTargetData ?? undefined"
          :existing-names="existingTargetNames"
          @submit="handleTargetFormSubmit"
          @cancel="closeTargetForm"
        />
      </template>
      <template v-else>
        <TargetList />
      </template>
    </aside>

    <!-- Template modal -->
    <TemplateModal v-model:open="templateModalVisible" />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { message } from 'ant-design-vue'
import { useMakefileEditorStore } from '../../stores/makefileEditor'
import type { Target } from '../../../bindings/github.com/Aliuyanfeng/happytools/backend/services/makefile/models.js'
import FilePanel from './FilePanel.vue'
import VariableList from './VariableList.vue'
import TargetList from './TargetList.vue'
import TargetForm from './TargetForm.vue'
import DependencyGraph from './DependencyGraph.vue'
import RawEditor from './RawEditor.vue'
import TemplateModal from './TemplateModal.vue'
import { serializeDoc } from './makefileSerializer'

const { t } = useI18n()
const store = useMakefileEditorStore()

// ── Mode toggle ────────────────────────────────────────────────────────────
const isRawMode = ref(false)
const rawText = ref('')

// Keep rawText in sync when visual mode changes (req 7.4)
watch(
  () => store.currentDoc,
  (doc) => {
    if (!isRawMode.value && doc) {
      rawText.value = serializeDoc(doc)
    }
  },
  { deep: true },
)

function toggleMode() {
  if (!isRawMode.value) {
    // Switching to raw: serialize current doc
    if (store.currentDoc) {
      rawText.value = serializeDoc(store.currentDoc)
    }
    isRawMode.value = true
  } else {
    // Switching back to visual: apply raw text first
    handleRawApply(rawText.value)
  }
}

async function handleRawApply(text: string) {
  // 有文件路径：保存到磁盘再重新解析
  if (store.currentPath) {
    try {
      const { SaveRawText, OpenFile } = await import(
        '../../../bindings/github.com/Aliuyanfeng/happytools/backend/services/makefile/makefileservice.js'
      )
      await SaveRawText(store.currentPath, text)
      const doc = await OpenFile(store.currentPath)
      store.currentDoc = doc
      store.isDirty = false
      isRawMode.value = false
    } catch (e: any) {
      message.error(e?.message ?? t('makefileEditor.applyToVisual'))
    }
    return
  }

  // 无文件路径（如模板新建未保存）：用后端解析原始文本
  if (store.currentDoc) {
    try {
      const { ParseRawText } = await import(
        '../../../bindings/github.com/Aliuyanfeng/happytools/backend/services/makefile/makefileservice.js'
      )
      const doc = await ParseRawText(text)
      store.currentDoc = doc
      store.isDirty = true
      isRawMode.value = false
    } catch (e: any) {
      message.error(e?.message ?? t('makefileEditor.applyToVisual'))
    }
    return
  }

  message.warning(t('makefileEditor.noFileOpen'))
}

// ── Right panel: TargetForm state ──────────────────────────────────────────
const showAddTarget = ref(false)

const selectedTargetData = computed<Target | null>(() => {
  if (!store.selectedTargetName || !store.currentDoc) return null
  return store.currentDoc.targets.find(t => t.name === store.selectedTargetName) ?? null
})

const existingTargetNames = computed<string[]>(() =>
  store.currentDoc?.targets.map(t => t.name) ?? [],
)

function handleSelectTarget(name: string) {
  store.setSelectedTarget(name)
  showAddTarget.value = false
}

function closeTargetForm() {
  store.setSelectedTarget('')
  showAddTarget.value = false
}

function handleTargetFormSubmit(target: Target) {
  if (store.selectedTargetName) {
    store.updateTarget(store.selectedTargetName, target)
    message.success(t('makefileEditor.editTarget') + ' ✓')
  } else {
    store.addTarget(target)
    message.success(t('makefileEditor.addTarget') + ' ✓')
  }
  closeTargetForm()
}

// ── Template modal ─────────────────────────────────────────────────────────
const templateModalVisible = ref(false)

// ── Resizable panels ───────────────────────────────────────────────────────
const editorRef = ref<HTMLElement>()
const leftWidth  = ref(260)
const rightWidth = ref(320)
const MIN_W = 160
const MAX_LEFT  = 480
const MAX_RIGHT = 560

let resizing: 'left' | 'right' | null = null
let startX = 0
let startW = 0

function startResize(side: 'left' | 'right', e: MouseEvent) {
  resizing = side
  startX = e.clientX
  startW = side === 'left' ? leftWidth.value : rightWidth.value
  document.body.style.cursor = 'col-resize'
  document.body.style.userSelect = 'none'
  window.addEventListener('mousemove', onMouseMove)
  window.addEventListener('mouseup', stopResize)
}

function onMouseMove(e: MouseEvent) {
  if (!resizing) return
  const delta = e.clientX - startX
  if (resizing === 'left') {
    leftWidth.value = Math.min(MAX_LEFT, Math.max(MIN_W, startW + delta))
  } else {
    rightWidth.value = Math.min(MAX_RIGHT, Math.max(MIN_W, startW - delta))
  }
}

function stopResize() {
  resizing = null
  document.body.style.cursor = ''
  document.body.style.userSelect = ''
  window.removeEventListener('mousemove', onMouseMove)
  window.removeEventListener('mouseup', stopResize)
}

onUnmounted(stopResize)
</script>

<style scoped>
.makefile-editor {
  display: flex;
  height: 100%;
  overflow: hidden;
  background: #f5f5f5;
  gap: 0;
}

/* Left sidebar */
.makefile-editor__sidebar {
  flex-shrink: 0;
  background: #fff;
  overflow-y: auto;
  padding: 12px;
  display: flex;
  flex-direction: column;
  min-width: 160px;
}

/* Center panel */
.makefile-editor__center {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  min-width: 0;
}

.center-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 16px;
  background: #fff;
  border-bottom: 1px solid #e8e8e8;
  flex-shrink: 0;
}

.center-title {
  font-size: 13px;
  font-weight: 600;
  color: #262626;
}

.center-body {
  flex: 1;
  overflow: hidden;
  padding: 12px;
}

/* Right panel */
.makefile-editor__right {
  flex-shrink: 0;
  background: #fff;
  overflow-y: auto;
  padding: 12px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  min-width: 160px;
}

.right-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 4px;
}

.section-title {
  font-size: 11px;
  font-weight: 600;
  color: #8c8c8c;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

/* ── Resizer ── */
.resizer {
  width: 6px;
  flex-shrink: 0;
  background: #e8e8e8;
  cursor: col-resize;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background 0.15s;
  position: relative;
  z-index: 10;
}

.resizer:hover {
  background: #1677ff;
}

.resizer-handle {
  width: 2px;
  height: 36px;
  border-radius: 2px;
  background: rgba(0, 0, 0, 0.2);
  transition: background 0.15s, height 0.15s;
  pointer-events: none;
}

.resizer:hover .resizer-handle {
  background: rgba(255, 255, 255, 0.9);
  height: 48px;
}
</style>