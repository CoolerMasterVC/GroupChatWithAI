package models

import "time"

// SendRequest представляет полное сообщение от прикладного уровня
// @Description Полное сообщение, которое транспортный уровень передаёт агентному уровню без изменений
type SendRequest struct {
	Username string    `json:"username" example:"alice"`                 // Отправитель
	Data     string    `json:"data" example:"Привет, как дела?"`         // Текст сообщения
	SendTime time.Time `json:"send_time" example:"2025-03-29T12:00:00Z"` // Время отправки
}
