package middleware

import (
	"food-court/utils"
	"net/http"
	"strings"
	"strconv"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.RespondWithError(w, http.StatusUnauthorized, "No authorization header")
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 {
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid token format")
			return
		}

		userID, err := utils.ValidateJWT(bearerToken[1])
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		// Add user ID to request context
		r = r.WithContext(utils.ContextWithUserID(r.Context(), strconv.FormatUint(uint64(userID), 10)))
		next.ServeHTTP(w, r)
	})
}