import axios from 'axios'

// API base URL with version prefix
// All API endpoints are versioned using path-based versioning (e.g., /api/v1/)
// This allows for future API versions without breaking existing client integrations
const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1'

class AuthService {
  constructor() {
    this.accessToken = localStorage.getItem('accessToken')
    this.refreshToken = localStorage.getItem('refreshToken')
    this.user = JSON.parse(localStorage.getItem('user') || 'null')
  }

  async login(username, password) {
    try {
      const response = await axios.post(`${API_URL}/auth/login`, {
        username,
        password
      })

      if (response.data.accessToken) {
        this.accessToken = response.data.accessToken
        this.refreshToken = response.data.refreshToken
        this.user = response.data.user
        localStorage.setItem('accessToken', this.accessToken)
        localStorage.setItem('refreshToken', this.refreshToken)
        localStorage.setItem('user', JSON.stringify(this.user))
        return { success: true, data: response.data }
      }
      
      return { success: false, message: 'Login failed' }
    } catch (error) {
      console.error('Login error:', error)
      return {
        success: false,
        message: error.response?.data?.message || 'Login failed'
      }
    }
  }

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
    
    this.accessToken = null
    this.refreshToken = null
    this.user = null
    localStorage.removeItem('accessToken')
    localStorage.removeItem('refreshToken')
    localStorage.removeItem('user')
  }

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
      
      return true
    } catch (error) {
      console.error('Token refresh error:', error)
      this.logout()
      return false
    }
  }

  getAccessToken() {
    return this.accessToken
  }

  getRefreshToken() {
    return this.refreshToken
  }

  getUser() {
    return this.user
  }

  isAuthenticated() {
    return !!this.accessToken
  }

  hasRole(role) {
    return this.user && this.user.role === role
  }

  isAdmin() {
    return this.hasRole('admin')
  }

  isViewer() {
    return this.hasRole('viewer')
  }
}

export default new AuthService()