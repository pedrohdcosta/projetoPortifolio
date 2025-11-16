<template>
  <div class="spinner-container" role="status" aria-busy="true" aria-live="polite">
    <div class="spinner"></div>
    <span v-if="label" class="spinner-label">{{ label }}</span>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  label?: string;
}>();
</script>

<style scoped>
.spinner-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--sp-2, 8px);
}

.spinner {
  width: 40px;
  height: 40px;
  border: 4px solid rgba(0, 0, 0, 0.1);
  border-left-color: var(--color-primary, #3b82f6);
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

.spinner-label {
  font-size: var(--fs-sm, 14px);
  color: var(--text-muted, #6b7280);
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
