<template>
  <div class="daily-report">
    <!-- 顶部工具栏 -->
    <div class="toolbar">
      <a-row :gutter="16" align="middle" justify="space-between">
        <a-col>
          <a-radio-group v-model:value="viewMode" button-style="solid" @change="handleViewModeChange">
            <a-radio-button value="week">周视图</a-radio-button>
            <a-radio-button value="month">月视图</a-radio-button>
          </a-radio-group>
        </a-col>
        <a-col>
          <a-space size="middle">
            <a-button @click="handlePrev" size="large">
              <LeftOutlined />
            </a-button>
            <div class="date-display">
              <CalendarOutlined class="mr-2" />
              {{ dateDisplay }}
            </div>
            <a-button @click="handleNext" size="large">
              <RightOutlined />
            </a-button>
            <a-button type="primary" @click="handleToday" size="large">
              今天
            </a-button>
            <a-button @click="handleOpenTagStats" size="large">
              <BarChartOutlined />
              月度工时统计
            </a-button>
          </a-space>
        </a-col>
      </a-row>
    </div>

    <!-- 日历视图 -->
    <div class="calendar-container">
      <!-- 周视图 -->
      <div v-if="viewMode === 'week'" class="week-view">
        <div class="week-row-centered">
          <div
            class="day-card"
            :class="{
              'is-today': day.isToday,
              'has-report': day.hasReport,
              'is-weekend': day.isWeekend,
              'is-holiday': day.isHoliday,
              'is-workday': day.isWorkday
            }"
            v-for="day in weekDays"
            :key="day.date"
            @click="handleDayClick(day)"
          >
            <div class="day-header">
              <div class="day-name">{{ day.dayName }}</div>
              <div class="day-date">{{ day.dayNumber }}</div>
              <div v-if="day.isToday" class="today-badge">今天</div>
              <div v-else-if="day.isHoliday" class="holiday-week-badge">假</div>
              <div v-else-if="day.isWorkday" class="workday-week-badge">班</div>
            </div>
            <div class="day-content">
              <!-- 节日名称 -->
              <div v-if="day.holidayName" class="week-holiday-name">{{ day.holidayName }}</div>
              <div v-if="day.hasReport" class="report-preview">
                <div class="status-badge success">
                  <CheckCircleOutlined />
                  <span>已记录</span>
                </div>
                <div class="preview-text">{{ day.previewText }}</div>
                <div class="report-meta">
                  <span v-if="day.report.summary" class="summary-tag">
                    {{ day.report.summary }}
                  </span>
                  <div v-if="day.report.tags && day.report.tags.length" class="tags-row">
                    <a-tag v-for="tag in day.report.tags" :key="tag" size="small">{{ tag }}</a-tag>
                  </div>
                </div>
              </div>
              <div v-else class="no-report">
                <div class="status-badge empty">
                  <EditOutlined />
                  <span>点击记录</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 月视图 -->
      <div v-else class="month-view">
        <div class="month-header">
          <div class="day-name-header" v-for="dayName in dayNames" :key="dayName">
            {{ dayName }}
          </div>
        </div>
        <div class="month-body">
          <div class="month-week" v-for="(week, weekIndex) in monthWeeks" :key="weekIndex">
            <a-tooltip
              v-for="day in week"
              :key="day.date"
              :title="day.holidayName || undefined"
              placement="top"
            >
              <div
                class="day-card-mini"
                :class="{
                  'is-today': day.isToday,
                  'has-report': day.hasReport,
                  'is-current-month': day.isCurrentMonth,
                  'is-weekend': day.isWeekend,
                  'is-holiday': day.isHoliday,
                  'is-workday': day.isWorkday
                }"
                @click="handleDayClick(day)"
              >
                <!-- 节假日角标 -->
                <span v-if="day.isHoliday && day.isCurrentMonth" class="holiday-badge">假</span>
                <span v-if="day.isWorkday && day.isCurrentMonth" class="workday-badge">班</span>

                <div class="day-number">{{ day.dayNumber }}</div>
                <div v-if="day.holidayName && day.isCurrentMonth" class="holiday-name">{{ day.holidayName }}</div>
                <div v-if="day.hasReport" class="report-indicator">
                  <CheckCircleFilled class="icon-success" />
                </div>
                <div v-if="day.isToday" class="today-dot"></div>
              </div>
            </a-tooltip>
          </div>
        </div>
      </div>
    </div>

    <!-- 日报编辑弹窗 -->
    <a-modal
      v-model:open="editModalVisible"
      :title="editModalTitle"
      width="600px"
      @cancel="handleCancelEdit"
      :bodyStyle="{ padding: '20px', maxHeight: '60vh', overflowY: 'auto' }"
    >
      <template #footer>
        <div style="display: flex; justify-content: space-between; align-items: center;">
          <a-button
            v-if="editForm.isEditing"
            danger
            @click="handleDeleteReport"
          >
            <DeleteOutlined /> 删除日报
          </a-button>
          <div v-else />
          <a-space>
            <a-button @click="handleCancelEdit">取消</a-button>
            <a-button type="primary" @click="handleSaveReport">保存</a-button>
          </a-space>
        </div>
      </template>
      <a-form layout="vertical">
        <a-row :gutter="12">
          <a-col :span="12">
            <a-form-item label="日期">
              <a-date-picker
                v-model:value="editForm.date"
                format="YYYY-MM-DD"
                valueFormat="YYYY-MM-DD"
                style="width: 100%"
                :disabled="editForm.isPeriod"
              />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="摘要">
              <a-input
                v-model:value="editForm.summary"
                placeholder="请输入日报摘要（可选）"
              />
            </a-form-item>
          </a-col>
        </a-row>
        
        <a-form-item>
          <a-checkbox v-model:checked="editForm.isPeriod">
            设置为周期日报
          </a-checkbox>
          <div class="period-hint" v-if="editForm.isPeriod">
            <InfoCircleOutlined class="mr-1" />
            周期日报将为选定日期范围内的所有日期创建相同的日报内容
          </div>
        </a-form-item>

        <a-form-item label="周期日期范围" v-if="editForm.isPeriod">
          <a-range-picker
            v-model:value="editForm.periodRange"
            format="YYYY-MM-DD"
            valueFormat="YYYY-MM-DD"
            style="width: 100%"
            :placeholder="['开始日期', '结束日期']"
          />
        </a-form-item>

        <a-form-item label="日报内容">
          <a-textarea
            v-model:value="editForm.content"
            :rows="6"
            placeholder="请输入今日工作内容..."
          />
        </a-form-item>
        <a-form-item label="标签">
          <a-select
            v-model:value="editForm.tags"
            mode="tags"
            placeholder="添加标签（可选）"
            :options="[...new Set([...allTags, ...editForm.tags])].map(t => ({ value: t, label: t }))"
            style="width: 100%"
          />
        </a-form-item>
      </a-form>
    </a-modal>
    <!-- 月度工时统计弹窗 -->
    <a-modal
      v-model:open="tagStatsVisible"
      title="月度工时统计 · 标签分布"
      width="800px"
      :footer="null"
      :bodyStyle="{ padding: '16px 24px 24px', maxHeight: '75vh', overflowY: 'auto' }"
    >
      <!-- 月份选择器 -->
      <div class="stats-month-picker">
        <a-month-picker
          v-model:value="statsMonth"
          format="YYYY-MM"
          valueFormat="YYYY-MM"
          :allow-clear="false"
          @change="handleStatsMonthChange"
          style="width: 160px"
        />
      </div>

      <div v-if="tagStatsLoading" class="stats-loading">
        <a-spin tip="加载中..." />
      </div>
      <div v-else-if="!currentMonthStat || currentMonthStat.total_days === 0" class="stats-empty">
        <a-empty description="本月暂无日报数据" />
      </div>
      <div v-else class="stats-content">
        <div class="month-stat-block">
          <!-- 月份标题 -->
          <div class="month-stat-header">
            <span class="month-stat-title">{{ currentMonthStat.month }}</span>
            <span class="month-stat-total">共 {{ currentMonthStat.total_days }} 天日报</span>
          </div>

          <!-- 标签统计 -->
          <div v-if="currentMonthStat.tag_stats && currentMonthStat.tag_stats.length" class="tag-stat-list">
            <div
              v-for="ts in currentMonthStat.tag_stats"
              :key="ts.tag"
              class="tag-stat-row"
            >
              <span class="tag-stat-name">
                <a-tag :color="getTagColor(ts.tag)">{{ ts.tag }}</a-tag>
              </span>
              <div class="tag-stat-bar-wrap">
                <div
                  class="tag-stat-bar"
                  :style="{ width: getBarWidth(ts.days, currentMonthStat.total_days), backgroundColor: getTagColorHex(ts.tag) }"
                ></div>
              </div>
              <span class="tag-stat-days">{{ formatDays(ts.days) }} 天</span>
            </div>
          </div>
          <div v-else class="no-tags-hint">本月所有日报均未打标签</div>

          <!-- 多标签天的工时比例编辑器 -->
          <div v-if="currentMonthStat.multi_tag_dates && currentMonthStat.multi_tag_dates.length" class="multi-tag-section">
            <div
              class="multi-tag-section-title"
              @click="toggleMultiTagExpand(currentMonthStat.month)"
            >
              <SlidersOutlined class="multi-tag-icon" />
              本月有 {{ currentMonthStat.multi_tag_dates.length }} 天包含多个标签，点击调整工时占比
              <DownOutlined :class="['expand-arrow', { expanded: expandedMonths.has(currentMonthStat.month) }]" />
            </div>
            <div v-if="expandedMonths.has(currentMonthStat.month)" class="multi-tag-editors">
              <TagRatioEditor
                v-for="item in currentMonthStat.multi_tag_dates"
                :key="item.date"
                :date="item.date"
                :tags="item.tags"
                :saved-ratios="item.savedRatios"
                @saved="handleRatioSaved(currentMonthStat.month)"
              />
            </div>
          </div>

          <!-- 未打标签日期 -->
          <div v-if="currentMonthStat.untagged_dates && currentMonthStat.untagged_dates.length" class="untagged-section">
            <div class="untagged-title">
              <ExclamationCircleOutlined class="untagged-icon" />
              未打标签日期（{{ currentMonthStat.untagged_dates.length }} 天）
            </div>
            <div class="untagged-dates">
              <a-tag
                v-for="d in currentMonthStat.untagged_dates"
                :key="d"
                class="untagged-date-tag"
                @click="handleClickUntaggedDate(d)"
              >{{ d }}</a-tag>
            </div>
          </div>
        </div>
      </div>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { message, Modal } from 'ant-design-vue'
import {
  LeftOutlined,
  RightOutlined,
  CheckCircleOutlined,
  CheckCircleFilled,
  EditOutlined,
  CalendarOutlined,
  InfoCircleOutlined,
  DeleteOutlined,
  BarChartOutlined,
  ExclamationCircleOutlined,
  SlidersOutlined,
  DownOutlined
} from '@ant-design/icons-vue'
import dayjs, { Dayjs } from 'dayjs'
import * as DailyReportService from '../../../bindings/github.com/Aliuyanfeng/happytools/backend/services/dailyreport/dailyreportservice'
import TagRatioEditor from './TagRatioEditor.vue'
import { getMonthHolidayMap, type DayHolidayInfo } from '../../utils/holidays'

// 视图模式
const viewMode = ref<'week' | 'month'>('week')

// 当前日期
const currentDate = ref<Dayjs>(dayjs())

// 日报数据
const reports = ref<Map<string, any>>(new Map())

// 节假日数据（月视图用）：date -> DayHolidayInfo
const holidayMap = ref<Map<string, DayHolidayInfo>>(new Map())

// 加载节假日数据（周视图可能跨月，需同时加载相邻月）
const loadHolidays = async () => {
  try {
    if (viewMode.value === 'month') {
      const year = currentDate.value.year()
      const month = currentDate.value.month() + 1
      holidayMap.value = await getMonthHolidayMap(year, month)
    } else {
      // 周视图：加载本周涉及的所有月份
      const startOfWeek = currentDate.value.startOf('week').add(1, 'day')
      const endOfWeek = startOfWeek.add(6, 'day')
      const maps = await Promise.all([
        getMonthHolidayMap(startOfWeek.year(), startOfWeek.month() + 1),
        // 若跨月则也加载结束月
        startOfWeek.month() !== endOfWeek.month()
          ? getMonthHolidayMap(endOfWeek.year(), endOfWeek.month() + 1)
          : Promise.resolve(new Map())
      ])
      const merged = new Map<string, DayHolidayInfo>()
      maps.forEach(m => m.forEach((v, k) => merged.set(k, v)))
      holidayMap.value = merged
    }
  } catch (e) {
    console.error('加载节假日数据失败', e)
    holidayMap.value = new Map()
  }
}

// 编辑弹窗
const editModalVisible = ref(false)
const editModalTitle = ref('编辑日报')
const editForm = ref({
  date: dayjs().format('YYYY-MM-DD'),
  originalDate: '',
  content: '',
  summary: '',
  tags: [] as string[],
  isPeriod: false,
  periodRange: [] as string[],
  isEditing: false,
  reportId: 0
})

// 全局标签列表
const allTags = ref<string[]>([])

// 月度工时统计弹窗
const tagStatsVisible = ref(false)
const tagStatsLoading = ref(false)
const monthlyTagStats = ref<any[]>([])   // 保留兼容，不再使用
const currentMonthStat = ref<any>(null)
const statsMonth = ref<string>(dayjs().format('YYYY-MM'))

// 标签颜色池
const TAG_COLORS = [
  { name: 'blue', hex: '#1890ff' },
  { name: 'green', hex: '#52c41a' },
  { name: 'orange', hex: '#fa8c16' },
  { name: 'purple', hex: '#722ed1' },
  { name: 'cyan', hex: '#13c2c2' },
  { name: 'red', hex: '#f5222d' },
  { name: 'gold', hex: '#faad14' },
  { name: 'lime', hex: '#a0d911' },
  { name: 'geekblue', hex: '#2f54eb' },
  { name: 'magenta', hex: '#eb2f96' },
]

const tagColorCache = new Map<string, number>()
let tagColorIndex = 0

const getTagColorIdx = (tag: string): number => {
  if (!tagColorCache.has(tag)) {
    tagColorCache.set(tag, tagColorIndex % TAG_COLORS.length)
    tagColorIndex++
  }
  return tagColorCache.get(tag)!
}

const getTagColor = (tag: string): string => TAG_COLORS[getTagColorIdx(tag)].name
const getTagColorHex = (tag: string): string => TAG_COLORS[getTagColorIdx(tag)].hex

const getBarWidth = (days: number, total: number): string => {
  if (!total) return '0%'
  return Math.round((days / total) * 100) + '%'
}

const handleOpenTagStats = async () => {
  tagStatsVisible.value = true
  await loadMonthTagStats(statsMonth.value)
}

const handleStatsMonthChange = async (month: string) => {
  await loadMonthTagStats(month)
}

const loadMonthTagStats = async (month: string) => {
  tagStatsLoading.value = true
  try {
    const ms = await DailyReportService.GetMonthTagStats(month)
    if (!ms) {
      currentMonthStat.value = null
      return
    }

    // 找出有 2+ 个标签的日期，加载已保存比例
    const startDate = `${month}-01`
    const [year, mon] = month.split('-')
    const lastDay = new Date(parseInt(year), parseInt(mon), 0).getDate()
    const endDate = `${month}-${String(lastDay).padStart(2, '0')}`
    const reports = await DailyReportService.GetRange(startDate, endDate)

    const multiTagDates: Array<{ date: string; tags: string[]; savedRatios: Record<string, number> }> = []
    for (const r of (reports || [])) {
      const validTags = (r.tags || []).filter((t: string) => t !== '')
      if (validTags.length >= 2) {
        const savedRatios = await DailyReportService.GetTagRatios(r.date)
        multiTagDates.push({ date: r.date, tags: validTags, savedRatios: savedRatios || {} })
      }
    }
    multiTagDates.sort((a, b) => a.date.localeCompare(b.date))

    currentMonthStat.value = { ...ms, multi_tag_dates: multiTagDates }
  } catch (e) {
    message.error('加载统计数据失败')
    console.error(e)
  } finally {
    tagStatsLoading.value = false
  }
}

// 展开/收起某月的多标签编辑器
const expandedMonths = ref<Set<string>>(new Set())
const toggleMultiTagExpand = (month: string) => {
  if (expandedMonths.value.has(month)) {
    expandedMonths.value.delete(month)
  } else {
    expandedMonths.value.add(month)
  }
  // 触发响应式更新
  expandedMonths.value = new Set(expandedMonths.value)
}

// 某天比例保存后，刷新当前月统计
const handleRatioSaved = async (month: string) => {
  await loadMonthTagStats(month)
  expandedMonths.value.add(month)
  expandedMonths.value = new Set(expandedMonths.value)
}

// 格式化天数：整数不显示小数，小数保留2位
const formatDays = (days: number): string => {
  if (days === 0) return '0'
  if (Math.abs(days - Math.round(days)) < 0.005) return Math.round(days).toString()
  return days.toFixed(2)
}

// 点击未打标签日期，跳转到对应日期并打开编辑弹窗
const handleClickUntaggedDate = async (date: string) => {
  tagStatsVisible.value = false
  currentDate.value = dayjs(date)
  await loadReports()
  // 找到该日期的日报并打开编辑
  const report = reports.value.get(date)
  handleDayClick({
    date,
    hasReport: !!report,
    report
  })
}

const loadAllTags = async () => {
  try {
    const tags = await DailyReportService.GetAllTags()
    allTags.value = tags || []
  } catch (e) {
    console.error('加载标签失败', e)
  }
}

// 星期名称
const dayNames = ['周一', '周二', '周三', '周四', '周五', '周六', '周日']

// 日期显示
const dateDisplay = computed(() => {
  if (viewMode.value === 'week') {
    const startOfWeek = currentDate.value.startOf('week').add(1, 'day')
    const endOfWeek = startOfWeek.add(6, 'day')
    return `${startOfWeek.format('YYYY年MM月DD日')} - ${endOfWeek.format('MM月DD日')}`
  } else {
    return currentDate.value.format('YYYY年MM月')
  }
})

// 周视图数据
const weekDays = computed(() => {
  const startOfWeek = currentDate.value.startOf('week').add(1, 'day')
  const days: any[] = []

  for (let i = 0; i < 7; i++) {
    const day = startOfWeek.add(i, 'day')
    const dateStr = day.format('YYYY-MM-DD')
    const report = reports.value.get(dateStr)
    const dayOfWeek = day.day()
    const holiday = holidayMap.value.get(dateStr)

    days.push({
      date: dateStr,
      dayName: dayNames[i],
      dayNumber: day.format('DD'),
      isToday: day.isSame(dayjs(), 'day'),
      isCurrentMonth: day.month() === currentDate.value.month(),
      isWeekend: dayOfWeek === 0 || dayOfWeek === 6,
      hasReport: !!report,
      previewText: report ? report.content.substring(0, 80) + (report.content.length > 80 ? '...' : '') : '',
      report: report,
      isHoliday: holiday?.isHoliday ?? false,
      isWorkday: holiday?.isWorkday ?? false,
      holidayName: holiday?.name ?? ''
    })
  }

  return days
})

// 月视图数据
const monthWeeks = computed(() => {
  const startOfMonth = currentDate.value.startOf('month')
  const endOfMonth = currentDate.value.endOf('month')
  const startDay = startOfMonth.day() || 7
  const daysInMonth = endOfMonth.date()

  const weeks: any[] = []
  let week: any[] = []

  // 填充月初空白
  for (let i = 1; i < startDay; i++) {
    const prevDay = startOfMonth.subtract(startDay - i, 'day')
    week.push({
      date: prevDay.format('YYYY-MM-DD'),
      dayNumber: prevDay.format('DD'),
      isToday: false,
      isCurrentMonth: false,
      isWeekend: false,
      hasReport: false,
      isHoliday: false,
      isWorkday: false,
      holidayName: ''
    })
  }

  // 填充本月日期
  for (let day = 1; day <= daysInMonth; day++) {
    const date = currentDate.value.date(day)
    const dateStr = date.format('YYYY-MM-DD')
    const report = reports.value.get(dateStr)
    const dayOfWeek = date.day()
    const holiday = holidayMap.value.get(dateStr)

    week.push({
      date: dateStr,
      dayNumber: day.toString(),
      isToday: date.isSame(dayjs(), 'day'),
      isCurrentMonth: true,
      isWeekend: dayOfWeek === 0 || dayOfWeek === 6,
      hasReport: !!report,
      report: report,
      isHoliday: holiday?.isHoliday ?? false,
      isWorkday: holiday?.isWorkday ?? false,
      holidayName: holiday?.name ?? ''
    })

    if (week.length === 7) {
      weeks.push(week)
      week = []
    }
  }

  // 填充月末空白
  if (week.length > 0) {
    const remaining = 7 - week.length
    for (let i = 1; i <= remaining; i++) {
      const nextDay = endOfMonth.add(i, 'day')
      week.push({
        date: nextDay.format('YYYY-MM-DD'),
        dayNumber: nextDay.format('DD'),
        isToday: false,
        isCurrentMonth: false,
        isWeekend: false,
        hasReport: false,
        isHoliday: false,
        isWorkday: false,
        holidayName: ''
      })
    }
    weeks.push(week)
  }

  return weeks
})

// 加载日报数据
const loadReports = async () => {
  try {
    let startDate: string
    let endDate: string

    if (viewMode.value === 'week') {
      startDate = currentDate.value.startOf('week').add(1, 'day').format('YYYY-MM-DD')
      endDate = currentDate.value.startOf('week').add(7, 'day').format('YYYY-MM-DD')
    } else {
      startDate = currentDate.value.startOf('month').format('YYYY-MM-DD')
      endDate = currentDate.value.endOf('month').format('YYYY-MM-DD')
    }

    const result = await DailyReportService.GetRange(startDate, endDate)
    
    reports.value.clear()
    result.forEach((report: any) => {
      reports.value.set(report.date, report)
    })
  } catch (error) {
    message.error('加载日报数据失败')
    console.error(error)
  }
}

// 处理视图模式切换
const handleViewModeChange = () => {
  loadReports()
  loadHolidays()
}

// 处理上一周/月
const handlePrev = () => {
  if (viewMode.value === 'week') {
    currentDate.value = currentDate.value.subtract(1, 'week')
  } else {
    currentDate.value = currentDate.value.subtract(1, 'month')
  }
  loadReports()
  loadHolidays()
}

// 处理下一周/月
const handleNext = () => {
  if (viewMode.value === 'week') {
    currentDate.value = currentDate.value.add(1, 'week')
  } else {
    currentDate.value = currentDate.value.add(1, 'month')
  }
  loadReports()
  loadHolidays()
}

// 处理回到今天
const handleToday = () => {
  currentDate.value = dayjs()
  loadReports()
  loadHolidays()
}

// 处理日期点击
const handleDayClick = async (day: any) => {
  editForm.value.date = day.date
  editForm.value.originalDate = day.date
  editForm.value.isPeriod = false
  editForm.value.periodRange = []
  
  if (day.hasReport && day.report) {
    editModalTitle.value = `编辑日报 - ${day.date}`
    editForm.value.content = day.report.content
    editForm.value.summary = day.report.summary || ''
    editForm.value.tags = day.report.tags ? [...day.report.tags] : []
    editForm.value.isEditing = true
    editForm.value.reportId = day.report.id
  } else {
    editModalTitle.value = `新建日报 - ${day.date}`
    editForm.value.content = ''
    editForm.value.summary = ''
    editForm.value.tags = []
    editForm.value.isEditing = false
    editForm.value.reportId = 0
  }
  
  editModalVisible.value = true
}

// 保存日报
const handleSaveReport = async () => {
  if (!editForm.value.content) {
    message.warning('请输入日报内容')
    return
  }

  try {
    // 如果是周期日报
    if (editForm.value.isPeriod && editForm.value.periodRange.length === 2) {
      const startDate = dayjs(editForm.value.periodRange[0])
      const endDate = dayjs(editForm.value.periodRange[1])
      const days = endDate.diff(startDate, 'day') + 1

      // 批量保存日报
      for (let i = 0; i < days; i++) {
        const date = startDate.add(i, 'day').format('YYYY-MM-DD')
        await DailyReportService.Save(
          date,
          editForm.value.content,
          editForm.value.summary,
          editForm.value.tags
        )
      }

      message.success(`已为 ${days} 天创建周期日报`)
    } else {
      // 单日日报
      await DailyReportService.Save(
        editForm.value.date,
        editForm.value.content,
        editForm.value.summary,
        editForm.value.tags
      )
      // 编辑模式下日期发生变化，删除原日期的日报（移动而非复制）
      if (editForm.value.isEditing && editForm.value.date !== editForm.value.originalDate) {
        await DailyReportService.Delete(editForm.value.reportId)
      }
      message.success('保存成功')
    }

    editModalVisible.value = false
    loadReports()
    loadAllTags()
  } catch (error) {
    message.error('保存失败')
    console.error(error)
  }
}

// 删除日报
const handleDeleteReport = () => {
  Modal.confirm({
    title: '确认删除',
    content: `确定要删除 ${editForm.value.date} 的日报吗？此操作不可恢复。`,
    okText: '删除',
    okType: 'danger',
    cancelText: '取消',
    async onOk() {
      try {
        await DailyReportService.Delete(editForm.value.reportId)
        message.success('删除成功')
        editModalVisible.value = false
        loadReports()
        loadAllTags()
      } catch (error) {
        message.error('删除失败')
        console.error(error)
      }
    }
  })
}

// 取消编辑
const handleCancelEdit = () => {
  editModalVisible.value = false
}

// 初始化
onMounted(() => {
  loadReports()
  loadAllTags()
  loadHolidays()
})
</script>

<style scoped>
.daily-report {
  height: calc(100vh - 120px);
  display: flex;
  flex-direction: column;
  padding: 0;
  background: #f5f7fa;
}

/* 顶部工具栏 */
.toolbar {
  background: #fff;
  padding: 20px 24px;
  border-bottom: 1px solid #e8e8e8;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  z-index: 10;
}

.date-display {
  font-size: 18px;
  font-weight: 600;
  color: #262626;
  min-width: 280px;
  text-align: center;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* 日历容器 */
.calendar-container {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
}

/* 周视图样式 */
.week-view {
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: flex-start;
}

.week-row-centered {
  display: flex;
  gap: 16px;
  justify-content: center;
  width: 100%;
  max-width: 1200px;
  padding: 0 16px;
}

.week-view .day-card {
  background: #fff;
  border: 2px solid #e8e8e8;
  border-radius: 12px;
  padding: 16px;
  width: 140px;
  min-width: 140px;
  height: calc(100vh - 220px);
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.week-view .day-card:hover {
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
  transform: translateY(-4px);
  border-color: #1890ff;
}

.week-view .day-card.is-today {
  border-color: #1890ff;
  background: linear-gradient(135deg, #e6f7ff 0%, #ffffff 100%);
  box-shadow: 0 4px 16px rgba(24, 144, 255, 0.2);
}

.week-view .day-card.has-report {
  border-color: #52c41a;
}

.week-view .day-card.is-weekend {
  background: #fafafa;
}

/* 周视图：法定假日 */
.week-view .day-card.is-holiday {
  border-color: #ffa39e;
  background: linear-gradient(135deg, #fff1f0 0%, #ffffff 100%);
}

.week-view .day-card.is-holiday .day-date {
  color: #cf1322;
}

/* 周视图：调休上班 */
.week-view .day-card.is-workday {
  border-color: #ffe58f;
  background: linear-gradient(135deg, #fffbe6 0%, #ffffff 100%);
}

.week-view .day-card.is-workday .day-date {
  color: #d46b08;
}

/* 周视图假/班角标 */
.holiday-week-badge,
.workday-week-badge {
  position: absolute;
  top: 0;
  right: 0;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 600;
}

.holiday-week-badge {
  background: #ff4d4f;
  color: #fff;
}

.workday-week-badge {
  background: #fa8c16;
  color: #fff;
}

/* 周视图节日名称 */
.week-holiday-name {
  display: inline-block;
  font-size: 11px;
  font-weight: 600;
  padding: 2px 8px;
  border-radius: 10px;
  margin-bottom: 8px;
  align-self: flex-start;
}

.is-holiday .week-holiday-name {
  background: #fff1f0;
  color: #cf1322;
  border: 1px solid #ffa39e;
}

.is-workday .week-holiday-name {
  background: #fffbe6;
  color: #d46b08;
  border: 1px solid #ffe58f;
}

.day-header {
  text-align: center;
  padding-bottom: 12px;
  border-bottom: 2px solid #f0f0f0;
  margin-bottom: 12px;
  position: relative;
}

.day-name {
  font-size: 13px;
  color: #8c8c8c;
  margin-bottom: 6px;
  font-weight: 500;
}

.day-date {
  font-size: 28px;
  font-weight: 700;
  color: #262626;
  line-height: 1;
}

.today-badge {
  position: absolute;
  top: 0;
  right: 0;
  background: #1890ff;
  color: #fff;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 600;
}

.day-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.status-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 10px;
  border-radius: 16px;
  font-size: 12px;
  font-weight: 600;
  margin-bottom: 8px;
}

.status-badge.success {
  background: #f6ffed;
  color: #52c41a;
  border: 1px solid #b7eb8f;
}

.status-badge.empty {
  background: #f5f5f5;
  color: #8c8c8c;
  border: 1px solid #d9d9d9;
}

.preview-text {
  font-size: 12px;
  color: #595959;
  line-height: 1.6;
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 5;
  -webkit-box-orient: vertical;
}

.report-meta {
  margin-top: 8px;
  padding-top: 8px;
  border-top: 1px solid #f0f0f0;
}

.tags-row {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  margin-top: 6px;
}

.summary-tag {
  display: inline-block;
  background: #f0f5ff;
  color: #1890ff;
  padding: 3px 8px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 500;
}

/* 月视图样式 */
.month-view {
  background: #fff;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.month-header {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 16px 0;
}

.day-name-header {
  text-align: center;
  font-weight: 600;
  color: #fff;
  font-size: 14px;
}

.month-body {
  padding: 8px;
}

.month-week {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 8px;
  margin-bottom: 8px;
}

.month-view .day-card-mini {
  background: #fff;
  border: 1px solid #e8e8e8;
  border-radius: 8px;
  padding: 12px;
  min-height: 80px;
  cursor: pointer;
  transition: all 0.3s;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  position: relative;
}

.month-view .day-card-mini:hover {
  background: #f5f5f5;
  border-color: #1890ff;
  transform: scale(1.05);
}

.month-view .day-card-mini.is-today {
  border-color: #1890ff;
  background: #e6f7ff;
  border-width: 2px;
}

.month-view .day-card-mini.is-current-month {
  color: #262626;
}

.month-view .day-card-mini:not(.is-current-month) {
  color: #bfbfbf;
  background: #fafafa;
}

.month-view .day-card-mini.is-weekend.is-current-month {
  background: #fafafa;
}

/* 法定假日：红色背景 */
.month-view .day-card-mini.is-holiday.is-current-month {
  background: #fff1f0;
  border-color: #ffa39e;
}

.month-view .day-card-mini.is-holiday.is-current-month .day-number {
  color: #cf1322;
}

/* 调休上班：黄色背景（周末但要上班） */
.month-view .day-card-mini.is-workday.is-current-month {
  background: #fffbe6;
  border-color: #ffe58f;
}

.month-view .day-card-mini.is-workday.is-current-month .day-number {
  color: #d46b08;
}

/* 假/班 角标 */
.holiday-badge,
.workday-badge {
  position: absolute;
  top: 4px;
  right: 4px;
  font-size: 10px;
  font-weight: 700;
  padding: 1px 4px;
  border-radius: 3px;
  line-height: 1.4;
}

.holiday-badge {
  background: #ff4d4f;
  color: #fff;
}

.workday-badge {
  background: #fa8c16;
  color: #fff;
}

/* 节日名称小字 */
.holiday-name {
  font-size: 10px;
  color: #cf1322;
  margin-top: 2px;
  margin-bottom: 2px;
  font-weight: 500;
  max-width: 100%;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  text-align: center;
}

.month-view .day-card-mini.is-workday.is-current-month .holiday-name {
  color: #d46b08;
}

.day-number {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 4px;
}

.report-indicator {
  margin-top: 4px;
}

.icon-success {
  font-size: 20px;
  color: #52c41a;
}

.today-dot {
  position: absolute;
  bottom: 8px;
  width: 6px;
  height: 6px;
  background: #1890ff;
  border-radius: 50%;
}

/* 滚动条样式 */
.calendar-container::-webkit-scrollbar {
  width: 8px;
}

.calendar-container::-webkit-scrollbar-track {
  background: #f5f5f5;
  border-radius: 4px;
}

.calendar-container::-webkit-scrollbar-thumb {
  background: #d9d9d9;
  border-radius: 4px;
}

.calendar-container::-webkit-scrollbar-thumb:hover {
  background: #bfbfbf;
}

/* 响应式调整 */
@media (max-width: 1400px) {
  .week-view .day-card {
    width: 130px;
    min-width: 130px;
    padding: 12px;
  }
  
  .day-date {
    font-size: 24px;
  }
  
  .preview-text {
    font-size: 11px;
  }
}

@media (max-width: 1200px) {
  .week-view .day-card {
    width: 120px;
    min-width: 120px;
    padding: 10px;
  }
  
  .day-date {
    font-size: 22px;
  }
  
  .day-name {
    font-size: 12px;
  }
}

/* 周期提示样式 */
.period-hint {
  margin-top: 8px;
  padding: 8px 12px;
  background: #e6f7ff;
  border: 1px solid #91d5ff;
  border-radius: 4px;
  font-size: 12px;
  color: #0050b3;
  display: flex;
  align-items: center;
}

.period-hint .mr-1 {
  margin-right: 4px;
}

/* 工时统计弹窗样式 */
.stats-month-picker {
  margin-bottom: 16px;
}

.stats-loading,
.stats-empty {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 40px 0;
}

.stats-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.month-stat-block {
  background: #fafafa;
  border: 1px solid #e8e8e8;
  border-radius: 10px;
  padding: 16px 20px;
}

.month-stat-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 14px;
  padding-bottom: 10px;
  border-bottom: 1px solid #e8e8e8;
}

.month-stat-title {
  font-size: 16px;
  font-weight: 700;
  color: #262626;
}

.month-stat-total {
  font-size: 13px;
  color: #8c8c8c;
}

.tag-stat-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.tag-stat-row {
  display: flex;
  align-items: center;
  gap: 10px;
}

.tag-stat-name {
  min-width: 120px;
  text-align: right;
}

.tag-stat-bar-wrap {
  flex: 1;
  height: 18px;
  background: #f0f0f0;
  border-radius: 9px;
  overflow: hidden;
}

.tag-stat-bar {
  height: 100%;
  border-radius: 9px;
  transition: width 0.4s ease;
  min-width: 4px;
}

.tag-stat-days {
  min-width: 44px;
  text-align: right;
  font-size: 13px;
  font-weight: 600;
  color: #595959;
}

.no-tags-hint {
  font-size: 13px;
  color: #bfbfbf;
  padding: 4px 0;
}

.untagged-section {
  margin-top: 14px;
  padding-top: 12px;
  border-top: 1px dashed #e8e8e8;
}

.untagged-title {
  font-size: 13px;
  color: #fa8c16;
  font-weight: 600;
  margin-bottom: 8px;
  display: flex;
  align-items: center;
  gap: 5px;
}

.untagged-icon {
  color: #fa8c16;
}

.untagged-dates {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.untagged-date-tag {
  cursor: pointer;
  border-color: #ffd591;
  color: #d46b08;
  background: #fff7e6;
  transition: all 0.2s;
}

.untagged-date-tag:hover {
  background: #ffd591;
  border-color: #fa8c16;
}

/* 多标签天编辑区域 */
.multi-tag-section {
  margin-top: 14px;
  padding-top: 12px;
  border-top: 1px dashed #d6e4ff;
}

.multi-tag-section-title {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: #1890ff;
  font-weight: 600;
  cursor: pointer;
  user-select: none;
  padding: 4px 0;
  transition: color 0.2s;
}

.multi-tag-section-title:hover {
  color: #096dd9;
}

.multi-tag-icon {
  color: #1890ff;
}

.expand-arrow {
  margin-left: auto;
  transition: transform 0.25s ease;
  font-size: 12px;
}

.expand-arrow.expanded {
  transform: rotate(180deg);
}

.multi-tag-editors {
  margin-top: 10px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}
</style>
