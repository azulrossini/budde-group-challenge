<script setup lang="ts">
import { NAlert, NButton } from 'naive-ui'
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
  <NAlert v-if="error" type="error" :title="error.message" :show-icon="true">
    <ul v-if="error.details?.length" class="details">
      <li v-for="detail in error.details" :key="detail.field">{{ detail.message }}</li>
    </ul>
    <p v-if="error.requestId" class="request-id">Request ID: {{ error.requestId }}</p>
    <NButton v-if="retryable" size="small" class="retry" @click="$emit('retry')">{{ MSG_RETRY }}</NButton>
  </NAlert>
</template>

<style scoped>
.details {
  margin: 0.4rem 0;
  padding-left: 1.1rem;
}

.request-id {
  margin: 0.2rem 0 0.5rem;
  font-size: 0.8rem;
  color: var(--color-text-muted);
}

.retry {
  margin-top: 0.2rem;
}
</style>
