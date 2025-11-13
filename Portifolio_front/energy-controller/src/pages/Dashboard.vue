<template>
  <div class="container">
    <div class="grid" style="gap: var(--sp-6)">
      <!-- Título + badge + seletor de dispositivo -->
      <header class="row" style="justify-content: space-between; align-items: center;">
        <div>
          <h2 style="font-size: var(--fs-xl); font-weight: 700;">Dashboard</h2>
        </div>
        <div class="row" style="gap: var(--sp-3); align-items: center;">
          <label for="device-select" style="font-size: var(--fs-sm);">Dispositivo:</label>
          <select 
            id="device-select" 
            v-model="selectedDeviceId" 
            @change="onDeviceChange"
            style="padding: var(--sp-2); border-radius: 4px; border: 1px solid #e5e7eb;"
          >
            <option :value="null">Selecione um dispositivo</option>
            <option v-for="device in devices" :key="device.id" :value="device.id">
              {{ device.name }} ({{ device.room }})
            </option>
          </select>
          <span class="badge">tempo real</span>
        </div>
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
            <strong style="font-size: var(--fs-2xl);">{{ devices.length }}</strong>
          </div>
        </article>
      </section>

      <!-- Gráfico de Consumo -->
      <section class="card" style="padding: var(--sp-5);">
        <div class="row" style="justify-content: space-between; margin-bottom: var(--sp-3);">
          <h3 class="text-muted small">Gráfico de Consumo</h3>
        </div>
        <div style="height: 300px;">
          <ConsumptionChart :labels="chartLabels" :series="series" />
        </div>
      </section>

      <!-- Tabela de Telemetria -->
      <section class="card" style="padding: var(--sp-5);">
        <TelemetryTable :data="telemetryData" @refresh="handleRefresh" />
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, ref, computed } from 'vue'
import ConsumptionChart from '../components/ConsumptionChart.vue'
import TelemetryTable from '../components/TelemetryTable.vue'
import type { TelemetryRow } from '../components/TelemetryTable.vue'
import { listDevices, type Device } from '../api/devices'
import { fetchTelemetry } from '../api/telemetry'

// State
const devices = ref<Device[]>([])
const selectedDeviceId = ref<number | null>(null)
const series = ref<number[]>([])
const currentPower = ref(120)
const telemetryData = ref<TelemetryRow[]>([])
const chartLabels = ref<string[]>([])

// Computed
const estimatedKwh = computed(() => (currentPower.value / 1000).toFixed(3))

// Timers
let pollingTimer: number | undefined

// Methods
async function loadDevices() {
  try {
    devices.value = await listDevices()
  } catch (error) {
    console.error('Error loading devices:', error)
    // If API fails, use mock data
    devices.value = [
      { id: 1, name: 'Ar Condicionado', room: 'Sala' },
      { id: 2, name: 'Geladeira', room: 'Cozinha' },
      { id: 3, name: 'Computador', room: 'Escritório' }
    ]
  }
}

async function loadTelemetry() {
  if (!selectedDeviceId.value) return

  try {
    const data = await fetchTelemetry({ 
      device_id: selectedDeviceId.value,
      limit: 100
    })
    
    telemetryData.value = data
    
    // Update chart data
    if (data.length > 0) {
      series.value = data.map(d => d.power_w)
      chartLabels.value = data.map(d => {
        const date = new Date(d.time)
        return `${date.getHours().toString().padStart(2, '0')}:${date.getMinutes().toString().padStart(2, '0')}:${date.getSeconds().toString().padStart(2, '0')}`
      })
      currentPower.value = data[data.length - 1]?.power_w || 120
    }
  } catch (error) {
    console.error('Error loading telemetry:', error)
    // If API fails, use mock data
    generateMockData()
  }
}

function generateMockData() {
  // Generate mock telemetry data
  const now = new Date()
  const mockData: TelemetryRow[] = []
  
  for (let i = 0; i < 30; i++) {
    const time = new Date(now.getTime() - (29 - i) * 1000)
    const power = 100 + Math.random() * 50
    mockData.push({
      time: time.toISOString(),
      power_w: parseFloat(power.toFixed(2)),
      voltage: parseFloat((220 + Math.random() * 10).toFixed(2)),
      current: parseFloat((power / 220).toFixed(2))
    })
  }
  
  telemetryData.value = mockData
  series.value = mockData.map(d => d.power_w)
  chartLabels.value = mockData.map(d => {
    const date = new Date(d.time)
    return `${date.getHours().toString().padStart(2, '0')}:${date.getMinutes().toString().padStart(2, '0')}:${date.getSeconds().toString().padStart(2, '0')}`
  })
  currentPower.value = mockData[mockData.length - 1]?.power_w || 120
}

function startPolling() {
  // Stop previous polling if any
  stopPolling()
  
  // Start polling every 1 second
  pollingTimer = window.setInterval(() => {
    loadTelemetry()
  }, 1000)
}

function stopPolling() {
  if (pollingTimer) {
    clearInterval(pollingTimer)
    pollingTimer = undefined
  }
}

function onDeviceChange() {
  if (selectedDeviceId.value) {
    loadTelemetry()
    startPolling()
  } else {
    stopPolling()
    telemetryData.value = []
    series.value = []
    chartLabels.value = []
  }
}

function handleRefresh() {
  loadTelemetry()
}

// Lifecycle
onMounted(async () => {
  await loadDevices()
  
  // Auto-select first device if available
  if (devices.value.length > 0) {
    selectedDeviceId.value = devices.value[0].id
    await loadTelemetry()
    startPolling()
  }
})

onUnmounted(() => {
  stopPolling()
})
</script>

<style scoped>
/* sem estilos extras: usamos theme.css (stat-grid, chip-grid, card, etc.) */
</style>