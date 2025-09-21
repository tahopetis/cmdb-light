<template>
  <div class="py-6">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between items-center">
        <h1 class="text-2xl font-semibold text-gray-900">CI Relationship Graph</h1>
        <div class="flex space-x-3">
          <button
            @click="refreshGraph"
            class="btn btn-secondary"
            :disabled="loading"
          >
            Refresh
          </button>
        </div>
      </div>
      
      <div class="mt-6">
        <div class="card">
          <div class="px-4 py-5 sm:p-6">
            <div v-if="loading" class="flex justify-center items-center h-96">
              <div class="text-center">
                <svg class="animate-spin h-12 w-12 text-primary-600 mx-auto" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                <p class="mt-2 text-sm text-gray-600">Loading graph data...</p>
              </div>
            </div>
            
            <div v-else-if="error" class="text-center py-12">
              <svg class="mx-auto h-12 w-12 text-gray-400" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
              </svg>
              <h3 class="mt-2 text-sm font-medium text-gray-900">Error loading graph</h3>
              <p class="mt-1 text-sm text-gray-500">{{ error }}</p>
              <div class="mt-6">
                <button
                  @click="refreshGraph"
                  class="btn btn-primary"
                >
                  Try again
                </button>
              </div>
            </div>
            
            <div v-else class="relative">
              <div class="bg-gray-50 rounded-lg border border-gray-200 p-4">
                <div class="flex justify-between items-center mb-4">
                  <h3 class="text-lg font-medium text-gray-900">Configuration Item Relationships</h3>
                  <div class="flex space-x-2">
                    <button
                      @click="zoomIn"
                      class="p-1 rounded-md hover:bg-gray-200"
                      title="Zoom in"
                    >
                      <svg class="h-5 w-5 text-gray-600" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
                      </svg>
                    </button>
                    <button
                      @click="zoomOut"
                      class="p-1 rounded-md hover:bg-gray-200"
                      title="Zoom out"
                    >
                      <svg class="h-5 w-5 text-gray-600" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18 12H6" />
                      </svg>
                    </button>
                    <button
                      @click="resetZoom"
                      class="p-1 rounded-md hover:bg-gray-200"
                      title="Reset zoom"
                    >
                      <svg class="h-5 w-5 text-gray-600" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                      </svg>
                    </button>
                  </div>
                </div>
                
                <div class="bg-white rounded-lg border border-gray-300 h-96">
                  <ForceDirectedGraph
                    :nodes="graphData.nodes"
                    :links="graphData.links"
                    :width="graphWidth"
                    :height="graphHeight"
                    :isLoading="loading"
                    :highlightedCIId="highlightedCIId"
                    :maxNodes="100"
                    :enablePerformanceOptimizations="true"
                    :enableLazyLoading="true"
                    @node-click="handleNodeClick"
                    @node-double-click="handleNodeDoubleClick"
                    @link-click="handleLinkClick"
                    @load-relationships="loadRelationships"
                  />
                </div>
                
                <div class="mt-4 grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
                  <div class="bg-white rounded-lg border border-gray-200 p-4">
                    <h4 class="text-sm font-medium text-gray-900">Total Nodes</h4>
                    <p class="mt-1 text-2xl font-semibold text-primary-600">{{ graphData.nodes.length || 0 }}</p>
                  </div>
                  <div class="bg-white rounded-lg border border-gray-200 p-4">
                    <h4 class="text-sm font-medium text-gray-900">Total Edges</h4>
                    <p class="mt-1 text-2xl font-semibold text-primary-600">{{ graphData.links.length || 0 }}</p>
                  </div>
                  <div class="bg-white rounded-lg border border-gray-200 p-4">
                    <h4 class="text-sm font-medium text-gray-900">CI Types</h4>
                    <p class="mt-1 text-2xl font-semibold text-primary-600">{{ graphData.types || 0 }}</p>
                  </div>
                  <div class="bg-white rounded-lg border border-gray-200 p-4">
                    <h4 class="text-sm font-medium text-gray-900">Relationship Types</h4>
                    <p class="mt-1 text-2xl font-semibold text-primary-600">{{ graphData.relationshipTypes || 0 }}</p>
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
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useCIStore } from '../stores/ci'
import { useUIStore } from '../stores/ui'
import ForceDirectedGraph from '../components/graphs/ForceDirectedGraph.vue'
import ciService from '../services/ciService'
import relationshipService from '../services/relationshipService'

export default {
  name: 'GraphView',
  components: {
    ForceDirectedGraph
  },
  setup() {
    const router = useRouter()
    const ciStore = useCIStore()
    const uiStore = useUIStore()
    
    const loading = ref(false)
    const error = ref(null)
    const highlightedCIId = ref(null)
    
    const graphData = reactive({
      nodes: [],
      links: [],
      types: 0,
      relationshipTypes: 0
    })
    
    // Graph dimensions
    const graphWidth = ref(800)
    const graphHeight = ref(384) // h-96 = 24rem = 384px
    
    // Fetch graph data from API
    const fetchGraphData = async () => {
      loading.value = true
      error.value = null
      
      try {
        // Get all CIs
        const cisResponse = await ciService.getAllCIs()
        const cis = cisResponse.data.data || []
        
        // Get all relationships
        const relationshipsResponse = await relationshipService.getAllRelationships()
        const relationships = relationshipsResponse.data.data || []
        
        // Transform CIs to nodes
        const nodes = cis.map(ci => ({
          id: ci.id,
          name: ci.name,
          type: ci.type
        }))
        
        // Transform relationships to links
        const links = relationships.map(rel => ({
          source: rel.source_id,
          target: rel.target_id,
          type: rel.type,
          id: rel.id
        }))
        
        // Count unique CI types
        const uniqueTypes = new Set(cis.map(ci => ci.type))
        
        // Count unique relationship types
        const uniqueRelationshipTypes = new Set(relationships.map(rel => rel.type))
        
        // Update graph data
        graphData.nodes = nodes
        graphData.links = links
        graphData.types = uniqueTypes.size
        graphData.relationshipTypes = uniqueRelationshipTypes.size
        
      } catch (err) {
        console.error('Error fetching graph data:', err)
        error.value = err.response?.data?.message || 'Failed to fetch graph data'
      } finally {
        loading.value = false
      }
    }
    
    // Refresh graph data
    const refreshGraph = async () => {
      await fetchGraphData()
    }
    
    // Handle node click
    const handleNodeClick = (node) => {
      highlightedCIId.value = highlightedCIId.value === node.id ? null : node.id
      uiStore.showInfo(`Selected CI: ${node.name} (${node.type})`)
    }
    
    // Handle node double click
    const handleNodeDoubleClick = (node) => {
      router.push({
        name: 'ci-detail',
        params: { id: node.id }
      })
    }
    
    // Handle link click
    const handleLinkClick = (link) => {
      uiStore.showInfo(`Relationship: ${link.type}`)
    }
    
    // Load relationships for a node (lazy loading)
    const loadRelationships = async (nodeId) => {
      try {
        // In a real implementation, this would fetch additional relationships
        // For now, we'll just show a notification
        uiStore.showInfo(`Loading relationships for CI ${nodeId}`)
      } catch (err) {
        console.error('Error loading relationships:', err)
        uiStore.showError('Failed to load relationships')
      }
    }
    
    // Zoom controls
    const zoomIn = () => {
      // This will be handled by the ForceDirectedGraph component
    }
    
    const zoomOut = () => {
      // This will be handled by the ForceDirectedGraph component
    }
    
    const resetZoom = () => {
      // This will be handled by the ForceDirectedGraph component
    }
    
    onMounted(() => {
      fetchGraphData()
    })
    
    return {
      loading,
      error,
      graphData,
      graphWidth,
      graphHeight,
      highlightedCIId,
      refreshGraph,
      zoomIn,
      zoomOut,
      resetZoom,
      handleNodeClick,
      handleNodeDoubleClick,
      handleLinkClick,
      loadRelationships
    }
  }
}
</script>
