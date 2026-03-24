<template>
  <div class="repo-list">
    <div class="list-header">
      <span class="header-title">{{ t('gitConfig.title') }}</span>
      <a-button size="small" type="primary" @click="addModalVisible = true">
        <template #icon><PlusOutlined /></template>
      </a-button>
    </div>

    <!-- 平台筛选 -->
    <div class="platform-filter">
      <a-tag
        v-for="p in platforms"
        :key="p.value"
        :color="activePlatform === p.value ? 'blue' : ''"
        class="filter-tag"
        @click="activePlatform = p.value"
      >
        {{ p.label }}
      </a-tag>
    </div>

    <!-- 仓库列表 -->
    <div class="repo-items">
      <div
        v-for="repo in filteredRepos"
        :key="repo.id"
        class="repo-item"
        :class="{ active: store.activeRepoID === repo.id }"
        @click="store.selectRepo(repo.id)"
      >
        <div class="repo-icon">{{ platformIcon(repo.platform) }}</div>
        <div class="repo-info">
          <div class="repo-name">{{ repo.name }}</div>
          <div class="repo-path">{{ shortPath(repo.path) }}</div>
        </div>
        <a-popconfirm
          :title="t('gitConfig.deleteRepoConfirm')"
          @confirm.stop="store.deleteRepo(repo.id)"
        >
          <a-button
            size="small"
            type="text"
            danger
            class="delete-btn"
            @click.stop
          >
            <template #icon><DeleteOutlined /></template>
          </a-button>
        </a-popconfirm>
      </div>

      <div v-if="filteredRepos.length === 0" class="empty-repos">
        <BranchesOutlined class="empty-icon" />
        <div>{{ t('gitConfig.addRepo') }}</div>
      </div>
    </div>

    <!-- 添加仓库 Modal -->
    <a-modal
      v-model:open="addModalVisible"
      :title="t('gitConfig.addRepo')"
      @ok="handleAddRepo"
      @cancel="resetForm"
      :confirm-loading="adding"
    >
      <a-form :model="form" layout="vertical" class="mt-2">
        <a-form-item :label="t('gitConfig.repoName')" required>
          <a-input v-model:value="form.name" :placeholder="t('gitConfig.repoName')" />
        </a-form-item>
        <a-form-item :label="t('gitConfig.repoPath')" required>
          <a-input-group compact>
            <a-input
              v-model:value="form.path"
              :placeholder="t('gitConfig.repoPath')"
              style="width: calc(100% - 100px)"
            />
            <a-button @click="selectDirectory">{{ t('gitConfig.selectDirectory') }}</a-button>
          </a-input-group>
        </a-form-item>
        <a-form-item :label="t('gitConfig.platform')">
          <a-select v-model:value="form.platform" style="width: 100%">
            <a-select-option value="github">GitHub</a-select-option>
            <a-select-option value="gitlab">GitLab</a-select-option>
            <a-select-option value="gitee">Gitee</a-select-option>
            <a-select-option value="custom">{{ t('gitConfig.platformCustom') }}</a-select-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, reactive } from 'vue'
import { useI18n } from 'vue-i18n'
import { PlusOutlined, DeleteOutlined, BranchesOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { useGitConfigStore } from '@/stores/gitconfig'
import { OpenDirectoryDialog } from '../../../bindings/github.com/Aliuyanfeng/happytools/backend/services/gitconfig/gitconfigservice.js'

const { t } = useI18n()
const store = useGitConfigStore()

const addModalVisible = ref(false)
const adding = ref(false)
const activePlatform = ref('all')

const form = reactive({ name: '', path: '', platform: 'github' })

const platforms = [
  { label: '全部', value: 'all' },
  { label: 'GitHub', value: 'github' },
  { label: 'GitLab', value: 'gitlab' },
  { label: 'Gitee', value: 'gitee' },
  { label: '自定义', value: 'custom' },
]

const filteredRepos = computed(() =>
  activePlatform.value === 'all'
    ? store.repos
    : store.repos.filter(r => r.platform === activePlatform.value)
)

function platformIcon(platform: string): string {
  const icons: Record<string, string> = {
    github: '🐙', gitlab: '🦊', gitee: '🐧', custom: '📁',
  }
  return icons[platform] ?? '📁'
}

function shortPath(path: string): string {
  if (path.length <= 30) return path
  return '...' + path.slice(-27)
}

async function selectDirectory() {
  const dir = await OpenDirectoryDialog()
  if (dir) {
    form.path = dir
    if (!form.name) {
      form.name = dir.split(/[\\/]/).pop() ?? ''
    }
  }
}

async function handleAddRepo() {
  if (!form.name.trim() || !form.path.trim()) {
    message.warning('请填写仓库名称和路径')
    return
  }
  adding.value = true
  try {
    await store.addRepo(form.name.trim(), form.path.trim(), form.platform)
    message.success(t('common.success'))
    addModalVisible.value = false
    resetForm()
  } catch (e: any) {
    message.error(e?.message ?? t('gitConfig.saveFailed'))
  } finally {
    adding.value = false
  }
}

function resetForm() {
  form.name = ''
  form.path = ''
  form.platform = 'github'
}
</script>

<style scoped>
.repo-list {
  width: 220px;
  flex-shrink: 0;
  background: #fff;
  border-right: 1px solid #e8e8e8;
  display: flex;
  flex-direction: column;
  height: 100%;
}

.list-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 12px 10px;
  border-bottom: 1px solid #f0f0f0;
}

.header-title { font-weight: 600; font-size: 14px; color: #262626; }

.platform-filter {
  padding: 8px 12px;
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  border-bottom: 1px solid #f0f0f0;
}

.filter-tag { cursor: pointer; margin: 0; }

.repo-items {
  flex: 1;
  overflow-y: auto;
  padding: 6px 0;
}

.repo-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  cursor: pointer;
  transition: background 0.15s;
  position: relative;
}

.repo-item:hover { background: #f5f5f5; }
.repo-item.active { background: #e6f7ff; border-right: 2px solid #1890ff; }

.repo-icon { font-size: 18px; flex-shrink: 0; }

.repo-info { flex: 1; min-width: 0; }
.repo-name { font-size: 13px; font-weight: 500; color: #262626; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.repo-path { font-size: 11px; color: #8c8c8c; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }

.delete-btn { opacity: 0; transition: opacity 0.15s; }
.repo-item:hover .delete-btn { opacity: 1; }

.empty-repos {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 40px 16px;
  color: #bfbfbf;
  font-size: 13px;
}

.empty-icon { font-size: 32px; }
</style>
