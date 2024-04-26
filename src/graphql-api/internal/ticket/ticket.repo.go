package ticket

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"graphql-api/internal/event"
	"graphql-api/pkg/data"
	"graphql-api/pkg/data/models"
)

// TicketRepo represents the repository for ticket operations
type TicketRepo struct {
	DB *data.DB
}

// NewTicketRepo creates a new instance of TicketRepo
func NewTicketRepo() *TicketRepo {
	db := data.NewDB()
	return &TicketRepo{DB: db}
}

// Get Tickets fetches tickets from the database with support for text search, limit, and offset
func (cr *TicketRepo) GetTicketsBySearchText(searchText string, limit, offset int) ([]*models.TicketModel, error) {
	var tickets []*models.TicketModel

	query := fmt.Sprintf(`
            SELECT * FROM ticket
             Where type like '%%%s%%'
            LIMIT ? OFFSET ?
        `, searchText)

	rows, err := cr.DB.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var ticket models.TicketModel
		err := rows.Scan(
			&ticket.TicketId,
			&ticket.Type,
			&ticket.Price,
			&ticket.EventId,
			&ticket.CreatedAt,
			&ticket.CreatedBy,
		)
		if err != nil {
			return nil, err
		}
		tickets = append(tickets, &ticket)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tickets, nil
}

// Get Tickets fetches tickets from the database with support for text search, limit, and offset
func (cr *TicketRepo) GetTicketEventsBySearchText(searchText string, limit, offset int) ([]models.TicketEventModel, error) {
	var tickets []models.TicketModel
	var eventIds []int

	query := fmt.Sprintf(`
            SELECT * FROM ticket
             Where type like '%%%s%%'
            LIMIT ? OFFSET ?
        `, searchText)

	rows, err := cr.DB.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	
	for rows.Next() {
		var ticket models.TicketModel
		err := rows.Scan(
			&ticket.TicketId,
			&ticket.Type,
			&ticket.Price,
			&ticket.EventId,
			&ticket.CreatedAt,
			&ticket.CreatedBy,
		)
		if err != nil {
			return nil, err
		}

		tickets = append(tickets, ticket)
		eventIds = append(eventIds, ticket.EventId)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	eventRepo := event.NewEventRepo()
	events, err := eventRepo.GetEventsByIDs(eventIds)

	if err != nil {
		return nil, err
	}

	// Create a map from Event slice by EventId
	eventList:= models.CopyEventSlice(events)
	events =nil 
	eventMap := models.CreateEventMap(eventList)
	ticketEvents := models.MapTicketsWithEvents(tickets, eventMap)
	eventMap =nil
	return ticketEvents, nil
}

// ดร สุวิทย์

// Get TicketByID retrieves a ticket by its ID from the database
func (cr *TicketRepo) GetTicketByID(id int) (*models.TicketModel, error) {
	var ticket models.TicketModel
	// Execute query to get a ticket by ID from the database
	row, err := cr.DB.QueryRow("SELECT * FROM ticket WHERE ticket_id = ?", id)

	if err != nil {
		return &ticket, nil
	}

	row.Scan(
		&ticket.TicketId,
		&ticket.Type,
		&ticket.Price,
		&ticket.EventId,
		&ticket.CreatedAt,
		&ticket.CreatedBy,
	)

	return &ticket, nil
}

// Insert Ticket inserts a new ticket into the database
func (cr *TicketRepo) InsertTicket(ticket *models.TicketModel) (int64, error) {
	// Execute insert query to insert a new ticket into the database
	result, err := cr.DB.Insert("INSERT INTO ticket (ticket_id,type,price,event_id) VALUES ({?,?,?,?})",
		ticket.TicketId, ticket.Type, ticket.Price, ticket.EventId)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// Update Ticket updates an existing ticket in the database
func (cr *TicketRepo) UpdateTicket(ticket *models.TicketModel) (int64, error) {
	// Execute update query to update an existing ticket in the database
	result, err := cr.DB.Update("UPDATE ticket SET ticket_id=?,type=?,price=?,event_id=? where ticket_id=?",
		ticket.TicketId, ticket.Type, ticket.Price, ticket.EventId, ticket.TicketId)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Delete Ticket deletes a ticket from the database
func (cr *TicketRepo) DeleteTicket(id int) (int64, error) {
	// Execute delete query to delete a ticket from the database
	result, err := cr.DB.Delete("DELETE FROM ticket WHERE ticket_id=?", id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
