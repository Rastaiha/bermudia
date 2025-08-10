package models

// APIResponse is the generic response format for all API endpoints
type APIResponse struct {
	OK     bool   `json:"ok"`
	Error  string `json:"error,omitempty"`
	Result any    `json:"result,omitempty"`
}
