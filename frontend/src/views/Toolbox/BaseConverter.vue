<template>
  <div class="base-converter">
    <div class="bc-header">
      <h2 class="bc-title">进制转换</h2>
      <p class="bc-sub">二进制 · 十进制 · 十六进制 相互转换，并展示转换步骤</p>
    </div>

    <!-- 输入区 -->
    <div class="input-section">
      <div class="input-row">
        <div class="input-group">
          <label class="input-label">输入值</label>
          <a-input
            v-model:value="inputValue"
            :placeholder="inputPlaceholder"
            size="large"
            allow-clear
            @input="onInput"
            class="main-input"
          />
        </div>
        <div class="input-group base-select-group">
          <label class="input-label">输入进制</label>
          <a-radio-group v-model:value="fromBase" size="large" button-style="solid" @change="onInput">
            <a-radio-button value="2">二进制</a-radio-button>
            <a-radio-button value="10">十进制</a-radio-button>
            <a-radio-button value="16">十六进制</a-radio-button>
          </a-radio-group>
        </div>
      </div>
      <div v-if="error" class="error-tip">{{ error }}</div>
    </div>

    <!-- 结果区 -->
    <div class="result-section" v-if="result">
      <div class="result-cards">
        <div class="result-card" :class="{ active: fromBase === '2' }">
          <div class="rc-label">二进制 (Base 2)</div>
          <div class="rc-value">{{ result.bin }}</div>
          <a-button size="small" type="text" @click="copy(result.bin)">复制</a-button>
        </div>
        <div class="result-card" :class="{ active: fromBase === '10' }">
          <div class="rc-label">十进制 (Base 10)</div>
          <div class="rc-value">{{ result.dec }}</div>
          <a-button size="small" type="text" @click="copy(result.dec)">复制</a-button>
        </div>
        <div class="result-card" :class="{ active: fromBase === '16' }">
          <div class="rc-label">十六进制 (Base 16)</div>
          <div class="rc-value">{{ result.hex.toUpperCase() }}</div>
          <a-button size="small" type="text" @click="copy(result.hex.toUpperCase())">复制</a-button>
        </div>
      </div>
    </div>

    <!-- 步骤展示 -->
    <div class="steps-section" v-if="steps.length">
      <div class="steps-header">
        <span class="steps-title">📖 转换步骤详解</span>
        <span class="steps-sub">帮助初学者理解转换原理</span>
      </div>
      <div class="steps-list">
        <div v-for="(group, gi) in steps" :key="gi" class="step-group">
          <div class="step-group-title">{{ group.title }}</div>
          <div v-for="(step, si) in group.items" :key="si" class="step-item">
            <div class="step-num">{{ si + 1 }}</div>
            <div class="step-body">
              <div class="step-desc">{{ step.desc }}</div>
              <div class="step-calc" v-if="step.calc">
                <code>{{ step.calc }}</code>
              </div>
              <div class="step-result" v-if="step.result">
                <span class="step-result-label">结果：</span>
                <span class="step-result-val">{{ step.result }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 空状态 -->
    <div class="empty-state" v-if="!result && !error">
      <NumberOutlined class="empty-icon" />
      <p>输入一个数值，选择进制，即可查看转换结果和详细步骤</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { message } from 'ant-design-vue'
import { NumberOutlined } from '@ant-design/icons-vue'

const inputValue = ref('')
const fromBase = ref<'2' | '10' | '16'>('10')
const error = ref('')

interface StepItem { desc: string; calc?: string; result?: string }
interface StepGroup { title: string; items: StepItem[] }

const result = ref<{ bin: string; dec: string; hex: string } | null>(null)
const steps = ref<StepGroup[]>([])

const inputPlaceholder = computed(() => {
  if (fromBase.value === '2')  return '输入二进制数，如：1010'
  if (fromBase.value === '16') return '输入十六进制数，如：FF 或 ff'
  return '输入十进制数，如：255'
})

function onInput() {
  error.value = ''
  result.value = null
  steps.value = []
  const raw = inputValue.value.trim()
  if (!raw) return

  // 验证输入
  const base = parseInt(fromBase.value)
  if (base === 2  && !/^[01]+$/.test(raw))  { error.value = '二进制只能包含 0 和 1'; return }
  if (base === 10 && !/^\d+$/.test(raw))     { error.value = '十进制只能包含数字 0-9'; return }
  if (base === 16 && !/^[0-9a-fA-F]+$/.test(raw)) { error.value = '十六进制只能包含 0-9 和 A-F'; return }

  const decimal = parseInt(raw, base)
  if (isNaN(decimal) || decimal < 0) { error.value = '无效的数值'; return }

  result.value = {
    bin: decimal.toString(2),
    dec: decimal.toString(10),
    hex: decimal.toString(16),
  }

  steps.value = buildSteps(raw, base, decimal)
}

function buildSteps(raw: string, base: number, decimal: number): StepGroup[] {
  const groups: StepGroup[] = []

  // Step 1: 转为十进制
  if (base !== 10) {
    const g: StepGroup = {
      title: base === 2 ? '第一步：二进制 → 十进制（按权展开法）' : '第一步：十六进制 → 十进制（按权展开法）',
      items: []
    }
    g.items.push({
      desc: `将每一位数字乘以对应的权值（${base} 的幂次），然后求和。`,
    })

    const digits = raw.toUpperCase().split('')
    const n = digits.length
    let calcParts: string[] = []
    let sum = 0

    digits.forEach((d, i) => {
      const power = n - 1 - i
      const digitVal = parseInt(d, base)
      const weight = Math.pow(base, power)
      const contribution = digitVal * weight
      sum += contribution
      calcParts.push(`${d}×${base}^${power}(${weight})=${contribution}`)
      g.items.push({
        desc: `第 ${i + 1} 位（从左数）：${d} × ${base}^${power}`,
        calc: `${d} × ${weight} = ${contribution}`,
      })
    })

    g.items.push({
      desc: '将所有结果相加：',
      calc: calcParts.map(p => p.split('=')[1]).join(' + ') + ' = ' + sum,
      result: `十进制结果 = ${sum}`,
    })
    groups.push(g)
  }

  // Step 2: 十进制 → 二进制（除2取余）
  if (base !== 2) {
    const g: StepGroup = {
      title: base === 10 ? '第一步：十进制 → 二进制（除 2 取余法）' : '第二步：十进制 → 二进制（除 2 取余法）',
      items: []
    }
    g.items.push({ desc: '反复将十进制数除以 2，记录余数，直到商为 0，余数从下往上读即为二进制。' })

    let n = decimal
    const remainders: number[] = []
    const divSteps: string[] = []
    while (n > 0) {
      const rem = n % 2
      divSteps.push(`${n} ÷ 2 = ${Math.floor(n / 2)} 余 ${rem}`)
      remainders.push(rem)
      n = Math.floor(n / 2)
    }
    if (remainders.length === 0) remainders.push(0)

    divSteps.forEach((s, i) => {
      g.items.push({ desc: `第 ${i + 1} 次除法：`, calc: s, result: `余数 = ${remainders[i]}` })
    })
    g.items.push({
      desc: '将余数从下往上排列（逆序）：',
      calc: [...remainders].reverse().join(''),
      result: `二进制结果 = ${decimal.toString(2)}`,
    })
    groups.push(g)
  }

  // Step 3: 十进制 → 十六进制（除16取余）
  if (base !== 16) {
    const hexMap = ['0','1','2','3','4','5','6','7','8','9','A','B','C','D','E','F']
    const g: StepGroup = {
      title: base === 10 ? '第二步：十进制 → 十六进制（除 16 取余法）' : '第三步：十进制 → 十六进制（除 16 取余法）',
      items: []
    }
    g.items.push({ desc: '反复将十进制数除以 16，记录余数（10-15 用 A-F 表示），余数从下往上读即为十六进制。' })

    let n = decimal
    const remainders: string[] = []
    while (n > 0) {
      const rem = n % 16
      g.items.push({
        desc: `${n} ÷ 16 = ${Math.floor(n / 16)} 余 ${rem}`,
        calc: `余数 ${rem} → ${hexMap[rem]}`,
      })
      remainders.push(hexMap[rem])
      n = Math.floor(n / 16)
    }
    if (remainders.length === 0) remainders.push('0')

    g.items.push({
      desc: '将余数从下往上排列（逆序）：',
      calc: [...remainders].reverse().join(''),
      result: `十六进制结果 = ${decimal.toString(16).toUpperCase()}`,
    })
    groups.push(g)
  }

  // 二进制 ↔ 十六进制 快捷方法
  if (base === 2 || base === 16) {
    const g: StepGroup = {
      title: '💡 快捷方法：二进制与十六进制直接转换（4位分组法）',
      items: [
        { desc: '每 4 位二进制对应 1 位十六进制，可以直接转换，无需经过十进制。' },
      ]
    }
    const binStr = decimal.toString(2).padStart(Math.ceil(decimal.toString(2).length / 4) * 4, '0')
    const groups4 = binStr.match(/.{1,4}/g) ?? []
    groups4.forEach(chunk => {
      const val = parseInt(chunk, 2)
      const hexChar = val.toString(16).toUpperCase()
      g.items.push({
        desc: `二进制 ${chunk}`,
        calc: `${chunk.split('').map((b, i) => `${b}×2^${3-i}`).join(' + ')} = ${val}`,
        result: `十六进制 ${hexChar}`,
      })
    })
    g.items.push({
      desc: '拼接所有十六进制位：',
      result: `${groups4.map(c => parseInt(c,2).toString(16).toUpperCase()).join('')}`,
    })
    groups.push(g)
  }

  return groups
}

async function copy(text: string) {
  try {
    await navigator.clipboard.writeText(text)
    message.success('已复制')
  } catch {
    message.error('复制失败')
  }
}
</script>

<style scoped>
.base-converter {
  height: 100%;
  overflow-y: auto;
  padding: 24px 28px;
  display: flex;
  flex-direction: column;
  gap: 20px;
  background: #f6f8fb;
}

.bc-header { flex-shrink: 0; }
.bc-title { font-size: 20px; font-weight: 700; color: #1e1b4b; margin: 0 0 4px; }
.bc-sub { font-size: 13px; color: #94a3b8; margin: 0; }

/* 输入区 */
.input-section { flex-shrink: 0; }
.input-row { display: flex; gap: 16px; align-items: flex-end; flex-wrap: wrap; }
.input-group { display: flex; flex-direction: column; gap: 6px; }
.input-group:first-child { flex: 1; min-width: 200px; }
.input-label { font-size: 12px; font-weight: 600; color: #64748b; }
.main-input { font-size: 16px; font-family: monospace; }
.error-tip { margin-top: 8px; color: #ef4444; font-size: 13px; }

/* 结果卡片 */
.result-section { flex-shrink: 0; }
.result-cards { display: flex; gap: 12px; flex-wrap: wrap; }
.result-card {
  flex: 1;
  min-width: 160px;
  background: #fff;
  border: 1px solid #e2e8f0;
  border-radius: 14px;
  padding: 14px 16px;
  display: flex;
  flex-direction: column;
  gap: 6px;
  transition: border-color 0.2s, box-shadow 0.2s;
}
.result-card.active {
  border-color: #6366f1;
  box-shadow: 0 0 0 3px rgba(99,102,241,0.1);
}
.rc-label { font-size: 11px; color: #94a3b8; font-weight: 600; text-transform: uppercase; }
.rc-value { font-size: 18px; font-weight: 700; color: #1e1b4b; font-family: monospace; word-break: break-all; }

/* 步骤区 */
.steps-section {
  flex: 1;
  background: #fff;
  border-radius: 16px;
  border: 1px solid #e2e8f0;
  padding: 18px 20px;
  overflow: auto;
}
.steps-header { display: flex; align-items: baseline; gap: 10px; margin-bottom: 16px; }
.steps-title { font-size: 14px; font-weight: 700; color: #1e1b4b; }
.steps-sub { font-size: 12px; color: #94a3b8; }

.step-group { margin-bottom: 20px; }
.step-group-title {
  font-size: 13px;
  font-weight: 700;
  color: #6366f1;
  margin-bottom: 10px;
  padding-bottom: 6px;
  border-bottom: 1px solid #eef2ff;
}

.step-item {
  display: flex;
  gap: 12px;
  margin-bottom: 10px;
  align-items: flex-start;
}
.step-num {
  width: 22px; height: 22px;
  border-radius: 50%;
  background: #eef2ff;
  color: #6366f1;
  font-size: 11px;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  margin-top: 1px;
}
.step-body { flex: 1; display: flex; flex-direction: column; gap: 4px; }
.step-desc { font-size: 13px; color: #334155; }
.step-calc {
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  padding: 6px 10px;
}
.step-calc code { font-family: monospace; font-size: 13px; color: #0f172a; }
.step-result { font-size: 12px; }
.step-result-label { color: #94a3b8; }
.step-result-val { font-weight: 700; color: #059669; font-family: monospace; }

/* 空状态 */
.empty-state {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  color: #cbd5e1;
}
.empty-icon { font-size: 48px; }
.empty-state p { font-size: 13px; text-align: center; max-width: 300px; }
</style>
