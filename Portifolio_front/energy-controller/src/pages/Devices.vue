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
          v-for="d in devices"
          :key="d.id"
          class="card"
          style="padding: var(--sp-5); display: grid; gap: var(--sp-3);"
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

          <div class="row" style="justify-content: space-between; margin-top: auto;">
            <span class="text-muted small">ID: {{ d.id }}</span>
            <div class="row" style="gap: var(--sp-2);">
              <button 
                class="btn btn--outline" 
                @click="toggleStatus(d.id)"
                :disabled="loading"
              >
                {{ d.status === 'online' ? 'Desativar' : 'Ativar' }}
              </button>
              <button 
                class="btn btn--outline" 
                @click="remove(d.id)"
                :disabled="loading"
              >
                Remover
              </button>
            </div>
          </div>
        </article>
      </section>

      <!-- Empty state -->
      <section v-else class="card" style="padding: var(--sp-5);">
        <p class="text-muted">Nenhum dispositivo cadastrado. Adicione seu primeiro dispositivo IoT acima.</p>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted } from 'vue'
import LoadingSpinner from '../components/LoadingSpinner.vue'
import { listDevices, createDevice, deleteDevice, type Device } from '../api/devices'

const devices = ref<Device[]>([])
const name = ref('')
const room = ref('')
const deviceType = ref('smart_plug')
const loading = ref(false)
const loadingList = ref(false)
const error = ref('')

const canAdd = computed(() => name.value.trim().length > 0)

function getDeviceTypeLabel(type?: string): string {
  switch (type) {
    case 'smart_plug': return '🔌 Tomada'
    case 'sensor': return '📡 Sensor'
    case 'meter': return '⚡ Medidor'
    default: return '📱 Dispositivo'
  }
}

async function loadDevices() {
  loadingList.value = true
  error.value = ''
  
  try {
    devices.value = await listDevices()
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
  } catch (e: any) {
    error.value = e?.response?.data?.error || 
      'Erro ao remover dispositivo. Tente novamente.'
  } finally {
    loading.value = false
  }
}

async function toggleStatus(id: number) {
  if (loading.value) return

  loading.value = true
  error.value = ''

  try {
    const device = devices.value.find(d => d.id === id)
    if (device) {
      device.status = device.status === 'online' ? 'offline' : 'online'
    }
  } catch (e: any) {
    error.value = e?.response?.data?.error || 
      'Erro ao alterar status do dispositivo. Tente novamente.'
  } finally {
    loading.value = false
  }
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
  .device-grid > .card { grid-column: span 4; } /* 3 por linha */
}
@media (min-width: 1400px) {
  .device-grid > .card { grid-column: span 3; } /* 4 por linha em telas grandes */
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
</style>
