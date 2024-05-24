package biding

import (
	"fmt"
	"graphql-api/pkg/data"
	"graphql-api/pkg/data/models"

	_ "github.com/mattn/go-sqlite3"
)

// BidingRepo represents the repository for contact operations
type BidingRoomRepo struct {
	DB *data.DB
}

// NewBidingRepo creates a new instance of BidingRepo
func NewBidingRoomRepo() *BidingRoomRepo {
	db := data.NewDB()
	return &BidingRoomRepo{DB: db}
}

// Get ContactByID retrieves a contact by its ID from the database
func (cr *BidingRoomRepo) GetRoomById(roomId int) (*models.BidingRoomModel, error) {
	var room models.BidingRoomModel
	// Execute query to get a contact by ID from the database
	row, err := cr.DB.QueryRow("SELECT * FROM biding_room WHERE room_id = ?", roomId)
	fmt.Printf("roomId:%v row:%v", roomId, row)
	if err != nil {
		return &room, nil
	}

	row.Scan(
		&room.RoomId,
		&room.StartDate,
		&room.EndDate,
		&room.ProductName,
		&room.FloorPrice,
		&room.ProductImage,
	)

	return &room, nil
}
