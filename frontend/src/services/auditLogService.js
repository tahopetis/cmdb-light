import apiClient from './api'

class AuditLogService {
  // Get all audit logs with pagination
  async getAllAuditLogs(params = {}) {
    try {
      const response = await apiClient.get('/audit-logs', { params })
      return response.data
    } catch (error) {
      console.error('Error fetching audit logs:', error)
      throw error
    }
  }

  // Get a specific audit log by ID
  async getAuditLog(id) {
    try {
      const response = await apiClient.get(`/audit-logs/${id}`)
      return response.data
    } catch (error) {
      console.error(`Error fetching audit log with ID ${id}:`, error)
      throw error
    }
  }

  // Get audit logs by entity type
  async getAuditLogsByEntityType(entityType, params = {}) {
    try {
      const response = await apiClient.get(`/audit-logs/entity-type/${entityType}`, { params })
      return response.data
    } catch (error) {
      console.error(`Error fetching audit logs for entity type ${entityType}:`, error)
      throw error
    }
  }

  // Get audit logs by entity ID
  async getAuditLogsByEntityID(entityId, params = {}) {
    try {
      const response = await apiClient.get(`/audit-logs/entity-id/${entityId}`, { params })
      return response.data
    } catch (error) {
      console.error(`Error fetching audit logs for entity ID ${entityId}:`, error)
      throw error
    }
  }

  // Get audit logs by the user who made changes
  async getAuditLogsByChangedBy(changedBy, params = {}) {
    try {
      const response = await apiClient.get(`/audit-logs/changed-by/${changedBy}`, { params })
      return response.data
    } catch (error) {
      console.error(`Error fetching audit logs for user ${changedBy}:`, error)
      throw error
    }
  }

  // Delete an audit log (admin only)
  async deleteAuditLog(id) {
    try {
      await apiClient.delete(`/audit-logs/${id}`)
      return { success: true }
    } catch (error) {
      console.error(`Error deleting audit log with ID ${id}:`, error)
      throw error
    }
  }
}

export default new AuditLogService()