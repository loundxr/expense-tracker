package core_http_response

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	core_errors "github.com/loundxr/expense-tracker/internal/core/errors"
)

type HTTPResponseHandler struct {
	log *slog.Logger
	w   http.ResponseWriter
}

func NewHTTPResponseHandler(w http.ResponseWriter, log *slog.Logger) *HTTPResponseHandler {
	return &HTTPResponseHandler{
		log: log,
		w:   w,
	}
}

func (h *HTTPResponseHandler) NoContentResponse() {
	h.w.WriteHeader(http.StatusNoContent)
}

func (h *HTTPResponseHandler) JSONResponse(responseBody any, statusCode int) {
	h.w.Header().Set("Content-Type", "application/json")
	h.w.WriteHeader(statusCode)
	if err := json.NewEncoder(h.w).Encode(responseBody); err != nil {
		h.log.Error("error encoding json response", slog.String("error", err.Error()))
	}
}

func (h *HTTPResponseHandler) ErrorResponse(err error, msg string) {
	var statusCode int

	switch {
	case errors.Is(err, core_errors.ErrInvalidArgument):
		statusCode = http.StatusBadRequest
	case errors.Is(err, core_errors.ErrAlreadyExists) ||
		errors.Is(err, core_errors.ErrConflict):
		statusCode = http.StatusConflict
	case errors.Is(err, core_errors.ErrNotFound):
		statusCode = http.StatusNotFound
	case errors.Is(err, core_errors.ErrForbidden):
		statusCode = http.StatusForbidden
	default:
		statusCode = http.StatusInternalServerError
	}

	if statusCode == http.StatusInternalServerError {
		h.log.Error(msg, slog.String("error", err.Error()))
	} else {
		h.log.Warn(msg, slog.String("error", err.Error()), slog.Int("status_code", statusCode))
	}
	h.errorResponse(statusCode, err, msg)
}

func (h *HTTPResponseHandler) errorResponse(code int, err error, msg string) {
	errStr := err.Error()
	if code == http.StatusInternalServerError {
		errStr = "internal server error"
	}
	resp := errorResponse{
		Error:   errStr,
		Message: msg,
	}
	h.JSONResponse(resp, code)
}
