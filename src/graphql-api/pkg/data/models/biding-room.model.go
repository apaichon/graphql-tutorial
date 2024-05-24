package models

import (
	"time"
)
// BidingRoom represents the biding_room table
type BidingRoomModel struct {
    RoomId      int       `json:"room_id"`
    StartDate   time.Time `json:"start_date"`
    EndDate     time.Time `json:"end_date"`
    ProductName string    `json:"product_name"`
    FloorPrice  float64   `json:"floor_price"`
    ProductImage string   `json:"product_image"`
}