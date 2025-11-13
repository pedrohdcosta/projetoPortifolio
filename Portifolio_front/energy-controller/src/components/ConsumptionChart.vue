<template>
  <div class="consumption-chart card">
    <div v-if="loading || isEmpty" class="chart-placeholder">
      <SkeletonCard v-if="loading" />
      <div v-else class="empty-state">
        <p class="text-muted">Nenhum dado de consumo disponível.</p>
      </div>
    </div>
    <div v-else class="chart-content">
      <h3 class="chart-title">{{ title }}</h3>
      <canvas ref="chartCanvas" role="img" :aria-label="ariaLabel" />
      <p class="text-muted small chart-note">
        Gráfico mostrando consumo ao longo do tempo
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch, onMounted, onUnmounted } from 'vue'
import SkeletonCard from './SkeletonCard.vue'

const props = defineProps<{
  labels: string[]
  series: number[]
  title?: string
  loading?: boolean
}>()

const chartCanvas = ref<HTMLCanvasElement | null>(null)
const isEmpty = computed(() => props.labels.length === 0 || props.series.length === 0)
const ariaLabel = computed(() => 
  `Gráfico de consumo com ${props.series.length} pontos de dados`
)

// Simple canvas rendering (placeholder for real chart library)
function renderChart() {
  if (!chartCanvas.value || isEmpty.value || props.loading) return
  
  const canvas = chartCanvas.value
  const ctx = canvas.getContext('2d')
  if (!ctx) return

  // Set canvas size
  canvas.width = canvas.offsetWidth * 2
  canvas.height = 300 * 2
  ctx.scale(2, 2)

  // Clear canvas
  ctx.clearRect(0, 0, canvas.width, canvas.height)

  // Draw simple line chart
  const padding = 40
  const width = canvas.width / 2 - padding * 2
  const height = canvas.height / 2 - padding * 2
  
  const max = Math.max(...props.series, 1)
  const min = Math.min(...props.series, 0)
  const range = max - min || 1

  ctx.strokeStyle = 'rgba(67, 160, 71, 0.8)'
  ctx.lineWidth = 2
  ctx.beginPath()

  props.series.forEach((value, index) => {
    const x = padding + (index / (props.series.length - 1 || 1)) * width
    const y = padding + height - ((value - min) / range) * height
    
    if (index === 0) {
      ctx.moveTo(x, y)
    } else {
      ctx.lineTo(x, y)
    }
  })

  ctx.stroke()

  // Draw axes
  ctx.strokeStyle = 'rgba(255, 255, 255, 0.2)'
  ctx.lineWidth = 1
  ctx.beginPath()
  ctx.moveTo(padding, padding)
  ctx.lineTo(padding, padding + height)
  ctx.lineTo(padding + width, padding + height)
  ctx.stroke()
}

watch([() => props.labels, () => props.series, () => props.loading], () => {
  if (!props.loading) {
    setTimeout(renderChart, 100)
  }
})

onMounted(() => {
  if (!props.loading && !isEmpty.value) {
    renderChart()
  }
  window.addEventListener('resize', renderChart)
})

onUnmounted(() => {
  window.removeEventListener('resize', renderChart)
})
</script>

<style scoped>
.consumption-chart {
  padding: var(--sp-5, 20px);
  min-height: 300px;
}

.chart-placeholder {
  min-height: 250px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.chart-content {
  display: flex;
  flex-direction: column;
  gap: var(--sp-3, 12px);
}

.chart-title {
  font-size: var(--fs-lg, 1.25rem);
  font-weight: 600;
  margin: 0;
}

canvas {
  width: 100%;
  height: 300px;
  border-radius: 8px;
  background: rgba(0, 0, 0, 0.2);
}

.chart-note {
  margin-top: var(--sp-2, 8px);
  font-size: 0.875rem;
}

.empty-state {
  padding: var(--sp-6, 24px);
  text-align: center;
}
</style>
