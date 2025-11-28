<template>
  <div class="skeleton-card card" role="status" aria-busy="true">
    <div class="skeleton-header">
      <div class="skeleton-line skeleton-line--title" />
      <div class="skeleton-line skeleton-line--subtitle" />
    </div>
    <div class="skeleton-body">
      <div class="skeleton-line skeleton-line--text" />
      <div class="skeleton-line skeleton-line--text" />
      <div class="skeleton-line skeleton-line--text skeleton-line--short" />
    </div>
    <span class="visually-hidden">Carregando conteúdo...</span>
  </div>
</template>

<script setup lang="ts">
// No props needed for basic skeleton
</script>

<style scoped>
.skeleton-card {
  padding: var(--sp-5, 20px);
  display: flex;
  flex-direction: column;
  gap: var(--sp-4, 16px);
}

.skeleton-header {
  display: flex;
  flex-direction: column;
  gap: var(--sp-2, 8px);
}

.skeleton-body {
  display: flex;
  flex-direction: column;
  gap: var(--sp-3, 12px);
}

.skeleton-line {
  height: 16px;
  background: linear-gradient(
    90deg,
    rgba(255, 255, 255, 0.05) 0%,
    rgba(255, 255, 255, 0.1) 50%,
    rgba(255, 255, 255, 0.05) 100%
  );
  background-size: 200% 100%;
  border-radius: 4px;
  animation: pulse 1.5s ease-in-out infinite;
}

.skeleton-line--title {
  height: 24px;
  width: 60%;
}

.skeleton-line--subtitle {
  height: 14px;
  width: 40%;
}

.skeleton-line--text {
  width: 100%;
}

.skeleton-line--short {
  width: 80%;
}

@keyframes pulse {
  0% {
    background-position: 200% 0;
  }
  100% {
    background-position: -200% 0;
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
  .skeleton-line {
    animation: none;
    opacity: 0.5;
  }
}
</style>
