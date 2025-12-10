<template>
    <header class="h">
        <strong>Energy Controller</strong>
        <div class="spacer" />
        <button @click="logout">Sair</button>
    </header>
</template>

<script setup lang="ts">
    import { useAuth } from '../stores/auth'
    import { useRouter } from 'vue-router'
    const auth = useAuth(); const router = useRouter()
    function logout(){ auth.logout(); router.push('/login') }
</script>

<style>
.h{display:flex;align-items:center;gap:12px;padding:10px;border-bottom:1px solid #e5e7eb}
.spacer{flex:1}
</style>