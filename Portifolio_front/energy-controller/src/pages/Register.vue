<template>
  <div class="page">
    <h1>Cadastrar</h1>
    <form @submit.prevent="onSubmit">
      <input v-model="name" placeholder="Nome" required />
      <input v-model="email" type="email" placeholder="Email" required />
      <input v-model="password" type="password" placeholder="Senha" required />
      <button type="submit">Criar conta</button>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";
import { useAuth } from "../stores/auth";

const name = ref("");
const email = ref("");
const password = ref("");
const auth = useAuth();
const router = useRouter();

async function onSubmit() {
  await auth.signup(name.value, email.value, password.value);
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
