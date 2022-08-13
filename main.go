package main

import (
	"log"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"gorm.io/gorm"
)

var DB *gorm.DB

func initRoutes(r *router.Router) {
	r.GET("/api/healthCheck", healthCheck)

	//Questions
	r.GET("/api/questions/byTech/{techType}", middleware(getQuestionsByTechType))
}

func main() {

	DB = getDb()
	r := router.New()
	initRoutes(r)

	log.Fatal(fasthttp.ListenAndServe(":8080", r.Handler))
}
