<script setup lang="ts">
import { NButton, NInput } from 'naive-ui'
import { formatCents } from '../utils/money'

defineProps<{
  name: string
  percentageInput: string
  amountCents: number | null
  nameError?: string
  percentageError?: string
  disabled: boolean
}>()

defineEmits<{
  'update:name': [name: string]
  'update:percentage': [value: string]
  remove: []
}>()
</script>

<template>
  <tr class="share-row">
    <td>
      <NInput
        :value="name"
        :status="nameError ? 'error' : undefined"
        :disabled="disabled"
        placeholder="Name"
        @update:value="(value: string) => $emit('update:name', value)"
      />
      <p v-if="nameError" class="field-error">{{ nameError }}</p>
    </td>
    <td>
      <NInput
        :value="percentageInput"
        :status="percentageError ? 'error' : undefined"
        :disabled="disabled"
        :input-props="{ inputmode: 'decimal' }"
        class="percentage-field"
        placeholder="0"
        @update:value="(value: string) => $emit('update:percentage', value)"
      />
      <p v-if="percentageError" class="field-error">{{ percentageError }}</p>
    </td>
    <td class="amount">{{ amountCents !== null ? formatCents(amountCents) : '—' }}</td>
    <td class="remove-cell">
      <NButton circle quaternary :disabled="disabled" aria-label="Remove person" @click="$emit('remove')">
        ✕
      </NButton>
    </td>
  </tr>
</template>

<style scoped>
.percentage-field {
  max-width: 6rem;
}

.field-error {
  margin: 0.25rem 0 0;
  font-size: 0.75rem;
  color: var(--color-error);
}

.amount {
  color: var(--color-text);
  white-space: nowrap;
}

.remove-cell {
  text-align: center;
}
</style>
