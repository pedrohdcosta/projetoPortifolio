<template>
  <main class="auth-wrap center">
    <section class="auth-card card">
      <header class="auth-header">
        <div class="logo-dot" />
        <h1>Criar conta</h1>
        <p class="subtitle">Acesse recursos do Energy IoT</p>
      </header>


      <form class="grid" @submit.prevent="onSubmit">
        <div class="grid">
          <label for="name">Nome</label>
          <input id="name" v-model.trim="name" class="input" placeholder="Seu nome" required />
        </div>
        <div class="grid">
          <label for="email">E-mail</label>
          <input id="email" v-model.trim="email" class="input" type="email" placeholder="voce@exemplo.com" required />
        </div>
        <div class="grid">
          <label for="password">Senha</label>
          <input id="password" v-model="password" class="input" type="password" placeholder="••••••••" required
            minlength="4" />
        </div>


        <button class="btn primary" :disabled="loading">
          <span v-if="!loading">Cadastrar</span>
          <span v-else class="spinner" />
        </button>
        <p v-if="error" class="error">{{ error }}</p>
        <p class="muted">Já tem conta? <router-link to="/login">Entrar</router-link></p>
      </form>
    </section>
  </main>
</template>
<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '../stores/auth'
const router = useRouter(); const auth = useAuth()
const name = ref(''); const email = ref(''); const password = ref('')
const loading = ref(false); const error = ref('')
async function onSubmit() {
  error.value = ''; loading.value = true
  try { await auth.signup(name.value, email.value, password.value); await auth.login(email.value, password.value); router.push('/app/dashboard') }
  catch (e: any) { error.value = e?.response?.data?.error || 'Falha no cadastro.' }
  finally { loading.value = false }
}
</script>
<style scoped>
@import url('../style/theme.css');

.auth-wrap {
  min-height: 100dvh;
  padding: 24px
}

.auth-card {
  width: min(520px, 92vw);
  padding: 28px
}

.auth-header {
  text-align: center;
  margin-bottom: 18px
}

.auth-header h1 {
  margin: 8px 0 2px;
  font-size: 1.5rem
}

.subtitle {
  margin: 0;
  color: var(--hint)
}

.spinner {
  display: inline-block;
  width: 18px;
  height: 18px;
  border: 2px solid rgba(255, 255, 255, .5);
  border-top-color: #fff;
  border-radius: 50%;
  animation: spin .8s linear infinite;
  vertical-align: -3px
}

@keyframes spin {
  to {
    transform: rotate(360deg)
  }
}

.error {
  margin: 6px 0 0;
  padding: 10px 12px;
  border-radius: 8px;
  background: #361520;
  color: #ffb8c5;
  border: 1px solid rgba(255, 99, 132, .35);
  font-size: .92rem
}

.muted {
  margin: 6px 0 0;
  color: var(--muted);
  font-size: .92rem;
  text-align: center
}
</style>
