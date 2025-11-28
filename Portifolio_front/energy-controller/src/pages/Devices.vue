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
            <div class="row" style="gap: var(--sp-2);">
              <span :class="['badge', d.status === 'online' ? 'badge--success' : 'badge--offline']">
                {{ d.status === 'online' ? '🟢 Online' : '⚫ Offline' }}
              </span>
            </div>
          </div>

          <div class="row" style="gap: var(--sp-2); flex-wrap: wrap;">
            <span v-if="d.room" class="badge">📍 {{ d.room }}</span>
            <span class="badge">{{ getDeviceTypeLabel(d.type) }}</span>
          </div>

          <!-- Status Toggle -->
          <div class="power-toggle-section">
            <button 
              :class="['power-toggle-btn', d.status === 'online' ? 'power-on' : 'power-off']"
              @click="toggle(d.id)"
              :disabled="toggling === d.id"
              :title="d.status === 'online' ? 'Desativar dispositivo' : 'Ativar dispositivo'"
            >
              <span v-if="toggling === d.id" class="toggle-loading">⏳</span>
              <span v-else class="toggle-icon">{{ d.status === 'online' ? '🔴' : '🟢' }}</span>
              <span class="toggle-text">{{ d.status === 'online' ? 'Desativar' : 'Ativar' }}</span>
            </button>
          </div>

          <!-- Telemetry Info -->
          <div class="telemetry-info">
            <div class="telemetry-power">
              <span class="power-label">Potência Atual</span>
              <div style="display: flex; gap: 12px; align-items: center;">
                <span class="power-value" :class="{ 'power-active': d.latestPower !== undefined && d.status === 'online' }">
                  {{ d.latestPower !== undefined ? d.latestPower.toFixed(1) : '--' }} W
                </span>
                <span :class="['status-badge', consumptionStatusClass(d.latestPower, d.id)]">
                  {{ consumptionStatusLabel(d.latestPower, d.id) }}
                </span>
              </div>
            </div>
            <div v-if="d.last_seen" class="last-seen">
              <span class="text-muted small">Última leitura: {{ formatLastSeen(d.last_seen) }}</span>
            </div>
          </div>

          <div class="row" style="justify-content: space-between; margin-top: auto;">
            <span class="text-muted small">ID: {{ d.id }}</span>
            <div class="row" style="gap: var(--sp-2); flex-wrap: wrap;">
              <button 
                class="btn btn--outline btn--sm btn--simulate" 
                @click="generateReading(d.id)"
                :disabled="!isSimulationEnabled(d.id) || loading || simulating === d.id"
                title="Gerar leitura simulada"
              >
                <span v-if="simulating === d.id">⏳</span>
                <span v-else>⚡ Simular</span>
              </button>
              <button
                class="btn btn--outline btn--sm"
                :class="{ 'btn--solid': isSimulationEnabled(d.id) }"
                @click.prevent="toggleSimulation(d.id)"
                title="Alternar modo: simulador / API"
              >
                {{ isSimulationEnabled(d.id) ? 'Modo: Simulado' : 'Modo: API' }}
              </button>
              <button
                class="btn btn--outline btn--sm"
                @click.prevent="openTapoConfig(d)"
                title="Configurar credenciais TAPO"
              >
                ⚙️ Config TAPO
              </button>
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

            <!-- Simulator Section -->
            <div class="simulator-section">
              <h4>🔌 Simulador de Tomada</h4>
              <p class="text-muted small">Gere dados simulados para testar o dispositivo sem hardware real.</p>
              <div class="simulator-controls">
                <button 
                  class="btn btn--outline btn--sm"
                  @click="generateBulkReadings"
                  :disabled="generatingBulk"
                >
                  <span v-if="generatingBulk">⏳ Gerando...</span>
                  <span v-else>📈 Gerar histórico (24h)</span>
                </button>
                <button 
                  class="btn btn--solid btn--sm"
                  @click="generateSingleReading"
                  :disabled="generatingSingle"
                >
                  <span v-if="generatingSingle">⏳</span>
                  <span v-else>⚡ Nova leitura</span>
                </button>
              </div>
              <p v-if="simulatorMessage" class="simulator-message text-muted small">{{ simulatorMessage }}</p>
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

    <!-- TAPO Config Modal -->
    <div v-if="editingTapoDeviceId !== null" class="modal-overlay" @click.self="cancelTapoConfig">
      <div class="modal-content card">
        <div class="modal-header">
          <h3>Configurar TAPO - Dispositivo {{ editingTapoDeviceId }}</h3>
          <button class="btn btn--outline btn--sm" @click="cancelTapoConfig">✕</button>
        </div>
        <div class="modal-body">
          <div class="row" style="gap:12px; flex-wrap:wrap;">
            <label class="text-muted">IP</label>
            <input class="input" v-model="tapoIp" placeholder="192.168.1.10" />
            <label class="text-muted">Username</label>
            <input class="input" v-model="tapoUser" placeholder="tapo email" />
            <label class="text-muted">Password</label>
            <input class="input" v-model="tapoPassword" type="password" placeholder="password" />
          </div>
          <p class="text-muted small">Atenção: as credenciais são salvas em `device.metadata` como JSON. Para produção considere um armazenamento seguro.</p>
          <div style="display:flex; gap:8px; margin-top:12px;">
            <button class="btn btn--solid" @click="saveTapoConfig" :disabled="savingTapo">{{ savingTapo ? 'Salvando...' : 'Salvar' }}</button>
            <button class="btn btn--outline" @click="cancelTapoConfig">Cancelar</button>
          </div>
          <div v-if="tapoError" style="color:#c92a2a; margin-top:8px">{{ tapoError }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted, reactive } from 'vue'
import { getThresholds } from '../utils/thresholds'
import LoadingSpinner from '../components/LoadingSpinner.vue'
import { 
  listDevices, 
  createDevice, 
  deleteDevice,
  updateDevice,
  toggleDevice,
  getLatestTelemetry,
  getDeviceTelemetry,
  getDeviceTelemetrySummary,
  simulateTelemetry,
  simulateBulkTelemetry,
  type Device,
  type TelemetryData,
  type TelemetrySummary
} from '../api/devices'
import { getAllSimulationModes, setSimulationMode, toggleSimulationMode } from '../utils/simulationMode'

// Configuration constants
const TELEMETRY_DISPLAY_LIMIT = 10

// Health thresholds (Watts) - can be tuned
const WARNING_THRESHOLD = 100
const DANGER_THRESHOLD = 500

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

// Simulator state
const simulating = ref<number | null>(null) // deviceId being simulated
const generatingSingle = ref(false)
const generatingBulk = ref(false)
const simulatorMessage = ref('')
const toggling = ref<number | null>(null) // deviceId being toggled
// TAPO config modal state
const editingTapoDeviceId = ref<number | null>(null)
const tapoIp = ref('')
const tapoUser = ref('')
const tapoPassword = ref('')
const savingTapo = ref(false)
const tapoError = ref('')

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

// local map of simulation modes for UI responsiveness
const simulationModes = reactive<Record<number, boolean>>({})

function isSimulationEnabled(deviceId: number): boolean {
  return !!simulationModes[deviceId]
}

function setSimulationForDevice(deviceId: number, enabled: boolean) {
  simulationModes[deviceId] = !!enabled
  setSimulationMode(deviceId, !!enabled)
}

function toggleSimulation(deviceId: number) {
  const newVal = toggleSimulationMode(deviceId)
  simulationModes[deviceId] = !!newVal
}

function consumptionStatusLevel(power?: number, deviceId?: number): 'ok' | 'warning' | 'danger' | 'unknown' {
  if (power === undefined || power === null) return 'unknown'
  const p = Number(power)
  if (Number.isNaN(p)) return 'unknown'

  // check per-device thresholds first
  if (deviceId !== undefined && deviceId !== null) {
    const t = getThresholds(deviceId)
    if (t) {
      if (p >= t.danger) return 'danger'
      if (p >= t.warning) return 'warning'
      return 'ok'
    }
  }

  // fallback to page-level defaults
  if (p >= DANGER_THRESHOLD) return 'danger'
  if (p >= WARNING_THRESHOLD) return 'warning'
  return 'ok'
}

function consumptionStatusLabel(power?: number, deviceId?: number): string {
  const lvl = consumptionStatusLevel(power, deviceId)
  if (lvl === 'ok') return 'OK'
  if (lvl === 'warning') return 'Atenção'
  if (lvl === 'danger') return 'Perigo'
  return '--'
}

function consumptionStatusClass(power?: number, deviceId?: number): string {
  return consumptionStatusLevel(power, deviceId)
}

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
    // load simulation modes
    const modes = getAllSimulationModes()
    for (const d of devicesData) {
      simulationModes[d.id] = !!modes[String(d.id)]
    }
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

async function toggle(deviceId: number) {
  if (toggling.value !== null) return
  
  toggling.value = deviceId
  error.value = ''
  
  try {
    const updatedDevice = await toggleDevice(deviceId)
    
    // Update device in the list
    const index = devices.value.findIndex(d => d.id === deviceId)
    if (index >= 0) {
      devices.value[index] = updatedDevice
    }
  } catch (e: any) {
    error.value = e?.response?.data?.error || 
      'Erro ao alternar estado do dispositivo. Tente novamente.'
  } finally {
    toggling.value = null
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

// Simulator functions

async function generateReading(deviceId: number) {
  simulating.value = deviceId
  
  try {
    const result = await simulateTelemetry(deviceId, {
      base_power: 150,
      variation: 0.2,
      base_voltage: 220
    })
    
    // Update the latest telemetry for this device
    if (result.power !== undefined) {
      const existing = latestTelemetry.value.findIndex(t => t.device_id === deviceId)
      const newReading: TelemetryData = {
        id: result.id || 0,
        device_id: deviceId,
        power: result.power,
        voltage: result.voltage,
        current: result.current,
        timestamp: result.timestamp || new Date().toISOString()
      }
      
      if (existing >= 0) {
        latestTelemetry.value[existing] = newReading
      } else {
        latestTelemetry.value.push(newReading)
      }
      
      // Update device status to online
      const device = devices.value.find(d => d.id === deviceId)
      if (device) {
        device.status = 'online'
        device.last_seen = new Date().toISOString()
      }
    }
  } catch (e: any) {
    error.value = e?.response?.data?.error || 'Erro ao simular leitura.'
  } finally {
    simulating.value = null
  }
}

async function generateSingleReading() {
  if (selectedDeviceId.value === null) return
  
  generatingSingle.value = true
  simulatorMessage.value = ''
  
  try {
    const result = await simulateTelemetry(selectedDeviceId.value, {
      base_power: 150,
      variation: 0.2,
      base_voltage: 220
    })
    
    simulatorMessage.value = `✅ Leitura gerada: ${result.power?.toFixed(1)} W`
    
    // Refresh telemetry data
    await loadDeviceTelemetry(selectedDeviceId.value)
  } catch (e: any) {
    simulatorMessage.value = `❌ ${e?.response?.data?.error || 'Erro ao simular'}`
  } finally {
    generatingSingle.value = false
  }
}

async function generateBulkReadings() {
  if (selectedDeviceId.value === null) return
  
  generatingBulk.value = true
  simulatorMessage.value = 'Gerando dados históricos...'
  
  try {
    const result = await simulateBulkTelemetry(selectedDeviceId.value, {
      base_power: 150,
      variation: 0.2,
      base_voltage: 220,
      count: 48, // 48 readings
      interval_sec: 1800 // 30 minutes apart = 24 hours of data
    })
    
    simulatorMessage.value = `✅ ${result.readings_created} leituras geradas com sucesso!`
    
    // Refresh telemetry data
    await loadDeviceTelemetry(selectedDeviceId.value)
    
    // Also refresh latest telemetry
    latestTelemetry.value = await getLatestTelemetry()
  } catch (e: any) {
    simulatorMessage.value = `❌ ${e?.response?.data?.error || 'Erro ao gerar histórico'}`
  } finally {
    generatingBulk.value = false
  }
}

function closeModal() {
  selectedDeviceId.value = null
  selectedTelemetry.value = []
  selectedSummary.value = null
  simulatorMessage.value = ''
}

function openTapoConfig(d: Device) {
  tapoError.value = ''
  editingTapoDeviceId.value = d.id
  // try to parse metadata
  try {
    const meta = d.metadata ? JSON.parse(d.metadata) : {}
    const t = meta?.tapo || {}
    tapoIp.value = t.ip || ''
    tapoUser.value = t.username || ''
    tapoPassword.value = t.password || ''
  } catch (e) {
    tapoIp.value = ''
    tapoUser.value = ''
    tapoPassword.value = ''
  }
}

async function saveTapoConfig() {
  if (editingTapoDeviceId.value === null) return
  const id = editingTapoDeviceId.value
  tapoError.value = ''
  savingTapo.value = true
  try {
    // fetch current device to merge other metadata fields
    const device = devices.value.find(d => d.id === id)
    if (!device) throw new Error('device not found')

    let meta: any = {}
    try { meta = device.metadata ? JSON.parse(device.metadata) : {} } catch { meta = {} }
    meta.tapo = { ip: tapoIp.value, username: tapoUser.value, password: tapoPassword.value }

    const updated = await updateDevice(id, { metadata: JSON.stringify(meta) })
    // update local list
    const idx = devices.value.findIndex(d => d.id === id)
    if (idx >= 0) devices.value[idx] = updated
    editingTapoDeviceId.value = null
  } catch (e: any) {
    tapoError.value = e?.response?.data?.error || e?.message || 'Erro ao salvar credenciais'
  } finally {
    savingTapo.value = false
  }
}

function cancelTapoConfig() {
  editingTapoDeviceId.value = null
  tapoIp.value = ''
  tapoUser.value = ''
  tapoPassword.value = ''
  tapoError.value = ''
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

.badge--offline {
  background: rgba(158, 158, 158, 0.2);
  color: #9e9e9e;
}

select.input {
  min-width: 150px;
}

.btn--sm {
  padding: 6px 12px;
  font-size: 0.85rem;
}

/* Power Toggle Section */
.power-toggle-section {
  display: flex;
  justify-content: center;
  padding: var(--sp-2) 0;
}

.power-toggle-btn {
  display: flex;
  align-items: center;
  gap: var(--sp-2);
  padding: var(--sp-3) var(--sp-5);
  border: 2px solid;
  border-radius: 50px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  background: transparent;
}

.power-toggle-btn.power-on {
  border-color: #f44336;
  color: #f44336;
}

.power-toggle-btn.power-on:hover:not(:disabled) {
  background: rgba(244, 67, 54, 0.15);
  box-shadow: 0 0 15px rgba(244, 67, 54, 0.3);
}

.power-toggle-btn.power-off {
  border-color: #4caf50;
  color: #4caf50;
}

.power-toggle-btn.power-off:hover:not(:disabled) {
  background: rgba(76, 175, 80, 0.15);
  box-shadow: 0 0 15px rgba(76, 175, 80, 0.3);
}

.power-toggle-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.toggle-icon {
  font-size: 1.2rem;
}

.toggle-loading {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
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

/* Simulator Section */
.simulator-section {
  background: rgba(0, 150, 136, 0.1);
  border: 1px solid rgba(0, 150, 136, 0.3);
  border-radius: 8px;
  padding: var(--sp-4);
  display: flex;
  flex-direction: column;
  gap: var(--sp-2);
}

.simulator-section h4 {
  margin: 0;
  font-size: 1rem;
  color: #00bfa5;
}

.simulator-controls {
  display: flex;
  gap: var(--sp-2);
  flex-wrap: wrap;
}

.simulator-message {
  margin-top: var(--sp-2);
  font-style: italic;
}

.btn--simulate {
  border-color: #00bfa5;
  color: #00bfa5;
}

.btn--simulate:hover {
  background: rgba(0, 150, 136, 0.1);
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

/* Status badge for consumption health */
.status-badge {
  display: inline-block;
  padding: 4px 10px;
  border-radius: 999px;
  font-weight: 600;
  font-size: 0.75rem;
}
.status-badge.ok {
  color: var(--ok, #1f8a3d);
  background: rgba(31,138,61,0.08);
  border: 1px solid rgba(31,138,61,0.12);
}
.status-badge.warning {
  color: var(--warning, #b36b00);
  background: rgba(179,107,0,0.08);
  border: 1px solid rgba(179,107,0,0.12);
}
.status-badge.danger {
  color: var(--danger, #c92a2a);
  background: rgba(201,42,42,0.08);
  border: 1px solid rgba(201,42,42,0.12);
}
.status-badge.unknown {
  color: var(--muted, #9e9e9e);
  background: rgba(158,158,158,0.06);
  border: 1px solid rgba(158,158,158,0.08);
}
</style>
