package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Firuz43/ecommerce/internal/auth"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Get the "Authorization" header
		// Format: "Bearer <token>"
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// 2. Remove the "Bearer " prefix
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		// 3. Validate the token
		claims, err := auth.ValidateToken(tokenStr)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// 4. (Optional) Inject UserID into the request context
		// so the next function knows WHO is making the request
		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
