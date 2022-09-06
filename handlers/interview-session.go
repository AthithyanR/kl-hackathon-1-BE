package handlers

import (
	"encoding/json"

	"github.com/AthithyanR/kl-hackathon-1-BE/db"
	"github.com/AthithyanR/kl-hackathon-1-BE/models"
	"github.com/AthithyanR/kl-hackathon-1-BE/utils"
	"github.com/valyala/fasthttp"
)

func GetInterviewSession(ctx *fasthttp.RequestCtx) {
	var interviewSession models.InterviewSession
	whereClause := &models.InterviewSession{
		SessionKey: ctx.UserValue("sessionKey").(string),
	}
	db.DB.Where(whereClause).Find(&interviewSession)
	sendSuccessResponse(ctx, interviewSession)
}

func GetInterviewSessions(ctx *fasthttp.RequestCtx) {
	var interviewSessions []models.InterviewSession
	db.DB.Find(&interviewSessions)
	sendSuccessResponse(ctx, interviewSessions)
}

func AddInterviewSession(ctx *fasthttp.RequestCtx) {
	var interviewSession models.InterviewSession
	err := json.Unmarshal(ctx.PostBody(), &interviewSession)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}
	interviewSession.Id = utils.CanonicId()
	interviewSession.SessionKey = utils.CanonicId()
	result := db.DB.Create(&interviewSession)
	if result.Error != nil {
		sendFailureResponse(ctx, result.Error.Error())
		return
	}
	sendSuccessResponse(ctx, interviewSession.Id)
}

func UpdateInterviewSession(ctx *fasthttp.RequestCtx) {
	var interviewSession models.InterviewSession
	err := json.Unmarshal(ctx.PostBody(), &interviewSession)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}
	result := db.DB.Updates(&interviewSession)
	if result.Error != nil {
		sendFailureResponse(ctx, result.Error.Error())
		return
	}
	sendSuccessResponse(ctx, nil)
}
