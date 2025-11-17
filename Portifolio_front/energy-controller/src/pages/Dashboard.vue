<template>
  <div class="container">
    <div class="grid" style="gap: var(--sp-6)">
      <!-- Título + badge -->
      <header class="row" style="justify-content: space-between;">
        <h2 style="font-size: var(--fs-xl); font-weight: 700;">Dashboard</h2>
        <span class="badge">tempo real</span>
      </header>

      <!-- Error banner -->
      <div v-if="errorDevices || errorTelemetry" class="error-banner">
        <div class="error-content">
          <span class="error-icon">⚠️</span>
          <div>
            <strong>Erro ao carregar dados</strong>
            <p>{{ errorDevices || errorTelemetry }}</p>
          </div>
        </div>
        <button class="btn btn--outline" @click="retryFetch">
          Tentar novamente
        </button>
      </div>

      <!-- Device selector -->
      <section class="card" style="padding: var(--sp-4);">
        <label for="device-select" class="form-label">Selecionar Dispositivo</label>
        <select 
          id="device-select"
          v-model="selectedDeviceId" 
          class="form-select"
          :disabled="loadingDevices"
        >
          <option :value="null">Selecione um dispositivo...</option>
          <option v-for="device in devices" :key="device.id" :value="device.id">
            {{ device.name }} ({{ device.status }})
          </option>
        </select>
      </section>

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

      <!-- Consumption Chart -->
      <section class="card" style="padding: var(--sp-5);">
        <div class="row" style="justify-content: space-between; margin-bottom: var(--sp-3);">
          <h3 class="text-muted small">Consumo ao Longo do Tempo</h3>
        </div>
        <ConsumptionChart :data="chartData" :loading="loadingTelemetry" />
      </section>

      <!-- Telemetry Table -->
      <section class="card" style="padding: var(--sp-5);">
        <TelemetryTable 
          :data="telemetryData" 
          :loading="loadingTelemetry"
          @refresh="handleRefresh"
        />
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, ref, computed, watch } from 'vue'
import { listDevices, fetchTelemetry, type Device, type TelemetryData } from '../api/devices'
import ConsumptionChart from '../components/ConsumptionChart.vue'
import TelemetryTable from '../components/TelemetryTable.vue'

// State
const devices = ref<Device[]>([])
const selectedDeviceId = ref<number | null>(null)
const telemetryData = ref<TelemetryData[]>([])

// Loading states
const loadingDevices = ref(false)
const loadingTelemetry = ref(false)

// Error states
const errorDevices = ref<string | null>(null)
const errorTelemetry = ref<string | null>(null)

// Computed
const currentPower = computed(() => {
  if (telemetryData.value.length > 0) {
    return telemetryData.value[telemetryData.value.length - 1].power
  }
  return 0
})

const estimatedKwh = computed(() => (currentPower.value / 1000).toFixed(3))

const activeDevices = computed(() => 
  devices.value.filter(d => d.status === 'online').length
)

const chartData = computed(() => 
  telemetryData.value.map(d => ({
    timestamp: d.timestamp,
    power: d.power,
  }))
)

// Functions
async function loadDevices() {
  if (loadingDevices.value) return
  
  loadingDevices.value = true
  errorDevices.value = null
  
  try {
    devices.value = await listDevices()
    // Auto-select first device if available
    if (devices.value.length > 0 && !selectedDeviceId.value) {
      selectedDeviceId.value = devices.value[0].id
    }
  } catch (err: any) {
    errorDevices.value = err.response?.data?.error || err.message || 'Erro ao carregar dispositivos'
    console.error('Error loading devices:', err)
  } finally {
    loadingDevices.value = false
  }
}

async function loadTelemetry() {
  if (!selectedDeviceId.value) {
    telemetryData.value = []
    return
  }
  
  // Prevent concurrent requests
  if (loadingTelemetry.value) return
  
  loadingTelemetry.value = true
  errorTelemetry.value = null
  
  try {
    telemetryData.value = await fetchTelemetry(selectedDeviceId.value)
  } catch (err: any) {
    errorTelemetry.value = err.response?.data?.error || err.message || 'Erro ao carregar telemetria'
    console.error('Error loading telemetry:', err)
  } finally {
    loadingTelemetry.value = false
  }
}

function handleRefresh() {
  loadTelemetry()
}

function retryFetch() {
  errorDevices.value = null
  errorTelemetry.value = null
  loadDevices()
  if (selectedDeviceId.value) {
    loadTelemetry()
  }
}

// Watch for device selection changes
watch(selectedDeviceId, (newDeviceId) => {
  if (newDeviceId) {
    loadTelemetry()
  } else {
    telemetryData.value = []
  }
})

// Polling setup (1 second interval, but prevent concurrent requests)
let pollingInterval: number | undefined

onMounted(async () => {
  await loadDevices()
  
  // Setup polling for telemetry data
  pollingInterval = window.setInterval(() => {
    if (selectedDeviceId.value && !loadingTelemetry.value) {
      loadTelemetry()
    }
  }, 1000)
})

onUnmounted(() => {
  if (pollingInterval) {
    clearInterval(pollingInterval)
  }
})
</script>

<style scoped>
.error-banner {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--sp-4, 16px);
  background-color: #920707;
  border: 1px solid #ffffff;
  border-radius: var(--radius, 8px);
  gap: var(--sp-3, 12px);
}

.error-content {
  display: flex;
  align-items: flex-start;
  gap: var(--sp-3, 12px);
}

.error-icon {
  font-size: 24px;
}

.error-content strong {
  color: #000000;
  display: block;
  margin-bottom: 4px;
}

.error-content p {
  color: #000000;
  margin: 0;
  font-size: var(--fs-sm, 14px);
}

.form-label {
  display: block;
  margin-bottom: var(--sp-2, 8px);
  font-weight: 600;
  font-size: var(--fs-sm, 14px);
}

.form-select {
  width: 100%;
  padding: var(--sp-2, 8px) var(--sp-3, 12px);
  border: 1px solid var(--border-color, #e5e7eb);
  border-radius: var(--radius, 8px);
  font-size: var(--fs-base, 16px);
  background-color: white;
  cursor: pointer;
}

.form-select:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
</style>