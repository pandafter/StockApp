export const formatCurrency = (value: number, currency: string): string => {
  if (value === undefined || value === null) return '—'
  return new Intl.NumberFormat('en-US', { style: 'currency', currency, minimumFractionDigits: 2 }).format(value)
}

export const formatPercent = (value: number): string => {
  if (value === undefined || value === null) return '—'
  const sign = value >= 0 ? '+' : ''
  return `${sign}${value.toFixed(2)}%`
}

export const formatChange = (current: number, prev: number): { value: number; percent: number } => {
  if (prev === 0) return { value: 0, percent: 0 }
  const value = current - prev
  const percent = (value / prev) * 100
  return { value, percent }
}
