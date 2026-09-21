<script setup lang="ts">
import type { RowError, ShareRowData } from '../composables/useShares'
import ShareRow from './ShareRow.vue'

defineProps<{
  rows: ShareRowData[]
  rowErrors: Record<string, RowError>
  amounts: Record<string, number>
  disabled: boolean
}>()

defineEmits<{
  'update:name': [key: string, name: string]
  'update:percentage': [key: string, value: string]
  remove: [key: string]
}>()
</script>

<template>
  <table class="shares-table">
    <thead>
      <tr>
        <th>Name</th>
        <th>Percentage</th>
        <th>Amount</th>
        <th></th>
      </tr>
    </thead>
    <tbody>
      <ShareRow
        v-for="row in rows"
        :key="row.key"
        :name="row.name"
        :percentage-input="row.percentageInput"
        :amount-cents="amounts[row.key] ?? null"
        :name-error="rowErrors[row.key]?.name"
        :percentage-error="rowErrors[row.key]?.percentage"
        :disabled="disabled"
        @update:name="(name) => $emit('update:name', row.key, name)"
        @update:percentage="(value) => $emit('update:percentage', row.key, value)"
        @remove="$emit('remove', row.key)"
      />
    </tbody>
  </table>
</template>

<style scoped>
.shares-table {
  width: 100%;
  border-collapse: collapse;
  border: 1px solid var(--color-border);
  border-radius: 8px;
  overflow: hidden;
  background: var(--color-surface);
}

th {
  text-align: left;
  padding: 0.6rem 0.75rem;
  font-size: 0.8rem;
  font-weight: 600;
  color: var(--color-text-muted);
  border-bottom: 1px solid var(--color-border);
}

:deep(td) {
  padding: 0.6rem 0.75rem;
  border-bottom: 1px solid var(--color-border);
  vertical-align: top;
}

:deep(tr:last-child td) {
  border-bottom: none;
}

:deep(tr:hover td) {
  background: var(--color-surface-raised);
}
</style>
