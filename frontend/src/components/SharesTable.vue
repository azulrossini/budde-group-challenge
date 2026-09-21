<script setup lang="ts">
import { NTable } from 'naive-ui'
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
  <NTable class="shares-table" :single-line="false">
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
  </NTable>
</template>

<style scoped>
.shares-table :deep(td) {
  vertical-align: top;
}
</style>
