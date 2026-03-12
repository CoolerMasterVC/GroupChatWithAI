package models

import "time"

// Segment представляет один сегмент сообщения от агентного уровня
// @Description Структура сегмента, отправляемого от агента на транспортный уровень

type Segment struct {
	SegmentNumber int       `json:"segment_number" example:"1"`               // Номер сегмента в сообщении
	TotalSegments int       `json:"total_segments" example:"5"`               // Общее количество сегментов
	Username      string    `json:"username" example:"alice"`                 // Имя отправителя
	SendTime      time.Time `json:"send_time" example:"2025-03-12T10:00:00Z"` // Время отправки (идентификатор сообщения)
	Payload       string    `json:"payload" example:"часть текста..."`        // Полезная нагрузка (часть сообщения)
}
