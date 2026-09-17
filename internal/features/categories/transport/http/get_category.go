package transport

import (
	"net/http"

	core_http_request "github.com/loundxr/expense-tracker/internal/core/transport/http/request"
	core_http_response "github.com/loundxr/expense-tracker/internal/core/transport/http/response"
)

type GetCategoryResponse CategoryDTOResponse

func (h *CategoriesHTTPHandler) GetCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rh := core_http_response.NewHTTPResponseHandler(w, h.logger)

	id, err := core_http_request.GetInt64IDPathValue(r, "id")
	if err != nil {
		rh.ErrorResponse(err, "invalid id format")
		return
	}

	category, err := h.categoriesService.GetCategory(ctx, id)
	if err != nil {
		rh.ErrorResponse(err, "failed to get category")
		return
	}

	response := GetCategoryResponse(categoryDTOFromDomain(category))
	rh.JSONResponse(response, http.StatusOK)
}
