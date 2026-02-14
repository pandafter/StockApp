import { defineStore } from 'pinia'
import apiClient from '../utils/api'
import type { Recommendation, Stock, StockPrice } from '../models/stock'

export const useStockStore = defineStore('stock', {
  state: () => ({
    stocks: [] as Stock[],
    currentStock: null as Stock | null,
    history: [] as StockPrice[],
    recommendation: null as Recommendation | null,

    loading: false,
    error: null as string | null,
  }),
  getters: {
    hasError: (state) => state.error !== null,
  },
  actions: {
    async fetchStocks(search = '', sortBy = '', watchlistOnly = false) {
      this.loading = true
      this.error = null
      try {
        const params = new URLSearchParams()
        if (search) params.append('search', search)
        if (sortBy) params.append('sort_by', sortBy)
        if (watchlistOnly) params.append('watchlist', 'true')

        const { data } = await apiClient.get<Stock[]>(`/stocks?${params.toString()}`)
        this.stocks = data
      } catch (err: unknown) {
        this.error = err instanceof Error ? err.message : 'Failed to fetch stocks'
      } finally {
        this.loading = false
      }
    },

    async fetchStockDetail(symbol: string) {
      this.loading = true
      this.error = null
      try {
        const { data } = await apiClient.get<{ stock: Stock; history: StockPrice[] }>(`/stocks/${symbol}`)
        this.currentStock = data.stock
        this.history = data.history
      } catch (err: unknown) {
        this.error = err instanceof Error ? err.message : 'Failed to fetch stock detail'
      } finally {
        this.loading = false
      }
    },

    async fetchRecommendation() {
      if (this.recommendation) return
      this.loading = true
      try {
        const { data } = await apiClient.get<Recommendation>('/recommendation')
        this.recommendation = data
      } catch (err: unknown) {
        this.error = err instanceof Error ? err.message : 'Failed to fetch recommendation'
      } finally {
        this.loading = false
      }
    },

    async toggleWatchlist(symbol: string) {
      this.error = null
      try {
        const { data } = await apiClient.post<Stock>(`/stocks/${symbol}/watchlist`)

        // Update stocks list
        const index = this.stocks.findIndex(s => s.symbol === symbol)
        if (index !== -1) {
          this.stocks[index] = data
        }

        // Update current stock if viewing detail
        if (this.currentStock?.symbol === symbol) {
          this.currentStock = data
        }

        // Update recommendation if it's the recommended stock
        if (this.recommendation?.recommendation.symbol === symbol) {
          this.recommendation.recommendation = data
        }
      } catch (err: unknown) {
        const errorMessage = err instanceof Error ? err.message : 'Failed to toggle watchlist'
        this.error = errorMessage
        console.error('Toggle watchlist error:', err)
      }
    }
  }
})