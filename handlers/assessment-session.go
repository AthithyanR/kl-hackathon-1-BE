package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/AthithyanR/kl-hackathon-1-BE/db"
	"github.com/AthithyanR/kl-hackathon-1-BE/models"
	"github.com/AthithyanR/kl-hackathon-1-BE/utils"
	"github.com/valyala/fasthttp"
)

func GetAssessmentSession(ctx *fasthttp.RequestCtx) {
	var assessmentSession models.AssessmentSession
	whereClause := &models.AssessmentSession{
		Id: ctx.UserValue("sessionKey").(string),
	}
	db.DB.Where(whereClause).Find(&assessmentSession)
	sendSuccessResponse(ctx, assessmentSession)
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

	var questionData models.QuestionData
	err = json.Unmarshal([]byte(assessmentSessionCreate.QuestionData), &questionData)
	if err != nil {
		fmt.Println(err)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	for k, v := range questionData {
		fmt.Println(k)
		for k2, v2 := range v {
			fmt.Println(k2, v2)
			if _, ok := v2.(int); ok {
				fmt.Println("v2 is an integer")
			} else if vtwo, ok := v2.(map[string]interface{}); ok {
				fmt.Println("v2 is a map")
				for k3, v3 := range vtwo {
					fmt.Println(k3, v3)
				}
			}

		}
	}

	candidateEmails := assessmentSessionCreate.CandidateEmails

	for _, candidateEmail := range candidateEmails {
		sessionId := utils.CanonicId()
		isEmailSent := utils.SendMail([]string{candidateEmail}, fmt.Sprintf("Please use this link- %s", sessionId))
		assessmentSession := &models.AssessmentSession{
			Id:                sessionId,
			CandidateEmail:    candidateEmail,
			QuestionData:      assessmentSessionCreate.QuestionData,
			TimeAllowedInMins: assessmentSessionCreate.TimeAllowedInMins,
			IsEmailSent:       isEmailSent,
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
	result := db.DB.Updates(&assessmentSession)
	if result.Error != nil {
		sendFailureResponse(ctx, result.Error.Error())
		return
	}
	sendSuccessResponse(ctx, nil)
}
