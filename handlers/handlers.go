package handlers

import (
	"encoding/json"

	"github.com/AthithyanR/kl-hackathon-1-BE/auth"
	"github.com/AthithyanR/kl-hackathon-1-BE/db"
	"github.com/AthithyanR/kl-hackathon-1-BE/models"
	"github.com/AthithyanR/kl-hackathon-1-BE/utils"
	"github.com/valyala/fasthttp"
	"gorm.io/gorm/clause"
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
		sendFailureResponse(ctx, err.Error())
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

func GetQuestion(ctx *fasthttp.RequestCtx) {
	var question models.Question
	whereClause := &models.Question{
		Id: ctx.UserValue("id").(string),
	}
	db.DB.Where(whereClause).Find(&question)
	sendSuccessResponse(ctx, question)
}

func AddQuestions(ctx *fasthttp.RequestCtx) {
	var questions []models.Question
	var response []string
	err := json.Unmarshal(ctx.PostBody(), &questions)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}
	for i := 0; i < len(questions); i++ {
		id := utils.CanonicId()
		questions[i].Id = id
		response = append(response, id)
	}
	result := db.DB.Create(&questions)
	if result.Error != nil {
		sendFailureResponse(ctx, result.Error.Error())
		return
	}
	sendSuccessResponse(ctx, response)
}

func UpdateQuestions(ctx *fasthttp.RequestCtx) {
	var questions []models.Question
	err := json.Unmarshal(ctx.PostBody(), &questions)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}
	result := db.DB.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&questions)
	if result.Error != nil {
		sendFailureResponse(ctx, result.Error.Error())
		return
	}
	sendSuccessResponse(ctx, nil)
}

func DeleteQuestions(ctx *fasthttp.RequestCtx) {
	var questionIds []string
	err := json.Unmarshal(ctx.PostBody(), &questionIds)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}
	result := db.DB.Delete(&models.Question{}, questionIds)
	if result.Error != nil {
		sendFailureResponse(ctx, result.Error.Error())
		return
	}
	sendSuccessResponse(ctx, nil)
}
