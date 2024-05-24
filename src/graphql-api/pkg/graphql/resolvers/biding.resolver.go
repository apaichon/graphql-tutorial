package resolvers

import (
	"fmt"
	"strconv"
	"time"
	"graphql-api/internal/biding"
	"graphql-api/internal/cache"

	"graphql-api/pkg/data/models"
	"github.com/samborkent/uuidv7"

	"github.com/graphql-go/graphql"
	"graphql-api/internal/subscription"
)

func GetBidingRoomByIdResolve(params graphql.ResolveParams) (interface{}, error) {
	id := params.Args["id"].(int)
	bidingRoomRepo := biding.NewBidingRoomRepo()

	// Fetch biding room from the database
	bidingRoom, err := bidingRoomRepo.GetRoomById(id)
	fmt.Printf("Biding room:%v", bidingRoom)
	if err != nil {
		return nil, err
	}
	go cache.SetCacheResolver(params, bidingRoom)
	return bidingRoom, nil
}

func GetBidingTop5Resolve(params graphql.ResolveParams) (interface{}, error) {
	id := params.Args["id"].(int)
	bidingRepo := biding.NewBidingRepo()

	// Fetch biding room from the database
	biding, err := bidingRepo.GetTop5(id)
	if err != nil {
		return nil, err
	}
	go cache.SetCacheResolver(params, biding)
	return biding, nil
}

//BidingResolve
func CreateBidingResolve(params graphql.ResolveParams) (interface{}, error) {
	// Map input fields to Biding struct
	input := params.Args["input"].(map[string]interface{})

	bidingInput := models.BidingModel{
		RoomId: input["room_id"].(int),
		Bidder:  input["bidder"].(string),
		BidPrice:  input["bid_price"].(float64),
		BidTime:       time.Now(),
	}

	bidingRepo := biding.NewBidingRepo()

	// Insert Biding to the database
	id, err := bidingRepo.InsertBiding(&bidingInput)
	if err != nil {
		return nil, err
	}

	bidingId := int(id)
	bidingInput.BidId = bidingId
	go cache.RemoveGetCacheResolver("BidingQueries")
	
	uuid := uuidv7.New().String()

	go subscription.SendMessage(subscription.SubscribeMessage{Id: uuid, Content: "Create Biding Id:" + strconv.Itoa( bidingId), Timestamp: time.Now()})
	return bidingInput, nil

}
