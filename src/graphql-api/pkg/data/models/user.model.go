package models

import (
	"time"
)

type UserModel struct {
	UserId    int       `json:"user_id"`
	Username  string    `json:"user_name"`
	Password  string    `json:"password"`
	Salt      string    `json:"salt"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
}
