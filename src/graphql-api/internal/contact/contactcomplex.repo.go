package contact

import (
	"fmt"
	"graphql-api/pkg/data/models"
	"graphql-api/internal"
	_ "github.com/mattn/go-sqlite3"
)

// ContactRepo represents the repository for contact operations
type ContactComplexRepo struct {
	internal.BaseRepo
}

// NewContactRepo creates a new instance of ContactRepo
func NewContactComplexRepo() *ContactComplexRepo {	
	return &ContactComplexRepo{BaseRepo: *internal.NewBaseRepo()}
}
const SEARCH_CONDITION = `name like '%%%s%%' OR first_name like '%%%s%%' OR last_name like '%%%s%%' OR email like '%%%s%%' OR phone like '%%%s%%' OR address like '%%%s%%' OR photo_path like '%%%s%%'`

func (cr *ContactComplexRepo) GetContactsBySearchText(searchText string, limit, offset int) ([]interface{}, error) {
	where :=  fmt.Sprintf(SEARCH_CONDITION,
	searchText, searchText, searchText, searchText, searchText, searchText, searchText)
	
	filter := models.FilterModel {
		Table: "contact",
		SearchText: searchText,
		Limit:limit,
		Offset: offset,
		FilterCondition: where,
	}
	
	contacts, err := cr.BaseRepo.FindBySearchText(models.ContactModel{},filter)
	if err != nil {
		return nil, err
	}
	return contacts, nil
}

