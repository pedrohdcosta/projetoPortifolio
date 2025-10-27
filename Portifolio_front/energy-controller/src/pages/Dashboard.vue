<template>
  <div class="container">
    <div class="grid" style="gap: var(--sp-6)">
      <!-- Título + badge -->
      <header class="row" style="justify-content: space-between;">
        <h2 style="font-size: var(--fs-xl); font-weight: 700;">Dashboard</h2>
        <span class="badge">tempo real (mock)</span>
      </header>

      <!-- KPIs -->
      <section class="stat-grid">
        <article class="card stat stat--1">
          <div class="col">
            <h3 class="text-muted small">Potência Atual</h3>
            <strong style="font-size: var(--fs-2xl);">{{ currentPower.toFixed(1) }} W</strong>
          </div>
        </article>

        <article class="card stat stat--2">
          <div class="col">
            <h3 class="text-muted small">Consumo Estimado/h</h3>
            <strong style="font-size: var(--fs-2xl);">{{ estimatedKwh }} kWh</strong>
          </div>
        </article>

        <article class="card stat stat--3">
          <div class="col">
            <h3 class="text-muted small">Dispositivos Ativos</h3>
            <strong style="font-size: var(--fs-2xl);">{{ activeDevices }}</strong>
          </div>
        </article>
      </section>

      <!-- Últimos 5 minutos -->
      <section class="card" style="padding: var(--sp-5);">
        <div class="row" style="justify-content: space-between; margin-bottom: var(--sp-3);">
          <h3 class="text-muted small">Últimos 5 minutos</h3>
          <button class="btn btn--outline">Exportar</button>
        </div>

        <div class="chip-grid">
          <span v-for="(p, i) in series" :key="i" class="chip">{{ p.toFixed(1) }}</span>
        </div>

        <p class="text-muted small" style="margin-top: var(--sp-3);">
          (MVP com lista; substitua por gráfico depois)
        </p>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, ref, computed } from 'vue'

const series = ref<number[]>([])
const currentPower = ref(120)
const activeDevices = 3

const estimatedKwh = computed(() => (currentPower.value / 1000).toFixed(3))

let t: number | undefined
onMounted(() => {
  t = window.setInterval(() => {
    const noise = (Math.random() - 0.5) * 10
    currentPower.value = Math.max(50, currentPower.value + noise)
    series.value.push(currentPower.value)
    // mantém ~5 min se 1 amostra/s => 300 pontos
    if (series.value.length > 300) series.value.shift()
  }, 1000)
})
onUnmounted(() => { if (t) clearInterval(t) })
</script>

<style scoped>
/* sem estilos extras: usamos theme.css (stat-grid, chip-grid, card, etc.) */
</style>