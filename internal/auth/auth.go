package auth

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
)

// BearerAuth ...
func BearerAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var header = r.Header.Get("Authorization")

		header = strings.TrimSpace(header)

		w.Header().Add("Content-Type", "application/json")
		if header != os.Getenv("AUTH_TOKEN") {
			w.WriteHeader(http.StatusForbidden)
			res, _ := json.Marshal(map[string]string{"error": "Invalid token"})
			w.Write(res)
			return
		}
		next.ServeHTTP(w, r)
	})
}
