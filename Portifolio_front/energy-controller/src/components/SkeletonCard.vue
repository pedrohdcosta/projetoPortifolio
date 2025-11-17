<template>
  <div class="skeleton-card" aria-busy="true" aria-label="Loading content">
    <div class="skeleton-header"></div>
    <div class="skeleton-line"></div>
    <div class="skeleton-line skeleton-line--short"></div>
  </div>
</template>

<script setup lang="ts">
// No props needed for basic skeleton
</script>

<style scoped>
.skeleton-card {
  padding: var(--sp-4, 16px);
  background: var(--bg-card, #fff);
  border-radius: var(--radius, 8px);
  border: 1px solid var(--border-color, #e5e7eb);
}

.skeleton-header {
  width: 60%;
  height: 24px;
  background: linear-gradient(90deg, #f0f0f0 25%, #e0e0e0 50%, #f0f0f0 75%);
  background-size: 200% 100%;
  border-radius: 4px;
  margin-bottom: var(--sp-3, 12px);
  animation: skeleton-loading 1.5s ease-in-out infinite;
}

.skeleton-line {
  width: 100%;
  height: 16px;
  background: linear-gradient(90deg, #f0f0f0 25%, #e0e0e0 50%, #f0f0f0 75%);
  background-size: 200% 100%;
  border-radius: 4px;
  margin-bottom: var(--sp-2, 8px);
  animation: skeleton-loading 1.5s ease-in-out infinite;
}

.skeleton-line--short {
  width: 40%;
}

@keyframes skeleton-loading {
  0% {
    background-position: 200% 0;
  }
  100% {
    background-position: -200% 0;
  }
}
</style>
