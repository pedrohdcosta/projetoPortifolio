<template>
  <div class="consumption-chart card">
    <div v-if="initialLoading || isEmpty" class="chart-placeholder">
      <SkeletonCard v-if="initialLoading" />
      <div v-else class="empty-state">
        <p class="text-muted">Nenhum dado de consumo disponível.</p>
      </div>
    </div>
    <div v-else class="chart-content">
      <div class="chart-header">
        <h3 class="chart-title">{{ title }}</h3>
        <span v-if="loading" class="updating-badge">Atualizando...</span>
      </div>
      <div class="chart-wrapper">
        <div class="y-axis">
          <span class="axis-label">{{ maxValue.toFixed(0) }} W</span>
          <span class="axis-label">{{ midValue.toFixed(0) }} W</span>
          <span class="axis-label">{{ minValue.toFixed(0) }} W</span>
        </div>
        <div class="chart-area">
          <canvas ref="chartCanvas" role="img" :aria-label="ariaLabel" />
          <div class="x-axis">
            <span v-for="(label, i) in xAxisLabels" :key="i" class="axis-label">{{ label }}</span>
          </div>
        </div>
      </div>
      <div class="chart-legend">
        <span class="legend-item">
          <span class="legend-color"></span>
          Consumo (Watts)
        </span>
      </div>
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
const hasLoadedOnce = ref(false)
const isEmpty = computed(() => props.labels.length === 0 || props.series.length === 0)
const initialLoading = computed(() => props.loading && !hasLoadedOnce.value)
const ariaLabel = computed(() => 
  `Gráfico de consumo com ${props.series.length} pontos de dados`
)

// Calculate Y-axis values using reduce for better performance with large arrays
const maxValue = computed(() => {
  if (props.series.length === 0) return 1
  return props.series.reduce((max, val) => Math.max(max, val), props.series[0])
})
const minValue = computed(() => {
  if (props.series.length === 0) return 0
  return props.series.reduce((min, val) => Math.min(min, val), props.series[0])
})
const midValue = computed(() => (maxValue.value + minValue.value) / 2)

// Calculate X-axis labels (show 5 labels evenly distributed)
const xAxisLabels = computed(() => {
  if (props.labels.length === 0) return []
  if (props.labels.length <= 5) return props.labels
  
  const step = Math.floor(props.labels.length / 4)
  return [
    props.labels[0],
    props.labels[step],
    props.labels[step * 2],
    props.labels[step * 3],
    props.labels[props.labels.length - 1]
  ]
})

// Canvas chart rendering
function renderChart() {
  if (!chartCanvas.value || isEmpty.value) return
  
  const canvas = chartCanvas.value
  const ctx = canvas.getContext('2d')
  if (!ctx) return

  // Set canvas size for retina
  const dpr = window.devicePixelRatio || 1
  const rect = canvas.getBoundingClientRect()
  canvas.width = rect.width * dpr
  canvas.height = rect.height * dpr
  ctx.scale(dpr, dpr)

  // Clear canvas
  ctx.clearRect(0, 0, rect.width, rect.height)

  const padding = { top: 10, right: 10, bottom: 10, left: 10 }
  const width = rect.width - padding.left - padding.right
  const height = rect.height - padding.top - padding.bottom
  
  const max = maxValue.value
  const min = minValue.value
  const range = max - min || 1

  // Pre-calculate all coordinates once for better performance
  const coordinates = props.series.map((value, index) => ({
    x: padding.left + (index / (props.series.length - 1 || 1)) * width,
    y: padding.top + height - ((value - min) / range) * height
  }))

  // Draw grid lines
  ctx.strokeStyle = 'rgba(255, 255, 255, 0.1)'
  ctx.lineWidth = 1
  for (let i = 0; i <= 4; i++) {
    const y = padding.top + (i / 4) * height
    ctx.beginPath()
    ctx.moveTo(padding.left, y)
    ctx.lineTo(padding.left + width, y)
    ctx.stroke()
  }

  // Draw line chart with gradient fill
  const gradient = ctx.createLinearGradient(0, padding.top, 0, padding.top + height)
  gradient.addColorStop(0, 'rgba(67, 160, 71, 0.4)')
  gradient.addColorStop(1, 'rgba(67, 160, 71, 0.05)')

  // Draw fill area
  ctx.beginPath()
  coordinates.forEach((coord, index) => {
    if (index === 0) {
      ctx.moveTo(coord.x, coord.y)
    } else {
      ctx.lineTo(coord.x, coord.y)
    }
  })
  // Close the path for fill
  ctx.lineTo(padding.left + width, padding.top + height)
  ctx.lineTo(padding.left, padding.top + height)
  ctx.closePath()
  ctx.fillStyle = gradient
  ctx.fill()

  // Draw line
  ctx.strokeStyle = 'rgba(67, 160, 71, 1)'
  ctx.lineWidth = 2
  ctx.lineCap = 'round'
  ctx.lineJoin = 'round'
  ctx.beginPath()
  
  coordinates.forEach((coord, index) => {
    if (index === 0) {
      ctx.moveTo(coord.x, coord.y)
    } else {
      ctx.lineTo(coord.x, coord.y)
    }
  })
  ctx.stroke()

  // Draw data points
  ctx.fillStyle = 'rgba(67, 160, 71, 1)'
  const pointRadius = props.series.length > 50 ? 0 : props.series.length > 20 ? 2 : 4
  if (pointRadius > 0) {
    coordinates.forEach((coord) => {
      ctx.beginPath()
      ctx.arc(coord.x, coord.y, pointRadius, 0, Math.PI * 2)
      ctx.fill()
    })
  }
}

// Mark as loaded once when data first arrives
watch(() => props.series, (newSeries) => {
  if (newSeries.length > 0) {
    hasLoadedOnce.value = true
  }
})

watch([() => props.labels, () => props.series], () => {
  requestAnimationFrame(renderChart)
}, { deep: true })

onMounted(() => {
  if (!isEmpty.value) {
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
  min-height: 350px;
}

.chart-placeholder {
  min-height: 300px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.chart-content {
  display: flex;
  flex-direction: column;
  gap: var(--sp-3, 12px);
}

.chart-header {
  display: flex;
  align-items: center;
  gap: var(--sp-3, 12px);
}

.chart-title {
  font-size: var(--fs-lg, 1.25rem);
  font-weight: 600;
  margin: 0;
}

.updating-badge {
  font-size: 0.75rem;
  padding: 2px 8px;
  background: rgba(67, 160, 71, 0.2);
  border-radius: 4px;
  color: rgba(67, 160, 71, 1);
  animation: pulse 1.5s infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.chart-wrapper {
  display: flex;
  gap: var(--sp-2, 8px);
}

.y-axis {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  padding: 10px 0;
  min-width: 60px;
  text-align: right;
}

.chart-area {
  flex: 1;
  display: flex;
  flex-direction: column;
}

canvas {
  width: 100%;
  height: 250px;
  border-radius: 8px;
  background: rgba(0, 0, 0, 0.2);
}

.x-axis {
  display: flex;
  justify-content: space-between;
  padding-top: var(--sp-2, 8px);
  padding-left: 10px;
  padding-right: 10px;
}

.axis-label {
  font-size: 0.75rem;
  color: var(--text-muted, rgba(255, 255, 255, 0.6));
}

.chart-legend {
  display: flex;
  justify-content: center;
  gap: var(--sp-4, 16px);
  padding-top: var(--sp-2, 8px);
}

.legend-item {
  display: flex;
  align-items: center;
  gap: var(--sp-2, 8px);
  font-size: 0.875rem;
  color: var(--text-muted, rgba(255, 255, 255, 0.7));
}

.legend-color {
  width: 12px;
  height: 12px;
  background: rgba(67, 160, 71, 1);
  border-radius: 2px;
}

.empty-state {
  padding: var(--sp-6, 24px);
  text-align: center;
}
</style>
