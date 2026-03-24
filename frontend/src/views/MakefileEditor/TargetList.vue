<template>
  <div class="target-list">
    <div class="list-header">
      <span class="section-title">{{ t('makefileEditor.targets') }}</span>
      <a-button
        size="small"
        type="primary"
        :disabled="!store.currentDoc"
        @click="openAddModal"
      >
        + {{ t('makefileEditor.addTarget') }}
      </a-button>
    </div>

    <a-table
      :columns="columns"
      :data-source="targets"
      :pagination="false"
      size="small"
      row-key="name"
      :locale="{ emptyText: t('makefileEditor.noRecentFiles') }"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'name'">
          <span class="target-name">{{ record.name }}</span>
          <a-tag v-if="record.isPhony" color="blue" style="margin-left: 6px">PHONY</a-tag>
        </template>
        <template v-else-if="column.key === 'deps'">
          <span class="count-badge">{{ record.deps?.length ?? 0 }}</span>
        </template>
        <template v-else-if="column.key === 'commands'">
          <span class="count-badge">{{ record.commands?.length ?? 0 }}</span>
        </template>
        <template v-else-if="column.key === 'actions'">
          <div class="action-btns">
            <a-button size="small" @click="openEditModal(record)">
              {{ t('makefileEditor.editTarget') }}
            </a-button>
            <a-popconfirm
              :title="t('makefileEditor.deleteTarget') + '?'"
              :ok-text="t('common.confirm')"
              :cancel-text="t('common.cancel')"
              @confirm="handleDelete(record.name)"
            >
              <a-button size="small" danger>
                {{ t('makefileEditor.deleteTarget') }}
              </a-button>
            </a-popconfirm>
          </div>
        </template>
      </template>
    </a-table>

    <!-- Add / Edit modal -->
    <a-modal
      v-model:open="modalVisible"
      :title="editingTarget ? t('makefileEditor.editTarget') : t('makefileEditor.addTarget')"
      :footer="null"
      destroy-on-close
      @cancel="closeModal"
    >
      <TargetForm
        :initial-value="editingTarget ?? undefined"
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
import type { Target } from '../../../bindings/github.com/Aliuyanfeng/happytools/backend/services/makefile/models.js'
import TargetForm from './TargetForm.vue'

const { t } = useI18n()
const store = useMakefileEditorStore()

const targets = computed(() => store.currentDoc?.targets ?? [])
const existingNames = computed(() => targets.value.map((t) => t.name))

const modalVisible = ref(false)
const editingTarget = ref<Target | null>(null)

const columns = [
  { title: t('makefileEditor.targetName'), key: 'name', width: '35%' },
  { title: t('makefileEditor.dependencies'), key: 'deps', width: '15%' },
  { title: t('makefileEditor.commands'), key: 'commands', width: '15%' },
  { title: '', key: 'actions', width: '200px' },
]

function openAddModal() {
  editingTarget.value = null
  modalVisible.value = true
}

function openEditModal(target: Target) {
  editingTarget.value = { ...target, deps: [...target.deps], commands: [...target.commands] }
  modalVisible.value = true
}

function closeModal() {
  modalVisible.value = false
  editingTarget.value = null
}

function handleFormSubmit(target: Target) {
  if (editingTarget.value) {
    store.updateTarget(editingTarget.value.name, target)
    message.success(t('makefileEditor.editTarget') + ' 成功')
  } else {
    store.addTarget(target)
    message.success(t('makefileEditor.addTarget') + ' 成功')
  }
  closeModal()
}

function handleDelete(name: string) {
  store.deleteTarget(name)
  message.success(t('makefileEditor.deleteTarget') + ' 成功')
}
</script>

<style scoped>
.target-list {
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

.target-name {
  font-family: monospace;
  font-size: 12px;
}

.count-badge {
  font-size: 12px;
  color: #595959;
}

.action-btns {
  display: flex;
  gap: 4px;
}
</style>
