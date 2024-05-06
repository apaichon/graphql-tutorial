package main

import (
	"fmt"
	"graphql-api/config"
	"graphql-api/internal/auth"

)

// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjowLCJ1c2VyX25hbWUiOiJhZG1pbiIsImV4cCI6MTcxNDkxMjUyMH0.WEfSq4P0AKfINYSAdb9lEt1Hx0TkTXN67p1G4Adl1Ss

func main() {
	config := config.NewConfig()
	var token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyLCJ1c2VyX25hbWUiOiJwdXBweSIsImV4cCI6MTcxNDkyMTA5OH0.mmLRoAM4aIKVAPZgrEeI6FGBL7X8-w5A45p2V0U5Q44"
	fmt.Println("Secrete:", config.SecretKey )
	jwt, err :=  auth.DecodeJWTToken(token,config.SecretKey )
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("UserId:", jwt.UserId )
	fmt.Println("UserName:", jwt.Username )

}