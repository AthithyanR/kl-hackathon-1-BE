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
	r.GET("/api/questions/{id}", utils.MiddlewareWithAuth(handlers.GetQuestion))
	r.POST("/api/questions", utils.MiddlewareWithAuth(handlers.AddQuestions))
	r.PUT("/api/questions", utils.MiddlewareWithAuth(handlers.UpdateQuestions))
	r.DELETE("/api/questions", utils.MiddlewareWithAuth(handlers.DeleteQuestions))

	//Interview Session
	r.GET("/api/interviewSession/{sessionKey}", utils.MiddlewareWithAuth(handlers.GetInterviewSession))
	r.GET("/api/interviewSession", utils.MiddlewareWithAuth(handlers.GetInterviewSessions))
	r.POST("/api/interviewSession", utils.MiddlewareWithAuth(handlers.AddInterviewSession))
	r.PUT("/api/interviewSession", utils.MiddlewareWithAuth(handlers.UpdateInterviewSession))

	return r
}
