<template>
  <div class="p">
    <h2>Dispositivos</h2>
    <form @submit.prevent="add">
      <input v-model="name" placeholder="Nome" required />
      <input v-model="room" placeholder="Ambiente" required />
      <button>Adicionar</button>
    </form>
    <ul>
      <li v-for="d in devices" :key="d.id">
        <strong>{{ d.name }}</strong> – {{ d.room }}
        <button @click="remove(d.id)">Remover</button>
      </li>
    </ul>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
const devices = ref<{ id: number; name: string; room: string }[]>([]);
const name = ref("");
const room = ref("");
let next = 1;
function add() {
  devices.value.push({ id: next++, name: name.value, room: room.value });
  name.value = "";
  room.value = "";
}
function remove(id: number) {
  devices.value = devices.value.filter((x) => x.id !== id);
}

</script>

<style>
.p {
  padding: 16px;
}
form {
  display: flex;
  gap: 8px;
  margin-bottom: 12px;
}
input,
button {
  padding: 8px;
}
</style>
