<!--
 * @Author: LiuYanFeng
 * @Date: 2026-02-12
 * @Description: 单位转换组件 - 包含字节转换、时间转换、时间差值计算
 -->
<template>
  <div class="unit-converter">
    <a-tabs v-model:activeKey="activeTab" type="card" class="converter-tabs">
      <!-- 字节转换标签页 -->
      <a-tab-pane key="byte" tab="字节转换">
        <div class="converter-section">
          <a-card title="字节单位转换" :bordered="false">
            <div class="converter-form">
              <a-row :gutter="16" align="middle">
                <a-col :span="8">
                  <a-input-number
                    v-model:value="byteValue"
                    :min="0"
                    :step="1"
                    :precision="6"
                    placeholder="输入数值"
                    style="width: 100%"
                    @change="handleByteConversion"
                  />
                </a-col>
                <a-col :span="6">
                  <a-select
                    v-model:value="byteFromUnit"
                    style="width: 100%"
                    @change="handleByteConversion"
                  >
                    <a-select-option v-for="unit in BYTE_UNITS" :key="unit.value" :value="unit.value">
                      {{ unit.label }}
                    </a-select-option>
                  </a-select>
                </a-col>
                <a-col :span="2" class="text-center">
                  <span class="converter-arrow">→</span>
                </a-col>
                <a-col :span="6">
                  <a-select
                    v-model:value="byteToUnit"
                    style="width: 100%"
                    @change="handleByteConversion"
                  >
                    <a-select-option v-for="unit in BYTE_UNITS" :key="unit.value" :value="unit.value">
                      {{ unit.label }}
                    </a-select-option>
                  </a-select>
                </a-col>
                <a-col :span="2">
                  <a-button type="primary" @click="swapByteUnits" :disabled="byteFromUnit === byteToUnit">
                    交换
                  </a-button>
                </a-col>
              </a-row>
              
              <div class="converter-result mt-4">
                <a-alert
                  v-if="byteResult !== null"
                  type="success"
                  :message="`转换结果: ${byteValue} ${byteFromUnit} = ${byteResult.toFixed(6)} ${byteToUnit}`"
                  show-icon
                />
                <a-alert
                  v-else
                  type="info"
                  message="请输入数值并选择单位进行转换"
                  show-icon
                />
              </div>
              
              <div class="quick-conversions mt-6">
                <h4>快速转换</h4>
                <a-row :gutter="8" class="mt-2">
                  <a-col :span="8" v-for="unit in BYTE_UNITS" :key="unit.value">
                    <a-button
                      block
                      @click="setByteToUnit(unit.value)"
                      :type="byteToUnit === unit.value ? 'primary' : 'default'"
                    >
                      {{ unit.label }}
                    </a-button>
                  </a-col>
                </a-row>
              </div>
            </div>
          </a-card>
        </div>
      </a-tab-pane>
      
      <!-- 时间转换标签页 -->
      <a-tab-pane key="time" tab="时间转换">
        <div class="converter-section">
          <a-card title="时间单位转换" :bordered="false">
            <div class="converter-form">
              <a-row :gutter="16" align="middle">
                <a-col :span="8">
                  <a-input-number
                    v-model:value="timeValue"
                    :min="0"
                    :step="1"
                    :precision="6"
                    placeholder="输入数值"
                    style="width: 100%"
                    @change="handleTimeConversion"
                  />
                </a-col>
                <a-col :span="6">
                  <a-select
                    v-model:value="timeFromUnit"
                    style="width: 100%"
                    @change="handleTimeConversion"
                  >
                    <a-select-option v-for="unit in TIME_UNITS" :key="unit.value" :value="unit.value">
                      {{ unit.label }}
                    </a-select-option>
                  </a-select>
                </a-col>
                <a-col :span="2" class="text-center">
                  <span class="converter-arrow">→</span>
                </a-col>
                <a-col :span="6">
                  <a-select
                    v-model:value="timeToUnit"
                    style="width: 100%"
                    @change="handleTimeConversion"
                  >
                    <a-select-option v-for="unit in TIME_UNITS" :key="unit.value" :value="unit.value">
                      {{ unit.label }}
                    </a-select-option>
                  </a-select>
                </a-col>
                <a-col :span="2">
                  <a-button type="primary" @click="swapTimeUnits" :disabled="timeFromUnit === timeToUnit">
                    交换
                  </a-button>
                </a-col>
              </a-row>
              
              <div class="converter-result mt-4">
                <a-alert
                  v-if="timeResult !== null"
                  type="success"
                  :message="`转换结果: ${timeValue} ${timeFromUnit} = ${timeResult.toFixed(6)} ${timeToUnit}`"
                  show-icon
                />
                <a-alert
                  v-else
                  type="info"
                  message="请输入数值并选择单位进行转换"
                  show-icon
                />
              </div>
              
              <div class="quick-conversions mt-6">
                <h4>快速转换</h4>
                <a-row :gutter="8" class="mt-2">
                  <a-col :span="8" v-for="unit in TIME_UNITS" :key="unit.value">
                    <a-button
                      block
                      @click="setTimeToUnit(unit.value)"
                      :type="timeToUnit === unit.value ? 'primary' : 'default'"
                    >
                      {{ unit.label }}
                    </a-button>
                  </a-col>
                </a-row>
              </div>
            </div>
          </a-card>
        </div>
      </a-tab-pane>
      
      <!-- 时间差值计算标签页 -->
      <a-tab-pane key="difference" tab="时间差值计算">
        <div class="converter-section">
          <a-card title="时间差值计算" :bordered="false">
            <div class="time-difference-form">
              <a-row :gutter="16" class="mb-4">
                <a-col :span="12">
                  <a-card title="开始时间" size="small">
                    <a-date-picker
                      v-model:value="startDate"
                      style="width: 100%"
                      placeholder="选择开始日期"
                      @change="calculateTimeDifference"
                    />
                    <a-time-picker
                      v-model:value="startTime"
                      style="width: 100%; margin-top: 8px"
                      placeholder="选择开始时间"
                      format="HH:mm:ss"
                      @change="calculateTimeDifference"
                    />
                    <div class="mt-2">
                      <a-input
                        v-model:value="startDateTimeString"
                        placeholder="或输入日期时间 (YYYY-MM-DD HH:mm:ss)"
                        @change="parseStartDateTime"
                      />
                    </div>
                  </a-card>
                </a-col>
                <a-col :span="12">
                  <a-card title="结束时间" size="small">
                    <a-date-picker
                      v-model:value="endDate"
                      style="width: 100%"
                      placeholder="选择结束日期"
                      @change="calculateTimeDifference"
                    />
                    <a-time-picker
                      v-model:value="endTime"
                      style="width: 100%; margin-top: 8px"
                      placeholder="选择结束时间"
                      format="HH:mm:ss"
                      @change="calculateTimeDifference"
                    />
                    <div class="mt-2">
                      <a-input
                        v-model:value="endDateTimeString"
                        placeholder="或输入日期时间 (YYYY-MM-DD HH:mm:ss)"
                        @change="parseEndDateTime"
                      />
                    </div>
                  </a-card>
                </a-col>
              </a-row>
              
              <div class="quick-actions mb-4">
                <a-space>
                  <a-button @click="setStartToNow">开始时间设为现在</a-button>
                  <a-button @click="setEndToNow">结束时间设为现在</a-button>
                  <a-button @click="swapTimes">交换时间</a-button>
                  <a-button type="primary" @click="calculateTimeDifference">计算差值</a-button>
                </a-space>
              </div>
              
              <div class="difference-results" v-if="timeDifferenceResult">
                <a-card title="计算结果" :bordered="false">
                  <a-descriptions bordered :column="2">
                    <a-descriptions-item label="总毫秒数">
                      {{ timeDifferenceResult.totalMilliseconds.toLocaleString() }} ms
                    </a-descriptions-item>
                    <a-descriptions-item label="格式化结果">
                      <span class="result-text">{{ timeDifferenceResult.formatted }}</span>
                    </a-descriptions-item>
                    <a-descriptions-item label="年">
                      {{ timeDifferenceResult.breakdown.years }}
                    </a-descriptions-item>
                    <a-descriptions-item label="月">
                      {{ timeDifferenceResult.breakdown.months }}
                    </a-descriptions-item>
                    <a-descriptions-item label="天">
                      {{ timeDifferenceResult.breakdown.days }}
                    </a-descriptions-item>
                    <a-descriptions-item label="小时">
                      {{ timeDifferenceResult.breakdown.hours }}
                    </a-descriptions-item>
                    <a-descriptions-item label="分钟">
                      {{ timeDifferenceResult.breakdown.minutes }}
                    </a-descriptions-item>
                    <a-descriptions-item label="秒">
                      {{ timeDifferenceResult.breakdown.seconds }}
                    </a-descriptions-item>
                    <a-descriptions-item label="毫秒">
                      {{ timeDifferenceResult.breakdown.milliseconds }}
                    </a-descriptions-item>
                  </a-descriptions>
                  
                  <div class="additional-info mt-4">
                    <a-row :gutter="16">
                      <a-col :span="8">
                        <a-statistic
                          title="总天数"
                          :value="Math.floor(timeDifferenceResult.totalMilliseconds / (1000 * 60 * 60 * 24))"
                          :precision="2"
                        />
                      </a-col>
                      <a-col :span="8">
                        <a-statistic
                          title="总小时数"
                          :value="Math.floor(timeDifferenceResult.totalMilliseconds / (1000 * 60 * 60))"
                          :precision="2"
                        />
                      </a-col>
                      <a-col :span="8">
                        <a-statistic
                          title="总分钟数"
                          :value="Math.floor(timeDifferenceResult.totalMilliseconds / (1000 * 60))"
                          :precision="2"
                        />
                      </a-col>
                    </a-row>
                  </div>
                </a-card>
              </div>
              
              <div class="difference-results" v-else>
                <a-alert
                  type="info"
                  message="请选择或输入开始时间和结束时间"
                  show-icon
                />
              </div>
            </div>
          </a-card>
        </div>
      </a-tab-pane>
    </a-tabs>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { 
  BYTE_UNITS, 
  TIME_UNITS, 
  convertBytes, 
  convertTime, 
  calculateTimeDifference as calcTimeDiff,
  parseTimeString,
  formatDateTime,
  type ByteUnit,
  type TimeUnit,
  type TimeDifferenceResult 
} from '@/utils/unitConverter';
import { Dayjs } from 'dayjs';
import dayjs from 'dayjs';

// ==================== 字节转换相关 ====================
const activeTab = ref<'byte' | 'time' | 'difference'>('byte');
const byteValue = ref<number>(1024);
const byteFromUnit = ref<ByteUnit>('KB');
const byteToUnit = ref<ByteUnit>('MB');
const byteResult = ref<number | null>(null);

// 字节转换处理
const handleByteConversion = () => {
  if (byteValue.value === null || byteValue.value === undefined) {
    byteResult.value = null;
    return;
  }
  
  try {
    byteResult.value = convertBytes(byteValue.value, byteFromUnit.value, byteToUnit.value);
  } catch (error) {
    console.error('字节转换错误:', error);
    byteResult.value = null;
  }
};

// 交换字节单位
const swapByteUnits = () => {
  const temp = byteFromUnit.value;
  byteFromUnit.value = byteToUnit.value;
  byteToUnit.value = temp;
  handleByteConversion();
};

// 设置目标字节单位
const setByteToUnit = (unit: ByteUnit) => {
  byteToUnit.value = unit;
  handleByteConversion();
};

// ==================== 时间转换相关 ====================
const timeValue = ref<number>(3600);
const timeFromUnit = ref<TimeUnit>('s');
const timeToUnit = ref<TimeUnit>('h');
const timeResult = ref<number | null>(null);

// 时间转换处理
const handleTimeConversion = () => {
  if (timeValue.value === null || timeValue.value === undefined) {
    timeResult.value = null;
    return;
  }
  
  try {
    timeResult.value = convertTime(timeValue.value, timeFromUnit.value, timeToUnit.value);
  } catch (error) {
    console.error('时间转换错误:', error);
    timeResult.value = null;
  }
};

// 交换时间单位
const swapTimeUnits = () => {
  const temp = timeFromUnit.value;
  timeFromUnit.value = timeToUnit.value;
  timeToUnit.value = temp;
  handleTimeConversion();
};

// 设置目标时间单位
const setTimeToUnit = (unit: TimeUnit) => {
  timeToUnit.value = unit;
  handleTimeConversion();
};

// ==================== 时间差值计算相关 ====================
const startDate = ref<Dayjs | null>(dayjs().subtract(1, 'day'));
const startTime = ref<Dayjs | null>(dayjs().set('hour', 9).set('minute', 0).set('second', 0));
const endDate = ref<Dayjs | null>(dayjs());
const endTime = ref<Dayjs | null>(dayjs());
const startDateTimeString = ref<string>('');
const endDateTimeString = ref<string>('');
const timeDifferenceResult = ref<TimeDifferenceResult | null>(null);

// 获取完整的开始时间
const getStartDateTime = (): Date => {
  if (startDate.value && startTime.value) {
    return startDate.value
      .set('hour', startTime.value.hour())
      .set('minute', startTime.value.minute())
      .set('second', startTime.value.second())
      .toDate();
  } else if (startDate.value) {
    return startDate.value.toDate();
  } else if (startTime.value) {
    const now = new Date();
    return startTime.value
      .set('year', now.getFullYear())
      .set('month', now.getMonth())
      .set('date', now.getDate())
      .toDate();
  } else {
    return new Date();
  }
};

// 获取完整的结束时间
const getEndDateTime = (): Date => {
  if (endDate.value && endTime.value) {
    return endDate.value
      .set('hour', endTime.value.hour())
      .set('minute', endTime.value.minute())
      .set('second', endTime.value.second())
      .toDate();
  } else if (endDate.value) {
    return endDate.value.toDate();
  } else if (endTime.value) {
    const now = new Date();
    return endTime.value
      .set('year', now.getFullYear())
      .set('month', now.getMonth())
      .set('date', now.getDate())
      .toDate();
  } else {
    return new Date();
  }
};

// 计算时间差值
const calculateTimeDifference = () => {
  try {
    const start = getStartDateTime();
    const end = getEndDateTime();
    timeDifferenceResult.value = calcTimeDiff(start, end);
    
    // 更新字符串显示
    startDateTimeString.value = formatDateTime(start);
    endDateTimeString.value = formatDateTime(end);
  } catch (error) {
    console.error('时间差值计算错误:', error);
    timeDifferenceResult.value = null;
  }
};

// 解析开始日期时间字符串
const parseStartDateTime = () => {
  try {
    const date = parseTimeString(startDateTimeString.value);
    startDate.value = dayjs(date);
    startTime.value = dayjs(date);
    calculateTimeDifference();
  } catch (error) {
    console.error('解析开始时间错误:', error);
  }
};

// 解析结束日期时间字符串
const parseEndDateTime = () => {
  try {
    const date = parseTimeString(endDateTimeString.value);
    endDate.value = dayjs(date);
    endTime.value = dayjs(date);
    calculateTimeDifference();
  } catch (error) {
    console.error('解析结束时间错误:', error);
  }
};

// 设置开始时间为现在
const setStartToNow = () => {
  const now = dayjs();
  startDate.value = now;
  startTime.value = now;
  calculateTimeDifference();
};

// 设置结束时间为现在
const setEndToNow = () => {
  const now = dayjs();
  endDate.value = now;
  endTime.value = now;
  calculateTimeDifference();
};

// 交换开始和结束时间
const swapTimes = () => {
  const tempDate = startDate.value;
  const tempTime = startTime.value;
  const tempString = startDateTimeString.value;
  
  startDate.value = endDate.value;
  startTime.value = endTime.value;
  startDateTimeString.value = endDateTimeString.value;
  
  endDate.value = tempDate;
  endTime.value = tempTime;
  endDateTimeString.value = tempString;
  
  calculateTimeDifference();
};

// 组件挂载时初始化计算
onMounted(() => {
  handleByteConversion();
  handleTimeConversion();
  calculateTimeDifference();
});
</script>

<style scoped lang="scss">
.unit-converter {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  
  .converter-tabs {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: hidden;
    
    :deep(.ant-tabs-nav) {
      margin-bottom: 16px;
      flex-shrink: 0;
    }
    
    :deep(.ant-tabs-content) {
      flex: 1;
      overflow: hidden;
      
      .ant-tabs-tabpane {
        height: 100%;
        overflow: hidden;
      }
    }
  }
  
  .converter-section {
    height: 100%;
    overflow: hidden;
    display: flex;
    flex-direction: column;
    
    .ant-card {
      flex: 1;
      display: flex;
      flex-direction: column;
      overflow: hidden;
      
      .ant-card-body {
        flex: 1;
        overflow: auto;
        padding: 16px;
        
        /* 隐藏滚动条但保持滚动功能 */
        &::-webkit-scrollbar {
          width: 6px;
          height: 6px;
        }
        
        &::-webkit-scrollbar-track {
          background: #f1f1f1;
          border-radius: 3px;
        }
        
        &::-webkit-scrollbar-thumb {
          background: #c1c1c1;
          border-radius: 3px;
          
          &:hover {
            background: #a8a8a8;
          }
        }
      }
    }
    
    .converter-form {
      .converter-arrow {
        font-size: 20px;
        color: #1890ff;
        font-weight: bold;
      }
      
      .converter-result {
        min-height: 60px;
      }
      
      .quick-conversions {
        h4 {
          margin-bottom: 12px;
          color: #666;
        }
      }
    }
    
    .time-difference-form {
      .quick-actions {
        display: flex;
        justify-content: center;
        margin: 16px 0;
      }
      
      .difference-results {
        .result-text {
          font-size: 18px;
          font-weight: bold;
          color: #1890ff;
        }
        
        .additional-info {
          padding: 16px;
          background: #fafafa;
          border-radius: 4px;
        }
      }
    }
  }
}
</style>