<template>
  <div class="consumption-chart">
    <SkeletonCard v-if="loading || !hasData" />
    <div v-else class="chart-container">
      <canvas ref="chartCanvas"></canvas>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted, computed } from 'vue';
import {
  Chart,
  LineController,
  LineElement,
  PointElement,
  LinearScale,
  TimeScale,
  Title,
  Tooltip,
  Legend,
  CategoryScale,
  Filler,
} from 'chart.js';
import SkeletonCard from './SkeletonCard.vue';

// Register Chart.js components
Chart.register(
  LineController,
  LineElement,
  PointElement,
  LinearScale,
  TimeScale,
  CategoryScale,
  Title,
  Tooltip,
  Legend,
  Filler
);

interface ChartData {
  timestamp: string;
  power: number;
}

const props = defineProps<{
  data: ChartData[];
  loading?: boolean;
}>();

const chartCanvas = ref<HTMLCanvasElement | null>(null);
let chartInstance: Chart | null = null;

const hasData = computed(() => props.data && props.data.length > 0);

function createChart() {
  if (!chartCanvas.value || !hasData.value) return;

  const labels = props.data.map((d) => {
    const date = new Date(d.timestamp);
    return date.toLocaleTimeString('pt-BR', { hour: '2-digit', minute: '2-digit', second: '2-digit' });
  });
  const values = props.data.map((d) => d.power);

  chartInstance = new Chart(chartCanvas.value, {
    type: 'line',
    data: {
      labels,
      datasets: [
        {
          label: 'Potência (W)',
          data: values,
          borderColor: 'rgb(59, 130, 246)',
          backgroundColor: 'rgba(59, 130, 246, 0.1)',
          fill: true,
          tension: 0.4,
          pointRadius: 2,
          pointHoverRadius: 5,
        },
      ],
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: {
          display: true,
          position: 'top',
        },
        tooltip: {
          mode: 'index',
          intersect: false,
        },
      },
      scales: {
        x: {
          display: true,
          title: {
            display: true,
            text: 'Tempo',
          },
        },
        y: {
          display: true,
          title: {
            display: true,
            text: 'Potência (W)',
          },
          beginAtZero: true,
        },
      },
      interaction: {
        mode: 'nearest',
        axis: 'x',
        intersect: false,
      },
    },
  });
}

function destroyChart() {
  if (chartInstance) {
    chartInstance.destroy();
    chartInstance = null;
  }
}

watch(
  () => [props.data, props.loading],
  () => {
    if (!props.loading && hasData.value) {
      destroyChart();
      setTimeout(() => createChart(), 0);
    }
  },
  { deep: true }
);

onMounted(() => {
  if (!props.loading && hasData.value) {
    createChart();
  }
});

onUnmounted(() => {
  destroyChart();
});
</script>

<style scoped>
.consumption-chart {
  width: 100%;
  min-height: 300px;
}

.chart-container {
  width: 100%;
  height: 300px;
  position: relative;
}
</style>
