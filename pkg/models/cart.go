package models

import "time"

type Cart struct {
	Id        int64     `json:"id" db:"id"`
	UserId    int64     `json:"user_id" db:"user_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
