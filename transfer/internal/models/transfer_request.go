package models

import "time"

type TransferRequest struct {
	SegmentNumber int       `json:"segment_number" example:"1"`
	TotalSegments int       `json:"total_segments" example:"3"`
	Username      string    `json:"username" example:"alice"`
	SendTime      time.Time `json:"send_time" example:"2025-04-25T12:00:00Z"`
	Payload       string    `json:"payload" example:"part of message"`
}
