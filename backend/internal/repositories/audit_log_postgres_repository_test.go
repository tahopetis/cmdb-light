package repositories

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/cmdb-lite/backend/internal/models"
	"github.com/cmdb-lite/backend/internal/testutils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuditLogPostgresRepository_Create(t *testing.T) {
	// Setup test data
	auditLogID := uuid.New()
	entityID := uuid.New()
	testAuditLog := &models.AuditLog{
		ID:         auditLogID,
		EntityType: "configuration_item",
		EntityID:   entityID,
		Action:     "create",
		ChangedBy:  "testuser",
		ChangedAt:  time.Now(),
		Details:    models.JSONBMap{"name": "Test CI", "type": "server"},
	}

	tests := []struct {
		name          string
		setupMock     func(*sqlx.DB)
		expectedError bool
	}{
		{
			name: "Successful audit log creation",
			setupMock: func(db *sqlx.DB) {
				db.ExpectExec("INSERT INTO audit_logs").
					WithArgs(
						testAuditLog.ID,
						testAuditLog.EntityType,
						testAuditLog.EntityID,
						testAuditLog.Action,
						testAuditLog.ChangedBy,
						testAuditLog.ChangedAt,
						testAuditLog.Details,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: false,
		},
		{
			name: "Database error",
			setupMock: func(db *sqlx.DB) {
				db.ExpectExec("INSERT INTO audit_logs").
					WithArgs(
						testAuditLog.ID,
						testAuditLog.EntityType,
						testAuditLog.EntityID,
						testAuditLog.Action,
						testAuditLog.ChangedBy,
						testAuditLog.ChangedAt,
						testAuditLog.Details,
					).
					WillReturnError(assert.AnError)
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock database
			db, mock := testutils.NewMockDB(t)
			
			// Setup mock expectations
			tt.setupMock(db)
			
			// Create repository
			repo := NewAuditLogPostgresRepository(db)
			
			// Call the method
			err := repo.Create(context.Background(), testAuditLog)
			
			// Assert results
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			
			// Assert that all expectations were met
			mock.AssertExpectations(t)
		})
	}
}

func TestAuditLogPostgresRepository_GetByID(t *testing.T) {
	// Setup test data
	auditLogID := uuid.New()
	entityID := uuid.New()
	testAuditLog := &models.AuditLog{
		ID:         auditLogID,
		EntityType: "configuration_item",
		EntityID:   entityID,
		Action:     "create",
		ChangedBy:  "testuser",
		ChangedAt:  time.Now(),
		Details:    models.JSONBMap{"name": "Test CI", "type": "server"},
	}

	tests := []struct {
		name           string
		setupMock      func(*sqlx.DB)
		expectedAuditLog *models.AuditLog
		expectedError  bool
	}{
		{
			name: "Successful audit log retrieval",
			setupMock: func(db *sqlx.DB) {
				rows := sqlmock.NewRows([]string{"id", "entity_type", "entity_id", "action", "changed_by", "changed_at", "details"}).
					AddRow(testAuditLog.ID, testAuditLog.EntityType, testAuditLog.EntityID, testAuditLog.Action, testAuditLog.ChangedBy, testAuditLog.ChangedAt, testAuditLog.Details)
				db.ExpectQuery("SELECT id, entity_type, entity_id, action, changed_by, changed_at, details FROM audit_logs WHERE id =").
					WithArgs(auditLogID).
					WillReturnRows(rows)
			},
			expectedAuditLog: testAuditLog,
			expectedError:   false,
		},
		{
			name: "Audit log not found",
			setupMock: func(db *sqlx.DB) {
				db.ExpectQuery("SELECT id, entity_type, entity_id, action, changed_by, changed_at, details FROM audit_logs WHERE id =").
					WithArgs(auditLogID).
					WillReturnError(sql.ErrNoRows)
			},
			expectedAuditLog: nil,
			expectedError:   true,
		},
		{
			name: "Database error",
			setupMock: func(db *sqlx.DB) {
				db.ExpectQuery("SELECT id, entity_type, entity_id, action, changed_by, changed_at, details FROM audit_logs WHERE id =").
					WithArgs(auditLogID).
					WillReturnError(assert.AnError)
			},
			expectedAuditLog: nil,
			expectedError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock database
			db, mock := testutils.NewMockDB(t)
			
			// Setup mock expectations
			tt.setupMock(db)
			
			// Create repository
			repo := NewAuditLogPostgresRepository(db)
			
			// Call the method
			auditLog, err := repo.GetByID(context.Background(), auditLogID)
			
			// Assert results
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, auditLog)
				assert.Equal(t, tt.expectedAuditLog.ID, auditLog.ID)
				assert.Equal(t, tt.expectedAuditLog.EntityType, auditLog.EntityType)
				assert.Equal(t, tt.expectedAuditLog.EntityID, auditLog.EntityID)
				assert.Equal(t, tt.expectedAuditLog.Action, auditLog.Action)
				assert.Equal(t, tt.expectedAuditLog.ChangedBy, auditLog.ChangedBy)
				assert.Equal(t, tt.expectedAuditLog.Details, auditLog.Details)
			}
			
			// Assert that all expectations were met
			mock.AssertExpectations(t)
		})
	}
}

func TestAuditLogPostgresRepository_GetAll(t *testing.T) {
	// Setup test data
	testAuditLogs := []*models.AuditLog{
		{
			ID:         uuid.New(),
			EntityType: "configuration_item",
			EntityID:   uuid.New(),
			Action:     "create",
			ChangedBy:  "testuser",
			ChangedAt:  time.Now(),
			Details:    models.JSONBMap{"name": "Test CI 1", "type": "server"},
		},
		{
			ID:         uuid.New(),
			EntityType: "configuration_item",
			EntityID:   uuid.New(),
			Action:     "update",
			ChangedBy:  "testuser",
			ChangedAt:  time.Now(),
			Details:    models.JSONBMap{"name": "Test CI 2", "type": "database"},
		},
	}

	tests := []struct {
		name           string
		setupMock      func(*sqlx.DB)
		expectedAuditLogs []*models.AuditLog
		expectedError  bool
	}{
		{
			name: "Successful audit log retrieval",
			setupMock: func(db *sqlx.DB) {
				rows := sqlmock.NewRows([]string{"id", "entity_type", "entity_id", "action", "changed_by", "changed_at", "details"}).
					AddRow(testAuditLogs[0].ID, testAuditLogs[0].EntityType, testAuditLogs[0].EntityID, testAuditLogs[0].Action, testAuditLogs[0].ChangedBy, testAuditLogs[0].ChangedAt, testAuditLogs[0].Details).
					AddRow(testAuditLogs[1].ID, testAuditLogs[1].EntityType, testAuditLogs[1].EntityID, testAuditLogs[1].Action, testAuditLogs[1].ChangedBy, testAuditLogs[1].ChangedAt, testAuditLogs[1].Details)
				db.ExpectQuery("SELECT id, entity_type, entity_id, action, changed_by, changed_at, details FROM audit_logs").
					WillReturnRows(rows)
			},
			expectedAuditLogs: testAuditLogs,
			expectedError:     false,
		},
		{
			name: "No audit logs found",
			setupMock: func(db *sqlx.DB) {
				rows := sqlmock.NewRows([]string{"id", "entity_type", "entity_id", "action", "changed_by", "changed_at", "details"})
				db.ExpectQuery("SELECT id, entity_type, entity_id, action, changed_by, changed_at, details FROM audit_logs").
					WillReturnRows(rows)
			},
			expectedAuditLogs: []*models.AuditLog{},
			expectedError:     false,
		},
		{
			name: "Database error",
			setupMock: func(db *sqlx.DB) {
				db.ExpectQuery("SELECT id, entity_type, entity_id, action, changed_by, changed_at, details FROM audit_logs").
					WillReturnError(assert.AnError)
			},
			expectedAuditLogs: nil,
			expectedError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock database
			db, mock := testutils.NewMockDB(t)
			
			// Setup mock expectations
			tt.setupMock(db)
			
			// Create repository
			repo := NewAuditLogPostgresRepository(db)
			
			// Call the method
			auditLogs, err := repo.GetAll(context.Background())
			
			// Assert results
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.expectedAuditLogs), len(auditLogs))
				
				if len(tt.expectedAuditLogs) > 0 {
					for i, expectedAuditLog := range tt.expectedAuditLogs {
						assert.Equal(t, expectedAuditLog.ID, auditLogs[i].ID)
						assert.Equal(t, expectedAuditLog.EntityType, auditLogs[i].EntityType)
						assert.Equal(t, expectedAuditLog.EntityID, auditLogs[i].EntityID)
						assert.Equal(t, expectedAuditLog.Action, auditLogs[i].Action)
						assert.Equal(t, expectedAuditLog.ChangedBy, auditLogs[i].ChangedBy)
						assert.Equal(t, expectedAuditLog.Details, auditLogs[i].Details)
					}
				}
			}
			
			// Assert that all expectations were met
			mock.AssertExpectations(t)
		})
	}
}

func TestAuditLogPostgresRepository_GetByEntityType(t *testing.T) {
	// Setup test data
	entityType := "configuration_item"
	testAuditLogs := []*models.AuditLog{
		{
			ID:         uuid.New(),
			EntityType: entityType,
			EntityID:   uuid.New(),
			Action:     "create",
			ChangedBy:  "testuser",
			ChangedAt:  time.Now(),
			Details:    models.JSONBMap{"name": "Test CI 1", "type": "server"},
		},
		{
			ID:         uuid.New(),
			EntityType: entityType,
			EntityID:   uuid.New(),
			Action:     "update",
			ChangedBy:  "testuser",
			ChangedAt:  time.Now(),
			Details:    models.JSONBMap{"name": "Test CI 2", "type": "database"},
		},
	}

	tests := []struct {
		name           string
		setupMock      func(*sqlx.DB)
		expectedAuditLogs []*models.AuditLog
		expectedError  bool
	}{
		{
			name: "Successful audit log retrieval by entity type",
			setupMock: func(db *sqlx.DB) {
				rows := sqlmock.NewRows([]string{"id", "entity_type", "entity_id", "action", "changed_by", "changed_at", "details"}).
					AddRow(testAuditLogs[0].ID, testAuditLogs[0].EntityType, testAuditLogs[0].EntityID, testAuditLogs[0].Action, testAuditLogs[0].ChangedBy, testAuditLogs[0].ChangedAt, testAuditLogs[0].Details).
					AddRow(testAuditLogs[1].ID, testAuditLogs[1].EntityType, testAuditLogs[1].EntityID, testAuditLogs[1].Action, testAuditLogs[1].ChangedBy, testAuditLogs[1].ChangedAt, testAuditLogs[1].Details)
				db.ExpectQuery("SELECT id, entity_type, entity_id, action, changed_by, changed_at, details FROM audit_logs WHERE entity_type =").
					WithArgs(entityType).
					WillReturnRows(rows)
			},
			expectedAuditLogs: testAuditLogs,
			expectedError:     false,
		},
		{
			name: "No audit logs found for entity type",
			setupMock: func(db *sqlx.DB) {
				rows := sqlmock.NewRows([]string{"id", "entity_type", "entity_id", "action", "changed_by", "changed_at", "details"})
				db.ExpectQuery("SELECT id, entity_type, entity_id, action, changed_by, changed_at, details FROM audit_logs WHERE entity_type =").
					WithArgs(entityType).
					WillReturnRows(rows)
			},
			expectedAuditLogs: []*models.AuditLog{},
			expectedError:     false,
		},
		{
			name: "Database error",
			setupMock: func(db *sqlx.DB) {
				db.ExpectQuery("SELECT id, entity_type, entity_id, action, changed_by, changed_at, details FROM audit_logs WHERE entity_type =").
					WithArgs(entityType).
					WillReturnError(assert.AnError)
			},
			expectedAuditLogs: nil,
			expectedError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock database
			db, mock := testutils.NewMockDB(t)
			
			// Setup mock expectations
			tt.setupMock(db)
			
			// Create repository
			repo := NewAuditLogPostgresRepository(db)
			
			// Call the method
			auditLogs, err := repo.GetByEntityType(context.Background(), entityType)
			
			// Assert results
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.expectedAuditLogs), len(auditLogs))
				
				if len(tt.expectedAuditLogs) > 0 {
					for i, expectedAuditLog := range tt.expectedAuditLogs {
						assert.Equal(t, expectedAuditLog.ID, auditLogs[i].ID)
						assert.Equal(t, expectedAuditLog.EntityType, auditLogs[i].EntityType)
						assert.Equal(t, expectedAuditLog.EntityID, auditLogs[i].EntityID)
						assert.Equal(t, expectedAuditLog.Action, auditLogs[i].Action)
						assert.Equal(t, expectedAuditLog.ChangedBy, auditLogs[i].ChangedBy)
						assert.Equal(t, expectedAuditLog.Details, auditLogs[i].Details)
					}
				}
			}
			
			// Assert that all expectations were met
			mock.AssertExpectations(t)
		})
	}
}

func TestAuditLogPostgresRepository_GetByEntityID(t *testing.T) {
	// Setup test data
	entityID := uuid.New()
	testAuditLogs := []*models.AuditLog{
		{
			ID:         uuid.New(),
			EntityType: "configuration_item",
			EntityID:   entityID,
			Action:     "create",
			ChangedBy:  "testuser",
			ChangedAt:  time.Now(),
			Details:    models.JSONBMap{"name": "Test CI", "type": "server"},
		},
		{
			ID:         uuid.New(),
			EntityType: "configuration_item",
			EntityID:   entityID,
			Action:     "update",
			ChangedBy:  "testuser",
			ChangedAt:  time.Now(),
			Details:    models.JSONBMap{"name": "Test CI Updated", "type": "server"},
		},
	}

	tests := []struct {
		name           string
		setupMock      func(*sqlx.DB)
		expectedAuditLogs []*models.AuditLog
		expectedError  bool
	}{
		{
			name: "Successful audit log retrieval by entity ID",
			setupMock: func(db *sqlx.DB) {
				rows := sqlmock.NewRows([]string{"id", "entity_type", "entity_id", "action", "changed_by", "changed_at", "details"}).
					AddRow(testAuditLogs[0].ID, testAuditLogs[0].EntityType, testAuditLogs[0].EntityID, testAuditLogs[0].Action, testAuditLogs[0].ChangedBy, testAuditLogs[0].ChangedAt, testAuditLogs[0].Details).
					AddRow(testAuditLogs[1].ID, testAuditLogs[1].EntityType, testAuditLogs[1].EntityID, testAuditLogs[1].Action, testAuditLogs[1].ChangedBy, testAuditLogs[1].ChangedAt, testAuditLogs[1].Details)
				db.ExpectQuery("SELECT id, entity_type, entity_id, action, changed_by, changed_at, details FROM audit_logs WHERE entity_id =").
					WithArgs(entityID).
					WillReturnRows(rows)
			},
			expectedAuditLogs: testAuditLogs,
			expectedError:     false,
		},
		{
			name: "No audit logs found for entity ID",
			setupMock: func(db *sqlx.DB) {
				rows := sqlmock.NewRows([]string{"id", "entity_type", "entity_id", "action", "changed_by", "changed_at", "details"})
				db.ExpectQuery("SELECT id, entity_type, entity_id, action, changed_by, changed_at, details FROM audit_logs WHERE entity_id =").
					WithArgs(entityID).
					WillReturnRows(rows)
			},
			expectedAuditLogs: []*models.AuditLog{},
			expectedError:     false,
		},
		{
			name: "Database error",
			setupMock: func(db *sqlx.DB) {
				db.ExpectQuery("SELECT id, entity_type, entity_id, action, changed_by, changed_at, details FROM audit_logs WHERE entity_id =").
					WithArgs(entityID).
					WillReturnError(assert.AnError)
			},
			expectedAuditLogs: nil,
			expectedError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock database
			db, mock := testutils.NewMockDB(t)
			
			// Setup mock expectations
			tt.setupMock(db)
			
			// Create repository
			repo := NewAuditLogPostgresRepository(db)
			
			// Call the method
			auditLogs, err := repo.GetByEntityID(context.Background(), entityID)
			
			// Assert results
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.expectedAuditLogs), len(auditLogs))
				
				if len(tt.expectedAuditLogs) > 0 {
					for i, expectedAuditLog := range tt.expectedAuditLogs {
						assert.Equal(t, expectedAuditLog.ID, auditLogs[i].ID)
						assert.Equal(t, expectedAuditLog.EntityType, auditLogs[i].EntityType)
						assert.Equal(t, expectedAuditLog.EntityID, auditLogs[i].EntityID)
						assert.Equal(t, expectedAuditLog.Action, auditLogs[i].Action)
						assert.Equal(t, expectedAuditLog.ChangedBy, auditLogs[i].ChangedBy)
						assert.Equal(t, expectedAuditLog.Details, auditLogs[i].Details)
					}
				}
			}
			
			// Assert that all expectations were met
			mock.AssertExpectations(t)
		})
	}
}

func TestAuditLogPostgresRepository_GetByChangedBy(t *testing.T) {
	// Setup test data
	changedBy := "testuser"
	testAuditLogs := []*models.AuditLog{
		{
			ID:         uuid.New(),
			EntityType: "configuration_item",
			EntityID:   uuid.New(),
			Action:     "create",
			ChangedBy:  changedBy,
			ChangedAt:  time.Now(),
			Details:    models.JSONBMap{"name": "Test CI 1", "type": "server"},
		},
		{
			ID:         uuid.New(),
			EntityType: "configuration_item",
			EntityID:   uuid.New(),
			Action:     "update",
			ChangedBy:  changedBy,
			ChangedAt:  time.Now(),
			Details:    models.JSONBMap{"name": "Test CI 2", "type": "database"},
		},
	}

	tests := []struct {
		name           string
		setupMock      func(*sqlx.DB)
		expectedAuditLogs []*models.AuditLog
		expectedError  bool
	}{
		{
			name: "Successful audit log retrieval by changed by",
			setupMock: func(db *sqlx.DB) {
				rows := sqlmock.NewRows([]string{"id", "entity_type", "entity_id", "action", "changed_by", "changed_at", "details"}).
					AddRow(testAuditLogs[0].ID, testAuditLogs[0].EntityType, testAuditLogs[0].EntityID, testAuditLogs[0].Action, testAuditLogs[0].ChangedBy, testAuditLogs[0].ChangedAt, testAuditLogs[0].Details).
					AddRow(testAuditLogs[1].ID, testAuditLogs[1].EntityType, testAuditLogs[1].EntityID, testAuditLogs[1].Action, testAuditLogs[1].ChangedBy, testAuditLogs[1].ChangedAt, testAuditLogs[1].Details)
				db.ExpectQuery("SELECT id, entity_type, entity_id, action, changed_by, changed_at, details FROM audit_logs WHERE changed_by =").
					WithArgs(changedBy).
					WillReturnRows(rows)
			},
			expectedAuditLogs: testAuditLogs,
			expectedError:     false,
		},
		{
			name: "No audit logs found for changed by",
			setupMock: func(db *sqlx.DB) {
				rows := sqlmock.NewRows([]string{"id", "entity_type", "entity_id", "action", "changed_by", "changed_at", "details"})
				db.ExpectQuery("SELECT id, entity_type, entity_id, action, changed_by, changed_at, details FROM audit_logs WHERE changed_by =").
					WithArgs(changedBy).
					WillReturnRows(rows)
			},
			expectedAuditLogs: []*models.AuditLog{},
			expectedError:     false,
		},
		{
			name: "Database error",
			setupMock: func(db *sqlx.DB) {
				db.ExpectQuery("SELECT id, entity_type, entity_id, action, changed_by, changed_at, details FROM audit_logs WHERE changed_by =").
					WithArgs(changedBy).
					WillReturnError(assert.AnError)
			},
			expectedAuditLogs: nil,
			expectedError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock database
			db, mock := testutils.NewMockDB(t)
			
			// Setup mock expectations
			tt.setupMock(db)
			
			// Create repository
			repo := NewAuditLogPostgresRepository(db)
			
			// Call the method
			auditLogs, err := repo.GetByChangedBy(context.Background(), changedBy)
			
			// Assert results
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.expectedAuditLogs), len(auditLogs))
				
				if len(tt.expectedAuditLogs) > 0 {
					for i, expectedAuditLog := range tt.expectedAuditLogs {
						assert.Equal(t, expectedAuditLog.ID, auditLogs[i].ID)
						assert.Equal(t, expectedAuditLog.EntityType, auditLogs[i].EntityType)
						assert.Equal(t, expectedAuditLog.EntityID, auditLogs[i].EntityID)
						assert.Equal(t, expectedAuditLog.Action, auditLogs[i].Action)
						assert.Equal(t, expectedAuditLog.ChangedBy, auditLogs[i].ChangedBy)
						assert.Equal(t, expectedAuditLog.Details, auditLogs[i].Details)
					}
				}
			}
			
			// Assert that all expectations were met
			mock.AssertExpectations(t)
		})
	}
}

func TestAuditLogPostgresRepository_Delete(t *testing.T) {
	// Setup test data
	auditLogID := uuid.New()

	tests := []struct {
		name          string
		setupMock     func(*sqlx.DB)
		expectedError bool
	}{
		{
			name: "Successful audit log deletion",
			setupMock: func(db *sqlx.DB) {
				db.ExpectExec("DELETE FROM audit_logs WHERE id =").
					WithArgs(auditLogID).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: false,
		},
		{
			name: "Audit log not found",
			setupMock: func(db *sqlx.DB) {
				db.ExpectExec("DELETE FROM audit_logs WHERE id =").
					WithArgs(auditLogID).
					WillReturnResult(sqlmock.NewResult(1, 0))
			},
			expectedError: true,
		},
		{
			name: "Database error",
			setupMock: func(db *sqlx.DB) {
				db.ExpectExec("DELETE FROM audit_logs WHERE id =").
					WithArgs(auditLogID).
					WillReturnError(assert.AnError)
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock database
			db, mock := testutils.NewMockDB(t)
			
			// Setup mock expectations
			tt.setupMock(db)
			
			// Create repository
			repo := NewAuditLogPostgresRepository(db)
			
			// Call the method
			err := repo.Delete(context.Background(), auditLogID)
			
			// Assert results
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			
			// Assert that all expectations were met
			mock.AssertExpectations(t)
		})
	}
}