<template>
  <div class="grid" style="padding:16px">
    <h2>Dispositivos</h2>


    <section class="card" style="padding:16px">
      <form class="row" @submit.prevent="add">
        <input v-model="name" class="input" placeholder="Nome" required />
        <input v-model="room" class="input" placeholder="Ambiente" required />
        <button class="btn primary">Adicionar</button>
      </form>
    </section>


    <section class="grid" style="grid-template-columns: repeat(auto-fit, minmax(260px,1fr));">
      <div v-for="d in devices" :key="d.id" class="card" style="padding:16px;display:grid;gap:8px">
        <div class="row" style="justify-content:space-between">
          <strong>{{ d.name }}</strong>
          <span class="badge">{{ d.room }}</span>
        </div>
        <div class="row" style="justify-content:space-between">
          <span class="text-muted small">ID: {{ d.id }}</span>
          <button class="btn ghost" @click="remove(d.id)">Remover</button>
        </div>
      </div>
    </section>


    <p v-if="devices.length === 0" class="text-muted">Nenhum dispositivo ainda. Adicione o primeiro acima.</p>
  </div>
</template>
<script setup lang="ts">
import { ref } from 'vue'
const devices = ref<{ id: number; name: string; room: string }[]>([])
const name = ref(''); const room = ref(''); let next = 1
function add() { devices.value.push({ id: next++, name: name.value, room: room.value }); name.value = ''; room.value = '' }
function remove(id: number) { devices.value = devices.value.filter(x => x.id !== id) }
</script>