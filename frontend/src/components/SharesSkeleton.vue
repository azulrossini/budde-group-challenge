<script setup lang="ts"></script>

<template>
  <div class="skeleton" aria-hidden="true">
    <div v-for="i in 3" :key="i" class="skeleton-row">
      <span class="skeleton-block name" />
      <span class="skeleton-block percentage" />
      <span class="skeleton-block amount" />
    </div>
  </div>
</template>

<style scoped>
.skeleton {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.skeleton-row {
  display: flex;
  gap: 0.75rem;
}

.skeleton-block {
  display: block;
  height: 2.25rem;
  border-radius: 8px;
  background: var(--color-surface-raised);
  animation: pulse 1.2s ease-in-out infinite;
}

.skeleton-block.name {
  flex: 2;
}

.skeleton-block.percentage {
  flex: 1;
}

.skeleton-block.amount {
  flex: 1;
}

@keyframes pulse {
  0%,
  100% {
    opacity: 1;
  }
  50% {
    opacity: 0.45;
  }
}
</style>
