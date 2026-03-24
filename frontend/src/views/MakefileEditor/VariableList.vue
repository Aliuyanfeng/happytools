<template>
  <div class="variable-list">
    <div class="list-header">
      <span class="section-title">{{ t('makefileEditor.variables') }}</span>
      <a-button
        size="small"
        type="primary"
        :disabled="!store.currentDoc"
        @click="openAddModal"
      >
        + {{ t('makefileEditor.addVariable') }}
      </a-button>
    </div>

    <a-table
      :columns="columns"
      :data-source="variables"
      :pagination="false"
      size="small"
      row-key="name"
      :locale="{ emptyText: t('makefileEditor.noRecentFiles') }"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'operator'">
          <a-tag>{{ record.operator }}</a-tag>
        </template>
        <template v-else-if="column.key === 'value'">
          <span class="value-cell">{{ record.value }}</span>
        </template>
        <template v-else-if="column.key === 'actions'">
          <div class="action-btns">
            <a-button size="small" @click="openEditModal(record)">
              {{ t('makefileEditor.editVariable') }}
            </a-button>
            <a-popconfirm
              :title="t('makefileEditor.deleteVariable') + '?'"
              :ok-text="t('common.confirm')"
              :cancel-text="t('common.cancel')"
              @confirm="handleDelete(record.name)"
            >
              <a-button size="small" danger>
                {{ t('makefileEditor.deleteVariable') }}
              </a-button>
            </a-popconfirm>
          </div>
        </template>
      </template>
    </a-table>

    <!-- Add / Edit modal -->
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
  gap: 8px;
}

.list-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.section-title {
  font-size: 11px;
  font-weight: 600;
  color: #8c8c8c;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.value-cell {
  font-family: monospace;
  font-size: 12px;
  word-break: break-all;
}

.action-btns {
  display: flex;
  gap: 4px;
}
</style>
