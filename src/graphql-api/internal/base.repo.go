package internal

import (
	"fmt"
	"graphql-api/pkg/data"
	"graphql-api/pkg/data/models"
	"reflect"
	"strings"
	"sync"
	"unicode"

	// _ "github.com/mattn/go-sqlite3"
	_ "modernc.org/sqlite"
)

var instance *BaseRepo
var once sync.Once

// BaseRepo represents the repository for contact operations
type BaseRepo struct {
	DB *data.DB
}

// NewBaseRepo creates a new instance of BaseRepo
func NewBaseRepo() *BaseRepo {
	db := data.NewDB()
	once.Do(func() {
		instance = &BaseRepo{DB: db}
	})
	return instance
}

// Get Generic fetches data from the database with support for text search, limit, and offset
func (b *BaseRepo) FindBySearchText(itf interface{}, filter models.FilterModel) ([]interface{}, error) {
	// Create a slice of empty interfaces to hold the data
	var result []interface{}
	where := ""

	if len(filter.FilterCondition) > 0 {
		where = " Where " + filter.FilterCondition
	}

	query := fmt.Sprintf(`
	SELECT * FROM %s
	%s
	LIMIT ? OFFSET ?`, filter.Table, where)

	// fmt.Println("Query:",query)

	rows, err := b.DB.Query(query, filter.Limit, filter.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error())
	}

	// Create a slice of empty interfaces to hold the data for a single row
	values := make([]interface{}, len(columns))
	for i := range values {
		var v interface{}
		values[i] = &v
	}

	// Iterate over the result set
	for rows.Next() {
		// Scan the row into the slice of empty interfaces
		err := rows.Scan(values...)

		if err != nil {
			panic(err.Error())
		}

		// Create a new struct instance using reflection
		s := reflect.New(reflect.TypeOf(itf)).Elem()

		// Map the data from the empty interfaces to the struct fields using reflection
		for i, v := range values {

			field := s.FieldByName(SnakeToPascalCase(columns[i]))
			// fmt.Println("field value:",field)
			if field.IsValid() {
				field.Set(reflect.ValueOf(*(v.(*interface{}))))
			}
		}

		// Append the struct to the result slice
		result = append(result, s.Interface())
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		panic(err.Error())
	}

	return result, nil
}

func titleCase(str string) string {
	var result strings.Builder
	prev := ' '
	for _, char := range str {
		if prev == ' ' || prev == '"' || prev == '(' || prev == '[' || prev == '`' || prev == '\'' {
			result.WriteRune(unicode.ToTitle(char))
		} else {
			result.WriteRune(unicode.ToLower(char))
		}
		prev = char
	}
	return result.String()
}

// SnakeToPascalCase converts a snake_case string to PascalCase
func SnakeToPascalCase(s string) string {
	// Split the string by underscores
	parts := strings.Split(s, "_")

	// Convert each part to PascalCase
	for i, part := range parts {
		parts[i] = titleCase(part)
	}

	// Join the parts together
	return strings.Join(parts, "")
}
