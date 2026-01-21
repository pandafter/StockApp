export const formatCurrency = (value: number, currency: string) => {
    if (value === undefined || value === null) return '-'
    return new Intl.NumberFormat('en-US', { style: 'currency', currency: currency }).format(value)
}
