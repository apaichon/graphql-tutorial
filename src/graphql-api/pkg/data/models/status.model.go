package models

type Status struct {
	StatusID int    `json:"status_id"`
	StatusText       string `json:"status_text"`
	Message string `json:"message"`
}