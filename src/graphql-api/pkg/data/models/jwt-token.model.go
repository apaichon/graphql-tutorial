package models

type JwtToken struct {
	Token string `json:"token"`
	ExpiredAt int64 `json:"expiredAt"`
	
}
