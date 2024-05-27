package event

import (
	"fmt"
	// _ "github.com/mattn/go-sqlite3"
	"graphql-api/pkg/data"
	"graphql-api/pkg/data/models"
	"strconv"
	"strings"

	_ "modernc.org/sqlite"
)

// EventRepo represents the repository for event operations
type EventRepo struct {
	DB *data.DB
}

// NewEventRepo creates a new instance of EventRepo
func NewEventRepo() *EventRepo {
	db := data.NewDB()
	return &EventRepo{DB: db}
}

// Get Events fetches events from the database with support for text search, limit, and offset
func (cr *EventRepo) GetEventsBySearchText(searchText string, limit, offset int) ([]*models.EventModel, error) {
	var events []*models.EventModel

	query := fmt.Sprintf(`
            SELECT * FROM event
             Where name like '%%%s%%' OR description like '%%%s%%'
            LIMIT ? OFFSET ?
        `, searchText, searchText)

	rows, err := cr.DB.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var event models.EventModel
		err := rows.Scan(
			&event.EventId,
			&event.ParentEventId,
			&event.Name,
			&event.Description,
			&event.StartDate,
			&event.EndDate,
			&event.LocationId,
			&event.CreatedAt,
			&event.CreatedBy,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, &event)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

// Get EventByID retrieves a event by its ID from the database
func (cr *EventRepo) GetEventByID(id int) (*models.EventModel, error) {
	var event models.EventModel
	// Execute query to get a event by ID from the database
	row, err := cr.DB.QueryRow("SELECT * FROM event WHERE event_id = ?", id)

	if err != nil {
		return &event, nil
	}

	row.Scan(
		&event.EventId,
		&event.ParentEventId,
		&event.Name,
		&event.Description,
		&event.StartDate,
		&event.EndDate,
		&event.LocationId,
		&event.CreatedAt,
		&event.CreatedBy,
	)

	return &event, nil
}

// Get Events fetches events from the database with support for text search, limit, and offset
func (cr *EventRepo) GetEventsByIDs(ids []int) ([]*models.EventModel, error) {
	var events []*models.EventModel
	// Convert each integer to a string
	strArr := make([]string, len(ids))
	for i, num := range ids {
		strArr[i] = strconv.Itoa(num)
	}
	// Join the array of strings with commas and add parentheses
	idList := fmt.Sprintf("(%s)", strings.Join(strArr, ","))

	query := fmt.Sprintf(`SELECT * FROM event
             Where event_id in %v`, idList)

	rows, err := cr.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var event models.EventModel
		err := rows.Scan(
			&event.EventId,
			&event.ParentEventId,
			&event.Name,
			&event.Description,
			&event.StartDate,
			&event.EndDate,
			&event.LocationId,
			&event.CreatedAt,
			&event.CreatedBy,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, &event)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

// Insert Event inserts a new event into the database
func (cr *EventRepo) InsertEvent(event *models.EventModel) (int64, error) {
	// Execute insert query to insert a new event into the database
	result, err := cr.DB.Insert("INSERT INTO event (event_id,parent_event_id,name,description,start_date,end_date,location_id) VALUES ({?,?,?,?,?,?,?})",
		event.EventId, event.ParentEventId, event.Name, event.Description, event.StartDate, event.EndDate, event.LocationId)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// Update Event updates an existing event in the database
func (cr *EventRepo) UpdateEvent(event *models.EventModel) (int64, error) {
	// Execute update query to update an existing event in the database
	result, err := cr.DB.Update("UPDATE event SET event_id=?,parent_event_id=?,name=?,description=?,start_date=?,end_date=?,location_id=? where event_id=?",
		event.EventId, event.ParentEventId, event.Name, event.Description, event.StartDate, event.EndDate, event.LocationId, event.EventId)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Delete Event deletes a event from the database
func (cr *EventRepo) DeleteEvent(id int) (int64, error) {
	// Execute delete query to delete a event from the database
	result, err := cr.DB.Delete("DELETE FROM event WHERE event_id=?", id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
