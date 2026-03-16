<template>
  <div class="p-6">
    <a-card title="PNG Chunk 注入工具">
      <a-steps :current="step" size="small" class="mb-6">
        <a-step title="选择文件" />
        <a-step title="查看 Chunks" />
        <a-step title="配置注入" />
        <a-step title="完成" />
      </a-steps>

      <!-- Step 0: 选择文件 -->
      <div v-if="step === 0" class="step-body">
        <div class="select-area" @click="selectFile">
          <FileImageOutlined class="select-icon" />
          <p class="mt-3 text-base text-gray-600">点击选择 PNG 文件</p>
          <p class="text-xs text-gray-400">仅支持 .png 格式</p>
        </div>
      </div>

      <!-- Step 1: Chunk 列表 -->
      <div v-if="step === 1" class="step-body">
        <div class="mb-3 flex items-center justify-between">
          <span class="text-sm text-gray-500">
            文件：<span class="font-medium text-gray-700">{{ fileName }}</span>
            &nbsp;·&nbsp;共 <span class="font-medium">{{ chunks.length }}</span> 个 chunk
          </span>
          <a-button size="small" @click="step = 0">重新选择</a-button>
        </div>

        <a-alert
          class="mb-3"
          :closable="false"
          show-icon
          type="info"
        >
          <template #message>
            点击「在此之前」或「在此之后」选择注入位置。
            <span class="text-gray-500">灰色按钮为非法位置（违反 PNG chunk 顺序约束）。</span>
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
              <span class="font-mono text-xs text-gray-500">{{ record.dataHex || '(空)' }}</span>
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
                    在此之前
                  </a-button>
                </a-tooltip>
                <a-tooltip :title="getDisabledReason('after', record)">
                  <a-button
                    type="link"
                    size="small"
                    :disabled="!isLegalPosition('after', record)"
                    @click="pickPosition('after', record)"
                  >
                    在此之后
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
            注入位置：在
            <a-tag :color="targetChunk?.critical ? 'red' : 'blue'" class="font-mono mx-1">
              {{ targetChunk?.type }}
            </a-tag>
            （索引 {{ targetChunk?.index }}）
            <strong>{{ position === 'before' ? '之前' : '之后' }}</strong>
          </template>
        </a-alert>

        <a-form :model="form" layout="vertical">
          <a-form-item label="Chunk 类型（4个字符）" required>
            <a-input
              v-model:value="form.chunkType"
              maxlength="4"
              show-count
              placeholder="例如: rPOC"
              style="width: 200px"
              :status="form.chunkType.length > 0 && form.chunkType.length !== 4 ? 'error' : ''"
            />
            <div class="text-xs text-gray-400 mt-1">
              首字母小写 = 辅助 chunk（推荐），大写 = 关键 chunk
            </div>
          </a-form-item>

          <a-form-item label="Payload" required>
            <a-textarea
              v-model:value="form.payload"
              :rows="8"
              placeholder="输入要注入的数据内容"
              class="font-mono"
            />
            <div class="text-xs text-gray-400 mt-1">
              {{ form.payload.length }} 字节
            </div>
          </a-form-item>

          <a-form-item label="输出文件路径" required>
            <a-input-group compact>
              <a-input
                v-model:value="form.outputPath"
                readonly
                placeholder="点击右侧按钮选择保存路径"
                style="width: calc(100% - 100px)"
              />
              <a-button @click="selectOutput">选择路径</a-button>
            </a-input-group>
          </a-form-item>

          <a-form-item>
            <a-space>
              <a-button type="primary" :loading="loading" @click="doInject">
                执行注入
              </a-button>
              <a-button @click="step = 1">返回</a-button>
            </a-space>
          </a-form-item>
        </a-form>
      </div>

      <!-- Step 3: 完成 -->
      <div v-if="step === 3" class="step-body">
        <a-result status="success" title="注入成功">
          <template #subTitle>
            文件已保存至：<span class="font-mono text-sm">{{ form.outputPath }}</span>
          </template>
          <template #extra>
            <a-button type="primary" @click="reset">注入新文件</a-button>
          </template>
        </a-result>

        <a-descriptions bordered size="small" class="mt-4" :column="2">
          <a-descriptions-item label="原始文件">{{ fileName }}</a-descriptions-item>
          <a-descriptions-item label="Chunk 类型">
            <span class="font-mono">{{ form.chunkType }}</span>
          </a-descriptions-item>
          <a-descriptions-item label="注入位置">
            {{ targetChunk?.type }}（索引 {{ targetChunk?.index }}）{{ position === 'before' ? ' 之前' : ' 之后' }}
          </a-descriptions-item>
          <a-descriptions-item label="Payload 大小">
            {{ formatBytes(form.payload.length) }}
          </a-descriptions-item>
          <a-descriptions-item label="输出路径" :span="2">
            <span class="font-mono text-xs">{{ form.outputPath }}</span>
          </a-descriptions-item>
        </a-descriptions>
      </div>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { message } from 'ant-design-vue'
import { FileImageOutlined } from '@ant-design/icons-vue'
import * as PNGInjectorService from '../../../bindings/github.com/Aliuyanfeng/happytools/backend/services/pnginjector/pnginjectorservice'

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

const chunkCols = [
  { title: '索引', dataIndex: 'index', key: 'index', width: 60 },
  { title: 'Chunk 类型', dataIndex: 'type', key: 'type', width: 130 },
  { title: '数据长度', dataIndex: 'length', key: 'length', width: 100 },
  { title: 'CRC', dataIndex: 'crc', key: 'crc', width: 110 },
  { title: '数据预览 (Hex)', dataIndex: 'dataHex', key: 'dataHex', ellipsis: true },
  { title: '注入位置', key: 'action', width: 160 }
]

// 选择 PNG 文件
const selectFile = async () => {
  try {
    const path = await PNGInjectorService.OpenFileDialog()
    if (!path) return

    filePath.value = path
    fileName.value = path.replace(/\\/g, '/').split('/').pop() || path

    loading.value = true
    const result = await PNGInjectorService.ParsePNG(path)
    chunks.value = result
    step.value = 1
    message.success(`解析完成，共 ${result.length} 个 chunk`)
  } catch (e: any) {
    message.error('解析失败：' + (e?.message || e))
  } finally {
    loading.value = false
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

  if (pos === 'before' && chunk.type === 'IHDR') return 'IHDR 必须是第一个 chunk，不能在其之前注入'
  if (pos === 'after' && chunk.type === 'IEND') return 'IEND 必须是最后一个 chunk，不能在其之后注入'
  if (pos === 'before' && chunk.type === 'IDAT' && prev?.type === 'IDAT') return '不能在连续 IDAT 块之间注入'
  if (pos === 'after' && chunk.type === 'IDAT' && next?.type === 'IDAT') return '不能在连续 IDAT 块之间注入'
  return '非法注入位置'
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
    message.error('选择路径失败：' + (e?.message || e))
  }
}

// 执行注入
const doInject = async () => {
  if (form.chunkType.length !== 4) {
    message.warning('Chunk 类型必须是 4 个字符')
    return
  }
  if (!form.payload) {
    message.warning('请输入 payload')
    return
  }
  if (!form.outputPath) {
    message.warning('请选择输出路径')
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
    message.error('注入失败：' + (e?.message || e))
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

.select-area:hover {
  border-color: #1890ff;
  background: #f0f8ff;
}

.select-icon {
  font-size: 56px;
  color: #bfbfbf;
}

.select-area:hover .select-icon {
  color: #1890ff;
}

.font-mono {
  font-family: 'Courier New', Courier, monospace;
}
</style>
