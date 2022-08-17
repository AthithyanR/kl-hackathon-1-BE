package handlers

import (
	"encoding/json"

	"github.com/AthithyanR/kl-hackathon-1-BE/db"
	"github.com/AthithyanR/kl-hackathon-1-BE/models"
	"github.com/AthithyanR/kl-hackathon-1-BE/utils"
	"github.com/valyala/fasthttp"
	"gorm.io/gorm/clause"
)

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
