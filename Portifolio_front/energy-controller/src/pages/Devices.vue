<template>
  <div class="container">
    <div class="col" style="gap: var(--sp-6)">

      <!-- Título -->
      <header class="row" style="justify-content: space-between;">
        <h2 style="font-size: var(--fs-xl); font-weight: 700;">Dispositivos</h2>
        <span class="badge">{{ devices.length }} dispositivos</span>
      </header>

      <!-- Error State -->
      <div v-if="error" class="error-card card">
        <p class="error-message">{{ error }}</p>
        <button class="btn btn--solid" @click="clearError">Fechar</button>
      </div>

      <!-- Formulário -->
      <section class="card" style="padding: var(--sp-5);">
        <form class="row" style="gap: var(--sp-3); align-items: stretch; flex-wrap: wrap;" @submit.prevent="add">
          <input
            v-model.trim="name"
            class="input"
            placeholder="Nome do dispositivo"
            aria-label="Nome do dispositivo"
            required
            autocomplete="off"
            :disabled="loading"
          />
          <input
            v-model.trim="room"
            class="input"
            placeholder="Ambiente (ex: Sala, Quarto)"
            aria-label="Ambiente do dispositivo"
            autocomplete="off"
            :disabled="loading"
          />
          <select v-model="deviceType" class="input" :disabled="loading">
            <option value="smart_plug">Tomada Inteligente</option>
            <option value="sensor">Sensor</option>
            <option value="meter">Medidor</option>
          </select>
          <button class="btn btn--solid" :disabled="!canAdd || loading">
            <span v-if="!loading">Adicionar</span>
            <LoadingSpinner v-else small />
          </button>
        </form>
      </section>

      <!-- Loading state -->
      <section v-if="loadingList" class="card" style="padding: var(--sp-5);">
        <LoadingSpinner />
        <p class="text-muted" style="text-align: center; margin-top: var(--sp-3);">Carregando dispositivos...</p>
      </section>

      <!-- Lista -->
      <section v-else-if="devices.length > 0" class="grid device-grid">
        <article
          v-for="d in devicesWithTelemetry"
          :key="d.id"
          class="card device-card"
        >
          <div class="row" style="justify-content: space-between;">
            <strong style="font-size: var(--fs-lg);">{{ d.name }}</strong>
            <span :class="['badge', d.status === 'online' ? 'badge--success' : '']">
              {{ d.status === 'online' ? '🟢 Online' : '⚫ Offline' }}
            </span>
          </div>

          <div class="row" style="gap: var(--sp-2); flex-wrap: wrap;">
            <span v-if="d.room" class="badge">📍 {{ d.room }}</span>
            <span class="badge">{{ getDeviceTypeLabel(d.type) }}</span>
          </div>

          <!-- Telemetry Info -->
          <div class="telemetry-info">
            <div class="telemetry-power">
              <span class="power-label">Potência Atual</span>
              <span class="power-value" :class="{ 'power-active': d.latestPower !== undefined }">
                {{ d.latestPower !== undefined ? d.latestPower.toFixed(1) : '--' }} W
              </span>
            </div>
            <div v-if="d.last_seen" class="last-seen">
              <span class="text-muted small">Última leitura: {{ formatLastSeen(d.last_seen) }}</span>
            </div>
          </div>

          <div class="row" style="justify-content: space-between; margin-top: auto;">
            <span class="text-muted small">ID: {{ d.id }}</span>
            <div class="row" style="gap: var(--sp-2);">
              <button 
                class="btn btn--outline btn--sm" 
                @click="viewTelemetry(d.id)"
                :disabled="loading"
              >
                📊 Ver dados
              </button>
              <button 
                class="btn btn--outline btn--sm" 
                @click="remove(d.id)"
                :disabled="loading"
              >
                🗑️ Remover
              </button>
            </div>
          </div>
        </article>
      </section>

      <!-- Empty state -->
      <section v-else class="card" style="padding: var(--sp-5);">
        <p class="text-muted">Nenhum dispositivo cadastrado. Adicione seu primeiro dispositivo IoT acima.</p>
      </section>

      <!-- Device Telemetry Modal -->
      <div v-if="selectedDeviceId !== null" class="modal-overlay" @click.self="closeModal">
        <div class="modal-content card">
          <div class="modal-header">
            <h3>Telemetria - {{ getDeviceName(selectedDeviceId) }}</h3>
            <button class="btn btn--outline btn--sm" @click="closeModal">✕</button>
          </div>
          
          <div v-if="loadingTelemetry" class="modal-body">
            <LoadingSpinner />
          </div>
          
          <div v-else class="modal-body">
            <!-- Summary Cards -->
            <div v-if="selectedSummary" class="summary-grid">
              <div class="summary-card">
                <span class="summary-label">Consumo Médio</span>
                <span class="summary-value">{{ selectedSummary.avg_power.toFixed(1) }} W</span>
              </div>
              <div class="summary-card">
                <span class="summary-label">Pico Máximo</span>
                <span class="summary-value">{{ selectedSummary.max_power.toFixed(1) }} W</span>
              </div>
              <div class="summary-card">
                <span class="summary-label">Consumo Mínimo</span>
                <span class="summary-value">{{ selectedSummary.min_power.toFixed(1) }} W</span>
              </div>
              <div class="summary-card">
                <span class="summary-label">Energia Total</span>
                <span class="summary-value">{{ selectedSummary.total_energy.toFixed(3) }} kWh</span>
              </div>
            </div>

            <!-- Period Selector -->
            <div class="period-selector">
              <button 
                :class="['btn', 'btn--sm', selectedPeriod === 'day' ? 'btn--solid' : 'btn--outline']"
                @click="changePeriod('day')"
              >Hoje</button>
              <button 
                :class="['btn', 'btn--sm', selectedPeriod === 'week' ? 'btn--solid' : 'btn--outline']"
                @click="changePeriod('week')"
              >Semana</button>
              <button 
                :class="['btn', 'btn--sm', selectedPeriod === 'month' ? 'btn--solid' : 'btn--outline']"
                @click="changePeriod('month')"
              >Mês</button>
            </div>

            <!-- Telemetry List -->
            <div v-if="selectedTelemetry.length > 0" class="telemetry-list">
              <div class="telemetry-row telemetry-header">
                <span>Horário</span>
                <span>Potência</span>
                <span>Tensão</span>
                <span>Corrente</span>
              </div>
              <div v-for="t in displayedTelemetry" :key="t.id" class="telemetry-row">
                <span>{{ formatTimestamp(t.timestamp) }}</span>
                <span>{{ t.power.toFixed(1) }} W</span>
                <span>{{ t.voltage?.toFixed(1) || '--' }} V</span>
                <span>{{ t.current?.toFixed(2) || '--' }} A</span>
              </div>
            </div>
            <p v-else class="text-muted">Nenhum dado de telemetria disponível.</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted } from 'vue'
import LoadingSpinner from '../components/LoadingSpinner.vue'
import { 
  listDevices, 
  createDevice, 
  deleteDevice, 
  getLatestTelemetry,
  getDeviceTelemetry,
  getDeviceTelemetrySummary,
  type Device,
  type TelemetryData,
  type TelemetrySummary
} from '../api/devices'

// Configuration constants
const TELEMETRY_DISPLAY_LIMIT = 10

interface DeviceWithTelemetry extends Device {
  latestPower?: number;
}

const devices = ref<Device[]>([])
const latestTelemetry = ref<TelemetryData[]>([])
const name = ref('')
const room = ref('')
const deviceType = ref('smart_plug')
const loading = ref(false)
const loadingList = ref(false)
const error = ref('')

// Modal state
const selectedDeviceId = ref<number | null>(null)
const selectedTelemetry = ref<TelemetryData[]>([])
const selectedSummary = ref<TelemetrySummary | null>(null)
const selectedPeriod = ref<'day' | 'week' | 'month'>('day')
const loadingTelemetry = ref(false)

// Computed property for displayed telemetry with configurable limit
const displayedTelemetry = computed(() => selectedTelemetry.value.slice(0, TELEMETRY_DISPLAY_LIMIT))

const canAdd = computed(() => name.value.trim().length > 0)

// Combine devices with their latest telemetry
const devicesWithTelemetry = computed<DeviceWithTelemetry[]>(() => {
  const telemetryMap = new Map<number, number>()
  for (const t of latestTelemetry.value) {
    telemetryMap.set(t.device_id, t.power)
  }
  
  return devices.value.map(d => ({
    ...d,
    latestPower: telemetryMap.get(d.id)
  }))
})

function getDeviceTypeLabel(type?: string): string {
  switch (type) {
    case 'smart_plug': return '🔌 Tomada'
    case 'sensor': return '📡 Sensor'
    case 'meter': return '⚡ Medidor'
    default: return '📱 Dispositivo'
  }
}

function getDeviceName(id: number): string {
  const device = devices.value.find(d => d.id === id)
  return device?.name || `Dispositivo ${id}`
}

function formatLastSeen(timestamp?: string): string {
  if (!timestamp) return '--'
  const date = new Date(timestamp)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  
  if (diff < 60000) return 'agora'
  if (diff < 3600000) return `${Math.floor(diff / 60000)} min atrás`
  if (diff < 86400000) return `${Math.floor(diff / 3600000)} h atrás`
  return date.toLocaleDateString('pt-BR')
}

function formatTimestamp(timestamp: string): string {
  return new Date(timestamp).toLocaleString('pt-BR')
}

async function loadDevices() {
  loadingList.value = true
  error.value = ''
  
  try {
    const [devicesData, telemetryData] = await Promise.all([
      listDevices(),
      getLatestTelemetry()
    ])
    devices.value = devicesData
    latestTelemetry.value = telemetryData
  } catch (e: any) {
    error.value = e?.response?.data?.error || 
      'Erro ao carregar dispositivos. Verifique sua conexão.'
  } finally {
    loadingList.value = false
  }
}

async function add() {
  if (!canAdd.value || loading.value) return
  
  loading.value = true
  error.value = ''
  
  try {
    const newDevice = await createDevice({
      name: name.value.trim(),
      room: room.value.trim() || undefined,
      type: deviceType.value
    })
    
    devices.value.unshift(newDevice)
    name.value = ''
    room.value = ''
    deviceType.value = 'smart_plug'
  } catch (e: any) {
    error.value = e?.response?.data?.error || 
      'Erro ao adicionar dispositivo. Tente novamente.'
  } finally {
    loading.value = false
  }
}

async function remove(id: number) {
  if (loading.value) return
  
  loading.value = true
  error.value = ''
  
  try {
    await deleteDevice(id)
    devices.value = devices.value.filter(x => x.id !== id)
    latestTelemetry.value = latestTelemetry.value.filter(t => t.device_id !== id)
  } catch (e: any) {
    error.value = e?.response?.data?.error || 
      'Erro ao remover dispositivo. Tente novamente.'
  } finally {
    loading.value = false
  }
}

async function viewTelemetry(deviceId: number) {
  selectedDeviceId.value = deviceId
  selectedPeriod.value = 'day'
  await loadDeviceTelemetry(deviceId)
}

async function loadDeviceTelemetry(deviceId: number) {
  loadingTelemetry.value = true
  
  try {
    const [telemetry, summary] = await Promise.all([
      getDeviceTelemetry(deviceId, 50),
      getDeviceTelemetrySummary(deviceId, selectedPeriod.value)
    ])
    selectedTelemetry.value = telemetry
    selectedSummary.value = summary
  } catch (e: any) {
    selectedTelemetry.value = []
    selectedSummary.value = null
  } finally {
    loadingTelemetry.value = false
  }
}

async function changePeriod(period: 'day' | 'week' | 'month') {
  if (selectedDeviceId.value === null) return
  selectedPeriod.value = period
  
  loadingTelemetry.value = true
  try {
    selectedSummary.value = await getDeviceTelemetrySummary(selectedDeviceId.value, period)
  } catch (e) {
    selectedSummary.value = null
  } finally {
    loadingTelemetry.value = false
  }
}

function closeModal() {
  selectedDeviceId.value = null
  selectedTelemetry.value = []
  selectedSummary.value = null
}

function clearError() {
  error.value = ''
}

onMounted(() => {
  loadDevices()
})
</script>

<style scoped>
/* grid fluido p/ cards de dispositivos (desktop de verdade) */
.device-grid {
  gap: var(--sp-4);
  grid-template-columns: repeat(12, 1fr);
}
.device-grid > .card {
  grid-column: span 12;
}
@media (min-width: 900px) {
  .device-grid > .card { grid-column: span 6; }
}
@media (min-width: 1400px) {
  .device-grid > .card { grid-column: span 4; }
}

.device-card {
  padding: var(--sp-5);
  display: grid;
  gap: var(--sp-3);
}

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

.badge--success {
  background: rgba(67, 160, 71, 0.2);
  color: #43a047;
}

select.input {
  min-width: 150px;
}

.btn--sm {
  padding: 6px 12px;
  font-size: 0.85rem;
}

/* Telemetry Info */
.telemetry-info {
  background: rgba(0, 0, 0, 0.2);
  border-radius: 8px;
  padding: var(--sp-3);
  display: flex;
  flex-direction: column;
  gap: var(--sp-2);
}

.telemetry-power {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.power-label {
  color: var(--muted);
  font-size: 0.9rem;
}

.power-value {
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--muted);
}

.power-value.power-active {
  color: #43a047;
}

.last-seen {
  text-align: right;
}

/* Modal */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: var(--sp-4);
}

.modal-content {
  width: 100%;
  max-width: 700px;
  max-height: 90vh;
  overflow-y: auto;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--sp-4);
  border-bottom: 1px solid var(--border);
}

.modal-header h3 {
  margin: 0;
  font-size: var(--fs-lg);
}

.modal-body {
  padding: var(--sp-4);
  display: flex;
  flex-direction: column;
  gap: var(--sp-4);
}

/* Summary Grid */
.summary-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: var(--sp-3);
}

@media (min-width: 600px) {
  .summary-grid {
    grid-template-columns: repeat(4, 1fr);
  }
}

.summary-card {
  background: rgba(0, 0, 0, 0.2);
  border-radius: 8px;
  padding: var(--sp-3);
  text-align: center;
}

.summary-label {
  display: block;
  font-size: 0.8rem;
  color: var(--muted);
  margin-bottom: 4px;
}

.summary-value {
  display: block;
  font-size: 1.2rem;
  font-weight: 600;
}

/* Period Selector */
.period-selector {
  display: flex;
  gap: var(--sp-2);
  justify-content: center;
}

/* Telemetry List */
.telemetry-list {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.telemetry-row {
  display: grid;
  grid-template-columns: 2fr 1fr 1fr 1fr;
  gap: var(--sp-2);
  padding: var(--sp-2);
  background: rgba(0, 0, 0, 0.1);
  border-radius: 4px;
  font-size: 0.85rem;
}

.telemetry-header {
  font-weight: 600;
  background: rgba(0, 0, 0, 0.3);
}
</style>
