package models

import (
	"github.com/dgrijalva/jwt-go"
)

type ClaimValues struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

type Claims struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	jwt.StandardClaims
}
