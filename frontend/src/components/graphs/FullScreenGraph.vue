<template>
  <div class="relative">
    <!-- Full Screen Toggle Button -->
    <button
      @click="toggleFullScreen"
      class="absolute top-4 right-4 z-10 bg-white rounded-md shadow-md p-2 hover:bg-gray-100"
      :title="isFullScreen ? 'Exit Full Screen' : 'Full Screen'"
    >
      <svg v-if="!isFullScreen" xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
        <path fill-rule="evenodd" d="M3 4a1 1 0 011-1h4a1 1 0 010 2H6.414l2.293 2.293a1 1 0 11-1.414 1.414L5 6.414V8a1 1 0 01-2 0V4zm9 1a1 1 0 110-2h4a1 1 0 011 1v4a1 1 0 11-2 0V6.414l-2.293 2.293a1 1 0 11-1.414-1.414L13.586 5H12zm-9 9a1 1 0 012 0v1.586l2.293-2.293a1 1 0 111.414 1.414L6.414 15H8a1 1 0 110 2H4a1 1 0 01-1-1v-4zm13-1a1 1 0 011 1v4a1 1 0 01-1 1h-4a1 1 0 110-2h1.586l-2.293-2.293a1 1 0 111.414-1.414L15 13.586V12a1 1 0 011-1z" clip-rule="evenodd" />
      </svg>
      <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
        <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
      </svg>
    </button>
    
    <!-- Graph Container -->
    <div
      ref="graphContainer"
      :class="[
        'relative bg-gray-50 transition-all duration-300',
        isFullScreen ? 'fixed inset-0 z-50' : 'h-96 w-full'
      ]"
    >
      <!-- Graph Controls Header -->
      <div
        v-if="isFullScreen"
        class="absolute top-0 left-0 right-0 bg-white bg-opacity-90 shadow-md p-4 z-10"
      >
        <div class="flex justify-between items-center">
          <h2 class="text-xl font-bold text-gray-800">Configuration Item Relationship Graph</h2>
          <div class="flex space-x-4">
            <!-- Search Input -->
            <div class="relative">
              <input
                v-model="searchQuery"
                type="text"
                placeholder="Search CIs..."
                class="pl-10 pr-4 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500"
                @keyup.enter="searchNodes"
              />
              <svg
                xmlns="http://www.w3.org/2000/svg"
                class="h-5 w-5 absolute left-3 top-2.5 text-gray-400"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
                />
              </svg>
            </div>
            
            <!-- Export Button -->
            <button
              @click="exportGraph"
              class="bg-primary-600 text-white px-4 py-2 rounded-md hover:bg-primary-700 flex items-center"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M3 17a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm3.293-7.707a1 1 0 011.414 0L9 10.586V3a1 1 0 112 0v7.586l1.293-1.293a1 1 0 111.414 1.414l-3 3a1 1 0 01-1.414 0l-3-3a1 1 0 010-1.414z" clip-rule="evenodd" />
              </svg>
              Export
            </button>
          </div>
        </div>
      </div>
      
      <!-- Legend -->
      <div
        v-if="isFullScreen && (ciTypes.length > 0 || relationshipTypes.length > 0)"
        class="absolute bottom-4 left-4 bg-white bg-opacity-90 rounded-md shadow-md p-4 z-10 max-w-xs"
      >
        <h3 class="font-semibold text-gray-700 mb-2">Legend</h3>
        
        <!-- CI Types -->
        <div v-if="ciTypes.length > 0" class="mb-3">
          <h4 class="text-sm font-medium text-gray-600 mb-1">CI Types</h4>
          <div class="space-y-1">
            <div v-for="type in ciTypes" :key="type" class="flex items-center">
              <div
                class="w-4 h-4 rounded-full mr-2"
                :style="{ backgroundColor: getColor(type) }"
              ></div>
              <span class="text-sm text-gray-600">{{ type }}</span>
            </div>
          </div>
        </div>
        
        <!-- Relationship Types -->
        <div v-if="relationshipTypes.length > 0">
          <h4 class="text-sm font-medium text-gray-600 mb-1">Relationship Types</h4>
          <div class="space-y-1">
            <div v-for="type in relationshipTypes" :key="type" class="flex items-center">
              <div class="w-6 h-0.5 bg-gray-400 mr-2"></div>
              <span class="text-sm text-gray-600">{{ type }}</span>
            </div>
          </div>
        </div>
      </div>
      
      <!-- Force Directed Graph Component -->
      <ForceDirectedGraph
        :nodes="graphData.nodes"
        :links="graphData.links"
        :width="graphWidth"
        :height="graphHeight"
        :is-loading="isLoading"
        @node-click="handleNodeClick"
        @node-double-click="handleNodeDoubleClick"
        @link-click="handleLinkClick"
      />
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useStore } from 'vuex'
import { transformToGraphData, getCITypes, getRelationshipTypes } from '@/utils/graphTransformations'
import ForceDirectedGraph from './ForceDirectedGraph.vue'

export default {
  name: 'FullScreenGraph',
  components: {
    ForceDirectedGraph
  },
  props: {
    cis: {
      type: Array,
      required: true
    },
    relationships: {
      type: Array,
      required: true
    },
    rootCIId: {
      type: String,
      default: null
    },
    isLoading: {
      type: Boolean,
      default: false
    }
  },
  emits: ['node-click', 'node-double-click', 'link-click'],
  setup(props, { emit }) {
    const store = useStore()
    const graphContainer = ref(null)
    const isFullScreen = ref(false)
    const searchQuery = ref('')
    
    // Colors for different CI types
    const typeColors = {
      'Server': '#3B82F6',
      'Application': '#10B981',
      'Database': '#8B5CF6',
      'Network Device': '#F59E0B',
      'Storage': '#6366F1',
      'Service': '#EC4899',
      'License': '#EF4444',
      'Other': '#6B7280'
    }
    
    // Get color for a CI type
    const getColor = (type) => {
      return typeColors[type] || typeColors['Other']
    }
    
    // Transform CI and relationship data to graph format
    const graphData = computed(() => {
      return transformToGraphData(props.cis, props.relationships, props.rootCIId)
    })
    
    // Get CI types for legend
    const ciTypes = computed(() => {
      return getCITypes(props.cis)
    })
    
    // Get relationship types for legend
    const relationshipTypes = computed(() => {
      return getRelationshipTypes(props.relationships)
    })
    
    // Calculate graph dimensions
    const graphWidth = computed(() => {
      if (isFullScreen.value) {
        return window.innerWidth
      }
      return graphContainer.value?.clientWidth || 800
    })
    
    const graphHeight = computed(() => {
      if (isFullScreen.value) {
        return window.innerHeight
      }
      return graphContainer.value?.clientHeight || 600
    })
    
    // Toggle full screen
    const toggleFullScreen = () => {
      isFullScreen.value = !isFullScreen.value
      
      if (isFullScreen.value) {
        document.body.style.overflow = 'hidden'
      } else {
        document.body.style.overflow = ''
      }
    }
    
    // Handle window resize
    const handleResize = () => {
      // Force re-render on resize
      if (graphContainer.value) {
        graphContainer.value.clientWidth // Trigger reflow
      }
    }
    
    // Handle node click
    const handleNodeClick = (node) => {
      emit('node-click', node)
    }
    
    // Handle node double click
    const handleNodeDoubleClick = (node) => {
      emit('node-double-click', node)
    }
    
    // Handle link click
    const handleLinkClick = (link) => {
      emit('link-click', link)
    }
    
    // Search nodes
    const searchNodes = () => {
      if (!searchQuery.value.trim()) return
      
      // Find matching CI
      const matchingCI = props.cis.find(ci => 
        ci.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
        ci.type.toLowerCase().includes(searchQuery.value.toLowerCase())
      )
      
      if (matchingCI) {
        // Set as root CI and emit node click
        emit('node-click', matchingCI)
      }
    }
    
    // Export graph as image
    const exportGraph = () => {
      // This would be implemented with a library like html2canvas or dom-to-image
      // For now, we'll just show a placeholder message
      store.dispatch('notification/show', {
        message: 'Export functionality would be implemented here',
        type: 'info'
      })
    }
    
    // Set up resize listener
    onMounted(() => {
      window.addEventListener('resize', handleResize)
    })
    
    // Clean up
    onUnmounted(() => {
      window.removeEventListener('resize', handleResize)
      document.body.style.overflow = ''
    })
    
    return {
      graphContainer,
      isFullScreen,
      searchQuery,
      graphData,
      ciTypes,
      relationshipTypes,
      graphWidth,
      graphHeight,
      getColor,
      toggleFullScreen,
      handleNodeClick,
      handleNodeDoubleClick,
      handleLinkClick,
      searchNodes,
      exportGraph
    }
  }
}
</script>

<style scoped>
/* Component-specific styles */
</style>