<template>
  <div class="cc-page">
    <div class="cc-header">
      <h2 class="cc-title">颜色转换</h2>
      <p class="cc-sub">HEX · RGB · HSL 互转，实时预览</p>
    </div>

    <!-- 颜色预览 + 选色器 -->
    <div class="preview-row">
      <div class="color-swatch" :style="{ background: hexValue || '#ffffff' }" />
      <div class="picker-wrap">
        <label class="input-label">颜色选择器</label>
        <input type="color" class="native-picker" :value="hexValue" @input="onPickerInput" />
      </div>
      <div class="preview-info" v-if="parsed">
        <span class="pi-item" :style="{ background: parsed.hex }">{{ parsed.hex }}</span>
        <span class="pi-item">{{ parsed.rgb }}</span>
        <span class="pi-item">{{ parsed.hsl }}</span>
      </div>
    </div>

    <!-- 三栏输入 -->
    <div class="input-grid">
      <!-- HEX -->
      <div class="input-card">
        <div class="ic-label">HEX</div>
        <a-input
          v-model:value="hexInput"
          placeholder="#43ad7f 或 43ad7f"
          allow-clear
          @input="fromHex"
        />
        <div class="error-tip" v-if="hexErr">{{ hexErr }}</div>
        <div class="ic-hint">支持 3位 / 6位 / 8位（含透明度）</div>
      </div>

      <!-- RGB -->
      <div class="input-card">
        <div class="ic-label">RGB</div>
        <a-input
          v-model:value="rgbInput"
          placeholder="rgb(67, 173, 127) 或 67,173,127"
          allow-clear
          @input="fromRgb"
        />
        <div class="error-tip" v-if="rgbErr">{{ rgbErr }}</div>
        <div class="ic-hint">每个通道 0–255</div>
      </div>

      <!-- HSL -->
      <div class="input-card">
        <div class="ic-label">HSL</div>
        <a-input
          v-model:value="hslInput"
          placeholder="hsl(152, 44%, 47%) 或 152,44,47"
          allow-clear
          @input="fromHsl"
        />
        <div class="error-tip" v-if="hslErr">{{ hslErr }}</div>
        <div class="ic-hint">H: 0–360，S/L: 0–100%</div>
      </div>
    </div>

    <!-- 结果卡片 -->
    <div class="result-section" v-if="parsed">
      <div class="result-card" v-for="item in resultItems" :key="item.label">
        <div class="rc-label">{{ item.label }}</div>
        <div class="rc-value">{{ item.value }}</div>
        <a-button size="small" type="text" @click="copy(item.value)">复制</a-button>
      </div>
    </div>

    <!-- 色板：常用颜色 -->
    <div class="palette-section">
      <div class="section-title">常用颜色</div>
      <div class="palette-grid">
        <div
          v-for="c in palette"
          :key="c.hex"
          class="palette-item"
          :style="{ background: c.hex }"
          :title="c.name"
          @click="applyPalette(c.hex)"
        >
          <span class="palette-name">{{ c.name }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { message } from 'ant-design-vue'

// ── 状态 ──────────────────────────────────────────────────────────────────
const hexInput = ref('')
const rgbInput = ref('')
const hslInput = ref('')
const hexErr   = ref('')
const rgbErr   = ref('')
const hslErr   = ref('')

interface Parsed { hex: string; rgb: string; hsl: string; r: number; g: number; b: number }
const parsed = ref<Parsed | null>(null)

const hexValue = computed(() => parsed.value?.hex ?? '#cccccc')

const resultItems = computed(() => parsed.value ? [
  { label: 'HEX',          value: parsed.value.hex },
  { label: 'HEX（大写）',  value: parsed.value.hex.toUpperCase() },
  { label: 'RGB',          value: parsed.value.rgb },
  { label: 'RGB（0–1）',   value: `rgb(${(parsed.value.r/255).toFixed(3)}, ${(parsed.value.g/255).toFixed(3)}, ${(parsed.value.b/255).toFixed(3)})` },
  { label: 'HSL',          value: parsed.value.hsl },
  { label: 'CSS var',      value: `--color: ${parsed.value.hex};` },
] : [])

// ── 转换核心 ──────────────────────────────────────────────────────────────
function hexToRgb(hex: string): { r: number; g: number; b: number } | null {
  let h = hex.replace('#', '')
  if (h.length === 3) h = h.split('').map(c => c + c).join('')
  if (h.length === 8) h = h.slice(0, 6) // 去掉 alpha
  if (h.length !== 6 || !/^[0-9a-fA-F]{6}$/.test(h)) return null
  return {
    r: parseInt(h.slice(0, 2), 16),
    g: parseInt(h.slice(2, 4), 16),
    b: parseInt(h.slice(4, 6), 16),
  }
}

function rgbToHsl(r: number, g: number, b: number): { h: number; s: number; l: number } {
  const rn = r / 255, gn = g / 255, bn = b / 255
  const max = Math.max(rn, gn, bn), min = Math.min(rn, gn, bn)
  let h = 0, s = 0
  const l = (max + min) / 2
  if (max !== min) {
    const d = max - min
    s = l > 0.5 ? d / (2 - max - min) : d / (max + min)
    switch (max) {
      case rn: h = ((gn - bn) / d + (gn < bn ? 6 : 0)) / 6; break
      case gn: h = ((bn - rn) / d + 2) / 6; break
      case bn: h = ((rn - gn) / d + 4) / 6; break
    }
  }
  return { h: Math.round(h * 360), s: Math.round(s * 100), l: Math.round(l * 100) }
}

function hslToRgb(h: number, s: number, l: number): { r: number; g: number; b: number } {
  const sn = s / 100, ln = l / 100
  const c = (1 - Math.abs(2 * ln - 1)) * sn
  const x = c * (1 - Math.abs((h / 60) % 2 - 1))
  const m = ln - c / 2
  let r = 0, g = 0, b = 0
  if (h < 60)       { r = c; g = x }
  else if (h < 120) { r = x; g = c }
  else if (h < 180) { g = c; b = x }
  else if (h < 240) { g = x; b = c }
  else if (h < 300) { r = x; b = c }
  else              { r = c; b = x }
  return { r: Math.round((r + m) * 255), g: Math.round((g + m) * 255), b: Math.round((b + m) * 255) }
}

function applyRgb(r: number, g: number, b: number) {
  const hex = '#' + [r, g, b].map(v => v.toString(16).padStart(2, '0')).join('')
  const { h, s, l } = rgbToHsl(r, g, b)
  parsed.value = {
    hex, r, g, b,
    rgb: `rgb(${r}, ${g}, ${b})`,
    hsl: `hsl(${h}, ${s}%, ${l}%)`,
  }
  hexInput.value = hex
  rgbInput.value = `rgb(${r}, ${g}, ${b})`
  hslInput.value = `hsl(${h}, ${s}%, ${l}%)`
}

// ── 从 HEX 输入 ───────────────────────────────────────────────────────────
function fromHex() {
  hexErr.value = ''
  const raw = hexInput.value.trim()
  if (!raw) { parsed.value = null; return }
  const rgb = hexToRgb(raw)
  if (!rgb) { hexErr.value = '无效的 HEX 颜色，请输入 3 或 6 位十六进制'; return }
  applyRgb(rgb.r, rgb.g, rgb.b)
}

// ── 从 RGB 输入 ───────────────────────────────────────────────────────────
function fromRgb() {
  rgbErr.value = ''
  const raw = rgbInput.value.trim()
  if (!raw) { parsed.value = null; return }
  // 支持 rgb(r,g,b) 或 r,g,b 或 r g b
  const nums = raw.replace(/rgb\(|\)/gi, '').split(/[,\s]+/).map(Number)
  if (nums.length !== 3 || nums.some(n => isNaN(n) || n < 0 || n > 255)) {
    rgbErr.value = '格式错误，示例：rgb(67, 173, 127) 或 67,173,127'; return
  }
  applyRgb(nums[0], nums[1], nums[2])
}

// ── 从 HSL 输入 ───────────────────────────────────────────────────────────
function fromHsl() {
  hslErr.value = ''
  const raw = hslInput.value.trim()
  if (!raw) { parsed.value = null; return }
  const nums = raw.replace(/hsl\(|\)|%/gi, '').split(/[,\s]+/).map(Number)
  if (nums.length !== 3 || isNaN(nums[0]) || isNaN(nums[1]) || isNaN(nums[2])) {
    hslErr.value = '格式错误，示例：hsl(152, 44%, 47%) 或 152,44,47'; return
  }
  const [h, s, l] = nums
  if (h < 0 || h > 360 || s < 0 || s > 100 || l < 0 || l > 100) {
    hslErr.value = 'H: 0–360，S/L: 0–100'; return
  }
  const { r, g, b } = hslToRgb(h, s, l)
  applyRgb(r, g, b)
}

// ── 原生选色器 ────────────────────────────────────────────────────────────
function onPickerInput(e: Event) {
  hexInput.value = (e.target as HTMLInputElement).value
  fromHex()
}

// ── 色板 ──────────────────────────────────────────────────────────────────
function applyPalette(hex: string) {
  hexInput.value = hex
  fromHex()
}

const palette = [
  { hex: '#ef4444', name: '红' },   { hex: '#f97316', name: '橙' },
  { hex: '#eab308', name: '黄' },   { hex: '#22c55e', name: '绿' },
  { hex: '#06b6d4', name: '青' },   { hex: '#3b82f6', name: '蓝' },
  { hex: '#8b5cf6', name: '紫' },   { hex: '#ec4899', name: '粉' },
  { hex: '#ffffff', name: '白' },   { hex: '#f1f5f9', name: '浅灰' },
  { hex: '#94a3b8', name: '灰' },   { hex: '#334155', name: '深灰' },
  { hex: '#0f172a', name: '黑' },   { hex: '#43ad7f', name: '翠绿' },
  { hex: '#f59e0b', name: '琥珀' }, { hex: '#6366f1', name: '靛蓝' },
]

async function copy(text: string) {
  try { await navigator.clipboard.writeText(text); message.success('已复制') }
  catch { message.error('复制失败') }
}
</script>

<style scoped>
.cc-page {
  height: 100%;
  overflow-y: auto;
  padding: 20px 24px;
  display: flex;
  flex-direction: column;
  gap: 18px;
  background: #f6f8fb;
}

.cc-header { flex-shrink: 0; }
.cc-title { font-size: 20px; font-weight: 700; color: #1e1b4b; margin: 0 0 4px; }
.cc-sub { font-size: 13px; color: #94a3b8; margin: 0; }

/* 预览行 */
.preview-row {
  display: flex;
  align-items: center;
  gap: 16px;
  background: #fff;
  border-radius: 16px;
  border: 1px solid #e2e8f0;
  padding: 16px 20px;
  flex-shrink: 0;
}
.color-swatch {
  width: 72px; height: 72px;
  border-radius: 14px;
  border: 1px solid rgba(0,0,0,0.08);
  flex-shrink: 0;
  transition: background 0.2s;
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
}
.picker-wrap { display: flex; flex-direction: column; gap: 4px; }
.input-label { font-size: 11px; font-weight: 600; color: #94a3b8; }
.native-picker {
  width: 48px; height: 36px;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  cursor: pointer;
  padding: 2px;
  background: none;
}
.preview-info { display: flex; flex-direction: column; gap: 6px; }
.pi-item {
  font-size: 12px;
  font-family: monospace;
  color: #1e293b;
  padding: 3px 10px;
  border-radius: 6px;
  background: #f1f5f9;
  white-space: nowrap;
}

/* 三栏输入 */
.input-grid {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;
  gap: 12px;
  flex-shrink: 0;
}
.input-card {
  background: #fff;
  border-radius: 14px;
  border: 1px solid #e2e8f0;
  padding: 14px 16px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.ic-label { font-size: 12px; font-weight: 700; color: #6366f1; }
.ic-hint  { font-size: 11px; color: #cbd5e1; }
.error-tip { font-size: 12px; color: #ef4444; }

/* 结果卡片 */
.result-section {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
  gap: 10px;
  flex-shrink: 0;
}
.result-card {
  background: #fff;
  border-radius: 12px;
  border: 1px solid #e2e8f0;
  padding: 12px 14px;
  display: flex;
  flex-direction: column;
  gap: 4px;
}
.rc-label { font-size: 10px; color: #94a3b8; text-transform: uppercase; letter-spacing: 0.5px; }
.rc-value { font-size: 13px; font-weight: 600; color: #1e293b; font-family: monospace; word-break: break-all; }

/* 色板 */
.palette-section {
  background: #fff;
  border-radius: 16px;
  border: 1px solid #e2e8f0;
  padding: 16px 20px;
  flex-shrink: 0;
}
.section-title { font-size: 13px; font-weight: 700; color: #1e1b4b; margin-bottom: 12px; }
.palette-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}
.palette-item {
  width: 48px; height: 48px;
  border-radius: 10px;
  cursor: pointer;
  border: 1px solid rgba(0,0,0,0.08);
  display: flex;
  align-items: flex-end;
  justify-content: center;
  padding-bottom: 4px;
  transition: transform 0.15s, box-shadow 0.15s;
  overflow: hidden;
}
.palette-item:hover { transform: scale(1.15); box-shadow: 0 4px 12px rgba(0,0,0,0.15); }
.palette-name {
  font-size: 9px;
  color: rgba(255,255,255,0.9);
  text-shadow: 0 1px 2px rgba(0,0,0,0.5);
  white-space: nowrap;
}
</style>
