<template>
  <div class="page">
    <h1>Entrar</h1>
    <form @submit.prevent="onSubmit">
      <input v-model="email" type="email" placeholder="Email" required />
      <input v-model="password" type="password" placeholder="Senha" required />
      <button type="submit">Login</button>
      <p>Sem conta? <router-link to="/register">Cadastrar</router-link></p>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";
import { useAuth } from "../stores/auth";

const email = ref("");
const password = ref("");
const auth = useAuth();
const router = useRouter();

async function onSubmit() {
  await auth.login(email.value, password.value);
  router.push("/app/dashboard");
}
</script>

<style>
.page {
  max-width: 360px;
  margin: 10vh auto;
  display: flex;
  flex-direction: column;
  gap: 12px;
}
input,
button {
  padding: 10px;
}
</style>
