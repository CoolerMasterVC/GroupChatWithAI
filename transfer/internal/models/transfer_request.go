package models

import "time"

// TransferRequest представляет один сегмент от агентного уровня
// @Description Сегмент разбитого сообщения
type TransferRequest struct {
	SegmentNumber int       `json:"segment_number" example:"1"`               // Номер сегмента (с 1)
	TotalSegments int       `json:"total_segments" example:"5"`               // Всего сегментов в сообщении
	Username      string    `json:"username" example:"alice"`                 // Имя отправителя
	SendTime      time.Time `json:"send_time" example:"2025-03-29T12:00:00Z"` // Время отправки (ID сообщения)
	Payload       string    `json:"payload" example:"часть текста..."`        // Данные сегмента
}
