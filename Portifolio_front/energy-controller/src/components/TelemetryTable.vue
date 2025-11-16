<template>
  <div class="telemetry-table-container">
    <div class="table-header">
      <h3>Dados de Telemetria</h3>
      <button 
        class="btn btn--outline" 
        @click="$emit('refresh')"
        :disabled="loading"
      >
        <span v-if="loading" class="spinner-small"></span>
        <span v-else>🔄</span>
        Atualizar
      </button>
    </div>

    <SkeletonCard v-if="loading" />
    
    <div v-else-if="paginatedData.length > 0" class="table-wrapper">
      <table class="telemetry-table">
        <thead>
          <tr>
            <th>Timestamp</th>
            <th>Potência (W)</th>
            <th>Tensão (V)</th>
            <th>Corrente (A)</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in paginatedData" :key="item.id">
            <td>{{ formatTimestamp(item.timestamp) }}</td>
            <td>{{ item.power.toFixed(2) }}</td>
            <td>{{ item.voltage?.toFixed(2) ?? 'N/A' }}</td>
            <td>{{ item.current?.toFixed(2) ?? 'N/A' }}</td>
          </tr>
        </tbody>
      </table>

      <div v-if="totalPages > 1" class="pagination">
        <button 
          class="btn btn--sm" 
          @click="prevPage" 
          :disabled="currentPage === 1"
        >
          ← Anterior
        </button>
        <span class="page-info">
          Página {{ currentPage }} de {{ totalPages }}
        </span>
        <button 
          class="btn btn--sm" 
          @click="nextPage" 
          :disabled="currentPage === totalPages"
        >
          Próxima →
        </button>
      </div>
    </div>

    <div v-else class="no-data">
      <p>Nenhum dado de telemetria disponível</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import SkeletonCard from './SkeletonCard.vue';

interface TelemetryData {
  id: number;
  deviceId: number;
  timestamp: string;
  power: number;
  voltage?: number;
  current?: number;
}

const props = defineProps<{
  data: TelemetryData[];
  loading?: boolean;
}>();

defineEmits<{
  refresh: [];
}>();

const currentPage = ref(1);
const itemsPerPage = 10;

const totalPages = computed(() => 
  Math.ceil(props.data.length / itemsPerPage)
);

const paginatedData = computed(() => {
  const start = (currentPage.value - 1) * itemsPerPage;
  const end = start + itemsPerPage;
  return props.data.slice(start, end);
});

function formatTimestamp(timestamp: string): string {
  const date = new Date(timestamp);
  return date.toLocaleString('pt-BR');
}

function nextPage() {
  if (currentPage.value < totalPages.value) {
    currentPage.value++;
  }
}

function prevPage() {
  if (currentPage.value > 1) {
    currentPage.value--;
  }
}
</script>

<style scoped>
.telemetry-table-container {
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
  font-size: var(--fs-lg, 18px);
  font-weight: 600;
}

.spinner-small {
  display: inline-block;
  width: 14px;
  height: 14px;
  border: 2px solid rgba(0, 0, 0, 0.1);
  border-left-color: currentColor;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.table-wrapper {
  overflow-x: auto;
}

.telemetry-table {
  width: 100%;
  border-collapse: collapse;
  font-size: var(--fs-sm, 14px);
}

.telemetry-table th,
.telemetry-table td {
  padding: var(--sp-3, 12px);
  text-align: left;
  border-bottom: 1px solid var(--border-color, #e5e7eb);
}

.telemetry-table th {
  background-color: var(--bg-subtle, #f9fafb);
  font-weight: 600;
  color: var(--text-muted, #6b7280);
  text-transform: uppercase;
  font-size: var(--fs-xs, 12px);
}

.telemetry-table tbody tr:hover {
  background-color: var(--bg-hover, #f9fafb);
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: var(--sp-3, 12px);
  padding: var(--sp-3, 12px) 0;
}

.page-info {
  font-size: var(--fs-sm, 14px);
  color: var(--text-muted, #6b7280);
}

.no-data {
  text-align: center;
  padding: var(--sp-6, 24px);
  color: var(--text-muted, #6b7280);
}
</style>
