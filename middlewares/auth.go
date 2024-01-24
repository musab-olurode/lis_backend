package middlewares

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/musab-olurode/lis_backend/constants"
	"github.com/musab-olurode/lis_backend/database"
	"github.com/musab-olurode/lis_backend/utils"
)

func Authenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearerToken, err := utils.GetBearerToken(r.Header)
		if err != nil {
			utils.RespondWithErr(w, http.StatusUnauthorized, err.Error())
			return
		}

		jwtClaims, err := utils.GetJWTClaims(bearerToken)
		if err != nil {
			utils.RespondWithErr(w, http.StatusUnauthorized, "invalid authorization")
			return
		}

		claims := *jwtClaims

		userID, err := uuid.Parse(claims["jti"].(string))
		if err != nil {
			utils.RespondWithErr(w, http.StatusUnauthorized, "invalid authorization")
			return
		}

		user, err := database.DB.GetUserByID(r.Context(), userID)
		if err != nil {
			utils.RespondWithErr(w, http.StatusUnauthorized, "invalid authorization")
			return
		}

		ctx := context.WithValue(r.Context(), constants.USER_CONTEXT_KEY, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func Admin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(constants.USER_CONTEXT_KEY).(database.User)
		if user.Role != database.UserRoleADMIN {
			utils.RespondWithErr(w, http.StatusForbidden, "forbidden")
			return
		}

		next.ServeHTTP(w, r)
	})
}
