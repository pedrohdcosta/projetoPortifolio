<!-- src/pages/Login.vue -->
<template>
  <main class="auth-wrap">
    <section class="auth-card">
      <header class="auth-header">
        <div class="logo-dot" />
        <h1>Energy Controller</h1>
        <p class="subtitle">Entre para acessar seu painel</p>
      </header>

      <form class="auth-form" @submit.prevent="onSubmit">
        <div class="field">
          <label for="email">E-mail</label>
          <input id="email" v-model.trim="email" type="email" inputmode="email" autocomplete="email"
            placeholder="voce@exemplo.com" required />
        </div>

        <div class="field">
          <label for="password">Senha</label>
          <div class="password-wrap">
            <input id="password" :type="showPass ? 'text' : 'password'" v-model="password"
              autocomplete="current-password" placeholder="••••••••" required minlength="4" />
            <button type="button" class="ghost" @click="showPass = !showPass">
              {{ showPass ? 'Ocultar' : 'Mostrar' }}
            </button>
          </div>
        </div>

        <button class="primary" type="submit" :disabled="loading">
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
/* Layout base */
.auth-wrap {
  min-height: 100dvh;
  display: grid;
  place-items: center;
  padding: 24px;
  background:
    radial-gradient(60% 40% at 10% 10%, rgba(67, 160, 71, 0.12), transparent 60%),
    radial-gradient(40% 30% at 90% 20%, rgba(25, 118, 210, 0.12), transparent 60%),
    linear-gradient(#0b1020, #0b1020) fixed;
  color: #e7eaf3;
}

/* Card */
.auth-card {
  width: min(440px, 92vw);
  background: rgba(16, 22, 44, 0.72);
  backdrop-filter: blur(8px);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 16px;
  padding: 28px;
  box-shadow:
    0 10px 30px rgba(0, 0, 0, 0.35),
    inset 0 1px 0 rgba(255, 255, 255, 0.06);
}

/* Header do card */
.auth-header {
  text-align: center;
  margin-bottom: 18px;
}

.auth-header h1 {
  margin: 8px 0 2px;
  font-size: 1.5rem;
  letter-spacing: 0.2px;
}

.subtitle {
  margin: 0;
  color: #95a2bf;
  font-size: 0.95rem;
}

/* Logo/ícone */
.logo-dot {
  width: 44px;
  height: 44px;
  margin: 0 auto;
  border-radius: 12px;
  background:
    conic-gradient(from 180deg at 50% 50%, #43a047, #1976d2, #43a047);
  box-shadow: 0 6px 20px rgba(76, 175, 80, 0.35);
}

/* Form */
.auth-form {
  display: grid;
  gap: 14px;
  margin-top: 10px;
}

.field {
  display: grid;
  gap: 8px;
}

label {
  font-size: 0.9rem;
  color: #b8c2dc;
}

input {
  width: 100%;
  background: #0f1731;
  color: #e9edf7;
  border: 1px solid #243153;
  border-radius: 10px;
  padding: 12px 14px;
  outline: none;
  transition: border-color 0.15s, box-shadow 0.15s;
}

input::placeholder {
  color: #8fa0c7;
}

input:focus {
  border-color: #4f8df7;
  box-shadow: 0 0 0 3px rgba(79, 141, 247, 0.18);
}

/* Senha + botão mostrar/ocultar */
.password-wrap {
  position: relative;
  display: flex;
  align-items: center;
  gap: 8px;
}

.password-wrap .ghost {
  margin-left: -4px;
  padding: 8px 10px;
  font-size: 0.85rem;
  border-radius: 8px;
  border: 1px dashed transparent;
  background: transparent;
  color: #a8b5d4;
  cursor: pointer;
}

.password-wrap .ghost:hover {
  color: #d6def0;
  border-color: rgba(255, 255, 255, 0.12);
}

/* Botões */
.primary {
  margin-top: 6px;
  width: 100%;
  padding: 12px 14px;
  border: 0;
  border-radius: 10px;
  background: linear-gradient(90deg, #43a047, #1976d2);
  color: white;
  font-weight: 600;
  letter-spacing: 0.2px;
  cursor: pointer;
  transition: transform 0.05s ease-in-out, filter 0.15s;
}

.primary:hover {
  filter: brightness(1.04);
}

.primary:active {
  transform: translateY(1px);
}

.primary:disabled {
  opacity: 0.65;
  cursor: not-allowed;
}

/* Loading spinner */
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

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

/* Mensagens */
.error {
  margin: 6px 0 0;
  padding: 10px 12px;
  border-radius: 8px;
  background: #361520;
  color: #ffb8c5;
  border: 1px solid rgba(255, 99, 132, 0.35);
  font-size: 0.92rem;
}

.muted {
  margin: 6px 0 0;
  color: #9aa9c9;
  font-size: 0.92rem;
  text-align: center;
}

.muted a {
  color: #e0ebff;
  text-decoration: underline;
}
</style>
