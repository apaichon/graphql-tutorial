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

func GetUserName(p graphql.ResolveParams)(models.JwtClaims, error) {
	tokenString, _ := p.Context.Value(userKey).(string)
		config := config.NewConfig()
		claim, err := DecodeJWTToken(tokenString, config.SecretKey)
		return *claim, err
}


