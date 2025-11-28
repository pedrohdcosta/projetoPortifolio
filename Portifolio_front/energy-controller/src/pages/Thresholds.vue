<template>
  <div class="container">
    <header class="row" style="justify-content: space-between;">
      <h2>Configurar Thresholds por Dispositivo</h2>
      <p class="text-muted">Defina limites de aviso e perigo (Watts) por dispositivo.</p>
    </header>

    <section class="card" style="padding: var(--sp-4);">
      <div v-if="loading" class="text-muted">Carregando dispositivos...</div>

      <div v-else>
        <div v-if="devices.length === 0" class="text-muted">Nenhum dispositivo encontrado.</div>

        <div v-for="d in devices" :key="d.id" class="threshold-row card" style="margin-bottom: var(--sp-3); padding: var(--sp-3);">
          <div class="row" style="justify-content: space-between; align-items: center;">
            <strong>{{ d.name }}</strong>
            <div style="display:flex; gap:8px; align-items:center;">
              <label class="text-muted" style="font-size:0.85rem">Aviso</label>
              <input type="number" v-model.number="model[d.id].warning" class="input" style="width:120px" />
              <label class="text-muted" style="font-size:0.85rem">Perigo</label>
              <input type="number" v-model.number="model[d.id].danger" class="input" style="width:120px" />
              <button class="btn btn--solid btn--sm" @click="save(d.id)">Salvar</button>
              <button class="btn btn--outline btn--sm" @click="reset(d.id)">Reset</button>
            </div>
          </div>
          <div class="text-muted small">ID: {{ d.id }} • Tipo: {{ getDeviceTypeLabel(d.type) }}</div>
          <div v-if="errors[d.id]" class="text-muted" style="color: #c92a2a">{{ errors[d.id] }}</div>
        </div>

        <div style="margin-top: var(--sp-4);">
          <button class="btn btn--solid" @click="saveAll">Salvar todos</button>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { listDevices, type Device } from '../api/devices'
import { getThresholds, setThresholds, getAllThresholds } from '../utils/thresholds'

const devices = ref<Device[]>([])
const loading = ref(false)
const errors = reactive<Record<number, string>>({})

// model maps deviceId -> {warning,danger}
const model = reactive<Record<number, {warning: number, danger: number}>>({})

function getDeviceTypeLabel(type?: string): string {
  switch (type) {
    case 'smart_plug': return '🔌 Tomada'
    case 'sensor': return '📡 Sensor'
    case 'meter': return '⚡ Medidor'
    default: return '📱 Dispositivo'
  }
}

async function load() {
  loading.value = true
  try {
    devices.value = await listDevices()
    // init model from stored thresholds
    const all = getAllThresholds()
    for (const d of devices.value) {
      const t = all[String(d.id)]
      model[d.id] = {
        warning: t?.warning ?? 100,
        danger: t?.danger ?? 500
      }
    }
  } finally {
    loading.value = false
  }
}

function validateValues(w: number, g: number): string | null {
  if (!isFinite(w) || w <= 0) return 'Valor de aviso inválido'
  if (!isFinite(g) || g <= 0) return 'Valor de perigo inválido'
  if (g <= w) return 'Perigo deve ser maior que Aviso'
  return null
}

function save(deviceId: number) {
  errors[deviceId] = ''
  const m = model[deviceId]
  const err = validateValues(m.warning, m.danger)
  if (err) { errors[deviceId] = err; return }
  setThresholds(deviceId, { warning: m.warning, danger: m.danger })
}

function reset(deviceId: number) {
  // remove saved and reset to defaults
  setThresholds(deviceId, { warning: 100, danger: 500 })
  model[deviceId].warning = 100
  model[deviceId].danger = 500
}

function saveAll() {
  for (const d of devices.value) save(d.id)
}

onMounted(load)
</script>

<style scoped>
.threshold-row { display: flex; flex-direction: column; gap: 6px }
.input { padding: 6px; }
</style>
