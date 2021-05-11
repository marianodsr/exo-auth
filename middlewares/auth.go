package middlewares

import (
	"fmt"
	"net/http"

	"github.com/marianodsr/nura-api/authentication"
)

//RequireAuth func
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header["Authentication"]
		if authHeader == nil {
			http.Error(w, "Unauthorized", 401)
			return
		}
		tokenString := r.Header["Authentication"][0]
		if tokenString == "" {
			http.Error(w, "Unauthorized", 401)
			return
		}
		token, err := authentication.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Unauthorized", 401)
			return
		}
		fmt.Printf("%+v", token["UserID"])
		next.ServeHTTP(w, r)
	})
}
