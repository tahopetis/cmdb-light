package repositories

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/cmdb-lite/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockDB is a mock implementation of sqlx.DB for testing
type MockDB struct {
	mock.Expecter
}

func TestUserPostgresRepository_Create(t *testing.T) {
	// Setup test data
	userID := uuid.New()
	testUser := &models.User{
		ID:           userID,
		Username:     "testuser",
		PasswordHash: "hashedpassword",
		Role:         "user",
		CreatedAt:    time.Now(),
	}

	tests := []struct {
		name          string
		setupMock     func(*sqlx.DB)
		expectedError bool
	}{
		{
			name: "Successful user creation",
			setupMock: func(db *sqlx.DB) {
				db.ExpectExec("INSERT INTO users").
					WithArgs(
						testUser.ID,
						testUser.Username,
						testUser.PasswordHash,
						testUser.Role,
						testUser.CreatedAt,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: false,
		},
		{
			name: "Database error",
			setupMock: func(db *sqlx.DB) {
				db.ExpectExec("INSERT INTO users").
					WithArgs(
						testUser.ID,
						testUser.Username,
						testUser.PasswordHash,
						testUser.Role,
						testUser.CreatedAt,
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
			repo := NewUserPostgresRepository(db)
			
			// Call the method
			err := repo.Create(context.Background(), testUser)
			
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

func TestUserPostgresRepository_GetByID(t *testing.T) {
	// Setup test data
	userID := uuid.New()
	testUser := &models.User{
		ID:           userID,
		Username:     "testuser",
		PasswordHash: "hashedpassword",
		Role:         "user",
		CreatedAt:    time.Now(),
	}

	tests := []struct {
		name          string
		setupMock     func(*sqlx.DB)
		expectedUser  *models.User
		expectedError bool
		expectedErrorType error
	}{
		{
			name: "Successful user retrieval",
			setupMock: func(db *sqlx.DB) {
				rows := sqlmock.NewRows([]string{"id", "username", "password_hash", "role", "created_at"}).
					AddRow(testUser.ID, testUser.Username, testUser.PasswordHash, testUser.Role, testUser.CreatedAt)
				db.ExpectQuery("SELECT id, username, password_hash, role, created_at FROM users WHERE id =").
					WithArgs(userID).
					WillReturnRows(rows)
			},
			expectedUser:  testUser,
			expectedError: false,
		},
		{
			name: "User not found",
			setupMock: func(db *sqlx.DB) {
				db.ExpectQuery("SELECT id, username, password_hash, role, created_at FROM users WHERE id =").
					WithArgs(userID).
					WillReturnError(sql.ErrNoRows)
			},
			expectedUser:  nil,
			expectedError: true,
			expectedErrorType: ErrUserNotFound,
		},
		{
			name: "Database error",
			setupMock: func(db *sqlx.DB) {
				db.ExpectQuery("SELECT id, username, password_hash, role, created_at FROM users WHERE id =").
					WithArgs(userID).
					WillReturnError(assert.AnError)
			},
			expectedUser:  nil,
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
			repo := NewUserPostgresRepository(db)
			
			// Call the method
			user, err := repo.GetByID(context.Background(), userID)
			
			// Assert results
			if tt.expectedError {
				assert.Error(t, err)
				if tt.expectedErrorType != nil {
					assert.ErrorIs(t, err, tt.expectedErrorType)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser.ID, user.ID)
				assert.Equal(t, tt.expectedUser.Username, user.Username)
				assert.Equal(t, tt.expectedUser.PasswordHash, user.PasswordHash)
				assert.Equal(t, tt.expectedUser.Role, user.Role)
			}
			
			// Assert that all expectations were met
			mock.AssertExpectations(t)
		})
	}
}

func TestUserPostgresRepository_GetByUsername(t *testing.T) {
	// Setup test data
	username := "testuser"
	testUser := &models.User{
		ID:           uuid.New(),
		Username:     username,
		PasswordHash: "hashedpassword",
		Role:         "user",
		CreatedAt:    time.Now(),
	}

	tests := []struct {
		name          string
		setupMock     func(*sqlx.DB)
		expectedUser  *models.User
		expectedError bool
		expectedErrorType error
	}{
		{
			name: "Successful user retrieval",
			setupMock: func(db *sqlx.DB) {
				rows := sqlmock.NewRows([]string{"id", "username", "password_hash", "role", "created_at"}).
					AddRow(testUser.ID, testUser.Username, testUser.PasswordHash, testUser.Role, testUser.CreatedAt)
				db.ExpectQuery("SELECT id, username, password_hash, role, created_at FROM users WHERE username =").
					WithArgs(username).
					WillReturnRows(rows)
			},
			expectedUser:  testUser,
			expectedError: false,
		},
		{
			name: "User not found",
			setupMock: func(db *sqlx.DB) {
				db.ExpectQuery("SELECT id, username, password_hash, role, created_at FROM users WHERE username =").
					WithArgs(username).
					WillReturnError(sql.ErrNoRows)
			},
			expectedUser:  nil,
			expectedError: true,
			expectedErrorType: ErrUserNotFound,
		},
		{
			name: "Database error",
			setupMock: func(db *sqlx.DB) {
				db.ExpectQuery("SELECT id, username, password_hash, role, created_at FROM users WHERE username =").
					WithArgs(username).
					WillReturnError(assert.AnError)
			},
			expectedUser:  nil,
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
			repo := NewUserPostgresRepository(db)
			
			// Call the method
			user, err := repo.GetByUsername(context.Background(), username)
			
			// Assert results
			if tt.expectedError {
				assert.Error(t, err)
				if tt.expectedErrorType != nil {
					assert.ErrorIs(t, err, tt.expectedErrorType)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser.ID, user.ID)
				assert.Equal(t, tt.expectedUser.Username, user.Username)
				assert.Equal(t, tt.expectedUser.PasswordHash, user.PasswordHash)
				assert.Equal(t, tt.expectedUser.Role, user.Role)
			}
			
			// Assert that all expectations were met
			mock.AssertExpectations(t)
		})
	}
}

func TestUserPostgresRepository_Update(t *testing.T) {
	// Setup test data
	userID := uuid.New()
	testUser := &models.User{
		ID:           userID,
		Username:     "testuser",
		PasswordHash: "hashedpassword",
		Role:         "user",
		CreatedAt:    time.Now(),
	}

	tests := []struct {
		name          string
		setupMock     func(*sqlx.DB)
		expectedError bool
		expectedErrorType error
	}{
		{
			name: "Successful user update",
			setupMock: func(db *sqlx.DB) {
				db.ExpectExec("UPDATE users").
					WithArgs(
						testUser.ID,
						testUser.Username,
						testUser.PasswordHash,
						testUser.Role,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: false,
		},
		{
			name: "User not found",
			setupMock: func(db *sqlx.DB) {
				db.ExpectExec("UPDATE users").
					WithArgs(
						testUser.ID,
						testUser.Username,
						testUser.PasswordHash,
						testUser.Role,
					).
					WillReturnResult(sqlmock.NewResult(1, 0))
			},
			expectedError: true,
			expectedErrorType: ErrUserNotFound,
		},
		{
			name: "Database error",
			setupMock: func(db *sqlx.DB) {
				db.ExpectExec("UPDATE users").
					WithArgs(
						testUser.ID,
						testUser.Username,
						testUser.PasswordHash,
						testUser.Role,
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
			repo := NewUserPostgresRepository(db)
			
			// Call the method
			err := repo.Update(context.Background(), testUser)
			
			// Assert results
			if tt.expectedError {
				assert.Error(t, err)
				if tt.expectedErrorType != nil {
					assert.ErrorIs(t, err, tt.expectedErrorType)
				}
			} else {
				assert.NoError(t, err)
			}
			
			// Assert that all expectations were met
			mock.AssertExpectations(t)
		})
	}
}

func TestUserPostgresRepository_Delete(t *testing.T) {
	// Setup test data
	userID := uuid.New()

	tests := []struct {
		name          string
		setupMock     func(*sqlx.DB)
		expectedError bool
		expectedErrorType error
	}{
		{
			name: "Successful user deletion",
			setupMock: func(db *sqlx.DB) {
				db.ExpectExec("DELETE FROM users WHERE id =").
					WithArgs(userID).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: false,
		},
		{
			name: "User not found",
			setupMock: func(db *sqlx.DB) {
				db.ExpectExec("DELETE FROM users WHERE id =").
					WithArgs(userID).
					WillReturnResult(sqlmock.NewResult(1, 0))
			},
			expectedError: true,
			expectedErrorType: ErrUserNotFound,
		},
		{
			name: "Database error",
			setupMock: func(db *sqlx.DB) {
				db.ExpectExec("DELETE FROM users WHERE id =").
					WithArgs(userID).
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
			repo := NewUserPostgresRepository(db)
			
			// Call the method
			err := repo.Delete(context.Background(), userID)
			
			// Assert results
			if tt.expectedError {
				assert.Error(t, err)
				if tt.expectedErrorType != nil {
					assert.ErrorIs(t, err, tt.expectedErrorType)
				}
			} else {
				assert.NoError(t, err)
			}
			
			// Assert that all expectations were met
			mock.AssertExpectations(t)
		})
	}
}