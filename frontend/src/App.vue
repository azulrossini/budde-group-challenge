<script setup lang="ts">
import { onMounted } from 'vue'
import ErrorBanner from './components/ErrorBanner.vue'
import LoadingSpinner from './components/LoadingSpinner.vue'
import SharesSkeleton from './components/SharesSkeleton.vue'
import SharesTable from './components/SharesTable.vue'
import TotalSummary from './components/TotalSummary.vue'
import { useShares } from './composables/useShares'
import { LoadStatus } from './types/status'
import { MSG_ADD_PERSON, MSG_SAVE, MSG_SAVED, MSG_SAVING } from './utils/messages'
import { formatCents } from './utils/money'

const {
  status,
  isSaving,
  bill,
  rows,
  rowErrors,
  listError,
  sumMessage,
  totalBasisPoints,
  canSave,
  amounts,
  loadError,
  saveError,
  saveSucceeded,
  load,
  save,
  addPerson,
  removePerson,
  updateName,
  updatePercentage,
} = useShares()

onMounted(load)
</script>

<template>
  <main class="app">
    <h1>Split Bill</h1>

    <SharesSkeleton v-if="status === LoadStatus.Loading" />

    <ErrorBanner v-else-if="status === LoadStatus.Error" :error="loadError" retryable @retry="load" />

    <template v-else-if="bill">
      <p class="bill-summary">
        <span class="bill-description">{{ bill.description }}</span>
        <span class="bill-total">{{ formatCents(bill.totalCents) }}</span>
      </p>

      <SharesTable
        :rows="rows"
        :row-errors="rowErrors"
        :amounts="amounts"
        :disabled="isSaving"
        @update:name="updateName"
        @update:percentage="updatePercentage"
        @remove="removePerson"
      />

      <button type="button" class="add-person" :disabled="isSaving" @click="addPerson">
        {{ MSG_ADD_PERSON }}
      </button>

      <TotalSummary :total-basis-points="totalBasisPoints" :warning="listError ?? sumMessage" />

      <ErrorBanner :error="saveError" />
      <p v-if="saveSucceeded" class="save-success">{{ MSG_SAVED }}</p>

      <button type="button" class="save-button" :disabled="!canSave" @click="save">
        <LoadingSpinner v-if="isSaving" />
        {{ isSaving ? MSG_SAVING : MSG_SAVE }}
      </button>
    </template>
  </main>
</template>

<style scoped>
.app {
  max-width: 640px;
  margin: 0 auto;
  padding: 2.5rem 1.5rem 4rem;
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

h1 {
  margin: 0;
  font-size: 1.5rem;
  font-weight: 700;
}

.bill-summary {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
  margin: 0;
  padding: 0.9rem 1rem;
  border-radius: 8px;
  border: 1px solid var(--color-border);
  background: var(--color-surface);
}

.bill-description {
  color: var(--color-text);
  font-weight: 600;
}

.bill-total {
  color: var(--color-text-muted);
}

.add-person {
  align-self: flex-start;
  padding: 0.5rem 1rem;
  border-radius: 6px;
  border: 1px solid var(--color-border);
  background: transparent;
  color: var(--color-text);
  cursor: pointer;
}

.add-person:hover:not(:disabled) {
  border-color: var(--color-accent);
  color: var(--color-accent);
}

.add-person:disabled {
  opacity: 0.5;
  cursor: default;
}

.save-button {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  align-self: flex-end;
  padding: 0.6rem 1.4rem;
  border-radius: 6px;
  border: none;
  background: var(--color-accent);
  color: var(--color-bg);
  font-weight: 600;
  cursor: pointer;
}

.save-button:hover:not(:disabled) {
  background: var(--color-accent-hover);
}

.save-button:disabled {
  opacity: 0.5;
  cursor: default;
  background: var(--color-accent-muted);
  color: var(--color-text-muted);
}

.save-success {
  align-self: flex-end;
  margin: 0;
  color: var(--color-success);
  font-weight: 600;
}
</style>
