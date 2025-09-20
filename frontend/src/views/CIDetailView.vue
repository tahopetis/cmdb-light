<template>
  <div class="py-6">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <!-- Header with title and actions -->
      <div class="flex justify-between items-center">
        <h1 class="text-2xl font-semibold text-gray-900">Configuration Item Details</h1>
        <div class="flex space-x-3">
          <router-link
            v-if="authStore.isAdmin"
            :to="`/cis/${ciId}/edit`"
            class="btn btn-secondary"
          >
            Edit
          </router-link>
          <button
            v-if="authStore.isAdmin"
            @click="confirmDelete"
            class="btn btn-primary"
          >
            Delete
          </button>
        </div>
      </div>
      <!-- End Header -->
    </div>
    
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 mt-6">
      <div v-if="ciStore.isLoading" class="flex justify-center py-12">
        <svg class="animate-spin h-8 w-8 text-primary-600" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
      </div>
      
      <div v-else-if="ciStore.error" class="card">
        <div class="px-4 py-5 sm:p-6">
          <div class="text-center py-12">
            <svg class="mx-auto h-12 w-12 text-red-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
            <h3 class="mt-2 text-sm font-medium text-gray-900">Error loading configuration item</h3>
            <p class="mt-1 text-sm text-gray-500">{{ ciStore.error }}</p>
            <div class="mt-6">
              <button @click="loadCI" class="btn btn-primary">
                Try Again
              </button>
            </div>
          </div>
        </div>
      </div>
      
      <div v-else-if="!ciStore.currentCI" class="card">
        <div class="px-4 py-5 sm:p-6">
          <div class="text-center py-12">
            <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
            <h3 class="mt-2 text-sm font-medium text-gray-900">Configuration item not found</h3>
            <p class="mt-1 text-sm text-gray-500">The requested configuration item could not be found.</p>
            <div class="mt-6">
              <router-link to="/cis" class="btn btn-primary">
                Back to Configuration Items
              </router-link>
            </div>
          </div>
        </div>
      </div>
      
      <!-- Main Content Grid -->
      <div v-else class="grid grid-cols-1 gap-6 lg:grid-cols-3">
        <!-- Main Information -->
        <div class="lg:col-span-2">
          <div class="card">
            <div class="px-4 py-5 sm:p-6">
              <div class="flex items-start justify-between">
                <div>
                  <h2 class="text-xl font-semibold text-gray-900">{{ ciStore.currentCI.name }}</h2>
                  <div class="mt-1 flex items-center space-x-4">
                    <span
                      class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                      :class="getStatusClass(ciStore.currentCI.status)"
                    >
                      {{ ciStore.currentCI.status }}
                    </span>
                    <span class="text-sm text-gray-500">{{ ciStore.currentCI.type }}</span>
                  </div>
                </div>
                
                <div
                  class="flex-shrink-0 h-12 w-12 rounded-md flex items-center justify-center"
                  :class="getTypeColorClass(ciStore.currentCI.type)"
                >
                  <svg
                    class="h-6 w-6 text-white"
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"
                    />
                  </svg>
                </div>
              </div>
              
              <div class="mt-6">
                <h3 class="text-sm font-medium text-gray-900">Description</h3>
                <div class="mt-1 text-sm text-gray-600">
                  {{ ciStore.currentCI.description || 'No description available' }}
                </div>
              </div>
              
              <div v-if="ciStore.currentCI.tags && ciStore.currentCI.tags.length > 0" class="mt-6">
                <h3 class="text-sm font-medium text-gray-900">Tags</h3>
                <div class="mt-2 flex flex-wrap gap-2">
                  <span
                    v-for="(tag, index) in ciStore.currentCI.tags"
                    :key="index"
                    class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-primary-100 text-primary-800"
                  >
                    {{ tag }}
                  </span>
                </div>
              </div>
              
              <div v-if="ciStore.currentCI.custom_attributes && ciStore.currentCI.custom_attributes.length > 0" class="mt-6">
                <h3 class="text-sm font-medium text-gray-900">Custom Attributes</h3>
                <div class="mt-2 grid grid-cols-1 gap-4 sm:grid-cols-2">
                  <div
                    v-for="(attr, index) in ciStore.currentCI.custom_attributes"
                    :key="index"
                    class="bg-gray-50 px-4 py-3 rounded-md"
                  >
                    <dt class="text-sm font-medium text-gray-500">{{ attr.key }}</dt>
                    <dd class="mt-1 text-sm text-gray-900">{{ attr.value }}</dd>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
        <!-- End Main Information -->
        
        <!-- Sidebar Information -->
        <div class="lg:col-span-1">
          <div class="card">
            <div class="px-4 py-5 sm:p-6">
              <h3 class="text-sm font-medium text-gray-900">Details</h3>
              <dl class="mt-2 space-y-2">
                <div>
                  <dt class="text-sm font-medium text-gray-500">ID</dt>
                  <dd class="mt-1 text-sm text-gray-900">{{ ciStore.currentCI.id }}</dd>
                </div>
                <div>
                  <dt class="text-sm font-medium text-gray-500">Created</dt>
                  <dd class="mt-1 text-sm text-gray-900">{{ formatDate(ciStore.currentCI.created_at) }}</dd>
                </div>
                <div>
                  <dt class="text-sm font-medium text-gray-500">Last Updated</dt>
                  <dd class="mt-1 text-sm text-gray-900">{{ formatDate(ciStore.currentCI.updated_at) }}</dd>
                </div>
              </dl>
            </div>
          </div>
          
          <div class="mt-6 card">
            <div class="px-4 py-5 sm:p-6">
              <h3 class="text-sm font-medium text-gray-900">Actions</h3>
              <div class="mt-4 space-y-3">
                <router-link
                  :to="`/cis`"
                  class="btn btn-secondary w-full"
                >
                  Back to Configuration Items
                </router-link>
                <router-link
                  v-if="authStore.isAdmin"
                  :to="`/cis/${ciId}/edit`"
                  class="btn btn-secondary w-full"
                >
                  Edit Configuration Item
                </router-link>
                <button
                  v-if="authStore.isAdmin"
                  @click="confirmDelete"
                  class="btn btn-primary w-full"
                >
                  Delete Configuration Item
                </button>
              </div>
            </div>
          </div>
        </div>
        <!-- End Sidebar Information -->
        
        <!-- Relationships Section -->
        <div class="lg:col-span-3 mt-6">
          <div class="card">
            <div class="px-4 py-5 sm:p-6">
              <div class="flex justify-between items-center">
                <h3 class="text-lg leading-6 font-medium text-gray-900">Relationships</h3>
                <button
                  v-if="authStore.isAdmin"
                  @click="showAddRelationship = true"
                  class="btn btn-secondary"
                >
                  Add Relationship
                </button>
              </div>
              
              <div v-if="ciStore.isLoading && ciStore.isLoadingRelationships" class="flex justify-center py-6">
                <svg class="animate-spin h-6 w-6 text-primary-600" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
              </div>
              
              <div v-else-if="ciStore.relationships.length === 0" class="text-center py-6">
                <svg class="mx-auto h-12 w-12 text-gray-400" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
                </svg>
                <h3 class="mt-2 text-sm font-medium text-gray-900">No relationships</h3>
                <p class="mt-1 text-sm text-gray-500">This configuration item has no relationships with other items.</p>
              </div>
              
              <div v-else class="mt-4">
                <div class="flow-root">
                  <ul class="-mb-8">
                    <li v-for="(relationship, index) in ciStore.relationships" :key="relationship.id">
                      <div class="relative pb-8">
                        <span v-if="index !== ciStore.relationships.length - 1" class="absolute top-4 left-4 -ml-px h-full w-0.5 bg-gray-200" aria-hidden="true"></span>
                        <div class="relative flex space-x-3">
                          <div>
                            <span class="h-8 w-8 rounded-full bg-primary-100 flex items-center justify-center ring-8 ring-white">
                              <svg class="h-5 w-5 text-primary-600" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" />
                              </svg>
                            </span>
                          </div>
                          <div class="min-w-0 flex-1 pt-1.5 flex justify-between space-x-4">
                            <div>
                              <p class="text-sm text-gray-900 font-medium">
                                <span v-if="isSource(relationship)" class="text-primary-600">
                                  {{ relationship.type }} →
                                </span>
                                <span v-else class="text-green-600">
                                  ← {{ relationship.type }}
                                </span>
                                <router-link
                                  :to="`/cis/${getRelatedCIId(relationship)}`"
                                  class="ml-1 hover:underline"
                                >
                                  {{ getRelatedCIName(relationship) }}
                                </router-link>
                              </p>
                              <p class="text-sm text-gray-500">
                                {{ getRelatedCIType(relationship) }}
                              </p>
                            </div>
                            <div class="text-right text-sm whitespace-nowrap text-gray-500">
                              <span>{{ formatDate(relationship.created_at) }}</span>
                              <button
                                v-if="authStore.isAdmin"
                                @click="confirmDeleteRelationship(relationship)"
                                class="ml-4 text-red-600 hover:text-red-900"
                                title="Delete relationship"
                              >
                                <svg class="h-4 w-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                                </svg>
                              </button>
                            </div>
                          </div>
                        </div>
                      </div>
                    </li>
                  </ul>
                </div>
              </div>
            </div>
          </div>
        </div>
        <!-- End Relationships Section -->
      </div>
    </div>
    
    <!-- Add Relationship Modal -->
    <div v-if="showAddRelationship" class="fixed z-10 inset-0 overflow-y-auto">
      <div class="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
        <div class="fixed inset-0 transition-opacity" aria-hidden="true">
          <div class="absolute inset-0 bg-gray-500 opacity-75"></div>
        </div>

        <span class="hidden sm:inline-block sm:align-middle sm:h-screen" aria-hidden="true">&#8203;</span>

        <div class="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full">
          <RelationshipForm
            :preselected-source="ciStore.currentCI.id"
            @submit="handleRelationshipAdded"
            @cancel="showAddRelationship = false"
          />
        </div>
      </div>
    </div>
    <!-- End Add Relationship Modal -->
  </div>
</template>

<script>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useCIStore } from '../stores/ci'
import { useAuthStore } from '../stores/auth'
import { useUIStore } from '../stores/ui'
import RelationshipForm from '../components/forms/RelationshipForm.vue'

export default {
  name: 'CIDetailView',
  components: {
    RelationshipForm
  },
  setup() {
    const route = useRoute()
    const router = useRouter()
    const ciStore = useCIStore()
    const authStore = useAuthStore()
    const uiStore = useUIStore()
    
    const showAddRelationship = ref(false)
    const relatedCIs = ref({})
    
    const ciId = computed(() => route.params.id)
    
    const loadCI = async () => {
      try {
        uiStore.setLoading(true)
        await ciStore.fetchCIById(ciId.value)
        await ciStore.fetchCIRelationships(ciId.value)
      } catch (error) {
        uiStore.showError('Failed to load configuration item')
      } finally {
        uiStore.setLoading(false)
      }
    }
    
    const confirmDelete = () => {
      uiStore.showConfirmDialog(
        'Delete Configuration Item',
        'Are you sure you want to delete this configuration item? This action cannot be undone.',
        async () => {
          try {
            uiStore.setLoading(true)
            const result = await ciStore.deleteCI(ciId.value)
            
            if (result.success) {
              uiStore.showSuccess('Configuration item deleted successfully')
              router.push('/cis')
            } else {
              uiStore.showError(result.error || 'Failed to delete configuration item')
            }
          } catch (error) {
            console.error('Error deleting CI:', error)
            uiStore.showError('An unexpected error occurred')
          } finally {
            uiStore.setLoading(false)
          }
        }
      )
    }
    
    const confirmDeleteRelationship = (relationship) => {
      uiStore.showConfirmDialog(
        'Delete Relationship',
        'Are you sure you want to delete this relationship? This action cannot be undone.',
        async () => {
          try {
            uiStore.setLoading(true)
            const result = await ciStore.deleteRelationship(relationship.id)
            
            if (result.success) {
              uiStore.showSuccess('Relationship deleted successfully')
            } else {
              uiStore.showError(result.error || 'Failed to delete relationship')
            }
          } catch (error) {
            console.error('Error deleting relationship:', error)
            uiStore.showError('Failed to delete relationship')
          } finally {
            uiStore.setLoading(false)
          }
        }
      )
    }
    
    const handleRelationshipAdded = async () => {
      showAddRelationship.value = false
      
      try {
        uiStore.setLoading(true)
        await ciStore.fetchCIRelationships(ciId.value)
        uiStore.showSuccess('Relationship added successfully')
      } catch (error) {
        console.error('Error refreshing relationships:', error)
        uiStore.showError('Failed to refresh relationships')
      } finally {
        uiStore.setLoading(false)
      }
    }
    
    const isSource = (relationship) => {
      return relationship.source_id === ciStore.currentCI.id
    }
    
    const getRelatedCIId = (relationship) => {
      return isSource(relationship) ? relationship.target_id : relationship.source_id
    }
    
    const getRelatedCIName = (relationship) => {
      const relatedCI = relatedCIs.value[getRelatedCIId(relationship)]
      return relatedCI ? relatedCI.name : 'Unknown CI'
    }
    
    const getRelatedCIType = (relationship) => {
      const relatedCI = relatedCIs.value[getRelatedCIId(relationship)]
      return relatedCI ? relatedCI.type : 'Unknown Type'
    }
    
    const loadRelatedCIs = async () => {
      if (!ciStore.relationships || ciStore.relationships.length === 0) return
      
      const relatedCIIds = new Set()
      ciStore.relationships.forEach(relationship => {
        relatedCIIds.add(relationship.source_id)
        relatedCIIds.add(relationship.target_id)
      })
      
      // Remove current CI ID
      relatedCIIds.delete(ciStore.currentCI.id)
      
      // Fetch related CIs
      for (const id of relatedCIIds) {
        if (!relatedCIs.value[id]) {
          try {
            const ci = await ciStore.getCIById(id)
            relatedCIs.value[id] = ci
          } catch (error) {
            console.error(`Error loading CI with ID ${id}:`, error)
          }
        }
      }
    }
    
    const getStatusClass = (status) => {
      const classes = {
        'Active': 'bg-green-100 text-green-800',
        'Inactive': 'bg-gray-100 text-gray-800',
        'Pending': 'bg-yellow-100 text-yellow-800',
        'Maintenance': 'bg-blue-100 text-blue-800',
        'Retired': 'bg-red-100 text-red-800'
      }
      return classes[status] || classes['Active']
    }
    
    const getTypeColorClass = (type) => {
      const classes = {
        'Server': 'bg-blue-500',
        'Application': 'bg-green-500',
        'Database': 'bg-purple-500',
        'Network Device': 'bg-yellow-500',
        'Storage': 'bg-indigo-500',
        'Service': 'bg-pink-500',
        'License': 'bg-red-500',
        'Other': 'bg-gray-500'
      }
      return classes[type] || classes['Other']
    }
    
    const formatDate = (dateString) => {
      if (!dateString) return 'N/A'
      const date = new Date(dateString)
      return date.toLocaleDateString()
    }
    
    onMounted(() => {
      loadCI()
    })
    
    // Watch for changes in relationships to load related CIs
    watch(() => ciStore.relationships, loadRelatedCIs, { immediate: true })
    
    // Watch for route changes to reload data when navigating between CIs
    watch(() => route.params.id, (newId) => {
      if (newId) {
        loadCI()
      }
    })
    
    return {
      ciStore,
      authStore,
      uiStore,
      ciId,
      showAddRelationship,
      loadCI,
      confirmDelete,
      confirmDeleteRelationship,
      handleRelationshipAdded,
      isSource,
      getRelatedCIId,
      getRelatedCIName,
      getRelatedCIType,
      getStatusClass,
      getTypeColorClass,
      formatDate
    }
  }
}
</script>