package handler

import (
	"encoding/base64"
	"net/http"
	"strings"

	"cmd/auth/internal/service"
)

type LoginHandler struct {
	userService service.UserService
}

func NewLoginHandler(userService service.UserService) *LoginHandler {
	return &LoginHandler{
		userService: userService,
	}
}

func (h *LoginHandler) Handle(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		http.Error(w, "Authorization header required", http.StatusUnauthorized)
		return
	}

	if !strings.HasPrefix(authHeader, "Basic ") {
		http.Error(w, "Only Basic authentication supported", http.StatusBadRequest)
		return
	}

	encodedCredentials := strings.TrimPrefix(authHeader, "Basic ")
	decoded, err := base64.StdEncoding.DecodeString(encodedCredentials)
	if err != nil {
		http.Error(w, "Invalid credentials encoding", http.StatusBadRequest)
		return
	}

	credentials := string(decoded)
	colonIndex := strings.Index(credentials, ":")
	if colonIndex == -1 {
		http.Error(w, "Invalid credentials format", http.StatusBadRequest)
		return
	}

	username := credentials[:colonIndex]
	password := credentials[colonIndex+1:]

	valid, err := h.userService.ValidateCredentials(username, password)
	if err != nil || !valid {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	token, err := h.userService.GenerateToken(username)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Authorization", "Bearer "+token)
	w.WriteHeader(http.StatusOK)

	response := `{"message":"Login successful","token_type":"Bearer"}`
	w.Write([]byte(response))
}
