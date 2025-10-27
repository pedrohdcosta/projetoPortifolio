<template>
  <main class="auth-wrap center">
    <section class="auth-card card">
      <header class="auth-header">
        <div class="logo-dot" />
        <h1>Criar conta</h1>
        <p class="subtitle">Acesse recursos do Energy IoT</p>
      </header>

      <form class="auth-form" @submit.prevent="onSubmit" novalidate>
        <div class="field">
          <label for="name">Nome</label>
          <input
            id="name"
            v-model.trim="name"
            class="input"
            placeholder="Seu nome"
            required
            autocomplete="name"
          />
        </div>

        <div class="field">
          <label for="email">E-mail</label>
          <input
            id="email"
            v-model.trim="email"
            class="input"
            type="email"
            placeholder="voce@exemplo.com"
            required
            autocomplete="email"
          />
        </div>

        <div class="field">
          <label for="password">Senha</label>
          <input
            id="password"
            v-model="password"
            class="input"
            type="password"
            placeholder="••••••••"
            required
            minlength="4"
            autocomplete="new-password"
          />
        </div>

        <button class="btn primary" :disabled="loading">
          <span v-if="!loading">Cadastrar</span>
          <span v-else class="spinner" aria-hidden="true" />
        </button>

        <p v-if="error" class="error">{{ error }}</p>

        <p class="muted">
          Já tem conta?
          <router-link to="/login">Entrar</router-link>
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

const name = ref('')
const email = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')

async function onSubmit() {
  error.value = ''
  loading.value = true
  try {
    await auth.signup(name.value, email.value, password.value)
    await auth.login(email.value, password.value)
    router.push('/app/dashboard')
  } catch (e: any) {
    error.value = e?.response?.data?.error || 'Falha no cadastro.'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
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

/* Card */
.auth-card {
  width: min(480px, 94vw);
  padding: 28px;
}

/* Header */
.auth-header {
  text-align: center;
  margin-bottom: 18px;
}
.auth-header h1 {
  margin: 8px 0 2px;
  font-size: var(--fs-xl);
}
.subtitle {
  margin: 0;
  color: var(--hint);
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

/* Spinner */
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
  border: 1px solid rgba(255, 99, 132, .35);
  font-size: .92rem;
}
.muted {
  margin: 6px 0 0;
  color: var(--muted);
  font-size: .92rem;
  text-align: center;
}
.muted a {
  color: #e0ebff;
  text-decoration: underline;
}
</style>
