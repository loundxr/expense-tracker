package core_http_response

import (
	"encoding/json"
	"log/slog"
	"net/http"
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

func (h *HTTPResponseHandler) JSONResponse(responseBody any, statusCode int) {
	h.w.Header().Set("Content-Type", "application/json")
	h.w.WriteHeader(statusCode)
	if err := json.NewEncoder(h.w).Encode(responseBody); err != nil {
		// TODO: write error method to logger
		//h.log.Error("Error encoding JSON response", "error", err)
	}
}

func (h *HTTPResponseHandler) ErrorResponse(err error, msg string) {
	// TODO: rework with logger and switch/case with status codes according to error type
}
