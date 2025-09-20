package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPasswordManager_HashPassword(t *testing.T) {
	// Setup test data
	password := "test-password-123"

	tests := []struct {
		name          string
		setupPassword func() string
		expectedError bool
	}{
		{
			name: "Successful password hashing",
			setupPassword: func() string {
				return password
			},
			expectedError: false,
		},
		{
			name: "Empty password",
			setupPassword: func() string {
				return ""
			},
			expectedError: false, // bcrypt can hash empty passwords
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup password manager
			passwordManager := NewPasswordManager()
			
			// Setup password
			password := tt.setupPassword()
			
			// Call the method
			hashedPassword, err := passwordManager.HashPassword(password)
			
			// Assert results
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, hashedPassword)
				assert.NotEqual(t, password, hashedPassword) // Hash should be different from original password
				
				// Verify that the hash starts with bcrypt's prefix
				assert.Contains(t, hashedPassword, "$2a$12$")
			}
		})
	}
}

func TestPasswordManager_CheckPassword(t *testing.T) {
	// Setup test data
	password := "test-password-123"
	wrongPassword := "wrong-password"
	
	// Create a password manager and hash the password
	passwordManager := NewPasswordManager()
	hashedPassword, err := passwordManager.HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	tests := []struct {
		name          string
		setupPassword func() (string, string)
		expectedError bool
		errorType     error
	}{
		{
			name: "Successful password verification",
			setupPassword: func() (string, string) {
				return password, hashedPassword
			},
			expectedError: false,
		},
		{
			name: "Incorrect password",
			setupPassword: func() (string, string) {
				return wrongPassword, hashedPassword
			},
			expectedError: true,
			errorType:     ErrInvalidPassword,
		},
		{
			name: "Empty password",
			setupPassword: func() (string, string) {
				return "", hashedPassword
			},
			expectedError: true,
			errorType:     ErrInvalidPassword,
		},
		{
			name: "Invalid hash",
			setupPassword: func() (string, string) {
				return password, "invalid-hash"
			},
			expectedError: true,
			errorType:     ErrInvalidPassword,
		},
		{
			name: "Empty hash",
			setupPassword: func() (string, string) {
				return password, ""
			},
			expectedError: true,
			errorType:     ErrInvalidPassword,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup password manager
			passwordManager := NewPasswordManager()
			
			// Setup password and hash
			password, hash := tt.setupPassword()
			
			// Call the method
			err := passwordManager.CheckPassword(password, hash)
			
			// Assert results
			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorType != nil {
					assert.ErrorIs(t, err, tt.errorType)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPasswordManager_HashAndCheck(t *testing.T) {
	// This test verifies that a password can be hashed and then verified
	
	// Setup test data
	passwords := []string{
		"simple",
		"more-complex-password123",
		"very!@#$%^&*()complex+{}|:\"<>?password",
		"password-with-unicode-测试",
		"", // Empty password
	}
	
	for _, password := range passwords {
		t.Run("Password: "+password, func(t *testing.T) {
			// Setup password manager
			passwordManager := NewPasswordManager()
			
			// Hash the password
			hashedPassword, err := passwordManager.HashPassword(password)
			require.NoError(t, err)
			require.NotEmpty(t, hashedPassword)
			
			// Verify the password
			err = passwordManager.CheckPassword(password, hashedPassword)
			assert.NoError(t, err)
			
			// Verify that a different password fails
			err = passwordManager.CheckPassword(password+"-wrong", hashedPassword)
			assert.Error(t, err)
			assert.ErrorIs(t, err, ErrInvalidPassword)
		})
	}
}

func TestPasswordManager_DifferentHashes(t *testing.T) {
	// This test verifies that hashing the same password multiple times produces different hashes
	
	// Setup test data
	password := "test-password-123"
	
	// Create a password manager
	passwordManager := NewPasswordManager()
	
	// Hash the password multiple times
	hash1, err := passwordManager.HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hash1)
	
	hash2, err := passwordManager.HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hash2)
	
	hash3, err := passwordManager.HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hash3)
	
	// Verify that all hashes are different (due to salt)
	assert.NotEqual(t, hash1, hash2)
	assert.NotEqual(t, hash2, hash3)
	assert.NotEqual(t, hash1, hash3)
	
	// Verify that all hashes can be verified with the original password
	err = passwordManager.CheckPassword(password, hash1)
	assert.NoError(t, err)
	
	err = passwordManager.CheckPassword(password, hash2)
	assert.NoError(t, err)
	
	err = passwordManager.CheckPassword(password, hash3)
	assert.NoError(t, err)
}