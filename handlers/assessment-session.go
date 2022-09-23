package handlers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/AthithyanR/kl-hackathon-1-BE/db"
	"github.com/AthithyanR/kl-hackathon-1-BE/models"
	"github.com/AthithyanR/kl-hackathon-1-BE/utils"
	"github.com/valyala/fasthttp"
)

func GetAssessmentSessionMeta(ctx *fasthttp.RequestCtx) {
	var assessmentSession models.AssessmentSession
	queryParams := ctx.QueryArgs()
	whereClause := &models.AssessmentSession{
		Id: string(queryParams.Peek("sessionKey")),
	}
	db.DB.Where(whereClause).Find(&assessmentSession)
	// to handle can start now?? logic
	if assessmentSession.Id == "" {
		sendSuccessResponse(ctx, nil)
		return
	}
	var assessmentSessionMeta models.AssessmentSessionMeta
	addQuestionMeta(&assessmentSessionMeta, &assessmentSession)
	assessmentSessionMeta.CandidateEmail = assessmentSession.CandidateEmail
	assessmentSessionMeta.QuestionsCount = assessmentSession.QuestionsCount
	assessmentSessionMeta.TimeAllowedInMins = assessmentSession.TimeAllowedInMins
	assessmentSessionMeta.StartTime = assessmentSession.StartTime

	if assessmentSession.StartTime == nil {
		currentTime := time.Now()
		assessmentSession.StartTime = &currentTime
		db.DB.Updates(&assessmentSession)
	}
	sendSuccessResponse(ctx, assessmentSessionMeta)
}

func GetAssessmentSessions(ctx *fasthttp.RequestCtx) {
	var assessmentSessions []models.AssessmentSession
	db.DB.Find(&assessmentSessions)
	sendSuccessResponse(ctx, assessmentSessions)
}

func AddAssessmentSession(ctx *fasthttp.RequestCtx) {
	var assessmentSessionCreate models.AssessmentSessionCreate
	err := json.Unmarshal(ctx.PostBody(), &assessmentSessionCreate)
	if err != nil {
		fmt.Println(err)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	if err := processAssessmentSessionCreateData(&assessmentSessionCreate); err != nil {
		fmt.Println(err)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	for _, candidateEmail := range assessmentSessionCreate.CandidateEmails {
		sessionId := utils.CanonicId()
		// isEmailSent := utils.SendMail([]string{candidateEmail}, fmt.Sprintf("Please use this link- %s", sessionId))
		assessmentSession := &models.AssessmentSession{
			Id:                sessionId,
			CandidateEmail:    candidateEmail,
			QuestionData:      assessmentSessionCreate.QuestionData,
			TimeAllowedInMins: assessmentSessionCreate.TimeAllowedInMins,
			IsEmailSent:       true,
			PossibleScore:     assessmentSessionCreate.PossibleScore,
			QuestionsCount:    assessmentSessionCreate.QuestionsCount,
		}
		result := db.DB.Create(&assessmentSession)
		if result.Error != nil {
			sendFailureResponse(ctx, result.Error.Error())
			return
		}
	}
	sendSuccessResponse(ctx, nil)
}

func UpdateAssessmentSession(ctx *fasthttp.RequestCtx) {
	var assessmentSession models.AssessmentSession
	err := json.Unmarshal(ctx.PostBody(), &assessmentSession)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}
	if err := processAssessmentSessionData(&assessmentSession); err != nil {
		fmt.Println(err)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}
	// isEmailSent := utils.SendMail([]string{assessmentSession.CandidateEmail}, fmt.Sprintf("Please use this link- %s", assessmentSession.Id))
	assessmentSession.IsEmailSent = true
	result := db.DB.Updates(&assessmentSession)
	if result.Error != nil {
		sendFailureResponse(ctx, result.Error.Error())
		return
	}
	sendSuccessResponse(ctx, nil)
}

func DeleteAssessmentSession(ctx *fasthttp.RequestCtx) {
	var sessionIds []string
	err := json.Unmarshal(ctx.PostBody(), &sessionIds)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}
	if len(sessionIds) == 0 {
		sendFailureResponse(ctx, "No ids provided")
		return
	}
	result := db.DB.Delete(&models.AssessmentSession{}, sessionIds)
	if result.Error != nil {
		sendFailureResponse(ctx, result.Error.Error())
		return
	}
	sendSuccessResponse(ctx, nil)
}
