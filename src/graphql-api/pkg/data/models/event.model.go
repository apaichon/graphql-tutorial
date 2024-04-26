package models
import ("time")

// Event represents a Event record in the database
type EventModel struct {
	EventId int `json:"event_id"`
	ParentEventId *int `json:"parent_event_id"`
	Name string `json:"name"`
	Description string `json:"description"`
	StartDate time.Time `json:"start_date"`
	EndDate time.Time `json:"end_date"`
	LocationId *int `json:"location_id"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string `json:"created_by"`
}