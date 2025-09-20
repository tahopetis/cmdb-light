/**
 * Transform CI items and relationships into D3.js graph format
 * @param {Array} cis - Array of CI items
 * @param {Array} relationships - Array of relationships between CIs
 * @param {String} rootCIId - Optional root CI ID to focus the graph on
 * @returns {Object} - Object with nodes and links arrays for D3.js
 */
export function transformToGraphData(cis, relationships, rootCIId = null) {
  // Create a map of CI IDs to CI objects for quick lookup
  const ciMap = new Map()
  cis.forEach(ci => {
    ciMap.set(ci.id, ci)
  })
  
  // If a root CI is specified, filter to only include related CIs
  let includedCIIds = new Set()
  
  if (rootCIId) {
    // Start with the root CI
    includedCIIds.add(rootCIId)
    
    // Find all relationships that involve the root CI
    const rootRelationships = relationships.filter(rel => 
      rel.source_id === rootCIId || rel.target_id === rootCIId
    )
    
    // Add all directly related CIs
    rootRelationships.forEach(rel => {
      includedCIIds.add(rel.source_id)
      includedCIIds.add(rel.target_id)
    })
    
    // For each added CI, find their relationships (one level deep)
    const secondLevelCIIds = new Set()
    includedCIIds.forEach(ciId => {
      if (ciId === rootCIId) return
      
      const relatedRelationships = relationships.filter(rel => 
        rel.source_id === ciId || rel.target_id === ciId
      )
      
      relatedRelationships.forEach(rel => {
        secondLevelCIIds.add(rel.source_id)
        secondLevelCIIds.add(rel.target_id)
      })
    })
    
    // Add second level CIs
    secondLevelCIIds.forEach(ciId => {
      includedCIIds.add(ciId)
    })
  } else {
    // If no root CI is specified, include all CIs
    cis.forEach(ci => {
      includedCIIds.add(ci.id)
    })
  }
  
  // Create nodes array
  const nodes = Array.from(includedCIIds).map(ciId => {
    const ci = ciMap.get(ciId)
    return {
      id: ci.id,
      name: ci.name,
      type: ci.type,
      isRoot: ci.id === rootCIId
    }
  })
  
  // Create links array
  const links = relationships
    .filter(rel => 
      includedCIIds.has(rel.source_id) && 
      includedCIIds.has(rel.target_id)
    )
    .map(rel => ({
      id: rel.id,
      source: rel.source_id,
      target: rel.target_id,
      type: rel.type,
      description: rel.description
    }))
  
  return {
    nodes,
    links
  }
}

/**
 * Transform CI items and relationships into a hierarchical tree format
 * @param {Array} cis - Array of CI items
 * @param {Array} relationships - Array of relationships between CIs
 * @param {String} rootCIId - Root CI ID to build the tree from
 * @returns {Object} - Hierarchical tree structure
 */
export function transformToTreeData(cis, relationships, rootCIId) {
  // Create a map of CI IDs to CI objects for quick lookup
  const ciMap = new Map()
  cis.forEach(ci => {
    ciMap.set(ci.id, ci)
  })
  
  // Create a map of relationships by source ID
  const relationshipMap = new Map()
  relationships.forEach(rel => {
    if (!relationshipMap.has(rel.source_id)) {
      relationshipMap.set(rel.source_id, [])
    }
    relationshipMap.get(rel.source_id).push(rel)
  })
  
  // Recursive function to build the tree
  function buildTree(ciId, depth = 0) {
    const ci = ciMap.get(ciId)
    if (!ci) return null
    
    const node = {
      id: ci.id,
      name: ci.name,
      type: ci.type,
      depth,
      children: []
    }
    
    // Get all relationships where this CI is the source
    const childRelationships = relationshipMap.get(ciId) || []
    
    // Add children for each relationship
    childRelationships.forEach(rel => {
      const childNode = buildTree(rel.target_id, depth + 1)
      if (childNode) {
        childNode.relationshipType = rel.type
        childNode.relationshipId = rel.id
        node.children.push(childNode)
      }
    })
    
    return node
  }
  
  // Build the tree starting from the root CI
  return buildTree(rootCIId)
}

/**
 * Get all CI types present in the data
 * @param {Array} cis - Array of CI items
 * @returns {Array} - Array of unique CI types
 */
export function getCITypes(cis) {
  const types = new Set()
  cis.forEach(ci => {
    types.add(ci.type)
  })
  return Array.from(types)
}

/**
 * Get all relationship types present in the data
 * @param {Array} relationships - Array of relationships
 * @returns {Array} - Array of unique relationship types
 */
export function getRelationshipTypes(relationships) {
  const types = new Set()
  relationships.forEach(rel => {
    types.add(rel.type)
  })
  return Array.from(types)
}

/**
 * Group CIs by type
 * @param {Array} cis - Array of CI items
 * @returns {Object} - Object with CI types as keys and arrays of CIs as values
 */
export function groupCIsByType(cis) {
  const groups = {}
  
  cis.forEach(ci => {
    if (!groups[ci.type]) {
      groups[ci.type] = []
    }
    groups[ci.type].push(ci)
  })
  
  return groups
}

/**
 * Group relationships by type
 * @param {Array} relationships - Array of relationships
 * @returns {Object} - Object with relationship types as keys and arrays of relationships as values
 */
export function groupRelationshipsByType(relationships) {
  const groups = {}
  
  relationships.forEach(rel => {
    if (!groups[rel.type]) {
      groups[rel.type] = []
    }
    groups[rel.type].push(rel)
  })
  
  return groups
}

/**
 * Find the shortest path between two CIs
 * @param {Array} relationships - Array of relationships
 * @param {String} sourceId - Source CI ID
 * @param {String} targetId - Target CI ID
 * @returns {Array} - Array of CI IDs representing the path
 */
export function findShortestPath(relationships, sourceId, targetId) {
  // Create adjacency list
  const graph = new Map()
  
  // Initialize graph with all CI IDs
  relationships.forEach(rel => {
    if (!graph.has(rel.source_id)) {
      graph.set(rel.source_id, [])
    }
    if (!graph.has(rel.target_id)) {
      graph.set(rel.target_id, [])
    }
    
    // Add bidirectional edges (undirected graph)
    graph.get(rel.source_id).push(rel.target_id)
    graph.get(rel.target_id).push(rel.source_id)
  })
  
  // BFS to find shortest path
  const queue = [[sourceId]]
  const visited = new Set()
  visited.add(sourceId)
  
  while (queue.length > 0) {
    const path = queue.shift()
    const currentId = path[path.length - 1]
    
    if (currentId === targetId) {
      return path
    }
    
    const neighbors = graph.get(currentId) || []
    for (const neighborId of neighbors) {
      if (!visited.has(neighborId)) {
        visited.add(neighborId)
        queue.push([...path, neighborId])
      }
    }
  }
  
  // No path found
  return []
}

/**
 * Find all CIs that are connected to a given CI (directly or indirectly)
 * @param {Array} relationships - Array of relationships
 * @param {String} ciId - CI ID to find connections for
 * @returns {Set} - Set of CI IDs that are connected
 */
export function findConnectedCIs(relationships, ciId) {
  // Create adjacency list
  const graph = new Map()
  
  // Initialize graph with all CI IDs
  relationships.forEach(rel => {
    if (!graph.has(rel.source_id)) {
      graph.set(rel.source_id, [])
    }
    if (!graph.has(rel.target_id)) {
      graph.set(rel.target_id, [])
    }
    
    // Add bidirectional edges (undirected graph)
    graph.get(rel.source_id).push(rel.target_id)
    graph.get(rel.target_id).push(rel.source_id)
  })
  
  // BFS to find all connected nodes
  const queue = [ciId]
  const visited = new Set()
  visited.add(ciId)
  
  while (queue.length > 0) {
    const currentId = queue.shift()
    
    const neighbors = graph.get(currentId) || []
    for (const neighborId of neighbors) {
      if (!visited.has(neighborId)) {
        visited.add(neighborId)
        queue.push(neighborId)
      }
    }
  }
  
  return visited
}