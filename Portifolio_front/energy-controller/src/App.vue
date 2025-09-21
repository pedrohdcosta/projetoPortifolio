
<template>
  <div class="app">
    <AppHeader v-if="isAuth" />
    <div class="layout">
      <AppSidebar v-if="isAuth" />
      <router-view />
    </div>
  </div>
</template>

<script setup lang="ts">
import {useAuth} from './stores/auth'
import AppHeader from './components/AppHeader.vue'
import AppSidebar from './components/AppSidebar.vue'
const auth = useAuth()
const isAuth = auth.isAuthenticated
</script>

<style>
body,html,#app{height:100%;margin:0;font-family:system-ui,Segoe UI,Arial}
.layout{display:flex;min-height:100vh}
</style>