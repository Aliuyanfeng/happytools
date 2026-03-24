import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import {
  AddRepository,
  ListRepositories,
  DeleteRepository,
  LoadConfig,
  SaveEntry,
  DeleteEntry,
  AddSection,
  DeleteSection,
  GetKnownKeys,
  GetQuickPanel,
  SaveQuickPanel,
  GetRemoteNames,
} from '../../bindings/github.com/Aliuyanfeng/happytools/backend/services/gitconfig/gitconfigservice.js'

export interface Repository {
  id: string
  name: string
  path: string
  platform: string
  createdAt: string
}

export interface ConfigEntry {
  key: string
  value: string
}

export interface ConfigSection {
  name: string
  subKey: string
  entries: ConfigEntry[]
}

export interface KnownKey {
  section: string
  key: string
  type: string
  default: string
  enumValues: string[]
  descZh: string
  descEn: string
}

export interface QuickPanelItem {
  section: string
  subKey: string
  key: string
  order: number
}

export const useGitConfigStore = defineStore('gitconfig', () => {
  const repos = ref<Repository[]>([])
  const activeRepoID = ref<string | null>(null)
  const sections = ref<ConfigSection[]>([])
  const quickPanel = ref<QuickPanelItem[]>([])
  const knownKeys = ref<KnownKey[]>([])
  const remoteNames = ref<string[]>([])
  const searchKeyword = ref('')
  const loading = ref(false)

  const activeRepo = computed(() => repos.value.find(r => r.id === activeRepoID.value) ?? null)

  const filteredSections = computed(() => {
    const kw = searchKeyword.value.trim().toLowerCase()
    if (!kw) return sections.value
    return sections.value
      .map(sec => ({
        ...sec,
        entries: sec.entries.filter(
          e => e.key.toLowerCase().includes(kw) || e.value.toLowerCase().includes(kw)
        ),
      }))
      .filter(sec => sec.entries.length > 0)
  })

  async function loadRepos() {
    const result = await ListRepositories()
    repos.value = (result ?? []).filter(Boolean) as Repository[]
  }

  async function selectRepo(id: string) {
    activeRepoID.value = id
    await Promise.all([loadConfig(id), loadQuickPanel(id), loadRemoteNames(id)])
  }

  async function loadConfig(repoID: string) {
    loading.value = true
    try {
      const result = await LoadConfig(repoID)
      sections.value = (result ?? []) as ConfigSection[]
    } finally {
      loading.value = false
    }
  }

  async function addRepo(name: string, path: string, platform: string) {
    const repo = await AddRepository(name, path, platform)
    if (repo) {
      repos.value.push(repo as Repository)
    }
    return repo
  }

  async function deleteRepo(id: string) {
    await DeleteRepository(id)
    repos.value = repos.value.filter(r => r.id !== id)
    if (activeRepoID.value === id) {
      activeRepoID.value = null
      sections.value = []
      quickPanel.value = []
    }
  }

  async function saveEntry(section: string, subKey: string, key: string, value: string) {
    if (!activeRepoID.value) return
    await SaveEntry(activeRepoID.value, section, subKey, key, value)
    await loadConfig(activeRepoID.value)
  }

  async function deleteEntry(section: string, subKey: string, key: string) {
    if (!activeRepoID.value) return
    await DeleteEntry(activeRepoID.value, section, subKey, key)
    await loadConfig(activeRepoID.value)
  }

  async function addSection(section: string, subKey: string) {
    if (!activeRepoID.value) return
    await AddSection(activeRepoID.value, section, subKey)
    await loadConfig(activeRepoID.value)
  }

  async function deleteSection(section: string, subKey: string) {
    if (!activeRepoID.value) return
    await DeleteSection(activeRepoID.value, section, subKey)
    await loadConfig(activeRepoID.value)
  }

  async function loadQuickPanel(repoID: string) {
    const result = await GetQuickPanel(repoID)
    quickPanel.value = (result ?? []) as QuickPanelItem[]
  }

  async function loadRemoteNames(repoID: string) {
    const result = await GetRemoteNames(repoID)
    remoteNames.value = (result ?? []) as string[]
  }

  async function saveQuickPanel(items: QuickPanelItem[]) {
    if (!activeRepoID.value) return
    await SaveQuickPanel(activeRepoID.value, items)
    quickPanel.value = items
  }

  async function loadKnownKeys() {
    if (knownKeys.value.length > 0) return
    const result = await GetKnownKeys()
    knownKeys.value = (result ?? []) as KnownKey[]
  }

  return {
    repos, activeRepoID, activeRepo, sections, quickPanel, knownKeys,
    remoteNames, searchKeyword, loading, filteredSections,
    loadRepos, selectRepo, loadConfig, addRepo, deleteRepo,
    saveEntry, deleteEntry, addSection, deleteSection,
    loadQuickPanel, saveQuickPanel, loadKnownKeys, loadRemoteNames,
  }
})
