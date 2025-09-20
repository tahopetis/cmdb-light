<template>
  <div class="card hover:shadow-md transition-shadow duration-200">
    <div class="px-4 py-5 sm:p-6">
      <div class="flex items-start justify-between">
        <div class="flex items-center">
          <div
            class="flex-shrink-0 h-10 w-10 rounded-md flex items-center justify-center"
            :class="getTypeColorClass(ci.type)"
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
          <div class="ml-4">
            <h3 class="text-lg font-medium text-gray-900">
              <router-link
                :to="`/cis/${ci.id}`"
                class="hover:text-primary-600 transition-colors duration-200"
              >
                {{ ci.name }}
              </router-link>
            </h3>
            <p class="text-sm text-gray-500">{{ ci.type }}</p>
          </div>
        </div>
        
        <div class="flex items-center">
          <span
            class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
            :class="getStatusClass(ci.status)"
          >
            {{ ci.status }}
          </span>
        </div>
      </div>
      
      <div class="mt-4">
        <p class="text-sm text-gray-600 line-clamp-2">
          {{ ci.description || 'No description available' }}
        </p>
      </div>
      
      <div v-if="ci.tags && ci.tags.length > 0" class="mt-4">
        <div class="flex flex-wrap gap-1">
          <span
            v-for="(tag, index) in ci.tags"
            :key="index"
            class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-gray-100 text-gray-800"
          >
            {{ tag }}
          </span>
        </div>
      </div>
      
      <div class="mt-6 flex justify-between items-center">
        <div class="text-sm text-gray-500">
          Last updated: {{ formatDate(ci.updated_at) }}
        </div>
        
        <div class="flex space-x-2">
          <router-link
            :to="`/cis/${ci.id}`"
            class="text-primary-600 hover:text-primary-900 text-sm font-medium"
          >
            View
          </router-link>
          
          <router-link
            v-if="authStore.isAdmin"
            :to="`/cis/${ci.id}/edit`"
            class="text-gray-600 hover:text-gray-900 text-sm font-medium"
          >
            Edit
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { computed } from 'vue'
import { useAuthStore } from '../../stores/auth'

export default {
  name: 'CICard',
  props: {
    ci: {
      type: Object,
      required: true
    }
  },
  setup(props) {
    const authStore = useAuthStore()
    
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
    
    return {
      authStore,
      getTypeColorClass,
      getStatusClass,
      formatDate
    }
  }
}
</script>