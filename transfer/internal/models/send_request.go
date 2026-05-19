package models

import "time"

type SendRequest struct {
	Sender    string    `json:"sender"`
	Timestamp time.Time `json:"timestamp"`
	Payload   string    `json:"payload"`
}
