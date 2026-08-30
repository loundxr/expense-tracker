package auth_transport_http

import (
	"net/http"

	core_http_request "github.com/loundxr/expense-tracker/internal/core/transport/http/request"
	core_http_response "github.com/loundxr/expense-tracker/internal/core/transport/http/response"
)

type SignInRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=64"`
}

type SignInResponse struct {
	Token string `json:"token"`
}

func (h *AuthHTTPHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rh := core_http_response.NewHTTPResponseHandler(w, h.logger)

	var req SignInRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		rh.ErrorResponse(err, "invalid request")
		return
	}

	token, err := h.authService.SignIn(ctx, req.Email, req.Password)
	if err != nil {
		rh.ErrorResponse(err, "authentication failed")
		return
	}

	resp := SignInResponse{
		Token: token,
	}
	rh.JSONResponse(resp, http.StatusOK)
}
