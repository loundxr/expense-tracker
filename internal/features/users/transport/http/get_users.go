package users_transport_http

import (
	"fmt"
	"net/http"

	core_http_request "github.com/loundxr/expense-tracker/internal/core/transport/http/request"
	core_http_response "github.com/loundxr/expense-tracker/internal/core/transport/http/response"
)

type GetUsersResponse []UserDTOResponse

func (h *UserHTTPHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rh := core_http_response.NewHTTPResponseHandler(w, h.logger)

	limit, offset, err := getLimitOffsetQueryParams(r)
	if err != nil {
		rh.ErrorResponse(err, "failed to get limit/offset query params")
		return
	}

	users, err := h.userService.GetUsers(ctx, limit, offset)
	if err != nil {
		rh.ErrorResponse(err, "failed to get users list")
		return
	}

	usersDTO := GetUsersResponse(usersDTOFromDomains(users))
	rh.JSONResponse(usersDTO, http.StatusOK)
}

func getLimitOffsetQueryParams(r *http.Request) (*int, *int, error) {
	const (
		limitQueryKey  = "limit"
		offsetQueryKey = "offset"
	)

	limit, err := core_http_request.GetIntQueryParams(r, limitQueryKey)
	if err != nil {
		return nil, nil, fmt.Errorf("get %s query param: %w", limitQueryKey, err)
	}

	offset, err := core_http_request.GetIntQueryParams(r, offsetQueryKey)
	if err != nil {
		return nil, nil, fmt.Errorf("get %s query param: %w", offsetQueryKey, err)
	}
	return limit, offset, nil
}
