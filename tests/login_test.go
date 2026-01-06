package tests

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"cmd/auth/internal/handler"
)

type mockUserService struct{}

func (m *mockUserService) ValidateCredentials(username, password string) (bool, error) {
	if username == "testuser" && password == "testpass" {
		return true, nil
	}
	return false, nil
}

func (m *mockUserService) GenerateToken(username string) (string, error) {
	return "mock-token", nil
}

func (m *mockUserService) RefreshToken(token string) (string, error) {
	return "", nil
}

func TestLoginHandler_Success(t *testing.T) {
	mockService := &mockUserService{}
	loginHandler := handler.NewLoginHandler(mockService)

	req := httptest.NewRequest("POST", "/login", nil)
	req.Header.Set("Authorization", "Basic "+basicAuth("testuser", "testpass"))

	rr := httptest.NewRecorder()
	loginHandler.Handle(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected status 200, got %v", status)
	}

	authHeader := rr.Header().Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		t.Errorf("expected Bearer token, got %v", authHeader)
	}
}

func TestLoginHandler_InvalidAuth(t *testing.T) {
	mockService := &mockUserService{}
	loginHandler := handler.NewLoginHandler(mockService)

	req := httptest.NewRequest("POST", "/login", nil)
	req.Header.Set("Authorization", "Bearer token")

	rr := httptest.NewRecorder()
	loginHandler.Handle(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("expected status 400, got %v", status)
	}
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
