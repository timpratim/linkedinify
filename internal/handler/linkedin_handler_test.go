package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/you/linkedinify/internal/handler"
	"github.com/you/linkedinify/internal/service"
)

// Helper to generate a JWT for testing
func generateTestToken(t *testing.T, userID uuid.UUID, secret []byte) string {
	t.Helper()
	claims := jwt.MapClaims{
		"sub": userID.String(),
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(secret)
	require.NoError(t, err)
	return signedToken
}

func TestLinkedInHandler_transform_Success(t *testing.T) {
	mockService := &service.LinkedInServiceInteractorMock{
		TransformFunc: func(ctx context.Context, userID uuid.UUID, text string) (string, error) {
			assert.Equal(t, "00000000-0000-0000-0000-000000000001", userID.String())
			assert.Equal(t, "some input text", text)
			return "transformed linkedin post", nil
		},
	}
	testUserID, _ := uuid.Parse("00000000-0000-0000-0000-000000000001")
	testSecret := []byte("your-test-jwt-secret")
	linkedinHandler := handler.NewLinkedIn(mockService)
	router := linkedinHandler.Routes(testSecret)
	server := httptest.NewServer(router)
	defer server.Close()

	requestBody := map[string]string{"text": "some input text"}
	jsonBody, err := json.Marshal(requestBody)
	require.NoError(t, err)
	authToken := generateTestToken(t, testUserID, testSecret)

	req, err := http.NewRequest(http.MethodPost, server.URL+"/", bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+authToken)

	resp, err := server.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var responseBody map[string]string
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	require.NoError(t, err)
	assert.Equal(t, "transformed linkedin post", responseBody["post"])
	assert.Len(t, mockService.TransformCalls(), 1)
	call := mockService.TransformCalls()[0]
	assert.Equal(t, testUserID, call.UserID)
	assert.Equal(t, "some input text", call.Text)
}

func TestLinkedInHandler_transform_BadRequest_EmptyText(t *testing.T) {
	mockService := &service.LinkedInServiceInteractorMock{}
	testUserID, _ := uuid.Parse("00000000-0000-0000-0000-000000000002")
	testSecret := []byte("your-test-jwt-secret")
	linkedinHandler := handler.NewLinkedIn(mockService)
	router := linkedinHandler.Routes(testSecret)
	server := httptest.NewServer(router)
	defer server.Close()

	requestBody := map[string]string{"text": ""}
	jsonBody, _ := json.Marshal(requestBody)
	authToken := generateTestToken(t, testUserID, testSecret)

	req, _ := http.NewRequest(http.MethodPost, server.URL+"/", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+authToken)

	resp, _ := server.Client().Do(req)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Len(t, mockService.TransformCalls(), 0)
}
