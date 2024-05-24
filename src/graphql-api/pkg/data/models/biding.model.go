package models

import (
	"time"
)
// Biding represents the biding table
type BidingModel struct {
    BidId    int       `json:"bid_id"`
    RoomId   int       `json:"room_id"`
    Bidder   string    `json:"bidder"`
    BidPrice float64   `json:"bid_price"`
    BidTime  time.Time `json:"bid_time"`
}