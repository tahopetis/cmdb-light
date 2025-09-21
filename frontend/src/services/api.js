import axios from 'axios'
import { useAuthStore } from '../stores/auth'

// API base URL with version prefix
// All API endpoints are versioned using path-based versioning (e.g., /api/v1/)
// This allows for future API versions without breaking existing client integrations
const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1'

// Create axios instance
const apiClient = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json'
  }
})

// Request interceptor to add JWT token to requests
apiClient.interceptors.request.use(
  (config) => {
    const authStore = useAuthStore()
    if (authStore.accessToken) {
      config.headers.Authorization = `Bearer ${authStore.accessToken}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor to handle common errors
apiClient.interceptors.response.use(
  (response) => {
    return response
  },
  async (error) => {
    const originalRequest = error.config
    
    // Handle 401 Unauthorized errors
    if (error.response && error.response.status === 401 && !originalRequest._retry) {
      const authStore = useAuthStore()
      
      // Try to refresh the token
      if (authStore.refreshToken) {
        originalRequest._retry = true
        
        try {
          await authStore.refreshAccessToken()
          
          // Retry the original request with the new token
          originalRequest.headers.Authorization = `Bearer ${authStore.accessToken}`
          return apiClient(originalRequest)
        } catch (refreshError) {
          // If refresh fails, logout and redirect to login
          authStore.logout()
          window.location.href = '/login'
          return Promise.reject(refreshError)
        }
      } else {
        // No refresh token, logout and redirect to login
        authStore.logout()
        window.location.href = '/login'
      }
    }
    
    return Promise.reject(error)
  }
)

export default apiClient