<script setup lang="ts">
import { formatCents } from '../utils/money'

defineProps<{
  name: string
  percentageInput: string
  amountCents: number | null
  nameError?: string
  percentageError?: string
  disabled: boolean
}>()

const emit = defineEmits<{
  'update:name': [name: string]
  'update:percentage': [value: string]
  remove: []
}>()

function onNameInput(event: Event) {
  emit('update:name', (event.target as HTMLInputElement).value)
}

function onPercentageInput(event: Event) {
  emit('update:percentage', (event.target as HTMLInputElement).value)
}
</script>

<template>
  <tr class="share-row">
    <td>
      <input
        type="text"
        class="field"
        :class="{ invalid: nameError }"
        :value="name"
        :disabled="disabled"
        placeholder="Name"
        @input="onNameInput"
      />
      <p v-if="nameError" class="field-error">{{ nameError }}</p>
    </td>
    <td>
      <input
        type="text"
        inputmode="decimal"
        class="field percentage"
        :class="{ invalid: percentageError }"
        :value="percentageInput"
        :disabled="disabled"
        placeholder="0"
        @input="onPercentageInput"
      />
      <p v-if="percentageError" class="field-error">{{ percentageError }}</p>
    </td>
    <td class="amount">{{ amountCents !== null ? formatCents(amountCents) : '—' }}</td>
    <td class="remove-cell">
      <button
        type="button"
        class="remove"
        :disabled="disabled"
        aria-label="Remove person"
        @click="$emit('remove')"
      >
        ✕
      </button>
    </td>
  </tr>
</template>

<style scoped>
.field {
  width: 100%;
  padding: 0.5rem 0.6rem;
  border-radius: 6px;
  border: 1px solid var(--color-border);
  background: var(--color-surface-raised);
  color: var(--color-text);
}

.field:focus {
  outline: none;
  border-color: var(--color-accent);
  box-shadow: 0 0 0 2px var(--color-accent-muted);
}

.field.invalid {
  border-color: var(--color-error);
}

.field.percentage {
  max-width: 6rem;
}

.field:disabled {
  opacity: 0.6;
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

.remove {
  width: 2rem;
  height: 2rem;
  border-radius: 6px;
  border: 1px solid var(--color-border);
  background: transparent;
  color: var(--color-text-muted);
  cursor: pointer;
}

.remove:hover:not(:disabled) {
  border-color: var(--color-error);
  color: var(--color-error);
}

.remove:disabled {
  opacity: 0.5;
  cursor: default;
}
</style>
