# Module 6: Authentication and Authorization
## Lab6.1 - Authentication with JWT
**Objective:** Understand how to create a login function with JWT and use it to control Api calls for basic security.
**Related files in this lab**
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
**** auth
***** authen.go
***** cryptography.go
***** genpass.go
***** user.repo.go
*** pkg
**** data
***** db.go
***** models
****** jwt-claim.model.go
****** jwt-token.model.go
****** user.model.go
@endmindmap

```
1. Create a user table database

```sql
CREATE TABLE IF NOT EXISTS user (
    user_id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_name TEXT NOT NULL,
    password TEXT NOT NULL,
    salt TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    created_by TEXT NOT NULL
);

```

2. Prepare functions for creating user and password. Prepare files as follows.

- *pkg/models/jwt-claim.model.go*

```go 
package models
import("github.com/dgrijalva/jwt-go")

type JwtClaims struct {
	UserId int `json:"user_id"`
	Username string `json:"user_name"`
	jwt.StandardClaims
}
```
- *pkg/models/jwt-token.model.go*

```go
package models

type JwtToken struct {
	Token string `json:"token"`
	ExpiredAt int64 `json:"expiredAt"`
	
}

```
- *pkg/models/user.model.go*

```go
package models
import("time")

type User struct {
	UserId int `json:"user_id"`
	Username string `json:"user_name"`
	Password string `json:"password"`
	Salt string `json:"salt"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string  `json:"created_by"`
}

```

- *internal/auth*

```go
// authen.go
package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"graphql-api/config"
	"graphql-api/pkg/data/models"
	"github.com/dgrijalva/jwt-go"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userRepo := NewUserRepo()
	config := config.NewConfig()

	userdb, ok := userRepo.GetUserByName(user.Username) // users[user.Username]
	if ok != nil {
		http.Error(w, "Invalid username", http.StatusUnauthorized)
		return
	}
	// fmt.Println("User")
	hash := HashString(user.Username + user.Password + userdb.Salt)
	// fmt.Println("UserDbPassword:", userdb.Password, "Hash:",hash, "UserNameDb:", userdb.Username, "username", user.Username )
	if userdb.Password != hash {
		http.Error(w, "Invalid username", http.StatusUnauthorized)
		return
	}

	// fmt.Println("Duration", time.Duration(config.TokenAge))
	expirationTime := time.Now().Add(time.Duration(config.TokenAge) * time.Minute)

	claims := &models.JwtClaims{
		UserId:  user.UserId,
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.SecretKey))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	response := models.JwtToken{
		Token:     tokenString,
		ExpiredAt: expirationTime.Unix(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func AuthenticationHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the JWT token is present in the request
		tokenString := getTokenFromRequest(r)
		// fmt.Println("tokenString:" + tokenString)
		if tokenString == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Validate the JWT token
		token, err := validateToken(tokenString)
		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Pass the request to the next handler if the token is valid
		next.ServeHTTP(w, r)
	})
}

func getTokenFromRequest(r *http.Request) string {
	// Check if the token is present in the request header
	token := r.Header.Get("Authorization")
	if token != "" {
		return token
	}

	// Check if the token is present in the request cookies
	cookie, err := r.Cookie("token")
	if err == nil {
		return cookie.Value
	}

	return ""
}

func validateToken(tokenString string) (*jwt.Token, error) {
	config := config.NewConfig()
	// Parse and validate the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.SecretKey), nil // Replace "your-secret-key" with your actual secret key
	})
	// fmt.Println("Token", token)
	if err != nil {
		return nil, err
	}

	return token, nil
}

```

```go
// cryptography.go
package auth

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	
)

func GenerateRandomSalt(length int) string {
	// Define the character set from which to generate the salt
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// Seed the random number generator
	//rand.Seed(time.Now().UnixNano())

	// Create a byte slice to store the salt characters
	salt := make([]byte, length)

	// Fill the byte slice with random characters from the charset
	for i := 0; i < length; i++ {
		salt[i] = charset[rand.Intn(len(charset))]
	}

	// Convert the byte slice to a string and return it
	return string(salt)
}

func HashString(input string) string {

	// Calculate the MD5 hash of the concatenated string
	hasher := md5.New()
	hasher.Write([]byte(input))
	hash := hasher.Sum(nil)

	// Convert the hash to a hexadecimal string
	hashString := hex.EncodeToString(hash)

	return hashString
}

```

```go
// user.repo.go
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
func (cr *UserRepo) GetUsersBySearchText(searchText string, limit, offset int) ([]*models.User, error) {
	var users []*models.User

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
		var user models.User
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
func (cr *UserRepo) GetUserByID(id int) (*models.User, error) {
	var user models.User
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
func (cr *UserRepo) GetUserByName(name string) (*models.User, error) {
	var user models.User
	// Execute query to get a user by ID from the database
	row, err := cr.DB.QueryRow("SELECT * FROM user WHERE user_name = ?", name)

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

// Insert User inserts a new user into the database
func (cr *UserRepo) InsertUser(user *models.User) (int64, error) {
	// Execute insert query to insert a new user into the database
	result, err := cr.DB.Insert("INSERT INTO user (user_id,user_name,password,salt,created_at) VALUES ({?,?,?,?,?})",
		user.UserId, user.Username, user.Password, user.Salt, user.CreatedAt)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// Update User updates an existing user in the database
func (cr *UserRepo) UpdateUser(user *models.User) (int64, error) {
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

```

- *pkg/genpass/main.go*

```go
package main

import (
	"fmt"
	"os"
	"graphql-api/internal/auth"
)

func main() {

	fmt.Println(os.Args)
	if len(os.Args) < 3 {
		fmt.Println("Please provide user and password")
		return
	}

	user := os.Args[1]
	password := os.Args[2]
	fmt.Printf("User :%s, Password: %s", user, password)

	salt := auth.GenerateRandomSalt(8)
	userhash := user+password +  salt

	fmt.Printf("\n Salt:%s, Hash:%s", salt, auth.HashString(userhash) )
}

```
3. Create user and password, run Command Line

```sh
# pkg/genpass/main.go
	go run main.go <user> <password>
```

4. Add user information in the database user_name, password, salt

```sql
INSERT INTO user (user_name, password, salt, created_by)
VALUES ('admin', '051504d937d4351e2c5f72a84bb9944f', 'tBwTPvHD', 'system');
```
5.Create an api for login at cmd/server/main.go file.

```go
import (
	"fmt"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"graphql-api/config"
	gql "graphql-api/pkg/graphql"
	"graphql-api/internal/auth" // Add library

)

http.HandleFunc("/login", auth.LoginHandler) // auth.LoginHandler
```

6. Make a High Order Function and cover it with the AuthenHandler function.

```go
http.Handle("/graphql", auth.AuthenticationHandler( graphqlHandler)) // Add auth.AuthenticationHandler
```

7. Run a test project Authentication with Thunder Client
 - Login
 **Request**
 ```json
{
  "user_name": "admin",
  "password": "P@ssw0rd"
}
 ```
 **Response**
 
 ```json
 {
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjowLCJ1c2VyX25hbWUiOiJhZG1pbiIsImV4cCI6MTcxNDcyOTM3M30.Gjy7R2lHK1QoV578CqC98sn6DwPFzKEY3pspLbuqSik",
  "expiredAt": 1714729373
}
 ```

 8. Test the function in GraphQL before inserting the token in the Header and after inserting the Authentication Header.


## Lab6.2 - Authorization with Role and Permission of Database
**Objective:** Set access rights at the function level
**Related files in this lab**
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
**** auth
***** authorize.go
***** role.repo.go
***** user.repo.go
*** pkg
**** data
***** db.go
***** models
****** jwt-claim.model.go
****** jwt-token.model.go
****** role-permission-view.model.go
@endmindmap
```

**ER-Diagram**

```plantuml
@startuml
' Entity-Relationship Diagram for Database Schema

' Table: User
entity "user" {
    + user_id : INTEGER (PK)  ' Primary Key
    -- 
    + user_name : TEXT
    + password : TEXT
    + salt : TEXT
    + created_at : DATETIME (Default: CURRENT_TIMESTAMP)
    + created_by : TEXT
}

' Table: Role
entity "role" {
    + role_id : INTEGER (PK)  ' Primary Key
    -- 
    + role_name : TEXT
    + role_desc : TEXT
    + is_super_admin : BOOLEAN
    + created_at : DATETIME
    + created_by : TEXT
    + status_id : INTEGER
}

' Table: User Roles
entity "user_roles" {
    + user_role_id : INTEGER (PK)  ' Primary Key
    -- 
    + role_id : INTEGER  ' Foreign Key
    + user_id : INTEGER  ' Foreign Key
    + created_at : DATETIME
    + created_by : TEXT
}

' Table: Role Permissions
entity "role_permissions" {
    + role_permission_id : INTEGER (PK)  ' Primary Key
    -- 
    + role_id : INTEGER  ' Foreign Key
    + resource_type_id : INTEGER  ' 1=GraphQLResolve, 2=GraphQLField, 3=API Function, 4=Web Menu, 5=Mobile Menu
    + resource_name : TEXT
    + can_execute : BOOLEAN
    + can_read : BOOLEAN
    + can_write : BOOLEAN
    + can_delete : BOOLEAN
    + created_at : DATETIME
    + created_by : TEXT
}

' Relationships
"user_roles" o-- "user" : References user_id
"user_roles" o-- "role" : References role_id
"role_permissions" o-- "role" : References role_id

@enduml

```

1. Create a database to store Roles, UserRoles, and RolePermissions.
```sql
CREATE TABLE IF NOT EXISTS role (
    role_id INTEGER PRIMARY KEY AUTOINCREMENT,
    role_name TEXT NOT NULL,
	role_desc TEXT,
	is_super_admin BOOLEAN,
	created_at DATETIME,
	created_by TEXT NOT NULL,
	status_id INTEGER
);


CREATE TABLE IF NOT EXISTS user_roles (
    user_role_id INTEGER PRIMARY KEY AUTOINCREMENT,
    role_id INTEGER NOT NULL,
	user_id INTEGER,
	created_at DATETIME,
	created_by TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS role_permissions (
    role_permission_id INTEGER PRIMARY KEY AUTOINCREMENT,
    role_id INTEGER,
    resource_type_id INTEGER, -- 1 GraphQLResolve, 2 GraphQLField, 3 API Function, 4 Web Menu, 5 Mobile Menu
	resource_name TEXT,
    can_execute BOOLEAN,
    can_read BOOLEAN,
    can_write BOOLEAN,
    can_delete BOOLEAN,
	created_at DATETIME,
	created_by TEXT NOT NULL
);

```

​2. Create the user you want to test from Command Line and SQL.
```sh
# pkg/genpass/main.go
go run main.go <user> <password>

```
**SQL**
```sql
INSERT INTO user (user_name, password, salt, created_by)
VALUES ('<user>', '<hash>', '<salt>', 'system');

```

3. Prepare Data Role, UserRole, Role Permission in Database.
```sql
CREATE TABLE IF NOT EXISTS role (
    role_id INTEGER PRIMARY KEY AUTOINCREMENT,
    role_name TEXT NOT NULL,
	role_desc TEXT,
	is_super_admin BOOLEAN,
	created_at DATETIME,
	created_by TEXT NOT NULL,
	status_id INTEGER
);


CREATE TABLE IF NOT EXISTS user_roles (
    user_role_id INTEGER PRIMARY KEY AUTOINCREMENT,
    role_id INTEGER NOT NULL,
	user_id INTEGER,
	created_at DATETIME,
	created_by TEXT NOT NULL
);
 

CREATE TABLE IF NOT EXISTS role_permissions (
    role_permission_id INTEGER PRIMARY KEY AUTOINCREMENT,
    role_id INTEGER,
    resource_type_id INTEGER, -- 1 GraphQLResolve, 2 GraphQLField, 3 API Function, 4 Web Menu, 5 Mobile Menu
	resource_name TEXT,
    can_execute BOOLEAN,
    can_read BOOLEAN,
    can_write BOOLEAN,
    can_delete BOOLEAN,
	created_at DATETIME,
	created_by TEXT NOT NULL
);

INSERT INTO role (role_name, role_desc, is_super_admin, created_at, created_by, status_id)
VALUES
    ('Admin', 'Administration role with full access', 1, '2023-05-05 10:00:00', 'system', 1),
    ('Manager', 'Manager role with limited access', 0, '2023-05-05 10:01:00', 'system', 1),
    ('User', 'Regular user role with read-only access', 0, '2023-05-05 10:02:00', 'system', 1);

INSERT INTO user_roles (role_id, user_id, created_at, created_by)
VALUES
    (1, 1, '2023-05-05 10:10:00', 'system'),
    (2, 2, '2023-05-05 10:11:00', 'system')
  

    INSERT INTO role_permissions (role_id, resource_type_id, resource_name, can_execute, can_read, can_write, can_delete, created_at, created_by)
    VALUES
   
    (2, 1, 'contacts.gets', 1, 1, 0, 0, '2023-05-05 10:23:00', 'system'),
    (2, 1, 'contacts.getById', 1, 1, 0, 0, '2023-05-05 10:24:00', 'system'),
    (2, 1, 'contacts.getPagination', 0, 1, 0, 0, '2023-05-05 10:26:00', 'system'),
    (2, 1, 'contactMutations.createContact', 1, 1, 0, 0, '2023-05-05 10:25:00', 'system'),
    (2, 1, 'contactMutations.updateContact', 0, 1, 0, 0, '2023-05-05 10:25:00', 'system'),
    (2, 1, 'contactMutations.deleteContact', 0, 1, 0, 0, '2023-05-05 10:25:00', 'system')
   
CREATE VIEW vw_user_role_permissions AS
SELECT 
    u.user_id,
    u.user_name,
    r.role_id,
    r.role_name,
    r.role_desc,
    r.is_super_admin,
    ur.user_role_id,
    ur.created_at AS user_role_created_at,
    rp.role_permission_id,
    rp.resource_type_id,
    rp.resource_name,
    rp.can_execute,
    rp.can_read,
    rp.can_write,
    rp.can_delete,
    rp.created_at AS role_permission_created_at
FROM
    user_roles ur
    JOIN user u ON ur.user_id = u.user_id
    JOIN role r ON ur.role_id = r.role_id
    LEFT JOIN role_permissions rp ON r.role_id = rp.role_id
```
4. Create role-permission-view.model.go for Binding with RolePermissions data from the database.
```go
package models

import "time"

// UserRolePermissions represents the vw_user_role_permissions view
type UserRolePermissionsViewModel struct {
	UserID                    int       `json:"user_id"`
	UserName                  string    `json:"user_name"`
	RoleID                    int       `json:"role_id"`
	RoleName                  string    `json:"role_name"`
	RoleDesc                  *string   `json:"role_desc,omitempty"` // Nullable field
	IsSuperAdmin              bool      `json:"is_super_admin"`
	UserRoleID                int       `json:"user_role_id"`
	UserRoleCreatedAt         time.Time `json:"user_role_created_at"`
	RolePermissionID          *int      `json:"role_permission_id,omitempty"` // Nullable field
	ResourceTypeID            *int      `json:"resource_type_id,omitempty"`  // Nullable field
	ResolveName               *string   `json:"resolve_name,omitempty"` // Nullable field
	CanExecute                *bool     `json:"can_execute,omitempty"` // Nullable field
	CanRead                   *bool     `json:"can_read,omitempty"`   // Nullable field
	CanWrite                  *bool     `json:"can_write,omitempty"`  // Nullable field
	CanDelete                 *bool     `json:"can_delete,omitempty"` // Nullable field
	RolePermissionCreatedAt   *time.Time `json:"role_permission_created_at,omitempty"` // Nullable field
}

// UserRolePermissions represents the vw_user_role_permissions view
type UserRoleResolvePermissionsModel struct {
	UserID                    int       `json:"user_id"`
	IsSuperAdmin              bool      `json:"is_super_admin"`
	ResolveName               *string   `json:"resolve_name,omitempty"` // Nullable field
	CanExecute                *bool     `json:"can_execute,omitempty"` // Nullable field
}
```
5. Add function in internal/auth/role.repo.go
```go
package auth

import (
	"time"
	"fmt"
	"graphql-api/pkg/data"
	"graphql-api/pkg/data/models"
	
)

// RoleRepo represents the repository for role operations
type RoleRepo struct {
	DB *data.DB
}

// NewRoleRepo creates a new instance of RoleRepo
func NewRoleRepo() *RoleRepo {
	db := data.NewDB()
	return &RoleRepo{DB: db}
}

// GetUserRoleResolveByID retrieves a user role by its ID
func (rr *RoleRepo) GetUserIsSuperAdminByUserID(userID int) (bool, error) {
	row, err := rr.DB.QueryRow(`select sum(is_super_admin) as total_super_admin
	from vw_user_role_permissions 
	where user_id =?`, userID)
	
	isSuperAdmin :=false
	if err != nil {
		return isSuperAdmin, err
	}
	
	totalSuperAdmin := 0
	
	
	err = row.Scan(
		&totalSuperAdmin,
	)

	if err != nil {
		return isSuperAdmin, err
	}

	if totalSuperAdmin >0 {
		isSuperAdmin =true
	}

	return isSuperAdmin, nil
}

// GetUserRoleResolveByID retrieves a user role by its ID
func (rr *RoleRepo) GetUserRoleResolvePermissionByUserID(userID int, resolveName string) (*models.UserRoleResolvePermissionsModel, error) {
	row, err := rr.DB.QueryRow(`select user_id, resource_name, sum(can_execute) as can_execute ,
	sum(is_super_admin) as is_super_admin
	from vw_user_role_permissions 
	where user_id =? and resource_name =? and resource_type_id =1
	group by user_id, resource_name`, userID, resolveName)
	
	if err != nil {
		return nil, err
	}
	
	var userRoleResolve models.UserRoleResolvePermissionsModel
	err = row.Scan(
		&userRoleResolve.UserID,
		&userRoleResolve.ResolveName,
		&userRoleResolve.CanExecute,
		&userRoleResolve.IsSuperAdmin,
	)

	if err != nil {
		return nil, err
	}
	return &userRoleResolve, nil
}

```
6. Create /internal/auth/authorize.go for Set Authorize to Resolver.
```go
package auth

import (
	"fmt"
	"errors"
	"graphql-api/config"
	"graphql-api/pkg/data/models"
	"github.com/graphql-go/graphql"
)

var roleRepo RoleRepo
var appConfig config.Config
type ContextKey string
const userKey = ContextKey("user")

func init() {
	roleRepo = *NewRoleRepo()
	appConfig = *config.NewConfig()
}

type AuthorizeWorkflow struct {
	errors []error     // List of errors encountered during the workflow
	result interface{} // The result of the workflow operations
}

func (auth *AuthorizeWorkflow) IsSuperAdmin(userId int) *AuthorizeWorkflow {
	// roleRepo := NewRoleRepo()
	isSuperAdmin, err := roleRepo.GetUserIsSuperAdminByUserID(userId)
	if err != nil {
		auth.addError(err)
	}
	auth.setResult(isSuperAdmin)
	return auth
}

func (auth *AuthorizeWorkflow) GetUserIDFromToken( p graphql.ResolveParams) *AuthorizeWorkflow {
	userKey := ContextKey("user")
	tokenString, _ := p.Context.Value(userKey).(string)
	claims, err := DecodeJWTToken(tokenString, appConfig.SecretKey)

	if err != nil {
		auth.addError(errors.New("token expired"))
	}
	auth.setResult(claims)
	return auth
}

func (auth *AuthorizeWorkflow) GetRosolvePermission(userId int, resolveName string) *AuthorizeWorkflow {
	permission, err := roleRepo.GetUserRoleResolvePermissionByUserID(userId, resolveName)

	if err != nil {
		auth.addError(fmt.Errorf("unauthorized: missing %s permission. error:%s", resolveName, err))
	}

	canExecute :=false
	if  permission.CanExecute != nil {
		canExecute = *permission.CanExecute
	}
	auth.setResult(canExecute)
	
	return auth
}

func (auth *AuthorizeWorkflow) GetResult() interface{} {
	return auth.result
}

func (auth *AuthorizeWorkflow) GetError() interface{} {
	joinedErr := errors.Join(auth.errors...)
	return joinedErr
}

func (auth *AuthorizeWorkflow) addError(err error) *AuthorizeWorkflow {
	if err != nil {
		auth.errors = append(auth.errors, err)
	}
	return auth
}

func (auth *AuthorizeWorkflow) setResult(res interface{}) *AuthorizeWorkflow {
	auth.result = res
	return auth
}


// Middleware to enforce authorization based on permission
func AuthorizeResolver(resovleName string, next func(p graphql.ResolveParams) (interface{}, error)) func(p graphql.ResolveParams) (interface{}, error) {
	return func(p graphql.ResolveParams) (interface{}, error) {

		tokenString, _ := p.Context.Value(userKey).(string)
		config := config.NewConfig()
		claims, err := DecodeJWTToken(tokenString, config.SecretKey)

		if err != nil {
			fmt.Printf("\nerror:%s", err)
			return nil, errors.New("token expired")
		}

		// Check if the user has the required permission
		roleRepo := NewRoleRepo()
		isSuperAdmin, err := roleRepo.GetUserIsSuperAdminByUserID(claims.UserId)
		
		if err != nil {
			return nil, err
		}

		if isSuperAdmin {
			return next(p)
		}

		permission, err := roleRepo.GetUserRoleResolvePermissionByUserID(claims.UserId, resovleName)

		if err != nil {
			return nil, fmt.Errorf("unauthorized: missing %s permission. error:%s", resovleName, err)
		}

		canExecute :=false
		if  permission.CanExecute != nil {
			canExecute = *permission.CanExecute
		}

		if !permission.IsSuperAdmin && !canExecute {
			return nil, fmt.Errorf("unauthorized: missing %s permission", resovleName)
		}

		// Execute the resolver if permission is granted
		return next(p)
	}
}


// Middleware to enforce authorization based on permission
func AuthorizeResolverClean(resovleName string, next func(p graphql.ResolveParams) (interface{}, error)) func(p graphql.ResolveParams) (interface{}, error) {
	return func(p graphql.ResolveParams) (interface{}, error) {

		wf := &AuthorizeWorkflow{}
		user := wf.GetUserIDFromToken(p).GetResult().(*models.JwtClaims)
		isSuperAdmin := wf.IsSuperAdmin(user.UserId).GetResult().(bool)

		if isSuperAdmin {
			return next(p)
		}

		canExecute := wf.GetRosolvePermission(user.UserId, resovleName).GetResult().(bool)

		if canExecute {
			return next(p)
		}
		wf.addError(fmt.Errorf("unauthorized: missing %s permission", resovleName))

		errors := wf.GetError().(error)

		if errors !=nil {
			return nil, errors
		}
		
		return next(p)
	}
}

```

7. Put the AuthorizeResolver function in Resolve of pkg/graphql.types.go in the Query and Mutation section to check permissions.
```go
// Define the ContactQueries type
var ContactQueriesType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ContactQueries",
	Fields: graphql.Fields{
		"gets": &graphql.Field{
			Type:    graphql.NewList(ContactGraphQLType),
			Args:    SearhTextQueryArgument,
			Resolve: auth.AuthorizeResolverClean("contacts.gets", resolvers.GetContactResolve),
		},
		"getPagination": &graphql.Field{
			Type:    ContactPaginationGraphQLType,
			Args:    SearhTextPaginationQueryArgument,
			Resolve: auth.AuthorizeResolverClean("contacts.getPagination",resolvers.GetContactsPaginationResolve),
		},
		"getById": &graphql.Field{
			Type:    ContactGraphQLType,
			Args:    IdArgument,
			Resolve: auth.AuthorizeResolverClean("contacts.getById",resolvers.GetContactByIdResolve),
		},

	},
})

// Define the ContactMutations type
var ContactMutationsType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ContactMutations",
	Fields: graphql.Fields{
		"createContact": &graphql.Field{
			Type:    ContactGraphQLType,
			Args:    CreateContactArgument,
			Resolve: auth.AuthorizeResolver("contactMutations.createContact",resolvers.CretateContactResolve),
		},
		"createContacts": &graphql.Field{
			Type:    StatusGraphQLType,
			Args:    CreateContactsArgument,
			Resolve: resolvers.CreateContactsResolve,
		},
		"updateContact": &graphql.Field{
			Type:    ContactGraphQLType,
			Args:    UpdateContactArgument,
			Resolve: auth.AuthorizeResolver("contactMutations.updateContact",resolvers.UpdateContactResolve),
		},
		"deleteContact": &graphql.Field{
			Type:    StatusGraphQLType,
			Args:    IdArgument,
			Resolve: auth.AuthorizeResolver("contactMutations.deleteContact",resolvers.DeleteContactResolve),
		},
	},
})

```

8. Test Permission from Query and Mutation

```graphql
{
  contacts {
    gets(searchText: "") {
      name
      first_name
      last_name
      gender_id
      dob
      email
      phone
      address
      photo_path
      created_at
      created_by
    }
    getById(id: 1) {
      name
      first_name
      last_name
      gender_id
      dob
      email
      phone
      address
      photo_path
      created_at
      created_by
    }
    getPagination(searchText: "", page: 1, pageSize: 50) {
      contacts {
        contact_id
        name
        first_name
        last_name
        gender_id
      }
      pagination {
        page
        totalItems
        totalPages
        hasNext
      }
    }
  }
}
```

**Response**
```json
{
  "data": {
    "contacts": {
      "getById": {
        "address": "123 Main St",
        "created_at": "2024-04-16T12:00:00Z",
        "created_by": "Admin",
        "dob": "1990-01-01T00:00:00Z",
        "email": "john@example.com",
        "first_name": "John",
        "gender_id": 1,
        "last_name": "Doe",
        "name": "John Doe",
        "phone": "123-456-7890",
        "photo_path": "path/to/photo.jpg"
      },
      "getPagination": null,
      "gets": [
        {
          "address": "123 Main St",
          "created_at": "2024-04-16T12:00:00Z",
          "created_by": "Admin",
          "dob": "1990-01-01T00:00:00Z",
          "email": "john@example.com",
          "first_name": "John",
          "gender_id": 1,
          "last_name": "Doe",
          "name": "John Doe",
          "phone": "123-456-7890",
          "photo_path": "path/to/photo.jpg"
        },
        {
          "address": "456 Oak St",
          "created_at": "2024-04-16T12:30:00Z",
          "created_by": "Admin",
          "dob": "1992-05-15T00:00:00Z",
          "email": "jane@example.com",
          "first_name": "Jane",
          "gender_id": 2,
          "last_name": "Smith",
          "name": "Jane Smith",
          "phone": "987-654-3210",
          "photo_path": "path/to/photo2.jpg"
        },
        {
          "address": "123 Main St",
          "created_at": "2024-04-28T17:20:11.438835+07:00",
          "created_by": "test-api",
          "dob": "1979-11-13T00:00:00Z",
          "email": "john@example.com",
          "first_name": "PUP",
          "gender_id": 1,
          "last_name": "Apaichon",
          "name": "Dr.PUP",
          "phone": "123-456-7890",
          "photo_path": "path/to/photo.jpg"
        },
        {
          "address": "123 Main St",
          "created_at": "2024-04-28T17:27:49.559621+07:00",
          "created_by": "test-api",
          "dob": "1979-11-13T00:00:00Z",
          "email": "john@example.com",
          "first_name": "PUP",
          "gender_id": 1,
          "last_name": "Apaichon",
          "name": "Dr.PUP",
          "phone": "123-456-7890",
          "photo_path": "path/to/photo.jpg"
        },
        {
          "address": "123 Main St",
          "created_at": "2024-05-01T23:16:56.437155+07:00",
          "created_by": "test-api",
          "dob": "1979-11-13T00:00:00Z",
          "email": "john@example.com",
          "first_name": "PUP",
          "gender_id": 1,
          "last_name": "Apaichon",
          "name": "Dr.PUP",
          "phone": "123-456-7890",
          "photo_path": "path/to/photo.jpg"
        },
        {
          "address": "123 Main St",
          "created_at": "2024-05-01T23:16:56.437156+07:00",
          "created_by": "test-api",
          "dob": "1979-11-13T00:00:00Z",
          "email": "john@example.com",
          "first_name": "PUP2",
          "gender_id": 1,
          "last_name": "Apaichon2",
          "name": "Dr.PUP2",
          "phone": "123-456-7890",
          "photo_path": "path/to/photo.jpg"
        },
        {
          "address": "123 Main St",
          "created_at": "2024-05-01T23:17:43.885492+07:00",
          "created_by": "test-api",
          "dob": "1979-11-13T00:00:00Z",
          "email": "john@example.com",
          "first_name": "PUP",
          "gender_id": 1,
          "last_name": "Apaichon",
          "name": "Dr.PUP",
          "phone": "123-456-7890",
          "photo_path": "path/to/photo.jpg"
        },
        {
          "address": "123 Main St",
          "created_at": "2024-05-01T23:17:43.885494+07:00",
          "created_by": "test-api",
          "dob": "1979-11-13T00:00:00Z",
          "email": "john@example.com",
          "first_name": "PUP2",
          "gender_id": 1,
          "last_name": "Apaichon2",
          "name": "Dr.PUP2",
          "phone": "123-456-7890",
          "photo_path": "path/to/photo.jpg"
        },
        {
          "address": "123 Main St",
          "created_at": "2024-05-01T23:18:55.852702+07:00",
          "created_by": "test-api",
          "dob": "1979-11-13T00:00:00Z",
          "email": "john@example.com",
          "first_name": "PUP13",
          "gender_id": 1,
          "last_name": "Apaichon",
          "name": "Dr.PUP13",
          "phone": "123-456-7890",
          "photo_path": "path/to/photo.jpg"
        },
        {
          "address": "123 Main St",
          "created_at": "2024-05-01T23:18:55.852705+07:00",
          "created_by": "test-api",
          "dob": "1979-11-13T00:00:00Z",
          "email": "john@example.com",
          "first_name": "PUP11",
          "gender_id": 1,
          "last_name": "Apaichon2",
          "name": "Dr.PUP11",
          "phone": "123-456-7890",
          "photo_path": "path/to/photo.jpg"
        }
      ]
    }
  },
  "errors": [
    {
      "message": "unauthorized: missing contacts.getPagination permission",
      "locations": [
        {
          "line": 29,
          "column": 5
        }
      ],
      "path": [
        "contacts",
        "getPagination"
      ]
    }
  ]
}
```





