<script setup lang="ts">
import { onMounted } from 'vue'
import { useStockStore } from '../stores/StockStore'
import { formatCurrency } from '../utils/format'

const store = useStockStore()

onMounted(() => {
  store.fetchRecommendation()
})

const toggleWatchlist = async (symbol: string) => {
  await store.toggleWatchlist(symbol)
}
</script>

<template>
  <div class="max-w-4xl mx-auto px-4 py-8 animate-fade-in">
    <div class="mb-10">
      <h2 class="text-xs font-bold uppercase tracking-widest text-indigo-400 mb-2">Today's Alpha Selection</h2>
      <h1 class="text-4xl md:text-5xl font-black bg-gradient-to-r from-white to-indigo-500 bg-clip-text text-indigo-800">
        Top Investment Pick
      </h1>
    </div>

    <div v-if="store.loading && !store.recommendation" class="flex justify-center py-20">
      <div class="loader"></div>
    </div>

    <div v-else-if="store.error" class="glass-alert-error">
      {{ store.error }}
    </div>

    <div v-else-if="store.recommendation" class="glass-card hover-glow overflow-hidden relative">
      <div class="absolute top-0 right-0 p-8 opacity-10 pointer-events-none">
        <svg class="w-48 h-48 fill-white" viewBox="0 0 24 24"><path d="M12 2L4.5 20.29l.71.71L12 18l6.79 3 .71-.71L12 2z"/></svg>
      </div>

      <div class="flex flex-col md:flex-row gap-10 items-center">
        <div class="flex-grow">
          <div class="flex items-center gap-3 mb-6">
            <span class="px-3 py-1 bg-indigo-500/20 text-white-700 rounded-full text-[10px] font-bold uppercase tracking-wider border border-indigo-500/20">
              High Confidence
            </span>
          </div>

          <div class="mb-8">
            <h2 class="text-7xl font-black tracking-tighter mb-1">{{ store.recommendation.recommendation.symbol }}</h2>
            <p class="text-2xl text-indigo-700 font-medium tracking-tight">{{ store.recommendation.recommendation.name }}</p>
          </div>

          <div class="space-y-4 mb-8">
            <div class="p-4 bg-white/5 rounded-xl border border-white/10">
              <h4 class="text-indigo-700 text-xs font-bold uppercase mb-2">Investment Thesis</h4>
              <p class="text-white italic text-lg leading-relaxed">"{{ store.recommendation.reason }}"</p>
            </div>
          </div>
        </div>

        <div class="w-full md:w-72 flex flex-col gap-4">
          <div class="p-8 bg-gradient-to-br from-green-500/10 to-emerald-500/90 rounded-3xl border border-green-500/20 text-center">
            <p class="text-white-400 text-xs uppercase tracking-widest font-bold mb-2">Potential ROI</p>
            <p class="text-6xl font-black text-green-400">+{{ store.recommendation.potential_gain_percent.toFixed(1) }}%</p>
          </div>

          <div class="p-6 bg-white/5 rounded-2xl border border-white/5 space-y-3">
            <div class="flex justify-between text-sm">
              <span class="text-white-500">Target</span>
              <span class="font-mono text-white">{{ formatCurrency(store.recommendation.recommendation.high_price, store.recommendation.recommendation.currency) }}</span>
            </div>
            <div class="flex justify-between text-sm">
              <span class="text-white-500">Current</span>
              <span class="font-mono text-white">{{ formatCurrency(store.recommendation.recommendation.current_price, store.recommendation.recommendation.currency) }}</span>
            </div>
          </div>
        </div>
      </div>

      <div class="mt-10 flex flex-col sm:flex-row gap-4">
        <router-link 
          :to="'/stock/' + store.recommendation.recommendation.symbol" 
          class="flex-grow action-btn py-4 bg-white text-black hover:bg-gray-200 text-center"
        >
          View Full Analysis
        </router-link>
        
        <button 
          @click="toggleWatchlist(store.recommendation.recommendation.symbol)"
          :class="['px-8 py-4 action-btn border border-indigo-500/50', store.recommendation.recommendation.in_watchlist ? 'bg-indigo-600 border-indigo-600 text-white' : 'text-indigo-700 bg-indigo-700/10 hover:bg-indigo-500 hover:text-white']"
        >
          {{ store.recommendation.recommendation.in_watchlist ? 'Joined Dashboard' : 'Add to Dashboard' }}
        </button>
      </div>
    </div>
  </div>
</template>
