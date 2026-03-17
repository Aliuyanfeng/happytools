<template>
  <div class="p-6 unit-converter">
    <a-card :title="t('toolbox.unitConverter.title')" class="mb-6">
      <a-tabs v-model:activeKey="activeTab">
        <!-- 字节转换 -->
        <a-tab-pane key="bytes" :tab="t('toolbox.unitConverter.tabBytes')">
          <a-form layout="vertical">
            <a-row :gutter="16">
              <a-col :span="12">
                <a-form-item :label="t('toolbox.unitConverter.inputValue')">
                  <a-input-number v-model:value="byteInput.value" :min="0" style="width: 100%" :placeholder="t('toolbox.unitConverter.inputPlaceholder')" />
                </a-form-item>
              </a-col>
              <a-col :span="12">
                <a-form-item :label="t('toolbox.unitConverter.unit')">
                  <a-select v-model:value="byteInput.unit" style="width: 100%">
                    <a-select-option value="B">{{ t('toolbox.unitConverter.byteB') }}</a-select-option>
                    <a-select-option value="KB">{{ t('toolbox.unitConverter.byteKB') }}</a-select-option>
                    <a-select-option value="MB">{{ t('toolbox.unitConverter.byteMB') }}</a-select-option>
                    <a-select-option value="GB">{{ t('toolbox.unitConverter.byteGB') }}</a-select-option>
                    <a-select-option value="TB">{{ t('toolbox.unitConverter.byteTB') }}</a-select-option>
                  </a-select>
                </a-form-item>
              </a-col>
            </a-row>
            <a-form-item>
              <a-button type="primary" @click="convertBytes" :loading="loading">{{ t('toolbox.unitConverter.convert') }}</a-button>
            </a-form-item>
          </a-form>
          <div v-if="byteResult" class="result-section">
            <a-divider>{{ t('toolbox.unitConverter.resultTitle') }}</a-divider>
            <a-row :gutter="[16, 16]">
              <a-col :span="12" v-for="(value, key) in byteResult" :key="key">
                <a-statistic :title="getByteLabel(key)" :value="value" :precision="2" />
              </a-col>
            </a-row>
          </div>
        </a-tab-pane>

        <!-- 长度转换 -->
        <a-tab-pane key="length" :tab="t('toolbox.unitConverter.tabLength')">
          <a-form layout="vertical">
            <a-row :gutter="16">
              <a-col :span="12">
                <a-form-item :label="t('toolbox.unitConverter.inputValue')">
                  <a-input-number v-model:value="lengthInput.value" :min="0" style="width: 100%" :placeholder="t('toolbox.unitConverter.inputPlaceholder')" />
                </a-form-item>
              </a-col>
              <a-col :span="12">
                <a-form-item :label="t('toolbox.unitConverter.unit')">
                  <a-select v-model:value="lengthInput.unit" style="width: 100%">
                    <a-select-option value="mm">{{ t('toolbox.unitConverter.lengthMm') }}</a-select-option>
                    <a-select-option value="cm">{{ t('toolbox.unitConverter.lengthCm') }}</a-select-option>
                    <a-select-option value="m">{{ t('toolbox.unitConverter.lengthM') }}</a-select-option>
                    <a-select-option value="km">{{ t('toolbox.unitConverter.lengthKm') }}</a-select-option>
                    <a-select-option value="in">{{ t('toolbox.unitConverter.lengthIn') }}</a-select-option>
                    <a-select-option value="ft">{{ t('toolbox.unitConverter.lengthFt') }}</a-select-option>
                    <a-select-option value="yd">{{ t('toolbox.unitConverter.lengthYd') }}</a-select-option>
                    <a-select-option value="mi">{{ t('toolbox.unitConverter.lengthMi') }}</a-select-option>
                  </a-select>
                </a-form-item>
              </a-col>
            </a-row>
            <a-form-item>
              <a-button type="primary" @click="convertLength" :loading="loading">{{ t('toolbox.unitConverter.convert') }}</a-button>
            </a-form-item>
          </a-form>
          <div v-if="lengthResult" class="result-section">
            <a-divider>{{ t('toolbox.unitConverter.resultTitle') }}</a-divider>
            <a-row :gutter="[16, 16]">
              <a-col :span="12" v-for="(value, key) in lengthResult" :key="key">
                <a-statistic :title="getLengthLabel(key)" :value="value" :precision="4" />
              </a-col>
            </a-row>
          </div>
        </a-tab-pane>

        <!-- 时间转换 -->
        <a-tab-pane key="time" :tab="t('toolbox.unitConverter.tabTime')">
          <a-form layout="vertical">
            <a-row :gutter="16">
              <a-col :span="12">
                <a-form-item :label="t('toolbox.unitConverter.inputValue')">
                  <a-input-number v-model:value="timeInput.value" :min="0" style="width: 100%" :placeholder="t('toolbox.unitConverter.inputPlaceholder')" />
                </a-form-item>
              </a-col>
              <a-col :span="12">
                <a-form-item :label="t('toolbox.unitConverter.unit')">
                  <a-select v-model:value="timeInput.unit" style="width: 100%">
                    <a-select-option value="ms">{{ t('toolbox.unitConverter.timeMs') }}</a-select-option>
                    <a-select-option value="s">{{ t('toolbox.unitConverter.timeS') }}</a-select-option>
                    <a-select-option value="min">{{ t('toolbox.unitConverter.timeMin') }}</a-select-option>
                    <a-select-option value="h">{{ t('toolbox.unitConverter.timeH') }}</a-select-option>
                    <a-select-option value="d">{{ t('toolbox.unitConverter.timeD') }}</a-select-option>
                    <a-select-option value="w">{{ t('toolbox.unitConverter.timeW') }}</a-select-option>
                    <a-select-option value="mon">{{ t('toolbox.unitConverter.timeMon') }}</a-select-option>
                    <a-select-option value="y">{{ t('toolbox.unitConverter.timeY') }}</a-select-option>
                  </a-select>
                </a-form-item>
              </a-col>
            </a-row>
            <a-form-item>
              <a-button type="primary" @click="convertTime" :loading="loading">{{ t('toolbox.unitConverter.convert') }}</a-button>
            </a-form-item>
          </a-form>
          <div v-if="timeResult" class="result-section">
            <a-divider>{{ t('toolbox.unitConverter.resultTitle') }}</a-divider>
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
import { useI18n } from 'vue-i18n'
import * as UnitConverterService from '../../../bindings/github.com/Aliuyanfeng/happytools/backend/services/unitconverter/unitconverterservice'

const { t } = useI18n()

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

const convertBytes = async () => {
  if (byteInput.value.value === 0) { message.warning(t('toolbox.unitConverter.pleaseInput')); return }
  loading.value = true
  try {
    byteResult.value = await UnitConverterService.ConvertBytes(byteInput.value.value, byteInput.value.unit)
  } catch (error: any) { message.error(error.message || t('toolbox.unitConverter.convertFailed')) }
  finally { loading.value = false }
}

const convertLength = async () => {
  if (lengthInput.value.value === 0) { message.warning(t('toolbox.unitConverter.pleaseInput')); return }
  loading.value = true
  try {
    lengthResult.value = await UnitConverterService.ConvertLength(lengthInput.value.value, lengthInput.value.unit)
  } catch (error: any) { message.error(error.message || t('toolbox.unitConverter.convertFailed')) }
  finally { loading.value = false }
}

const convertTime = async () => {
  if (timeInput.value.value === 0) { message.warning(t('toolbox.unitConverter.pleaseInput')); return }
  loading.value = true
  try {
    timeResult.value = await UnitConverterService.ConvertTime(timeInput.value.value, timeInput.value.unit)
  } catch (error: any) { message.error(error.message || t('toolbox.unitConverter.convertFailed')) }
  finally { loading.value = false }
}

const getByteLabel = (key: string | number) => {
  const keyStr = String(key)
  const map: Record<string, string> = {
    bytes: t('toolbox.unitConverter.byteB'),
    kilobytes: t('toolbox.unitConverter.byteKB'),
    megabytes: t('toolbox.unitConverter.byteMB'),
    gigabytes: t('toolbox.unitConverter.byteGB'),
    terabytes: t('toolbox.unitConverter.byteTB'),
  }
  return map[keyStr] || keyStr
}

const getLengthLabel = (key: string | number) => {
  const keyStr = String(key)
  const map: Record<string, string> = {
    millimeters: t('toolbox.unitConverter.lengthMm'),
    centimeters: t('toolbox.unitConverter.lengthCm'),
    meters: t('toolbox.unitConverter.lengthM'),
    kilometers: t('toolbox.unitConverter.lengthKm'),
    inches: t('toolbox.unitConverter.lengthIn'),
    feet: t('toolbox.unitConverter.lengthFt'),
    yards: t('toolbox.unitConverter.lengthYd'),
    miles: t('toolbox.unitConverter.lengthMi'),
  }
  return map[keyStr] || keyStr
}

const getTimeLabel = (key: string | number) => {
  const keyStr = String(key)
  const map: Record<string, string> = {
    milliseconds: t('toolbox.unitConverter.timeMs'),
    seconds: t('toolbox.unitConverter.timeS'),
    minutes: t('toolbox.unitConverter.timeMin'),
    hours: t('toolbox.unitConverter.timeH'),
    days: t('toolbox.unitConverter.timeD'),
    weeks: t('toolbox.unitConverter.timeW'),
    months: t('toolbox.unitConverter.timeMon'),
    years: t('toolbox.unitConverter.timeY'),
  }
  return map[keyStr] || keyStr
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
