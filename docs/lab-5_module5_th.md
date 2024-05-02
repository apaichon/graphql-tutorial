# Module 5: Mutation
## Lab5.1 - Create the contact data.
**Objective:** เข้าใจการสร้างฟังก์ชั่น Mutation เพื่อเพิ่มข้อมูลในตาราง Contact ด้วย GraphQL
**ไฟล์ทีี่เกี่ยวข้องใน Lab นี้**
```plantuml

@startmindmap
* data
** event.db
* src
** graphql-api
*** cmd
**** server
***** .env
***** main.go
*** internal
**** contact
***** contact.repo.go
*** pkg
**** data
***** db.go
***** models
****** contact.model.go
**** graphql
***** resolvers
****** contact.resolver.go
***** agruments.go
***** schema.go
***** types.go
@endmindmap

```
1. เพิ่มฟังก์ชั่น *InsertContact* ในไฟล์ *internal/contact/contact.repo.go* ดังนี้
```go
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
```
2. เพื่ม *CreateContactInput* และ *CreateContactArgument* ในไฟล์ */pkg/graphql/agruments.go*
```go
var CreateContactInput =  graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "CreateContactsInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"name":  {Type: graphql.NewNonNull(graphql.String)},
		"first_name":  {Type: graphql.NewNonNull(graphql.String)},
		"last_name":  {Type: graphql.NewNonNull(graphql.String)},
		"gender_id":  {Type: graphql.NewNonNull(graphql.Int)},
		"dob":  {Type: graphql.NewNonNull(graphql.DateTime)},
		"email": {Type: graphql.NewNonNull(graphql.String)},
		"address": {Type: graphql.String},
		"phone": {Type: graphql.String},
		"photo_path": {Type: graphql.String},
	},
})

var CreateContactArgument = graphql.FieldConfigArgument{
	"input": &graphql.ArgumentConfig{
		Type: CreateContactInput,
	},
}
```
3. เพิ่มฟังก์ชั้น CreateContactResolve ใน pkg/graphql/resolvers/contact.resolver.go
```go
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
```
4. เพิ่ม *ContactMutation* ในไฟล์ *pkg/graphql/types.go*
```go
type ContactMutations struct {
	CreateContact func(map[string]interface{}) (*models.ContactModel, error) `json:"createContact"`
}
// Define the ContactMutations type
var ContactMutationsType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ContactMutations",
	Fields: graphql.Fields{
		"createContact": &graphql.Field{
			Type:    ContactGraphQLType,
			Args:    CreateContactArgument,
			Resolve: resolvers.CretateContactResolve,
		},
	},
})
```
5. เพื่มฟิลด์ *ContactMutations* ในไฟล์ *pkg/graphql/schema.go*
```go
// RootQuery represents the root GraphQL query.
var RootMutation = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "RootMutation",
		Fields: graphql.Fields{
			"contactMutations": &graphql.Field{
				Type: ContactMutationsType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return &ContactMutations{}, nil
				},
			},
			// Add other queries as needed
		},
	},
)
```
6. ทดสอบ insert ข้อมูลจากการส่ง GraphQL Query.

**Query**
```graphql
mutation CreateContact($input: CreateContactInput) {
  contactMutations {
    createContact(input: $input) {
      contact_id
    }
  }
}
```
**Variables**
```graphql
{
  "input": {
    "name": "Dr.PUP",
    "first_name": "PUP",
    "last_name": "Apaichon",
    "gender_id": 1,
    "dob": "1979-11-13T00:00:00Z",
    "email": "john@example.com",
    "phone": "123-456-7890",
    "address": "123 Main St",
    "photo_path": "path/to/photo.jpg"
  }
}
```


![SQL Tools](./images/lab5/figure5.1.png)
*Figure 5.1 Create Contact Mutation.*

## Lab5.2 - Create the contact data with Validation.
**Objective:** เข้าใจการเขียนฟังก์ชั่นสำหรับ Validation Input ก่อนส่งไปเพิ่มข้อมูลที่ฐานข้อมูล Mutation
**ไฟล์ทีี่เกี่ยวข้องใน Lab นี้**
```plantuml

@startmindmap
* data
** event.db
* src
** graphql-api
*** cmd
**** server
***** .env
***** main.go
*** internal
**** contact
***** contact.repo.go
*** pkg
**** data
***** db.go
***** models
****** contact.model.go
**** graphql
***** resolvers
****** contact.resolver.go
***** agruments.go
***** schema.go
***** types.go
@endmindmap

```

1. เพิ่ม Library github.com/go-playground/validator/v10 สำหรับ นำมาช่วย validate input ไปที่ root project.
```sh
go get github.com/go-playground/validator/v10
```
2. เพิ่มฟังก์ชั่น validateContact และแทรกฟังก์ชั่นก่อนทำการ insert ข้อมูล

```go
import (
	"fmt"
	"graphql-api/internal/contact"
	"graphql-api/pkg/data/models"
	"time"

	"github.com/go-playground/validator/v10" // add validator
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

func CreateContactResolve(params graphql.ResolveParams) (interface{}, error) {
	// Map input fields to Contact struct
	input := params.Args["input"].(map[string]interface{})
	invalids := validateContact(input) // Add in CreateContactResolve.

	if len(invalids) > 0 { // Check error validation.
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
```
3. ทดสอบลอง insert ข้อมูลจาก GraphQL ลองไม่ป้อนค่าใน name, first_name, last_name หรือ email ผิด format

**Query**
```graphql
mutation CreateContact($input: CreateContactInput) {
  contactMutations {
    createContact(input: $input) {
      contact_id
    }
  }
}
```
**Variables**
```graphql
{
  "input": {
    "name": "PUP",
    "first_name": "PUP",
    "last_name": "Apaichon",
    "gender_id": 1,
    "dob": "1979-11-13T00:00:00Z",
    "email": "john@example.com",
    "phone": "123-456-7890",
    "address": "123 Main St",
    "photo_path": "path/to/photo.jpg"
  }
}
```
![SQL Tools](./images/lab5/figure5.2.png)
*Figure 5.2 Create Contact Mutation with validation.*


## Lab5.3 - Update the contact data.
**Objective:** เข้าใจการเขียนฟังก์ชั่นสำหรับ แก้ไขข้อมูลในฐานข้อมูลด้วย Mutation
**ไฟล์ทีี่เกี่ยวข้องใน Lab นี้**

```plantuml
@startmindmap
* data
** event.db
* src
** graphql-api
*** cmd
**** server
***** .env
***** main.go
*** internal
**** contact
***** contact.repo.go
*** pkg
**** data
***** db.go
***** models
****** contact.model.go
**** graphql
***** resolvers
****** contact.resolver.go
***** agruments.go
***** schema.go
***** types.go
@endmindmap

```

1. เพิ่มฟังก์ชั่น *UpdateContact* ในไฟล์ *internal/contact/contact.repo.go* ดังนี้
```go
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
```
2. เพื่ม *UpdateContactInput* และ *CreateContactArgument* ในไฟล์ */pkg/graphql/agruments.go*
```go
var UpdateContactArgument = graphql.FieldConfigArgument{
	"input": &graphql.ArgumentConfig{
		Type: graphql.NewInputObject(graphql.InputObjectConfig{
			Name: "UpdateContactInput",
			Fields: graphql.InputObjectConfigFieldMap{
				"contact_id":  {Type: graphql.NewNonNull(graphql.Int)},
				"name":  {Type: graphql.NewNonNull(graphql.String)},
				"first_name":  {Type: graphql.NewNonNull(graphql.String)},
				"last_name":  {Type: graphql.NewNonNull(graphql.String)},
				"gender_id":  {Type: graphql.NewNonNull(graphql.Int)},
				"dob":  {Type: graphql.NewNonNull(graphql.DateTime)},
				"email": {Type: graphql.NewNonNull(graphql.String)},
				"address": {Type: graphql.String},
				"phone": {Type: graphql.String},
				"photo_path": {Type: graphql.String},
			},
		}),
	},
}
```
3. เพิ่มฟังก์ชั้น UpdateContactResolve ใน pkg/graphql/resolvers/contact.resolver.go
```go
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
```
4. เพิ่ม *ContactMutation* ในไฟล์ *pkg/graphql/types.go*
```go
type ContactMutations struct {
	CreateContact func(map[string]interface{}) (*models.ContactModel, error) `json:"createContact"`
	UpdateContact func(map[string]interface{}) (*models.ContactModel, error) `json:"updateContact"`	
}

// Define the ContactMutations type
var ContactMutationsType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ContactMutations",
	Fields: graphql.Fields{
		"createContact": &graphql.Field{
			Type:    ContactGraphQLType,
			Args:    CreateContactArgument,
			Resolve: resolvers.CretateContactResolve,
		},
		"updateContact": &graphql.Field{
			Type:    ContactGraphQLType,
			Args:    UpdateContactArgument,
			Resolve: resolvers.UpdateContactResolve,
		},
	},
})

```
5. ทดสอบ insert ข้อมูลจากการส่ง GraphQL Query.

**Query**
```graphql
mutation UpdateContact($input: UpdateContactInput) {
  contactMutations {
 updateContact(input: $input) {
    contact_id
    name
    first_name
    }
  }
}
```
**Variables**
```graphql
{
  "input": {
    "contact_id": 5,
    "name": "Apaichon555",
    "first_name": "John555",
    "last_name": "Doe555",
    "gender_id": 1,
    "dob": "1990-01-01T00:00:00Z",
    "email": "john@example.com",
    "phone": "123-456-7890",
    "address": "123 Main St",
    "photo_path": "path/to/photox.jpg"
  }
}
```


## Lab5.4 - Delete the contact data.
**Objective:** เข้าใจการเขียนฟังก์ชั่นสำหรับ ลบข้อมูลในฐานข้อมูลด้วย Mutation
**ไฟล์ทีี่เกี่ยวข้องใน Lab นี้**
```plantuml
@startmindmap
* data
** event.db
* src
** graphql-api
*** cmd
**** server
***** .env
***** main.go
*** internal
**** contact
***** contact.repo.go
*** pkg
**** data
***** db.go
***** models
****** status.model.go
**** graphql
***** resolvers
****** contact.resolver.go
***** agruments.go
***** schema.go
***** types.go
@endmindmap

```

1. เพิ่มฟังก์ชั่น *DeleteContact* ในไฟล์ *internal/contact/contact.repo.go* ดังนี้
```go
func DeleteContactResolve(params graphql.ResolveParams) (interface{}, error) {
	// Map input fields to Contact struct
	contact_id := params.Args["id"].(int)

	contactRepo := contact.NewContactRepo()

	// Insert Contact to the database
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
```
2. เพิ่มฟังก์ชั้น DeleteContactResolve ใน pkg/graphql/resolvers/contact.resolver.go
```go
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
```
3. เพิ่ม *ContactMutation* ในไฟล์ *pkg/graphql/types.go*
```go
type ContactMutations struct {
	CreateContact func(map[string]interface{}) (*models.ContactModel, error) `json:"createContact"`
	UpdateContact func(map[string]interface{}) (*models.ContactModel, error) `json:"updateContact"`
	DeleteContact func(int) (*models.Status, error) `json:"deleteContact"`		
}

// Define the ContactMutations type
var ContactMutationsType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ContactMutations",
	Fields: graphql.Fields{
		"createContact": &graphql.Field{
			Type:    ContactGraphQLType,
			Args:    CreateContactArgument,
			Resolve: resolvers.CretateContactResolve,
		},
		"updateContact": &graphql.Field{
			Type:    ContactGraphQLType,
			Args:    UpdateContactArgument,
			Resolve: resolvers.UpdateContactResolve,
		},
		"deleteContact": &graphql.Field{
			Type:    StatusGraphQLType,
			Args:    IdArgument,
			Resolve: resolvers.DeleteContactResolve,
		},
	},
})

```
4. ทดสอบ insert ข้อมูลจากการส่ง GraphQL Query.

**Query**
```graphql
mutation {
  contactMutations {
    deleteContact(id: 6) {
      status_id
      status
      message
    }
  }
}
```

## Lab5.5 - Multiple insert data with exist mutation.
**Objective:** เข้าใจการเขียน GraphQuery Mutation เพื่อ insert ข้อมูลหลายๆรายการจาก  Mutation ที่มีอยู่แล้ว
**ไฟล์ทีี่เกี่ยวข้องใน Lab นี้**
ใช้ไฟล์กลุ่มเดิมจากตัวอย่าง insert

**Query**
```graphql
mutation CreateContact($input: CreateContactInput, $input2: CreateContactInput) {
  contactMutations {
    contact1: createContact(input: $input) {
      contact_id
    }
    contact2: createContact(input: $input2) {
      contact_id
    }
  }
}
```

**Variables**
```graphql
{
  "input": {
    "name": "",
    "first_name": "PUP",
    "last_name": "Apaichon",
    "gender_id": 1,
    "dob": "1979-11-13T00:00:00Z",
    "email": "john@example.com",
    "phone": "123-456-7890",
    "address": "123 Main St",
    "photo_path": "path/to/photo.jpg"
  },
  "input2": {
    "name": "Dr.PUP2",
    "first_name": "PUP2",
    "last_name": "Apaichon2",
    "gender_id": 1,
    "dob": "1979-11-13T00:00:00Z",
    "email": "john@example.com",
    "phone": "123-456-7890",
    "address": "123 Main St",
    "photo_path": "path/to/photo.jpg"
  }
}
```


## Lab5.6 - Multiple insert data with best performance for database.
**Objective:** เข้าใจการเขียน GraphQuery Mutation เพื่อ insert ข้อมูลหลายๆรายการ ด้วยการเพิ่ม Mutation ให้รอบรับสำหรับการเพิ่มข้อมูลปริมาณมากที่เหมาะกับการทำงานแบบ Batch ของฐานข้อมูล
**ไฟล์ทีี่เกี่ยวข้องใน Lab นี้**
ใช้ไฟล์กลุ่มเดิมจากตัวอย่าง insert เพิ่มฟังก์ชั่นในส่วนต่างๆ

1. เพิ่มฟังก์ชั่น *InsertContacts* ในไฟล์ *internal/contact/contact.repo.go* ดังนี้
```go
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
```
2. เพื่ม *CreateContactsArgument* ในไฟล์ */pkg/graphql/agruments.go*
```go
var CreateContactsArgument = graphql.FieldConfigArgument{
	"contacts": &graphql.ArgumentConfig{
		Type: graphql.NewNonNull(graphql.NewList(CreateContactInput)),
	},
}

```
3. เพิ่มฟังก์ชั้น CreateContactsResolve ใน pkg/graphql/resolvers/contact.resolver.go
```go
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
```

4. เพิ่ม CreateContacts ใน *ContactMutation* ในไฟล์ *pkg/graphql/types.go*
```go
type ContactMutations struct {
	CreateContact func(map[string]interface{}) (*models.ContactModel, error) `json:"createContact"`
	CreateContacts func(map[string]interface{}) ([]*models.ContactModel, error) `json:"createContacts"`
	UpdateContact func(map[string]interface{}) (*models.ContactModel, error) `json:"updateContact"`
	DeleteContact func(int) (*models.Status, error) `json:"deleteContact"`		
}

// Define the ContactMutations type
var ContactMutationsType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ContactMutations",
	Fields: graphql.Fields{
		"createContact": &graphql.Field{
			Type:    ContactGraphQLType,
			Args:    CreateContactArgument,
			Resolve: resolvers.CretateContactResolve,
		},
		"createContacts": &graphql.Field{
			Type:    StatusGraphQLType,
			Args:    CreateContactsArgument,
			Resolve: resolvers.CreateContactsResolve,
		},
		"updateContact": &graphql.Field{
			Type:    ContactGraphQLType,
			Args:    UpdateContactArgument,
			Resolve: resolvers.UpdateContactResolve,
		},
		"deleteContact": &graphql.Field{
			Type:    StatusGraphQLType,
			Args:    IdArgument,
			Resolve: resolvers.DeleteContactResolve,
		},
	},
})

```
5. ทดสอบ insert ข้อมูลจากการส่ง GraphQL Query.

**Query**
```graphql
mutation CreateContact($contacts: [CreateContactInput]!) {
  contactMutations {
    createContacts(contacts: $contacts) {
      message
    }
  }
}
```

**Variables**
```graphql
{
  "contacts": [
    {
      "name": "Dr.PUP13",
      "first_name": "PUP13",
      "last_name": "Apaichon",
      "gender_id": 1,
      "dob": "1979-11-13T00:00:00Z",
      "email": "john@example.com",
      "phone": "123-456-7890",
      "address": "123 Main St",
      "photo_path": "path/to/photo.jpg"
    },
    {
      "name": "Dr.PUP11",
      "first_name": "PUP11",
      "last_name": "Apaichon2",
      "gender_id": 1,
      "dob": "1979-11-13T00:00:00Z",
      "email": "john@example.com",
      "phone": "123-456-7890",
      "address": "123 Main St",
      "photo_path": "path/to/photo.jpg"
    }
  ]
}
```
