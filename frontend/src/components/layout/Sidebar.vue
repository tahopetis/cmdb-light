<template>
  <div
    class="hidden md:flex md:flex-shrink-0"
    :class="{ 'fixed inset-0 z-40': sidebarOpen }"
  >
    <div class="flex flex-col w-64 border-r border-gray-200 bg-white">
      <div class="h-0 flex-1 flex flex-col pt-5 pb-4 overflow-y-auto">
        <!-- Logo and app name -->
        <div class="flex items-center flex-shrink-0 px-4">
          <img class="h-8 w-auto" src="/logo.svg" alt="CMDB Lite" />
          <span class="ml-2 text-xl font-bold text-gray-900">CMDB Lite</span>
        </div>
        
        <!-- Navigation -->
        <nav class="mt-8 flex-1 px-2 bg-white space-y-1">
          <router-link
            v-for="item in navigation"
            :key="item.name"
            :to="item.href"
            class="group flex items-center px-2 py-2 text-sm font-medium rounded-md"
            :class="[
              route.path === item.href
                ? 'bg-primary-100 text-primary-900'
                : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'
            ]"
          >
            <component
              :is="item.icon"
              class="mr-3 flex-shrink-0 h-6 w-6"
              :class="[
                route.path === item.href
                  ? 'text-primary-500'
                  : 'text-gray-400 group-hover:text-gray-500'
              ]"
              aria-hidden="true"
            />
            {{ item.name }}
          </router-link>
        </nav>
        
        <!-- CI Type Filters -->
        <div class="px-4 py-4 border-t border-gray-200">
          <h3 class="px-3 text-xs font-semibold text-gray-500 uppercase tracking-wider">
            Filter by Type
          </h3>
          <div class="mt-2 space-y-1">
            <button
              v-for="type in ciStore.ciTypes"
              :key="type"
              @click="filterByType(type)"
              class="group w-full flex items-center px-3 py-2 text-sm font-medium rounded-md"
              :class="[
                ciStore.filters.type === type
                  ? 'bg-primary-100 text-primary-900'
                  : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'
              ]"
            >
              <span
                class="inline-block h-2 w-2 rounded-full mr-3"
                :class="getTypeColor(type)"
              ></span>
              {{ type }}
              <span
                class="ml-auto inline-block py-0.5 px-2 text-xs rounded-full"
                :class="[
                  ciStore.filters.type === type
                    ? 'bg-primary-200 text-primary-800'
                    : 'bg-gray-100 text-gray-800'
                ]"
              >
                {{ ciStore.getCITypesCount[type] || 0 }}
              </span>
            </button>
            
            <button
              v-if="ciStore.filters.type"
              @click="clearTypeFilter"
              class="group w-full flex items-center px-3 py-2 text-sm font-medium rounded-md text-gray-600 hover:bg-gray-50 hover:text-gray-900"
            >
              <svg
                class="mr-3 flex-shrink-0 h-5 w-5 text-gray-400 group-hover:text-gray-500"
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 20 20"
                fill="currentColor"
                aria-hidden="true"
              >
                <path
                  fill-rule="evenodd"
                  d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
                  clip-rule="evenodd"
                />
              </svg>
              Clear Filter
            </button>
          </div>
        </div>
      </div>
      
      <!-- User profile -->
      <div class="flex-shrink-0 flex border-t border-gray-200 p-4">
        <div class="flex-shrink-0 w-full group block">
          <div class="flex items-center">
            <div class="ml-3 flex-1">
              <p class="text-sm font-medium text-gray-700 group-hover:text-gray-900">
                {{ authStore.username }}
              </p>
              <p class="text-xs text-gray-500 group-hover:text-gray-700">
                {{ authStore.user?.role || 'User' }}
              </p>
            </div>
            <button
              @click="handleLogout"
              class="ml-3 flex-shrink-0 text-gray-400 hover:text-gray-500"
            >
              <svg
                class="h-5 w-5"
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 20 20"
                fill="currentColor"
                aria-hidden="true"
              >
                <path
                  fill-rule="evenodd"
                  d="M3 3a1 1 0 00-1 1v12a1 1 0 102 0V4a1 1 0 00-1-1zm10.293 9.293a1 1 0 001.414 1.414l3-3a1 1 0 000-1.414l-3-3a1 1 0 10-1.414 1.414L14.586 9H7a1 1 0 100 2h7.586l-1.293 1.293z"
                  clip-rule="evenodd"
                />
              </svg>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { computed, h } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '../../stores/auth'
import { useUIStore } from '../../stores/ui'
import { useCIStore } from '../../stores/ci'

export default {
  name: 'Sidebar',
  props: {
    sidebarOpen: {
      type: Boolean,
      default: false
    }
  },
  setup(props) {
    const route = useRoute()
    const router = useRouter()
    const authStore = useAuthStore()
    const uiStore = useUIStore()
    const ciStore = useCIStore()
    
    const navigation = computed(() => {
      const items = [
        {
          name: 'Dashboard',
          href: '/',
          icon: {
            render() {
              return h('svg', { 
                xmlns: 'http://www.w3.org/2000/svg', 
                fill: 'none', 
                viewBox: '0 0 24 24', 
                stroke: 'currentColor' 
              }, [
                h('path', { 
                  'stroke-linecap': 'round', 
                  'stroke-linejoin': 'round', 
                  'stroke-width': '2', 
                  d: 'M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6' 
                })
              ])
            }
          }
        },
        {
          name: 'Configuration Items',
          href: '/cis',
          icon: {
            render() {
              return h('svg', { 
                xmlns: 'http://www.w3.org/2000/svg', 
                fill: 'none', 
                viewBox: '0 0 24 24', 
                stroke: 'currentColor' 
              }, [
                h('path', { 
                  'stroke-linecap': 'round', 
                  'stroke-linejoin': 'round', 
                  'stroke-width': '2', 
                  d: 'M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4' 
                })
              ])
            }
          }
        }
      ]
      
      // Add admin-only navigation items
      if (authStore.isAdmin) {
        items.push({
          name: 'Audit Logs',
          href: '/audit-logs',
          icon: {
            render() {
              return h('svg', { 
                xmlns: 'http://www.w3.org/2000/svg', 
                fill: 'none', 
                viewBox: '0 0 24 24', 
                stroke: 'currentColor' 
              }, [
                h('path', { 
                  'stroke-linecap': 'round', 
                  'stroke-linejoin': 'round', 
                  'stroke-width': '2', 
                  d: 'M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z' 
                })
              ])
            }
          }
        })
      }
      
      return items
    })
    
    const getTypeColor = (type) => {
      const colors = {
        'Server': 'bg-blue-500',
        'Application': 'bg-green-500',
        'Database': 'bg-purple-500',
        'Network Device': 'bg-yellow-500',
        'Storage': 'bg-indigo-500',
        'Service': 'bg-pink-500',
        'License': 'bg-red-500',
        'Other': 'bg-gray-500'
      }
      return colors[type] || colors['Other']
    }
    
    const filterByType = (type) => {
      ciStore.setFilters({ type })
      
      // If not on CI list page, navigate to it
      if (route.path !== '/cis') {
        router.push('/cis')
      }
    }
    
    const clearTypeFilter = () => {
      ciStore.setFilters({ type: '' })
    }
    
    const handleLogout = () => {
      uiStore.showConfirmDialog(
        'Confirm Logout',
        'Are you sure you want to logout?',
        () => {
          authStore.logout()
          router.push('/login')
          uiStore.showSuccess('You have been logged out')
        }
      )
    }
    
    return {
      route,
      authStore,
      ciStore,
      navigation,
      getTypeColor,
      filterByType,
      clearTypeFilter,
      handleLogout
    }
  }
}
</script>
