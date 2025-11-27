<template>
  <div class="container">
    <div class="grid" style="gap: var(--sp-6)">
      <!-- Título + badge -->
      <header class="row" style="justify-content: space-between;">
        <h2 style="font-size: var(--fs-xl); font-weight: 700;">Dashboard</h2>
        <span class="badge">tempo real</span>
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
            <strong style="font-size: var(--fs-2xl);">{{ activeDevicesCount }}</strong>
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
        :data="telemetryTableData" 
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
import { listDevices, listTelemetry, type Device, type TelemetryData as ApiTelemetry } from '../api/devices'

const series = ref<number[]>([])
const labels = ref<string[]>([])
const currentPower = ref(0)
const devices = ref<Device[]>([])
const telemetryRaw = ref<ApiTelemetry[]>([])

const loadingDevices = ref(false)
const loadingTelemetry = ref(false)
const errorDevices = ref('')
const errorTelemetry = ref('')

const estimatedKwh = computed(() => (currentPower.value / 1000).toFixed(3))
const activeDevicesCount = computed(() => devices.value.filter(d => d.status === 'online').length)

// Transform API telemetry data to table format - showing ALL readings
const telemetryTableData = computed<TelemetryData[]>(() => {
  // Create device lookup map for O(1) access - O(m)
  const deviceLookup = new Map<number, Device>()
  for (const d of devices.value) {
    deviceLookup.set(d.id, d)
  }
  
  // Map all telemetry readings to table format - O(n)
  return telemetryRaw.value.map(t => {
    const device = deviceLookup.get(t.device_id)
    return {
      id: t.id,
      deviceName: device?.name || `Dispositivo ${t.device_id}`,
      power: t.power,
      timestamp: t.timestamp
    }
  })
})

let pollingTimer: number | undefined
let isPolling = false

async function loadDevices() {
  if (loadingDevices.value) return
  
  loadingDevices.value = true
  errorDevices.value = ''
  
  try {
    devices.value = await listDevices()
  } catch (e: any) {
    errorDevices.value = e?.response?.data?.error || 
      'Erro ao carregar dispositivos. Verifique sua conexão e tente novamente.'
  } finally {
    loadingDevices.value = false
  }
}

async function loadTelemetry() {
  if (loadingTelemetry.value) return
  
  loadingTelemetry.value = true
  errorTelemetry.value = ''
  
  try {
    const data = await listTelemetry(undefined, 100)
    telemetryRaw.value = data
    
    // Calculate current total power
    if (data.length > 0) {
      // Sum the most recent reading from each device
      const deviceLatest = new Map<number, number>()
      for (const t of data) {
        if (!deviceLatest.has(t.device_id)) {
          deviceLatest.set(t.device_id, t.power)
        }
      }
      currentPower.value = Array.from(deviceLatest.values()).reduce((sum, p) => sum + p, 0)
    }
    
    // Update chart series with latest power
    series.value.push(currentPower.value)
    labels.value.push(new Date().toLocaleTimeString('pt-BR'))
    
    // Keep last 300 points (~5 min at 1 sample/sec)
    if (series.value.length > 300) {
      series.value.shift()
      labels.value.shift()
    }
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
    if (!loadingTelemetry.value) {
      loadTelemetry()
    }
  }, 5000) // Poll every 5 seconds
}

function stopPolling() {
  isPolling = false
  if (pollingTimer) {
    clearInterval(pollingTimer)
    pollingTimer = undefined
  }
}

onMounted(async () => {
  await Promise.all([
    loadDevices(),
    loadTelemetry()
  ])
  
  startPolling()
})

onUnmounted(() => {
  stopPolling()
})
</script>

<style scoped>
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