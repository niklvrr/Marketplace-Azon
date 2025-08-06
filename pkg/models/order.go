package models

import "time"

type Order struct {
	Id       int64     `json:"id" db:"id"`
	UserId   int64     `json:"user_id" db:"user_id"`
	Status   string    `json:"status" db:"status"`
	Total    float64   `json:"total" db:"total"`
	CreateAt time.Time `json:"create_at" db:"create_at"`
}
