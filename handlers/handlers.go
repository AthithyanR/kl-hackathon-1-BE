package handlers

import (
	"encoding/json"

	"github.com/AthithyanR/kl-hackathon-1-BE/db"
	"github.com/AthithyanR/kl-hackathon-1-BE/models"
	"github.com/valyala/fasthttp"
)

type BaseResponse struct {
	Success bool `json:"success"`
	Data    any  `json:"data"`
}

func sendSuccessResponse(ctx *fasthttp.RequestCtx, data any) {
	json.NewEncoder(ctx).Encode(&BaseResponse{Success: true, Data: data})
}

func HealthCheck(ctx *fasthttp.RequestCtx) {
	sendSuccessResponse(ctx, nil)
}

// techType
func GetAllTechTypes(ctx *fasthttp.RequestCtx) {
	var techTypes []models.TechType
	db.DB.Find(&techTypes)
	sendSuccessResponse(ctx, techTypes)
}

// questions
func GetQuestions(ctx *fasthttp.RequestCtx) {
	var questions []models.Question
	whereClause := &models.Question{
		TechType:     string(ctx.QueryArgs().Peek("techType")),
		QuestionType: string(ctx.QueryArgs().Peek("questionType")),
	}
	db.DB.Where(whereClause).Find(&questions)
	sendSuccessResponse(ctx, questions)
}
