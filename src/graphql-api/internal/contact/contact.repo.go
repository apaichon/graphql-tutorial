package contact

import (
	"fmt"
	"strings"
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
	result, err := cr.DB.Insert("INSERT INTO contact (name,first_name,last_name,gender_id,dob,email,phone,address,photo_path,created_at,created_by) VALUES (?,?,?,?,?,?,?,?,?,?,?)",
		 contact.Name, contact.FirstName, contact.LastName, contact.GenderId, contact.Dob, contact.Email, contact.Phone, contact.Address, contact.PhotoPath, contact.CreatedAt, contact.CreatedBy)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (cr *ContactRepo) InsertContacts(contacts []*models.ContactModel) ( int64, error) {
  
	// Step 1: Query the count of records before insertion
    var countBefore int64
    row,err := cr.DB.QueryRow("SELECT COUNT(*) FROM contact")
	
    if err != nil {
        return 0, fmt.Errorf("failed to query count before insertion: %v", err)
    }
	row.Scan(&countBefore)

    // Begin a transaction by starting a deferred transaction
    _, err = cr.DB.Exec("BEGIN TRANSACTION")
    if err != nil {
        return 0, err
    }
    defer func() {
        // Rollback the transaction if there's an error and it hasn't been committed
        if err != nil {
            _, rollbackErr := cr.DB.Exec("ROLLBACK")
            if rollbackErr != nil {
                err = fmt.Errorf("rollback failed: %v, original error: %v", rollbackErr, err)
            }
            return
        }
        // Commit the transaction if no error occurred
        _, commitErr := cr.DB.Exec("COMMIT")
        if commitErr != nil {
            err = fmt.Errorf("commit failed: %v", commitErr)
        }
    }()

    // Prepare the SQL statement for batch insertion
    stmt, err := cr.DB.Prepare("INSERT INTO contact (name, first_name, last_name, gender_id, dob, email, phone, address, photo_path, created_at, created_by) VALUES " + placeholders(len(contacts)))
    if err != nil {
        return 0, err
    }
    defer stmt.Close()

    // Prepare the slice to hold the arguments for the prepared statement
    args := make([]interface{}, 0, len(contacts)*11)

    // Flatten the contacts into a single slice of values
    for _, contact := range contacts {
        args = append(args, contact.Name, contact.FirstName, contact.LastName, contact.GenderId, contact.Dob, contact.Email, contact.Phone, contact.Address, contact.PhotoPath, contact.CreatedAt, contact.CreatedBy)
    }

    // Execute the prepared statement with the concatenated values
    _, err = stmt.Exec(args...)
    if err != nil {
        return 0, err
    }

    var countAfter int64
    row, err = cr.DB.QueryRow("SELECT COUNT(*) FROM contact")
	row.Scan(&countAfter)
    if err != nil {
        return 0, fmt.Errorf("failed to query count after insertion: %v", err)
    }
	if countAfter != (countBefore + int64(len(contacts))) {
		return 0, fmt.Errorf("insert batch is not completed: %v", err)
	}

    return  countAfter, nil
}

// placeholders returns a string with n question marks separated by commas, for use in a SQL statement.
func placeholders(n int) string {
    if n <= 0 {
        return ""
    }
    return strings.Repeat("(?,?,?,?,?,?,?,?,?,?,?),", n-1) + "(?,?,?,?,?,?,?,?,?,?,?)"
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
