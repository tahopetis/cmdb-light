package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/cmdb-lite/backend/internal/middleware"
	"github.com/cmdb-lite/backend/internal/models"
	"github.com/cmdb-lite/backend/internal/repositories"
	"github.com/cmdb-lite/backend/internal/validation"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// CIHandler handles HTTP requests for CIs
type CIHandler struct {
	ciRepo    repositories.CIRepository
	relRepo   repositories.RelationshipRepository
	auditRepo repositories.AuditLogRepository
	validator *validation.Validator
}

// NewCIHandler creates a new CIHandler
func NewCIHandler(
	ciRepo repositories.CIRepository,
	relRepo repositories.RelationshipRepository,
	auditRepo repositories.AuditLogRepository,
) *CIHandler {
	return &CIHandler{
		ciRepo:    ciRepo,
		relRepo:   relRepo,
		auditRepo: auditRepo,
		validator: validation.NewValidator(),
	}
}

// CreateCI handles the creation of a new CI
// @Summary Create a new CI
// @Description Create a new configuration item
// @Tags cis
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param ci body models.CI true "CI object"
// @Success 201 {object} models.CI
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /cis [post]
func (h *CIHandler) CreateCI(w http.ResponseWriter, r *http.Request) {
	// Get the username from the context
	username, ok := middleware.GetUsernameFromContext(r.Context())
	if !ok {
		middleware.RespondWithUnauthorizedError(w, "User not authenticated", nil)
		return
	}

	var ci models.CI
	if err := json.NewDecoder(r.Body).Decode(&ci); err != nil {
		middleware.RespondWithValidationError(w, "Invalid request body", nil)
		return
	}

	// Validate CI data using the validator
	if validationError := h.validator.Validate(ci); validationError != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(models.GetHTTPStatusForError(models.ErrorTypeValidation))
		json.NewEncoder(w).Encode(validationError)
		return
	}

	// Set default values
	if ci.ID == uuid.Nil {
		ci.ID = uuid.New()
	}
	now := time.Now()
	ci.CreatedAt = now
	ci.UpdatedAt = now

	// Create the CI
	if err := h.ciRepo.Create(r.Context(), &ci); err != nil {
		middleware.RespondWithInternalError(w, "Failed to create CI", nil)
		return
	}

	// Create audit log
	auditLog := &models.AuditLog{
		ID:         uuid.New(),
		EntityType: "configuration_item",
		EntityID:   ci.ID,
		Action:     "create",
		ChangedBy:  username,
		ChangedAt:  time.Now(),
		Details:    models.JSONBMap{"name": ci.Name, "type": ci.Type},
	}
	if err := h.auditRepo.Create(r.Context(), auditLog); err != nil {
		// Log the error but don't fail the request
		// In a real application, you would use a proper logger
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ci)
}

// GetCI handles retrieving a CI by ID
// @Summary Get a CI by ID
// @Description Get a configuration item by its ID
// @Tags cis
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "CI ID"
// @Success 200 {object} models.CI
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /cis/{id} [get]
func (h *CIHandler) GetCI(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL parameters
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		middleware.RespondWithValidationError(w, "ID parameter is required", nil)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		middleware.RespondWithValidationError(w, "Invalid ID format", nil)
		return
	}

	// Get the CI
	ci, err := h.ciRepo.GetByID(r.Context(), id)
	if err != nil {
		middleware.RespondWithNotFoundError(w, "CI not found", nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ci)
}

// GetAllCIs handles retrieving all CIs with pagination
// @Summary Get all CIs
// @Description Get all configuration items with pagination
// @Tags cis
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /cis [get]
func (h *CIHandler) GetAllCIs(w http.ResponseWriter, r *http.Request) {
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

	// Get all CIs
	cis, err := h.ciRepo.GetAll(r.Context())
	if err != nil {
		middleware.RespondWithInternalError(w, "Failed to get CIs", nil)
		return
	}

	// Calculate pagination
	total := len(cis)
	start := (page - 1) * limit
	end := start + limit

	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	paginatedCIs := cis[start:end]

	// Create response
	response := map[string]interface{}{
		"data": paginatedCIs,
		"pagination": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// UpdateCI handles updating an existing CI
// @Summary Update a CI
// @Description Update an existing configuration item
// @Tags cis
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "CI ID"
// @Param ci body models.CI true "Updated CI object"
// @Success 200 {object} models.CI
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /cis/{id} [put]
func (h *CIHandler) UpdateCI(w http.ResponseWriter, r *http.Request) {
	// Get the username from the context
	username, ok := middleware.GetUsernameFromContext(r.Context())
	if !ok {
		middleware.RespondWithUnauthorizedError(w, "User not authenticated", nil)
		return
	}

	// Extract ID from URL parameters
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		middleware.RespondWithValidationError(w, "ID parameter is required", nil)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		middleware.RespondWithValidationError(w, "Invalid ID format", nil)
		return
	}

	// Get the existing CI
	existingCI, err := h.ciRepo.GetByID(r.Context(), id)
	if err != nil {
		middleware.RespondWithNotFoundError(w, "CI not found", nil)
		return
	}

	// Decode the request body
	var updatedCI models.CI
	if err := json.NewDecoder(r.Body).Decode(&updatedCI); err != nil {
		middleware.RespondWithValidationError(w, "Invalid request body", nil)
		return
	}

	// Validate CI data using the validator
	if validationError := h.validator.Validate(updatedCI); validationError != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(models.GetHTTPStatusForError(models.ErrorTypeValidation))
		json.NewEncoder(w).Encode(validationError)
		return
	}

	// Update the CI
	existingCI.Name = updatedCI.Name
	existingCI.Type = updatedCI.Type
	existingCI.Attributes = updatedCI.Attributes
	existingCI.Tags = updatedCI.Tags
	existingCI.UpdatedAt = time.Now()

	if err := h.ciRepo.Update(r.Context(), existingCI); err != nil {
		middleware.RespondWithInternalError(w, "Failed to update CI", nil)
		return
	}

	// Create audit log
	auditLog := &models.AuditLog{
		ID:         uuid.New(),
		EntityType: "configuration_item",
		EntityID:   existingCI.ID,
		Action:     "update",
		ChangedBy:  username,
		ChangedAt:  time.Now(),
		Details:    models.JSONBMap{"name": existingCI.Name, "type": existingCI.Type},
	}
	if err := h.auditRepo.Create(r.Context(), auditLog); err != nil {
		// Log the error but don't fail the request
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existingCI)
}

// DeleteCI handles deleting a CI
// @Summary Delete a CI
// @Description Delete a configuration item
// @Tags cis
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "CI ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /cis/{id} [delete]
func (h *CIHandler) DeleteCI(w http.ResponseWriter, r *http.Request) {
	// Get the username from the context
	username, ok := middleware.GetUsernameFromContext(r.Context())
	if !ok {
		middleware.RespondWithUnauthorizedError(w, "User not authenticated", nil)
		return
	}

	// Extract ID from URL parameters
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		middleware.RespondWithValidationError(w, "ID parameter is required", nil)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		middleware.RespondWithValidationError(w, "Invalid ID format", nil)
		return
	}

	// Get the CI
	ci, err := h.ciRepo.GetByID(r.Context(), id)
	if err != nil {
		middleware.RespondWithNotFoundError(w, "CI not found", nil)
		return
	}

	// Delete the CI
	if err := h.ciRepo.Delete(r.Context(), id); err != nil {
		middleware.RespondWithInternalError(w, "Failed to delete CI", nil)
		return
	}

	// Create audit log
	auditLog := &models.AuditLog{
		ID:         uuid.New(),
		EntityType: "configuration_item",
		EntityID:   ci.ID,
		Action:     "delete",
		ChangedBy:  username,
		ChangedAt:  time.Now(),
		Details:    models.JSONBMap{"name": ci.Name, "type": ci.Type},
	}
	if err := h.auditRepo.Create(r.Context(), auditLog); err != nil {
		// Log the error but don't fail the request
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "CI deleted successfully"})
}

// GetCIGraph handles retrieving related CIs (graph traversal)
// @Summary Get CI graph
// @Description Get related configuration items (graph traversal)
// @Tags cis
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "CI ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /cis/{id}/graph [get]
func (h *CIHandler) GetCIGraph(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL parameters
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		middleware.RespondWithValidationError(w, "ID parameter is required", nil)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		middleware.RespondWithValidationError(w, "Invalid ID format", nil)
		return
	}

	// Get the CI
	ci, err := h.ciRepo.GetByID(r.Context(), id)
	if err != nil {
		middleware.RespondWithNotFoundError(w, "CI not found", nil)
		return
	}

	// Get relationships where this CI is the source
	sourceRels, err := h.relRepo.GetBySourceCI(r.Context(), id)
	if err != nil {
		middleware.RespondWithInternalError(w, "Failed to get relationships", nil)
		return
	}

	// Get relationships where this CI is the target
	targetRels, err := h.relRepo.GetByTargetCI(r.Context(), id)
	if err != nil {
		middleware.RespondWithInternalError(w, "Failed to get relationships", nil)
		return
	}

	// Create response
	response := map[string]interface{}{
		"node": ci,
		"relationships": map[string]interface{}{
			"outgoing": sourceRels,
			"incoming": targetRels,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
