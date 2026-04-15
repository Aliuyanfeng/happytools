<template>
  <div class="nuclei-page">
    <!-- 左侧：输入区 -->
    <div class="left-pane">
      <div class="pane-header">
        <span class="pane-title">模板输入</span>
        <div class="header-actions">
          <a-select
            v-model:value="currentTheme"
            size="small"
            style="width: 140px"
            @change="changeTheme"
          >
            <a-select-option v-for="name in themeNames" :key="name" :value="name">
              {{ name }}
            </a-select-option>
          </a-select>
          <a-button size="small" @click="openFile"><FolderOpenOutlined /> 打开文件</a-button>
          <a-button size="small" @click="clearAll">清空</a-button>
        </div>
      </div>

      <!-- 拖拽区 / 编辑器 -->
      <div
        class="editor-wrap"
        :class="{ dragging: isDragging }"
        id="nuclei-drop-zone"
        @dragover.prevent="isDragging = true"
        @dragleave="isDragging = false"
      >
        <div ref="editorEl" class="cm-host" />
        <div v-if="!yamlContent" class="drop-hint">
          <div class="drop-icon">📄</div>
          <div>拖拽 .yaml 文件到此处</div>
          <div class="drop-sub">或点击上方「打开文件」按钮</div>
        </div>
      </div>

      <div class="parse-bar">
        <span v-if="parseError" class="parse-error"><CloseCircleOutlined /> {{ parseError }}</span>
        <span v-else-if="template" class="parse-ok"><CheckCircleOutlined /> 解析成功</span>
        <div class="parse-actions">
          <a-button size="small" @click="formatYaml" :disabled="!yamlContent.trim()" title="格式化 YAML">
            <AlignLeftOutlined /> 格式化
          </a-button>
          <a-button type="primary" :loading="parsing" @click="parse" :disabled="!yamlContent.trim()">
            解析模板
          </a-button>
        </div>
      </div>
    </div>

    <!-- 右侧：解析结果 -->
    <div class="right-pane" v-if="template">
      <!-- 基本信息 -->
      <div class="result-section">
        <div class="section-title">基本信息</div>
        <div class="info-grid">
          <div class="info-row">
            <span class="info-label">ID</span>
            <span class="info-value mono">{{ template.id }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">名称</span>
            <span class="info-value">{{ template.info.name }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">作者</span>
            <span class="info-value">{{ template.author_str || '-' }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">严重性</span>
            <span class="info-value">
              <span class="severity-badge" :class="template.info.severity">
                {{ template.info.severity || '-' }}
              </span>
            </span>
          </div>
          <div class="info-row">
            <span class="info-label">协议</span>
            <span class="info-value">
              <span class="proto-badge">{{ template.protocol }}</span>
            </span>
          </div>
          <div class="info-row" v-if="template.info.description">
            <span class="info-label">描述</span>
            <span class="info-value desc">{{ template.info.description }}</span>
          </div>
          <div class="info-row" v-if="template.info.remediation">
            <span class="info-label">修复建议</span>
            <span class="info-value desc">{{ template.info.remediation }}</span>
          </div>
          <div class="info-row" v-if="template.tag_list?.length">
            <span class="info-label">标签</span>
            <span class="info-value">
              <a-tag v-for="t in template.tag_list" :key="t" size="small" class="mr-1">{{ t }}</a-tag>
            </span>
          </div>
          <div class="info-row" v-if="template.info.cve_id">
            <span class="info-label">CVE</span>
            <span class="info-value mono">{{ template.info.cve_id }}</span>
          </div>
          <div class="info-row" v-if="template.info.cvss_score">
            <span class="info-label">CVSS 分数</span>
            <span class="info-value mono">{{ template.info.cvss_score }}</span>
          </div>
          <div class="info-row" v-if="template.ref_list?.length">
            <span class="info-label">参考链接</span>
            <div class="info-value">
              <div v-for="ref in template.ref_list" :key="ref" class="ref-link">{{ ref }}</div>
            </div>
          </div>
        </div>
      </div>

      <!-- HTTP 请求 -->
      <div class="result-section" v-if="template.http?.length">
        <div class="section-title">HTTP 请求 <span class="count-badge">{{ template.http.length }}</span></div>
        <div v-for="(req, i) in template.http" :key="i" class="request-block">
          <div class="req-header">
            <span class="req-index">#{{ i + 1 }}</span>
            <span class="req-method">{{ normalizeMethod(req.method) }}</span>
            <span v-if="req.attack" class="req-attack">{{ req.attack }}</span>
          </div>

          <!-- 路径 -->
          <div v-if="req.path?.length" class="req-group">
            <div class="req-group-label">Path</div>
            <div v-for="p in req.path" :key="p" class="code-line">{{ p }}</div>
          </div>

          <!-- Raw -->
          <div v-if="req.raw?.length" class="req-group">
            <div class="req-group-label">Raw</div>
            <pre v-for="(r, ri) in req.raw" :key="ri" class="raw-block">{{ r }}</pre>
          </div>

          <!-- Body -->
          <div v-if="req.body" class="req-group">
            <div class="req-group-label">Body</div>
            <pre class="raw-block">{{ req.body }}</pre>
          </div>

          <!-- Matchers -->
          <div v-if="req.matchers?.length" class="req-group">
            <div class="req-group-label">
              Matchers
              <span v-if="req.matchers_condition" class="condition-badge">{{ req.matchers_condition }}</span>
            </div>
            <div v-for="(m, mi) in req.matchers" :key="mi" class="matcher-item">
              <span class="matcher-type">{{ m.type }}</span>
              <span v-if="m.part" class="matcher-part">{{ m.part }}</span>
              <span v-if="m.negative" class="matcher-neg">NOT</span>
              <div class="matcher-values">
                <span v-for="w in m.words" :key="w" class="match-val">{{ w }}</span>
                <span v-for="r in m.regex" :key="r" class="match-val regex">{{ r }}</span>
                <span v-for="d in m.dsl" :key="d" class="match-val dsl">{{ d }}</span>
                <span v-for="s in m.status" :key="s" class="match-val status">{{ s }}</span>
              </div>
            </div>
          </div>

          <!-- Extractors -->
          <div v-if="req.extractors?.length" class="req-group">
            <div class="req-group-label">Extractors</div>
            <div v-for="(e, ei) in req.extractors" :key="ei" class="matcher-item">
              <span class="matcher-type">{{ e.type }}</span>
              <span v-if="e.name" class="matcher-part">{{ e.name }}</span>
              <span v-if="e.part" class="matcher-part">{{ e.part }}</span>
              <div class="matcher-values">
                <span v-for="r in e.regex" :key="r" class="match-val regex">{{ r }}</span>
                <span v-for="j in e.json" :key="j" class="match-val dsl">{{ j }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- DNS -->
      <div class="result-section" v-if="template.dns?.length">
        <div class="section-title">DNS 请求 <span class="count-badge">{{ template.dns.length }}</span></div>
        <div v-for="(req, i) in template.dns" :key="i" class="request-block">
          <div class="req-header">
            <span class="req-index">#{{ i + 1 }}</span>
            <span class="req-method">{{ req.type }}</span>
          </div>
          <div v-if="req.name" class="req-group">
            <div class="req-group-label">Name</div>
            <div class="code-line">{{ req.name }}</div>
          </div>
          <div v-if="req.matchers?.length" class="req-group">
            <div class="req-group-label">Matchers</div>
            <div v-for="(m, mi) in req.matchers" :key="mi" class="matcher-item">
              <span class="matcher-type">{{ m.type }}</span>
              <div class="matcher-values">
                <span v-for="w in m.words" :key="w" class="match-val">{{ w }}</span>
                <span v-for="r in m.regex" :key="r" class="match-val regex">{{ r }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Code -->
      <div class="result-section" v-if="template.code?.length">
        <div class="section-title">Code 请求 <span class="count-badge">{{ template.code.length }}</span></div>
        <div v-for="(req, i) in template.code" :key="i" class="request-block">
          <div class="req-header">
            <span class="req-index">#{{ i + 1 }}</span>
            <span v-for="eng in req.engine" :key="eng" class="req-method">{{ eng }}</span>
          </div>

          <!-- 源码类型判断 -->
          <div class="req-group">
            <div class="req-group-label">
              Source
              <span v-if="isExternalFile(req.source)" class="source-type-badge file">外部文件</span>
              <span v-else class="source-type-badge inline">内联代码</span>
            </div>
            <template v-if="isExternalFile(req.source)">
              <div class="code-line">{{ req.source }}</div>
            </template>
            <template v-else>
              <pre class="raw-block code-source">{{ req.source }}</pre>
            </template>
          </div>

          <!-- 参数 -->
          <div v-if="req.args?.length" class="req-group">
            <div class="req-group-label">Args</div>
            <div v-for="arg in req.args" :key="arg" class="code-line">{{ arg }}</div>
          </div>

          <!-- Matchers -->
          <div v-if="req.matchers?.length" class="req-group">
            <div class="req-group-label">
              Matchers
              <span v-if="req.matchers_condition" class="condition-badge">{{ req.matchers_condition }}</span>
            </div>
            <div v-for="(m, mi) in req.matchers" :key="mi" class="matcher-item">
              <span class="matcher-type">{{ m.type }}</span>
              <span v-if="m.part" class="matcher-part">{{ m.part }}</span>
              <span v-if="m.negative" class="matcher-neg">NOT</span>
              <div class="matcher-values">
                <span v-for="w in m.words" :key="w" class="match-val">{{ w }}</span>
                <span v-for="r in m.regex" :key="r" class="match-val regex">{{ r }}</span>
                <span v-for="d in m.dsl" :key="d" class="match-val dsl">{{ d }}</span>
              </div>
            </div>
          </div>

          <!-- Extractors -->
          <div v-if="req.extractors?.length" class="req-group">
            <div class="req-group-label">Extractors</div>
            <div v-for="(e, ei) in req.extractors" :key="ei" class="matcher-item">
              <span class="matcher-type">{{ e.type }}</span>
              <span v-if="e.name" class="matcher-part">{{ e.name }}</span>
              <div class="matcher-values">
                <span v-for="r in e.regex" :key="r" class="match-val regex">{{ r }}</span>
                <span v-for="j in e.json" :key="j" class="match-val dsl">{{ j }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Variables -->
      <div class="result-section" v-if="template.variables && Object.keys(template.variables).length">
        <div class="section-title">变量</div>
        <div class="info-grid">
          <div v-for="(val, key) in template.variables" :key="key" class="info-row">
            <span class="info-label mono">{{ key }}</span>
            <span class="info-value mono">{{ val }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 右侧空态 -->
    <div class="right-pane empty-pane" v-else>
      <div class="empty-state">
        <div class="empty-icon-wrap">🔍</div>
        <div class="empty-title">解析结果将在此显示</div>
        <div class="empty-sub">支持 HTTP / DNS / TCP 协议的 Nuclei 模板</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { message } from 'ant-design-vue'
import { FolderOpenOutlined, CheckCircleOutlined, CloseCircleOutlined, AlignLeftOutlined } from '@ant-design/icons-vue'
import { Events } from '@wailsio/runtime'
import * as NucleiService from '../../../bindings/github.com/Aliuyanfeng/happytools/backend/services/nuclei/nucleiservice'
import { EditorView, keymap, placeholder as cmPlaceholder } from '@codemirror/view'
import { EditorState, Compartment } from '@codemirror/state'
import { yaml } from '@codemirror/lang-yaml'
import { oneDark } from '@codemirror/theme-one-dark'
import { defaultKeymap, indentWithTab } from '@codemirror/commands'
import { indentOnInput, syntaxHighlighting, defaultHighlightStyle } from '@codemirror/language'
import * as YAML from 'yaml'
import {
  abcdef, androidstudio, atomone, aura, bbedit, bespin, darcula, dracula,
  duotoneDark, duotoneLight, eclipse, githubDark, githubLight, gruvboxDark,
  gruvboxLight, materialDark, materialLight, monokai, monokaiDimmed,
  noctisLilac, nord, okaidia, quietlight, solarizedDark,
  solarizedLight, sublime, tokyoNight, tokyoNightDay, tokyoNightStorm,
  tomorrowNightBlue, vscodeDark, vscodeLight, xcodeDark, xcodeLight,
  abyss, andromeda, basicDark, basicLight, copilot, kimbie, red,
  whiteDark, whiteLight, consoleDark, consoleLight
} from '@uiw/codemirror-themes-all'

// 主题列表
const THEMES: Record<string, any> = {
  'One Dark': oneDark,
  'VS Code Dark': vscodeDark,
  'VS Code Light': vscodeLight,
  'GitHub Dark': githubDark,
  'GitHub Light': githubLight,
  'Dracula': dracula,
  'Tokyo Night': tokyoNight,
  'Tokyo Night Storm': tokyoNightStorm,
  'Tokyo Night Day': tokyoNightDay,
  'Nord': nord,
  'Monokai': monokai,
  'Monokai Dimmed': monokaiDimmed,
  'Material Dark': materialDark,
  'Material Light': materialLight,
  'Gruvbox Dark': gruvboxDark,
  'Gruvbox Light': gruvboxLight,
  'Solarized Dark': solarizedDark,
  'Solarized Light': solarizedLight,
  'Aura': aura,
  'Darcula': darcula,
  'Atom One': atomone,
  'Android Studio': androidstudio,
  'Sublime': sublime,
  'Okaidia': okaidia,
  'Eclipse': eclipse,
  'BBEdit': bbedit,
  'Quiet Light': quietlight,
  'Xcode Dark': xcodeDark,
  'Xcode Light': xcodeLight,
  'Tomorrow Night Blue': tomorrowNightBlue,
  'Noctis Lilac': noctisLilac,
  'Duotone Dark': duotoneDark,
  'Duotone Light': duotoneLight,
  'Bespin': bespin,
  'Abcdef': abcdef,
  'Abyss': abyss,
  'Andromeda': andromeda,
  'Basic Dark': basicDark,
  'Basic Light': basicLight,
  'Copilot': copilot,
  'Kimbie': kimbie,
  'Red': red,
  'White Dark': whiteDark,
  'White Light': whiteLight,
  'Console Dark': consoleDark,
  'Console Light': consoleLight,
}

const themeNames = Object.keys(THEMES)
const currentTheme = ref('One Dark')
const themeCompartment = new Compartment()

const yamlContent = ref('')
const template = ref<any>(null)
const parsing = ref(false)
const parseError = ref('')
const isDragging = ref(false)
const editorEl = ref<HTMLElement | null>(null)
let editorView: EditorView | null = null

let debounceTimer: ReturnType<typeof setTimeout> | null = null

// 初始化 CodeMirror
function initEditor() {
  if (!editorEl.value) return

  const updateListener = EditorView.updateListener.of((update) => {
    if (update.docChanged) {
      yamlContent.value = update.state.doc.toString()
      onInput()
    }
  })

  const state = EditorState.create({
    doc: yamlContent.value,
    extensions: [
      keymap.of([...defaultKeymap, indentWithTab]),
      yaml(),
      syntaxHighlighting(defaultHighlightStyle),
      themeCompartment.of(THEMES[currentTheme.value]),
      indentOnInput(),
      cmPlaceholder('粘贴 Nuclei YAML 模板内容，或拖拽 .yaml 文件到此处...'),
      EditorView.lineWrapping,
      updateListener,
      EditorView.theme({
        '&': { height: '100%', fontSize: '12.5px' },
        '.cm-scroller': { overflow: 'auto', fontFamily: "'Consolas', 'Monaco', monospace" },
        '.cm-content': { padding: '12px 0' },
      }),
    ],
  })

  editorView = new EditorView({ state, parent: editorEl.value })
}

function changeTheme(name: string) {
  currentTheme.value = name
  if (!editorView) return
  editorView.dispatch({
    effects: themeCompartment.reconfigure(THEMES[name])
  })
  localStorage.setItem('nuclei-editor-theme', name)
}

// 从外部更新编辑器内容
function setEditorContent(content: string) {
  if (!editorView) return
  editorView.dispatch({
    changes: { from: 0, to: editorView.state.doc.length, insert: content }
  })
}

function onInput() {
  parseError.value = ''
  template.value = null
  if (debounceTimer) clearTimeout(debounceTimer)
  if (yamlContent.value.trim()) {
    debounceTimer = setTimeout(() => parse(), 800)
  }
}

function formatYaml() {
  if (!yamlContent.value.trim()) return
  try {
    const parsed = YAML.parse(yamlContent.value)
    const formatted = YAML.stringify(parsed, { indent: 2, lineWidth: 0 })
    setEditorContent(formatted)
    yamlContent.value = formatted
  } catch (e: any) {
    message.error('格式化失败，请检查 YAML 语法: ' + (e?.message || e))
  }
}

async function parse() {
  if (!yamlContent.value.trim()) return
  parsing.value = true
  parseError.value = ''
  try {
    const result = await NucleiService.ParseContent(yamlContent.value)
    template.value = result
  } catch (e: any) {
    parseError.value = e?.message || '解析失败'
    template.value = null
  } finally {
    parsing.value = false
  }
}

async function openFile() {
  try {
    const path = await NucleiService.OpenFile()
    if (!path) return
    const result = await NucleiService.ParseFile(path)
    template.value = result
    const content = result?.raw_yaml ?? ''
    yamlContent.value = content
    setEditorContent(content)
    parseError.value = ''
  } catch (e: any) {
    message.error('打开失败: ' + (e?.message || e))
  }
}

function clearAll() {
  yamlContent.value = ''
  template.value = null
  parseError.value = ''
}

function normalizeMethod(method: any): string {
  if (!method) return 'GET'
  if (typeof method === 'string') return method.toUpperCase()
  if (Array.isArray(method)) return method.map((m: string) => m.toUpperCase()).join(' / ')
  return String(method).toUpperCase()
}

// source 是外部文件路径（不含换行、以路径形式结尾）
function isExternalFile(source: string): boolean {
  if (!source) return false
  const s = source.trim()
  // 多行 = 内联代码
  if (s.includes('\n')) return false
  // 以常见脚本扩展名结尾 = 外部文件
  return /\.(py|sh|ps1|rb|js|lua|go|pl|php|bash|zsh)$/i.test(s)
}

// Wails 文件拖拽
let offDrop: (() => void) | null = null
onMounted(() => {
  const saved = localStorage.getItem('nuclei-editor-theme')
  if (saved && THEMES[saved]) currentTheme.value = saved
  initEditor()
  offDrop = Events.On('wails:file-drop', async (event: any) => {
    const data = event?.data
    const files: string[] = Array.isArray(data?.files) ? data.files : []
    const yaml = files.find(f => f.toLowerCase().endsWith('.yaml') || f.toLowerCase().endsWith('.yml'))
    if (!yaml) return
    isDragging.value = false
    try {
      const result = await NucleiService.ParseFile(yaml)
      template.value = result
      const content = result?.raw_yaml ?? ''
      yamlContent.value = content
      setEditorContent(content)
      parseError.value = ''
    } catch (e: any) {
      message.error('解析失败: ' + (e?.message || e))
    }
  })
})
onUnmounted(() => {
  offDrop?.()
  offDrop = null
  editorView?.destroy()
  editorView = null
})
</script>

<style scoped>
.nuclei-page {
  display: flex;
  height: calc(100vh - 120px);
  gap: 0;
  background: #f7f8fa;
  overflow: hidden;
  align-items: stretch;
}

/* ── 左侧 ── */
.left-pane {
  width: 420px;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  background: #fff;
  border-right: 1px solid #f0f0f0;
  overflow: hidden;
}

.pane-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 16px;
  border-bottom: 1px solid #f5f5f5;
  flex-shrink: 0;
}

.pane-title {
  font-size: 14px;
  font-weight: 600;
  color: #1a1a1a;
}

.header-actions {
  display: flex;
  gap: 6px;
}

.editor-wrap {
  flex: 1;
  position: relative;
  overflow: hidden;
  transition: background 0.15s;
}

.editor-wrap.dragging {
  background: #e6f4ff;
}

.cm-host {
  width: 100%;
  height: 100%;
}

.drop-hint {
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  pointer-events: none;
  color: #ccc;
  font-size: 13px;
}

.drop-icon { font-size: 36px; }
.drop-sub { font-size: 11px; color: #ddd; }

.parse-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 16px;
  border-top: 1px solid #f5f5f5;
  flex-shrink: 0;
  gap: 8px;
}

.parse-actions {
  display: flex;
  gap: 6px;
  flex-shrink: 0;
}

.parse-ok { font-size: 12px; color: #52c41a; display: flex; align-items: center; gap: 4px; }
.parse-error { font-size: 12px; color: #ff4d4f; display: flex; align-items: center; gap: 4px; flex: 1; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }

/* ── 右侧 ── */
.right-pane {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
  padding: 20px 24px;
  min-width: 0;
  align-self: stretch;
}

.right-pane > * + * {
  margin-top: 16px;
}

.empty-pane {
  display: flex;
  align-items: center;
  justify-content: center;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
  color: #ccc;
}

.empty-icon-wrap { font-size: 48px; }
.empty-title { font-size: 15px; color: #bbb; font-weight: 500; }
.empty-sub { font-size: 12px; color: #ccc; }

/* ── 结果区块 ── */
.result-section {
  background: #fff;
  border-radius: 10px;
  border: 1px solid #f0f0f0;
  overflow: hidden;
}

.section-title {
  font-size: 13px;
  font-weight: 600;
  color: #555;
  padding: 10px 16px;
  border-bottom: 1px solid #f5f5f5;
  display: flex;
  align-items: center;
  gap: 6px;
  background: #fafafa;
}

.count-badge {
  font-size: 11px;
  background: #e6f4ff;
  color: #1677ff;
  padding: 1px 6px;
  border-radius: 8px;
  font-weight: 500;
}

/* 信息网格 */
.info-grid { padding: 4px 0; }

.info-row {
  display: flex;
  align-items: flex-start;
  padding: 7px 16px;
  border-bottom: 1px solid #fafafa;
  gap: 12px;
  min-height: 34px;
}
.info-row:last-child { border-bottom: none; }

.info-label {
  width: 80px;
  flex-shrink: 0;
  font-size: 12px;
  color: #999;
  padding-top: 2px;
}

.info-value {
  flex: 1;
  font-size: 13px;
  color: #1a1a1a;
  word-break: break-all;
  line-height: 1.5;
}

.info-value.mono { font-family: 'Consolas', monospace; }
.info-value.desc { color: #555; line-height: 1.6; }

.ref-link {
  font-size: 12px;
  color: #1677ff;
  word-break: break-all;
  line-height: 1.6;
}

/* 严重性徽章 */
.severity-badge {
  display: inline-block;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
}
.severity-badge.critical { background: #fff0f0; color: #cf1322; }
.severity-badge.high     { background: #fff2e8; color: #d4380d; }
.severity-badge.medium   { background: #fffbe6; color: #d48806; }
.severity-badge.low      { background: #f6ffed; color: #389e0d; }
.severity-badge.info     { background: #e6f4ff; color: #0958d9; }
.severity-badge.unknown  { background: #f5f5f5; color: #8c8c8c; }

.proto-badge {
  display: inline-block;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  background: #f0f5ff;
  color: #2f54eb;
}

/* 请求块 */
.request-block {
  border-bottom: 1px solid #f5f5f5;
  padding: 12px 16px;
}
.request-block:last-child { border-bottom: none; }

.req-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 10px;
}

.req-index {
  font-size: 11px;
  color: #999;
  background: #f5f5f5;
  padding: 1px 6px;
  border-radius: 4px;
}

.req-method {
  font-size: 12px;
  font-weight: 700;
  color: #1677ff;
  background: #e6f4ff;
  padding: 2px 8px;
  border-radius: 4px;
  font-family: monospace;
}

.req-attack {
  font-size: 11px;
  color: #d46b08;
  background: #fff7e6;
  padding: 1px 6px;
  border-radius: 4px;
}

.req-group { margin-bottom: 10px; }
.req-group:last-child { margin-bottom: 0; }

.req-group-label {
  font-size: 11px;
  font-weight: 600;
  color: #999;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin-bottom: 5px;
  display: flex;
  align-items: center;
  gap: 6px;
}

.condition-badge {
  font-size: 10px;
  background: #f9f0ff;
  color: #722ed1;
  padding: 1px 5px;
  border-radius: 3px;
  text-transform: none;
  font-weight: 500;
}

.code-line {
  font-family: 'Consolas', monospace;
  font-size: 12px;
  color: #1a1a1a;
  background: #f7f8fa;
  padding: 4px 8px;
  border-radius: 4px;
  margin-bottom: 3px;
  word-break: break-all;
}

.raw-block {
  font-family: 'Consolas', monospace;
  font-size: 11.5px;
  color: #1a1a1a;
  background: #f7f8fa;
  padding: 8px 10px;
  border-radius: 6px;
  margin: 0 0 4px;
  white-space: pre-wrap;
  word-break: break-all;
  line-height: 1.5;
}

/* Matcher */
.matcher-item {
  display: flex;
  flex-wrap: wrap;
  align-items: flex-start;
  gap: 6px;
  padding: 6px 0;
  border-bottom: 1px solid #fafafa;
}
.matcher-item:last-child { border-bottom: none; }

.matcher-type {
  font-size: 11px;
  font-weight: 600;
  background: #e6f4ff;
  color: #1677ff;
  padding: 2px 7px;
  border-radius: 4px;
  flex-shrink: 0;
}

.matcher-part {
  font-size: 11px;
  background: #f5f5f5;
  color: #555;
  padding: 2px 7px;
  border-radius: 4px;
  font-family: monospace;
  flex-shrink: 0;
}

.matcher-neg {
  font-size: 10px;
  background: #fff1f0;
  color: #ff4d4f;
  padding: 2px 6px;
  border-radius: 4px;
  font-weight: 600;
  flex-shrink: 0;
}

.matcher-values {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  flex: 1;
}

.match-val {
  font-size: 11.5px;
  font-family: 'Consolas', monospace;
  background: #f7f8fa;
  color: #1a1a1a;
  padding: 2px 7px;
  border-radius: 4px;
  border: 1px solid #f0f0f0;
  word-break: break-all;
}

.match-val.regex  { background: #fff7e6; color: #d46b08; border-color: #ffd591; }
.match-val.dsl    { background: #f9f0ff; color: #531dab; border-color: #d3adf7; }
.match-val.status { background: #f6ffed; color: #389e0d; border-color: #b7eb8f; }

/* Code source */
.source-type-badge {
  font-size: 10px;
  font-weight: 600;
  padding: 1px 6px;
  border-radius: 3px;
  text-transform: none;
}
.source-type-badge.file   { background: #e6f4ff; color: #0958d9; }
.source-type-badge.inline { background: #f6ffed; color: #389e0d; }

.code-source {
  max-height: 320px;
  overflow-y: auto;
  font-size: 12px;
  line-height: 1.6;
  tab-size: 2;
}

.mr-1 { margin-right: 4px; }
</style>
