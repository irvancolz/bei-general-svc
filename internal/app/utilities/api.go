package utilities

type APIResponse struct {
	Code    int                      `json:"code"`
	Message string                   `json:"message"`
	Data    []map[string]interface{} `json:"data"`
}
type APIResponseInterface struct {
	Code    int                      `json:"code"`
	Message string                   `json:"message"`
	Data    map[string]interface{} `json:"data"`
}
