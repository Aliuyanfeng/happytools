<template>
  <div class="variable-list">
    <div class="list-header">
      <span class="section-title">{{ t('makefileEditor.variables') }}</span>
      <a-button
        size="small"
        type="primary"
        :disabled="!store.currentDoc"
        @click="openAddModal"
      >+</a-button>
    </div>

    <div class="var-items">
      <div v-for="v in variables" :key="v.name" class="var-item">
        <div class="var-main">
          <span class="var-name">{{ v.name }}</span>
          <span class="var-op">{{ v.operator }}</span>
          <span class="var-val">{{ v.value }}</span>
        </div>
        <div class="var-actions">
          <a-button size="small" type="text" @click="openEditModal(v)">✎</a-button>
          <a-popconfirm
            :title="t('makefileEditor.deleteVariable') + '?'"
            :ok-text="t('common.confirm')"
            :cancel-text="t('common.cancel')"
            @confirm="handleDelete(v.name)"
          >
            <a-button size="small" type="text" danger>✕</a-button>
          </a-popconfirm>
        </div>
      </div>
      <div v-if="variables.length === 0" class="empty-hint">
        {{ t('makefileEditor.noRecentFiles') }}
      </div>
    </div>

    <a-modal
      v-model:open="modalVisible"
      :title="editingVariable ? t('makefileEditor.editVariable') : t('makefileEditor.addVariable')"
      :footer="null"
      destroy-on-close
      @cancel="closeModal"
    >
      <VariableForm
        :initial-value="editingVariable ?? undefined"
        :existing-names="existingNames"
        @submit="handleFormSubmit"
        @cancel="closeModal"
      />
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { message } from 'ant-design-vue'
import { useMakefileEditorStore } from '../../stores/makefileEditor'
import type { Variable } from '../../../bindings/github.com/Aliuyanfeng/happytools/backend/services/makefile/models.js'
import VariableForm from './VariableForm.vue'

const { t } = useI18n()
const store = useMakefileEditorStore()

const variables = computed(() => store.currentDoc?.variables ?? [])
const existingNames = computed(() => variables.value.map((v) => v.name))

const modalVisible = ref(false)
const editingVariable = ref<Variable | null>(null)

const columns = [
  { title: t('makefileEditor.variableName'), dataIndex: 'name', key: 'name', width: '35%' },
  { title: t('makefileEditor.operator'), dataIndex: 'operator', key: 'operator', width: '15%' },
  { title: t('makefileEditor.value'), dataIndex: 'value', key: 'value' },
  { title: '', key: 'actions', width: '180px' },
]

function openAddModal() {
  editingVariable.value = null
  modalVisible.value = true
}

function openEditModal(variable: Variable) {
  editingVariable.value = { ...variable }
  modalVisible.value = true
}

function closeModal() {
  modalVisible.value = false
  editingVariable.value = null
}

function handleFormSubmit(variable: Variable) {
  if (editingVariable.value) {
    store.updateVariable(editingVariable.value.name, variable)
    message.success(t('makefileEditor.editVariable') + ' 成功')
  } else {
    store.addVariable(variable)
    message.success(t('makefileEditor.addVariable') + ' 成功')
  }
  closeModal()
}

function handleDelete(name: string) {
  store.deleteVariable(name)
  message.success(t('makefileEditor.deleteVariable') + ' 成功')
}
</script>

<style scoped>
.variable-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
  min-height: 0;
}

.list-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-shrink: 0;
}

.section-title {
  font-size: 11px;
  font-weight: 600;
  color: #8c8c8c;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.var-items {
  display: flex;
  flex-direction: column;
  gap: 3px;
}

.var-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 4px;
  padding: 5px 7px;
  border-radius: 6px;
  background: #fafafa;
  border: 1px solid #f0f0f0;
  min-width: 0;
}
.var-item:hover { background: #f0f5ff; border-color: #d6e4ff; }

.var-main {
  display: flex;
  align-items: center;
  gap: 4px;
  min-width: 0;
  flex: 1;
  overflow: hidden;
}

.var-name {
  font-size: 11px;
  font-weight: 600;
  color: #262626;
  font-family: monospace;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 80px;
}

.var-op {
  font-size: 10px;
  color: #1677ff;
  font-family: monospace;
  flex-shrink: 0;
  background: #e6f4ff;
  padding: 0 3px;
  border-radius: 3px;
}

.var-val {
  font-size: 11px;
  color: #595959;
  font-family: monospace;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  flex: 1;
  min-width: 0;
}

.var-actions {
  display: flex;
  gap: 0;
  flex-shrink: 0;
}

.empty-hint {
  font-size: 12px;
  color: #bfbfbf;
  text-align: center;
  padding: 8px 0;
}
</style>
