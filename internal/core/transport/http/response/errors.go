package core_http_response

type errorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
