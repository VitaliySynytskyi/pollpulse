package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// UserClaims defines custom claims for the JWT
type UserClaims struct {
	UserID   string   `json:"user_id"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Roles    []string `json:"roles"`
	jwt.RegisteredClaims
}

// Auth middleware for validating JWT tokens
func Auth(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header is required", http.StatusUnauthorized)
				return
			}

			// Bearer token format check
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Authorization header format must be Bearer {token}", http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]
			claims := &UserClaims{}

			// Parse the JWT token
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				// Validate the algorithm
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(jwtSecret), nil
			})

			if err != nil {
				if errors.Is(err, jwt.ErrTokenExpired) {
					http.Error(w, "Token expired", http.StatusUnauthorized)
					return
				}
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			if !token.Valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// Create a context with user claims
			ctx := context.WithValue(r.Context(), "user", claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GenerateJWT creates a new JWT token for a user
func GenerateJWT(userID, username, email string, roles []string, jwtSecret string, expirationTime time.Duration) (string, error) {
	claims := UserClaims{
		UserID:   userID,
		Username: username,
		Email:    email,
		Roles:    roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expirationTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "pollpulse",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

// GetUserFromContext extracts the user information from the context
func GetUserFromContext(ctx context.Context) (*UserClaims, error) {
	user, ok := ctx.Value("user").(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("no user found in context")
	}
	return user, nil
}

// CheckRole verifies if the user has the required role
func CheckRole(ctx context.Context, requiredRole string) bool {
	user, err := GetUserFromContext(ctx)
	if err != nil {
		return false
	}

	for _, role := range user.Roles {
		if role == requiredRole {
			return true
		}
	}
	return false
}
