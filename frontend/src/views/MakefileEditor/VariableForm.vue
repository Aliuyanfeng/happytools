<template>
  <a-form
    ref="formRef"
    :model="formState"
    :rules="rules"
    layout="vertical"
    @finish="handleSubmit"
  >
    <a-form-item :label="t('makefileEditor.variableName')" name="name">
      <a-input
        v-model:value="formState.name"
        :placeholder="t('makefileEditor.variableName')"
        allow-clear
      />
    </a-form-item>

    <a-form-item :label="t('makefileEditor.operator')" name="operator">
      <a-select v-model:value="formState.operator" style="width: 100%">
        <a-select-option value="=">=</a-select-option>
        <a-select-option value=":=">:=</a-select-option>
        <a-select-option value="?=">?=</a-select-option>
        <a-select-option value="+=">+=</a-select-option>
      </a-select>
    </a-form-item>

    <a-form-item :label="t('makefileEditor.value')" name="value">
      <a-input
        v-model:value="formState.value"
        :placeholder="t('makefileEditor.value')"
        allow-clear
      />
    </a-form-item>

    <div class="form-actions">
      <a-button @click="emit('cancel')">{{ t('common.cancel') }}</a-button>
      <a-button type="primary" html-type="submit">{{ t('common.save') }}</a-button>
    </div>
  </a-form>
</template>

<script setup lang="ts">
import { reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import type { FormInstance, Rule } from 'ant-design-vue/es/form'
import type { Variable } from '../../../bindings/github.com/Aliuyanfeng/happytools/backend/services/makefile/models.js'

const props = defineProps<{
  initialValue?: Variable
  existingNames?: string[]
}>()

const emit = defineEmits<{
  (e: 'submit', variable: Variable): void
  (e: 'cancel'): void
}>()

const { t } = useI18n()
const formRef = ref<FormInstance>()

const formState = reactive<Variable>({
  name: props.initialValue?.name ?? '',
  operator: props.initialValue?.operator ?? '=',
  value: props.initialValue?.value ?? '',
})

watch(
  () => props.initialValue,
  (val) => {
    formState.name = val?.name ?? ''
    formState.operator = val?.operator ?? '='
    formState.value = val?.value ?? ''
  },
)

const rules: Record<string, Rule[]> = {
  name: [
    {
      required: true,
      whitespace: true,
      message: t('makefileEditor.variableNameRequired'),
      trigger: 'blur',
    },
    {
      validator: (_rule, value: string) => {
        const trimmed = value?.trim() ?? ''
        if (!trimmed) return Promise.resolve()
        // Skip duplicate check for the original name when editing
        const others = (props.existingNames ?? []).filter(
          (n) => n !== props.initialValue?.name,
        )
        if (others.includes(trimmed)) {
          return Promise.reject(t('makefileEditor.variableNameDuplicate'))
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
    emit('submit', {
      name: formState.name.trim(),
      operator: formState.operator,
      value: formState.value,
    })
  } catch {
    // validation failed — do nothing
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
</style>
