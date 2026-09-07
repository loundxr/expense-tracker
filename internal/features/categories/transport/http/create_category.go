package categories_transport_http

import (
	"net/http"

	"github.com/loundxr/expense-tracker/internal/core/domain"
	core_ctx "github.com/loundxr/expense-tracker/internal/core/transport/http/context"
	core_http_request "github.com/loundxr/expense-tracker/internal/core/transport/http/request"
	core_http_response "github.com/loundxr/expense-tracker/internal/core/transport/http/response"
)

type CreateCategoryRequest struct {
	Name string `json:"name" validate:"required,min=2,max=30"`
}

type CreateCategoryResponse CategoryDTOResponse

func (h *CategoriesHTTPHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rh := core_http_response.NewHTTPResponseHandler(w, h.logger)

	var req CreateCategoryRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		rh.ErrorResponse(err, "failed to decode and validate request")
		return
	}

	uid := core_ctx.GetUserID(ctx)
	uninitCat := domain.NewUninitializedCategory(req.Name, &uid)

	category, err := h.categoriesService.CreateCategory(ctx, uninitCat)
	if err != nil {
		rh.ErrorResponse(err, "failed to create user category")
		return
	}

	response := CreateCategoryResponse(categoryDTOFromDomain(category))
	rh.JSONResponse(response, http.StatusOK)
}
