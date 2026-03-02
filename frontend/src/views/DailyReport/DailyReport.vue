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
              'is-weekend': day.isWeekend
            }"
            v-for="day in weekDays"
            :key="day.date"
            @click="handleDayClick(day)"
          >
            <div class="day-header">
              <div class="day-name">{{ day.dayName }}</div>
              <div class="day-date">{{ day.dayNumber }}</div>
              <div v-if="day.isToday" class="today-badge">今天</div>
            </div>
            <div class="day-content">
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
            <div
              class="day-card-mini"
              :class="{
                'is-today': day.isToday,
                'has-report': day.hasReport,
                'is-current-month': day.isCurrentMonth,
                'is-weekend': day.isWeekend
              }"
              v-for="day in week"
              :key="day.date"
              @click="handleDayClick(day)"
            >
              <div class="day-number">{{ day.dayNumber }}</div>
              <div v-if="day.hasReport" class="report-indicator">
                <CheckCircleFilled class="icon-success" />
              </div>
              <div v-if="day.isToday" class="today-dot"></div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 日报编辑弹窗 -->
    <a-modal
      v-model:open="editModalVisible"
      :title="editModalTitle"
      width="600px"
      @ok="handleSaveReport"
      @cancel="handleCancelEdit"
      :bodyStyle="{ padding: '20px', maxHeight: '60vh', overflowY: 'auto' }"
    >
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
            style="width: 100%"
          />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import {
  LeftOutlined,
  RightOutlined,
  CheckCircleOutlined,
  CheckCircleFilled,
  EditOutlined,
  CalendarOutlined,
  InfoCircleOutlined
} from '@ant-design/icons-vue'
import dayjs, { Dayjs } from 'dayjs'
import * as DailyReportService from '../../../bindings/github.com/Aliuyanfeng/happytools/backend/services/dailyreport/dailyreportservice'

// 视图模式
const viewMode = ref<'week' | 'month'>('week')

// 当前日期
const currentDate = ref<Dayjs>(dayjs())

// 日报数据
const reports = ref<Map<string, any>>(new Map())

// 编辑弹窗
const editModalVisible = ref(false)
const editModalTitle = ref('编辑日报')
const editForm = ref({
  date: dayjs().format('YYYY-MM-DD'),
  content: '',
  summary: '',
  tags: [] as string[],
  isPeriod: false,
  periodRange: [] as string[]
})

// 星期名称
const dayNames = ['周一', '周二', '周三', '周四', '周五', '周六', '周日']

// 日期显示
const dateDisplay = computed(() => {
  if (viewMode.value === 'week') {
    const startOfWeek = currentDate.value.startOf('week')
    const endOfWeek = currentDate.value.endOf('week')
    return `${startOfWeek.format('YYYY年MM月DD日')} - ${endOfWeek.format('MM月DD日')}`
  } else {
    return currentDate.value.format('YYYY年MM月')
  }
})

// 周视图数据
const weekDays = computed(() => {
  const startOfWeek = currentDate.value.startOf('week')
  const days: any[] = []

  for (let i = 0; i < 7; i++) {
    const day = startOfWeek.add(i, 'day')
    const dateStr = day.format('YYYY-MM-DD')
    const report = reports.value.get(dateStr)
    const dayOfWeek = day.day()

    days.push({
      date: dateStr,
      dayName: dayNames[i],
      dayNumber: day.format('DD'),
      isToday: day.isSame(dayjs(), 'day'),
      isCurrentMonth: day.month() === currentDate.value.month(),
      isWeekend: dayOfWeek === 0 || dayOfWeek === 6,
      hasReport: !!report,
      previewText: report ? report.content.substring(0, 80) + (report.content.length > 80 ? '...' : '') : '',
      report: report
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
      hasReport: false
    })
  }

  // 填充本月日期
  for (let day = 1; day <= daysInMonth; day++) {
    const date = currentDate.value.date(day)
    const dateStr = date.format('YYYY-MM-DD')
    const report = reports.value.get(dateStr)
    const dayOfWeek = date.day()

    week.push({
      date: dateStr,
      dayNumber: day.toString(),
      isToday: date.isSame(dayjs(), 'day'),
      isCurrentMonth: true,
      isWeekend: dayOfWeek === 0 || dayOfWeek === 6,
      hasReport: !!report,
      report: report
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
        hasReport: false
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
      startDate = currentDate.value.startOf('week').format('YYYY-MM-DD')
      endDate = currentDate.value.endOf('week').format('YYYY-MM-DD')
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
}

// 处理上一周/月
const handlePrev = () => {
  if (viewMode.value === 'week') {
    currentDate.value = currentDate.value.subtract(1, 'week')
  } else {
    currentDate.value = currentDate.value.subtract(1, 'month')
  }
  loadReports()
}

// 处理下一周/月
const handleNext = () => {
  if (viewMode.value === 'week') {
    currentDate.value = currentDate.value.add(1, 'week')
  } else {
    currentDate.value = currentDate.value.add(1, 'month')
  }
  loadReports()
}

// 处理回到今天
const handleToday = () => {
  currentDate.value = dayjs()
  loadReports()
}

// 处理日期点击
const handleDayClick = async (day: any) => {
  editForm.value.date = day.date
  editForm.value.isPeriod = false
  editForm.value.periodRange = []
  
  if (day.hasReport && day.report) {
    editModalTitle.value = `编辑日报 - ${day.date}`
    editForm.value.content = day.report.content
    editForm.value.summary = day.report.summary || ''
    editForm.value.tags = day.report.tags || []
  } else {
    editModalTitle.value = `新建日报 - ${day.date}`
    editForm.value.content = ''
    editForm.value.summary = ''
    editForm.value.tags = []
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
      message.success('保存成功')
    }

    editModalVisible.value = false
    loadReports()
  } catch (error) {
    message.error('保存失败')
    console.error(error)
  }
}

// 取消编辑
const handleCancelEdit = () => {
  editModalVisible.value = false
}

// 初始化
onMounted(() => {
  loadReports()
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
  gap: 8px;
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
  height: calc(100vh - 280px);
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
</style>
