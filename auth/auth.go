package auth

import (
	"errors"
	"os"
	"time"

	"github.com/AthithyanR/kl-hackathon-1-BE/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/valyala/fasthttp"
)

var (
	JwtAuthPrefix = []byte("Bearer ")
)

func GenerateToken(creds *models.ClaimValues) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &models.Claims{
		Id:    creds.Id,
		Email: creds.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	return tokenString, err
}

func ValidateToken(tknStr string) error {
	claims := &models.Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return err
	}
	if !tkn.Valid {
		return errors.New("token expired")
	}
	return nil
}

func GetClaimsFromCtx(ctx *fasthttp.RequestCtx) *models.Claims {
	tknStr := string(ctx.Request.Header.Peek("Authorization")[len(JwtAuthPrefix):])
	claims := &models.Claims{}
	jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	return claims
}
