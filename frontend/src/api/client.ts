import axios from 'axios'
import router from '@/router'
import { useAuthStore } from '@/stores/auth'

const apiClient = axios.create({
  baseURL: '/api',
  withCredentials: true,
  headers: {
    'Content-Type': 'application/json'
  }
})

let isRefreshing = false
let failedQueue: Array<{ resolve: Function; reject: Function }> = []

const processQueue = (error: any = null) => {
  failedQueue.forEach(prom => {
    if (error) {
      prom.reject(error)
    } else {
      prom.resolve()
    }
  })
  failedQueue = []
}

apiClient.interceptors.response.use(
  response => response,
  async error => {
    const originalRequest = error.config

    if (error.response?.status === 401 && !originalRequest._retry) {
      // Don't retry if this was the refresh endpoint itself
      if (originalRequest.url?.includes('/auth/refresh')) {
        const authStore = useAuthStore()
        authStore.clearAuth()
        // Don't redirect to login - just let the request fail silently
        // The router guard will handle redirects for protected routes
        return Promise.reject(error)
      }

      if (isRefreshing) {
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject })
        }).then(() => {
          return apiClient(originalRequest)
        }).catch(err => {
          return Promise.reject(err)
        })
      }

      originalRequest._retry = true
      isRefreshing = true

      try {
        await apiClient.post('/auth/refresh', {})

        // Refresh succeeded - update auth store
        const authStore = useAuthStore()
        try {
          await authStore.fetchUser()
        } catch (e) {
          // If fetching user fails after refresh, something is wrong
          console.error('Failed to fetch user after refresh:', e)
        }

        processQueue()
        isRefreshing = false

        return apiClient(originalRequest)
      } catch (refreshError) {
        processQueue(refreshError)
        isRefreshing = false

        const authStore = useAuthStore()
        authStore.clearAuth()

        // Only redirect to login if we're trying to access a protected route
        // Check if the original request was for a protected endpoint
        const protectedEndpoints = ['/auth/me', '/books/favorites', '/books/downloads', '/books/upload', '/books/favorite', '/books/download', '/books/ghost-mode', '/books/delete', '/books/metadata']
        const isProtectedRequest = protectedEndpoints.some(endpoint => originalRequest.url?.includes(endpoint))

        if (isProtectedRequest && router.currentRoute.value.meta.requiresAuth) {
          router.push('/login')
        }

        return Promise.reject(refreshError)
      }
    }

    return Promise.reject(error)
  }
)

export default apiClient
