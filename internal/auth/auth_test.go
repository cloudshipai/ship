package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestValidateToken(t *testing.T) {
	tests := []struct {
		name          string
		token         string
		responseCode  int
		responseBody  interface{}
		expectError   bool
		expectedOrgID string
	}{
		{
			name:         "Valid token",
			token:        "valid-token",
			responseCode: http.StatusOK,
			responseBody: AuthResponse{
				OrgID:       "org-123",
				UserID:      "user-456",
				Permissions: []string{"push", "investigate"},
				Valid:       true,
			},
			expectError:   false,
			expectedOrgID: "org-123",
		},
		{
			name:         "Invalid token",
			token:        "invalid-token",
			responseCode: http.StatusUnauthorized,
			responseBody: ErrorResponse{
				Error:   "invalid_token",
				Message: "Token is invalid or expired",
			},
			expectError: true,
		},
		{
			name:         "Server error",
			token:        "any-token",
			responseCode: http.StatusInternalServerError,
			responseBody: ErrorResponse{
				Error:   "server_error",
				Message: "Internal server error",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Check method
				if r.Method != "POST" {
					t.Errorf("Expected POST method, got %s", r.Method)
				}

				// Check path
				if r.URL.Path != "/auth/validate" {
					t.Errorf("Expected /auth/validate path, got %s", r.URL.Path)
				}

				// Check auth header
				authHeader := r.Header.Get("Authorization")
				expectedHeader := "Bearer " + tt.token
				if authHeader != expectedHeader {
					t.Errorf("Expected auth header %s, got %s", expectedHeader, authHeader)
				}

				// Send response
				w.WriteHeader(tt.responseCode)
				json.NewEncoder(w).Encode(tt.responseBody)
			}))
			defer server.Close()

			// Create client
			client := NewClient(server.URL)

			// Test validate token
			resp, err := client.ValidateToken(tt.token)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if resp.OrgID != tt.expectedOrgID {
					t.Errorf("Expected org ID %s, got %s", tt.expectedOrgID, resp.OrgID)
				}
			}
		})
	}
}

func TestLogout(t *testing.T) {
	tests := []struct {
		name         string
		token        string
		responseCode int
		expectError  bool
	}{
		{
			name:         "Successful logout",
			token:        "valid-token",
			responseCode: http.StatusOK,
			expectError:  false,
		},
		{
			name:         "Successful logout with no content",
			token:        "valid-token",
			responseCode: http.StatusNoContent,
			expectError:  false,
		},
		{
			name:         "Failed logout",
			token:        "invalid-token",
			responseCode: http.StatusUnauthorized,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Check method
				if r.Method != "POST" {
					t.Errorf("Expected POST method, got %s", r.Method)
				}

				// Check path
				if r.URL.Path != "/auth/logout" {
					t.Errorf("Expected /auth/logout path, got %s", r.URL.Path)
				}

				// Send response
				w.WriteHeader(tt.responseCode)
				if tt.expectError {
					json.NewEncoder(w).Encode(ErrorResponse{
						Error:   "logout_failed",
						Message: "Logout failed",
					})
				}
			}))
			defer server.Close()

			// Create client
			client := NewClient(server.URL)

			// Test logout
			err := client.Logout(tt.token)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}
