<template>
  <div class="metadata-parser">
    <a-card :title="t('pdfParser.metadata.title')" :bordered="false">
      <!-- 文件选择区 -->
      <div class="file-select-section">
        <a-space>
          <a-button type="primary" @click="openFile" :loading="loading">
            <template #icon><FolderOpenOutlined /></template>
            {{ t('pdfParser.metadata.selectFile') }}
          </a-button>
          <span v-if="filePath" class="file-path">{{ filePath }}</span>
        </a-space>
      </div>

      <a-divider v-if="metadata" />

      <!-- 元数据展示与编辑区 -->
      <div v-if="metadata" class="metadata-content">
        <!-- 基本信息 -->
        <a-collapse v-model:activeKey="activePanels" :bordered="false">
          <!-- 标准元数据 -->
          <a-collapse-panel key="standard" :header="t('pdfParser.metadata.standardMetadata')">
            <a-form layout="vertical" class="metadata-form">
              <a-row :gutter="16">
                <a-col :span="12">
                  <a-form-item :label="t('pdfParser.metadata.title_field')">
                    <a-input v-model:value="editForm.title" :placeholder="t('pdfParser.metadata.titlePlaceholder')" />
                  </a-form-item>
                </a-col>
                <a-col :span="12">
                  <a-form-item :label="t('pdfParser.metadata.author')">
                    <a-input v-model:value="editForm.author" :placeholder="t('pdfParser.metadata.authorPlaceholder')" />
                  </a-form-item>
                </a-col>
              </a-row>
              <a-row :gutter="16">
                <a-col :span="24">
                  <a-form-item :label="t('pdfParser.metadata.subject')">
                    <a-textarea v-model:value="editForm.subject" :placeholder="t('pdfParser.metadata.subjectPlaceholder')" :rows="2" />
                  </a-form-item>
                </a-col>
              </a-row>
              <a-row :gutter="16">
                <a-col :span="12">
                  <a-form-item :label="t('pdfParser.metadata.creator')">
                    <a-input :value="metadata.creator" disabled />
                  </a-form-item>
                </a-col>
                <a-col :span="12">
                  <a-form-item :label="t('pdfParser.metadata.producer')">
                    <a-input :value="metadata.producer" disabled />
                  </a-form-item>
                </a-col>
              </a-row>
              <a-row :gutter="16">
                <a-col :span="12">
                  <a-form-item :label="t('pdfParser.metadata.creationDate')">
                    <a-input :value="metadata.creationDate" disabled />
                  </a-form-item>
                </a-col>
                <a-col :span="12">
                  <a-form-item :label="t('pdfParser.metadata.modificationDate')">
                    <a-input :value="metadata.modificationDate" disabled />
                  </a-form-item>
                </a-col>
              </a-row>
              <a-row :gutter="16">
                <a-col :span="12">
                  <a-form-item :label="t('pdfParser.metadata.pageCount')">
                    <a-input :value="metadata.pageCount" disabled />
                  </a-form-item>
                </a-col>
                <a-col :span="12">
                  <a-form-item :label="t('pdfParser.metadata.version')">
                    <a-input :value="'PDF ' + metadata.version" disabled />
                  </a-form-item>
                </a-col>
              </a-row>
              <a-button type="primary" @click="saveStandardMetadata" :loading="saving">
                <template #icon><SaveOutlined /></template>
                {{ t('pdfParser.metadata.saveStandard') }}
              </a-button>
            </a-form>
          </a-collapse-panel>

          <!-- 关键词 -->
          <a-collapse-panel key="keywords" :header="t('pdfParser.metadata.keywords')">
            <div class="keywords-section">
              <div class="keyword-tags">
                <a-tag
                  v-for="(keyword, index) in editForm.keywords"
                  :key="index"
                  closable
                  @close="removeKeyword(index)"
                  color="blue"
                >
                  {{ keyword }}
                </a-tag>
                <a-tag v-if="editForm.keywords.length === 0" color="default">
                  {{ t('pdfParser.metadata.noKeywords') }}
                </a-tag>
              </div>
              <div class="keyword-input-row">
                <a-input
                  v-model:value="newKeyword"
                  :placeholder="t('pdfParser.metadata.addKeywordPlaceholder')"
                  @pressEnter="addKeyword"
                  style="width: 300px"
                />
                <a-button type="dashed" @click="addKeyword" :disabled="!newKeyword.trim()">
                  <template #icon><PlusOutlined /></template>
                  {{ t('pdfParser.metadata.addKeyword') }}
                </a-button>
              </div>
              <a-button type="primary" @click="saveKeywords" :loading="saving" style="margin-top: 12px">
                <template #icon><SaveOutlined /></template>
                {{ t('pdfParser.metadata.saveKeywords') }}
              </a-button>
            </div>
          </a-collapse-panel>

          <!-- 自定义属性 -->
          <a-collapse-panel key="properties" :header="t('pdfParser.metadata.properties')">
            <div class="properties-section">
              <a-table
                :columns="propertyColumns"
                :data-source="propertyTableData"
                :pagination="false"
                size="small"
                :row-key="(record: PropertyRow) => record.key"
              >
                <template #bodyCell="{ column, record }">
                  <template v-if="column.key === 'action'">
                    <a-button type="link" danger size="small" @click="removeProperty(record.key)">
                      {{ t('pdfParser.metadata.removeProperty') }}
                    </a-button>
                  </template>
                </template>
              </a-table>
              <a-empty v-if="propertyTableData.length === 0" :description="t('pdfParser.metadata.noProperties')" />

              <a-divider />
              <div class="add-property-form">
                <h4>{{ t('pdfParser.metadata.addProperty') }}</h4>
                <a-row :gutter="8">
                  <a-col :span="8">
                    <a-input v-model:value="newPropertyKey" :placeholder="t('pdfParser.metadata.propertyKeyPlaceholder')" />
                  </a-col>
                  <a-col :span="12">
                    <a-input v-model:value="newPropertyValue" :placeholder="t('pdfParser.metadata.propertyValuePlaceholder')" />
                  </a-col>
                  <a-col :span="4">
                    <a-button type="dashed" @click="addProperty" :disabled="!newPropertyKey.trim() || !newPropertyValue.trim()" block>
                      <template #icon><PlusOutlined /></template>
                    </a-button>
                  </a-col>
                </a-row>
              </div>
              <a-button type="primary" @click="saveProperties" :loading="saving" style="margin-top: 12px">
                <template #icon><SaveOutlined /></template>
                {{ t('pdfParser.metadata.saveProperties') }}
              </a-button>
            </div>
          </a-collapse-panel>

          <!-- 文件信息 -->
          <a-collapse-panel key="fileinfo" :header="t('pdfParser.metadata.fileInfo')">
            <a-descriptions :column="2" size="small" bordered>
              <a-descriptions-item :label="t('pdfParser.metadata.tagged')">
                <a-tag :color="metadata.tagged ? 'green' : 'default'">
                  {{ metadata.tagged ? t('pdfParser.metadata.yes') : t('pdfParser.metadata.no') }}
                </a-tag>
              </a-descriptions-item>
              <a-descriptions-item :label="t('pdfParser.metadata.encrypted')">
                <a-tag :color="metadata.encrypted ? 'red' : 'green'">
                  {{ metadata.encrypted ? t('pdfParser.metadata.yes') : t('pdfParser.metadata.no') }}
                </a-tag>
              </a-descriptions-item>
            </a-descriptions>
          </a-collapse-panel>
        </a-collapse>
      </div>

      <!-- 空状态 -->
      <a-empty v-if="!metadata && !loading" :description="t('pdfParser.metadata.emptyHint')" />
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { message } from 'ant-design-vue'
import {
  FolderOpenOutlined,
  SaveOutlined,
  PlusOutlined,
} from '@ant-design/icons-vue'
import * as pdfService from '../../../bindings/github.com/Aliuyanfeng/happytools/backend/services/pdf/pdfservice'
import type { PDFMetadata } from '../../../bindings/github.com/Aliuyanfeng/happytools/backend/services/pdf/models.js'

const { t } = useI18n()

// 状态
const filePath = ref('')
const metadata = ref<PDFMetadata | null>(null)
const loading = ref(false)
const saving = ref(false)
const activePanels = ref(['standard', 'keywords', 'properties', 'fileinfo'])
const newKeyword = ref('')
const newPropertyKey = ref('')
const newPropertyValue = ref('')

// 编辑表单
const editForm = reactive({
  title: '',
  author: '',
  subject: '',
  keywords: [] as string[],
  properties: {} as Record<string, string>,
})

// 属性表格列定义
const propertyColumns = computed(() => [
  { title: t('pdfParser.metadata.propertyKey'), dataIndex: 'key', key: 'key' },
  { title: t('pdfParser.metadata.propertyValue'), dataIndex: 'value', key: 'value' },
  { title: t('pdfParser.metadata.action'), key: 'action', width: 100 },
])

// 属性表格数据
interface PropertyRow {
  key: string
  value: string
}

const propertyTableData = computed<PropertyRow[]>(() => {
  return Object.entries(editForm.properties).map(([key, value]) => ({ key, value }))
})

// 打开文件
const openFile = async () => {
  try {
    loading.value = true
    const path = await pdfService.OpenFileDialog()
    if (!path) return

    filePath.value = path
    const info = await pdfService.GetMetadata(path)
    if (!info) return
    metadata.value = info

    // 填充编辑表单
    editForm.title = info.title || ''
    editForm.author = info.author || ''
    editForm.subject = info.subject || ''
    editForm.keywords = [...(info.keywords || [])]
    editForm.properties = { ...(info.properties || {}) }
  } catch (err: any) {
    message.error(t('pdfParser.metadata.loadFailed') + ': ' + (err.message || err))
  } finally {
    loading.value = false
  }
}

// 保存标准元数据
const saveStandardMetadata = async () => {
  try {
    saving.value = true
    await pdfService.SetStandardMetadata(filePath.value, '', editForm.title, editForm.author, editForm.subject)
    message.success(t('pdfParser.metadata.saveSuccess'))
    // 刷新元数据
    const info = await pdfService.GetMetadata(filePath.value)
    if (info) metadata.value = info
  } catch (err: any) {
    message.error(t('pdfParser.metadata.saveFailed') + ': ' + (err.message || err))
  } finally {
    saving.value = false
  }
}

// 关键词操作
const addKeyword = () => {
  const kw = newKeyword.value.trim()
  if (!kw) return
  if (editForm.keywords.includes(kw)) {
    message.warning(t('pdfParser.metadata.keywordExists'))
    return
  }
  editForm.keywords.push(kw)
  newKeyword.value = ''
}

const removeKeyword = (index: number) => {
  editForm.keywords.splice(index, 1)
}

const saveKeywords = async () => {
  try {
    saving.value = true
    const originalKeywords = metadata.value?.keywords || []
    const toRemove = originalKeywords.filter(k => !editForm.keywords.includes(k))
    const toAdd = editForm.keywords.filter(k => !originalKeywords.includes(k))

    // 先移除旧关键词
    if (toRemove.length > 0) {
      await pdfService.RemoveKeywords(filePath.value, '', toRemove)
    }
    // 再添加新关键词
    if (toAdd.length > 0) {
      await pdfService.AddKeywords(filePath.value, '', toAdd)
    }

    message.success(t('pdfParser.metadata.saveSuccess'))
    // 刷新元数据
    const info = await pdfService.GetMetadata(filePath.value)
    if (info) {
      metadata.value = info
      editForm.keywords = [...(info.keywords || [])]
    }
  } catch (err: any) {
    message.error(t('pdfParser.metadata.saveFailed') + ': ' + (err.message || err))
  } finally {
    saving.value = false
  }
}

// 自定义属性操作
const addProperty = () => {
  const key = newPropertyKey.value.trim()
  const value = newPropertyValue.value.trim()
  if (!key || !value) return
  if (editForm.properties[key] !== undefined) {
    message.warning(t('pdfParser.metadata.propertyExists'))
    return
  }
  editForm.properties[key] = value
  newPropertyKey.value = ''
  newPropertyValue.value = ''
}

const removeProperty = (key: string) => {
  delete editForm.properties[key]
  // 触发响应式更新
  editForm.properties = { ...editForm.properties }
}

const saveProperties = async () => {
  try {
    saving.value = true
    const originalProps = metadata.value?.properties || {}
    const originalKeys = Object.keys(originalProps)
    const newKeys = Object.keys(editForm.properties)

    const toRemove = originalKeys.filter(k => !newKeys.includes(k))
    const toAdd: Record<string, string> = {}
    for (const [key, value] of Object.entries(editForm.properties)) {
      if (originalProps[key] !== value) {
        toAdd[key] = value
      }
    }

    // 先移除旧属性
    if (toRemove.length > 0) {
      await pdfService.RemoveProperties(filePath.value, '', toRemove)
    }
    // 再添加/更新属性
    if (Object.keys(toAdd).length > 0) {
      await pdfService.AddProperties(filePath.value, '', toAdd)
    }

    message.success(t('pdfParser.metadata.saveSuccess'))
    // 刷新元数据
    const info = await pdfService.GetMetadata(filePath.value)
    if (info) {
      metadata.value = info
      editForm.properties = { ...(info.properties || {}) }
    }
  } catch (err: any) {
    message.error(t('pdfParser.metadata.saveFailed') + ': ' + (err.message || err))
  } finally {
    saving.value = false
  }
}
</script>

<style scoped>
.metadata-parser {
  max-width: 900px;
  margin: 0 auto;
}

.file-select-section {
  margin-bottom: 8px;
}

.file-path {
  color: #666;
  font-size: 13px;
  max-width: 500px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  display: inline-block;
  vertical-align: middle;
}

.metadata-form {
  max-width: 800px;
}

.keywords-section {
  padding: 4px 0;
}

.keyword-tags {
  margin-bottom: 12px;
  min-height: 32px;
}

.keyword-tags .ant-tag {
  margin-bottom: 4px;
}

.keyword-input-row {
  display: flex;
  gap: 8px;
  align-items: center;
}

.properties-section {
  padding: 4px 0;
}

.add-property-form h4 {
  margin-bottom: 12px;
  font-weight: 500;
}

:deep(.ant-collapse-header) {
  font-weight: 600;
}

:deep(.ant-descriptions-item-label) {
  font-weight: 500;
}
</style>