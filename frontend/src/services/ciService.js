import apiClient from './api'

class CIService {
  // Get all CIs with pagination and filtering
  async getCIs(params = {}) {
    try {
      const response = await apiClient.get('/cis', { params })
      return response.data
    } catch (error) {
      console.error('Error fetching CIs:', error)
      throw error
    }
  }

  // Get a single CI by ID
  async getCIById(id) {
    try {
      const response = await apiClient.get(`/cis/${id}`)
      return response.data
    } catch (error) {
      console.error(`Error fetching CI with ID ${id}:`, error)
      throw error
    }
  }

  // Create a new CI
  async createCI(ciData) {
    try {
      const response = await apiClient.post('/cis', ciData)
      return response.data
    } catch (error) {
      console.error('Error creating CI:', error)
      throw error
    }
  }

  // Update an existing CI
  async updateCI(id, ciData) {
    try {
      const response = await apiClient.put(`/cis/${id}`, ciData)
      return response.data
    } catch (error) {
      console.error(`Error updating CI with ID ${id}:`, error)
      throw error
    }
  }

  // Delete a CI
  async deleteCI(id) {
    try {
      await apiClient.delete(`/cis/${id}`)
      return { success: true }
    } catch (error) {
      console.error(`Error deleting CI with ID ${id}:`, error)
      throw error
    }
  }

  // Get CI relationships
  async getCIRelationships(id) {
    try {
      const response = await apiClient.get(`/cis/${id}/relationships`)
      return response.data
    } catch (error) {
      console.error(`Error fetching relationships for CI with ID ${id}:`, error)
      throw error
    }
  }
}

export default new CIService()