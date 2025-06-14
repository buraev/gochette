package auth

import (
	"fmt"
	"lightweight-cache-proxy-service/internal/apis/secrets"
	"net/http"
	"strings"
)

func IsAuthorized(w http.ResponseWriter, r *http.Request) bool {
	validTokens := strings.Fields(secrets.ENV.ValidTokens)

	givenToken := r.Header.Get("Authorization")
	authorized := false
	for _, token := range validTokens {
		if givenToken == fmt.Sprintf("Bearer %s", token) {
			authorized = true
			break
		}
	}

	if !authorized {
		http.Error(w, "Invalid bearer auth token", http.StatusUnauthorized)
	}
	return authorized
}
