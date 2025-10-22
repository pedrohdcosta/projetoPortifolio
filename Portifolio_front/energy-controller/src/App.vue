<template>
  <div class="app">
    <header class="top card" v-if="isAuth">
      <div class="row">
        <div class="logo-dot" />
        <strong>Energy IoT</strong>
        <span class="badge">MVP</span>
      </div>
      <div class="row">
        <span class="text-muted small">{{ auth.user?.email }}</span>
        <button class="btn" @click="logout">Sair</button>
      </div>
    </header>


    <div class="layout">
      <aside class="side card" v-if="isAuth">
        <RouterLink to="/app/dashboard" class="nav" active-class="active">📊 Dashboard</RouterLink>
        <RouterLink to="/app/devices" class="nav" active-class="active">🔌 Dispositivos</RouterLink>
        <RouterLink to="/app/profile" class="nav" active-class="active">👤 Perfil</RouterLink>
      </aside>


      <main class="content">
        <RouterView />
      </main>
    </div>
  </div>
</template>
<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
import { useAuth } from './stores/auth'
const auth = useAuth()
const isAuth = auth.isAuthenticated
function logout() { auth.logout() }
</script>
<style scoped>
.top {
  position: sticky;
  top: 0;
  margin: 12px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px
}

.logo-dot {
  width: 24px;
  height: 24px;
  border-radius: 8px;
  background: conic-gradient(from 180deg, #43a047, #1976d2, #43a047);
  box-shadow: 0 4px 14px rgba(76, 175, 80, .35)
}

.layout {
  display: grid;
  grid-template-columns: 220px 1fr;
  gap: 12px;
  padding: 12px
}

.side {
  padding: 12px;
  height: calc(100dvh - 84px);
  position: sticky;
  top: 72px
}

.nav {
  display: block;
  padding: 10px 12px;
  border-radius: 10px;
  text-decoration: none;
  color: inherit
}

.nav:hover {
  background: #0f1731
}

.active {
  background: #0f1731;
  border: 1px solid var(--border)
}

.content {
  min-height: calc(100dvh - 84px)
}

@media (max-width: 900px) {
  .layout {
    grid-template-columns: 1fr
  }

  .side {
    display: none
  }
}
</style>