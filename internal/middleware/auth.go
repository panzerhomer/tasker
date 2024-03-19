package middleware

import (
	"context"
	"net/http"

	myjwt "rest/internal/lib/jwt"
)

func Auth(handler http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("jwt")
		switch err {
		case nil:
		case http.ErrNoCookie:
			w.WriteHeader(http.StatusUnauthorized)
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		user, ok := myjwt.ValidateToken(c.Value)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if user == nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		idCtx := context.WithValue(r.Context(), "userID", user.ID)

		handler.ServeHTTP(w, r.WithContext(idCtx))
	}

	return http.HandlerFunc(fn)
}
