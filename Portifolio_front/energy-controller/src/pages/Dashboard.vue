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

      <!-- Period Filter -->
      <section class="card filter-section">
        <div class="filter-header">
          <h3>Filtro de Período</h3>
          <span class="badge">{{ filteredTelemetry.length }} leituras</span>
        </div>
        <div class="period-selector">
          <button 
            :class="['btn', 'btn--sm', selectedPeriod === '24h' ? 'btn--solid' : 'btn--outline']"
            @click="changePeriod('24h')"
          >Últimas 24h</button>
          <button 
            :class="['btn', 'btn--sm', selectedPeriod === 'month' ? 'btn--solid' : 'btn--outline']"
            @click="changePeriod('month')"
          >Este Mês</button>
          <button 
            :class="['btn', 'btn--sm', selectedPeriod === 'year' ? 'btn--solid' : 'btn--outline']"
            @click="changePeriod('year')"
          >Este Ano</button>
        </div>
      </section>

      <!-- Telemetry Chart -->
      <ConsumptionChart 
        :labels="chartLabels" 
        :series="chartSeries" 
        :loading="loadingTelemetry"
        :title="chartTitle"
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

type PeriodFilter = '24h' | 'month' | 'year'

const currentPower = ref(0)
const devices = ref<Device[]>([])
const telemetryRaw = ref<ApiTelemetry[]>([])
const selectedPeriod = ref<PeriodFilter>('24h')

const loadingDevices = ref(false)
const loadingTelemetry = ref(false)
const errorDevices = ref('')
const errorTelemetry = ref('')

const estimatedKwh = computed(() => (currentPower.value / 1000).toFixed(3))
const activeDevicesCount = computed(() => devices.value.filter(d => d.status === 'online').length)

// Chart title based on selected period
const chartTitle = computed(() => {
  switch (selectedPeriod.value) {
    case '24h': return 'Consumo - Últimas 24 horas'
    case 'month': return 'Consumo - Este Mês'
    case 'year': return 'Consumo - Este Ano'
    default: return 'Consumo'
  }
})

// Filter telemetry data based on selected period
const filteredTelemetry = computed(() => {
  const now = new Date()
  
  return telemetryRaw.value.filter(t => {
    const timestamp = new Date(t.timestamp)
    
    switch (selectedPeriod.value) {
      case '24h': {
        const hoursAgo24 = new Date(now.getTime() - 24 * 60 * 60 * 1000)
        return timestamp >= hoursAgo24
      }
      case 'month': {
        const startOfMonth = new Date(now.getFullYear(), now.getMonth(), 1)
        return timestamp >= startOfMonth
      }
      case 'year': {
        const startOfYear = new Date(now.getFullYear(), 0, 1)
        return timestamp >= startOfYear
      }
      default:
        return true
    }
  })
})

// Create device lookup map for O(1) access
const deviceLookup = computed(() => {
  const lookup = new Map<number, Device>()
  for (const d of devices.value) {
    lookup.set(d.id, d)
  }
  return lookup
})

// Transform filtered telemetry data to table format
const telemetryTableData = computed<TelemetryData[]>(() => {
  return filteredTelemetry.value.map(t => {
    const device = deviceLookup.value.get(t.device_id)
    return {
      id: t.id,
      deviceName: device?.name || `Dispositivo ${t.device_id}`,
      power: t.power,
      timestamp: t.timestamp
    }
  })
})

// Sorted telemetry for chart (ascending by timestamp) - shared to avoid duplicate sorting
const sortedTelemetryForChart = computed(() => {
  return [...filteredTelemetry.value].sort((a, b) => 
    new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime()
  )
})

// Generate chart data from sorted telemetry
const chartSeries = computed(() => {
  return sortedTelemetryForChart.value.map(t => t.power)
})

const chartLabels = computed(() => {
  return sortedTelemetryForChart.value.map(t => {
    const date = new Date(t.timestamp)
    switch (selectedPeriod.value) {
      case '24h':
        return date.toLocaleTimeString('pt-BR', { hour: '2-digit', minute: '2-digit' })
      case 'month':
        return date.toLocaleDateString('pt-BR', { day: '2-digit', month: '2-digit' })
      case 'year':
        return date.toLocaleDateString('pt-BR', { month: 'short' })
      default:
        return date.toLocaleTimeString('pt-BR')
    }
  })
})

let pollingTimer: number | undefined
let isPolling = false

function changePeriod(period: PeriodFilter) {
  selectedPeriod.value = period
}

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
    // Load more data to support filtering by month/year
    const data = await listTelemetry(undefined, 1000)
    telemetryRaw.value = data
    
    // Calculate current total power from most recent readings
    if (data.length > 0) {
      const deviceLatest = new Map<number, number>()
      for (const t of data) {
        if (!deviceLatest.has(t.device_id)) {
          deviceLatest.set(t.device_id, t.power)
        }
      }
      currentPower.value = Array.from(deviceLatest.values()).reduce((sum, p) => sum + p, 0)
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
  }, 30000) // Poll every 30 seconds (less frequent since we're loading more data)
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

.filter-section {
  padding: var(--sp-4, 16px);
}

.filter-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--sp-3, 12px);
}

.filter-header h3 {
  margin: 0;
  font-size: var(--fs-base, 1rem);
  font-weight: 600;
}

.period-selector {
  display: flex;
  gap: var(--sp-2, 8px);
  flex-wrap: wrap;
}

.btn--sm {
  padding: 8px 16px;
  font-size: 0.875rem;
}
</style>