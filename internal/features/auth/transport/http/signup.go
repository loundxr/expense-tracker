package auth_transport_http

import (
	"net/http"
	"time"

	core_http_request "github.com/loundxr/expense-tracker/internal/core/transport/http/request"
	core_http_response "github.com/loundxr/expense-tracker/internal/core/transport/http/response"
)

type SignUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpResponse struct {
	UserID    int       `json:"user_id"`
	Version   int       `json:"version"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rh := core_http_response.NewHTTPResponseHandler(w, h.logger)

	var req SignUpRequest

	if err := core_http_request.DecodeJSONRequest(r, &req); err != nil {
		rh.ErrorResponse(err, "failed to decode json request")
		return
	}

	u, err := h.authService.SignUp(ctx, req.Email, req.Password)
	if err != nil {
		rh.ErrorResponse(err, "failed to create user")
		return
	}

	resp := SignUpResponse{
		UserID:    u.ID,
		Version:   u.Version,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}
	rh.JSONResponse(resp, http.StatusCreated)
}
