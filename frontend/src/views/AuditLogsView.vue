<template>
  <div class="py-6">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between items-center">
        <h1 class="text-2xl font-semibold text-gray-900">Audit Logs</h1>
        <div class="flex space-x-3">
          <button
            @click="refreshLogs"
            class="btn btn-secondary"
            :disabled="loading"
          >
            Refresh
          </button>
          <button
            @click="exportLogs"
            class="btn btn-primary"
            :disabled="loading"
          >
            Export
          </button>
        </div>
      </div>
      
      <div class="mt-6">
        <div class="card">
          <div class="px-4 py-5 sm:p-6">
            <!-- Filters -->
            <div class="mb-6 grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
              <div>
                <label for="action-filter" class="block text-sm font-medium text-gray-700">Action</label>
                <select
                  id="action-filter"
                  v-model="filters.action"
                  class="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-primary-500 focus:border-primary-500 sm:text-sm rounded-md"
                >
                  <option value="">All Actions</option>
                  <option value="CREATE">Create</option>
                  <option value="UPDATE">Update</option>
                  <option value="DELETE">Delete</option>
                  <option value="LOGIN">Login</option>
                  <option value="LOGOUT">Logout</option>
                </select>
              </div>
              
              <div>
                <label for="entity-type-filter" class="block text-sm font-medium text-gray-700">Entity Type</label>
                <select
                  id="entity-type-filter"
                  v-model="filters.entityType"
                  class="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-primary-500 focus:border-primary-500 sm:text-sm rounded-md"
                >
                  <option value="">All Types</option>
                  <option value="CI">Configuration Item</option>
                  <option value="USER">User</option>
                  <option value="RELATIONSHIP">Relationship</option>
                </select>
              </div>
              
              <div>
                <label for="user-filter" class="block text-sm font-medium text-gray-700">User</label>
                <input
                  type="text"
                  id="user-filter"
                  v-model="filters.user"
                  placeholder="Filter by username"
                  class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-primary-500 focus:border-primary-500 sm:text-sm"
                />
              </div>
              
              <div>
                <label for="date-filter" class="block text-sm font-medium text-gray-700">Date Range</label>
                <select
                  id="date-filter"
                  v-model="filters.dateRange"
                  class="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-primary-500 focus:border-primary-500 sm:text-sm rounded-md"
                >
                  <option value="today">Today</option>
                  <option value="week">Last 7 Days</option>
                  <option value="month">Last 30 Days</option>
                  <option value="all">All Time</option>
                </select>
              </div>
            </div>
            
            <!-- Logs table -->
            <div v-if="loading" class="flex justify-center items-center h-96">
              <div class="text-center">
                <svg class="animate-spin h-12 w-12 text-primary-600 mx-auto" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                <p class="mt-2 text-sm text-gray-600">Loading audit logs...</p>
              </div>
            </div>
            
            <div v-else-if="error" class="text-center py-12">
              <svg class="mx-auto h-12 w-12 text-gray-400" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
              </svg>
              <h3 class="mt-2 text-sm font-medium text-gray-900">Error loading audit logs</h3>
              <p class="mt-1 text-sm text-gray-500">{{ error }}</p>
              <div class="mt-6">
                <button
                  @click="refreshLogs"
                  class="btn btn-primary"
                >
                  Try again
                </button>
              </div>
            </div>
            
            <div v-else>
              <div class="overflow-hidden border border-gray-200 rounded-lg">
                <table class="min-w-full divide-y divide-gray-200">
                  <thead class="bg-gray-50">
                    <tr>
                      <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Timestamp
                      </th>
                      <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        User
                      </th>
                      <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Action
                      </th>
                      <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Entity Type
                      </th>
                      <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Entity ID
                      </th>
                      <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Details
                      </th>
                    </tr>
                  </thead>
                  <tbody class="bg-white divide-y divide-gray-200">
                    <tr v-for="log in filteredLogs" :key="log.id">
                      <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {{ formatDate(log.timestamp) }}
                      </td>
                      <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                        {{ log.user }}
                      </td>
                      <td class="px-6 py-4 whitespace-nowrap">
                        <span
                          class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full"
                          :class="getActionClass(log.action)"
                        >
                          {{ log.action }}
                        </span>
                      </td>
                      <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {{ log.entityType }}
                      </td>
                      <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {{ log.entityId }}
                      </td>
                      <td class="px-6 py-4 text-sm text-gray-500 max-w-xs truncate">
                        {{ log.details }}
                      </td>
                    </tr>
                    
                    <tr v-if="filteredLogs.length === 0">
                      <td colspan="6" class="px-6 py-4 text-center text-sm text-gray-500">
                        No audit logs found matching your filters
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>
              
              <!-- Pagination -->
              <div class="bg-white px-4 py-3 flex items-center justify-between border-t border-gray-200 sm:px-6 mt-4">
                <div class="flex-1 flex justify-between sm:hidden">
                  <button
                    @click="prevPage"
                    :disabled="pagination.page === 1"
                    class="relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50"
                    :class="{ 'opacity-50 cursor-not-allowed': pagination.page === 1 }"
                  >
                    Previous
                  </button>
                  <button
                    @click="nextPage"
                    :disabled="pagination.page >= pagination.totalPages"
                    class="ml-3 relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50"
                    :class="{ 'opacity-50 cursor-not-allowed': pagination.page >= pagination.totalPages }"
                  >
                    Next
                  </button>
                </div>
                <div class="hidden sm:flex-1 sm:flex sm:items-center sm:justify-between">
                  <div>
                    <p class="text-sm text-gray-700">
                      Showing
                      <span class="font-medium">{{ (pagination.page - 1) * pagination.limit + 1 }}</span>
                      to
                      <span class="font-medium">
                        {{ Math.min(pagination.page * pagination.limit, pagination.total) }}
                      </span>
                      of
                      <span class="font-medium">{{ pagination.total }}</span>
                      results
                    </p>
                  </div>
                  <div>
                    <nav class="relative z-0 inline-flex rounded-md shadow-sm -space-x-px" aria-label="Pagination">
                      <button
                        @click="prevPage"
                        :disabled="pagination.page === 1"
                        class="relative inline-flex items-center px-2 py-2 rounded-l-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50"
                        :class="{ 'opacity-50 cursor-not-allowed': pagination.page === 1 }"
                      >
                        <span class="sr-only">Previous</span>
                        <svg class="h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                          <path fill-rule="evenodd" d="M12.707 5.293a1 1 0 010 1.414L9.414 10l3.293 3.293a1 1 0 01-1.414 1.414l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 0z" clip-rule="evenodd" />
                        </svg>
                      </button>
                      
                      <button
                        v-for="page in visiblePages"
                        :key="page"
                        @click="goToPage(page)"
                        class="relative inline-flex items-center px-4 py-2 border text-sm font-medium"
                        :class="[
                          page === pagination.page
                            ? 'z-10 bg-primary-50 border-primary-500 text-primary-600'
                            : 'bg-white border-gray-300 text-gray-500 hover:bg-gray-50'
                        ]"
                      >
                        {{ page }}
                      </button>
                      
                      <button
                        @click="nextPage"
                        :disabled="pagination.page >= pagination.totalPages"
                        class="relative inline-flex items-center px-2 py-2 rounded-r-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50"
                        :class="{ 'opacity-50 cursor-not-allowed': pagination.page >= pagination.totalPages }"
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
  </div>
</template>

<script>
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useUIStore } from '../stores/ui'

export default {
  name: 'AuditLogsView',
  setup() {
    const uiStore = useUIStore()
    
    const loading = ref(false)
    const error = ref(null)
    const logs = ref([])
    
    const filters = reactive({
      action: '',
      entityType: '',
      user: '',
      dateRange: 'week'
    })
    
    const pagination = reactive({
      page: 1,
      limit: 10,
      total: 0,
      totalPages: 0
    })
    
    const visiblePages = computed(() => {
      const current = pagination.page
      const total = pagination.totalPages
      const delta = 2
      
      const range = []
      const rangeWithDots = []
      let l
      
      range.push(1)
      
      for (let i = current - delta; i <= current + delta; i++) {
        if (i > 1 && i < total) {
          range.push(i)
        }
      }
      
      if (total > 1) {
        range.push(total)
      }
      
      range.forEach(i => {
        if (l) {
          if (i - l === 2) {
            rangeWithDots.push(l + 1)
          } else if (i - l !== 1) {
            rangeWithDots.push('...')
          }
        }
        rangeWithDots.push(i)
        l = i
      })
      
      return rangeWithDots
    })
    
    const filteredLogs = computed(() => {
      let result = [...logs.value]
      
      // Apply filters
      if (filters.action) {
        result = result.filter(log => log.action === filters.action)
      }
      
      if (filters.entityType) {
        result = result.filter(log => log.entityType === filters.entityType)
      }
      
      if (filters.user) {
        const searchLower = filters.user.toLowerCase()
        result = result.filter(log => log.user.toLowerCase().includes(searchLower))
      }
      
      // Apply date range filter
      if (filters.dateRange !== 'all') {
        const now = new Date()
        let startDate
        
        switch (filters.dateRange) {
          case 'today':
            startDate = new Date(now.getFullYear(), now.getMonth(), now.getDate())
            break
          case 'week':
            startDate = new Date(now)
            startDate.setDate(now.getDate() - 7)
            break
          case 'month':
            startDate = new Date(now)
            startDate.setDate(now.getDate() - 30)
            break
        }
        
        if (startDate) {
          result = result.filter(log => new Date(log.timestamp) >= startDate)
        }
      }
      
      return result
    })
    
    const fetchLogs = async () => {
      loading.value = true
      error.value = null
      
      try {
        // In a real implementation, this would fetch audit logs from the API
        // For now, we'll simulate it with mock data
        await new Promise(resolve => setTimeout(resolve, 1000))
        
        // Mock data
        logs.value = [
          {
            id: 1,
            timestamp: new Date(Date.now() - 3600000).toISOString(),
            user: 'admin',
            action: 'CREATE',
            entityType: 'CI',
            entityId: 'ci-123',
            details: 'Created new configuration item: Web Server'
          },
          {
            id: 2,
            timestamp: new Date(Date.now() - 7200000).toISOString(),
            user: 'admin',
            action: 'UPDATE',
            entityType: 'CI',
            entityId: 'ci-456',
            details: 'Updated configuration item: Database Server'
          },
          {
            id: 3,
            timestamp: new Date(Date.now() - 86400000).toISOString(),
            user: 'viewer',
            action: 'LOGIN',
            entityType: 'USER',
            entityId: 'user-789',
            details: 'User logged in'
          },
          {
            id: 4,
            timestamp: new Date(Date.now() - 172800000).toISOString(),
            user: 'admin',
            action: 'DELETE',
            entityType: 'CI',
            entityId: 'ci-101',
            details: 'Deleted configuration item: Old Server'
          },
          {
            id: 5,
            timestamp: new Date(Date.now() - 259200000).toISOString(),
            user: 'admin',
            action: 'CREATE',
            entityType: 'RELATIONSHIP',
            entityId: 'rel-202',
            details: 'Created relationship between Web Server and Database Server'
          }
        ]
        
        pagination.total = logs.value.length
        pagination.totalPages = Math.ceil(pagination.total / pagination.limit)
      } catch (err) {
        console.error('Error fetching audit logs:', err)
        error.value = err.response?.data?.message || 'Failed to fetch audit logs'
      } finally {
        loading.value = false
      }
    }
    
    const refreshLogs = () => {
      fetchLogs()
    }
    
    const exportLogs = () => {
      uiStore.showInfo('Export functionality would be implemented here')
    }
    
    const prevPage = () => {
      if (pagination.page > 1) {
        pagination.page--
      }
    }
    
    const nextPage = () => {
      if (pagination.page < pagination.totalPages) {
        pagination.page++
      }
    }
    
    const goToPage = (page) => {
      if (typeof page === 'number') {
        pagination.page = page
      }
    }
    
    const getActionClass = (action) => {
      const classes = {
        'CREATE': 'bg-green-100 text-green-800',
        'UPDATE': 'bg-blue-100 text-blue-800',
        'DELETE': 'bg-red-100 text-red-800',
        'LOGIN': 'bg-purple-100 text-purple-800',
        'LOGOUT': 'bg-yellow-100 text-yellow-800'
      }
      return classes[action] || 'bg-gray-100 text-gray-800'
    }
    
    const formatDate = (dateString) => {
      if (!dateString) return 'N/A'
      const date = new Date(dateString)
      return date.toLocaleString()
    }
    
    // Watch for filter changes
    watch(filters, () => {
      pagination.page = 1
    }, { deep: true })
    
    onMounted(() => {
      fetchLogs()
    })
    
    return {
      loading,
      error,
      logs,
      filters,
      pagination,
      visiblePages,
      filteredLogs,
      refreshLogs,
      exportLogs,
      prevPage,
      nextPage,
      goToPage,
      getActionClass,
      formatDate
    }
  }
}
</script>