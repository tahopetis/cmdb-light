import { defineStore } from 'pinia'
import axios from 'axios'

// API base URL with version prefix
// All API endpoints are versioned using path-based versioning (e.g., /api/v1/)
// This allows for future API versions without breaking existing client integrations
const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null,
    accessToken: localStorage.getItem('accessToken') || null,
    refreshToken: localStorage.getItem('refreshToken') || null,
    isAuthenticated: false,
    isLoading: false
  }),
  
  getters: {
    isAdmin: (state) => {
      return state.user && state.user.role === 'admin'
    },
    
    username: (state) => {
      return state.user ? state.user.username : ''
    }
  },
  
  actions: {
    async login(username, password) {
      this.isLoading = true
      
      try {
        const response = await axios.post(`${API_URL}/auth/login`, {
          username,
          password
        })
        
        const { accessToken, refreshToken, user } = response.data
        
        // Store tokens in localStorage
        localStorage.setItem('accessToken', accessToken)
        localStorage.setItem('refreshToken', refreshToken)
        
        // Update state
        this.accessToken = accessToken
        this.refreshToken = refreshToken
        this.user = user
        this.isAuthenticated = true
        
        return { success: true }
      } catch (error) {
        console.error('Login error:', error)
        
        let message = 'Login failed'
        if (error.response && error.response.data) {
          message = error.response.data.message || message
        }
        
        return { success: false, message }
      } finally {
        this.isLoading = false
      }
    },
    
    async logout() {
      try {
        // Call the logout endpoint to revoke the refresh token
        if (this.accessToken) {
          await axios.post(`${API_URL}/auth/logout`, {}, {
            headers: {
              Authorization: `Bearer ${this.accessToken}`
            }
          })
        }
      } catch (error) {
        console.error('Logout error:', error)
        // Continue with local logout even if the API call fails
      }
      
      // Remove tokens from localStorage
      localStorage.removeItem('accessToken')
      localStorage.removeItem('refreshToken')
      
      // Reset state
      this.accessToken = null
      this.refreshToken = null
      this.user = null
      this.isAuthenticated = false
    },
    
    checkAuth() {
      // Check if access token exists and is not expired
      if (!this.accessToken) {
        this.isAuthenticated = false
        return false
      }
      
      try {
        // Decode JWT token to check expiration
        const payload = JSON.parse(atob(this.accessToken.split('.')[1]))
        const currentTime = Math.floor(Date.now() / 1000)
        
        if (payload.exp < currentTime) {
          // Access token expired, try to refresh
          if (this.refreshToken) {
            this.refreshAccessToken()
            return false // Will be updated after refresh
          } else {
            // No refresh token, logout
            this.logout()
            return false
          }
        }
        
        // Set user from token payload
        this.user = {
          id: payload.sub,
          username: payload.username,
          role: payload.role
        }
        
        this.isAuthenticated = true
        return true
      } catch (error) {
        console.error('Token validation error:', error)
        this.logout()
        return false
      }
    },
    
    async refreshAccessToken() {
      if (!this.refreshToken) {
        this.logout()
        return false
      }
      
      try {
        const response = await axios.post(`${API_URL}/auth/refresh`, {
          refreshToken: this.refreshToken
        })
        
        const { accessToken, refreshToken } = response.data
        
        // Store new tokens in localStorage
        localStorage.setItem('accessToken', accessToken)
        localStorage.setItem('refreshToken', refreshToken)
        
        // Update state
        this.accessToken = accessToken
        this.refreshToken = refreshToken
        
        // Update user info from the new access token
        const payload = JSON.parse(atob(accessToken.split('.')[1]))
        this.user = {
          id: payload.sub,
          username: payload.username,
          role: payload.role
        }
        
        this.isAuthenticated = true
        return true
      } catch (error) {
        console.error('Token refresh error:', error)
        this.logout()
        return false
      }
    },
    
    async fetchUserProfile() {
      if (!this.accessToken) return
      
      try {
        const response = await axios.get(`${API_URL}/auth/validate`, {
          headers: {
            Authorization: `Bearer ${this.accessToken}`
          }
        })
        
        this.user = response.data
        this.isAuthenticated = true
      } catch (error) {
        console.error('Error fetching user profile:', error)
        this.logout()
      }
    }
  }
})