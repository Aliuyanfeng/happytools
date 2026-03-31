<template>
  <div class="ts-page">

    <!-- ── 顶部：时间戳 ↔ 本地时间 ── -->
    <div class="section">
      <div class="section-title">时间戳 ↔ 本地时间</div>

      <div class="convert-grid">
        <!-- 左：时间戳 → 时间 -->
        <div class="conv-card">
          <div class="conv-label">时间戳 → 可读时间</div>
          <div class="conv-row">
            <a-input
              v-model:value="tsInput"
              placeholder="输入时间戳，如 1700000000"
              allow-clear
              size="small"
              @input="tsToTime"
              style="flex:1"
            />
            <a-radio-group v-model:value="tsUnit" button-style="solid" size="small" @change="tsToTime">
              <a-radio-button value="s" style="width:44px;text-align:center;font-size:12px">秒</a-radio-button>
              <a-radio-button value="ms" style="width:44px;text-align:center;font-size:12px">毫秒</a-radio-button>
            </a-radio-group>
          </div>
          <div class="conv-row">
            <a-select v-model:value="tsTimezone" style="width:100%" size="small" show-search @change="tsToTime">
              <a-select-option v-for="tz in timezones" :key="tz.value" :value="tz.value">
                {{ tz.label }}
              </a-select-option>
            </a-select>
          </div>
          <div class="result-box" v-if="tsResult">
            <div class="result-row"><span class="rk">本地时间</span><span class="rv">{{ tsResult.local }}</span></div>
            <div class="result-row"><span class="rk">UTC 时间</span><span class="rv">{{ tsResult.utc }}</span></div>
            <div class="result-row"><span class="rk">ISO 8601</span><span class="rv">{{ tsResult.iso }}</span></div>
            <div class="result-row"><span class="rk">星期</span><span class="rv">{{ tsResult.weekday }}</span></div>
          </div>
          <div class="error-tip" v-if="tsError">{{ tsError }}</div>
        </div>

        <!-- 右：时间 → 时间戳 -->
        <div class="conv-card">
          <div class="conv-label">可读时间 → 时间戳</div>
          <div class="conv-row">
            <a-input
              v-model:value="timeInput"
              placeholder="如 2024-01-01 12:00:00"
              allow-clear
              @input="timeToTs"
            />
          </div>
          <div class="conv-row">
            <a-select v-model:value="timeTimezone" style="width:100%" size="small" show-search @change="timeToTs">
              <a-select-option v-for="tz in timezones" :key="tz.value" :value="tz.value">
                {{ tz.label }}
              </a-select-option>
            </a-select>
          </div>
          <div class="result-box" v-if="timeResult">
            <div class="result-row">
              <span class="rk">秒级时间戳</span>
              <span class="rv mono">{{ timeResult.sec }}</span>
              <a-button size="small" type="text" @click="copy(String(timeResult.sec))">复制</a-button>
            </div>
            <div class="result-row">
              <span class="rk">毫秒时间戳</span>
              <span class="rv mono">{{ timeResult.ms }}</span>
              <a-button size="small" type="text" @click="copy(String(timeResult.ms))">复制</a-button>
            </div>
          </div>
          <div class="error-tip" v-if="timeError">{{ timeError }}</div>

          <!-- 当前时间戳 -->
          <div class="now-box">
            <span class="now-label">当前时间戳</span>
            <span class="now-sec mono">{{ nowSec }}</span>
            <span class="now-unit">秒</span>
            <span class="now-ms mono">{{ nowMs }}</span>
            <span class="now-unit">毫秒</span>
            <a-button size="small" @click="fillNow">填入</a-button>
          </div>
        </div>
      </div>
    </div>

    <!-- ── 时间差值计算 ── -->
    <div class="section">
      <div class="section-title">时间差值计算</div>
      <div class="diff-grid">
        <div class="diff-input-group">
          <label class="conv-label">开始时间</label>
          <a-input
            v-model:value="diffStart"
            placeholder="支持多种格式，如 2026-03-30 17:04:47"
            allow-clear
            @input="calcDiff"
          />
          <div class="parsed-hint" v-if="diffStartParsed">→ {{ diffStartParsed }}</div>
          <div class="error-tip" v-if="diffStartErr">{{ diffStartErr }}</div>
        </div>
        <div class="diff-swap" @click="swapDiff" title="交换">⇅</div>
        <div class="diff-input-group">
          <label class="conv-label">结束时间</label>
          <a-input
            v-model:value="diffEnd"
            placeholder="支持多种格式，如 2026/03/31 09:00:00"
            allow-clear
            @input="calcDiff"
          />
          <div class="parsed-hint" v-if="diffEndParsed">→ {{ diffEndParsed }}</div>
          <div class="error-tip" v-if="diffEndErr">{{ diffEndErr }}</div>
        </div>
      </div>

      <!-- 结果 -->
      <div class="diff-result" v-if="diffResult">
        <div class="diff-result-main">
          <span class="diff-sign" :class="diffResult.sign === '-' ? 'neg' : 'pos'">
            {{ diffResult.sign === '-' ? '结束早于开始' : '结束晚于开始' }}
          </span>
        </div>
        <div class="diff-cards">
          <div class="diff-card">
            <div class="dc-val">{{ diffResult.days }}</div>
            <div class="dc-unit">天</div>
          </div>
          <div class="diff-card">
            <div class="dc-val">{{ diffResult.hours }}</div>
            <div class="dc-unit">小时</div>
          </div>
          <div class="diff-card">
            <div class="dc-val">{{ diffResult.minutes }}</div>
            <div class="dc-unit">分钟</div>
          </div>
          <div class="diff-card">
            <div class="dc-val">{{ diffResult.seconds }}</div>
            <div class="dc-unit">秒</div>
          </div>
        </div>
        <div class="diff-total">
          <span class="dt-item"><b>{{ diffResult.totalSeconds }}</b> 秒</span>
          <span class="dt-sep">·</span>
          <span class="dt-item"><b>{{ diffResult.totalMinutes }}</b> 分钟</span>
          <span class="dt-sep">·</span>
          <span class="dt-item"><b>{{ diffResult.totalHours }}</b> 小时</span>
          <span class="dt-sep">·</span>
          <span class="dt-item"><b>{{ diffResult.totalDays }}</b> 天</span>
        </div>
      </div>
    </div>

    <!-- ── 世界时钟 ── -->
    <div class="section">
      <div class="section-header">
        <div class="section-title">世界时钟</div>
        <a-select
          v-model:value="addTz"
          style="width:220px"
          size="small"
          show-search
          placeholder="添加时区..."
          @change="addClock"
        >
          <a-select-option v-for="tz in availableToAdd" :key="tz.value" :value="tz.value">
            {{ tz.label }}
          </a-select-option>
        </a-select>
      </div>
      <div class="clock-grid">
        <div class="clock-card" v-for="tz in clockZones" :key="tz.value">
          <div class="clock-remove" @click="removeClock(tz.value)">×</div>
          <div class="clock-city">{{ tz.city }}</div>
          <div class="clock-time">{{ clockTimes[tz.value] }}</div>
          <div class="clock-date">{{ clockDates[tz.value] }}</div>
          <div class="clock-offset">{{ tz.label }}</div>
        </div>
      </div>
    </div>

  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { message } from 'ant-design-vue'

// ── 时区列表 ──────────────────────────────────────────────────────────────
const timezones = [
  { value: 'Asia/Shanghai',     label: 'UTC+8  北京/上海' },
  { value: 'Asia/Tokyo',        label: 'UTC+9  东京' },
  { value: 'Asia/Seoul',        label: 'UTC+9  首尔' },
  { value: 'Asia/Singapore',    label: 'UTC+8  新加坡' },
  { value: 'Asia/Hong_Kong',    label: 'UTC+8  香港' },
  { value: 'Asia/Kolkata',      label: 'UTC+5:30 孟买' },
  { value: 'Asia/Dubai',        label: 'UTC+4  迪拜' },
  { value: 'Europe/London',     label: 'UTC+0  伦敦' },
  { value: 'Europe/Paris',      label: 'UTC+1  巴黎' },
  { value: 'Europe/Berlin',     label: 'UTC+1  柏林' },
  { value: 'Europe/Moscow',     label: 'UTC+3  莫斯科' },
  { value: 'America/New_York',  label: 'UTC-5  纽约' },
  { value: 'America/Chicago',   label: 'UTC-6  芝加哥' },
  { value: 'America/Denver',    label: 'UTC-7  丹佛' },
  { value: 'America/Los_Angeles', label: 'UTC-8  洛杉矶' },
  { value: 'America/Sao_Paulo', label: 'UTC-3  圣保罗' },
  { value: 'Pacific/Auckland',  label: 'UTC+12 奥克兰' },
  { value: 'Pacific/Sydney',    label: 'UTC+10 悉尼' },
  { value: 'UTC',               label: 'UTC+0  协调世界时' },
]

const cityMap: Record<string, string> = {
  'Asia/Shanghai': '北京', 'Asia/Tokyo': '东京', 'Asia/Seoul': '首尔',
  'Asia/Singapore': '新加坡', 'Asia/Hong_Kong': '香港', 'Asia/Kolkata': '孟买',
  'Asia/Dubai': '迪拜', 'Europe/London': '伦敦', 'Europe/Paris': '巴黎',
  'Europe/Berlin': '柏林', 'Europe/Moscow': '莫斯科', 'America/New_York': '纽约',
  'America/Chicago': '芝加哥', 'America/Denver': '丹佛', 'America/Los_Angeles': '洛杉矶',
  'America/Sao_Paulo': '圣保罗', 'Pacific/Auckland': '奥克兰', 'Pacific/Sydney': '悉尼',
  'UTC': 'UTC',
}

// ── 时间戳 → 时间 ─────────────────────────────────────────────────────────
const tsInput = ref('')
const tsUnit = ref<'s' | 'ms'>('s')
const tsTimezone = ref('Asia/Shanghai')
const tsResult = ref<{ local: string; utc: string; iso: string; weekday: string } | null>(null)
const tsError = ref('')

const weekdays = ['周日', '周一', '周二', '周三', '周四', '周五', '周六']

function tsToTime() {
  tsError.value = ''
  tsResult.value = null
  const raw = tsInput.value.trim()
  if (!raw) return
  if (!/^\d+$/.test(raw)) { tsError.value = '请输入纯数字时间戳'; return }
  let ms = parseInt(raw)
  if (tsUnit.value === 's') ms *= 1000
  if (ms > 9999999999999) { tsError.value = '时间戳超出范围'; return }
  const d = new Date(ms)
  if (isNaN(d.getTime())) { tsError.value = '无效时间戳'; return }

  const fmt = (date: Date, tz: string) =>
    date.toLocaleString('zh-CN', { timeZone: tz, hour12: false,
      year: 'numeric', month: '2-digit', day: '2-digit',
      hour: '2-digit', minute: '2-digit', second: '2-digit' })

  tsResult.value = {
    local: fmt(d, tsTimezone.value),
    utc:   fmt(d, 'UTC'),
    iso:   d.toISOString(),
    weekday: weekdays[d.getDay()],
  }
}

// ── 时间 → 时间戳 ─────────────────────────────────────────────────────────
const timeInput = ref('')
const timeTimezone = ref('Asia/Shanghai')
const timeResult = ref<{ sec: number; ms: number } | null>(null)
const timeError = ref('')

function timeToTs() {
  timeError.value = ''
  timeResult.value = null
  const raw = timeInput.value.trim()
  if (!raw) return
  // 尝试解析，支持多种格式
  const d = new Date(raw)
  if (isNaN(d.getTime())) { timeError.value = '无法解析时间，请使用 YYYY-MM-DD HH:mm:ss 格式'; return }
  const ms = d.getTime()
  timeResult.value = { sec: Math.floor(ms / 1000), ms }
}

// ── 当前时间戳 ────────────────────────────────────────────────────────────
const nowSec = ref(Math.floor(Date.now() / 1000))
const nowMs  = ref(Date.now())

function fillNow() {
  tsInput.value = String(nowSec.value)
  tsUnit.value = 's'
  tsToTime()
}

// ── 世界时钟 ──────────────────────────────────────────────────────────────
const defaultClocks = ['Asia/Shanghai', 'America/New_York', 'Europe/London', 'Asia/Tokyo']
const clockZones = ref(
  defaultClocks.map(v => ({ value: v, city: cityMap[v] ?? v, label: timezones.find(t => t.value === v)?.label ?? v }))
)
const clockTimes = ref<Record<string, string>>({})
const clockDates = ref<Record<string, string>>({})
const addTz = ref<string | undefined>(undefined)

const availableToAdd = computed(() =>
  timezones.filter(t => !clockZones.value.some(c => c.value === t.value))
)

function addClock(val: string) {
  const tz = timezones.find(t => t.value === val)
  if (!tz) return
  clockZones.value.push({ value: tz.value, city: cityMap[tz.value] ?? tz.value, label: tz.label })
  addTz.value = undefined
}

function removeClock(val: string) {
  clockZones.value = clockZones.value.filter(c => c.value !== val)
}

function updateClocks() {
  const now = new Date()
  nowSec.value = Math.floor(now.getTime() / 1000)
  nowMs.value  = now.getTime()
  clockZones.value.forEach(tz => {
    clockTimes.value[tz.value] = now.toLocaleTimeString('zh-CN', { timeZone: tz.value, hour12: false })
    clockDates.value[tz.value] = now.toLocaleDateString('zh-CN', { timeZone: tz.value, month: '2-digit', day: '2-digit', weekday: 'short' })
  })
}

let timer: ReturnType<typeof setInterval>
onMounted(() => { updateClocks(); timer = setInterval(updateClocks, 1000) })
onUnmounted(() => clearInterval(timer))

// ── 时间差值计算 ──────────────────────────────────────────────────────────
const diffStart = ref('')
const diffEnd   = ref('')
const diffStartParsed = ref('')
const diffEndParsed   = ref('')
const diffStartErr = ref('')
const diffEndErr   = ref('')

interface DiffResult {
  sign: '+' | '-'
  days: number; hours: number; minutes: number; seconds: number
  totalSeconds: number; totalMinutes: string; totalHours: string; totalDays: string
}
const diffResult = ref<DiffResult | null>(null)

// 兼容多种时间格式的解析器
function parseFlexTime(raw: string): { date: Date | null; formatted: string; err: string } {
  const s = raw.trim()
  if (!s) return { date: null, formatted: '', err: '' }

  // 预处理：统一分隔符
  let normalized = s
    .replace(/\//g, '-')          // 2026/03/30 → 2026-03-30
    .replace(/\./g, '-')          // 2026.03.30 → 2026-03-30
    .replace(/年/g, '-').replace(/月/g, '-').replace(/日/g, ' ')  // 中文日期
    .replace(/时/g, ':').replace(/分/g, ':').replace(/秒/g, '')
    .replace(/T/g, ' ')           // ISO 8601
    .replace(/\s+/g, ' ')
    .trim()

  // 补全时间部分
  if (/^\d{4}-\d{1,2}-\d{1,2}$/.test(normalized)) normalized += ' 00:00:00'
  if (/^\d{4}-\d{1,2}-\d{1,2} \d{1,2}:\d{1,2}$/.test(normalized)) normalized += ':00'

  // 尝试解析
  const d = new Date(normalized)
  if (!isNaN(d.getTime())) {
    const fmt = d.toLocaleString('zh-CN', { hour12: false,
      year: 'numeric', month: '2-digit', day: '2-digit',
      hour: '2-digit', minute: '2-digit', second: '2-digit' })
    return { date: d, formatted: fmt, err: '' }
  }

  // 尝试纯时间戳
  if (/^\d{10}$/.test(s)) {
    const d2 = new Date(parseInt(s) * 1000)
    if (!isNaN(d2.getTime())) {
      return { date: d2, formatted: d2.toLocaleString('zh-CN', { hour12: false }), err: '' }
    }
  }
  if (/^\d{13}$/.test(s)) {
    const d2 = new Date(parseInt(s))
    if (!isNaN(d2.getTime())) {
      return { date: d2, formatted: d2.toLocaleString('zh-CN', { hour12: false }), err: '' }
    }
  }

  return { date: null, formatted: '', err: '无法识别的时间格式' }
}

function calcDiff() {
  diffResult.value = null
  diffStartErr.value = ''
  diffEndErr.value = ''
  diffStartParsed.value = ''
  diffEndParsed.value = ''

  const r1 = parseFlexTime(diffStart.value)
  const r2 = parseFlexTime(diffEnd.value)

  if (r1.err) { diffStartErr.value = r1.err; return }
  if (r2.err) { diffEndErr.value = r2.err; return }
  if (r1.formatted) diffStartParsed.value = r1.formatted
  if (r2.formatted) diffEndParsed.value = r2.formatted
  if (!r1.date || !r2.date) return

  const diffMs = r2.date.getTime() - r1.date.getTime()
  const sign = diffMs >= 0 ? '+' : '-'
  const abs = Math.abs(diffMs)

  const totalSec = Math.floor(abs / 1000)
  const days    = Math.floor(totalSec / 86400)
  const hours   = Math.floor((totalSec % 86400) / 3600)
  const minutes = Math.floor((totalSec % 3600) / 60)
  const seconds = totalSec % 60

  diffResult.value = {
    sign,
    days, hours, minutes, seconds,
    totalSeconds: totalSec,
    totalMinutes: (abs / 60000).toFixed(2),
    totalHours:   (abs / 3600000).toFixed(4),
    totalDays:    (abs / 86400000).toFixed(4),
  }
}

function swapDiff() {
  const tmp = diffStart.value
  diffStart.value = diffEnd.value
  diffEnd.value = tmp
  calcDiff()
}

async function copy(text: string) {
  try { await navigator.clipboard.writeText(text); message.success('已复制') }
  catch { message.error('复制失败') }
}
</script>

<style scoped>
.ts-page {
  height: 100%;
  overflow-y: auto;
  padding: 20px 24px;
  display: flex;
  flex-direction: column;
  gap: 20px;
  background: #f6f8fb;
}

/* ── section ── */
.section {
  background: #fff;
  border-radius: 16px;
  border: 1px solid #e2e8f0;
  padding: 18px 20px;
  display: flex;
  flex-direction: column;
  gap: 14px;
}
.section-header { display: flex; align-items: center; justify-content: space-between; }
.section-title { font-size: 14px; font-weight: 700; color: #1e1b4b; }

/* ── 转换网格 ── */
.convert-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; }
@media (max-width: 700px) { .convert-grid { grid-template-columns: 1fr; } }

.conv-card {
  background: #f8fafc;
  border-radius: 12px;
  border: 1px solid #e2e8f0;
  padding: 14px 16px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}
.conv-label { font-size: 12px; font-weight: 600; color: #6366f1; }
.conv-row { display: flex; gap: 8px; align-items: center; }

/* 结果 */
.result-box {
  background: #fff;
  border-radius: 10px;
  border: 1px solid #e2e8f0;
  padding: 10px 12px;
  display: flex;
  flex-direction: column;
  gap: 6px;
}
.result-row { display: flex; align-items: center; gap: 8px; flex-wrap: wrap; }
.rk { font-size: 11px; color: #94a3b8; width: 80px; flex-shrink: 0; }
.rv { font-size: 12px; color: #1e293b; font-weight: 500; flex: 1; }
.mono { font-family: monospace; }

.error-tip { font-size: 12px; color: #ef4444; }

/* 当前时间戳 */
.now-box {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  background: #eef2ff;
  border-radius: 10px;
  flex-wrap: wrap;
}
.now-label { font-size: 11px; color: #6366f1; font-weight: 600; }
.now-sec, .now-ms { font-family: monospace; font-size: 13px; font-weight: 700; color: #1e1b4b; }
.now-unit { font-size: 11px; color: #94a3b8; }

/* ── 时间差值 ── */
.diff-grid {
  display: flex;
  align-items: flex-start;
  gap: 10px;
}
.diff-input-group {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 6px;
}
.parsed-hint {
  font-size: 11px;
  color: #059669;
  font-family: monospace;
}
.diff-swap {
  font-size: 22px;
  color: #6366f1;
  cursor: pointer;
  padding: 6px 4px;
  margin-top: 22px;
  transition: transform 0.2s;
  user-select: none;
}
.diff-swap:hover { transform: scale(1.3); }

.diff-result { display: flex; flex-direction: column; gap: 10px; }
.diff-result-main { display: flex; align-items: center; gap: 8px; }
.diff-sign {
  font-size: 12px;
  font-weight: 600;
  padding: 2px 10px;
  border-radius: 8px;
}
.diff-sign.pos { background: #dcfce7; color: #16a34a; }
.diff-sign.neg { background: #fee2e2; color: #dc2626; }

.diff-cards {
  display: flex;
  gap: 10px;
}
.diff-card {
  flex: 1;
  background: linear-gradient(145deg, #eef2ff, #e0e7ff);
  border-radius: 12px;
  padding: 12px 8px;
  text-align: center;
  border: 1px solid rgba(99,102,241,0.15);
}
.dc-val { font-size: 24px; font-weight: 700; color: #3730a3; font-variant-numeric: tabular-nums; }
.dc-unit { font-size: 11px; color: #6366f1; margin-top: 2px; }

.diff-total {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
  padding: 10px 14px;
  background: #f8fafc;
  border-radius: 10px;
  border: 1px solid #e2e8f0;
}
.dt-item { font-size: 12px; color: #475569; }
.dt-item b { color: #1e1b4b; font-family: monospace; }
.dt-sep { color: #cbd5e1; }

/* ── 世界时钟 ── */
.clock-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  gap: 12px;
}
.clock-card {
  position: relative;
  background: linear-gradient(145deg, #eef2ff, #e0e7ff);
  border-radius: 14px;
  padding: 14px 14px 12px;
  border: 1px solid rgba(99,102,241,0.15);
  transition: box-shadow 0.2s;
}
.clock-card:hover { box-shadow: 0 6px 20px rgba(99,102,241,0.18); }
.clock-remove {
  position: absolute;
  top: 8px; right: 10px;
  font-size: 14px;
  color: #94a3b8;
  cursor: pointer;
  line-height: 1;
  transition: color 0.15s;
}
.clock-remove:hover { color: #ef4444; }
.clock-city { font-size: 12px; font-weight: 700; color: #3730a3; margin-bottom: 4px; }
.clock-time { font-size: 22px; font-weight: 700; color: #1e1b4b; font-variant-numeric: tabular-nums; letter-spacing: 1px; }
.clock-date { font-size: 11px; color: #6366f1; margin-top: 2px; }
.clock-offset { font-size: 10px; color: #94a3b8; margin-top: 4px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
</style>
