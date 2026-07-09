package auth

import (
	"context"
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v3"
	"github.com/auth0/go-jwt-middleware/v3/validator"
)

type contextKey string

const orgIDKey contextKey = "orgID"

type UserDetail struct {
	UserID string
	OrgID  string
}

func UserDetailsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := jwtmiddleware.GetClaims[*validator.ValidatedClaims](r.Context())
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"message":"Unauthorized."}`))
			return
		}

		customClaims, ok := claims.CustomClaims.(*CustomClaims)
		if !ok {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`{"message":"Issue with getting org_id."}`))
			return
		}
		userDetails := UserDetail{UserID: customClaims.UserID, OrgID: customClaims.OrgID}
		ctx := context.WithValue(r.Context(), orgIDKey, userDetails)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserDetails(ctx context.Context) UserDetail {
	userDetail, _ := ctx.Value(orgIDKey).(UserDetail)
	return userDetail
}
