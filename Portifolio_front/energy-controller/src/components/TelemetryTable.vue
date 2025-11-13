<template>
  <div class="telemetry-table">
    <div class="table-header">
      <h3>Telemetria</h3>
      <button class="btn btn--primary" @click="handleRefresh">Atualizar</button>
    </div>
    
    <div class="table-wrapper">
      <table class="table">
        <thead>
          <tr>
            <th>Tempo</th>
            <th>Potência (W)</th>
            <th>Tensão (V)</th>
            <th>Corrente (A)</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="paginatedData.length === 0">
            <td colspan="4" class="text-center text-muted">Nenhum dado disponível</td>
          </tr>
          <tr v-for="(row, index) in paginatedData" :key="index">
            <td>{{ row.time }}</td>
            <td>{{ row.power_w }}</td>
            <td>{{ row.voltage }}</td>
            <td>{{ row.current }}</td>
          </tr>
        </tbody>
      </table>
    </div>

    <div class="pagination" v-if="totalPages > 1">
      <button 
        class="btn btn--outline btn--sm" 
        @click="previousPage" 
        :disabled="currentPage === 1"
      >
        Anterior
      </button>
      <span class="page-info">Página {{ currentPage }} de {{ totalPages }}</span>
      <button 
        class="btn btn--outline btn--sm" 
        @click="nextPage" 
        :disabled="currentPage === totalPages"
      >
        Próxima
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'

export interface TelemetryRow {
  time: string
  power_w: number
  voltage: number
  current: number
}

// Props definition
const props = defineProps<{
  data: TelemetryRow[]
}>()

// Emits definition
const emit = defineEmits<{
  refresh: []
}>()

// Pagination state
const currentPage = ref(1)
const perPage = 10

// Computed properties
const totalPages = computed(() => Math.ceil(props.data.length / perPage))

const paginatedData = computed(() => {
  const start = (currentPage.value - 1) * perPage
  const end = start + perPage
  return props.data.slice(start, end)
})

// Methods
function nextPage() {
  if (currentPage.value < totalPages.value) {
    currentPage.value++
  }
}

function previousPage() {
  if (currentPage.value > 1) {
    currentPage.value--
  }
}

function handleRefresh() {
  currentPage.value = 1
  emit('refresh')
}
</script>

<style scoped>
.telemetry-table {
  display: flex;
  flex-direction: column;
  gap: var(--sp-3, 12px);
}

.table-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.table-header h3 {
  margin: 0;
  font-size: var(--fs-lg, 1.125rem);
  font-weight: 600;
}

.table-wrapper {
  overflow-x: auto;
}

.table {
  width: 100%;
  border-collapse: collapse;
  font-size: var(--fs-sm, 0.875rem);
}

.table th,
.table td {
  padding: var(--sp-2, 8px) var(--sp-3, 12px);
  text-align: left;
  border-bottom: 1px solid var(--color-border, #e5e7eb);
}

.table th {
  background-color: var(--color-bg-muted, #f9fafb);
  font-weight: 600;
  color: var(--color-text-muted, #6b7280);
  text-transform: uppercase;
  font-size: var(--fs-xs, 0.75rem);
  letter-spacing: 0.05em;
}

.table tbody tr:hover {
  background-color: var(--color-bg-hover, #f9fafb);
}

.text-center {
  text-align: center;
}

.text-muted {
  color: var(--color-text-muted, #6b7280);
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: var(--sp-3, 12px);
  margin-top: var(--sp-2, 8px);
}

.page-info {
  font-size: var(--fs-sm, 0.875rem);
  color: var(--color-text-muted, #6b7280);
}

.btn {
  padding: var(--sp-2, 8px) var(--sp-3, 12px);
  border: none;
  border-radius: var(--radius, 4px);
  cursor: pointer;
  font-size: var(--fs-sm, 0.875rem);
  font-weight: 500;
  transition: all 0.2s;
}

.btn--primary {
  background-color: var(--color-primary, #3b82f6);
  color: white;
}

.btn--primary:hover {
  background-color: var(--color-primary-dark, #2563eb);
}

.btn--outline {
  background-color: transparent;
  border: 1px solid var(--color-border, #e5e7eb);
  color: var(--color-text, #1f2937);
}

.btn--outline:hover:not(:disabled) {
  background-color: var(--color-bg-muted, #f9fafb);
}

.btn--sm {
  padding: var(--sp-1, 4px) var(--sp-2, 8px);
  font-size: var(--fs-xs, 0.75rem);
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>
