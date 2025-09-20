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

func TestCIPostgresRepository_Create(t *testing.T) {
	// Setup test data
	ciID := uuid.New()
	testCI := &models.CI{
		ID:         ciID,
		Name:       "Test CI",
		Type:       "server",
		Attributes: models.JSONBMap{"os": "linux", "cpu": "4"},
		Tags:       []string{"production", "web"},
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	tests := []struct {
		name          string
		setupMock     func(*sqlx.DB)
		expectedError bool
	}{
		{
			name: "Successful CI creation",
			setupMock: func(db *sqlx.DB) {
				db.ExpectExec("INSERT INTO configuration_items").
					WithArgs(
						testCI.ID,
						testCI.Name,
						testCI.Type,
						testCI.Attributes,
						testCI.Tags,
						testCI.CreatedAt,
						testCI.UpdatedAt,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: false,
		},
		{
			name: "Database error",
			setupMock: func(db *sqlx.DB) {
				db.ExpectExec("INSERT INTO configuration_items").
					WithArgs(
						testCI.ID,
						testCI.Name,
						testCI.Type,
						testCI.Attributes,
						testCI.Tags,
						testCI.CreatedAt,
						testCI.UpdatedAt,
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
			repo := NewCIPostgresRepository(db)
			
			// Call the method
			err := repo.Create(context.Background(), testCI)
			
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

func TestCIPostgresRepository_GetByID(t *testing.T) {
	// Setup test data
	ciID := uuid.New()
	testCI := &models.CI{
		ID:         ciID,
		Name:       "Test CI",
		Type:       "server",
		Attributes: models.JSONBMap{"os": "linux", "cpu": "4"},
		Tags:       []string{"production", "web"},
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	tests := []struct {
		name          string
		setupMock     func(*sqlx.DB)
		expectedCI    *models.CI
		expectedError bool
	}{
		{
			name: "Successful CI retrieval",
			setupMock: func(db *sqlx.DB) {
				rows := sqlmock.NewRows([]string{"id", "name", "type", "attributes", "tags", "created_at", "updated_at"}).
					AddRow(testCI.ID, testCI.Name, testCI.Type, testCI.Attributes, testCI.Tags, testCI.CreatedAt, testCI.UpdatedAt)
				db.ExpectQuery("SELECT id, name, type, attributes, tags, created_at, updated_at FROM configuration_items WHERE id =").
					WithArgs(ciID).
					WillReturnRows(rows)
			},
			expectedCI:    testCI,
			expectedError: false,
		},
		{
			name: "CI not found",
			setupMock: func(db *sqlx.DB) {
				db.ExpectQuery("SELECT id, name, type, attributes, tags, created_at, updated_at FROM configuration_items WHERE id =").
					WithArgs(ciID).
					WillReturnError(sql.ErrNoRows)
			},
			expectedCI:    nil,
			expectedError: true,
		},
		{
			name: "Database error",
			setupMock: func(db *sqlx.DB) {
				db.ExpectQuery("SELECT id, name, type, attributes, tags, created_at, updated_at FROM configuration_items WHERE id =").
					WithArgs(ciID).
					WillReturnError(assert.AnError)
			},
			expectedCI:    nil,
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
			repo := NewCIPostgresRepository(db)
			
			// Call the method
			ci, err := repo.GetByID(context.Background(), ciID)
			
			// Assert results
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, ci)
				assert.Equal(t, tt.expectedCI.ID, ci.ID)
				assert.Equal(t, tt.expectedCI.Name, ci.Name)
				assert.Equal(t, tt.expectedCI.Type, ci.Type)
				assert.Equal(t, tt.expectedCI.Attributes, ci.Attributes)
				assert.Equal(t, tt.expectedCI.Tags, ci.Tags)
			}
			
			// Assert that all expectations were met
			mock.AssertExpectations(t)
		})
	}
}

func TestCIPostgresRepository_GetAll(t *testing.T) {
	// Setup test data
	testCIs := []*models.CI{
		{
			ID:         uuid.New(),
			Name:       "Test CI 1",
			Type:       "server",
			Attributes: models.JSONBMap{"os": "linux", "cpu": "4"},
			Tags:       []string{"production", "web"},
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			ID:         uuid.New(),
			Name:       "Test CI 2",
			Type:       "database",
			Attributes: models.JSONBMap{"engine": "postgres", "version": "12"},
			Tags:       []string{"production"},
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	tests := []struct {
		name          string
		setupMock     func(*sqlx.DB)
		expectedCIs   []*models.CI
		expectedError bool
	}{
		{
			name: "Successful CI retrieval",
			setupMock: func(db *sqlx.DB) {
				rows := sqlmock.NewRows([]string{"id", "name", "type", "attributes", "tags", "created_at", "updated_at"}).
					AddRow(testCIs[0].ID, testCIs[0].Name, testCIs[0].Type, testCIs[0].Attributes, testCIs[0].Tags, testCIs[0].CreatedAt, testCIs[0].UpdatedAt).
					AddRow(testCIs[1].ID, testCIs[1].Name, testCIs[1].Type, testCIs[1].Attributes, testCIs[1].Tags, testCIs[1].CreatedAt, testCIs[1].UpdatedAt)
				db.ExpectQuery("SELECT id, name, type, attributes, tags, created_at, updated_at FROM configuration_items").
					WillReturnRows(rows)
			},
			expectedCIs:   testCIs,
			expectedError: false,
		},
		{
			name: "No CIs found",
			setupMock: func(db *sqlx.DB) {
				rows := sqlmock.NewRows([]string{"id", "name", "type", "attributes", "tags", "created_at", "updated_at"})
				db.ExpectQuery("SELECT id, name, type, attributes, tags, created_at, updated_at FROM configuration_items").
					WillReturnRows(rows)
			},
			expectedCIs:   []*models.CI{},
			expectedError: false,
		},
		{
			name: "Database error",
			setupMock: func(db *sqlx.DB) {
				db.ExpectQuery("SELECT id, name, type, attributes, tags, created_at, updated_at FROM configuration_items").
					WillReturnError(assert.AnError)
			},
			expectedCIs:   nil,
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
			repo := NewCIPostgresRepository(db)
			
			// Call the method
			cis, err := repo.GetAll(context.Background())
			
			// Assert results
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.expectedCIs), len(cis))
				
				if len(tt.expectedCIs) > 0 {
					for i, expectedCI := range tt.expectedCIs {
						assert.Equal(t, expectedCI.ID, cis[i].ID)
						assert.Equal(t, expectedCI.Name, cis[i].Name)
						assert.Equal(t, expectedCI.Type, cis[i].Type)
						assert.Equal(t, expectedCI.Attributes, cis[i].Attributes)
						assert.Equal(t, expectedCI.Tags, cis[i].Tags)
					}
				}
			}
			
			// Assert that all expectations were met
			mock.AssertExpectations(t)
		})
	}
}

func TestCIPostgresRepository_GetByType(t *testing.T) {
	// Setup test data
	ciType := "server"
	testCIs := []*models.CI{
		{
			ID:         uuid.New(),
			Name:       "Test CI 1",
			Type:       ciType,
			Attributes: models.JSONBMap{"os": "linux", "cpu": "4"},
			Tags:       []string{"production", "web"},
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			ID:         uuid.New(),
			Name:       "Test CI 2",
			Type:       ciType,
			Attributes: models.JSONBMap{"os": "windows", "cpu": "8"},
			Tags:       []string{"production"},
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	tests := []struct {
		name          string
		setupMock     func(*sqlx.DB)
		expectedCIs   []*models.CI
		expectedError bool
	}{
		{
			name: "Successful CI retrieval by type",
			setupMock: func(db *sqlx.DB) {
				rows := sqlmock.NewRows([]string{"id", "name", "type", "attributes", "tags", "created_at", "updated_at"}).
					AddRow(testCIs[0].ID, testCIs[0].Name, testCIs[0].Type, testCIs[0].Attributes, testCIs[0].Tags, testCIs[0].CreatedAt, testCIs[0].UpdatedAt).
					AddRow(testCIs[1].ID, testCIs[1].Name, testCIs[1].Type, testCIs[1].Attributes, testCIs[1].Tags, testCIs[1].CreatedAt, testCIs[1].UpdatedAt)
				db.ExpectQuery("SELECT id, name, type, attributes, tags, created_at, updated_at FROM configuration_items WHERE type =").
					WithArgs(ciType).
					WillReturnRows(rows)
			},
			expectedCIs:   testCIs,
			expectedError: false,
		},
		{
			name: "No CIs found for type",
			setupMock: func(db *sqlx.DB) {
				rows := sqlmock.NewRows([]string{"id", "name", "type", "attributes", "tags", "created_at", "updated_at"})
				db.ExpectQuery("SELECT id, name, type, attributes, tags, created_at, updated_at FROM configuration_items WHERE type =").
					WithArgs(ciType).
					WillReturnRows(rows)
			},
			expectedCIs:   []*models.CI{},
			expectedError: false,
		},
		{
			name: "Database error",
			setupMock: func(db *sqlx.DB) {
				db.ExpectQuery("SELECT id, name, type, attributes, tags, created_at, updated_at FROM configuration_items WHERE type =").
					WithArgs(ciType).
					WillReturnError(assert.AnError)
			},
			expectedCIs:   nil,
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
			repo := NewCIPostgresRepository(db)
			
			// Call the method
			cis, err := repo.GetByType(context.Background(), ciType)
			
			// Assert results
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.expectedCIs), len(cis))
				
				if len(tt.expectedCIs) > 0 {
					for i, expectedCI := range tt.expectedCIs {
						assert.Equal(t, expectedCI.ID, cis[i].ID)
						assert.Equal(t, expectedCI.Name, cis[i].Name)
						assert.Equal(t, expectedCI.Type, cis[i].Type)
						assert.Equal(t, expectedCI.Attributes, cis[i].Attributes)
						assert.Equal(t, expectedCI.Tags, cis[i].Tags)
					}
				}
			}
			
			// Assert that all expectations were met
			mock.AssertExpectations(t)
		})
	}
}

func TestCIPostgresRepository_Update(t *testing.T) {
	// Setup test data
	ciID := uuid.New()
	testCI := &models.CI{
		ID:         ciID,
		Name:       "Test CI",
		Type:       "server",
		Attributes: models.JSONBMap{"os": "linux", "cpu": "4"},
		Tags:       []string{"production", "web"},
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	tests := []struct {
		name          string
		setupMock     func(*sqlx.DB)
		expectedError bool
	}{
		{
			name: "Successful CI update",
			setupMock: func(db *sqlx.DB) {
				db.ExpectExec("UPDATE configuration_items").
					WithArgs(
						testCI.ID,
						testCI.Name,
						testCI.Type,
						testCI.Attributes,
						testCI.Tags,
						testCI.UpdatedAt,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: false,
		},
		{
			name: "CI not found",
			setupMock: func(db *sqlx.DB) {
				db.ExpectExec("UPDATE configuration_items").
					WithArgs(
						testCI.ID,
						testCI.Name,
						testCI.Type,
						testCI.Attributes,
						testCI.Tags,
						testCI.UpdatedAt,
					).
					WillReturnResult(sqlmock.NewResult(1, 0))
			},
			expectedError: true,
		},
		{
			name: "Database error",
			setupMock: func(db *sqlx.DB) {
				db.ExpectExec("UPDATE configuration_items").
					WithArgs(
						testCI.ID,
						testCI.Name,
						testCI.Type,
						testCI.Attributes,
						testCI.Tags,
						testCI.UpdatedAt,
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
			repo := NewCIPostgresRepository(db)
			
			// Call the method
			err := repo.Update(context.Background(), testCI)
			
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

func TestCIPostgresRepository_Delete(t *testing.T) {
	// Setup test data
	ciID := uuid.New()

	tests := []struct {
		name          string
		setupMock     func(*sqlx.DB)
		expectedError bool
	}{
		{
			name: "Successful CI deletion",
			setupMock: func(db *sqlx.DB) {
				db.ExpectExec("DELETE FROM configuration_items WHERE id =").
					WithArgs(ciID).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: false,
		},
		{
			name: "CI not found",
			setupMock: func(db *sqlx.DB) {
				db.ExpectExec("DELETE FROM configuration_items WHERE id =").
					WithArgs(ciID).
					WillReturnResult(sqlmock.NewResult(1, 0))
			},
			expectedError: true,
		},
		{
			name: "Database error",
			setupMock: func(db *sqlx.DB) {
				db.ExpectExec("DELETE FROM configuration_items WHERE id =").
					WithArgs(ciID).
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
			repo := NewCIPostgresRepository(db)
			
			// Call the method
			err := repo.Delete(context.Background(), ciID)
			
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