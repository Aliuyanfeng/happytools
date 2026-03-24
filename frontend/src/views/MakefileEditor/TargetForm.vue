<template>
  <a-form
    ref="formRef"
    :model="formState"
    :rules="rules"
    layout="vertical"
    @finish="handleSubmit"
  >
    <!-- Target Name -->
    <a-form-item :label="t('makefileEditor.targetName')" name="name">
      <a-input
        v-model:value="formState.name"
        :placeholder="t('makefileEditor.targetName')"
        allow-clear
      />
    </a-form-item>

    <!-- PhonyTarget toggle -->
    <a-form-item :label="t('makefileEditor.isPhony')" name="isPhony">
      <a-switch v-model:checked="formState.isPhony" />
    </a-form-item>

    <!-- Dependencies multi-select -->
    <a-form-item :label="t('makefileEditor.dependencies')" name="deps">
      <a-select
        v-model:value="formState.deps"
        mode="multiple"
        :options="depOptions"
        :placeholder="t('makefileEditor.dependencies')"
        allow-clear
        style="width: 100%"
        @change="handleDepsChange"
      />
      <a-alert
        v-if="cycleDetected"
        type="warning"
        :message="t('makefileEditor.cycleWarning')"
        :description="cycleDescription"
        show-icon
        style="margin-top: 8px"
      />
    </a-form-item>

    <!-- Commands textarea -->
    <a-form-item :label="t('makefileEditor.commands')" name="commandsText">
      <a-textarea
        v-model:value="commandsText"
        :placeholder="t('makefileEditor.commands')"
        :auto-size="{ minRows: 4, maxRows: 12 }"
        style="font-family: monospace; font-size: 12px"
        @keydown="handleTextareaKeydown"
      />
      <div class="commands-hint">{{ t('makefileEditor.commandsHint') }}</div>
    </a-form-item>

    <div class="form-actions">
      <a-button @click="emit('cancel')">{{ t('common.cancel') }}</a-button>
      <a-button type="primary" html-type="submit">{{ t('common.save') }}</a-button>
    </div>
  </a-form>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import type { FormInstance, Rule } from 'ant-design-vue/es/form'
import { useMakefileEditorStore } from '../../stores/makefileEditor'
import type { Target } from '../../../bindings/github.com/Aliuyanfeng/happytools/backend/services/makefile/models.js'

const props = defineProps<{
  initialValue?: Target
  existingNames?: string[]
}>()

const emit = defineEmits<{
  (e: 'submit', target: Target): void
  (e: 'cancel'): void
}>()

const { t } = useI18n()
const store = useMakefileEditorStore()
const formRef = ref<FormInstance>()

// Commands are stored as string[] but edited as a single textarea string
const commandsText = ref<string>(
  (props.initialValue?.commands ?? []).join('\n'),
)

const formState = reactive<Target>({
  name: props.initialValue?.name ?? '',
  isPhony: props.initialValue?.isPhony ?? false,
  deps: [...(props.initialValue?.deps ?? [])],
  commands: [...(props.initialValue?.commands ?? [])],
})

// Sync when initialValue changes (e.g. modal reuse)
watch(
  () => props.initialValue,
  (val) => {
    formState.name = val?.name ?? ''
    formState.isPhony = val?.isPhony ?? false
    formState.deps = [...(val?.deps ?? [])]
    formState.commands = [...(val?.commands ?? [])]
    commandsText.value = (val?.commands ?? []).join('\n')
  },
)

// Build dep options from existingNames, excluding the target being edited
const depOptions = computed(() =>
  (props.existingNames ?? [])
    .filter((n) => n !== props.initialValue?.name)
    .map((n) => ({ label: n, value: n })),
)

// Cycle detection state
const cycleDetected = ref(false)
const cycleDescription = ref('')

async function handleDepsChange() {
  // Temporarily apply deps to a snapshot doc for cycle checking
  if (!store.currentDoc) return

  // Build a temporary doc with the current form's deps applied
  const tempTargets = store.currentDoc.targets.map((t) => {
    if (t.name === (props.initialValue?.name ?? formState.name)) {
      return { ...t, deps: [...formState.deps] }
    }
    return t
  })

  // If adding a new target (not yet in doc), add it temporarily
  const isNew = !store.currentDoc.targets.find(
    (t) => t.name === formState.name,
  )
  if (isNew && formState.name) {
    tempTargets.push({ ...formState, commands: [] })
  }

  const tempDoc = { ...store.currentDoc, targets: tempTargets }

  // Use the store's ValidateDependencies via a temporary override
  try {
    const { ValidateDependencies } = await import(
      '../../../bindings/github.com/Aliuyanfeng/happytools/backend/services/makefile/makefileservice.js'
    )
    const cycles = await ValidateDependencies(tempDoc)
    if (cycles && cycles.length > 0) {
      cycleDetected.value = true
      cycleDescription.value = cycles.map((c: string[]) => c.join(' â†?')).join('; ')
    } else {
      cycleDetected.value = false
      cycleDescription.value = ''
    }
  } catch {
    cycleDetected.value = false
    cycleDescription.value = ''
  }
}

// Tab-indent on Enter key in the commands textarea
function handleTextareaKeydown(e: KeyboardEvent) {
  if (e.key === 'Enter') {
    e.preventDefault()
    const textarea = e.target as HTMLTextAreaElement
    const start = textarea.selectionStart
    const end = textarea.selectionEnd
    const value = commandsText.value
    // Insert newline + tab at cursor position
    commandsText.value = value.substring(0, start) + '\n\t' + value.substring(end)
    // Move cursor after the inserted tab
    const newPos = start + 2
    // Use nextTick-equivalent: set selection after DOM update
    setTimeout(() => {
      textarea.selectionStart = newPos
      textarea.selectionEnd = newPos
    }, 0)
  }
}

const rules: Record<string, Rule[]> = {
  name: [
    {
      required: true,
      whitespace: true,
      message: t('makefileEditor.targetNameRequired'),
      trigger: 'blur',
    },
    {
      validator: (_rule, value: string) => {
        const trimmed = value?.trim() ?? ''
        if (!trimmed) return Promise.resolve()
        const others = (props.existingNames ?? []).filter(
          (n) => n !== props.initialValue?.name,
        )
        if (others.includes(trimmed)) {
          return Promise.reject(t('makefileEditor.targetNameDuplicate'))
        }
        return Promise.resolve()
      },
      trigger: 'blur',
    },
  ],
}

async function handleSubmit() {
  try {
    await formRef.value?.validate()
    // Parse commands from textarea: split by newline, strip leading tab
    const commands = commandsText.value
      .split('\n')
      .map((line) => line.replace(/^\t/, ''))
      .filter((line) => line.trim() !== '')

    emit('submit', {
      name: formState.name.trim(),
      isPhony: formState.isPhony,
      deps: [...formState.deps],
      commands,
    })
  } catch {
    // validation failed â€?do nothing
  }
}
</script>

<style scoped>
.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  margin-top: 4px;
}

.commands-hint {
  font-size: 11px;
  color: #8c8c8c;
  margin-top: 4px;
}
</style>
