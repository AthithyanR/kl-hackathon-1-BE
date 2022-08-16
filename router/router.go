package router

import (
	"github.com/AthithyanR/kl-hackathon-1-BE/handlers"
	"github.com/AthithyanR/kl-hackathon-1-BE/utils"
	"github.com/fasthttp/router"
)

func InitRouter() *router.Router {
	r := router.New()

	r.GET("/api/healthCheck", handlers.HealthCheck)

	//Auth
	r.POST("/api/authenticate", utils.Middleware(handlers.Authenticate))

	//TechTypes
	r.GET("/api/techTypes", utils.MiddlewareWithAuth(handlers.GetAllTechTypes))

	//Questions
	r.GET("/api/questions", utils.MiddlewareWithAuth(handlers.GetQuestions))

	return r
}
