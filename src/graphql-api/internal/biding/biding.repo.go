package biding

import (
	"graphql-api/pkg/data"
	"graphql-api/pkg/data/models"
	_ "github.com/mattn/go-sqlite3"
)

// BidingRepo represents the repository for contact operations
type BidingRepo struct {
	DB *data.DB
}

// NewBidingRepo creates a new instance of BidingRepo
func NewBidingRepo() *BidingRepo {
	db := data.NewDB()
	return &BidingRepo{DB: db}
}

// Get Contacts fetches contacts from the database with support for text search, limit, and offset
func (cr *BidingRepo) GetTop5(roomId  int) ([]*models.BidingModel, error) {
	var top5 []*models.BidingModel

	query := `SELECT * FROM biding
             Where  room_id = ? order by bid_price desc, bid_time desc limit 5`

	rows, err := cr.DB.Query(query,roomId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var biding models.BidingModel
		err := rows.Scan(
			&biding.BidId,
			&biding.RoomId,
			&biding.Bidder,
			&biding.BidPrice,
			&biding.BidTime,
		
		)
		if err != nil {
			return nil, err
		}
		top5 = append(top5, &biding)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return top5, nil
}

// Insert Contact inserts a new contact into the database
func (bd *BidingRepo) InsertBiding(biding *models.BidingModel) (int64, error) {
	// Execute insert query to insert a new biding into the database

	result, err := bd.DB.Insert("INSERT INTO biding (room_id,bidder,bid_price, bid_time) VALUES (?,?,?,?)",
		 biding.RoomId,biding.Bidder, biding.BidPrice, biding.BidTime)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

