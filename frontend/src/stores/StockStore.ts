import { defineStore } from 'pinia'
import axios from 'axios'

export interface Stock {
    id: string
    symbol: string
    name: string
    type: string
    currency: string
    current_price: number
    high_price: number
    low_price: number
    open_price: number
    prev_close: number
    in_watchlist: boolean
    updated_at: string
}

export interface StockPrice {
    id: number
    stock_id: string
    price: number
    timestamp: string
}

export interface Recommendation {
    recommendation: Stock
    potential_gain_percent: number
    reason: string
}

export const useStockStore = defineStore('stock', {
    state: () => ({
        stocks: [] as Stock[],
        currentStock: null as Stock | null,
        history: [] as StockPrice[],
        recommendation: null as Recommendation | null,
        loading: false,
        error: null as string | null,
    }),
    actions: {
        watchEffect() {
            this.fetchStocks(),
                this.fetchRecommendation()
        },

        async fetchStocks(search: string = '', sortBy: string = '', watchlistOnly: boolean = false) {
            this.loading = true
            this.error = null
            try {
                const params = new URLSearchParams()
                if (search) params.append('search', search)
                if (sortBy) params.append('sort_by', sortBy)
                if (watchlistOnly) params.append('watchlist', 'true')

                const response = await axios.get(`http://localhost:8081/api/stocks?${params.toString()}`)
                this.stocks = response.data
            } catch (err: any) {
                this.error = err.message || 'Failed to fetch stocks'
            } finally {
                this.loading = false
            }
        },

        async fetchStockDetail(symbol: string) {
            this.loading = true
            this.error = null
            try {
                const response = await axios.get(`http://localhost:8081/api/stocks/${symbol}`)
                this.currentStock = response.data.stock
                this.history = response.data.history
            } catch (err: any) {
                this.error = err.message || 'Failed to fetch stock detail'
            } finally {
                this.loading = false
            }
        },

        async fetchRecommendation() {
            if (this.recommendation) return
            this.loading = true
            try {
                const response = await axios.get(`http://localhost:8081/api/recommendation`)
                this.recommendation = response.data
            } catch (err: any) {
                this.error = err.message || 'Failed to fetch recommendation'
            } finally {
                this.loading = false
            }
        },

        async toggleWatchlist(symbol: string) {
            try {
                const response = await axios.post(`http://localhost:8081/api/stocks/${symbol}/watchlist`)
                const updatedStock = response.data


                const index = this.stocks.findIndex(s => s.symbol === symbol)
                if (index !== -1) {
                    this.stocks[index] = updatedStock
                }


                if (this.recommendation && this.recommendation.recommendation.symbol === symbol) {
                    this.recommendation.recommendation = updatedStock
                }
            } catch (err: any) {
                this.error = err.message || 'Failed to toggle watchlist'
            }
        }
    }
})
