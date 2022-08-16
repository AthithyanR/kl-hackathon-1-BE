package handlers

import (
	"encoding/json"

	"github.com/AthithyanR/kl-hackathon-1-BE/auth"
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

func sendFailureResponse(ctx *fasthttp.RequestCtx, data any) {
	json.NewEncoder(ctx).Encode(&BaseResponse{Success: false, Data: data})
}

func HealthCheck(ctx *fasthttp.RequestCtx) {
	sendSuccessResponse(ctx, nil)
}

//Auth

func Authenticate(ctx *fasthttp.RequestCtx) {
	var requestBody models.User
	err := json.Unmarshal(ctx.PostBody(), &requestBody)
	if err != nil || requestBody.Email == "" || requestBody.Password == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	var existingUser models.User
	whereClause := &models.User{Email: requestBody.Email}
	db.DB.Where(whereClause).Find(&existingUser)

	if existingUser.Id == "" || existingUser.Password != requestBody.Password {
		sendFailureResponse(ctx, "Invalid username or password")
		return
	}

	token, err := auth.GenerateToken(&models.ClaimValues{Id: existingUser.Id, Email: existingUser.Email})

	if err != nil {
		sendFailureResponse(ctx, "Unable to generate token")
		return
	}

	sendSuccessResponse(ctx, token)
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
