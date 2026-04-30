<template>
  <div class="p-6 json-formatter">
    <a-card title="JSON 格式化" class="main-card">
      <!-- 工具栏 -->
      <div class="toolbar">
        <a-space wrap>
          <a-button type="primary" @click="handleFormat">
            <template #icon><FormatPainterOutlined /></template>
            格式化
          </a-button>
          <a-button @click="handleCompress">
            <template #icon><CompressOutlined /></template>
            压缩
          </a-button>
          <a-button @click="handleEscape">
            <template #icon><CodeOutlined /></template>
            转义
          </a-button>
          <a-button @click="handleUnescape">
            <template #icon><CodeOutlined /></template>
            去转义
          </a-button>
          <a-divider type="vertical" />
          <a-select v-model:value="indentSize" style="width: 120px" @change="handleFormat">
            <a-select-option :value="2">缩进 2 空格</a-select-option>
            <a-select-option :value="4">缩进 4 空格</a-select-option>
            <a-select-option :value="0">Tab 缩进</a-select-option>
          </a-select>
          <!-- 主题选择 -->
          <a-select
            v-model:value="selectedTheme"
            style="width: 160px"
            @change="handleThemeChange"
          >
            <a-select-opt-group v-for="group in themeGroups" :key="group.label" :label="group.label">
              <a-select-option v-for="t in group.themes" :key="t.value" :value="t.value">
                {{ t.label }}
              </a-select-option>
            </a-select-opt-group>
          </a-select>
          <a-divider type="vertical" />
          <a-button @click="handleClear" danger>
            <template #icon><ClearOutlined /></template>
            清空
          </a-button>
          <a-button @click="handleCopy">
            <template #icon><CopyOutlined /></template>
            复制结果
          </a-button>
        </a-space>

        <!-- 状态提示 -->
        <div class="status-bar">
          <a-tag v-if="status === 'valid'" color="success">
            <CheckCircleOutlined /> JSON 合法
          </a-tag>
          <a-tag v-else-if="status === 'error'" color="error">
            <CloseCircleOutlined /> {{ errorMsg }}
          </a-tag>
          <span v-else class="status-idle">输入 JSON 后点击格式化</span>
        </div>
      </div>

      <!-- 编辑区：左输入 右输出 -->
      <div class="editor-layout">
        <!-- 输入区 -->
        <div class="editor-pane">
          <div class="pane-header">
            <span class="pane-title">输入</span>
            <a-button size="small" type="text" @click="handlePasteFromClipboard">
              <template #icon><SnippetsOutlined /></template>
              粘贴
            </a-button>
          </div>
          <div class="editor-wrap" ref="inputWrapRef"></div>
        </div>

        <!-- 分隔线 + 箭头 -->
        <div class="editor-divider">
          <ArrowRightOutlined class="divider-arrow" />
        </div>

        <!-- 输出区 -->
        <div class="editor-pane">
          <div class="pane-header">
            <span class="pane-title">输出</span>
            <a-tag v-if="outputStats" color="default" style="font-size:11px">
              {{ outputStats }}
            </a-tag>
          </div>
          <div class="editor-wrap" ref="outputWrapRef"></div>
        </div>
      </div>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { message } from 'ant-design-vue'
import {
  FormatPainterOutlined,
  CodeOutlined,
  ClearOutlined,
  CopyOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
  SnippetsOutlined,
  ArrowRightOutlined,
  CompressOutlined,
} from '@ant-design/icons-vue'
import { EditorView, keymap } from '@codemirror/view'
import { EditorState, Compartment } from '@codemirror/state'
import { json } from '@codemirror/lang-json'
import { oneDark } from '@codemirror/theme-one-dark'
import { defaultKeymap, indentWithTab, history, historyKeymap } from '@codemirror/commands'
import { indentOnInput, syntaxHighlighting, defaultHighlightStyle, bracketMatching } from '@codemirror/language'
import {
  lineNumbers, highlightActiveLineGutter, highlightSpecialChars,
  drawSelection, dropCursor, rectangularSelection, crosshairCursor, highlightActiveLine,
} from '@codemirror/view'
import { closeBrackets, closeBracketsKeymap } from '@codemirror/autocomplete'

// 按需导入主题（只导入常用的，避免包体积过大）
import { dracula }          from '@uiw/codemirror-theme-dracula'
import { atomone }          from '@uiw/codemirror-theme-atomone'
import { androidstudio }    from '@uiw/codemirror-theme-androidstudio'
import { darcula }          from '@uiw/codemirror-theme-darcula'
import { monokai }          from '@uiw/codemirror-theme-monokai'
import { nord }             from '@uiw/codemirror-theme-nord'
import { tokyoNight }       from '@uiw/codemirror-theme-tokyo-night'
import { tokyoNightStorm }  from '@uiw/codemirror-theme-tokyo-night-storm'
import { tokyoNightDay }    from '@uiw/codemirror-theme-tokyo-night-day'
import { vscodeDark, vscodeLight } from '@uiw/codemirror-theme-vscode'
import { githubDark, githubLight } from '@uiw/codemirror-theme-github'
import { sublime }          from '@uiw/codemirror-theme-sublime'
import { gruvboxDark }      from '@uiw/codemirror-theme-gruvbox-dark'
import { aura }             from '@uiw/codemirror-theme-aura'
import { copilot }          from '@uiw/codemirror-theme-copilot'
import { eclipse }          from '@uiw/codemirror-theme-eclipse'
import { quietlight }       from '@uiw/codemirror-theme-quietlight'
import { solarizedDark, solarizedLight } from '@uiw/codemirror-theme-solarized'
import { tomorrowNightBlue } from '@uiw/codemirror-theme-tomorrow-night-blue'
import { xcodeDark, xcodeLight } from '@uiw/codemirror-theme-xcode'
import { Extension } from '@codemirror/state'

// ── 主题注册表 ────────────────────────────────────────
const THEME_MAP: Record<string, Extension> = {
  'one-dark':            oneDark,
  'dracula':             dracula,
  'atom-one':            atomone,
  'android-studio':      androidstudio,
  'darcula':             darcula,
  'monokai':             monokai,
  'nord':                nord,
  'tokyo-night':         tokyoNight,
  'tokyo-night-storm':   tokyoNightStorm,
  'tokyo-night-day':     tokyoNightDay,
  'vscode-dark':         vscodeDark,
  'vscode-light':        vscodeLight,
  'github-dark':         githubDark,
  'github-light':        githubLight,
  'sublime':             sublime,
  'gruvbox-dark':        gruvboxDark,
  'aura':                aura,
  'copilot':             copilot,
  'tomorrow-night-blue': tomorrowNightBlue,
  'solarized-dark':      solarizedDark,
  'solarized-light':     solarizedLight,
  'eclipse':             eclipse,
  'quietlight':          quietlight,
  'xcode-dark':          xcodeDark,
  'xcode-light':         xcodeLight,
}

// 主题分组（用于下拉菜单）
const themeGroups = [
  {
    label: '深色',
    themes: [
      { value: 'one-dark',            label: 'One Dark' },
      { value: 'dracula',             label: 'Dracula' },
      { value: 'atom-one',            label: 'Atom One' },
      { value: 'android-studio',      label: 'Android Studio' },
      { value: 'darcula',             label: 'Darcula' },
      { value: 'monokai',             label: 'Monokai' },
      { value: 'nord',                label: 'Nord' },
      { value: 'tokyo-night',         label: 'Tokyo Night' },
      { value: 'tokyo-night-storm',   label: 'Tokyo Night Storm' },
      { value: 'vscode-dark',         label: 'VS Code Dark' },
      { value: 'github-dark',         label: 'GitHub Dark' },
      { value: 'sublime',             label: 'Sublime' },
      { value: 'gruvbox-dark',        label: 'Gruvbox Dark' },
      { value: 'aura',                label: 'Aura' },
      { value: 'copilot',             label: 'Copilot' },
      { value: 'tomorrow-night-blue', label: 'Tomorrow Night Blue' },
      { value: 'solarized-dark',      label: 'Solarized Dark' },
      { value: 'xcode-dark',          label: 'Xcode Dark' },
    ],
  },
  {
    label: '浅色',
    themes: [
      { value: 'tokyo-night-day',  label: 'Tokyo Night Day' },
      { value: 'vscode-light',     label: 'VS Code Light' },
      { value: 'github-light',     label: 'GitHub Light' },
      { value: 'eclipse',          label: 'Eclipse' },
      { value: 'quietlight',       label: 'Quiet Light' },
      { value: 'solarized-light',  label: 'Solarized Light' },
      { value: 'xcode-light',      label: 'Xcode Light' },
    ],
  },
]

// ── 状态 ──────────────────────────────────────────────
const indentSize = ref<number>(2)
const status = ref<'idle' | 'valid' | 'error'>('idle')
const errorMsg = ref('')
const outputStats = ref('')
const selectedTheme = ref<string>('one-dark')

// ── CodeMirror 实例 ───────────────────────────────────
const inputWrapRef = ref<HTMLElement | null>(null)
const outputWrapRef = ref<HTMLElement | null>(null)
let inputEditor: EditorView | null = null
let outputEditor: EditorView | null = null

// 用 Compartment 实现运行时动态切换主题
const themeCompartment = new Compartment()

const baseTheme = EditorView.theme({
  '&': { height: '100%', fontSize: '13px' },
  '.cm-scroller': { overflow: 'auto', fontFamily: '"JetBrains Mono", "Fira Code", monospace' },
  '.cm-content': { padding: '12px 0' },
})

const baseExtensions = [
  lineNumbers(),
  highlightActiveLineGutter(),
  highlightSpecialChars(),
  history(),
  drawSelection(),
  dropCursor(),
  rectangularSelection(),
  crosshairCursor(),
  highlightActiveLine(),
  bracketMatching(),
  closeBrackets(),
  syntaxHighlighting(defaultHighlightStyle, { fallback: true }),
  indentOnInput(),
  keymap.of([...defaultKeymap, ...historyKeymap, ...closeBracketsKeymap, indentWithTab]),
]

onMounted(() => {
  const initTheme = THEME_MAP[selectedTheme.value]
  inputEditor = new EditorView({
    state: EditorState.create({
      extensions: [...baseExtensions, json(), themeCompartment.of(initTheme), baseTheme],
    }),
    parent: inputWrapRef.value!,
  })
  outputEditor = new EditorView({
    state: EditorState.create({
      extensions: [
        ...baseExtensions, json(),
        themeCompartment.of(initTheme),
        baseTheme,
        EditorView.editable.of(false),
      ],
    }),
    parent: outputWrapRef.value!,
  })
})

onBeforeUnmount(() => {
  inputEditor?.destroy()
  outputEditor?.destroy()
})

// ── 主题切换 ──────────────────────────────────────────
const handleThemeChange = (val: string) => {
  const theme = THEME_MAP[val] ?? oneDark
  // 同时更新两个编辑器
  inputEditor?.dispatch({ effects: themeCompartment.reconfigure(theme) })
  outputEditor?.dispatch({ effects: themeCompartment.reconfigure(theme) })
  // 持久化到 localStorage
  localStorage.setItem('json-formatter-theme', val)
}

// ── 工具函数 ──────────────────────────────────────────
const getInput = () => inputEditor?.state.doc.toString() ?? ''

const setOutput = (text: string) => {
  outputEditor?.dispatch({
    changes: { from: 0, to: outputEditor.state.doc.length, insert: text },
  })
}

const getIndent = () => (indentSize.value === 0 ? '\t' : indentSize.value)

// ── 格式化 ────────────────────────────────────────────
const handleFormat = () => {
  const raw = getInput().trim()
  if (!raw) { status.value = 'idle'; setOutput(''); outputStats.value = ''; return }
  try {
    const parsed = JSON.parse(raw)
    const formatted = JSON.stringify(parsed, null, getIndent())
    setOutput(formatted)
    status.value = 'valid'
    errorMsg.value = ''
    const keys = countKeys(parsed)
    outputStats.value = `${formatted.split('\n').length} 行 · ${keys} 个键`
  } catch (e: any) {
    status.value = 'error'
    errorMsg.value = parseErrorMsg(e.message)
    setOutput('')
    outputStats.value = ''
  }
}

// ── 压缩 ──────────────────────────────────────────────
const handleCompress = () => {
  const raw = getInput().trim()
  if (!raw) return
  try {
    const parsed = JSON.parse(raw)
    const compressed = JSON.stringify(parsed)
    setOutput(compressed)
    status.value = 'valid'
    errorMsg.value = ''
    outputStats.value = `${compressed.length} 字符`
  } catch (e: any) {
    status.value = 'error'
    errorMsg.value = parseErrorMsg(e.message)
    setOutput('')
    outputStats.value = ''
  }
}

// ── 转义 ──────────────────────────────────────────────
const handleEscape = () => {
  const raw = getInput()
  if (!raw) return
  const escaped = raw
    .replace(/\\/g, '\\\\')
    .replace(/"/g, '\\"')
    .replace(/\n/g, '\\n')
    .replace(/\r/g, '\\r')
    .replace(/\t/g, '\\t')
  setOutput(escaped)
  status.value = 'idle'
  outputStats.value = ''
}

// ── 去转义 ────────────────────────────────────────────
const handleUnescape = () => {
  const raw = getInput()
  if (!raw) return
  try {
    const unescaped = JSON.parse(`"${raw.replace(/^"|"$/g, '')}"`)
    setOutput(unescaped)
    status.value = 'idle'
    outputStats.value = ''
  } catch {
    const unescaped = raw
      .replace(/\\n/g, '\n')
      .replace(/\\r/g, '\r')
      .replace(/\\t/g, '\t')
      .replace(/\\"/g, '"')
      .replace(/\\\\/g, '\\')
    setOutput(unescaped)
    status.value = 'idle'
    outputStats.value = ''
  }
}

// ── 清空 ──────────────────────────────────────────────
const handleClear = () => {
  inputEditor?.dispatch({
    changes: { from: 0, to: inputEditor.state.doc.length, insert: '' },
  })
  setOutput('')
  status.value = 'idle'
  errorMsg.value = ''
  outputStats.value = ''
}

// ── 粘贴 ──────────────────────────────────────────────
const handlePasteFromClipboard = async () => {
  try {
    const text = await navigator.clipboard.readText()
    inputEditor?.dispatch({
      changes: { from: 0, to: inputEditor.state.doc.length, insert: text },
    })
    handleFormat()
  } catch {
    message.error('读取剪贴板失败')
  }
}

// ── 复制结果 ──────────────────────────────────────────
const handleCopy = async () => {
  const text = outputEditor?.state.doc.toString() ?? ''
  if (!text) { message.warning('输出为空'); return }
  try {
    await navigator.clipboard.writeText(text)
    message.success('已复制到剪贴板')
  } catch {
    message.error('复制失败')
  }
}

// ── 辅助函数 ──────────────────────────────────────────
const countKeys = (obj: any): number => {
  if (typeof obj !== 'object' || obj === null) return 0
  let count = 0
  for (const key in obj) {
    count++
    if (typeof obj[key] === 'object' && obj[key] !== null) count += countKeys(obj[key])
  }
  return count
}

const parseErrorMsg = (msg: string): string => {
  const pos = msg.match(/position (\d+)/)
  if (pos) return `语法错误（位置 ${pos[1]}）`
  return '语法错误'
}

// ── 初始化：恢复上次选择的主题 ────────────────────────
const savedTheme = localStorage.getItem('json-formatter-theme')
if (savedTheme && THEME_MAP[savedTheme]) {
  selectedTheme.value = savedTheme
}
</script>

<style scoped>
.json-formatter {
  height: 100%;
  box-sizing: border-box;
}

.main-card {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.main-card :deep(.ant-card-body) {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 16px;
  overflow: hidden;
}

/* 工具栏 */
.toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 14px;
  padding-bottom: 12px;
  border-bottom: 1px solid #f0f0f0;
}

.status-bar {
  display: flex;
  align-items: center;
  min-width: 160px;
}

.status-idle {
  font-size: 12px;
  color: #bfbfbf;
}

/* 编辑区布局 */
.editor-layout {
  flex: 1;
  display: flex;
  gap: 0;
  overflow: hidden;
  border-radius: 8px;
  border: 1px solid #2d2d2d;
  background: #282c34;
}

.editor-pane {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  min-width: 0;
}

.pane-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 6px 12px;
  background: #21252b;
  border-bottom: 1px solid #2d2d2d;
  flex-shrink: 0;
}

.pane-title {
  font-size: 12px;
  font-weight: 600;
  color: #abb2bf;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.editor-wrap {
  flex: 1;
  overflow: hidden;
}

.editor-wrap :deep(.cm-editor) {
  height: 100%;
}

/* 中间分隔箭头 */
.editor-divider {
  width: 32px;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #21252b;
  border-left: 1px solid #2d2d2d;
  border-right: 1px solid #2d2d2d;
}

.divider-arrow {
  color: #528bff;
  font-size: 16px;
}
</style>
