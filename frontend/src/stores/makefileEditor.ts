import { defineStore } from 'pinia'
import { ref } from 'vue'
import {
  OpenFile,
  NewFile,
  NewFromTemplate,
  SaveFile,
  GetRecentFiles,
  GetTemplates,
  ValidateDependencies,
  RemoveRecentFile,
  ValidateMakefileFormat,
} from '../../bindings/github.com/Aliuyanfeng/happytools/backend/services/makefile/makefileservice.js'
import type { MakefileDoc, Target, Variable, Template } from '../../bindings/github.com/Aliuyanfeng/happytools/backend/services/makefile/models.js'

export const useMakefileEditorStore = defineStore('makefileEditor', () => {
  // ── State ──────────────────────────────────────────────────────────────────
  const currentDoc = ref<MakefileDoc | null>(null)
  const currentPath = ref<string>('')
  const isDirty = ref<boolean>(false)
  const recentFiles = ref<string[]>([])
  const templates = ref<Template[]>([])
  const selectedTargetName = ref<string>('')
  const cycleWarnings = ref<string[][]>([])

  // ── Cycle detection ────────────────────────────────────────────────────────
  async function checkCycles() {
    if (!currentDoc.value) return
    try {
      const result = await ValidateDependencies(currentDoc.value)
      cycleWarnings.value = result ?? []
    } catch {
      cycleWarnings.value = []
    }
  }

  // ── File actions ───────────────────────────────────────────────────────────
  async function loadFile(path: string) {
    const doc = await OpenFile(path)
    currentDoc.value = doc
    currentPath.value = path
    isDirty.value = false
    await fetchRecentFiles()
  }

  async function newFile(dir: string) {
    const doc = await NewFile(dir)
    currentDoc.value = doc
    currentPath.value = dir.endsWith('/') ? `${dir}Makefile` : `${dir}/Makefile`
    isDirty.value = false
  }

  async function newFromTemplate(dir: string, templateID: string) {
    const doc = await NewFromTemplate(dir, templateID)
    currentDoc.value = doc
    currentPath.value = dir.endsWith('/') ? `${dir}Makefile` : `${dir}/Makefile`
    isDirty.value = false
  }

  async function saveFile() {
    if (!currentDoc.value || !currentPath.value) return
    await SaveFile(currentPath.value, currentDoc.value)
    isDirty.value = false
  }

  // ── Target actions ─────────────────────────────────────────────────────────
  function setSelectedTarget(name: string) {
    selectedTargetName.value = name
  }

  function addTarget(target: Target) {
    if (!currentDoc.value) return
    currentDoc.value.targets.push(target)
    isDirty.value = true
    checkCycles()
  }

  function updateTarget(name: string, updated: Target) {
    if (!currentDoc.value) return
    const idx = currentDoc.value.targets.findIndex(t => t.name === name)
    if (idx !== -1) {
      currentDoc.value.targets[idx] = updated
      isDirty.value = true
      checkCycles()
    }
  }

  function deleteTarget(name: string) {
    if (!currentDoc.value) return
    currentDoc.value.targets = currentDoc.value.targets.filter(t => t.name !== name)
    isDirty.value = true
    // Clear selection if the deleted target was selected
    if (selectedTargetName.value === name) {
      selectedTargetName.value = ''
    }
  }

  // ── Variable actions ───────────────────────────────────────────────────────
  function addVariable(variable: Variable) {
    if (!currentDoc.value) return
    currentDoc.value.variables.push(variable)
    isDirty.value = true
  }

  function updateVariable(name: string, updated: Variable) {
    if (!currentDoc.value) return
    const idx = currentDoc.value.variables.findIndex(v => v.name === name)
    if (idx !== -1) {
      currentDoc.value.variables[idx] = updated
      isDirty.value = true
    }
  }

  function deleteVariable(name: string) {
    if (!currentDoc.value) return
    currentDoc.value.variables = currentDoc.value.variables.filter(v => v.name !== name)
    isDirty.value = true
  }

  // ── Template actions ───────────────────────────────────────────────────────
  async function mergeTemplate(templateID: string) {
    if (!currentDoc.value) return
    const tmpl = templates.value.find(t => t.id === templateID)
    if (!tmpl) return

    const existingTargetNames = new Set(currentDoc.value.targets.map(t => t.name))
    const existingVarNames = new Set(currentDoc.value.variables.map(v => v.name))

    for (const target of tmpl.doc.targets) {
      if (!existingTargetNames.has(target.name)) {
        currentDoc.value.targets.push(target)
      }
    }
    for (const variable of tmpl.doc.variables) {
      if (!existingVarNames.has(variable.name)) {
        currentDoc.value.variables.push(variable)
      }
    }
    isDirty.value = true
  }

  // ── Fetch helpers ──────────────────────────────────────────────────────────
  async function fetchRecentFiles() {
    const result = await GetRecentFiles()
    recentFiles.value = result ?? []
  }

  async function removeRecentFile(path: string) {
    await RemoveRecentFile(path)
    recentFiles.value = recentFiles.value.filter(p => p !== path)
    // 如果移除的是当前打开的文件，不影响编辑状态
  }

  function closeFile() {
    currentDoc.value = null
    currentPath.value = ''
    isDirty.value = false
    selectedTargetName.value = ''
    cycleWarnings.value = []
  }

  async function validateAndLoadFile(path: string) {
    await ValidateMakefileFormat(path)  // 抛出错误则中断
    await loadFile(path)
  }

  async function fetchTemplates() {
    const result = await GetTemplates()
    templates.value = result ?? []
  }

  return {
    // State
    currentDoc,
    currentPath,
    isDirty,
    recentFiles,
    templates,
    selectedTargetName,
    cycleWarnings,
    // Actions
    loadFile,
    newFile,
    newFromTemplate,
    saveFile,
    setSelectedTarget,
    addTarget,
    updateTarget,
    deleteTarget,
    addVariable,
    updateVariable,
    deleteVariable,
    mergeTemplate,
    checkCycles,
    fetchRecentFiles,
    fetchTemplates,
    removeRecentFile,
    validateAndLoadFile,
    closeFile,
  }
})
