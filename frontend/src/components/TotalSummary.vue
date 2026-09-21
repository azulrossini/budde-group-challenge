<script setup lang="ts">
import { computed } from 'vue'
import { FULL_PERCENT_BASIS_POINTS } from '../utils/constants'
import { formatPercentage } from '../utils/money'

const props = defineProps<{
  totalBasisPoints: number
  warning: string | null
}>()

const isValid = computed(() => props.totalBasisPoints === FULL_PERCENT_BASIS_POINTS)
</script>

<template>
  <div class="total-summary" :class="{ invalid: !isValid }">
    <span class="total-value">
      {{ formatPercentage(totalBasisPoints) }} / {{ formatPercentage(FULL_PERCENT_BASIS_POINTS) }}
    </span>
    <p v-if="warning" class="warning">{{ warning }}</p>
  </div>
</template>

<style scoped>
.total-summary {
  display: flex;
  flex-wrap: wrap;
  align-items: baseline;
  gap: 0.6rem;
  padding: 0.9rem 1rem;
  border-radius: 8px;
  border: 1px solid var(--color-border);
  background: var(--color-surface);
}

.total-value {
  font-size: 1.05rem;
  font-weight: 600;
  color: var(--color-text);
}

.total-summary.invalid .total-value {
  color: var(--color-error);
}

.warning {
  margin: 0;
  font-size: 0.85rem;
  color: var(--color-error);
}
</style>
