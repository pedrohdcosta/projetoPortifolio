<template>
  <div class="container">
    <div class="col" style="gap: var(--sp-6)">

      <!-- Título -->
      <header class="row" style="justify-content: space-between;">
        <h2 style="font-size: var(--fs-xl); font-weight: 700;">Dispositivos</h2>
        <span class="badge">MVP</span>
      </header>

      <!-- Error State -->
      <div v-if="error" class="error-card card">
        <p class="error-message">{{ error }}</p>
        <button class="btn btn--solid" @click="clearError">Fechar</button>
      </div>

      <!-- Formulário -->
      <section class="card" style="padding: var(--sp-5);">
        <form class="row" style="gap: var(--sp-3); align-items: stretch;" @submit.prevent="add">
          <input
            v-model.trim="name"
            class="input"
            placeholder="Nome"
            aria-label="Nome do dispositivo"
            required
            autocomplete="off"
            :disabled="loading"
          />
          <input
            v-model.trim="room"
            class="input"
            placeholder="Ambiente"
            aria-label="Ambiente do dispositivo"
            required
            autocomplete="off"
            :disabled="loading"
          />
          <button class="btn btn--solid" :disabled="!canAdd || loading">
            <span v-if="!loading">Adicionar</span>
            <LoadingSpinner v-else small />
          </button>
        </form>
      </section>

      <!-- Lista -->
      <section v-if="devices.length > 0" class="grid device-grid">
        <article
          v-for="d in devices"
          :key="d.id"
          class="card"
          style="padding: var(--sp-5); display: grid; gap: var(--sp-3);"
        >
          <div class="row" style="justify-content: space-between;">
            <strong style="font-size: var(--fs-lg);">{{ d.name }}</strong>
            <span class="badge">{{ d.room }}</span>
          </div>

          <div class="row" style="justify-content: space-between;">
            <span class="text-muted small">ID: {{ d.id }}</span>
            <button 
              class="btn btn--outline" 
              @click="remove(d.id)"
              :disabled="loading"
            >
              Remover
            </button>
          </div>
        </article>
      </section>

      <!-- Empty state -->
      <section v-else class="card" style="padding: var(--sp-5);">
        <p class="text-muted">Nenhum dispositivo ainda. Adicione o primeiro acima.</p>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import LoadingSpinner from '../components/LoadingSpinner.vue'
// Uncomment when API is ready:
// import api from '../api/axios'

type Device = { id: number; name: string; room: string }

const devices = ref<Device[]>([])
const name = ref('')
const room = ref('')
const loading = ref(false)
const error = ref('')

let next = 1
const canAdd = computed(() => name.value.trim().length > 0 && room.value.trim().length > 0)

async function add() {
  if (!canAdd.value || loading.value) return
  
  loading.value = true
  error.value = ''
  
  try {
    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 500))
    // const { data } = await api.post('/devices', {
    //   name: name.value.trim(),
    //   room: room.value.trim()
    // })
    
    // For MVP, add locally
    devices.value.push({ 
      id: next++, 
      name: name.value.trim(), 
      room: room.value.trim() 
    })
    name.value = ''
    room.value = ''
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
    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 300))
    // await api.delete(`/devices/${id}`)
    
    devices.value = devices.value.filter(x => x.id !== id)
  } catch (e: any) {
    error.value = e?.response?.data?.error || 
      'Erro ao remover dispositivo. Tente novamente.'
  } finally {
    loading.value = false
  }
}

function clearError() {
  error.value = ''
}
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
</style>
