import axios from 'axios'

const apiClient = axios.create({
    baseURL: '/api',
    headers: {
        'Content-Type': 'application/json',
    },
    withCredentials: false,
})

// Request interceptor
apiClient.interceptors.request.use(
    (config) => {
        return config
    },
    (error) => {
        return Promise.reject(error)
    }
)

// Response interceptor
apiClient.interceptors.response.use(
    (response) => response,
    (error) => {
        if (error.response) {
            // Server responded with error status
            const { status, data } = error.response
            console.error(`API Error ${status}:`, data)

            if (status === 403) {
                console.error('Forbidden: Check API permissions and CORS configuration')
            }
        } else if (error.request) {
            // Request made but no response
            console.error('Network error: No response received from server')
        } else {
            console.error('Error:', error.message)
        }
        return Promise.reject(error)
    }
)

export default apiClient
