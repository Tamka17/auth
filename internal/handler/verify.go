package handler

import (
	"net/http"
	"strings"

	"cmd/auth/internal/service"
)

type VerifyHandler struct {
	userService service.UserService
}

func NewVerifyHandler(userService service.UserService) *VerifyHandler {
	return &VerifyHandler{
		userService: userService,
	}
}

func (h *VerifyHandler) Handle(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header required", http.StatusUnauthorized)
		return
	}

	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		http.Error(w, "Authorization scheme must be Bearer", http.StatusBadRequest)
		return
	}

	token := strings.TrimSpace(strings.TrimPrefix(authHeader, bearerPrefix))
	if token == "" {
		http.Error(w, "Token cannot be empty", http.StatusBadRequest)
		return
	}

	newToken, err := h.userService.RefreshToken(token)
	if err != nil {
		if strings.Contains(err.Error(), "expired") ||
			strings.Contains(err.Error(), "invalid") {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Authorization", "Bearer "+newToken)
	w.WriteHeader(http.StatusOK)

	response := `{"message":"Token refreshed successfully","token_type":"Bearer"}`
	w.Write([]byte(response))
}
