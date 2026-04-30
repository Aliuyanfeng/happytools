<template>
  <!-- 多标签工时比例编辑器：仅在该天有 2+ 个标签时显示 -->
  <div class="ratio-editor">
    <div class="ratio-header">
      <span class="ratio-date">{{ date }}</span>
      <span class="ratio-hint">拖动滑块调整工时占比，总和始终 = 1 天</span>
    </div>

    <!-- 可视化色块条 -->
    <div class="ratio-bar">
      <div
        v-for="(tag, idx) in tags"
        :key="tag"
        class="ratio-bar-segment"
        :style="{ width: (ratios[idx] * 100).toFixed(1) + '%', backgroundColor: getTagColorHex(tag) }"
        :title="tag + ': ' + formatDays(ratios[idx])"
      ></div>
    </div>

    <!-- 每个标签的滑块行 -->
    <div class="ratio-rows">
      <div v-for="(tag, idx) in tags" :key="tag" class="ratio-row">
        <span class="ratio-tag-label">
          <a-tag :color="getTagColor(tag)">{{ tag }}</a-tag>
        </span>
        <div class="ratio-slider-wrap">
          <a-slider
            :value="Math.round(ratios[idx] * 100)"
            :min="0"
            :max="100"
            :step="5"
            :tooltip-formatter="(v: number) => v + '%'"
            @change="(v: number) => handleSliderChange(idx, v)"
            :style="{ flex: 1 }"
          />
        </div>
        <span class="ratio-days-label">{{ formatDays(ratios[idx]) }} 天</span>
      </div>
    </div>

    <div class="ratio-actions">
      <a-button size="small" @click="handleReset">均分重置</a-button>
      <a-button size="small" type="primary" :loading="saving" @click="handleSave">保存</a-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { message } from 'ant-design-vue'
import * as DailyReportService from '../../../bindings/github.com/Aliuyanfeng/happytools/backend/services/dailyreport/dailyreportservice'

const props = defineProps<{
  date: string
  tags: string[]
  // 已保存的比例 map，没有则传 {}
  savedRatios: Record<string, number>
}>()

const emit = defineEmits<{
  (e: 'saved'): void
}>()

// 内部比例数组，与 props.tags 顺序对应，值为 0~1
const ratios = ref<number[]>([])
const saving = ref(false)

// 初始化 / 当 props 变化时重置
const initRatios = () => {
  const n = props.tags.length
  if (n === 0) return
  const arr: number[] = props.tags.map(tag => {
    const v = props.savedRatios[tag]
    return typeof v === 'number' ? v : 1 / n
  })
  // 归一化，防止保存数据有误差
  const sum = arr.reduce((a, b) => a + b, 0)
  ratios.value = sum > 0 ? arr.map(v => v / sum) : arr.map(() => 1 / n)
}

watch(() => [props.date, props.tags, props.savedRatios], initRatios, { immediate: true, deep: true })

// 拖动某个标签的滑块
const handleSliderChange = (idx: number, newPct: number) => {
  const n = props.tags.length
  if (n <= 1) return

  const newRatio = newPct / 100
  const oldRatio = ratios.value[idx]
  const delta = newRatio - oldRatio

  if (Math.abs(delta) < 0.001) return

  // 其余标签等比例缩放，保证总和 = 1
  const others = ratios.value.map((v, i) => (i === idx ? 0 : v))
  const othersSum = others.reduce((a, b) => a + b, 0)

  const newArr = ratios.value.map((v, i) => {
    if (i === idx) return newRatio
    if (othersSum < 0.0001) return (1 - newRatio) / (n - 1)
    return Math.max(0, v - delta * (v / othersSum))
  })

  // 修正浮点误差，确保总和精确为 1
  const total = newArr.reduce((a, b) => a + b, 0)
  ratios.value = newArr.map(v => v / total)
}

const handleReset = () => {
  const n = props.tags.length
  ratios.value = props.tags.map(() => 1 / n)
}

const handleSave = async () => {
  saving.value = true
  try {
    const ratioMap: Record<string, number> = {}
    props.tags.forEach((tag, i) => {
      ratioMap[tag] = ratios.value[i]
    })
    await DailyReportService.SaveTagRatios(props.date, ratioMap)
    message.success('已保存')
    emit('saved')
  } catch (e) {
    message.error('保存失败')
    console.error(e)
  } finally {
    saving.value = false
  }
}

// 格式化天数显示
const formatDays = (ratio: number): string => {
  const d = ratio * 1
  if (d === 0) return '0'
  if (Math.abs(d - Math.round(d)) < 0.01) return Math.round(d).toString()
  return d.toFixed(2)
}

// 颜色池（与父组件保持一致）
const TAG_COLORS = [
  { name: 'blue',     hex: '#1890ff' },
  { name: 'green',    hex: '#52c41a' },
  { name: 'orange',   hex: '#fa8c16' },
  { name: 'purple',   hex: '#722ed1' },
  { name: 'cyan',     hex: '#13c2c2' },
  { name: 'red',      hex: '#f5222d' },
  { name: 'gold',     hex: '#faad14' },
  { name: 'lime',     hex: '#a0d911' },
  { name: 'geekblue', hex: '#2f54eb' },
  { name: 'magenta',  hex: '#eb2f96' },
]
const tagColorCache = new Map<string, number>()
let tagColorIndex = 0
const getTagColorIdx = (tag: string) => {
  if (!tagColorCache.has(tag)) {
    tagColorCache.set(tag, tagColorIndex % TAG_COLORS.length)
    tagColorIndex++
  }
  return tagColorCache.get(tag)!
}
const getTagColor    = (tag: string) => TAG_COLORS[getTagColorIdx(tag)].name
const getTagColorHex = (tag: string) => TAG_COLORS[getTagColorIdx(tag)].hex
</script>

<style scoped>
.ratio-editor {
  background: #f8faff;
  border: 1px solid #d6e4ff;
  border-radius: 8px;
  padding: 12px 16px;
  margin-bottom: 4px;
}

.ratio-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 10px;
}

.ratio-date {
  font-weight: 700;
  font-size: 13px;
  color: #262626;
}

.ratio-hint {
  font-size: 11px;
  color: #8c8c8c;
}

/* 色块进度条 */
.ratio-bar {
  display: flex;
  height: 10px;
  border-radius: 5px;
  overflow: hidden;
  margin-bottom: 12px;
  background: #f0f0f0;
}

.ratio-bar-segment {
  height: 100%;
  transition: width 0.2s ease;
  min-width: 2px;
}

/* 滑块行 */
.ratio-rows {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.ratio-row {
  display: flex;
  align-items: center;
  gap: 10px;
}

.ratio-tag-label {
  min-width: 80px;
  text-align: right;
}

.ratio-slider-wrap {
  flex: 1;
  display: flex;
  align-items: center;
}

.ratio-days-label {
  min-width: 48px;
  text-align: right;
  font-size: 12px;
  font-weight: 600;
  color: #595959;
}

.ratio-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  margin-top: 10px;
}
</style>
