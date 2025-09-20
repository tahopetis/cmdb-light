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

func TestRelationshipPostgresRepository_Create(t *testing.T) {
	// Setup test data
	relationshipID := uuid.New()
	sourceID := uuid.New()
	targetID := uuid.New()
	testRelationship := &models.Relationship{
		ID:        relationshipID,
		SourceID:  sourceID,
		TargetID:  targetID,
		Type:      "depends_on",
		CreatedAt: time.Now(),
	}

	tests := []struct {
		name          string
		setupMock     func(*sqlx.DB)
		expectedError bool
	}{
		{
			name: "Successful relationship creation",
			setupMock: func(db *sqlx.DB) {
				db.ExpectExec("INSERT INTO relationships").
					WithArgs(
						testRelationship.ID,
						testRelationship.SourceID,
						testRelationship.TargetID,
						testRelationship.Type,
						testRelationship.CreatedAt,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: false,
		},
		{
			name: "Database error",
			setupMock: func(db *sqlx.DB) {
				db.ExpectExec("INSERT INTO relationships").
					WithArgs(
						testRelationship.ID,
						testRelationship.SourceID,
						testRelationship.TargetID,
						testRelationship.Type,
						testRelationship.CreatedAt,
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
			repo := NewRelationshipPostgresRepository(db)
			
			// Call the method
			err := repo.Create(context.Background(), testRelationship)
			
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

func TestRelationshipPostgresRepository_GetByID(t *testing.T) {
	// Setup test data
	relationshipID := uuid.New()
	sourceID := uuid.New()
	targetID := uuid.New()
	testRelationship := &models.Relationship{
		ID:        relationshipID,
		SourceID:  sourceID,
		TargetID:  targetID,
		Type:      "depends_on",
		CreatedAt: time.Now(),
	}

	tests := []struct {
		name               string
		setupMock          func(*sqlx.DB)
		expectedRelationship *models.Relationship
		expectedError      bool
	}{
		{
			name: "Successful relationship retrieval",
			setupMock: func(db *sqlx.DB) {
				rows := sqlmock.NewRows([]string{"id", "source_id", "target_id", "type", "created_at"}).
					AddRow(testRelationship.ID, testRelationship.SourceID, testRelationship.TargetID, testRelationship.Type, testRelationship.CreatedAt)
				db.ExpectQuery("SELECT id, source_id, target_id, type, created_at FROM relationships WHERE id =").
					WithArgs(relationshipID).
					WillReturnRows(rows)
			},
			expectedRelationship: testRelationship,
			expectedError:        false,
		},
		{
			name: "Relationship not found",
			setupMock: func(db *sqlx.DB) {
				db.ExpectQuery("SELECT id, source_id, target_id, type, created_at FROM relationships WHERE id =").
					WithArgs(relationshipID).
					WillReturnError(sql.ErrNoRows)
			},
			expectedRelationship: nil,
			expectedError:        true,
		},
		{
			name: "Database error",
			setupMock: func(db *sqlx.DB) {
				db.ExpectQuery("SELECT id, source_id, target_id, type, created_at FROM relationships WHERE id =").
					WithArgs(relationshipID).
					WillReturnError(assert.AnError)
			},
			expectedRelationship: nil,
			expectedError:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock database
			db, mock := testutils.NewMockDB(t)
			
			// Setup mock expectations
			tt.setupMock(db)
			
			// Create repository
			repo := NewRelationshipPostgresRepository(db)
			
			// Call the method
			relationship, err := repo.GetByID(context.Background(), relationshipID)
			
			// Assert results
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, relationship)
				assert.Equal(t, tt.expectedRelationship.ID, relationship.ID)
				assert.Equal(t, tt.expectedRelationship.SourceID, relationship.SourceID)
				assert.Equal(t, tt.expectedRelationship.TargetID, relationship.TargetID)
				assert.Equal(t, tt.expectedRelationship.Type, relationship.Type)
			}
			
			// Assert that all expectations were met
			mock.AssertExpectations(t)
		})
	}
}

func TestRelationshipPostgresRepository_GetAll(t *testing.T) {
	// Setup test data
	testRelationships := []*models.Relationship{
		{
			ID:        uuid.New(),
			SourceID:  uuid.New(),
			TargetID:  uuid.New(),
			Type:      "depends_on",
			CreatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			SourceID:  uuid.New(),
			TargetID:  uuid.New(),
			Type:      "connected_to",
			CreatedAt: time.Now(),
		},
	}

	tests := []struct {
		name               string
		setupMock          func(*sqlx.DB)
		expectedRelationships []*models.Relationship
		expectedError      bool
	}{
		{
			name: "Successful relationship retrieval",
			setupMock: func(db *sqlx.DB) {
				rows := sqlmock.NewRows([]string{"id", "source_id", "target_id", "type", "created_at"}).
					AddRow(testRelationships[0].ID, testRelationships[0].SourceID, testRelationships[0].TargetID, testRelationships[0].Type, testRelationships[0].CreatedAt).
					AddRow(testRelationships[1].ID, testRelationships[1].SourceID, testRelationships[1].TargetID, testRelationships[1].Type, testRelationships[1].CreatedAt)
				db.ExpectQuery("SELECT id, source_id, target_id, type, created_at FROM relationships").
					WillReturnRows(rows)
			},
			expectedRelationships: testRelationships,
			expectedError:         false,
		},
		{
			name: "No relationships found",
			setupMock: func(db *sqlx.DB) {
				rows := sqlmock.NewRows([]string{"id", "source_id", "target_id", "type", "created_at"})
				db.ExpectQuery("SELECT id, source_id, target_id, type, created_at FROM relationships").
					WillReturnRows(rows)
			},
			expectedRelationships: []*models.Relationship{},
			expectedError:         false,
		},
		{
			name: "Database error",
			setupMock: func(db *sqlx.DB) {
				db.ExpectQuery("SELECT id, source_id, target_id, type, created_at FROM relationships").
					WillReturnError(assert.AnError)
			},
			expectedRelationships: nil,
			expectedError:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock database
			db, mock := testutils.NewMockDB(t)
			
			// Setup mock expectations
			tt.setupMock(db)
			
			// Create repository
			repo := NewRelationshipPostgresRepository(db)
			
			// Call the method
			relationships, err := repo.GetAll(context.Background())
			
			// Assert results
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.expectedRelationships), len(relationships))
				
				if len(tt.expectedRelationships) > 0 {
					for i, expectedRelationship := range tt.expectedRelationships {
						assert.Equal(t, expectedRelationship.ID, relationships[i].ID)
						assert.Equal(t, expectedRelationship.SourceID, relationships[i].SourceID)
						assert.Equal(t, expectedRelationship.TargetID, relationships[i].TargetID)
						assert.Equal(t, expectedRelationship.Type, relationships[i].Type)
					}
				}
			}
			
			// Assert that all expectations were met
			mock.AssertExpectations(t)
		})
	}
}

func TestRelationshipPostgresRepository_GetBySourceID(t *testing.T) {
	// Setup test data
	sourceID := uuid.New()
	testRelationships := []*models.Relationship{
		{
			ID:        uuid.New(),
			SourceID:  sourceID,
			TargetID:  uuid.New(),
			Type:      "depends_on",
			CreatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			SourceID:  sourceID,
			TargetID:  uuid.New(),
			Type:      "connected_to",
			CreatedAt: time.Now(),
		},
	}

	tests := []struct {
		name               string
		setupMock          func(*sqlx.DB)
		expectedRelationships []*models.Relationship
		expectedError      bool
	}{
		{
			name: "Successful relationship retrieval by source ID",
			setupMock: func(db *sqlx.DB) {
				rows := sqlmock.NewRows([]string{"id", "source_id", "target_id", "type", "created_at"}).
					AddRow(testRelationships[0].ID, testRelationships[0].SourceID, testRelationships[0].TargetID, testRelationships[0].Type, testRelationships[0].CreatedAt).
					AddRow(testRelationships[1].ID, testRelationships[1].SourceID, testRelationships[1].TargetID, testRelationships[1].Type, testRelationships[1].CreatedAt)
				db.ExpectQuery("SELECT id, source_id, target_id, type, created_at FROM relationships WHERE source_id =").
					WithArgs(sourceID).
					WillReturnRows(rows)
			},
			expectedRelationships: testRelationships,
			expectedError:         false,
		},
		{
			name: "No relationships found for source ID",
			setupMock: func(db *sqlx.DB) {
				rows := sqlmock.NewRows([]string{"id", "source_id", "target_id", "type", "created_at"})
				db.ExpectQuery("SELECT id, source_id, target_id, type, created_at FROM relationships WHERE source_id =").
					WithArgs(sourceID).
					WillReturnRows(rows)
			},
			expectedRelationships: []*models.Relationship{},
			expectedError:         false,
		},
		{
			name: "Database error",
			setupMock: func(db *sqlx.DB) {
				db.ExpectQuery("SELECT id, source_id, target_id, type, created_at FROM relationships WHERE source_id =").
					WithArgs(sourceID).
					WillReturnError(assert.AnError)
			},
			expectedRelationships: nil,
			expectedError:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock database
			db, mock := testutils.NewMockDB(t)
			
			// Setup mock expectations
			tt.setupMock(db)
			
			// Create repository
			repo := NewRelationshipPostgresRepository(db)
			
			// Call the method
			relationships, err := repo.GetBySourceID(context.Background(), sourceID)
			
			// Assert results
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.expectedRelationships), len(relationships))
				
				if len(tt.expectedRelationships) > 0 {
					for i, expectedRelationship := range tt.expectedRelationships {
						assert.Equal(t, expectedRelationship.ID, relationships[i].ID)
						assert.Equal(t, expectedRelationship.SourceID, relationships[i].SourceID)
						assert.Equal(t, expectedRelationship.TargetID, relationships[i].TargetID)
						assert.Equal(t, expectedRelationship.Type, relationships[i].Type)
					}
				}
			}
			
			// Assert that all expectations were met
			mock.AssertExpectations(t)
		})
	}
}

func TestRelationshipPostgresRepository_GetByTargetID(t *testing.T) {
	// Setup test data
	targetID := uuid.New()
	testRelationships := []*models.Relationship{
		{
			ID:        uuid.New(),
			SourceID:  uuid.New(),
			TargetID:  targetID,
			Type:      "depends_on",
			CreatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			SourceID:  uuid.New(),
			TargetID:  targetID,
			Type:      "connected_to",
			CreatedAt: time.Now(),
		},
	}

	tests := []struct {
		name               string
		setupMock          func(*sqlx.DB)
		expectedRelationships []*models.Relationship
		expectedError      bool
	}{
		{
			name: "Successful relationship retrieval by target ID",
			setupMock: func(db *sqlx.DB) {
				rows := sqlmock.NewRows([]string{"id", "source_id", "target_id", "type", "created_at"}).
					AddRow(testRelationships[0].ID, testRelationships[0].SourceID, testRelationships[0].TargetID, testRelationships[0].Type, testRelationships[0].CreatedAt).
					AddRow(testRelationships[1].ID, testRelationships[1].SourceID, testRelationships[1].TargetID, testRelationships[1].Type, testRelationships[1].CreatedAt)
				db.ExpectQuery("SELECT id, source_id, target_id, type, created_at FROM relationships WHERE target_id =").
					WithArgs(targetID).
					WillReturnRows(rows)
			},
			expectedRelationships: testRelationships,
			expectedError:         false,
		},
		{
			name: "No relationships found for target ID",
			setupMock: func(db *sqlx.DB) {
				rows := sqlmock.NewRows([]string{"id", "source_id", "target_id", "type", "created_at"})
				db.ExpectQuery("SELECT id, source_id, target_id, type, created_at FROM relationships WHERE target_id =").
					WithArgs(targetID).
					WillReturnRows(rows)
			},
			expectedRelationships: []*models.Relationship{},
			expectedError:         false,
		},
		{
			name: "Database error",
			setupMock: func(db *sqlx.DB) {
				db.ExpectQuery("SELECT id, source_id, target_id, type, created_at FROM relationships WHERE target_id =").
					WithArgs(targetID).
					WillReturnError(assert.AnError)
			},
			expectedRelationships: nil,
			expectedError:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock database
			db, mock := testutils.NewMockDB(t)
			
			// Setup mock expectations
			tt.setupMock(db)
			
			// Create repository
			repo := NewRelationshipPostgresRepository(db)
			
			// Call the method
			relationships, err := repo.GetByTargetID(context.Background(), targetID)
			
			// Assert results
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.expectedRelationships), len(relationships))
				
				if len(tt.expectedRelationships) > 0 {
					for i, expectedRelationship := range tt.expectedRelationships {
						assert.Equal(t, expectedRelationship.ID, relationships[i].ID)
						assert.Equal(t, expectedRelationship.SourceID, relationships[i].SourceID)
						assert.Equal(t, expectedRelationship.TargetID, relationships[i].TargetID)
						assert.Equal(t, expectedRelationship.Type, relationships[i].Type)
					}
				}
			}
			
			// Assert that all expectations were met
			mock.AssertExpectations(t)
		})
	}
}

func TestRelationshipPostgresRepository_Update(t *testing.T) {
	// Setup test data
	relationshipID := uuid.New()
	sourceID := uuid.New()
	targetID := uuid.New()
	testRelationship := &models.Relationship{
		ID:        relationshipID,
		SourceID:  sourceID,
		TargetID:  targetID,
		Type:      "depends_on",
		CreatedAt: time.Now(),
	}

	tests := []struct {
		name          string
		setupMock     func(*sqlx.DB)
		expectedError bool
	}{
		{
			name: "Successful relationship update",
			setupMock: func(db *sqlx.DB) {
				db.ExpectExec("UPDATE relationships").
					WithArgs(
						testRelationship.ID,
						testRelationship.SourceID,
						testRelationship.TargetID,
						testRelationship.Type,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: false,
		},
		{
			name: "Relationship not found",
			setupMock: func(db *sqlx.DB) {
				db.ExpectExec("UPDATE relationships").
					WithArgs(
						testRelationship.ID,
						testRelationship.SourceID,
						testRelationship.TargetID,
						testRelationship.Type,
					).
					WillReturnResult(sqlmock.NewResult(1, 0))
			},
			expectedError: true,
		},
		{
			name: "Database error",
			setupMock: func(db *sqlx.DB) {
				db.ExpectExec("UPDATE relationships").
					WithArgs(
						testRelationship.ID,
						testRelationship.SourceID,
						testRelationship.TargetID,
						testRelationship.Type,
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
			repo := NewRelationshipPostgresRepository(db)
			
			// Call the method
			err := repo.Update(context.Background(), testRelationship)
			
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

func TestRelationshipPostgresRepository_Delete(t *testing.T) {
	// Setup test data
	relationshipID := uuid.New()

	tests := []struct {
		name          string
		setupMock     func(*sqlx.DB)
		expectedError bool
	}{
		{
			name: "Successful relationship deletion",
			setupMock: func(db *sqlx.DB) {
				db.ExpectExec("DELETE FROM relationships WHERE id =").
					WithArgs(relationshipID).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: false,
		},
		{
			name: "Relationship not found",
			setupMock: func(db *sqlx.DB) {
				db.ExpectExec("DELETE FROM relationships WHERE id =").
					WithArgs(relationshipID).
					WillReturnResult(sqlmock.NewResult(1, 0))
			},
			expectedError: true,
		},
		{
			name: "Database error",
			setupMock: func(db *sqlx.DB) {
				db.ExpectExec("DELETE FROM relationships WHERE id =").
					WithArgs(relationshipID).
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
			repo := NewRelationshipPostgresRepository(db)
			
			// Call the method
			err := repo.Delete(context.Background(), relationshipID)
			
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