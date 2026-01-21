<script setup lang="ts">
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
import { Line } from 'vue-chartjs'

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

const route = useRoute()
const router = useRouter()
const store = useStockStore()
const symbol = route.params.symbol as string

import { formatCurrency } from '../utils/format'

onMounted(() => {
  store.fetchStockDetail(symbol)
})

const chartData = computed(() => {
  const hist = [...store.history].sort((a, b) => new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime())
  
  return {
    labels: hist.map(p => new Date(p.timestamp).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })),
    datasets: [
      {
        label: 'Price',
        backgroundColor: 'rgba(79, 70, 229, 0.2)',
        borderColor: '#4f46e5',
        data: hist.map(p => p.price),
        fill: true,
        tension: 0.4
      }
    ]
  }
})

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: { display: false }
  },
  scales: {
    x: { display: true },
    y: { display: true }
  }
}
</script>

<template>
  <div>
    <button @click="router.back()" class="mb-6 flex items-center text-indigo-600 hover:text-indigo-800 font-medium">
      <svg class="w-5 h-5 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"></path></svg>
      Back
    </button>

    <div v-if="store.loading" class="text-center py-10">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600 mx-auto"></div>
    </div>

    <div v-else-if="store.currentStock" class="space-y-6">
      <!-- Header -->
      <div class="bg-white p-6 rounded-lg shadow-sm border border-gray-100 flex justify-between items-start">
        <div>
          <h1 class="text-4xl font-bold text-gray-900">{{ store.currentStock.symbol }}</h1>
          <p class="text-xl text-gray-500">{{ store.currentStock.name }}</p>
          <span class="inline-block mt-2 px-3 py-1 bg-gray-100 rounded-full text-sm text-gray-600 uppercase">{{ store.currentStock.type }}</span>
        </div>
        <div class="text-right">
          <p class="text-5xl font-bold text-indigo-600">{{ formatCurrency(store.currentStock.current_price, store.currentStock.currency) }}</p>
          <p class="text-gray-500 mt-1">Last Updated: {{ new Date(store.currentStock.updated_at).toLocaleString() }}</p>
        </div>
      </div>

      <!-- Stats Grid -->
      <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
        <div class="bg-white p-4 rounded-lg shadow-sm border border-gray-100 text-center">
          <p class="text-gray-500 text-sm uppercase">Open</p>
          <p class="text-xl text-indigo-600 font-bold">{{ formatCurrency(store.currentStock.open_price, store.currentStock.currency) }}</p>
        </div>
        <div class="bg-white p-4 rounded-lg shadow-sm border border-gray-100 text-center">
           <p class="text-gray-500 text-sm uppercase">Prev Close</p>
          <p class="text-xl text-indigo-600 font-bold">{{ formatCurrency(store.currentStock.prev_close, store.currentStock.currency) }}</p>
        </div>
         <div class="bg-white p-4 rounded-lg shadow-sm border border-gray-100 text-center">
           <p class="text-gray-500 text-sm uppercase">High</p>
          <p class="text-xl font-bold text-green-600">{{ formatCurrency(store.currentStock.high_price, store.currentStock.currency) }}</p>
        </div>
         <div class="bg-white p-4 rounded-lg shadow-sm border border-gray-100 text-center">
           <p class="text-gray-500 text-sm uppercase">Low</p>
          <p class="text-xl font-bold text-red-600">{{ formatCurrency(store.currentStock.low_price, store.currentStock.currency) }}</p>
        </div>
      </div>

      <!-- Chart -->
      <div class="bg-white p-6 rounded-lg shadow-md h-96">
        <h3 class="text-lg font-bold mb-4 text-gray-800">Price History</h3>
        <Line v-if="store.history.length > 0" :data="chartData" :options="chartOptions" />
        <div v-else class="flex items-center justify-center h-full text-gray-400 italic">
          No history data available for this session.
        </div>
      </div>
    </div>
  </div>
</template>
