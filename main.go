package main

import (
	"log"

	"github.com/valyala/fasthttp"

	"github.com/AthithyanR/kl-hackathon-1-BE/db"
	"github.com/AthithyanR/kl-hackathon-1-BE/router"
	"github.com/AthithyanR/kl-hackathon-1-BE/utils"
)

func main() {

	SERVER_PORT := utils.Getenv("port", ":8080")

	db.InitDb()
	r := router.InitRouter()

	log.Println("Application started on PORT ", SERVER_PORT)
	log.Fatal(fasthttp.ListenAndServe(SERVER_PORT, r.Handler))
}
