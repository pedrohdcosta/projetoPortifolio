<template>
  <div class="grid" style="padding:16px">
    <h2 class="row" style="justify-content:space-between">
      <span>Dashboard</span>
      <span class="badge">tempo real (mock)</span>
    </h2>


    <section class="grid" style="grid-template-columns: repeat(auto-fit, minmax(220px,1fr));">
      <div class="card" style="padding:16px">
        <h3 class="text-muted small">Potência Atual</h3>
        <div style="font-size:1.6rem;font-weight:700">{{ currentPower.toFixed(1) }} W</div>
      </div>
      <div class="card" style="padding:16px">
        <h3 class="text-muted small">Consumo Estimado/h</h3>
        <div style="font-size:1.6rem;font-weight:700">{{ (currentPower / 1000).toFixed(3) }} kWh</div>
      </div>
      <div class="card" style="padding:16px">
        <h3 class="text-muted small">Dispositivos Ativos</h3>
        <div style="font-size:1.6rem;font-weight:700">{{ activeDevices }}</div>
      </div>
    </section>


    <section class="card" style="padding:16px">
      <h3 class="text-muted small">Últimos 5 minutos</h3>
      <ul class="series">
        <li v-for="(p, i) in series" :key="i">{{ p.toFixed(1) }}</li>
      </ul>
      <p class="text-muted small">(lista simples no MVP; substitua por gráfico depois)</p>
    </section>
  </div>
</template>
<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
const series = ref<number[]>([])
const currentPower = ref(120)
const activeDevices = 3
let t: number | undefined
onMounted(() => { t = window.setInterval(() => { const noise = (Math.random() - 0.5) * 10; currentPower.value = Math.max(50, currentPower.value + noise); series.value.push(currentPower.value); if (series.value.length > 300) series.value.shift() }, 1000) })
onUnmounted(() => { if (t) clearInterval(t) })
</script>
<style scoped>
.series {
  display: grid;
  grid-template-columns: repeat(6, 1fr);
  gap: 8px;
  list-style: none;
  padding: 0;
  margin: 12px 0 0
}

.series li {
  background: #0f1731;
  padding: 8px;
  border: 1px solid var(--border);
  border-radius: 10px;
  text-align: center
}
</style>