<template>
  <div class="p">
    <h2>Perfil</h2>
    <pre>{{ user }}</pre>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from "vue";
import { useAuth } from "../stores/auth";
const auth = useAuth();
const user = auth.$state.user;
onMounted(() => {
  if (!user) auth.fetchMe();
});
</script>

<style>
.p {
  padding: 16px;
}
</style>
