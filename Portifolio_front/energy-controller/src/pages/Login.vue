<!-- src/pages/Login.vue -->
<template>
  <main class="auth-wrap">
    <section class="auth-card card">
      <header class="auth-header">
        <div class="logo-dot" />
        <h1>Energy Controller</h1>
        <p class="subtitle">Entre para acessar seu painel</p>
      </header>

      <form class="auth-form" @submit.prevent="onSubmit" novalidate>
        <div class="field">
          <label for="email">E-mail</label>
          <input
            id="email"
            v-model.trim="email"
            class="input"
            type="email"
            inputmode="email"
            autocomplete="email"
            placeholder="voce@exemplo.com"
            required
          />
        </div>

        <div class="field">
          <label for="password">Senha</label>
          <div class="password-wrap">
            <input
              id="password"
              :type="showPass ? 'text' : 'password'"
              v-model.trim="password"
              class="input"
              autocomplete="current-password"
              placeholder="••••••••"
              required
              minlength="4"
            />
            <button
              type="button"
              class="btn btn--outline btn-eye"
              @click="showPass = !showPass"
              :aria-pressed="showPass"
            >
              {{ showPass ? 'Ocultar' : 'Mostrar' }}
            </button>
          </div>
        </div>

        <button class="btn primary" type="submit" :disabled="loading">
          <span v-if="!loading">Entrar</span>
          <span v-else class="spinner" aria-hidden="true" />
        </button>

        <p v-if="error" class="error">{{ error }}</p>
        <p class="muted">
          Sem conta?
          <router-link to="/register">Criar cadastro</router-link>
        </p>
      </form>
    </section>
  </main>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '../stores/auth'

const router = useRouter()
const auth = useAuth()

const email = ref('')
const password = ref('')
const showPass = ref(false)
const loading = ref(false)
const error = ref('')

async function onSubmit() {
  error.value = ''
  loading.value = true
  try {
    await auth.login(email.value, password.value)
    router.push('/app/dashboard')
  } catch (e: any) {
    error.value =
      e?.response?.data?.error ||
      'Falha no login. Verifique e-mail/senha e tente novamente.'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
/* Layout base (usa teu fundo/cores de theme.css) */
.auth-wrap {
  min-height: 100dvh;
  display: grid;
  place-items: center;
  padding: var(--sp-6);
  background:
    radial-gradient(60% 40% at 10% 10%, rgba(67, 160, 71, 0.12), transparent 60%),
    radial-gradient(40% 30% at 90% 20%, rgba(25, 118, 210, 0.12), transparent 60%),
    linear-gradient(var(--bg), var(--bg)) fixed;
  color: var(--text);
}

/* Card aproveita .card do tema (vidro, sombra, raio) e só define largura/padding */
.auth-card {
  width: min(480px, 94vw);
  padding: 28px;
}

/* Header do card */
.auth-header {
  text-align: center;
  margin-bottom: 18px;
}
.auth-header h1 {
  margin: 8px 0 2px;
  font-size: var(--fs-xl);
  letter-spacing: 0.2px;
}
.subtitle {
  margin: 0;
  color: var(--hint);
  font-size: 0.95rem;
}

/* Logo/ícone (igual ao resto do app) */
.logo-dot {
  width: 44px;
  height: 44px;
  margin: 0 auto;
  border-radius: 12px;
  background: conic-gradient(from 180deg at 50% 50%, #43a047, #1976d2, #43a047);
  box-shadow: 0 6px 20px rgba(76, 175, 80, 0.35);
}

/* Form */
.auth-form {
  display: grid;
  gap: 14px;
  margin-top: 10px;
}
.field { display: grid; gap: 8px; }
label { font-size: 0.9rem; color: #b8c2dc; }

/* Senha + botão mostrar/ocultar */
.password-wrap {
  position: relative;
  display: flex;
  align-items: center;
  gap: 8px;
}
.btn-eye {
  padding: 8px 10px;
  font-size: 0.85rem;
  border-radius: 8px;
}

/* Loading spinner (reaproveitado) */
.spinner {
  display: inline-block;
  width: 18px;
  height: 18px;
  border: 2px solid rgba(255, 255, 255, .5);
  border-top-color: #fff;
  border-radius: 50%;
  animation: spin .8s linear infinite;
  vertical-align: -3px;
}
@keyframes spin { to { transform: rotate(360deg); } }

/* Mensagens */
.error {
  margin: 6px 0 0;
  padding: 10px 12px;
  border-radius: 8px;
  background: #361520;
  color: var(--warn);
  border: 1px solid rgba(255, 99, 132, 0.35);
  font-size: 0.92rem;
}
.muted {
  margin: 6px 0 0;
  color: var(--muted);
  font-size: 0.92rem;
  text-align: center;
}
.muted a {
  color: #e0ebff;
  text-decoration: underline;
}

/* Acessibilidade/UX */
@media (prefers-reduced-motion: reduce) {
  .btn, input { transition: none !important; }
}
</style>
