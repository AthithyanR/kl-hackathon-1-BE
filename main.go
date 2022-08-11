package main

import (
	"log"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

func main() {
	r := router.New()
	r.GET("/api/healthCheck", healthCheck)

	log.Fatal(fasthttp.ListenAndServe(":8080", r.Handler))
}
