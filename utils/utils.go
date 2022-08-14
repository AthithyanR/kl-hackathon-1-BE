package utils

import (
	"log"
	"os"

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
)

func Middleware(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.SetCanonical(strContentType, strApplicationJSON)
		log.Println(ctx)
		h(ctx)
	}
}
