package tests

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"cmd/auth/internal/handler"
)

type mockVerifyService struct {
	shouldFail bool
}

func (m *mockVerifyService) ValidateCredentials(username, password string) (bool, error) {
	return false, nil
}

func (m *mockVerifyService) GenerateToken(username string) (string, error) {
	return "", nil
}

func (m *mockVerifyService) RefreshToken(token string) (string, error) {
	if m.shouldFail {
		return "", errors.New("token expired")
	}
	return "new-token", nil
}

func TestVerifyHandler_Success(t *testing.T) {
	mockService := &mockVerifyService{shouldFail: false}
	verifyHandler := handler.NewVerifyHandler(mockService)

	req := httptest.NewRequest("POST", "/verify", nil)
	req.Header.Set("Authorization", "Bearer valid-token")

	rr := httptest.NewRecorder()
	verifyHandler.Handle(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected status 200, got %v", status)
	}
}

func TestVerifyHandler_InvalidToken(t *testing.T) {
	mockService := &mockVerifyService{shouldFail: true}
	verifyHandler := handler.NewVerifyHandler(mockService)

	req := httptest.NewRequest("POST", "/verify", nil)
	req.Header.Set("Authorization", "Bearer expired-token")

	rr := httptest.NewRecorder()
	verifyHandler.Handle(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %v", status)
	}
}

func TestVerifyHandler_NoAuthHeader(t *testing.T) {
	mockService := &mockVerifyService{shouldFail: false}
	verifyHandler := handler.NewVerifyHandler(mockService)

	req := httptest.NewRequest("POST", "/verify", nil)

	rr := httptest.NewRecorder()
	verifyHandler.Handle(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %v", status)
	}
}
