package gateways

import (
	"export-service/internal/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHTTPAuthService_Login_Responses(t *testing.T) {
	tests := []struct {
		name          string
		mockResponse  server.HttpResponse
		expectError   bool
		expectedToken string
	}{
		{
			name: "Should handle login",
			mockResponse: server.HttpResponse{
				StatusCode: 201,
				Body: map[string]any{
					"token": "token",
				},
			},
			expectError:   false,
			expectedToken: "token",
		},
		{
			name: "Should handle login error",
			mockResponse: server.HttpResponse{
				StatusCode: 400,
			},
			expectError: true,
		},
		{
			name: "Should handle wrong response",
			mockResponse: server.HttpResponse{
				StatusCode: 201,
				Body:       []map[string]any{},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, m := getHTTPAuthService()

			m.Expect("POST", "/authenticate", tt.mockResponse)

			response, err := a.Login("email", "password", "company")

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedToken, response["token"])
			}
		})
	}
}

func TestHTTPAuthService_Login_Logic(t *testing.T) {
	a, m := getHTTPAuthService()

	m.Expect("POST", "/authenticate", server.HttpResponse{
		StatusCode: 201,
		Body: map[string]any{
			"token": "token",
		},
	})

	response, err := a.Login("email", "password", "company")

	assert.NoError(t, err)
	assert.Equal(t, "token", response["token"])

	request, ok := m.VerifyRequest("POST", "/authenticate")
	require.True(t, ok)
	assert.Contains(t, request.URL, "/authenticate?company=company")
	expectedBody := map[string]any{
		"email":    "email",
		"password": "password",
	}
	assert.Equal(t, expectedBody, request.Body)
}

func TestHTTPAuthService_GetUserByToken_Responses(t *testing.T) {
	tests := []struct {
		name         string
		mockResponse server.HttpResponse
		expectError  bool
		expectedUser AuthUser
	}{
		{
			name: "Should get user by token",
			mockResponse: server.HttpResponse{
				StatusCode: 200,
				Body: map[string]any{
					"id":          "123",
					"workspaceId": "456",
				},
			},
			expectError: false,
			expectedUser: AuthUser{
				ID:          "123",
				WorkspaceID: "456",
			},
		},
		{
			name: "Should handle wrong status",
			mockResponse: server.HttpResponse{
				StatusCode: 401,
				Body:       map[string]any{},
			},
			expectError: true,
		},
		{
			name: "Should handle wrong response",
			mockResponse: server.HttpResponse{
				StatusCode: 200,
				Body: map[string]any{
					"user": map[string]any{},
				},
			},
			expectError: true,
		},
		{
			name: "Should handle wrong type of response",
			mockResponse: server.HttpResponse{
				StatusCode: 200,
				Body:       []map[string]any{},
			},
			expectError: true,
		},
		{
			name: "Should handle wrong user field type",
			mockResponse: server.HttpResponse{
				StatusCode: 200,
				Body: map[string]any{
					"id": 123,
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, m := getHTTPAuthService()

			m.Expect("GET", "/current-user", tt.mockResponse)

			user, err := a.GetUserByToken(map[string]any{})

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser, user)
			}
		})
	}
}

func TestHTTPAuthService_GetUserByToken_Logic(t *testing.T) {
	a, m := getHTTPAuthService()

	m.Expect("GET", "/current-user", server.HttpResponse{
		StatusCode: 200,
		Body: map[string]any{
			"id":          "123",
			"workspaceId": "456",
		},
	})

	expectedUser := AuthUser{
		ID:          "123",
		WorkspaceID: "456",
	}
	user, err := a.GetUserByToken(map[string]any{"Authorization": "token"})
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)

	request, ok := m.VerifyRequest("GET", "/current-user")
	require.True(t, ok)
	assert.Contains(t, request.Headers, "Authorization")
	assert.Equal(t, request.Headers["Authorization"], "token")
}

func getHTTPAuthService() (*HTTPAuthService, *server.MockHttpClient) {
	m := server.NewMockHttpClient()
	a := &HTTPAuthService{HttpClient: m}
	return a, m
}
