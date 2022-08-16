package utils

import (
	"log"
	"os"

	"github.com/AthithyanR/kl-hackathon-1-BE/auth"
	"github.com/valyala/fasthttp"
)

func Getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

var (
	strContentType     = []byte("Content-Type")
	strApplicationJSON = []byte("application/json")
	JwtAuthPrefix      = []byte("Bearer ")
)

func Middleware(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.SetCanonical(strContentType, strApplicationJSON)
		log.Println(ctx)
		h(ctx)
	}
}

func MiddlewareWithAuth(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		authHeader := ctx.Request.Header.Peek("Authorization")
		if len(authHeader) == 0 {
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			return
		}
		tokenString := string(authHeader[len(JwtAuthPrefix):])
		err := auth.ValidateToken(tokenString)
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			return
		}
		ctx.Response.Header.SetCanonical(strContentType, strApplicationJSON)
		log.Println(ctx)
		h(ctx)
	}
}
