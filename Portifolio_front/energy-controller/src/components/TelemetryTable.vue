<template>
  <div class="telemetry-table-wrapper">
    <div v-if="loading" class="loading-state">
      <SkeletonCard />
    </div>
    <div v-else class="telemetry-table card">
      <div class="table-header">
        <h3>Telemetria de Dispositivos</h3>
        <button 
          class="btn btn--outline btn--sm" 
          @click="$emit('refresh')"
          :disabled="loading"
        >
          Atualizar
        </button>
      </div>
      
      <div v-if="totalItems === 0" class="empty-state">
        <p class="text-muted">Nenhum dado de telemetria disponível.</p>
      </div>

      <table v-else>
        <thead>
          <tr>
            <th>Dispositivo</th>
            <th>Potência (W)</th>
            <th>Timestamp</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in paginatedData" :key="item.id">
            <td>{{ item.deviceName }}</td>
            <td>{{ item.power.toFixed(2) }}</td>
            <td>{{ formatTimestamp(item.timestamp) }}</td>
          </tr>
        </tbody>
      </table>

      <div class="pagination" v-if="totalPages > 1">
        <button class="btn btn--outline btn--sm" @click="goPrev" :disabled="currentPage === 1">Anterior</button>

        <div class="page-numbers">
          <button
            v-for="p in totalPages"
            :key="p"
            class="btn btn--sm page-btn"
            @click="changePage(p)"
            :class="{ active: currentPage === p }"
          >
            {{ p }}
          </button>
        </div>

        <button class="btn btn--outline btn--sm" @click="goNext" :disabled="currentPage === totalPages">Próximo</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import SkeletonCard from './SkeletonCard.vue'

export interface TelemetryData {
  id: number
  deviceName: string
  power: number
  timestamp: string
}

const props = defineProps<{
  data: TelemetryData[]
  loading?: boolean
}>()

const emit = defineEmits<{
  refresh: []
}>()

// Pagination
const pageSize = 10
const currentPage = ref(1)
const totalItems = computed(() => (props.data ?? []).length)
const totalPages = computed(() => Math.max(1, Math.ceil(totalItems.value / pageSize)))

const paginatedData = computed(() => {
  const start = (currentPage.value - 1) * pageSize
  return (props.data || []).slice(start, start + pageSize)
})

function changePage(p: number) {
  if (p >= 1 && p <= totalPages.value) currentPage.value = p
}

function goPrev() {
  if (currentPage.value > 1) currentPage.value--
}

function goNext() {
  if (currentPage.value < totalPages.value) currentPage.value++
}

// Reset page when data changes
watch(() => props.data, () => { currentPage.value = 1 })

function formatTimestamp(ts: string): string {
  return new Date(ts).toLocaleString('pt-BR')
}
</script>

<style scoped>
.telemetry-table-wrapper {
  width: 100%;
}

.loading-state {
  min-height: 200px;
}

.telemetry-table {
  padding: var(--sp-5, 20px);
}

.table-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--sp-4, 16px);
}

.table-header h3 {
  font-size: var(--fs-lg, 1.25rem);
  font-weight: 600;
}

.btn--sm {
  padding: 8px 16px;
  font-size: 0.875rem;
}

table {
  width: 100%;
  border-collapse: collapse;
}

thead {
  border-bottom: 1px solid var(--border, rgba(255, 255, 255, 0.1));
}

th {
  text-align: left;
  padding: var(--sp-3, 12px);
  font-weight: 600;
  color: var(--text-muted, rgba(255, 255, 255, 0.7));
  font-size: 0.875rem;
}

td {
  padding: var(--sp-3, 12px);
  border-bottom: 1px solid var(--border, rgba(255, 255, 255, 0.05));
}

.empty-state {
  padding: var(--sp-6, 24px);
  text-align: center;
}

.pagination {
  display: flex;
  gap: 12px;
  align-items: center;
  justify-content: flex-end;
  margin-top: 12px;
}

.page-numbers {
  display: flex;
  gap: 6px;
  align-items: center;
}

.page-btn {
  min-width: 36px;
  padding: 6px 10px;
  border-radius: 6px;
  background: transparent;
  border: 1px solid var(--border, rgba(255,255,255,0.06));
}

.page-btn.active {
  background: var(--accent, rgba(0, 123, 255, 0.12));
  border-color: var(--accent, rgba(0, 123, 255, 0.25));
}
</style>
