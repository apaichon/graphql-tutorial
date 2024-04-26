package contact

import (
	"fmt"
	"graphql-api/pkg/data"
	"graphql-api/pkg/data/models"
	_ "github.com/mattn/go-sqlite3"
)

// ContactRepo represents the repository for contact operations
type ContactRepo struct {
	DB *data.DB
}

// NewContactRepo creates a new instance of ContactRepo
func NewContactRepo() *ContactRepo {
	db := data.NewDB()
	return &ContactRepo{DB: db}
}

// Get Contacts fetches contacts from the database with support for text search, limit, and offset
func (cr *ContactRepo) GetContactsBySearchText(searchText string, limit, offset int) ([]*models.ContactModel, error) {
	var contacts []*models.ContactModel

	query := fmt.Sprintf(`
            SELECT * FROM contact
             Where name like '%%%s%%' OR first_name like '%%%s%%' OR last_name like '%%%s%%' OR email like '%%%s%%' OR phone like '%%%s%%' OR address like '%%%s%%' OR photo_path like '%%%s%%'
            LIMIT ? OFFSET ?
        `, searchText, searchText, searchText, searchText, searchText, searchText, searchText)

	rows, err := cr.DB.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var contact models.ContactModel
		err := rows.Scan(
			&contact.ContactId,
			&contact.Name,
			&contact.FirstName,
			&contact.LastName,
			&contact.GenderId,
			&contact.Dob,
			&contact.Email,
			&contact.Phone,
			&contact.Address,
			&contact.PhotoPath,
			&contact.CreatedAt,
			&contact.CreatedBy,
		)
		if err != nil {
			return nil, err
		}
		contacts = append(contacts, &contact)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return contacts, nil
}


// Get Contacts fetches contacts from the database with support for text search, limit, and offset
func (cr *ContactRepo) GetContactsBySearchTextPagination(searchText string, page, pageSize int) ([]*models.ContactModel, *models.PaginationModel, error) {
	var contacts []*models.ContactModel

	query := fmt.Sprintf(`
            SELECT * FROM contact
             Where name like '%%%s%%' OR first_name like '%%%s%%' OR last_name like '%%%s%%' OR email like '%%%s%%' OR phone like '%%%s%%' OR address like '%%%s%%' OR photo_path like '%%%s%%'
        `, searchText, searchText, searchText, searchText, searchText, searchText, searchText)
	offset := (page - 1) * pageSize
	limit := pageSize

	pagination := data.NewPagination(page, pageSize, query, limit, offset)

    pager, err := pagination.GetPageData(cr.DB)
    if err != nil {
        return nil, nil, err
    }
	
	query = query + " LIMIT ? OFFSET ?"
	
	rows, err := cr.DB.Query(query, limit, offset)
	if err != nil {
		return nil,nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var contact models.ContactModel
		err := rows.Scan(
			&contact.ContactId,
			&contact.Name,
			&contact.FirstName,
			&contact.LastName,
			&contact.GenderId,
			&contact.Dob,
			&contact.Email,
			&contact.Phone,
			&contact.Address,
			&contact.PhotoPath,
			&contact.CreatedAt,
			&contact.CreatedBy,
		)
		if err != nil {
			return nil,nil, err
		}
		contacts = append(contacts, &contact)
	}

	if err := rows.Err(); err != nil {
		return nil,nil, err
	}


	return contacts, pager, nil
}

// Get ContactByID retrieves a contact by its ID from the database
func (cr *ContactRepo) GetContactByID(id int) (*models.ContactModel, error) {
	var contact models.ContactModel
	// Execute query to get a contact by ID from the database
	row, err := cr.DB.QueryRow("SELECT * FROM contact WHERE contact_id = ?", id)

	if err != nil {
		return &contact, nil
	}

	row.Scan(
		&contact.ContactId,
		&contact.Name,
		&contact.FirstName,
		&contact.LastName,
		&contact.GenderId,
		&contact.Dob,
		&contact.Email,
		&contact.Phone,
		&contact.Address,
		&contact.PhotoPath,
		&contact.CreatedAt,
		&contact.CreatedBy,

	)

	return &contact, nil
}

// Insert Contact inserts a new contact into the database
func (cr *ContactRepo) InsertContact(contact *models.ContactModel) (int64, error) {
	// Execute insert query to insert a new contact into the database
	result, err := cr.DB.Insert("INSERT INTO contact (contact_id,name,first_name,last_name,gender_id,dob,email,phone,address,photo_path) VALUES ({?,?,?,?,?,?,?,?,?,?})",
		contact.ContactId, contact.Name, contact.FirstName, contact.LastName, contact.GenderId, contact.Dob, contact.Email, contact.Phone, contact.Address, contact.PhotoPath)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// Update Contact updates an existing contact in the database
func (cr *ContactRepo) UpdateContact(contact *models.ContactModel) (int64, error) {
	// Execute update query to update an existing contact in the database
	result, err := cr.DB.Update("UPDATE contact SET contact_id=?,name=?,first_name=?,last_name=?,gender_id=?,dob=?,email=?,phone=?,address=?,photo_path=? where contact_id=?",
		contact.ContactId, contact.Name, contact.FirstName, contact.LastName, contact.GenderId, contact.Dob, contact.Email, contact.Phone, contact.Address, contact.PhotoPath, contact.ContactId)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Delete Contact deletes a contact from the database
func (cr *ContactRepo) DeleteContact(id int) (int64, error) {
	// Execute delete query to delete a contact from the database
	result, err := cr.DB.Delete("DELETE FROM contact WHERE contact_id=?", id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
