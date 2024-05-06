package test

import (
	"fmt"
	"graphql-api/config"
	"graphql-api/internal/auth"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJWTTokenParse(t *testing.T) {

	config := config.NewConfig()
	var token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyLCJ1c2VyX25hbWUiOiJwdXBweSIsImV4cCI6MTcxNDk3OTc2N30.z8HeCicn3XFUYLqu3POx2gDNcEy9EExSHvftq37aPJE"
	fmt.Println("Secrete:", config.SecretKey)
	jwt, err := auth.DecodeJWTToken(token, config.SecretKey)
	if err != nil {
		fmt.Println(err)
	}
	assert.Greater(t, jwt.UserId, 0)
	assert.Greater(t, len(jwt.Username), 0)
}
