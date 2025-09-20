<template>
  <div class="relative w-full h-64 bg-gray-50 rounded-md overflow-hidden">
    <div ref="graphContainer" class="w-full h-full"></div>
    
    <!-- Loading Indicator -->
    <div v-if="isLoading" class="absolute inset-0 flex items-center justify-center bg-white bg-opacity-70">
      <div class="text-center">
        <svg class="animate-spin h-8 w-8 text-primary-600 mx-auto" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
        <p class="mt-1 text-sm text-gray-600">Loading relationships...</p>
      </div>
    </div>
    
    <!-- No Relationships Indicator -->
    <div v-if="!isLoading && (!nodes || nodes.length === 0)" class="absolute inset-0 flex items-center justify-center">
      <p class="text-gray-500 text-sm">No relationships found</p>
    </div>
    
    <!-- View Full Graph Button -->
    <button
      v-if="!isLoading && nodes && nodes.length > 0"
      @click="$emit('view-full-graph')"
      class="absolute bottom-2 right-2 bg-white bg-opacity-80 hover:bg-opacity-100 rounded-md shadow-sm px-3 py-1 text-sm text-primary-600 hover:text-primary-700 transition-all"
    >
      View Full Graph
    </button>
  </div>
</template>

<script>
import { ref, onMounted, onUnmounted, watch, computed } from 'vue'
import * as d3 from 'd3'

export default {
  name: 'MiniGraph',
  props: {
    nodes: {
      type: Array,
      required: true
    },
    links: {
      type: Array,
      required: true
    },
    isLoading: {
      type: Boolean,
      default: false
    }
  },
  emits: ['node-click', 'view-full-graph'],
  setup(props, { emit }) {
    const graphContainer = ref(null)
    let svg, simulation, link, node, label
    
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
    
    // Initialize the graph
    const initGraph = () => {
      if (!graphContainer.value || props.nodes.length === 0) return
      
      // Clear any existing graph
      d3.select(graphContainer.value).selectAll('*').remove()
      
      const width = graphContainer.value.clientWidth
      const height = graphContainer.value.clientHeight
      
      // Create SVG element
      svg = d3.select(graphContainer.value)
        .append('svg')
        .attr('width', width)
        .attr('height', height)
      
      // Create container group
      const g = svg.append('g')
      
      // Create arrow markers for directed links
      svg.append('defs').selectAll('marker')
        .data(['end'])
        .enter().append('marker')
        .attr('id', 'arrow')
        .attr('viewBox', '0 -5 10 10')
        .attr('refX', 20)
        .attr('refY', 0)
        .attr('markerWidth', 5)
        .attr('markerHeight', 5)
        .attr('orient', 'auto')
        .append('path')
        .attr('d', 'M0,-5L10,0L0,5')
        .attr('fill', '#999')
      
      // Create force simulation
      simulation = d3.forceSimulation(props.nodes)
        .force('link', d3.forceLink(props.links).id(d => d.id).distance(80))
        .force('charge', d3.forceManyBody().strength(-200))
        .force('center', d3.forceCenter(width / 2, height / 2))
        .force('collision', d3.forceCollide().radius(25))
      
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
      
      // Create nodes
      node = g.append('g')
        .attr('class', 'nodes')
        .selectAll('circle')
        .data(props.nodes)
        .enter().append('circle')
        .attr('r', d => d.isRoot ? 15 : 12)
        .attr('fill', d => getColor(d.type))
        .attr('stroke', '#fff')
        .attr('stroke-width', d => d.isRoot ? 3 : 2)
        .on('click', (event, d) => {
          emit('node-click', d)
        })
      
      // Add node icons
      const nodeIcons = g.append('g')
        .attr('class', 'node-icons')
        .selectAll('text')
        .data(props.nodes)
        .enter().append('text')
        .attr('font-family', 'FontAwesome')
        .attr('font-size', d => d.isRoot ? '10px' : '8px')
        .attr('fill', '#fff')
        .attr('text-anchor', 'middle')
        .attr('dominant-baseline', 'central')
        .text(d => getNodeIcon(d.type))
      
      // Create labels (only for root node)
      label = g.append('g')
        .attr('class', 'labels')
        .selectAll('text')
        .data(props.nodes.filter(d => d.isRoot))
        .enter().append('text')
        .attr('text-anchor', 'middle')
        .attr('dy', 20)
        .text(d => d.name)
        .attr('font-size', '10px')
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
    
    // Update graph when data changes
    const updateGraph = () => {
      if (!svg || props.nodes.length === 0) {
        initGraph()
        return
      }
      
      // Update simulation with new data
      simulation.nodes(props.nodes)
      simulation.force('link').links(props.links)
      simulation.alpha(1).restart()
      
      // Update links
      link = svg.select('.links').selectAll('line')
        .data(props.links)
      
      link.exit().remove()
      
      link = link.enter().append('line')
        .attr('stroke', '#999')
        .attr('stroke-opacity', 0.6)
        .attr('stroke-width', 1.5)
        .attr('marker-end', 'url(#arrow)')
        .merge(link)
      
      // Update nodes
      node = svg.select('.nodes').selectAll('circle')
        .data(props.nodes)
      
      node.exit().remove()
      
      node = node.enter().append('circle')
        .attr('stroke', '#fff')
        .on('click', (event, d) => {
          emit('node-click', d)
        })
        .merge(node)
        .attr('r', d => d.isRoot ? 15 : 12)
        .attr('fill', d => getColor(d.type))
        .attr('stroke-width', d => d.isRoot ? 3 : 2)
      
      // Update node icons
      const nodeIcons = svg.select('.node-icons').selectAll('text')
        .data(props.nodes)
      
      nodeIcons.exit().remove()
      
      nodeIcons = nodeIcons.enter().append('text')
        .attr('font-family', 'FontAwesome')
        .attr('fill', '#fff')
        .attr('text-anchor', 'middle')
        .attr('dominant-baseline', 'central')
        .merge(nodeIcons)
        .attr('font-size', d => d.isRoot ? '10px' : '8px')
        .text(d => getNodeIcon(d.type))
      
      // Update labels (only for root node)
      label = svg.select('.labels').selectAll('text')
        .data(props.nodes.filter(d => d.isRoot))
      
      label.exit().remove()
      
      label = label.enter().append('text')
        .attr('text-anchor', 'middle')
        .attr('dy', 20)
        .merge(label)
        .text(d => d.name)
        .attr('font-size', '10px')
        .attr('fill', '#333')
      
      // Add title tooltips
      node.append('title')
        .text(d => `${d.name} (${d.type})`)
    }
    
    // Watch for changes in props
    watch([() => props.nodes, () => props.links], () => {
      if (props.nodes.length > 0) {
        updateGraph()
      }
    }, { deep: true })
    
    // Initialize graph on mount
    onMounted(() => {
      if (props.nodes.length > 0) {
        initGraph()
      }
    })
    
    // Clean up on unmount
    onUnmounted(() => {
      if (simulation) {
        simulation.stop()
      }
    })
    
    return {
      graphContainer
    }
  }
}
</script>

<style scoped>
/* Component-specific styles */
</style>