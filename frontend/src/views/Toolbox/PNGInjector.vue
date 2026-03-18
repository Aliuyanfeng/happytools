<template>
  <div class="p-6">
    <a-card :title="t('toolbox.pngInjector.title')">
      <a-steps :current="step" size="small" class="mb-6">
        <a-step :title="t('toolbox.pngInjector.stepSelectFile')" />
        <a-step :title="t('toolbox.pngInjector.stepViewChunks')" />
        <a-step :title="t('toolbox.pngInjector.stepConfig')" />
        <a-step :title="t('toolbox.pngInjector.stepDone')" />
      </a-steps>

      <!-- Step 0: 选择文件 -->
      <div v-if="step === 0" class="step-body">
        <div
          id="png-injector-drop"
          data-file-drop-target
          class="select-area"
          @click="selectFile"
        >
          <FileImageOutlined class="select-icon" />
          <p class="mt-3 text-base text-gray-600">{{ t('toolbox.pngInjector.dropHint') }}</p>
          <p class="text-xs text-gray-400">{{ t('toolbox.pngInjector.dropHintSub') }}</p>
        </div>
      </div>

      <!-- Step 1: Chunk 列表 -->
      <div v-if="step === 1" class="step-body">
        <div class="mb-3 flex items-center justify-between">
          <span class="text-sm text-gray-500">
            {{ t('toolbox.pngInjector.fileInfo') }}<span class="font-medium text-gray-700">{{ fileName }}</span>
            &nbsp;·&nbsp;{{ t('toolbox.pngInjector.chunkCount', { count: chunks.length }) }}
          </span>
          <a-button size="small" @click="step = 0">{{ t('toolbox.pngInjector.reselect') }}</a-button>
        </div>

        <a-alert
          class="mb-3"
          :closable="false"
          show-icon
          type="info"
        >
          <template #message>
            {{ t('toolbox.pngInjector.positionHint') }}
            <span class="text-gray-500">{{ t('toolbox.pngInjector.positionHintGray') }}</span>
          </template>
        </a-alert>

        <a-table
          :columns="chunkCols"
          :data-source="chunks"
          :pagination="false"
          :scroll="{ y: 380 }"
          row-key="index"
          size="small"
          bordered
        >
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'type'">
              <a-tag :color="record.critical ? 'red' : 'blue'" class="font-mono">
                {{ record.type }}
              </a-tag>
            </template>
            <template v-if="column.key === 'length'">
              {{ record.length === 0 ? '0 B' : formatBytes(record.length) }}
            </template>
            <template v-if="column.key === 'crc'">
              <span class="font-mono text-xs">{{ record.crc }}</span>
            </template>
            <template v-if="column.key === 'dataHex'">
              <span class="font-mono text-xs text-gray-500">{{ record.dataHex || t('toolbox.pngInjector.empty') }}</span>
            </template>
            <template v-if="column.key === 'action'">
              <a-space>
                <a-tooltip :title="getDisabledReason('before', record)">
                  <a-button
                    type="link"
                    size="small"
                    :disabled="!isLegalPosition('before', record)"
                    @click="pickPosition('before', record)"
                  >
                    {{ t('toolbox.pngInjector.insertBefore') }}
                  </a-button>
                </a-tooltip>
                <a-tooltip :title="getDisabledReason('after', record)">
                  <a-button
                    type="link"
                    size="small"
                    :disabled="!isLegalPosition('after', record)"
                    @click="pickPosition('after', record)"
                  >
                    {{ t('toolbox.pngInjector.insertAfter') }}
                  </a-button>
                </a-tooltip>
              </a-space>
            </template>
          </template>
        </a-table>
      </div>

      <!-- Step 2: 配置注入参数 -->
      <div v-if="step === 2" class="step-body">
        <a-alert class="mb-4" type="success" show-icon :closable="false">
          <template #message>
            {{ t('toolbox.pngInjector.injectPositionLabel') }}
            <a-tag :color="targetChunk?.critical ? 'red' : 'blue'" class="font-mono mx-1">
              {{ targetChunk?.type }}
            </a-tag>
            （索引 {{ targetChunk?.index }}）
            <strong>{{ position === 'before' ? t('toolbox.pngInjector.before') : t('toolbox.pngInjector.after') }}</strong>
          </template>
        </a-alert>

        <a-form :model="form" layout="vertical">
          <a-form-item :label="t('toolbox.pngInjector.chunkTypeLabel')" required>
            <a-input
              v-model:value="form.chunkType"
              maxlength="4"
              show-count
              :placeholder="t('toolbox.pngInjector.chunkTypePlaceholder')"
              style="width: 200px"
              :status="form.chunkType.length > 0 && form.chunkType.length !== 4 ? 'error' : ''"
            />
            <div class="text-xs text-gray-400 mt-1">
              {{ t('toolbox.pngInjector.chunkTypeHint') }}
            </div>
          </a-form-item>

          <a-form-item :label="t('toolbox.pngInjector.payloadLabel')" required>
            <a-textarea
              v-model:value="form.payload"
              :rows="8"
              :placeholder="t('toolbox.pngInjector.payloadPlaceholder')"
              class="font-mono"
            />
            <div class="text-xs text-gray-400 mt-1">
              {{ t('toolbox.pngInjector.payloadBytes', { count: form.payload.length }) }}
            </div>
          </a-form-item>

          <a-form-item :label="t('toolbox.pngInjector.outputPathLabel')" required>
            <a-input-group compact>
              <a-input
                v-model:value="form.outputPath"
                readonly
                :placeholder="t('toolbox.pngInjector.outputPathPlaceholder')"
                style="width: calc(100% - 100px)"
              />
              <a-button @click="selectOutput">{{ t('toolbox.pngInjector.selectOutput') }}</a-button>
            </a-input-group>
          </a-form-item>

          <a-form-item>
            <a-space>
              <a-button type="primary" :loading="loading" @click="doInject">
                {{ t('toolbox.pngInjector.inject') }}
              </a-button>
              <a-button @click="step = 1">{{ t('toolbox.pngInjector.back') }}</a-button>
            </a-space>
          </a-form-item>
        </a-form>
      </div>

      <!-- Step 3: 完成 -->
      <div v-if="step === 3" class="step-body">
        <a-result status="success" :title="t('toolbox.pngInjector.successTitle')">
          <template #subTitle>
            {{ t('toolbox.pngInjector.successSub') }}<span class="font-mono text-sm">{{ form.outputPath }}</span>
          </template>
          <template #extra>
            <a-button type="primary" @click="reset">{{ t('toolbox.pngInjector.injectNew') }}</a-button>
          </template>
        </a-result>

        <a-descriptions bordered size="small" class="mt-4" :column="2">
          <a-descriptions-item :label="t('toolbox.pngInjector.descOriginalFile')">{{ fileName }}</a-descriptions-item>
          <a-descriptions-item :label="t('toolbox.pngInjector.descChunkType')">
            <span class="font-mono">{{ form.chunkType }}</span>
          </a-descriptions-item>
          <a-descriptions-item :label="t('toolbox.pngInjector.descInjectPosition')">
            {{ targetChunk?.type }}（索引 {{ targetChunk?.index }}）{{ position === 'before' ? ' ' + t('toolbox.pngInjector.before') : ' ' + t('toolbox.pngInjector.after') }}
          </a-descriptions-item>
          <a-descriptions-item :label="t('toolbox.pngInjector.descPayloadSize')">
            {{ formatBytes(form.payload.length) }}
          </a-descriptions-item>
          <a-descriptions-item :label="t('toolbox.pngInjector.descOutputPath')" :span="2">
            <span class="font-mono text-xs">{{ form.outputPath }}</span>
          </a-descriptions-item>
        </a-descriptions>
      </div>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, onUnmounted } from 'vue'
import { message } from 'ant-design-vue'
import { FileImageOutlined } from '@ant-design/icons-vue'
import { useI18n } from 'vue-i18n'
import { Events } from '@wailsio/runtime'
import * as PNGInjectorService from '../../../bindings/github.com/Aliuyanfeng/happytools/backend/services/pnginjector/pnginjectorservice'

const { t } = useI18n()

interface PNGChunk {
  index: number
  type: string
  length: number
  dataHex: string
  crc: string
  offset: number
  critical: boolean
}

const step = ref(0)
const loading = ref(false)
const filePath = ref('')
const fileName = ref('')
const chunks = ref<PNGChunk[]>([])
const targetChunk = ref<PNGChunk | null>(null)
const position = ref<'before' | 'after'>('after')

const form = reactive({
  chunkType: 'rPOC',
  payload: '',
  outputPath: ''
})

// 处理拖拽进来的文件路径
const handleDroppedFile = async (path: string) => {
  if (!path.toLowerCase().endsWith('.png')) {
    message.warning(t('toolbox.pngInjector.onlyPng'))
    return
  }
  filePath.value = path
  fileName.value = path.replace(/\\/g, '/').split('/').pop() || path
  loading.value = true
  try {
    const result = await PNGInjectorService.ParsePNG(path)
    chunks.value = result
    step.value = 1
    message.success(t('toolbox.pngInjector.parseDone', { count: result.length }))
  } catch (e: any) {
    message.error(t('toolbox.pngInjector.parseFailed') + '：' + (e?.message || e))
  } finally {
    loading.value = false
  }
}

let offFileDrop: (() => void) | null = null

onMounted(() => {
  offFileDrop = Events.On('wails:file-drop', (event: any) => {
    const data = event?.data
    const droppedFiles: string[] = Array.isArray(data?.files) ? data.files : []
    const targetId: string = data?.target ?? ''

    // 只在 step 0 且拖到本区域（或 body）时处理
    if (step.value !== 0) return
    if (targetId && targetId !== 'png-injector-drop') return
    if (droppedFiles.length === 0) return

    handleDroppedFile(droppedFiles[0])
  })
})

onUnmounted(() => {
  if (offFileDrop) {
    offFileDrop()
    offFileDrop = null
  }
})

const chunkCols = computed(() => [
  { title: t('toolbox.pngInjector.colIndex'), dataIndex: 'index', key: 'index', width: 60 },
  { title: t('toolbox.pngInjector.colType'), dataIndex: 'type', key: 'type', width: 130 },
  { title: t('toolbox.pngInjector.colLength'), dataIndex: 'length', key: 'length', width: 100 },
  { title: t('toolbox.pngInjector.colCrc'), dataIndex: 'crc', key: 'crc', width: 110 },
  { title: t('toolbox.pngInjector.colDataHex'), dataIndex: 'dataHex', key: 'dataHex', ellipsis: true },
  { title: t('toolbox.pngInjector.colAction'), key: 'action', width: 160 }
])

// 选择 PNG 文件
const selectFile = async () => {
  try {
    const path = await PNGInjectorService.OpenFileDialog()
    if (!path) return
    await handleDroppedFile(path)
  } catch (e: any) {
    message.error(t('toolbox.pngInjector.parseFailed') + '：' + (e?.message || e))
  }
}

// 判断某个注入位置是否合法
// 规则：
//   非法：IHDR 之前
//   非法：IEND 之后
//   非法：两个相邻 IDAT 之间（即 IDAT 之后且下一个也是 IDAT，或 IDAT 之前且上一个也是 IDAT）
const isLegalPosition = (pos: 'before' | 'after', chunk: PNGChunk): boolean => {
  const idx = chunk.index
  const prev = chunks.value[idx - 1]
  const next = chunks.value[idx + 1]

  if (pos === 'before') {
    // IHDR 之前 → 非法
    if (chunk.type === 'IHDR') return false
    // 上一个是 IDAT 且当前也是 IDAT → 非法（IDAT 块之间）
    if (chunk.type === 'IDAT' && prev?.type === 'IDAT') return false
    return true
  } else {
    // IEND 之后 → 非法
    if (chunk.type === 'IEND') return false
    // 当前是 IDAT 且下一个也是 IDAT → 非法（IDAT 块之间）
    if (chunk.type === 'IDAT' && next?.type === 'IDAT') return false
    return true
  }
}

// 返回禁用原因（用于 tooltip）
const getDisabledReason = (pos: 'before' | 'after', chunk: PNGChunk): string => {
  if (isLegalPosition(pos, chunk)) return ''
  const idx = chunk.index
  const prev = chunks.value[idx - 1]
  const next = chunks.value[idx + 1]

  if (pos === 'before' && chunk.type === 'IHDR') return t('toolbox.pngInjector.disabledIHDRBefore')
  if (pos === 'after' && chunk.type === 'IEND') return t('toolbox.pngInjector.disabledIENDAfter')
  if (pos === 'before' && chunk.type === 'IDAT' && prev?.type === 'IDAT') return t('toolbox.pngInjector.disabledIDATBetween')
  if (pos === 'after' && chunk.type === 'IDAT' && next?.type === 'IDAT') return t('toolbox.pngInjector.disabledIDATBetween')
  return t('toolbox.pngInjector.disabledIllegal')
}
const pickPosition = (pos: 'before' | 'after', chunk: PNGChunk) => {
  position.value = pos
  targetChunk.value = chunk

  // 生成默认输出路径
  const base = filePath.value.replace(/\\/g, '/')
  const dir = base.substring(0, base.lastIndexOf('/'))
  const name = (fileName.value || 'output').replace(/\.png$/i, '')
  form.outputPath = `${dir}/${name}_injected.png`

  step.value = 2
}

// 选择输出路径
const selectOutput = async () => {
  try {
    const defaultName = (fileName.value || 'output').replace(/\.png$/i, '') + '_injected.png'
    const path = await PNGInjectorService.SaveFileDialog(defaultName)
    if (path) form.outputPath = path
  } catch (e: any) {
    message.error(t('toolbox.pngInjector.selectOutputFailed') + '：' + (e?.message || e))
  }
}

// 执行注入
const doInject = async () => {
  if (form.chunkType.length !== 4) {
    message.warning(t('toolbox.pngInjector.chunkTypeMustBe4'))
    return
  }
  if (!form.payload) {
    message.warning(t('toolbox.pngInjector.payloadRequired'))
    return
  }
  if (!form.outputPath) {
    message.warning(t('toolbox.pngInjector.outputRequired'))
    return
  }

  loading.value = true
  try {
    await PNGInjectorService.InjectChunk(
      filePath.value,
      form.outputPath,
      form.chunkType,
      form.payload,
      position.value,
      targetChunk.value!.index
    )
    step.value = 3
  } catch (e: any) {
    message.error(t('toolbox.pngInjector.injectFailed') + '：' + (e?.message || e))
  } finally {
    loading.value = false
  }
}

// 重置
const reset = () => {
  step.value = 0
  filePath.value = ''
  fileName.value = ''
  chunks.value = []
  targetChunk.value = null
  position.value = 'after'
  form.chunkType = 'rPOC'
  form.payload = ''
  form.outputPath = ''
}

const formatBytes = (n: number) => {
  if (n === 0) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(n) / Math.log(1024))
  return (n / Math.pow(1024, i)).toFixed(i === 0 ? 0 : 1) + ' ' + units[i]
}
</script>

<style scoped>
.step-body {
  min-height: 420px;
  padding: 8px 0;
}

.select-area {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 280px;
  border: 2px dashed #d9d9d9;
  border-radius: 8px;
  cursor: pointer;
  transition: border-color 0.2s, background 0.2s;
}

.select-area:hover,
.select-area.file-drop-target-active {
  border-color: #1890ff;
  background: #f0f8ff;
}

.select-icon {
  font-size: 56px;
  color: #bfbfbf;
}

.select-area:hover .select-icon,
.select-area.file-drop-target-active .select-icon {
  color: #1890ff;
}

.font-mono {
  font-family: 'Courier New', Courier, monospace;
}
</style>
