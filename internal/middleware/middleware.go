package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/syafae/femProject/internal/store"
	"github.com/syafae/femProject/internal/tokens"
	"github.com/syafae/femProject/internal/utils"
)

type UserMiddleware struct {
	UserStore store.UserStore
}
type contextKey string

const UserContextKey = contextKey("user")

func SetUser(r *http.Request, user *store.User) *http.Request {
	ctx := context.WithValue(r.Context(), UserContextKey, user)
	return r.WithContext(ctx)
}

func GetUser(r *http.Request) *store.User {
	user, ok := r.Context().Value((UserContextKey)).(*store.User)
	if !ok {
		panic("missing user in request")
	}
	return user
}

func (um *UserMiddleware) Authenicate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			r = SetUser(r, store.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || strings.ToLower(headerParts[0]) != "bearer" {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid authorization header"})
			return
		}
		token := headerParts[1]
		user, err := um.UserStore.GetUserToken(tokens.ScopeAuth, token)
		if err != nil {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid token"})
			return
		}
		if user == nil {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid or expired token"})
			return
		}
		r = SetUser(r, user)
		next.ServeHTTP(w, r)

	})
}

func (um *UserMiddleware) RequireUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUser(r)
		if user.IsAnonymous() {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "you must be logged in to access this resource"})
			return
		}
		next.ServeHTTP(w, r)
	})
}
