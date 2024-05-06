package models
import("github.com/dgrijalva/jwt-go")

type JwtClaims struct {
	UserId int `json:"user_id"`
	Username string `json:"user_name"`
	jwt.StandardClaims
}
