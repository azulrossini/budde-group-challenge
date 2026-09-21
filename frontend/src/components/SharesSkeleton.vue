<script setup lang="ts">
import { NSkeleton } from 'naive-ui'
</script>

<template>
  <div class="skeleton" aria-hidden="true">
    <NSkeleton height="40px" :repeat="3" animated />
  </div>
</template>

<style scoped>
.skeleton :deep(.n-skeleton) {
  margin-bottom: 0.75rem;
  border-radius: 8px;
}

.skeleton :deep(.n-skeleton:last-child) {
  margin-bottom: 0;
}
</style>
