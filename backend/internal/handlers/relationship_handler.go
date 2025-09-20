package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/cmdb-lite/backend/internal/middleware"
	"github.com/cmdb-lite/backend/internal/models"
	"github.com/cmdb-lite/backend/internal/repositories"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// RelationshipHandler handles HTTP requests for relationships
type RelationshipHandler struct {
	relRepo   repositories.RelationshipRepository
	auditRepo repositories.AuditLogRepository
}

// NewRelationshipHandler creates a new RelationshipHandler
func NewRelationshipHandler(
	relRepo repositories.RelationshipRepository,
	auditRepo repositories.AuditLogRepository,
) *RelationshipHandler {
	return &RelationshipHandler{
		relRepo:   relRepo,
		auditRepo: auditRepo,
	}
}

// CreateRelationship handles the creation of a new relationship
// @Summary Create a new relationship
// @Description Create a new relationship between CIs
// @Tags relationships
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param relationship body models.Relationship true "Relationship object"
// @Success 201 {object} models.Relationship
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /relationships [post]
func (h *RelationshipHandler) CreateRelationship(w http.ResponseWriter, r *http.Request) {
	// Get the username from the context
	username, ok := middleware.GetUsernameFromContext(r.Context())
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	var relationship models.Relationship
	if err := json.NewDecoder(r.Body).Decode(&relationship); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate relationship data
	if relationship.SourceID == uuid.Nil || relationship.TargetID == uuid.Nil || relationship.Type == "" {
		http.Error(w, "Source ID, target ID, and type are required", http.StatusBadRequest)
		return
	}

	// Set default values
	if relationship.ID == uuid.Nil {
		relationship.ID = uuid.New()
	}
	relationship.CreatedAt = time.Now()

	// Create the relationship
	if err := h.relRepo.Create(r.Context(), &relationship); err != nil {
		http.Error(w, "Failed to create relationship", http.StatusInternalServerError)
		return
	}

	// Create audit log
	auditLog := &models.AuditLog{
		ID:         uuid.New(),
		EntityType: "relationship",
		EntityID:   relationship.ID,
		Action:     "create",
		ChangedBy:  username,
		ChangedAt:  time.Now(),
		Details:    models.JSONBMap{"source_id": relationship.SourceID, "target_id": relationship.TargetID, "type": relationship.Type},
	}
	if err := h.auditRepo.Create(r.Context(), auditLog); err != nil {
		// Log the error but don't fail the request
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(relationship)
}

// GetRelationship handles retrieving a relationship by ID
// @Summary Get a relationship by ID
// @Description Get a relationship by its ID
// @Tags relationships
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Relationship ID"
// @Success 200 {object} models.Relationship
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /relationships/{id} [get]
func (h *RelationshipHandler) GetRelationship(w http.ResponseWriter, r *http.Request) {
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

	// Get the relationship
	relationship, err := h.relRepo.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Relationship not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(relationship)
}

// GetAllRelationships handles retrieving all relationships
// @Summary Get all relationships
// @Description Get all relationships
// @Tags relationships
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} []models.Relationship
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /relationships [get]
func (h *RelationshipHandler) GetAllRelationships(w http.ResponseWriter, r *http.Request) {
	// Get all relationships
	relationships, err := h.relRepo.GetAll(r.Context())
	if err != nil {
		http.Error(w, "Failed to get relationships", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(relationships)
}

// UpdateRelationship handles updating an existing relationship
// @Summary Update a relationship
// @Description Update an existing relationship
// @Tags relationships
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Relationship ID"
// @Param relationship body models.Relationship true "Updated relationship object"
// @Success 200 {object} models.Relationship
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /relationships/{id} [put]
func (h *RelationshipHandler) UpdateRelationship(w http.ResponseWriter, r *http.Request) {
	// Get the username from the context
	username, ok := middleware.GetUsernameFromContext(r.Context())
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
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

	// Get the existing relationship
	existingRel, err := h.relRepo.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Relationship not found", http.StatusNotFound)
		return
	}

	// Decode the request body
	var updatedRel models.Relationship
	if err := json.NewDecoder(r.Body).Decode(&updatedRel); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate relationship data
	if updatedRel.SourceID == uuid.Nil || updatedRel.TargetID == uuid.Nil || updatedRel.Type == "" {
		http.Error(w, "Source ID, target ID, and type are required", http.StatusBadRequest)
		return
	}

	// Update the relationship
	existingRel.SourceID = updatedRel.SourceID
	existingRel.TargetID = updatedRel.TargetID
	existingRel.Type = updatedRel.Type

	if err := h.relRepo.Update(r.Context(), existingRel); err != nil {
		http.Error(w, "Failed to update relationship", http.StatusInternalServerError)
		return
	}

	// Create audit log
	auditLog := &models.AuditLog{
		ID:         uuid.New(),
		EntityType: "relationship",
		EntityID:   existingRel.ID,
		Action:     "update",
		ChangedBy:  username,
		ChangedAt:  time.Now(),
		Details:    models.JSONBMap{"source_id": existingRel.SourceID, "target_id": existingRel.TargetID, "type": existingRel.Type},
	}
	if err := h.auditRepo.Create(r.Context(), auditLog); err != nil {
		// Log the error but don't fail the request
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existingRel)
}

// DeleteRelationship handles deleting a relationship
// @Summary Delete a relationship
// @Description Delete a relationship
// @Tags relationships
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Relationship ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /relationships/{id} [delete]
func (h *RelationshipHandler) DeleteRelationship(w http.ResponseWriter, r *http.Request) {
	// Get the username from the context
	username, ok := middleware.GetUsernameFromContext(r.Context())
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
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

	// Get the relationship
	rel, err := h.relRepo.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Relationship not found", http.StatusNotFound)
		return
	}

	// Delete the relationship
	if err := h.relRepo.Delete(r.Context(), id); err != nil {
		http.Error(w, "Failed to delete relationship", http.StatusInternalServerError)
		return
	}

	// Create audit log
	auditLog := &models.AuditLog{
		ID:         uuid.New(),
		EntityType: "relationship",
		EntityID:   rel.ID,
		Action:     "delete",
		ChangedBy:  username,
		ChangedAt:  time.Now(),
		Details:    models.JSONBMap{"source_id": rel.SourceID, "target_id": rel.TargetID, "type": rel.Type},
	}
	if err := h.auditRepo.Create(r.Context(), auditLog); err != nil {
		// Log the error but don't fail the request
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Relationship deleted successfully"})
}