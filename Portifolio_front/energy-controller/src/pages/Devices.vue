<template>
  <div class="container">
    <div class="col" style="gap: var(--sp-6)">

      <!-- Título -->
      <header class="row" style="justify-content: space-between;">
        <h2 style="font-size: var(--fs-xl); font-weight: 700;">Dispositivos</h2>
        <span class="badge">MVP</span>
      </header>

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
          />
          <input
            v-model.trim="room"
            class="input"
            placeholder="Ambiente"
            aria-label="Ambiente do dispositivo"
            required
            autocomplete="off"
          />
          <button class="btn btn--solid" :disabled="!canAdd">
            Adicionar
          </button>
        </form>
      </section>

      <!-- Lista -->
      <section class="grid device-grid">
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
            <button class="btn btn--outline" @click="remove(d.id)">Remover</button>
          </div>
        </article>
      </section>

      <!-- Empty state -->
      <section v-if="devices.length === 0" class="card" style="padding: var(--sp-5);">
        <p class="text-muted">Nenhum dispositivo ainda. Adicione o primeiro acima.</p>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'

type Device = { id: number; name: string; room: string }

const devices = ref<Device[]>([])
const name = ref('')
const room = ref('')

let next = 1
const canAdd = computed(() => name.value.trim().length > 0 && room.value.trim().length > 0)

function add() {
  if (!canAdd.value) return
  devices.value.push({ id: next++, name: name.value.trim(), room: room.value.trim() })
  name.value = ''
  room.value = ''
}

function remove(id: number) {
  devices.value = devices.value.filter(x => x.id !== id)
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
</style>
