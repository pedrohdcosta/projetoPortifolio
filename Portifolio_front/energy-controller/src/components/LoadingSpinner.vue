<template>
  <div 
    class="loading-spinner" 
    :class="{ 'loading-spinner--small': small }"
    role="status" 
    aria-busy="true"
  >
    <div class="spinner" aria-hidden="true" />
    <span class="visually-hidden">Carregando...</span>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  small?: boolean
}>()
</script>

<style scoped>
.loading-spinner {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: var(--sp-4, 16px);
}

.loading-spinner--small {
  padding: var(--sp-2, 8px);
}

.spinner {
  display: inline-block;
  width: 40px;
  height: 40px;
  border: 3px solid rgba(255, 255, 255, 0.2);
  border-top-color: rgba(67, 160, 71, 0.9);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

.loading-spinner--small .spinner {
  width: 20px;
  height: 20px;
  border-width: 2px;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.visually-hidden {
  position: absolute;
  width: 1px;
  height: 1px;
  padding: 0;
  margin: -1px;
  overflow: hidden;
  clip: rect(0, 0, 0, 0);
  white-space: nowrap;
  border-width: 0;
}

@media (prefers-reduced-motion: reduce) {
  .spinner {
    animation-duration: 1.5s;
  }
}
</style>
