<template>
  <div class="container">
    <div class="grid" style="gap: var(--sp-6)">
      <!-- Título + badge -->
      <header class="row" style="justify-content: space-between;">
        <h2 style="font-size: var(--fs-xl); font-weight: 700;">Dashboard</h2>
        <span class="badge">tempo real (mock)</span>
      </header>

      <!-- Error States -->
      <div v-if="errorDevices" class="error-card card">
        <p class="error-message">{{ errorDevices }}</p>
        <button class="btn btn--solid" @click="loadDevices">Tentar novamente</button>
      </div>

      <div v-if="errorTelemetry" class="error-card card">
        <p class="error-message">{{ errorTelemetry }}</p>
        <button class="btn btn--solid" @click="loadTelemetry">Tentar novamente</button>
      </div>

      <!-- KPIs -->
      <section v-if="loadingDevices" class="stat-grid">
        <SkeletonCard />
        <SkeletonCard />
        <SkeletonCard />
      </section>
      <section v-else class="stat-grid">
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

      <!-- Telemetry Chart -->
      <ConsumptionChart 
        :labels="labels" 
        :series="series" 
        :loading="loadingTelemetry"
        title="Consumo - Últimos 5 minutos"
      />

      <!-- Telemetry Table -->
      <TelemetryTable 
        :data="telemetryData" 
        :loading="loadingTelemetry"
        @refresh="loadTelemetry"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, ref, computed } from 'vue'
import SkeletonCard from '../components/SkeletonCard.vue'
import ConsumptionChart from '../components/ConsumptionChart.vue'
import TelemetryTable, { type TelemetryData } from '../components/TelemetryTable.vue'
// Uncomment when API is ready:
// import api from '../api/axios'

const series = ref<number[]>([])
const labels = ref<string[]>([])
const currentPower = ref(120)
const activeDevices = 3

const telemetryData = ref<TelemetryData[]>([])
const loadingDevices = ref(false)
const loadingTelemetry = ref(false)
const errorDevices = ref('')
const errorTelemetry = ref('')

const estimatedKwh = computed(() => (currentPower.value / 1000).toFixed(3))

let pollingTimer: number | undefined
let isPolling = false

async function loadDevices() {
  // Prevent parallel requests
  if (loadingDevices.value) return
  
  loadingDevices.value = true
  errorDevices.value = ''
  
  try {
    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 1000))
    // const { data } = await api.get('/devices')
    // Process device data here
  } catch (e: any) {
    errorDevices.value = e?.response?.data?.error || 
      'Erro ao carregar dispositivos. Verifique sua conexão e tente novamente.'
  } finally {
    loadingDevices.value = false
  }
}

async function loadTelemetry() {
  // Prevent parallel requests during polling
  if (loadingTelemetry.value) return
  
  loadingTelemetry.value = true
  errorTelemetry.value = ''
  
  try {
    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 800))
    // const { data } = await api.get('/telemetry')
    
    // Mock data for demonstration
    const noise = (Math.random() - 0.5) * 10
    currentPower.value = Math.max(50, currentPower.value + noise)
    series.value.push(currentPower.value)
    labels.value.push(new Date().toLocaleTimeString('pt-BR'))
    
    // Keep last 300 points (~5 min at 1 sample/sec)
    if (series.value.length > 300) {
      series.value.shift()
      labels.value.shift()
    }
    
    // Update telemetry table
    telemetryData.value = [
      {
        id: 1,
        deviceName: 'Ar Condicionado',
        power: currentPower.value * 0.6,
        timestamp: new Date().toISOString()
      },
      {
        id: 2,
        deviceName: 'Geladeira',
        power: currentPower.value * 0.3,
        timestamp: new Date().toISOString()
      },
      {
        id: 3,
        deviceName: 'Iluminação',
        power: currentPower.value * 0.1,
        timestamp: new Date().toISOString()
      }
    ]
  } catch (e: any) {
    errorTelemetry.value = e?.response?.data?.error || 
      'Erro ao carregar dados de telemetria. Verifique sua conexão e tente novamente.'
  } finally {
    loadingTelemetry.value = false
  }
}

function startPolling() {
  if (isPolling) return
  isPolling = true
  
  pollingTimer = window.setInterval(() => {
    // Only poll if not currently loading (avoid parallel requests)
    if (!loadingTelemetry.value) {
      loadTelemetry()
    }
  }, 2000) // Poll every 2 seconds
}

function stopPolling() {
  isPolling = false
  if (pollingTimer) {
    clearInterval(pollingTimer)
    pollingTimer = undefined
  }
}

onMounted(async () => {
  // Initial load
  await Promise.all([
    loadDevices(),
    loadTelemetry()
  ])
  
  // Start polling for telemetry updates
  startPolling()
})

onUnmounted(() => {
  stopPolling()
})
</script>

<style scoped>
/* sem estilos extras: usamos theme.css (stat-grid, chip-grid, card, etc.) */
.error-card {
  padding: var(--sp-5, 20px);
  background: rgba(255, 99, 132, 0.1);
  border: 1px solid rgba(255, 99, 132, 0.3);
  display: flex;
  flex-direction: column;
  gap: var(--sp-3, 12px);
  align-items: center;
}

.error-message {
  color: var(--warn, #ff6384);
  margin: 0;
  text-align: center;
}
</style>