package models

// SuccessResponse успешный ответ
type SuccessResponse struct {
	Status string `json:"status" example:"accepted"` // Статус операции
}

// ErrorResponse ответ с ошибкой
type ErrorResponse struct {
	Error string `json:"error" example:"invalid json"` // Текст ошибки
}
