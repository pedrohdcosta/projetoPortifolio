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
      
      <div v-if="data.length === 0" class="empty-state">
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
          <tr v-for="item in data" :key="item.id">
            <td>{{ item.deviceName }}</td>
            <td>{{ item.power.toFixed(2) }}</td>
            <td>{{ formatTimestamp(item.timestamp) }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import SkeletonCard from './SkeletonCard.vue'

export interface TelemetryData {
  id: number
  deviceName: string
  power: number
  timestamp: string
}

defineProps<{
  data: TelemetryData[]
  loading?: boolean
}>()

defineEmits<{
  refresh: []
}>()

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
</style>
