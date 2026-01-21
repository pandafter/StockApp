<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useStockStore } from '../stores/StockStore'
import { useRouter } from 'vue-router'
import { formatCurrency } from '../utils/format'

const store = useStockStore()
const router = useRouter()

const searchQuery = ref('')
const sortBy = ref('symbol')
const activeTab = ref('dashboard')


onMounted(() => {
  store.fetchStocks('', sortBy.value, activeTab.value === 'dashboard')
})

const changeTab = (tab: string) => {
  activeTab.value = tab
  store.fetchStocks(searchQuery.value, sortBy.value, activeTab.value === 'dashboard')
}

let timeout: any = null
watch([searchQuery, sortBy], () => {
  if (timeout) clearTimeout(timeout)
  timeout = setTimeout(() => {
    store.fetchStocks(searchQuery.value, sortBy.value, activeTab.value === 'dashboard')
  }, 300)
})

const toggleWatchlist = async (symbol: string) => {
  await store.toggleWatchlist(symbol)
  if (activeTab.value === 'dashboard') {
    store.fetchStocks(searchQuery.value, sortBy.value, true)
  }
}

const getChangeColor = (current: number, prev: number) => {
  if (current > prev) return 'text-green-400'
  if (current < prev) return 'text-red-400'
  return 'text-gray-400'
}



</script>

<template>
  <div class="animate-fade-in max-w-7xl mx-auto px-4 py-8">
    <div class="flex flex-col lg:flex-row justify-between items-start lg:items-center gap-8 mb-12">
        <div class="relative flex-grow sm:w-80 bg-white rounded-xl">
          <input 
            v-model="searchQuery" 
            type="text"
            placeholder="Search Symbols & More" 
            class="glass-input w-full pl-10 text-indigo-500 focus:border-indigo-600 focus:outline-2 focus:outline-offset-2 focus:outline-indigo-600 active:outline-indigo-600 placeholder:text-indigo-600 focus:placeholder:text-indigo-300 caret-indigo-600 "
          />
        </div>
      <div class="flex  gap-2 items-center">
          <p class="text-indigo-700">Sort By:</p>
          <select v-model="sortBy" class="bg-indigo-600 text-white p-3 cursor-pointer custom-select">
            <option class="bg-indigo-600 text-white" value="symbol">Symbol</option>
            <option class="bg-indigo-600 text-white" value="current_price">Price</option>
            <option class="bg-indigo-600 text-white" value="name">Name</option>
          </select>
        </div>

    </div>
    

    <div class="flex justify-center mb-10">
      <div class="p-1 glass-card-transparent  flex gap-1 shadow-2xl">
        <button 
          @click="changeTab('dashboard')" 
          :class="['px-8 py-2.5 rounded-xl font-bold transition-all', activeTab === 'dashboard' ? 'bg-indigo-600 text-white shadow-lg' : 'text-indigo-700 hover:text-white']"
        >
          My Dashboard
        </button>
        <button 
          @click="changeTab('explore')" 
          :class="['px-8 py-2.5 rounded-xl font-bold transition-all', activeTab === 'explore' ? 'bg-indigo-600 text-white shadow-lg' : 'text-indigo-700 hover:text-white']"
        >
          Explore All
        </button>
      </div>
    </div>

    <div v-if="store.loading && store.stocks.length === 0" class="flex flex-col items-center py-20">
      <div class="loader mb-4"></div>
      <p class="text-sm font-medium text-gray-400">Synchronizing market data...</p>
    </div>

    <div v-else-if="store.error" class="glass-alert-error mb-12">
      {{ store.error }}
    </div>

    <div v-else-if="store.stocks.length === 0" class="glass-card text-center py-20 flex flex-col items-center">
      <div class="text-5xl mb-6">📉</div>
      <h3 class="text-2xl font-bold mb-2">No results matched your view</h3>
      <p class="text-white max-w-xs mb-8">Try adjusting your filters or search for something new in the <span class="text-indigo-200">"Explore All"</span>"Explore" tab.</p>
      <button v-if="activeTab === 'dashboard'" @click="changeTab('explore')" class="text-indigo-200 font-bold hover:underline">
        Go to Market Explorer →
      </button>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
      <div 
        v-for="stock in store.stocks" 
        :key="stock.symbol" 
        class="glass-card hover-glow group flex flex-col justify-between h-full"
      >
        <div>
          <div class="flex justify-between items-start mb-6">
            <div @click="router.push(`/stock/${stock.symbol}`)" class="cursor-pointer">
              <h3 class="text-2xl font-bold tracking-tighter group-hover:text-indigo-600 transition-colors">{{ stock.symbol }}</h3>
              <p class="text-xs text-white uppercase font-bold tracking-widest mt-1">{{ stock.type }}</p>
            </div>
            <div class="text-right">
              <p class="text-3xl font-mono font-bold leading-none">{{ formatCurrency(stock.current_price, stock.currency) }}</p>
              <p :class="['text-[10px] font-bold mt-1', getChangeColor(stock.current_price, stock.prev_close)]">
                {{ stock.current_price >= stock.prev_close ? '▲' : '▼' }} 
                {{ Math.abs(((stock.current_price - stock.prev_close) / stock.prev_close) * 100).toFixed(2) }}%
              </p>
            </div>
          </div>
          <div class="bg-indigo-400">

          <div class="flex items-start bg-indigo-800 w-1 h-5 absolute">
          </div>
          <p class="text-sm text-white font-medium mb-6 line-clamp-1 border-l-2 border-indigo-500/30 pl-3">
            {{ stock.name }}
          </p>
          </div>

        </div>

        <div class="flex gap-2 pt-4 border-t border-white/5">
          <router-link 
            :to="'/stock/' + stock.symbol" 
            class="action-btn flex-grow text-center py-2.5 bg-white/5 text-white hover:bg-white/10"
          >
            Monitor
          </router-link>
          
          <button 
            @click.stop="toggleWatchlist(stock.symbol)"
            :class="['action-btn px-4 py-2.5', stock.in_watchlist ? 'bg-indigo-600 text-white' : 'bg-indigo-600 text-white border border-indigo-500/30 text-indigo-400 hover:bg-indigo-500 hover:text-white']"
          >
            {{ stock.in_watchlist ? '✓' : '+' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
