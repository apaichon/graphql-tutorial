package auth

import (
	"fmt"
	"graphql-api/pkg/data"
	"graphql-api/pkg/data/models"
	_ "github.com/mattn/go-sqlite3"
)

// UserRepo represents the repository for user operations
type UserRepo struct {
	DB *data.DB
}

// NewUserRepo creates a new instance of UserRepo
func NewUserRepo() *UserRepo {
	db := data.NewDB()
	return &UserRepo{DB: db}
}

// Get Users fetches users from the database with support for text search, limit, and offset
func (cr *UserRepo) GetUsersBySearchText(searchText string, limit, offset int) ([]*models.UserModel, error) {
	var users []*models.UserModel

	query := fmt.Sprintf(`
            SELECT * FROM user
             Where user_name like '%%%s%%' OR password like '%%%s%%' OR salt like '%%%s%%'
            LIMIT ? OFFSET ?
        `, searchText, searchText, searchText)

	rows, err := cr.DB.Query(query, limit, offset)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.UserModel
		err := rows.Scan(
			&user.UserId,
			&user.Username,
			&user.Password,
			&user.Salt,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// Get UserByID retrieves a user by its ID from the database
func (cr *UserRepo) GetUserByID(id int) (*models.UserModel, error) {
	var user models.UserModel
	// Execute query to get a user by ID from the database
	row, err := cr.DB.QueryRow("SELECT * FROM user WHERE user_id = ?", id)

	if err != nil {
		return &user, nil
	}

	row.Scan(
		&user.UserId,
		&user.Username,
		&user.Password,
		&user.Salt,
		&user.CreatedAt,
		&user.CreatedBy,
	)

	return &user, nil
}

// Get UserByID retrieves a user by its ID from the database
func (cr *UserRepo) GetUserByName(name string) (*models.UserModel, error) {
	var user models.UserModel
	// Execute query to get a user by ID from the database
	row, err := cr.DB.QueryRow("SELECT * FROM user WHERE user_name = ?", name)

	if err != nil {
		return nil, fmt.Errorf("error query: %v", err)
	}

	row.Scan(
		&user.UserId,
		&user.Username,
		&user.Password,
		&user.Salt,
		&user.CreatedAt,
		&user.CreatedBy,
	)

	return &user, nil
}

// Insert User inserts a new user into the database
func (cr *UserRepo) InsertUser(user *models.UserModel) (int64, error) {
	// Execute insert query to insert a new user into the database
	result, err := cr.DB.Insert("INSERT INTO user (user_id,user_name,password,salt,created_at) VALUES ({?,?,?,?,?})",
		user.UserId, user.Username, user.Password, user.Salt, user.CreatedAt)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// Update User updates an existing user in the database
func (cr *UserRepo) UpdateUser(user *models.UserModel) (int64, error) {
	// Execute update query to update an existing user in the database
	result, err := cr.DB.Update("UPDATE user SET user_id=?,user_name=?,password=?,salt=? where user_id=?",
		user.UserId, user.Username, user.Password, user.Salt, user.UserId)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Delete User deletes a user from the database
func (cr *UserRepo) DeleteUser(id int) (int64, error) {
	// Execute delete query to delete a user from the database
	result, err := cr.DB.Delete("DELETE FROM user WHERE user_id=?", id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
