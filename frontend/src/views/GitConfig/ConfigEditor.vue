<template>
  <div class="config-editor">
    <!-- 无仓库选中 -->
    <div v-if="!store.activeRepoID" class="empty-state">
      <BranchesOutlined class="empty-icon" />
      <p>{{ t('gitConfig.addRepo') }}</p>
    </div>

    <template v-else>
      <!-- 工具栏 -->
      <div class="editor-toolbar">
        <div class="repo-meta">
          <a-tag :color="platformColor(store.activeRepo?.platform)">
            {{ store.activeRepo?.platform?.toUpperCase() }}
          </a-tag>
          <span class="repo-path-text">{{ store.activeRepo?.path }}</span>
        </div>
        <div class="toolbar-right">
          <a-input
            v-model:value="store.searchKeyword"
            :placeholder="t('gitConfig.searchPlaceholder')"
            allow-clear
            style="width: 200px"
          >
            <template #prefix><SearchOutlined /></template>
          </a-input>
          <a-button type="primary" ghost @click="addSectionModal = true">
            <template #icon><PlusOutlined /></template>
            {{ t('gitConfig.addSection') }}
          </a-button>
        </div>
      </div>

      <!-- QuickPanel -->
      <div class="editor-body">
        <QuickPanel />

        <!-- Section 列表 -->
        <a-spin :spinning="store.loading">
          <SectionCard
            v-for="sec in store.filteredSections"
            :key="sec.name + sec.subKey"
            :section="sec"
          />
          <div v-if="store.filteredSections.length === 0 && !store.loading" class="no-sections">
            暂无配置节
          </div>
        </a-spin>
      </div>
    </template>

    <!-- 添加节 Modal -->
    <a-modal
      v-model:open="addSectionModal"
      :title="t('gitConfig.addSection')"
      @ok="handleAddSection"
      @cancel="addSectionModal = false"
    >
      <a-form layout="vertical" class="mt-2">
        <a-form-item :label="t('gitConfig.sectionName')" required>
          <a-input v-model:value="newSection.name" placeholder="例如: remote" />
        </a-form-item>
        <a-form-item :label="t('gitConfig.subKey')">
          <a-input v-model:value="newSection.subKey" placeholder="例如: origin（可选）" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useI18n } from 'vue-i18n'
import { PlusOutlined, SearchOutlined, BranchesOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { useGitConfigStore } from '@/stores/gitconfig'
import QuickPanel from './QuickPanel.vue'
import SectionCard from './SectionCard.vue'

const { t } = useI18n()
const store = useGitConfigStore()

const addSectionModal = ref(false)
const newSection = reactive({ name: '', subKey: '' })

function platformColor(platform?: string): string {
  const colors: Record<string, string> = {
    github: 'default', gitlab: 'orange', gitee: 'red', custom: 'purple',
  }
  return colors[platform ?? ''] ?? 'default'
}

async function handleAddSection() {
  if (!newSection.name.trim()) {
    message.warning(t('gitConfig.sectionEmpty'))
    return
  }
  try {
    await store.addSection(newSection.name.trim(), newSection.subKey.trim())
    message.success(t('gitConfig.saveSuccess'))
    addSectionModal.value = false
    newSection.name = ''
    newSection.subKey = ''
  } catch (e: any) {
    message.error(e?.message ?? t('gitConfig.saveFailed'))
  }
}
</script>

<style scoped>
.config-editor {
  flex: 1;
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
}

.empty-state {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #bfbfbf;
  gap: 12px;
}

.empty-icon { font-size: 64px; }

.editor-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  background: #fff;
  border-bottom: 1px solid #e8e8e8;
  flex-shrink: 0;
}

.repo-meta { display: flex; align-items: center; gap: 8px; }
.repo-path-text { font-size: 12px; color: #8c8c8c; font-family: monospace; }
.toolbar-right { display: flex; align-items: center; gap: 8px; }

.editor-body {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  background: #f5f7fa;
}

.no-sections {
  text-align: center;
  color: #bfbfbf;
  padding: 40px;
  font-size: 14px;
}
</style>
