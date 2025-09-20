<template>
  <div class="py-6">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between items-center">
        <h1 class="text-2xl font-semibold text-gray-900">Configuration Items</h1>
        <router-link v-if="authStore.isAdmin" to="/cis/create" class="btn btn-primary">
          Add New CI
        </router-link>
      </div>
    </div>
    
    <!-- Filters -->
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 mt-6">
      <div class="card">
        <div class="px-4 py-5 sm:p-6">
          <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
            <div>
              <label for="search" class="block text-sm font-medium text-gray-700">Search</label>
              <div class="mt-1">
                <input
                  type="text"
                  id="search"
                  v-model="filters.search"
                  @input="applyFilters"
                  class="form-input"
                  placeholder="Search by name or description"
                />
              </div>
            </div>
            
            <div>
              <label for="type" class="block text-sm font-medium text-gray-700">Type</label>
              <div class="mt-1">
                <select
                  id="type"
                  v-model="filters.type"
                  @change="applyFilters"
                  class="form-input"
                >
                  <option value="">All Types</option>
                  <option
                    v-for="type in ciStore.ciTypes"
                    :key="type"
                    :value="type"
                  >
                    {{ type }}
                  </option>
                </select>
              </div>
            </div>
            
            <div>
              <label for="status" class="block text-sm font-medium text-gray-700">Status</label>
              <div class="mt-1">
                <select
                  id="status"
                  v-model="filters.status"
                  @change="applyFilters"
                  class="form-input"
                >
                  <option value="">All Statuses</option>
                  <option value="Active">Active</option>
                  <option value="Inactive">Inactive</option>
                  <option value="Pending">Pending</option>
                  <option value="Maintenance">Maintenance</option>
                  <option value="Retired">Retired</option>
                </select>
              </div>
            </div>
            
            <div class="flex items-end">
              <button
                @click="resetFilters"
                class="btn btn-secondary w-full"
              >
                Reset Filters
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- CI List -->
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 mt-6">
      <div class="card">
        <div class="px-4 py-5 sm:p-6">
          <div v-if="ciStore.isLoading" class="flex justify-center py-12">
            <svg class="animate-spin h-8 w-8 text-primary-600" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
          </div>
          
          <div v-else-if="ciStore.error" class="text-center py-12">
            <svg class="mx-auto h-12 w-12 text-red-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
            <h3 class="mt-2 text-sm font-medium text-gray-900">Error loading configuration items</h3>
            <p class="mt-1 text-sm text-gray-500">{{ ciStore.error }}</p>
            <div class="mt-6">
              <button @click="loadCIs" class="btn btn-primary">
                Try Again
              </button>
            </div>
          </div>
          
          <div v-else-if="ciStore.filteredCIs.length === 0" class="text-center py-12">
            <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
            </svg>
            <h3 class="mt-2 text-sm font-medium text-gray-900">No configuration items found</h3>
            <p class="mt-1 text-sm text-gray-500">Try adjusting your filters or create a new configuration item.</p>
            <div v-if="authStore.isAdmin" class="mt-6">
              <router-link to="/cis/create" class="btn btn-primary">
                Create Configuration Item
              </router-link>
            </div>
          </div>
          
          <div v-else class="flex flex-col">
            <div class="-my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
              <div class="py-2 align-middle inline-block min-w-full sm:px-6 lg:px-8">
                <div class="shadow overflow-hidden border-b border-gray-200 sm:rounded-lg">
                  <table class="min-w-full divide-y divide-gray-200">
                    <thead class="bg-gray-50">
                      <tr>
                        <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                          Name
                        </th>
                        <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                          Type
                        </th>
                        <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                          Status
                        </th>
                        <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                          Created
                        </th>
                        <th scope="col" class="relative px-6 py-3">
                          <span class="sr-only">Actions</span>
                        </th>
                      </tr>
                    </thead>
                    <tbody class="bg-white divide-y divide-gray-200">
                      <tr v-for="ci in ciStore.paginatedCIs" :key="ci.id">
                        <td class="px-6 py-4 whitespace-nowrap">
                          <div class="text-sm font-medium text-gray-900">{{ ci.name }}</div>
                          <div v-if="ci.description" class="text-sm text-gray-500 truncate max-w-xs">{{ ci.description }}</div>
                        </td>
                        <td class="px-6 py-4 whitespace-nowrap">
                          <div class="text-sm text-gray-900">{{ ci.type }}</div>
                        </td>
                        <td class="px-6 py-4 whitespace-nowrap">
                          <span
                            class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                            :class="getStatusClass(ci.status)"
                          >
                            {{ ci.status }}
                          </span>
                        </td>
                        <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                          {{ formatDate(ci.created_at) }}
                        </td>
                        <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                          <router-link
                            :to="`/cis/${ci.id}`"
                            class="text-primary-600 hover:text-primary-900 mr-3"
                          >
                            View
                          </router-link>
                          <router-link
                            v-if="authStore.isAdmin"
                            :to="`/cis/${ci.id}/edit`"
                            class="text-gray-600 hover:text-gray-900"
                          >
                            Edit
                          </router-link>
                        </td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>
            </div>
            
            <!-- Pagination -->
            <div v-if="ciStore.pagination.totalPages > 1" class="mt-6 flex items-center justify-between">
              <div class="flex-1 flex justify-between sm:hidden">
                <button
                  @click="prevPage"
                  :disabled="ciStore.pagination.page === 1"
                  class="btn btn-secondary"
                  :class="{ 'opacity-50 cursor-not-allowed': ciStore.pagination.page === 1 }"
                >
                  Previous
                </button>
                <button
                  @click="nextPage"
                  :disabled="ciStore.pagination.page === ciStore.pagination.totalPages"
                  class="btn btn-secondary ml-3"
                  :class="{ 'opacity-50 cursor-not-allowed': ciStore.pagination.page === ciStore.pagination.totalPages }"
                >
                  Next
                </button>
              </div>
              <div class="hidden sm:flex-1 sm:flex sm:items-center sm:justify-between">
                <div>
                  <p class="text-sm text-gray-700">
                    Showing
                    <span class="font-medium">
                      {{ (ciStore.pagination.page - 1) * ciStore.pagination.limit + 1 }}
                    </span>
                    to
                    <span class="font-medium">
                      {{ Math.min(ciStore.pagination.page * ciStore.pagination.limit, ciStore.pagination.total) }}
                    </span>
                    of
                    <span class="font-medium">{{ ciStore.pagination.total }}</span>
                    results
                  </p>
                </div>
                <div>
                  <nav class="relative z-0 inline-flex rounded-md shadow-sm -space-x-px" aria-label="Pagination">
                    <button
                      @click="prevPage"
                      :disabled="ciStore.pagination.page === 1"
                      class="btn btn-secondary rounded-l-md"
                      :class="{ 'opacity-50 cursor-not-allowed': ciStore.pagination.page === 1 }"
                    >
                      <span class="sr-only">Previous</span>
                      <svg class="h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                        <path fill-rule="evenodd" d="M12.707 5.293a1 1 0 010 1.414L9.414 10l3.293 3.293a1 1 0 01-1.414 1.414l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 0z" clip-rule="evenodd" />
                      </svg>
                    </button>
                    
                    <template v-for="page in visiblePages" :key="page">
                      <button
                        v-if="page !== '...'"
                        @click="goToPage(page)"
                        class="btn"
                        :class="[
                          page === ciStore.pagination.page
                            ? 'btn-primary z-10'
                            : 'btn-secondary'
                        ]"
                      >
                        {{ page }}
                      </button>
                      <span
                        v-else
                        class="btn btn-secondary"
                      >
                        ...
                      </span>
                    </template>
                    
                    <button
                      @click="nextPage"
                      :disabled="ciStore.pagination.page === ciStore.pagination.totalPages"
                      class="btn btn-secondary rounded-r-md"
                      :class="{ 'opacity-50 cursor-not-allowed': ciStore.pagination.page === ciStore.pagination.totalPages }"
                    >
                      <span class="sr-only">Next</span>
                      <svg class="h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                        <path fill-rule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clip-rule="evenodd" />
                      </svg>
                    </button>
                  </nav>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useCIStore } from '../stores/ci'
import { useAuthStore } from '../stores/auth'
import { useUIStore } from '../stores/ui'

export default {
  name: 'CIListView',
  setup() {
    const ciStore = useCIStore()
    const authStore = useAuthStore()
    const uiStore = useUIStore()
    
    const filters = reactive({
      search: '',
      type: '',
      status: ''
    })
    
    const loadCIs = async () => {
      try {
        uiStore.setLoading(true)
        await ciStore.fetchCIs()
      } catch (error) {
        uiStore.showError('Failed to load configuration items')
      } finally {
        uiStore.setLoading(false)
      }
    }
    
    const applyFilters = () => {
      ciStore.setFilters(filters)
      loadCIs()
    }
    
    const resetFilters = () => {
      filters.search = ''
      filters.type = ''
      filters.status = ''
      ciStore.resetFilters()
      loadCIs()
    }
    
    const prevPage = () => {
      if (ciStore.pagination.page > 1) {
        ciStore.setPage(ciStore.pagination.page - 1)
        loadCIs()
      }
    }
    
    const nextPage = () => {
      if (ciStore.pagination.page < ciStore.pagination.totalPages) {
        ciStore.setPage(ciStore.pagination.page + 1)
        loadCIs()
      }
    }
    
    const goToPage = (page) => {
      ciStore.setPage(page)
      loadCIs()
    }
    
    const visiblePages = computed(() => {
      const currentPage = ciStore.pagination.page
      const totalPages = ciStore.pagination.totalPages
      
      if (totalPages <= 7) {
        return Array.from({ length: totalPages }, (_, i) => i + 1)
      }
      
      if (currentPage <= 3) {
        return [1, 2, 3, 4, '...', totalPages]
      }
      
      if (currentPage >= totalPages - 2) {
        return [1, '...', totalPages - 3, totalPages - 2, totalPages - 1, totalPages]
      }
      
      return [1, '...', currentPage - 1, currentPage, currentPage + 1, '...', totalPages]
    })
    
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
    
    const formatDate = (dateString) => {
      if (!dateString) return 'N/A'
      const date = new Date(dateString)
      return date.toLocaleDateString()
    }
    
    onMounted(() => {
      loadCIs()
    })
    
    return {
      ciStore,
      authStore,
      uiStore,
      filters,
      loadCIs,
      applyFilters,
      resetFilters,
      prevPage,
      nextPage,
      goToPage,
      visiblePages,
      getStatusClass,
      formatDate
    }
  }
}
</script>