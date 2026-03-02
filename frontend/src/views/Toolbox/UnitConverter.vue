<template>
  <div class="p-6 unit-converter">
    <a-card title="单位转换工具" class="mb-6">
      <a-tabs v-model:activeKey="activeTab">
        <!-- 字节转换 -->
        <a-tab-pane key="bytes" tab="字节转换">
          <a-form layout="vertical">
            <a-row :gutter="16">
              <a-col :span="12">
                <a-form-item label="输入值">
                  <a-input-number
                    v-model:value="byteInput.value"
                    :min="0"
                    style="width: 100%"
                    placeholder="请输入数值"
                  />
                </a-form-item>
              </a-col>
              <a-col :span="12">
                <a-form-item label="单位">
                  <a-select v-model:value="byteInput.unit" style="width: 100%">
                    <a-select-option value="B">B (字节)</a-select-option>
                    <a-select-option value="KB">KB (千字节)</a-select-option>
                    <a-select-option value="MB">MB (兆字节)</a-select-option>
                    <a-select-option value="GB">GB (吉字节)</a-select-option>
                    <a-select-option value="TB">TB (太字节)</a-select-option>
                  </a-select>
                </a-form-item>
              </a-col>
            </a-row>
            <a-form-item>
              <a-button type="primary" @click="convertBytes" :loading="loading">
                转换
              </a-button>
            </a-form-item>
          </a-form>

          <div v-if="byteResult" class="result-section">
            <a-divider>转换结果</a-divider>
            <a-row :gutter="[16, 16]">
              <a-col :span="12" v-for="(value, key) in byteResult" :key="key">
                <a-statistic :title="getByteLabel(key)" :value="value" :precision="2" />
              </a-col>
            </a-row>
          </div>
        </a-tab-pane>

        <!-- 长度转换 -->
        <a-tab-pane key="length" tab="长度转换">
          <a-form layout="vertical">
            <a-row :gutter="16">
              <a-col :span="12">
                <a-form-item label="输入值">
                  <a-input-number
                    v-model:value="lengthInput.value"
                    :min="0"
                    style="width: 100%"
                    placeholder="请输入数值"
                  />
                </a-form-item>
              </a-col>
              <a-col :span="12">
                <a-form-item label="单位">
                  <a-select v-model:value="lengthInput.unit" style="width: 100%">
                    <a-select-option value="mm">mm (毫米)</a-select-option>
                    <a-select-option value="cm">cm (厘米)</a-select-option>
                    <a-select-option value="m">m (米)</a-select-option>
                    <a-select-option value="km">km (千米)</a-select-option>
                    <a-select-option value="in">in (英寸)</a-select-option>
                    <a-select-option value="ft">ft (英尺)</a-select-option>
                    <a-select-option value="yd">yd (码)</a-select-option>
                    <a-select-option value="mi">mi (英里)</a-select-option>
                  </a-select>
                </a-form-item>
              </a-col>
            </a-row>
            <a-form-item>
              <a-button type="primary" @click="convertLength" :loading="loading">
                转换
              </a-button>
            </a-form-item>
          </a-form>

          <div v-if="lengthResult" class="result-section">
            <a-divider>转换结果</a-divider>
            <a-row :gutter="[16, 16]">
              <a-col :span="12" v-for="(value, key) in lengthResult" :key="key">
                <a-statistic :title="getLengthLabel(key)" :value="value" :precision="4" />
              </a-col>
            </a-row>
          </div>
        </a-tab-pane>

        <!-- 时间转换 -->
        <a-tab-pane key="time" tab="时间转换">
          <a-form layout="vertical">
            <a-row :gutter="16">
              <a-col :span="12">
                <a-form-item label="输入值">
                  <a-input-number
                    v-model:value="timeInput.value"
                    :min="0"
                    style="width: 100%"
                    placeholder="请输入数值"
                  />
                </a-form-item>
              </a-col>
              <a-col :span="12">
                <a-form-item label="单位">
                  <a-select v-model:value="timeInput.unit" style="width: 100%">
                    <a-select-option value="ms">ms (毫秒)</a-select-option>
                    <a-select-option value="s">s (秒)</a-select-option>
                    <a-select-option value="min">min (分钟)</a-select-option>
                    <a-select-option value="h">h (小时)</a-select-option>
                    <a-select-option value="d">d (天)</a-select-option>
                    <a-select-option value="w">w (周)</a-select-option>
                    <a-select-option value="mon">mon (月)</a-select-option>
                    <a-select-option value="y">y (年)</a-select-option>
                  </a-select>
                </a-form-item>
              </a-col>
            </a-row>
            <a-form-item>
              <a-button type="primary" @click="convertTime" :loading="loading">
                转换
              </a-button>
            </a-form-item>
          </a-form>

          <div v-if="timeResult" class="result-section">
            <a-divider>转换结果</a-divider>
            <a-row :gutter="[16, 16]">
              <a-col :span="12" v-for="(value, key) in timeResult" :key="key">
                <a-statistic :title="getTimeLabel(key)" :value="value" :precision="4" />
              </a-col>
            </a-row>
          </div>
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { message } from 'ant-design-vue'
import * as UnitConverterService from '../../../bindings/github.com/Aliuyanfeng/happytools/backend/services/unitconverter/unitconverterservice'

const activeTab = ref('bytes')
const loading = ref(false)

// 字节转换
const byteInput = ref({ value: 0, unit: 'B' })
const byteResult = ref<any>(null)

// 长度转换
const lengthInput = ref({ value: 0, unit: 'm' })
const lengthResult = ref<any>(null)

// 时间转换
const timeInput = ref({ value: 0, unit: 's' })
const timeResult = ref<any>(null)

// 字节转换
const convertBytes = async () => {
  if (byteInput.value.value === 0) {
    message.warning('请输入数值')
    return
  }

  loading.value = true
  try {
    const result = await UnitConverterService.ConvertBytes(byteInput.value.value, byteInput.value.unit)
    byteResult.value = result
  } catch (error: any) {
    message.error(error.message || '转换失败')
  } finally {
    loading.value = false
  }
}

// 长度转换
const convertLength = async () => {
  if (lengthInput.value.value === 0) {
    message.warning('请输入数值')
    return
  }

  loading.value = true
  try {
    const result = await UnitConverterService.ConvertLength(lengthInput.value.value, lengthInput.value.unit)
    lengthResult.value = result
  } catch (error: any) {
    message.error(error.message || '转换失败')
  } finally {
    loading.value = false
  }
}

// 时间转换
const convertTime = async () => {
  if (timeInput.value.value === 0) {
    message.warning('请输入数值')
    return
  }

  loading.value = true
  try {
    const result = await UnitConverterService.ConvertTime(timeInput.value.value, timeInput.value.unit)
    timeResult.value = result
  } catch (error: any) {
    message.error(error.message || '转换失败')
  } finally {
    loading.value = false
  }
}

// 获取字节单位标签
const getByteLabel = (key: string | number) => {
  const keyStr = String(key)
  const labels: Record<string, string> = {
    bytes: 'B (字节)',
    kilobytes: 'KB (千字节)',
    megabytes: 'MB (兆字节)',
    gigabytes: 'GB (吉字节)',
    terabytes: 'TB (太字节)'
  }
  return labels[keyStr] || keyStr
}

// 获取长度单位标签
const getLengthLabel = (key: string | number) => {
  const keyStr = String(key)
  const labels: Record<string, string> = {
    millimeters: 'mm (毫米)',
    centimeters: 'cm (厘米)',
    meters: 'm (米)',
    kilometers: 'km (千米)',
    inches: 'in (英寸)',
    feet: 'ft (英尺)',
    yards: 'yd (码)',
    miles: 'mi (英里)'
  }
  return labels[keyStr] || keyStr
}

// 获取时间单位标签
const getTimeLabel = (key: string | number) => {
  const keyStr = String(key)
  const labels: Record<string, string> = {
    milliseconds: 'ms (毫秒)',
    seconds: 's (秒)',
    minutes: 'min (分钟)',
    hours: 'h (小时)',
    days: 'd (天)',
    weeks: 'w (周)',
    months: 'mon (月)',
    years: 'y (年)'
  }
  return labels[keyStr] || keyStr
}
</script>

<style scoped>
.unit-converter {
  max-width: 1200px;
  margin: 0 auto;
}

.result-section {
  margin-top: 24px;
  padding: 16px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 8px;
}
</style>
