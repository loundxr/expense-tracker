package categories_transport_http

import (
	"net/http"

	"github.com/loundxr/expense-tracker/internal/core/domain"
	core_http_request "github.com/loundxr/expense-tracker/internal/core/transport/http/request"
	core_http_response "github.com/loundxr/expense-tracker/internal/core/transport/http/response"
	core_types "github.com/loundxr/expense-tracker/internal/core/types"
)

type PatchCategoryRequest struct {
	Name core_types.Nullable[string] `json:"name"`
}

type PatchCategoryResponse CategoryDTOResponse

func (h *CategoriesHTTPHandler) PatchCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rh := core_http_response.NewHTTPResponseHandler(w, h.logger)

	id, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil || id <= 0 {
		rh.ErrorResponse(err, "invalid 'id' format")
		return
	}

	var req PatchCategoryRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		rh.ErrorResponse(err, "failed to decode and validate request")
		return
	}

	patch := categoryPatchFromRequest(req)
	categoryDomain, err := h.categoriesService.PatchCategory(ctx, id, patch)
	if err != nil {
		rh.ErrorResponse(err, "failed to patch category")
		return
	}

	response := PatchCategoryResponse(categoryDTOFromDomain(categoryDomain))
	rh.JSONResponse(response, http.StatusOK)
}

func categoryPatchFromRequest(req PatchCategoryRequest) domain.CategoryPatch {
	return domain.NewCategoryPatch(req.Name.ToDomain())
}
