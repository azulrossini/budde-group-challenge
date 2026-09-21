<script setup lang="ts">
import type { ApiError } from '../types/errors'
import { MSG_RETRY } from '../utils/messages'

defineProps<{
  error: ApiError | null
  retryable?: boolean
}>()

defineEmits<{
  retry: []
}>()
</script>

<template>
  <div v-if="error" class="error-banner" role="alert">
    <p class="message">{{ error.message }}</p>
    <ul v-if="error.details?.length" class="details">
      <li v-for="detail in error.details" :key="detail.field">{{ detail.message }}</li>
    </ul>
    <p v-if="error.requestId" class="request-id">Request ID: {{ error.requestId }}</p>
    <button v-if="retryable" type="button" class="retry" @click="$emit('retry')">{{ MSG_RETRY }}</button>
  </div>
</template>

<style scoped>
.error-banner {
  display: flex;
  flex-direction: column;
  gap: 0.4rem;
  padding: 0.9rem 1rem;
  border-radius: 8px;
  border: 1px solid var(--color-error);
  background: var(--color-error-muted);
  color: var(--color-text);
}

.message {
  margin: 0;
  font-weight: 600;
  color: var(--color-error);
}

.details {
  margin: 0;
  padding-left: 1.1rem;
  color: var(--color-text);
}

.request-id {
  margin: 0;
  font-size: 0.8rem;
  color: var(--color-text-muted);
}

.retry {
  align-self: flex-start;
  margin-top: 0.2rem;
  padding: 0.4rem 0.9rem;
  border-radius: 6px;
  border: 1px solid var(--color-error);
  background: transparent;
  color: var(--color-error);
  cursor: pointer;
}

.retry:hover {
  background: var(--color-error-muted);
}
</style>
