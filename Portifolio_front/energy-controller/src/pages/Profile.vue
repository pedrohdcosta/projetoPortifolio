<template>
  <div class="container">
    <div class="col" style="gap: var(--sp-6)">
      <!-- Cabeçalho -->
      <header class="row" style="justify-content: space-between;">
        <h2 style="font-size: var(--fs-xl); font-weight: 700;">Perfil</h2>
        <span class="badge">MVP</span>
      </header>

      <!-- Error State -->
      <div v-if="error" class="error-card card">
        <p class="error-message">{{ error }}</p>
        <button class="btn btn--solid" @click="refresh">Tentar novamente</button>
      </div>

      <!-- Loading -->
      <section v-if="!auth.user && loading" class="card" style="padding: var(--sp-5);">
        <SkeletonCard />
      </section>

      <!-- Card principal -->
      <section v-else class="card" style="padding: var(--sp-5); display: grid; gap: var(--sp-4);">
        <!-- Info principal -->
        <div class="row" style="align-items: center; gap: var(--sp-3);">
          <div class="avatar">{{ initials }}</div>
          <div class="col" style="gap: 6px;">
            <strong style="font-size: var(--fs-lg);">
              {{ auth.user?.name || 'Usuário' }}
            </strong>
            <div class="row" style="align-items: center; gap: var(--sp-2);">
              <span class="text-muted small">{{ auth.user?.email || '—' }}</span>
              <button
                v-if="auth.user?.email"
                class="btn btn--outline btn-xs"
                @click="copyEmail"
                aria-label="Copiar e-mail"
              >
                Copiar
              </button>
            </div>
          </div>
        </div>

        <div class="hr" />

        <!-- Apenas uma info extra genérica -->
        <div class="grid" style="grid-template-columns: repeat(auto-fit, minmax(220px, 1fr)); gap: var(--sp-3);">
          <article class="card" style="padding: var(--sp-4); text-align: center;">
            <div class="text-muted small">Tipo de Conta</div>
            <div>Usuário padrão</div>
          </article>
        </div>

        <!-- Ações -->
        <div class="row" style="justify-content: flex-end; gap: var(--sp-3);">
          <button class="btn btn--outline" @click="refresh" :disabled="loading">
            <span v-if="!loading">Atualizar dados</span>
            <LoadingSpinner v-else small />
          </button>
        </div>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useAuth } from '../stores/auth'
import SkeletonCard from '../components/SkeletonCard.vue'
import LoadingSpinner from '../components/LoadingSpinner.vue'

const auth = useAuth()
const loading = ref(false)
const error = ref('')

onMounted(async () => {
  if (!auth.user) {
    loading.value = true
    error.value = ''
    try { 
      await auth.fetchMe() 
    } catch (e: any) {
      error.value = e?.response?.data?.error || 
        'Erro ao carregar dados do usuário. Verifique sua conexão e tente novamente.'
    } finally { 
      loading.value = false 
    }
  }
})

const initials = computed(() => {
  const n = (auth.user?.name || 'U I').trim()
  return n.split(/\s+/).map(x => x[0]).slice(0, 2).join('').toUpperCase()
})

function copyEmail() {
  if (!auth.user?.email) return
  navigator.clipboard.writeText(auth.user.email).catch(() => {})
}

async function refresh() {
  loading.value = true
  error.value = ''
  try { 
    await auth.fetchMe() 
  } catch (e: any) {
    error.value = e?.response?.data?.error || 
      'Erro ao atualizar dados do usuário. Verifique sua conexão e tente novamente.'
  } finally { 
    loading.value = false 
  }
}
</script>

<style scoped>
.avatar {
  width: 56px;
  height: 56px;
  border-radius: 12px;
  background: #0f1731;
  border: 1px solid var(--border);
  display: grid;
  place-items: center;
  font-weight: 700;
}

.btn-xs {
  padding: 6px 10px;
  font-size: var(--fs-sm);
  border-radius: 8px;
}

.error-card {
  padding: var(--sp-5, 20px);
  background: rgba(255, 99, 132, 0.1);
  border: 1px solid rgba(255, 99, 132, 0.3);
  display: flex;
  flex-direction: column;
  gap: var(--sp-3, 12px);
  align-items: center;
}

.error-message {
  color: var(--warn, #ff6384);
  margin: 0;
  text-align: center;
}
</style>
