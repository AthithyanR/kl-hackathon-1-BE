package router

import (
	"github.com/AthithyanR/kl-hackathon-1-BE/handlers"
	"github.com/AthithyanR/kl-hackathon-1-BE/utils"
	"github.com/fasthttp/router"
)

func InitRouter() *router.Router {
	r := router.New()

	r.GET("/api/healthCheck", handlers.HealthCheck)

	//TechTypes
	r.GET("/api/techTypes", utils.Middleware(handlers.GetAllTechTypes))

	//Questions
	r.GET("/api/questions", utils.Middleware(handlers.GetQuestions))

	return r
}
