package models

// SuccessResponse — стандартный успешный ответ
type SuccessResponse struct {
	Status string `json:"status" example:"accepted"`
}

// ErrorResponse — ответ с ошибкой
type ErrorResponse struct {
	Error string `json:"error" example:"invalid json"`
}
