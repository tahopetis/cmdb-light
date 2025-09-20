<template>
  <header class="bg-white shadow">
    <div class="px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between h-16">
        <div class="flex">
          <!-- Mobile menu button -->
          <button
            @click="toggleSidebar"
            type="button"
            class="-ml-0.5 -mt-0.5 h-12 w-12 inline-flex items-center justify-center rounded-md text-gray-500 hover:text-gray-900 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-primary-500 md:hidden"
          >
            <span class="sr-only">Open sidebar</span>
            <svg class="h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
            </svg>
          </button>
          
          <div class="flex-shrink-0 flex items-center">
            <img class="h-8 w-auto" src="/logo.svg" alt="CMDB Lite" />
            <span class="ml-2 text-xl font-bold text-gray-900">CMDB Lite</span>
          </div>
          
          <nav class="ml-6 flex space-x-8" aria-label="Global">
            <router-link
              v-for="item in navigation"
              :key="item.name"
              :to="item.href"
              class="inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium"
              :class="[
                route.path === item.href
                  ? 'border-primary-500 text-gray-900'
                  : 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700'
              ]"
            >
              {{ item.name }}
            </router-link>
          </nav>
        </div>
        
        <div class="flex items-center">
          <!-- Search bar -->
          <div class="flex-1 max-w-lg mx-4">
            <div class="relative">
              <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <svg class="h-5 w-5 text-gray-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                  <path fill-rule="evenodd" d="M8 4a4 4 0 100 8 4 4 0 000-8zM2 8a6 6 0 1110.89 3.476l4.817 4.817a1 1 0 01-1.414 1.414l-4.816-4.816A6 6 0 012 8z" clip-rule="evenodd" />
                </svg>
              </div>
              <input
                v-model="searchQuery"
                @keyup.enter="handleSearch"
                type="search"
                name="search"
                id="search"
                class="block w-full pl-10 pr-3 py-2 border border-gray-300 rounded-md leading-5 bg-white placeholder-gray-500 focus:outline-none focus:placeholder-gray-400 focus:ring-1 focus:ring-primary-500 focus:border-primary-500 sm:text-sm"
                placeholder="Search"
              />
            </div>
          </div>
          
          <!-- User dropdown -->
          <div class="ml-3 relative">
            <div>
              <button
                @click="toggleUserMenu"
                class="flex text-sm rounded-full focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
                id="user-menu-button"
                aria-expanded="false"
                aria-haspopup="true"
              >
                <span class="sr-only">Open user menu</span>
                <div class="h-8 w-8 rounded-full bg-primary-100 flex items-center justify-center">
                  <span class="text-primary-800 font-medium">
                    {{ userInitial }}
                  </span>
                </div>
              </button>
            </div>
            
            <div
              v-if="userMenuOpen"
              class="origin-top-right absolute right-0 mt-2 w-48 rounded-md shadow-lg py-1 bg-white ring-1 ring-black ring-opacity-5 focus:outline-none z-10"
              role="menu"
              aria-orientation="vertical"
              aria-labelledby="user-menu-button"
            >
              <div class="px-4 py-2 border-b">
                <p class="text-sm font-medium text-gray-900">{{ authStore.username }}</p>
                <p class="text-xs text-gray-500">{{ authStore.user?.role || 'User' }}</p>
              </div>
              
              <router-link
                to="/settings"
                class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                role="menuitem"
              >
                Settings
              </router-link>
              
              <button
                @click="handleLogout"
                class="block w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                role="menuitem"
              >
                Sign out
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </header>
</template>

<script>
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '../../stores/auth'
import { useUIStore } from '../../stores/ui'
import { useCIStore } from '../../stores/ci'

export default {
  name: 'Header',
  setup() {
    const route = useRoute()
    const router = useRouter()
    const authStore = useAuthStore()
    const uiStore = useUIStore()
    const ciStore = useCIStore()
    
    const userMenuOpen = ref(false)
    const searchQuery = ref('')
    
    const navigation = computed(() => {
      const items = [
        { name: 'Dashboard', href: '/' },
        { name: 'Configuration Items', href: '/cis' }
      ]
      
      // Add admin-only navigation items
      if (authStore.isAdmin) {
        items.push({ name: 'Audit Logs', href: '/audit-logs' })
      }
      
      return items
    })
    
    const userInitial = computed(() => {
      const username = authStore.username
      return username ? username.charAt(0).toUpperCase() : 'U'
    })
    
    const toggleSidebar = () => {
      uiStore.toggleSidebar()
    }
    
    const toggleUserMenu = () => {
      userMenuOpen.value = !userMenuOpen.value
    }
    
    const handleSearch = () => {
      if (!searchQuery.value.trim()) return
      
      // If not on CI list page, navigate to it
      if (route.path !== '/cis') {
        router.push('/cis')
      }
      
      // Set search filter
      ciStore.setFilters({ search: searchQuery.value })
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
      navigation,
      userInitial,
      userMenuOpen,
      searchQuery,
      toggleSidebar,
      toggleUserMenu,
      handleSearch,
      handleLogout
    }
  }
}
</script>