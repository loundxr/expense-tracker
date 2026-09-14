package core_http_middleware

import (
	"log/slog"
	"net/http"
	"strings"

	core_jwt "github.com/loundxr/expense-tracker/internal/core/auth/jwt"
	"github.com/loundxr/expense-tracker/internal/core/domain"
	core_errors "github.com/loundxr/expense-tracker/internal/core/errors"
	core_ctx "github.com/loundxr/expense-tracker/internal/core/transport/http/context"
	core_http_response "github.com/loundxr/expense-tracker/internal/core/transport/http/response"
)

func Auth(secret string, logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rh := core_http_response.NewHTTPResponseHandler(w, logger)

			header := r.Header.Get("Authorization")
			if header == "" {
				rh.ErrorResponse(core_errors.ErrUnauthorized, "'Authorization' header missing")
				return
			}

			parts := strings.Split(header, " ")
			if parts[0] != "Bearer" || len(parts) != 2 {
				rh.ErrorResponse(core_errors.ErrUnauthorized, "invalid 'Authorization' header format")
				return
			}

			userClaims, err := core_jwt.ParseToken(parts[1], secret)
			if err != nil {
				rh.ErrorResponse(core_errors.ErrUnauthorized, "invalid token")
				return
			}

			ctx := r.Context()
			ctx = core_ctx.SetUserID(ctx, userClaims.ID)
			ctx = core_ctx.SetUserRole(ctx, userClaims.Role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AdminOnly(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role := core_ctx.GetUserRole(r.Context())
			uid := core_ctx.GetUserID(r.Context())

			if role != domain.RoleAdmin {
				logger.Warn(
					"user unsuccessfully tried to access users list",
					slog.Int64("id", uid),
				)
				rh := core_http_response.NewHTTPResponseHandler(w, logger)
				rh.ErrorResponse(core_errors.ErrForbidden, "admin access required")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
