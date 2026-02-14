<script lang="ts">
import { ref, computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useStockStore } from '../stores/StockStore'
import { formatCurrency } from '../utils/format'
import type { Stock } from '../models/stock'
import { storeToRefs } from 'pinia'


export default {
  setup() {

      const router = useRouter()
      const stockStore = useStockStore()

      const {stocks, currentStock, history, recommendation, loading, error} = storeToRefs(stockStore)


      const searchQuery = ref('')
      const sortBy = ref<'symbol' | 'current_price' | 'name'>('symbol')
      const activeTab = ref<'dashboard' | 'explore'>('explore')

      const watchlistOnly = computed(() => activeTab.value === 'dashboard')

      let debounceTimer: ReturnType<typeof setTimeout> | null = null
      watch([searchQuery, sortBy, activeTab], () => {
        if (debounceTimer) clearTimeout(debounceTimer)
        debounceTimer = setTimeout(fetchStocks, 300)
      }, { immediate: true })

      function fetchStocks() {
        stockStore.fetchStocks(searchQuery.value, sortBy.value, watchlistOnly.value)
      }

      function setTab(tab: 'dashboard' | 'explore') {
        activeTab.value = tab
      }

      async function toggleWatchlist(symbol: string) {
        await stockStore.toggleWatchlist(symbol)
        if (watchlistOnly.value) fetchStocks()
      }

      function goToDetail(symbol: string) {
        router.push(`/stock/${symbol}`)
      }

      function getChangeClass(stock: Stock) {
        const curr = stock.current_price
        const prev = stock.prev_close
        if (curr > prev) return 'text-emerald-500'
        if (curr < prev) return 'text-rose-500'
        return 'text-slate-500'
      }

      function getChangePercent(stock: Stock) {
        const curr = stock.current_price
        const prev = stock.prev_close
        if (prev === 0) return '0.00%'
        const pct = ((curr - prev) / prev) * 100
        const sign = pct >= 0 ? '+' : ''
        return `${sign}${pct.toFixed(2)}%`
      }

      return {stocks, currentStock, history, recommendation, loading, error, setTab, toggleWatchlist, goToDetail, getChangeClass, getChangePercent, formatCurrency, searchQuery, sortBy, activeTab}
  }
}

</script>

<template>
  <div class="animate-fade-in">
    <!-- Filters -->
    <div class="flex flex-col sm:flex-row gap-4 mb-8">
      <div class="relative flex-1 max-w-sm">
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Search symbols..."
          class="w-full h-10 pl-4 pr-4 bg-slate-900 border border-slate-700 rounded-lg text-slate-100 placeholder-slate-500 focus:outline-none focus:border-slate-600 focus:ring-1 focus:ring-slate-600 transition-colors"
        />
      </div>
      <div class="flex items-center gap-3">
        <span class="text-sm text-slate-500">Sort</span>
        <select
          v-model="sortBy"
          class="h-10 px-4 bg-slate-900 border border-slate-700 rounded-lg text-slate-300 text-sm focus:outline-none focus:border-slate-600 cursor-pointer"
        >
          <option value="symbol">Symbol</option>
          <option value="current_price">Price</option>
          <option value="name">Name</option>
        </select>
      </div>
    </div>

    <!-- Tabs -->
    <div class="flex gap-1 p-1 bg-slate-900/50 rounded-lg max-w-xs mb-10">
      <button
        :class="[
          'flex-1 py-2 text-sm font-medium rounded-md transition-colors',
          activeTab === 'dashboard'
            ? 'bg-slate-800 text-white'
            : 'text-slate-500 hover:text-slate-300',
        ]"
        @click="setTab('dashboard')"
      >
        Watchlist
      </button>
      <button
        :class="[
          'flex-1 py-2 text-sm font-medium rounded-md transition-colors',
          activeTab === 'explore'
            ? 'bg-slate-800 text-white'
            : 'text-slate-500 hover:text-slate-300',
        ]"
        @click="setTab('explore')"
      >
        All
      </button>
    </div>

    <!-- Loading -->
    <div
      v-if="loading && stocks.length === 0"
      class="flex flex-col items-center py-20 text-slate-500"
    >
      <div class="w-8 h-8 border-2 border-slate-700 border-t-slate-400 rounded-full animate-spin" />
      <p class="mt-4 text-sm">Loading market data...</p>
    </div>

    <!-- Error -->
    <div
      v-else-if="error"
      class="py-4 px-4 bg-rose-500/10 border border-rose-500/30 rounded-lg text-rose-400 text-sm"
    >
      {{ error }}
    </div>

    <!-- Empty -->
    <div
      v-else-if="stocks.length === 0"
      class="py-20 text-center"
    >
      <p class="text-slate-500 text-lg mb-2">No stocks found</p>
      <p class="text-slate-600 text-sm mb-6">
        {{ activeTab === 'dashboard' ? 'Add stocks from All to your watchlist' : 'Try adjusting your search' }}
      </p>
      <button
        v-if="activeTab === 'dashboard'"
        class="text-sm text-slate-400 hover:text-slate-300 transition-colors"
        @click="setTab('explore')"
      >
        Browse all stocks →
      </button>
    </div>

    <!-- Grid -->
    <div
      v-else
      class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4"
    >
      <article
        v-for="stock in stocks"
        :key="stock.symbol"
        class="group p-4 bg-slate-900/50 border border-slate-800 rounded-lg hover:border-slate-700 transition-stock cursor-pointer"
        @click="goToDetail(stock.symbol)"
      >
        <div class="flex justify-between items-start mb-3">
          <div>
            <h3 class="text-lg font-semibold text-slate-100 group-hover:text-white transition-colors">
              {{ stock.symbol }}
            </h3>
            <p class="text-xs text-slate-500 uppercase tracking-wider mt-0.5">
              {{ stock.type }}
            </p>
          </div>
          <button
            class="w-8 h-8 flex items-center justify-center rounded-lg border border-slate-700 text-slate-400 hover:bg-slate-800 hover:text-white transition-colors shrink-0 addButton"
            title="Toggle watchlist"
            :class="{active: stock.in_watchlist}"
            @click.stop="toggleWatchlist(stock.symbol)"
          >
            {{ stock.in_watchlist ? '✓' : '+' }}
          </button>
        </div>
        <p class="text-sm text-slate-500 line-clamp-1 mb-4">
          {{ stock.name }}
        </p>
        <div class="flex justify-between items-baseline">
          <span class="text-xl font-mono font-semibold font-mono-num text-slate-100">
            {{ formatCurrency(stock.current_price, stock.currency) }}
          </span>
          <span :class="['text-sm font-mono font-mono-num', getChangeClass(stock)]">
            {{ getChangePercent(stock) }}
          </span>
        </div>
        <div class="mt-4 pt-4 border-t border-slate-800 flex gap-2">
          <button
            class="flex-1 py-2 text-sm font-medium text-slate-400 hover:text-slate-200 bg-slate-800/50 hover:bg-slate-800 rounded-lg transition-colors"
            @click.stop="goToDetail(stock.symbol)"
          >
            Details
          </button>
        </div>
      </article>
    </div>
  </div>
</template>

<style scoped>
.addButton.active {
  background: #039e77;
  color: white;
}
</style>
