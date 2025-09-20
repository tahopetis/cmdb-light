import { defineStore } from 'pinia'
import axios from 'axios'

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null,
    token: localStorage.getItem('token') || null,
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
        
        const { token, user } = response.data
        
        // Store token in localStorage
        localStorage.setItem('token', token)
        
        // Update state
        this.token = token
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
    
    logout() {
      // Remove token from localStorage
      localStorage.removeItem('token')
      
      // Reset state
      this.token = null
      this.user = null
      this.isAuthenticated = false
    },
    
    checkAuth() {
      // Check if token exists and is not expired
      if (!this.token) {
        this.isAuthenticated = false
        return false
      }
      
      try {
        // Decode JWT token to check expiration
        const payload = JSON.parse(atob(this.token.split('.')[1]))
        const currentTime = Math.floor(Date.now() / 1000)
        
        if (payload.exp < currentTime) {
          // Token expired
          this.logout()
          return false
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
    
    async fetchUserProfile() {
      if (!this.token) return
      
      try {
        const response = await axios.get(`${API_URL}/auth/profile`, {
          headers: {
            Authorization: `Bearer ${this.token}`
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