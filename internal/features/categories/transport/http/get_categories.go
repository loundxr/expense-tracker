package transport

import (
	"net/http"

	core_http_response "github.com/loundxr/expense-tracker/internal/core/transport/http/response"
)

type GetCategoriesResponse struct {
	Categories []CategoryDTOResponse `json:"categories"`
}

func (h *CategoriesHTTPHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rh := core_http_response.NewHTTPResponseHandler(w, h.logger)

	categories, err := h.categoriesService.GetCategories(ctx)
	if err != nil {
		rh.ErrorResponse(err, "failed to get categories")
		return
	}

	response := NewGetCategoriesResponse(categoryDTOsFromDomains(categories))
	rh.JSONResponse(response, http.StatusOK)
}

func NewGetCategoriesResponse(categories []CategoryDTOResponse) GetCategoriesResponse {
	if categories == nil {
		categories = make([]CategoryDTOResponse, 0)
	}
	return GetCategoriesResponse{
		Categories: categories,
	}
}
