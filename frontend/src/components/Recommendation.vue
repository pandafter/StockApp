<script lang="ts">
import { onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useStockStore } from '../stores/StockStore'
import { formatCurrency } from '../utils/format'
import { storeToRefs } from 'pinia';

export default {
  setup() {
    const router = useRouter()
    const stockStore = useStockStore()

    const {stocks, currentStock, history, recommendation, loading, error} = storeToRefs(stockStore)

    const rec = computed(() => stockStore.recommendation)


    onMounted(() => {
      stockStore.fetchRecommendation()
    })

    async function toggleWatchlist(symbol: string) {
      await stockStore.toggleWatchlist(symbol)
    }

    function goToDetail(symbol: string) {
      router.push(`/stock/${symbol}`)
    }

    return { stockStore, stocks, currentStock, history, recommendation, loading, error,router, rec, toggleWatchlist, goToDetail, formatCurrency}
  }
}


</script>

<template>
  <div class="max-w-2xl animate-fade-in">
    <div class="mb-10">
      <p class="text-xs font-medium uppercase tracking-widest text-slate-500 mb-1">
        Top Pick
      </p>
      <h1 class="text-3xl font-bold text-slate-100">
        Investment Recommendation
      </h1>
    </div>

    <!-- Loading -->
    <div
      v-if="loading && !recommendation"
      class="flex flex-col items-center py-20 text-slate-500"
    >
      <div class="w-8 h-8 border-2 border-slate-700 border-t-slate-400 rounded-full animate-spin" />
      <p class="mt-4 text-sm">Analyzing market data...</p>
    </div>

    <!-- Error -->
    <div
      v-else-if="error"
      class="py-4 px-4 bg-rose-500/10 border border-rose-500/30 rounded-lg text-rose-400 text-sm"
    >
      {{ error }}
    </div>

    <!-- Card -->
    <div
      v-else-if="rec"
      class="border border-slate-800 rounded-lg bg-slate-900/50 overflow-hidden"
    >
      <div class="p-6 sm:p-8">
        <div class="flex flex-col sm:flex-row gap-8">
          <div class="flex-1">
            <div class="inline-flex items-center px-2.5 py-1 rounded bg-emerald-500/10 text-emerald-400 text-xs font-medium mb-4">
              High Potential
            </div>
            <h2 class="text-4xl font-bold text-slate-100 mb-1">
              {{ rec.recommendation.symbol }}
            </h2>
            <p class="text-slate-500 mb-6">
              {{ rec.recommendation.name }}
            </p>
            <div class="p-4 bg-slate-800/50 rounded-lg border border-slate-700/50">
              <p class="text-slate-400 text-sm mb-1">Thesis</p>
              <p class="text-slate-200 italic">"{{ rec.reason }}"</p>
            </div>
          </div>
          <div class="sm:w-48 flex flex-col gap-4">
            <div class="p-6 bg-emerald-500/10 border border-emerald-500/20 rounded-lg text-center">
              <p class="text-emerald-500/80 text-xs font-medium uppercase tracking-wider mb-1">
                Potential
              </p>
              <p class="text-3xl font-bold text-emerald-400 font-mono-num">
                +{{ rec.potential_gain_percent.toFixed(1) }}%
              </p>
            </div>
            <div class="p-4 bg-slate-800/50 rounded-lg border border-slate-700/50 space-y-2 text-sm">
              <div class="flex justify-between">
                <span class="text-slate-500">Current</span>
                <span class="font-mono text-slate-200 font-mono-num">
                  {{ formatCurrency(rec.recommendation.current_price, rec.recommendation.currency) }}
                </span>
              </div>
              <div class="flex justify-between">
                <span class="text-slate-500">Target</span>
                <span class="font-mono text-slate-200 font-mono-num">
                  {{ formatCurrency(rec.recommendation.high_price, rec.recommendation.currency) }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
      <div class="px-6 sm:px-8 py-4 bg-slate-800/30 border-t border-slate-800 flex flex-col sm:flex-row gap-3">
        <button
          class="flex-1 py-3 px-4 bg-slate-100 text-slate-900 font-medium rounded-lg hover:bg-white transition-colors text-sm"
          @click="goToDetail(rec.recommendation.symbol)"
        >
          View Analysis
        </button>
        <button
          :class="[
            'py-3 px-6 rounded-lg font-medium text-sm transition-colors border',
            rec.recommendation.in_watchlist
              ? 'bg-slate-700 border-slate-600 text-slate-300'
              : 'border-slate-600 text-slate-300 hover:bg-slate-800',
          ]"
          @click="toggleWatchlist(rec.recommendation.symbol)"
        >
          {{ rec.recommendation.in_watchlist ? 'In Watchlist' : 'Add to Watchlist' }}
        </button>
      </div>
    </div>
  </div>
</template>
