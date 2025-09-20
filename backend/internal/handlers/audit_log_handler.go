package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cmdb-lite/backend/internal/middleware"
	"github.com/cmdb-lite/backend/internal/models"
	"github.com/cmdb-lite/backend/internal/repositories"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// AuditLogHandler handles HTTP requests for audit logs
type AuditLogHandler struct {
	auditRepo repositories.AuditLogRepository
}

// NewAuditLogHandler creates a new AuditLogHandler
func NewAuditLogHandler(auditRepo repositories.AuditLogRepository) *AuditLogHandler {
	return &AuditLogHandler{auditRepo: auditRepo}
}

// GetAuditLog handles retrieving an audit log by ID
// @Summary Get an audit log by ID
// @Description Get an audit log by its ID
// @Tags audit-logs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Audit log ID"
// @Success 200 {object} models.AuditLog
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /audit-logs/{id} [get]
func (h *AuditLogHandler) GetAuditLog(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL parameters
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	// Get the audit log
	auditLog, err := h.auditRepo.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Audit log not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(auditLog)
}

// GetAllAuditLogs handles retrieving all audit logs with pagination
// @Summary Get all audit logs
// @Description Get all audit logs with pagination
// @Tags audit-logs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Param entity_type query string false "Filter by entity type"
// @Param entity_id query string false "Filter by entity ID"
// @Param changed_by query string false "Filter by user who made the change"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /audit-logs [get]
func (h *AuditLogHandler) GetAllAuditLogs(w http.ResponseWriter, r *http.Request) {
	// Get pagination parameters from query string
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	entityType := r.URL.Query().Get("entity_type")
	entityIDStr := r.URL.Query().Get("entity_id")
	changedBy := r.URL.Query().Get("changed_by")

	// Set default values
	page := 1
	limit := 10

	// Parse page parameter
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	// Parse limit parameter
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	// Get audit logs based on filters
	var auditLogs []*models.AuditLog
	var err error

	if entityType != "" {
		// Filter by entity type
		auditLogs, err = h.auditRepo.GetByEntityType(r.Context(), entityType)
	} else if entityIDStr != "" {
		// Filter by entity ID
		entityID, err := uuid.Parse(entityIDStr)
		if err != nil {
			http.Error(w, "Invalid entity ID format", http.StatusBadRequest)
			return
		}
		auditLogs, err = h.auditRepo.GetByEntityID(r.Context(), entityID)
	} else if changedBy != "" {
		// Filter by changed by
		auditLogs, err = h.auditRepo.GetByChangedBy(r.Context(), changedBy)
	} else {
		// Get all audit logs
		auditLogs, err = h.auditRepo.GetAll(r.Context())
	}

	if err != nil {
		http.Error(w, "Failed to get audit logs", http.StatusInternalServerError)
		return
	}

	// Calculate pagination
	total := len(auditLogs)
	start := (page - 1) * limit
	end := start + limit

	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	paginatedAuditLogs := auditLogs[start:end]

	// Create response
	response := map[string]interface{}{
		"data":       paginatedAuditLogs,
		"pagination": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetAuditLogsByEntityType handles retrieving audit logs by entity type
// @Summary Get audit logs by entity type
// @Description Get audit logs filtered by entity type
// @Tags audit-logs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param entity_type path string true "Entity type"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /audit-logs/entity-type/{entity_type} [get]
func (h *AuditLogHandler) GetAuditLogsByEntityType(w http.ResponseWriter, r *http.Request) {
	// Extract entity type from URL parameters
	vars := mux.Vars(r)
	entityType, ok := vars["entity_type"]
	if !ok {
		http.Error(w, "Entity type parameter is required", http.StatusBadRequest)
		return
	}

	// Get pagination parameters from query string
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	// Set default values
	page := 1
	limit := 10

	// Parse page parameter
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	// Parse limit parameter
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	// Get audit logs by entity type
	auditLogs, err := h.auditRepo.GetByEntityType(r.Context(), entityType)
	if err != nil {
		http.Error(w, "Failed to get audit logs", http.StatusInternalServerError)
		return
	}

	// Calculate pagination
	total := len(auditLogs)
	start := (page - 1) * limit
	end := start + limit

	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	paginatedAuditLogs := auditLogs[start:end]

	// Create response
	response := map[string]interface{}{
		"data":       paginatedAuditLogs,
		"pagination": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetAuditLogsByEntityID handles retrieving audit logs by entity ID
// @Summary Get audit logs by entity ID
// @Description Get audit logs filtered by entity ID
// @Tags audit-logs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param entity_id path string true "Entity ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /audit-logs/entity-id/{entity_id} [get]
func (h *AuditLogHandler) GetAuditLogsByEntityID(w http.ResponseWriter, r *http.Request) {
	// Extract entity ID from URL parameters
	vars := mux.Vars(r)
	entityIDStr, ok := vars["entity_id"]
	if !ok {
		http.Error(w, "Entity ID parameter is required", http.StatusBadRequest)
		return
	}

	entityID, err := uuid.Parse(entityIDStr)
	if err != nil {
		http.Error(w, "Invalid entity ID format", http.StatusBadRequest)
		return
	}

	// Get pagination parameters from query string
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	// Set default values
	page := 1
	limit := 10

	// Parse page parameter
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	// Parse limit parameter
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	// Get audit logs by entity ID
	auditLogs, err := h.auditRepo.GetByEntityID(r.Context(), entityID)
	if err != nil {
		http.Error(w, "Failed to get audit logs", http.StatusInternalServerError)
		return
	}

	// Calculate pagination
	total := len(auditLogs)
	start := (page - 1) * limit
	end := start + limit

	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	paginatedAuditLogs := auditLogs[start:end]

	// Create response
	response := map[string]interface{}{
		"data":       paginatedAuditLogs,
		"pagination": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetAuditLogsByChangedBy handles retrieving audit logs by the user who made the change
// @Summary Get audit logs by changed by
// @Description Get audit logs filtered by the user who made the change
// @Tags audit-logs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param changed_by path string true "Username"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /audit-logs/changed-by/{changed_by} [get]
func (h *AuditLogHandler) GetAuditLogsByChangedBy(w http.ResponseWriter, r *http.Request) {
	// Extract changed by from URL parameters
	vars := mux.Vars(r)
	changedBy, ok := vars["changed_by"]
	if !ok {
		http.Error(w, "Changed by parameter is required", http.StatusBadRequest)
		return
	}

	// Get pagination parameters from query string
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	// Set default values
	page := 1
	limit := 10

	// Parse page parameter
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	// Parse limit parameter
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	// Get audit logs by changed by
	auditLogs, err := h.auditRepo.GetByChangedBy(r.Context(), changedBy)
	if err != nil {
		http.Error(w, "Failed to get audit logs", http.StatusInternalServerError)
		return
	}

	// Calculate pagination
	total := len(auditLogs)
	start := (page - 1) * limit
	end := start + limit

	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	paginatedAuditLogs := auditLogs[start:end]

	// Create response
	response := map[string]interface{}{
		"data":       paginatedAuditLogs,
		"pagination": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DeleteAuditLog handles deleting an audit log
// @Summary Delete an audit log
// @Description Delete an audit log
// @Tags audit-logs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Audit log ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /audit-logs/{id} [delete]
func (h *AuditLogHandler) DeleteAuditLog(w http.ResponseWriter, r *http.Request) {
	// Only admin users can delete audit logs
	userRole, ok := middleware.GetUserRoleFromContext(r.Context())
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	if userRole != "admin" {
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	}

	// Extract ID from URL parameters
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	// Delete the audit log
	if err := h.auditRepo.Delete(r.Context(), id); err != nil {
		http.Error(w, "Failed to delete audit log", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Audit log deleted successfully"})
}