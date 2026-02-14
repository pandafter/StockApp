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
