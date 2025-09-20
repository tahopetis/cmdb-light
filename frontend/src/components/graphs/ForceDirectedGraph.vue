<template>
  <div class="relative w-full h-full">
    <div ref="graphContainer" class="w-full h-full"></div>
    
    <!-- Graph Controls -->
    <div class="absolute top-4 right-4 bg-white rounded-md shadow-md p-2 flex flex-col space-y-2">
      <button
        @click="zoomIn"
        class="p-2 rounded-md hover:bg-gray-100"
        title="Zoom In"
      >
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
          <path fill-rule="evenodd" d="M10 5a1 1 0 011 1v3h3a1 1 0 110 2h-3v3a1 1 0 11-2 0v-3H6a1 1 0 110-2h3V6a1 1 0 011-1z" clip-rule="evenodd" />
        </svg>
      </button>
      <button
        @click="zoomOut"
        class="p-2 rounded-md hover:bg-gray-100"
        title="Zoom Out"
      >
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
          <path fill-rule="evenodd" d="M5 10a1 1 0 011-1h8a1 1 0 110 2H6a1 1 0 01-1-1z" clip-rule="evenodd" />
        </svg>
      </button>
      <button
        @click="resetZoom"
        class="p-2 rounded-md hover:bg-gray-100"
        title="Reset Zoom"
      >
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
          <path fill-rule="evenodd" d="M4 2a1 1 0 011 1v2.101a7.002 7.002 0 0111.601 2.566 1 1 0 11-1.885.666A5.002 5.002 0 005.999 7H9a1 1 0 010 2H4a1 1 0 01-1-1V3a1 1 0 011-1zm.008 9.057a1 1 0 011.276.61A5.002 5.002 0 0014.001 13H11a1 1 0 110-2h5a1 1 0 011 1v5a1 1 0 11-2 0v-2.101a7.002 7.002 0 01-11.601-2.566 1 1 0 01.61-1.276z" clip-rule="evenodd" />
        </svg>
      </button>
      <button
        @click="expandAll"
        class="p-2 rounded-md hover:bg-gray-100"
        title="Expand All"
      >
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
          <path fill-rule="evenodd" d="M12 13a1 1 0 100 2h5a1 1 0 001-1V9a1 1 0 10-2 0v2.586l-4.293-4.293a1 1 0 00-1.414 0L8 9.586 3.707 5.293a1 1 0 00-1.414 1.414l5 5a1 1 0 001.414 0L11 9.414 14.586 13H12z" clip-rule="evenodd" />
        </svg>
      </button>
      <button
        @click="collapseAll"
        class="p-2 rounded-md hover:bg-gray-100"
        title="Collapse All"
      >
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
          <path fill-rule="evenodd" d="M14.707 12.707a1 1 0 01-1.414 0L10 9.414l-3.293 3.293a1 1 0 01-1.414-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 010 1.414z" clip-rule="evenodd" />
        </svg>
      </button>
    </div>
    
    <!-- Loading Indicator -->
    <div v-if="isLoading" class="absolute inset-0 flex items-center justify-center bg-white bg-opacity-70">
      <div class="text-center">
        <svg class="animate-spin h-12 w-12 text-primary-600 mx-auto" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
        <p class="mt-2 text-gray-600">Loading graph...</p>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, onUnmounted, watch, computed } from 'vue'
import { useRouter } from 'vue-router'
import { debounce } from 'lodash'
import * as d3 from 'd3'

export default {
  name: 'ForceDirectedGraph',
  props: {
    nodes: {
      type: Array,
      required: true
    },
    links: {
      type: Array,
      required: true
    },
    width: {
      type: Number,
      default: 800
    },
    height: {
      type: Number,
      default: 600
    },
    isLoading: {
      type: Boolean,
      default: false
    },
    highlightedCIId: {
      type: String,
      default: null
    },
    maxNodes: {
      type: Number,
      default: 100
    },
    enablePerformanceOptimizations: {
      type: Boolean,
      default: true
    },
    enableLazyLoading: {
      type: Boolean,
      default: true
    }
  },
  emits: ['node-click', 'node-double-click', 'link-click', 'load-relationships'],
  setup(props, { emit }) {
    const router = useRouter()
    const graphContainer = ref(null)
    let svg, simulation, link, node, label, zoom
    
    // Track expanded/collapsed state of nodes
    const expandedNodes = ref(new Set())
    
    // Track which relationships have been loaded
    const loadedRelationships = ref(new Set())
    
    // Track all nodes and links for expand/collapse functionality
    const allNodes = ref([...props.nodes])
    const allLinks = ref([...props.links])
    
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
    
    // Initialize the graph
    const initGraph = () => {
      if (!graphContainer.value) return
      
      // Clear any existing graph
      d3.select(graphContainer.value).selectAll('*').remove()
      
      // Create SVG element
      svg = d3.select(graphContainer.value)
        .append('svg')
        .attr('width', props.width)
        .attr('height', props.height)
      
      // Create zoom behavior
      zoom = d3.zoom()
        .scaleExtent([0.1, 10])
        .on('zoom', (event) => {
          g.attr('transform', event.transform)
        })
      
      svg.call(zoom)
      
      // Create container group for zoom
      const g = svg.append('g')
      
      // Create arrow markers for directed links
      svg.append('defs').selectAll('marker')
        .data(['end'])
        .enter().append('marker')
        .attr('id', 'arrow')
        .attr('viewBox', '0 -5 10 10')
        .attr('refX', 25)
        .attr('refY', 0)
        .attr('markerWidth', 6)
        .attr('markerHeight', 6)
        .attr('orient', 'auto')
        .append('path')
        .attr('d', 'M0,-5L10,0L0,5')
        .attr('fill', '#999')
      
      // Create force simulation
      simulation = d3.forceSimulation(props.nodes)
        .force('link', d3.forceLink(props.links).id(d => d.id).distance(100))
        .force('charge', d3.forceManyBody().strength(-300))
        .force('center', d3.forceCenter(props.width / 2, props.height / 2))
        .force('collision', d3.forceCollide().radius(30))
      
      // Create links
      link = g.append('g')
        .attr('class', 'links')
        .selectAll('line')
        .data(props.links)
        .enter().append('line')
        .attr('stroke', '#999')
        .attr('stroke-opacity', 0.6)
        .attr('stroke-width', 1.5)
        .attr('marker-end', 'url(#arrow)')
        .on('click', (event, d) => {
          emit('link-click', d)
        })
      
      // Create link labels
      const linkLabels = g.append('g')
        .attr('class', 'link-labels')
        .selectAll('text')
        .data(props.links)
        .enter().append('text')
        .attr('font-size', '10px')
        .attr('fill', '#666')
        .attr('text-anchor', 'middle')
        .attr('dy', -5)
        .text(d => d.type)
      
      // Create nodes
      node = g.append('g')
        .attr('class', 'nodes')
        .selectAll('circle')
        .data(props.nodes)
        .enter().append('circle')
        .attr('r', 20)
        .attr('fill', d => getColor(d.type))
        .attr('stroke', '#fff')
        .attr('stroke-width', 2)
        .call(drag(simulation))
        .on('click', (event, d) => {
          event.stopPropagation()
          toggleNodeExpansion(d.id)
          emit('node-click', d)
        })
        .on('dblclick', (event, d) => {
          event.stopPropagation()
          navigateToCI(d)
          emit('node-double-click', d)
        })
      
      // Add node icons
      const nodeIcons = g.append('g')
        .attr('class', 'node-icons')
        .selectAll('text')
        .data(props.nodes)
        .enter().append('text')
        .attr('font-family', 'FontAwesome')
        .attr('font-size', '14px')
        .attr('fill', '#fff')
        .attr('text-anchor', 'middle')
        .attr('dominant-baseline', 'central')
        .text(d => getNodeIcon(d.type))
      
      // Create labels
      label = g.append('g')
        .attr('class', 'labels')
        .selectAll('text')
        .data(props.nodes)
        .enter().append('text')
        .attr('text-anchor', 'middle')
        .attr('dy', 30)
        .text(d => d.name)
        .attr('font-size', '12px')
        .attr('fill', '#333')
      
      // Add title tooltips
      node.append('title')
        .text(d => `${d.name} (${d.type})`)
      
      // Update positions on tick
      simulation.on('tick', () => {
        link
          .attr('x1', d => d.source.x)
          .attr('y1', d => d.source.y)
          .attr('x2', d => d.target.x)
          .attr('y2', d => d.target.y)
        
        linkLabels
          .attr('x', d => (d.source.x + d.target.x) / 2)
          .attr('y', d => (d.source.y + d.target.y) / 2)
        
        node
          .attr('cx', d => d.x)
          .attr('cy', d => d.y)
        
        nodeIcons
          .attr('x', d => d.x)
          .attr('y', d => d.y)
        
        label
          .attr('x', d => d.x)
          .attr('y', d => d.y)
      })
    }
    
    // Get icon for a CI type
    const getNodeIcon = (type) => {
      const icons = {
        'Server': '\uf233',
        'Application': '\uf109',
        'Database': '\uf1c0',
        'Network Device': '\uf6ff',
        'Storage': '\uf1b0',
        'Service': '\uf013',
        'License': '\uf15c',
        'Other': '\uf1b2'
      }
      return icons[type] || icons['Other']
    }
    
    // Drag behavior for nodes
    const drag = (simulation) => {
      function dragstarted(event, d) {
        if (!event.active) simulation.alphaTarget(0.3).restart()
        d.fx = d.x
        d.fy = d.y
      }
      
      function dragged(event, d) {
        d.fx = event.x
        d.fy = event.y
      }
      
      function dragended(event, d) {
        if (!event.active) simulation.alphaTarget(0)
        d.fx = null
        d.fy = null
      }
      
      return d3.drag()
        .on('start', dragstarted)
        .on('drag', dragged)
        .on('end', dragended)
    }
    
    // Zoom in
    const zoomIn = () => {
      svg.transition().call(
        zoom.scaleBy,
        1.3,
        d3.zoomTransform(svg.node()).invert([props.width / 2, props.height / 2])
      )
    }
    
    // Zoom out
    const zoomOut = () => {
      svg.transition().call(
        zoom.scaleBy,
        0.7,
        d3.zoomTransform(svg.node()).invert([props.width / 2, props.height / 2])
      )
    }
    
    // Reset zoom
    const resetZoom = () => {
      svg.transition().call(
        zoom.transform,
        d3.zoomIdentity
      )
    }
    
    // Toggle node expansion
    const toggleNodeExpansion = (nodeId) => {
      if (expandedNodes.value.has(nodeId)) {
        expandedNodes.value.delete(nodeId)
      } else {
        expandedNodes.value.add(nodeId)
        
        // If lazy loading is enabled and relationships haven't been loaded yet
        if (props.enableLazyLoading && !loadedRelationships.value.has(nodeId)) {
          // Emit event to load relationships
          emit('load-relationships', nodeId)
          // Mark as loaded to prevent duplicate requests
          loadedRelationships.value.add(nodeId)
        }
      }
      updateVisibleGraph()
    }
    
    // Expand all nodes
    const expandAll = () => {
      props.nodes.forEach(node => {
        expandedNodes.value.add(node.id)
      })
      updateVisibleGraph()
    }
    
    // Collapse all nodes
    const collapseAll = () => {
      expandedNodes.value.clear()
      // Keep root nodes visible
      if (props.nodes.length > 0) {
        expandedNodes.value.add(props.nodes[0].id)
      }
      updateVisibleGraph()
    }
    
    // Create a debounced version of updateGraph for performance
    const debouncedUpdateGraph = debounce(updateGraph, 100)
    
    // Update visible graph based on expanded nodes
    const updateVisibleGraph = () => {
      let filteredNodes, filteredLinks
      
      if (expandedNodes.value.size === 0) {
        // If no nodes are expanded, show all nodes
        filteredNodes = [...props.nodes]
        filteredLinks = [...props.links]
      } else {
        // Start with expanded nodes
        const visibleNodeIds = new Set(expandedNodes.value)
        
        // Add nodes that are directly connected to expanded nodes
        const linksToAdd = []
        props.links.forEach(link => {
          if (visibleNodeIds.has(link.source.id || link.source) ||
              visibleNodeIds.has(link.target.id || link.target)) {
            linksToAdd.push(link)
            visibleNodeIds.add(link.source.id || link.source)
            visibleNodeIds.add(link.target.id || link.target)
          }
        })
        
        // Filter nodes and links
        filteredNodes = props.nodes.filter(node => visibleNodeIds.has(node.id))
        filteredLinks = linksToAdd
      }
      
      // Limit nodes and links for performance
      const limited = limitNodesAndLinks(filteredNodes, filteredLinks)
      allNodes.value = limited.nodes
      allLinks.value = limited.links
      
      // Update the graph with filtered data (use debounced version for performance)
      if (props.enablePerformanceOptimizations) {
        debouncedUpdateGraph()
      } else {
        updateGraph()
      }
    }
    
    // Navigate to CI detail view
    const navigateToCI = (node) => {
      router.push({
        name: 'ci-detail',
        params: { id: node.id }
      })
    }
    
    // Check if a node should be highlighted
    const isNodeHighlighted = (nodeId) => {
      if (!props.highlightedCIId) return false
      
      // The highlighted CI itself
      if (nodeId === props.highlightedCIId) return true
      
      // Find all directly related CIs
      const relatedNodeIds = new Set()
      props.links.forEach(link => {
        const sourceId = link.source.id || link.source
        const targetId = link.target.id || link.target
        
        if (sourceId === props.highlightedCIId) {
          relatedNodeIds.add(targetId)
        } else if (targetId === props.highlightedCIId) {
          relatedNodeIds.add(sourceId)
        }
      })
      
      return relatedNodeIds.has(nodeId)
    }
    
    // Check if a link should be highlighted
    const isLinkHighlighted = (link) => {
      if (!props.highlightedCIId) return false
      
      const sourceId = link.source.id || link.source
      const targetId = link.target.id || link.target
      
      return sourceId === props.highlightedCIId || targetId === props.highlightedCIId
    }
    
    // Limit the number of nodes and links for performance
    const limitNodesAndLinks = (nodes, links) => {
      if (!props.enablePerformanceOptimizations || nodes.length <= props.maxNodes) {
        return { nodes, links }
      }
      
      // If we have a highlighted CI, prioritize it and its related nodes
      if (props.highlightedCIId) {
        const highlightedNode = nodes.find(n => n.id === props.highlightedCIId)
        if (highlightedNode) {
          // Get all directly related nodes
          const relatedNodeIds = new Set([props.highlightedCIId])
          const relatedLinks = []
          
          links.forEach(link => {
            const sourceId = link.source.id || link.source
            const targetId = link.target.id || link.target
            
            if (sourceId === props.highlightedCIId || targetId === props.highlightedCIId) {
              relatedLinks.push(link)
              relatedNodeIds.add(sourceId)
              relatedNodeIds.add(targetId)
            }
          })
          
          // Get related nodes
          const relatedNodes = nodes.filter(n => relatedNodeIds.has(n.id))
          
          // If we still have too many nodes, add more nodes based on connectivity
          if (relatedNodes.length < props.maxNodes) {
            const remainingNodes = nodes.filter(n => !relatedNodeIds.has(n.id))
            
            // Sort by number of connections (degree)
            const nodeDegrees = new Map()
            links.forEach(link => {
              const sourceId = link.source.id || link.source
              const targetId = link.target.id || link.target
              
              nodeDegrees.set(sourceId, (nodeDegrees.get(sourceId) || 0) + 1)
              nodeDegrees.set(targetId, (nodeDegrees.get(targetId) || 0) + 1)
            })
            
            // Sort by degree (highest first)
            remainingNodes.sort((a, b) => {
              return (nodeDegrees.get(b.id) || 0) - (nodeDegrees.get(a.id) || 0)
            })
            
            // Add nodes until we reach the max
            const nodesToAdd = Math.min(props.maxNodes - relatedNodes.length, remainingNodes.length)
            for (let i = 0; i < nodesToAdd; i++) {
              relatedNodes.push(remainingNodes[i])
            }
            
            // Get all links between the selected nodes
            const selectedNodeIds = new Set(relatedNodes.map(n => n.id))
            const selectedLinks = links.filter(link => {
              const sourceId = link.source.id || link.source
              const targetId = link.target.id || link.target
              return selectedNodeIds.has(sourceId) && selectedNodeIds.has(targetId)
            })
            
            return { nodes: relatedNodes, links: selectedLinks }
          }
          
          return { nodes: relatedNodes, links: relatedLinks }
        }
      }
      
      // If no highlighted CI, just take the first maxNodes nodes
      const limitedNodes = nodes.slice(0, props.maxNodes)
      const limitedNodeIds = new Set(limitedNodes.map(n => n.id))
      
      // Get all links between the selected nodes
      const limitedLinks = links.filter(link => {
        const sourceId = link.source.id || link.source
        const targetId = link.target.id || link.target
        return limitedNodeIds.has(sourceId) && limitedNodeIds.has(targetId)
      })
      
      return { nodes: limitedNodes, links: limitedLinks }
    }
    
    // Update graph when data changes
    const updateGraph = () => {
      if (!svg) {
        initGraph()
        return
      }
      
      // Update simulation with new data
      simulation.nodes(allNodes.value)
      simulation.force('link').links(allLinks.value)
      simulation.alpha(1).restart()
      
      // Update links
      link = svg.select('.links').selectAll('line')
        .data(allLinks.value)
      
      link.exit().remove()
      
      link = link.enter().append('line')
        .attr('stroke', d => isLinkHighlighted(d) ? '#F59E0B' : '#999')
        .attr('stroke-opacity', d => isLinkHighlighted(d) ? 1 : 0.6)
        .attr('stroke-width', d => isLinkHighlighted(d) ? 2.5 : 1.5)
        .attr('marker-end', 'url(#arrow)')
        .on('click', (event, d) => {
          emit('link-click', d)
        })
        .merge(link)
      
      // Update link labels
      const linkLabels = svg.select('.link-labels').selectAll('text')
        .data(allLinks.value)
      
      linkLabels.exit().remove()
      
      linkLabels = linkLabels.enter().append('text')
        .attr('font-size', '10px')
        .attr('fill', '#666')
        .attr('text-anchor', 'middle')
        .attr('dy', -5)
        .merge(linkLabels)
        .text(d => d.type)
      
      // Update nodes
      node = svg.select('.nodes').selectAll('circle')
        .data(allNodes.value)
      
      node.exit().remove()
      
      node = node.enter().append('circle')
        .attr('r', 20)
        .attr('stroke', '#fff')
        .attr('stroke-width', 2)
        .call(drag(simulation))
        .on('click', (event, d) => {
          event.stopPropagation()
          toggleNodeExpansion(d.id)
          emit('node-click', d)
        })
        .on('dblclick', (event, d) => {
          event.stopPropagation()
          navigateToCI(d)
          emit('node-double-click', d)
        })
        .merge(node)
        .attr('fill', d => {
          if (isNodeHighlighted(d.id)) {
            return d.id === props.highlightedCIId ? '#F59E0B' : '#FCD34D'
          }
          return getColor(d.type)
        })
        .attr('stroke', d => isNodeHighlighted(d.id) ? '#D97706' : '#fff')
        .attr('stroke-width', d => {
          if (isNodeHighlighted(d.id)) {
            return expandedNodes.value.has(d.id) ? 4 : 3
          }
          return expandedNodes.value.has(d.id) ? 3 : 2
        })
      
      // Update node icons
      const nodeIcons = svg.select('.node-icons').selectAll('text')
        .data(allNodes.value)
      
      nodeIcons.exit().remove()
      
      nodeIcons = nodeIcons.enter().append('text')
        .attr('font-family', 'FontAwesome')
        .attr('font-size', '14px')
        .attr('fill', '#fff')
        .attr('text-anchor', 'middle')
        .attr('dominant-baseline', 'central')
        .merge(nodeIcons)
        .text(d => getNodeIcon(d.type))
      
      // Add expand/collapse indicator
      const expandIndicators = svg.select('.nodes').selectAll('text.expand-indicator')
        .data(allNodes.value)
      
      expandIndicators.exit().remove()
      
      expandIndicators = expandIndicators.enter().append('text')
        .attr('class', 'expand-indicator')
        .attr('font-family', 'FontAwesome')
        .attr('font-size', '10px')
        .attr('fill', '#fff')
        .attr('text-anchor', 'middle')
        .attr('dominant-baseline', 'central')
        .merge(expandIndicators)
        .attr('x', d => d.x + 12)
        .attr('y', d => d.y - 12)
        .text(d => {
          // Check if node has relationships
          const hasRelationships = props.links.some(link =>
            (link.source.id || link.source) === d.id ||
            (link.target.id || link.target) === d.id
          )
          return hasRelationships ? (expandedNodes.value.has(d.id) ? '\uf068' : '\uf067') : ''
        })
      
      // Update labels
      label = svg.select('.labels').selectAll('text')
        .data(allNodes.value)
      
      label.exit().remove()
      
      label = label.enter().append('text')
        .attr('text-anchor', 'middle')
        .attr('dy', 30)
        .merge(label)
        .text(d => d.name)
        .attr('font-size', '12px')
        .attr('fill', '#333')
      
      // Add title tooltips
      node.append('title')
        .text(d => `${d.name} (${d.type})`)
    }
    
    // Watch for changes in props
    watch([() => props.nodes, () => props.links], () => {
      // Update all nodes and links when props change
      allNodes.value = [...props.nodes]
      allLinks.value = [...props.links]
      
      // Reset loaded relationships when data changes
      loadedRelationships.value.clear()
      
      // Initialize expanded nodes with the first node
      if (props.nodes.length > 0 && expandedNodes.value.size === 0) {
        expandedNodes.value.add(props.nodes[0].id)
      }
      
      // Use debounced update if performance optimizations are enabled
      if (props.enablePerformanceOptimizations) {
        debouncedUpdateGraph()
      } else {
        updateGraph()
      }
    }, { deep: true, immediate: true })
    
    // Initialize graph on mount
    onMounted(() => {
      initGraph()
    })
    
    // Clean up on unmount
    onUnmounted(() => {
      if (simulation) {
        simulation.stop()
      }
    })
    
    return {
      graphContainer,
      zoomIn,
      zoomOut,
      resetZoom,
      expandAll,
      collapseAll,
      toggleNodeExpansion
    }
  }
}
</script>

<style scoped>
/* Component-specific styles */
</style>