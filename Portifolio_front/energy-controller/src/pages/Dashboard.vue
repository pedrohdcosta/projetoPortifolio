<template>
  <div class="p">
    <h2>Dashboard</h2>
    <div class="cards">
      <div class="card">
        <h3>Potência Atual</h3>
        <p>{{ currentPower.toFixed(1) }} W</p>
      </div>
      <div class="card">
        <h3>Consumo Estimado/h</h3>
        <p>{{ (currentPower / 1000).toFixed(3) }} kWh</p>
      </div>
    </div>
    <div>
      <h3>Últimos 5 minutos</h3>
      <ul class="series">
        <li v-for="(p, i) in series" :key="i">{{ p.toFixed(1) }}</li>
      </ul>
      <small>(lista simples para o MVP; troca por gráfico depois)</small>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, ref } from "vue";
const series = ref<number[]>([]);
const currentPower = ref(120);
let t: number | undefined;

onMounted(() => {
  t = window.setInterval(() => {
    const noise = (Math.random() - 0.5) * 10;
    currentPower.value = Math.max(50, currentPower.value + noise);
    series.value.push(currentPower.value);
    if (series.value.length > 300) series.value.shift();
  }, 1000);
});

onUnmounted(() => {
  if (t) clearInterval(t);
});
</script>
<style>
.p {
  padding: 16px;
}
.cards {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
}
.card {
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 12px;
  min-width: 180px;
}
.series {
  display: grid;
  grid-template-columns: repeat(6, 1fr);
  gap: 8px;
  list-style: none;
  padding: 0;
}
.series li {
  background: #f3f4f6;
  padding: 6px;
  border-radius: 6px;
  text-align: center;
}
</style>
