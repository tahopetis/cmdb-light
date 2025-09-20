import { defineStore } from 'pinia'
import ciService from '../services/ciService'
import relationshipService from '../services/relationshipService'

export const useCIStore = defineStore('ci', {
  state: () => ({
    cis: [],
    currentCI: null,
    ciTypes: [
      'Server',
      'Application',
      'Database',
      'Network Device',
      'Storage',
      'Service',
      'License',
      'Other'
    ],
    relationships: [],
    filters: {
      search: '',
      type: '',
      status: '',
      tags: []
    },
    pagination: {
      page: 1,
      limit: 10,
      total: 0,
      totalPages: 0
    },
    isLoading: false,
    isLoadingRelationships: false,
    error: null
  }),
  
  getters: {
    getCIsCount: (state) => {
      return state.cis.length
    },
    
    getCITypesCount: (state) => {
      const counts = {}
      
      state.ciTypes.forEach(type => {
        counts[type] = state.cis.filter(ci => ci.type === type).length
      })
      
      return counts
    },
    
    filteredCIs: (state) => {
      let result = [...state.cis]
      
      // Apply search filter
      if (state.filters.search) {
        const searchLower = state.filters.search.toLowerCase()
        result = result.filter(ci => 
          ci.name.toLowerCase().includes(searchLower) ||
          ci.description.toLowerCase().includes(searchLower) ||
          (ci.tags && ci.tags.some(tag => tag.toLowerCase().includes(searchLower)))
        )
      }
      
      // Apply type filter
      if (state.filters.type) {
        result = result.filter(ci => ci.type === state.filters.type)
      }
      
      // Apply status filter
      if (state.filters.status) {
        result = result.filter(ci => ci.status === state.filters.status)
      }
      
      // Apply tags filter
      if (state.filters.tags && state.filters.tags.length > 0) {
        result = result.filter(ci => 
          ci.tags && state.filters.tags.every(tag => ci.tags.includes(tag))
        )
      }
      
      return result
    },
    
    paginatedCIs: (state) => {
      const filtered = state.filteredCIs
      const start = (state.pagination.page - 1) * state.pagination.limit
      const end = start + state.pagination.limit
      
      return filtered.slice(start, end)
    }
  },
  
  actions: {
    async fetchCIs() {
      this.isLoading = true
      this.error = null
      
      try {
        const params = {
          page: this.pagination.page,
          limit: this.pagination.limit,
          search: this.filters.search,
          type: this.filters.type,
          status: this.filters.status,
          tags: this.filters.tags.join(',')
        }
        
        const response = await ciService.getCIs(params)
        
        this.cis = response.items || []
        this.pagination.total = response.total || 0
        this.pagination.totalPages = response.totalPages || 0
      } catch (error) {
        console.error('Error fetching CIs:', error)
        this.error = error.response?.data?.message || 'Failed to fetch configuration items'
        throw error
      } finally {
        this.isLoading = false
      }
    },
    
    async fetchCIById(id) {
      this.isLoading = true
      this.error = null
      
      try {
        const ci = await ciService.getCIById(id)
        this.currentCI = ci
      } catch (error) {
        console.error('Error fetching CI:', error)
        this.error = error.response?.data?.message || 'Failed to fetch configuration item'
        throw error
      } finally {
        this.isLoading = false
      }
    },
    
    async createCI(ciData) {
      this.isLoading = true
      this.error = null
      
      try {
        const newCI = await ciService.createCI(ciData)
        
        // Add new CI to the list
        this.cis.unshift(newCI)
        
        return { success: true, data: newCI }
      } catch (error) {
        console.error('Error creating CI:', error)
        this.error = error.response?.data?.message || 'Failed to create configuration item'
        return { success: false, error: this.error }
      } finally {
        this.isLoading = false
      }
    },
    
    async updateCI(id, ciData) {
      this.isLoading = true
      this.error = null
      
      try {
        const updatedCI = await ciService.updateCI(id, ciData)
        
        // Update CI in the list
        const index = this.cis.findIndex(ci => ci.id === id)
        if (index !== -1) {
          this.cis[index] = updatedCI
        }
        
        // Update current CI if it's the one being edited
        if (this.currentCI && this.currentCI.id === id) {
          this.currentCI = updatedCI
        }
        
        return { success: true, data: updatedCI }
      } catch (error) {
        console.error('Error updating CI:', error)
        this.error = error.response?.data?.message || 'Failed to update configuration item'
        return { success: false, error: this.error }
      } finally {
        this.isLoading = false
      }
    },
    
    async deleteCI(id) {
      this.isLoading = true
      this.error = null
      
      try {
        await ciService.deleteCI(id)
        
        // Remove CI from the list
        this.cis = this.cis.filter(ci => ci.id !== id)
        
        return { success: true }
      } catch (error) {
        console.error('Error deleting CI:', error)
        this.error = error.response?.data?.message || 'Failed to delete configuration item'
        return { success: false, error: this.error }
      } finally {
        this.isLoading = false
      }
    },
    
    async fetchCIRelationships(id) {
      this.isLoadingRelationships = true
      this.error = null
      
      try {
        const relationships = await relationshipService.getCIRelationships(id)
        this.relationships = relationships || []
      } catch (error) {
        console.error('Error fetching CI relationships:', error)
        this.error = error.response?.data?.message || 'Failed to fetch configuration item relationships'
        throw error
      } finally {
        this.isLoadingRelationships = false
      }
    },
    
    async createRelationship(relationshipData) {
      this.isLoading = true
      this.error = null
      
      try {
        const newRelationship = await relationshipService.createRelationship(relationshipData)
        
        // Add new relationship to the list
        this.relationships.push(newRelationship)
        
        return { success: true, data: newRelationship }
      } catch (error) {
        console.error('Error creating relationship:', error)
        this.error = error.response?.data?.message || 'Failed to create relationship'
        return { success: false, error: this.error }
      } finally {
        this.isLoading = false
      }
    },
    
    async updateRelationship(id, relationshipData) {
      this.isLoading = true
      this.error = null
      
      try {
        const updatedRelationship = await relationshipService.updateRelationship(id, relationshipData)
        
        // Update relationship in the list
        const index = this.relationships.findIndex(rel => rel.id === id)
        if (index !== -1) {
          this.relationships[index] = updatedRelationship
        }
        
        return { success: true, data: updatedRelationship }
      } catch (error) {
        console.error('Error updating relationship:', error)
        this.error = error.response?.data?.message || 'Failed to update relationship'
        return { success: false, error: this.error }
      } finally {
        this.isLoading = false
      }
    },
    
    async deleteRelationship(id) {
      this.isLoading = true
      this.error = null
      
      try {
        await relationshipService.deleteRelationship(id)
        
        // Remove relationship from the list
        this.relationships = this.relationships.filter(rel => rel.id !== id)
        
        return { success: true }
      } catch (error) {
        console.error('Error deleting relationship:', error)
        this.error = error.response?.data?.message || 'Failed to delete relationship'
        return { success: false, error: this.error }
      } finally {
        this.isLoading = false
      }
    },
    
    setFilters(filters) {
      this.filters = { ...this.filters, ...filters }
      this.pagination.page = 1 // Reset to first page when filters change
    },
    
    setPage(page) {
      this.pagination.page = page
    },
    
    setLimit(limit) {
      this.pagination.limit = limit
      this.pagination.page = 1 // Reset to first page when limit changes
    },
    
    resetFilters() {
      this.filters = {
        search: '',
        type: '',
        status: '',
        tags: []
      }
      this.pagination.page = 1
    }
  }
})