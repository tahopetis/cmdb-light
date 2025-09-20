import apiClient from './api'

class RelationshipService {
  // Get all relationships
  async getRelationships() {
    try {
      const response = await apiClient.get('/relationships')
      return response.data
    } catch (error) {
      console.error('Error fetching relationships:', error)
      throw error
    }
  }

  // Get a single relationship by ID
  async getRelationshipById(id) {
    try {
      const response = await apiClient.get(`/relationships/${id}`)
      return response.data
    } catch (error) {
      console.error(`Error fetching relationship with ID ${id}:`, error)
      throw error
    }
  }

  // Create a new relationship
  async createRelationship(relationshipData) {
    try {
      const response = await apiClient.post('/relationships', relationshipData)
      return response.data
    } catch (error) {
      console.error('Error creating relationship:', error)
      throw error
    }
  }

  // Update an existing relationship
  async updateRelationship(id, relationshipData) {
    try {
      const response = await apiClient.put(`/relationships/${id}`, relationshipData)
      return response.data
    } catch (error) {
      console.error(`Error updating relationship with ID ${id}:`, error)
      throw error
    }
  }

  // Delete a relationship
  async deleteRelationship(id) {
    try {
      await apiClient.delete(`/relationships/${id}`)
      return { success: true }
    } catch (error) {
      console.error(`Error deleting relationship with ID ${id}:`, error)
      throw error
    }
  }

  // Get relationships for a specific CI
  async getRelationshipsByCI(ciId) {
    try {
      const response = await apiClient.get(`/cis/${ciId}/relationships`)
      return response.data
    } catch (error) {
      console.error(`Error fetching relationships for CI with ID ${ciId}:`, error)
      throw error
    }
  }

  // Get graph data for visualization
  async getGraphData(params = {}) {
    try {
      const response = await apiClient.get('/relationships/graph', { params })
      return response.data
    } catch (error) {
      console.error('Error fetching graph data:', error)
      throw error
    }
  }
}

export default new RelationshipService()