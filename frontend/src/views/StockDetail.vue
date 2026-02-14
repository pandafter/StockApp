<script lang="ts">
import { onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useStockStore } from '../stores/StockStore'
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler
} from 'chart.js'
import { Line as LineChart } from 'vue-chartjs'
import { formatCurrency } from '../utils/format'
import { storeToRefs } from 'pinia'

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler
)

export default {
  components: {
    LineChart
  },
  setup() {
        
    const route = useRoute()
    const router = useRouter()
    const stockStore = useStockStore()

    const {stocks, currentStock, history, recommendation, loading, error} = storeToRefs(stockStore)


    const symbol = route.params.symbol as string

    onMounted(() => {
      stockStore.fetchStockDetail(symbol)
    })

    const chartData = computed(() => {
      const hist = [...stockStore.history].sort(
        (a, b) => new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime()
      )
      return {
        labels: hist.map(p =>
          new Date(p.timestamp).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
        ),
        datasets: [
          {
            label: 'Price',
            backgroundColor: 'rgba(34, 197, 94, 0.1)',
            borderColor: '#22c55e',
            data: hist.map(p => p.price),
            fill: true,
            tension: 0.3,
            pointRadius: 0,
            borderWidth: 2,
          },
        ],
      }
    })

    const chartOptions = {
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: { display: false },
      },
      scales: {
        x: {
          display: true,
          grid: { color: 'rgba(148, 163, 184, 0.1)' },
          ticks: { color: '#64748b', font: { size: 11 } },
        },
        y: {
          display: true,
          grid: { color: 'rgba(148, 163, 184, 0.1)' },
          ticks: { color: '#64748b', font: { size: 11 } },
        },
      },
    }

    function getChangeClass(current: number, prev: number) {
      if (current > prev) return 'text-emerald-500'
      if (current < prev) return 'text-rose-500'
      return 'text-slate-500'
    }

    function getChangePercent(current: number, prev: number) {
      if (prev === 0) return '0.00%'
      const pct = ((current - prev) / prev) * 100
      const sign = pct >= 0 ? '+' : ''
      return `${sign}${pct.toFixed(2)}%`
    }

    async function toggleWatchlist() {
      if (currentStock.value) {
        await stockStore.toggleWatchlist(currentStock.value.symbol)
      }
    }

    return { stockStore, stocks, currentStock, history, recommendation, loading, error,router, formatCurrency, chartOptions, chartData, LineChart, getChangeClass, getChangePercent, toggleWatchlist}
  }
}

</script>

<template>
  <div class="animate-fade-in">
    <button
      class="mb-6 flex items-center gap-2 text-slate-500 hover:text-slate-300 text-sm font-medium transition-colors"
      @click="router.back()"
    >
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
      </svg>
      Back
    </button>

    <!-- Loading -->
    <div v-if="loading" class="flex flex-col items-center py-20 text-slate-500">
      <div class="w-8 h-8 border-2 border-slate-700 border-t-slate-400 rounded-full animate-spin" />
      <p class="mt-4 text-sm">Loading...</p>
    </div>

    <!-- Content -->
    <div v-else-if="currentStock" class="space-y-6">
      <!-- Header -->
      <div class="flex flex-col sm:flex-row sm:items-start sm:justify-between gap-4 p-6 bg-slate-900/50 border border-slate-800 rounded-lg">
        <div>
          <h1 class="text-3xl font-bold text-slate-100">{{ currentStock.symbol }}</h1>
          <p class="text-slate-500 mt-1">{{ currentStock.name }}</p>
          <span class="inline-block mt-2 px-2.5 py-0.5 bg-slate-800 text-slate-400 text-xs font-medium rounded uppercase">
            {{ currentStock.type }}
          </span>
        </div>
        <div class="text-left sm:text-right">
          <p class="text-3xl font-bold font-mono font-mono-num text-slate-100">
            {{ formatCurrency(currentStock.current_price, currentStock.currency) }}
          </p>
          <p :class="['text-sm font-mono font-mono-num mt-1', getChangeClass(currentStock.current_price, currentStock.prev_close)]">
            {{ getChangePercent(currentStock.current_price, currentStock.prev_close) }}
          </p>
          <p class="text-slate-600 text-xs mt-1">
            {{ new Date(currentStock.updated_at).toLocaleString() }}
          </p>
          <button
            :class="[
              'mt-3 px-4 py-2 text-sm font-medium rounded-lg transition-colors border',
              currentStock.in_watchlist
                ? 'bg-slate-700 border-slate-600 text-slate-300'
                : 'border-slate-600 text-slate-300 hover:bg-slate-800',
            ]"
            @click="toggleWatchlist"
          >
            {{ currentStock.in_watchlist ? '✓ In Watchlist' : '+ Add to Watchlist' }}
          </button>
        </div>
      </div>

      <!-- Stats -->
      <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
        <div class="p-4 bg-slate-900/50 border border-slate-800 rounded-lg">
          <p class="text-slate-500 text-xs font-medium uppercase mb-1">Open</p>
          <p class="font-mono font-semibold text-slate-200 font-mono-num">
            {{ formatCurrency(currentStock.open_price, currentStock.currency) }}
          </p>
        </div>
        <div class="p-4 bg-slate-900/50 border border-slate-800 rounded-lg">
          <p class="text-slate-500 text-xs font-medium uppercase mb-1">Prev Close</p>
          <p class="font-mono font-semibold text-slate-200 font-mono-num">
            {{ formatCurrency(currentStock.prev_close, currentStock.currency) }}
          </p>
        </div>
        <div class="p-4 bg-slate-900/50 border border-slate-800 rounded-lg">
          <p class="text-slate-500 text-xs font-medium uppercase mb-1">High</p>
          <p class="font-mono font-semibold text-emerald-500 font-mono-num">
            {{ formatCurrency(currentStock.high_price, currentStock.currency) }}
          </p>
        </div>
        <div class="p-4 bg-slate-900/50 border border-slate-800 rounded-lg">
          <p class="text-slate-500 text-xs font-medium uppercase mb-1">Low</p>
          <p class="font-mono font-semibold text-rose-500 font-mono-num">
            {{ formatCurrency(currentStock.low_price, currentStock.currency) }}
          </p>
        </div>
      </div>

      <!-- Chart -->
      <div class="p-6 bg-slate-900/50 border border-slate-800 rounded-lg">
        <h3 class="text-sm font-medium text-slate-400 mb-4">Price History</h3>
        <div class="h-72">
          <LineChart v-if="history.length > 0" :data="chartData" :options="chartOptions" />
          <div v-else class="flex items-center justify-center h-full text-slate-500 text-sm">
            No history available
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
