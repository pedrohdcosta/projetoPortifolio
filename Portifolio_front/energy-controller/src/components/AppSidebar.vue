<template>
  <aside class="s">
    <router-link to="/app/dashboard">Dashboard</router-link>
    <router-link to="/app/devices">Dispositivos</router-link>
    <router-link to="/app/profile">Perfil</router-link>
    <router-link to="/app/thresholds">Limiares</router-link>
  </aside>
</template>
<script setup lang="ts"></script>
<style>
.s {
  width: 220px;
  border-right: 1px solid #e5e7eb;
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 12px;
}
.s a {
  padding: 8px;
  border-radius: 6px;
  text-decoration: none;
  color: inherit;
}
.s a.router-link-active {
  background: #f3f4f6;
}
</style>
