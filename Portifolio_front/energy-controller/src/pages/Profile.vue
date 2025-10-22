<template>
  <div class="grid" style="padding:16px">
    <h2>Perfil</h2>


    <section class="card" style="padding:16px;display:grid;gap:12px">
      <div class="row" style="align-items:center;gap:12px">
        <div class="avatar">{{ initials }}</div>
        <div>
          <div style="font-weight:700">{{ auth.user?.name || 'Usuário' }}</div>
          <div class="text-muted small">{{ auth.user?.email }}</div>
        </div>
      </div>


      <div class="hr" />
      <div class="grid" style="grid-template-columns:repeat(auto-fit,minmax(220px,1fr));gap:12px">
        <div class="card" style="padding:12px">
          <div class="text-muted small">Função</div>
          <div>user</div>
        </div>
        <div class="card" style="padding:12px">
          <div class="text-muted small">Conta criada</div>
          <div>—</div>
        </div>
      </div>
    </section>
  </div>
</template>
<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useAuth } from '../stores/auth'
const auth = useAuth()
onMounted(() => { if (!auth.user) auth.fetchMe() })
const initials = computed(() => {
  const n = auth.user?.name || 'U I';
  return n.split(' ').map(x => x[0]).slice(0, 2).join('').toUpperCase()
})
</script>
<style scoped>
.avatar {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  background: #0f1731;
  border: 1px solid var(--border);
  display: grid;
  place-items: center;
  font-weight: 700
}
</style>