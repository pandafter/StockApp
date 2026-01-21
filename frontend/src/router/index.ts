import { createRouter, createWebHistory } from 'vue-router'
import StockList from '../components/StockList.vue'
import Recommendation from '../components/Recommendation.vue'

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [
        {
            path: '/',
            name: 'home',
            component: StockList
        },
        {
            path: '/recommendation',
            name: 'recommendation',
            component: Recommendation
        },
        {
            path: '/stock/:symbol',
            name: 'stock-detail',
            component: () => import('../views/StockDetail.vue')
        }
    ]
})

export default router
