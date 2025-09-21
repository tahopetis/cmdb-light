package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/cmdb-lite/backend/internal/models"
)

// ErrorHandler is a middleware that handles errors and returns consistent error responses
func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a response recorder to capture the status code and response
		recorder := &responseRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK, // Default status code
		}

		// Call the next handler
		next.ServeHTTP(recorder, r)

		// If the status code is an error code (4xx or 5xx), format the error response
		if recorder.statusCode >= 400 {
			// Get the error message from the response body
			var errorResponse map[string]interface{}
			if err := json.Unmarshal(recorder.data, &errorResponse); err == nil {
				// If the response is already in our error format, just write it
				if _, hasCode := errorResponse["code"]; hasCode {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(recorder.statusCode)
					w.Write(recorder.data)
					return
				}
				
				// Otherwise, convert it to our standard format
				message := "An error occurred"
				if msg, ok := errorResponse["message"].(string); ok {
					message = msg
				} else if msg, ok := errorResponse["error"].(string); ok {
					message = msg
				}
				
				// Determine error type based on status code
				var errorType models.ErrorType
				switch recorder.statusCode {
				case http.StatusBadRequest:
					errorType = models.ErrorTypeValidation
				case http.StatusUnauthorized:
					errorType = models.ErrorTypeUnauthorized
				case http.StatusForbidden:
					errorType = models.ErrorTypeForbidden
				case http.StatusNotFound:
					errorType = models.ErrorTypeNotFound
				default:
					errorType = models.ErrorTypeInternal
				}
				
				// Create standardized error response
			 standardizedError := models.NewErrorResponse(errorType, message, nil)
				
				// Write the response
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(recorder.statusCode)
				json.NewEncoder(w).Encode(standardizedError)
			} else {
				// If we can't parse the response, create a generic error
				errorType := models.ErrorTypeInternal
				if recorder.statusCode == http.StatusBadRequest {
					errorType = models.ErrorTypeValidation
				} else if recorder.statusCode == http.StatusUnauthorized {
					errorType = models.ErrorTypeUnauthorized
				} else if recorder.statusCode == http.StatusForbidden {
					errorType = models.ErrorTypeForbidden
				} else if recorder.statusCode == http.StatusNotFound {
					errorType = models.ErrorTypeNotFound
				}
				
				standardizedError := models.NewErrorResponse(errorType, "An error occurred", nil)
				
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(recorder.statusCode)
				json.NewEncoder(w).Encode(standardizedError)
			}
		} else {
			// If not an error, just write the response as-is
			w.Header().Set("Content-Type", recorder.Header().Get("Content-Type"))
			w.WriteHeader(recorder.statusCode)
			w.Write(recorder.data)
		}
	})
}

// responseRecorder is a wrapper around http.ResponseWriter that captures the status code and response data
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
	data       []byte
}

// WriteHeader captures the status code
func (r *responseRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

// Write captures the response data
func (r *responseRecorder) Write(data []byte) (int, error) {
	r.data = append(r.data, data...)
	return r.ResponseWriter.Write(data)
}

// RespondWithError is a helper function to send a standardized error response
func RespondWithError(w http.ResponseWriter, errorType models.ErrorType, message string, details interface{}) {
	errorResponse := models.NewErrorResponse(errorType, message, details)
	statusCode := models.GetHTTPStatusForError(errorType)
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorResponse)
}

// RespondWithValidationError is a helper function for validation errors
func RespondWithValidationError(w http.ResponseWriter, message string, details interface{}) {
	RespondWithError(w, models.ErrorTypeValidation, message, details)
}

// RespondWithUnauthorizedError is a helper function for unauthorized errors
func RespondWithUnauthorizedError(w http.ResponseWriter, message string, details interface{}) {
	RespondWithError(w, models.ErrorTypeUnauthorized, message, details)
}

// RespondWithForbiddenError is a helper function for forbidden errors
func RespondWithForbiddenError(w http.ResponseWriter, message string, details interface{}) {
	RespondWithError(w, models.ErrorTypeForbidden, message, details)
}

// RespondWithNotFoundError is a helper function for not found errors
func RespondWithNotFoundError(w http.ResponseWriter, message string, details interface{}) {
	RespondWithError(w, models.ErrorTypeNotFound, message, details)
}

// RespondWithInternalError is a helper function for internal server errors
func RespondWithInternalError(w http.ResponseWriter, message string, details interface{}) {
	RespondWithError(w, models.ErrorTypeInternal, message, details)
}