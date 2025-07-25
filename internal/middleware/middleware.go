package middleware

import (
	"context"
	"github.com/KhalefAhmed/fem-server/internal/store"
	tokens "github.com/KhalefAhmed/fem-server/internal/token"
	"github.com/KhalefAhmed/fem-server/internal/utils"
	"net/http"
	"strings"
)

type contextKey string

const UserContextKey = contextKey("user")

type UserMiddleware struct {
	UserStore store.UserStore
}

func SetUser(r *http.Request, user *store.User) *http.Request {
	ctx := context.WithValue(r.Context(), UserContextKey, user)
	return r.WithContext(ctx)
}

func GetUser(r *http.Request) *store.User {
	user, ok := r.Context().Value(UserContextKey).(*store.User)
	if !ok {
		panic("missing user in request")
	}
	return user
}

func (um *UserMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Add("Vary", "Authorization")
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			r = SetUser(r, store.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			utils.WriteJson(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid authorization header"})
			return
		}

		token := headerParts[1]
		user, err := um.UserStore.GetUserToken(tokens.ScopeAuth, token)
		if err != nil {
			utils.WriteJson(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid token"})
			return
		}

		if user == nil {
			utils.WriteJson(w, http.StatusUnauthorized, utils.Envelope{"error": "token expired or invalid"})
			return
		}

		r = SetUser(r, user)
		next.ServeHTTP(w, r)
		return
	})
}

func (um *UserMiddleware) RequireUser(next http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		user := GetUser(request)
		if user.IsAnonymous() {
			utils.WriteJson(writer, http.StatusUnauthorized, utils.Envelope{"error": "you must be logged in to access this route"})
			return
		}

		next.ServeHTTP(writer, request)
	}
}
