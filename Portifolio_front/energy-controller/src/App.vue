<template>
  <div class="app">
    <!-- Header aparece só quando autenticado -->
    <header class="top card" v-if="isAuth">
      <div class="row">
        <div class="logo-dot" />
        <strong>Energy IoT</strong>
        <span class="badge">MVP</span>
      </div>
      <div class="row">
        <span class="text-muted small">{{ auth.user?.email }}</span>
        <button class="btn btn--outline" @click="logout">Sair</button>
      </div>
    </header>

    <!-- Layout autenticado: sidebar + conteúdo -->
    <div v-if="isAuth" class="layout">
      <aside class="side card">
        <RouterLink to="/app/dashboard" class="nav" active-class="active">📊 Dashboard</RouterLink>
        <RouterLink to="/app/devices"  class="nav" active-class="active">🔌 Dispositivos</RouterLink>
        <RouterLink to="/app/profile"  class="nav" active-class="active">👤 Perfil</RouterLink>
      </aside>

      <main class="content">
        <RouterView />
      </main>
    </div>

    <!-- Layout não-autenticado: viewport livre para centralizar Login/Register -->
    <main v-else class="auth-viewport">
      <RouterView />
    </main>
  </div>
</template>

<script setup lang="ts">
import { RouterLink, RouterView, useRouter } from 'vue-router'
import { computed } from 'vue'
import { useAuth } from './stores/auth'

const auth = useAuth()
const router = useRouter()
const isAuth = computed(() => auth.isAuthenticated)

function logout() {
  auth.logout()
  // navegar para /login para forçar o router a renderizar a tela de autenticação
  router.push('/login').catch(() => {})
}
</script>

<style scoped>
/* ===== Tokens locais ===== */
:host {
  --header-h: 64px;
  --gap-lg: 24px;
  --gap-md: 16px;
}

/* Casca do app: não limite largura aqui */
.app { width: 100%; }

/* ===== Header sticky centrado ===== */
.top {
  position: sticky;
  top: 0;
  z-index: 10;
  margin: 16px auto;
  padding: 12px 16px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  max-width: 1440px;
  min-height: var(--header-h);
}
@media (min-width: 1600px) {
  .top { max-width: 1600px; }
}

.logo-dot {
  width: 24px;
  height: 24px;
  border-radius: 8px;
  background: conic-gradient(from 180deg, #43a047, #1976d2, #43a047);
  box-shadow: 0 4px 14px rgba(76, 175, 80, .35);
}

/* ===== Layout autenticado (sidebar + conteúdo) ===== */
.layout {
  display: grid;
  grid-template-columns: 280px minmax(0, 1fr);
  gap: var(--gap-lg);
  padding: 0 var(--gap-lg) var(--gap-lg);
  max-width: 1440px;
  margin: 0 auto;
}
@media (min-width: 1600px) {
  .layout { max-width: 1600px; }
}

.side {
  padding: 16px;
  height: calc(100dvh - (var(--header-h) + 48px));
  position: sticky;
  top: calc(var(--header-h) + 16px);
  overflow: auto;
  overscroll-behavior: contain;
}

/* Navegação */
.nav {
  display: block;
  padding: 10px 12px;
  border-radius: 10px;
  text-decoration: none;
  color: inherit;
  outline: none;
}
.nav:hover { background: #0f1731; }
.active { background: #0f1731; border: 1px solid var(--border); }
.nav:focus-visible { box-shadow: 0 0 0 3px rgba(25,118,210,.28); }

/* Conteúdo */
.content {
  min-height: calc(100dvh - (var(--header-h) + 48px));
  min-width: 0;
}

/* ===== Layout não-autenticado (Login/Register centralizados) ===== */
.auth-viewport {
  min-height: 100dvh;
  display: grid;
  place-items: center;         /* centraliza o conteúdo do RouterView */
  padding: var(--gap-lg);
}

/* ===== Responsivo ===== */
@media (max-width: 1024px) {
  .layout {
    grid-template-columns: 1fr;
    gap: var(--gap-md);
    padding: 0 var(--gap-md) var(--gap-md);
  }
  .side { display: none; }
  .top  { margin: 12px 12px 0; }
}

/* ===== Acessibilidade/UX ===== */
@media (prefers-reduced-motion: reduce) {
  .nav, .btn { transition: none !important; }
}
</style>