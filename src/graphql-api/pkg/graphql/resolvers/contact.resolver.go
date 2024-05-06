package resolvers

import (
	"fmt"
	"graphql-api/internal/contact"
	"graphql-api/pkg/data/models"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/graphql-go/graphql"
)

var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

func validateContact(contactInput map[string]interface{}) map[string]interface{} {
	rules := map[string]interface{}{
		"name":       "required",
		"first_name": "required",
		"last_name":  "required",
		"email":      "omitempty,required,email"}
	errs := validate.ValidateMap(contactInput, rules)
	return errs
}

func CretateContactResolve(params graphql.ResolveParams) (interface{}, error) {
	// Map input fields to Contact struct
	input := params.Args["input"].(map[string]interface{})
	invalids := validateContact(input)

	if len(invalids) > 0 {
		return nil, fmt.Errorf("%v", invalids)
	}

	contactInput := models.ContactModel{

		Name:      input["name"].(string),
		FirstName: input["first_name"].(string),
		LastName:  input["last_name"].(string),
		GenderId:  input["gender_id"].(int),
		Dob:       input["dob"].(time.Time),
		Email:     input["email"].(string),
		Phone:     input["phone"].(string),
		Address:   input["address"].(string),
		PhotoPath: input["photo_path"].(string),
		CreatedBy: "test-api",
		CreatedAt: time.Now(),
	}

	contactRepo := contact.NewContactRepo()

	// Insert Contact to the database
	id, err := contactRepo.InsertContact(&contactInput)
	if err != nil {
		return nil, err
	}
	contactInput.ContactId = int64(id)
	return contactInput, nil
}

func CreateContactsResolve(params graphql.ResolveParams) (interface{}, error) {
	// Map input fields to Contact struct
	contactsArg, ok := params.Args["contacts"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("contacts argument not provided or has wrong type")
	}

	var contacts []*models.ContactModel
	for _, contact := range contactsArg {
		contactMap, ok := contact.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("contact has wrong type")
		}
		contacts = append(contacts, &models.ContactModel{
			// Parse other fields as needed
			Name:      contactMap["name"].(string),
			FirstName: contactMap["first_name"].(string),
			LastName:  contactMap["last_name"].(string),
			GenderId:  contactMap["gender_id"].(int),
			Dob:       contactMap["dob"].(time.Time),
			Email:     contactMap["email"].(string),
			Phone:     contactMap["phone"].(string),
			Address:   contactMap["address"].(string),
			PhotoPath: contactMap["photo_path"].(string),
			CreatedBy: "test-api",
			CreatedAt: time.Now(),
		})
	}

	contactRepo := contact.NewContactRepo()

	// Insert Contact to the database
	total, err := contactRepo.InsertContacts(contacts)
	fmt.Println("total", total)
	if err != nil {
		return nil, err
	}
	result:= models.Status {
		StatusID: 200,
		StatusText: "OK",
		Message: fmt.Sprintf("Total Items:%v", total),
	}

	return result, nil
}

func GetContactResolve(params graphql.ResolveParams) (interface{}, error) {
	// Update limit and offset if provided
	limit, ok := params.Args["limit"].(int)
	if !ok {
		limit = 10
	}

	offset, ok := params.Args["offset"].(int)
	if !ok {
		offset = 0
	}

	searchText, ok := params.Args["searchText"].(string)
	if !ok {
		searchText = ""
	}
	contactRepo := contact.NewContactRepo()

	// Fetch contacts from the database
	contacts, err := contactRepo.GetContactsBySearchText(searchText, limit, offset)
	if err != nil {
		return nil, err
	}
	return contacts, nil
}

func GetContactsPaginationResolve(params graphql.ResolveParams) (interface{}, error) {
	// Update limit and offset if provided
	page, ok := params.Args["page"].(int)
	if !ok {
		page = 1
	}

	pageSize, ok := params.Args["pageSize"].(int)
	if !ok {
		pageSize = 10
	}

	searchText, ok := params.Args["searchText"].(string)
	if !ok {
		searchText = ""
	}
	contactRepo := contact.NewContactRepo()

	// Fetch contacts from the database
	contacts, pager, err := contactRepo.GetContactsBySearchTextPagination(searchText, page, pageSize)
	var contactPagination = models.ContactPaginationModel{
		Contacts:   contacts,
		Pagination: pager,
	}

	if err != nil {
		return nil, err
	}
	return contactPagination, nil
}

func GetContactByIdResolve(params graphql.ResolveParams) (interface{}, error) {
	id := params.Args["id"].(int)
	contactRepo := contact.NewContactRepo()

	// Fetch contacts from the database
	contact, err := contactRepo.GetContactByID(id)
	if err != nil {
		return nil, err
	}
	return contact, nil
}

func UpdateContactResolve(params graphql.ResolveParams) (interface{}, error) {
	// Map input fields to Contact struct
	input := params.Args["input"].(map[string]interface{})
	invalids := validateContact(input)

	if len(invalids) > 0 {
		return nil, fmt.Errorf("%v", invalids)
	}

	contactInput := models.ContactModel{

		ContactId: int64(input["contact_id"].(int)),
		Name:      input["name"].(string),
		FirstName: input["first_name"].(string),
		LastName:  input["last_name"].(string),
		GenderId:  input["gender_id"].(int),
		Dob:       input["dob"].(time.Time),
		Email:     input["email"].(string),
		Phone:     input["phone"].(string),
		Address:   input["address"].(string),
		PhotoPath: input["photo_path"].(string),
		CreatedBy: "test-api",
		CreatedAt: time.Now(),
	}

	contactRepo := contact.NewContactRepo()

	// Update Contact to the database
	_, err := contactRepo.UpdateContact(&contactInput)
	if err != nil {
		return nil, err
	}

	return contactInput, nil
}

func DeleteContactResolve(params graphql.ResolveParams) (interface{}, error) {
	contact_id := params.Args["id"].(int)

	contactRepo := contact.NewContactRepo()

	// Delete Contact to the database
	_, err := contactRepo.DeleteContact(contact_id)
	if err != nil {
		return nil, err
	}
	result:= models.Status {
		StatusID: 200,
		StatusText: "OK",
		Message: fmt.Sprintf("Delete id:%v successfully.", contact_id),
	}

	return result, nil
}
