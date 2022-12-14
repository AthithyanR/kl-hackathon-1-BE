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
	r.GET("/api/techTypes", utils.Middleware(handlers.GetAllTechTypes))
	r.POST("/api/techTypes", utils.MiddlewareWithAuth(handlers.AddTechTypes))
	r.PUT("/api/techTypes", utils.MiddlewareWithAuth(handlers.UpdateTechTypes))
	r.DELETE("/api/techTypes", utils.MiddlewareWithAuth(handlers.DeleteTechTypes))

	//Questions
	r.GET("/api/questions", utils.MiddlewareWithAuth(handlers.GetQuestions))
	r.GET("/api/questions/{id}", utils.MiddlewareWithAuth(handlers.GetQuestion))
	r.POST("/api/questions", utils.MiddlewareWithAuth(handlers.AddQuestions))
	r.PUT("/api/questions", utils.MiddlewareWithAuth(handlers.UpdateQuestions))
	r.DELETE("/api/questions", utils.MiddlewareWithAuth(handlers.DeleteQuestions))

	//Assessment Session
	r.GET("/api/assessmentSession/meta", utils.Middleware(handlers.GetAssessmentSessionMeta))
	r.GET("/api/assessmentSession/question", utils.Middleware(handlers.GetAssessmentSessionQuestion))
	r.POST("/api/assessmentSession/evaluateAnswer", utils.Middleware(handlers.EvaluateAnswer))
	r.GET("/api/assessmentSession", utils.MiddlewareWithAuth(handlers.GetAssessmentSessions))
	r.GET("/api/assessmentSession/{id}", utils.MiddlewareWithAuth(handlers.GetAssessmentSessionDetailsById))
	r.POST("/api/assessmentSession", utils.MiddlewareWithAuth(handlers.AddAssessmentSession))
	r.PUT("/api/assessmentSession", utils.MiddlewareWithAuth(handlers.UpdateAssessmentSession))
	r.PUT("/api/assessmentSession/complete", utils.Middleware(handlers.CompleteAssessmentSession))
	r.DELETE("/api/assessmentSession", utils.MiddlewareWithAuth(handlers.DeleteAssessmentSession))

	return r
}
