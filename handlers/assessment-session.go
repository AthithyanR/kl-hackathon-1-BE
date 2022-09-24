package handlers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/AthithyanR/kl-hackathon-1-BE/db"
	"github.com/AthithyanR/kl-hackathon-1-BE/models"
	"github.com/AthithyanR/kl-hackathon-1-BE/utils"
	"github.com/valyala/fasthttp"
)

func GetAssessmentSessionMeta(ctx *fasthttp.RequestCtx) {
	var assessmentSession models.AssessmentSession
	queryParams := ctx.QueryArgs()
	sessionKey := string(queryParams.Peek("sessionKey"))
	if sessionKey == "" {
		sendSuccessResponse(ctx, nil)
		return
	}
	whereClause := &models.AssessmentSession{
		Id: sessionKey,
	}
	db.DB.Where(whereClause).Find(&assessmentSession)
	// to handle can start now?? logic
	if assessmentSession.Id == "" || assessmentSession.EndTime != nil {
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

func GetAssessmentSessionQuestion(ctx *fasthttp.RequestCtx) {
	var assessmentSession models.AssessmentSession
	queryParams := ctx.QueryArgs()
	assessmentSessionWhereClause := &models.AssessmentSession{
		Id: string(queryParams.Peek("sessionKey")),
	}
	techTypeIdQs := string(queryParams.Peek("techTypeId"))
	questionTypeQs := string(queryParams.Peek("questionType"))
	questionNumberQs, _ := strconv.Atoi(string(queryParams.Peek("questionNumber")))
	db.DB.Where(assessmentSessionWhereClause).Find(&assessmentSession)
	if assessmentSession.EndTime != nil {
		sendSuccessResponse(ctx, nil)
		return
	}
	questionData, err := getParsedQuestionData(assessmentSession.QuestionData)
	if err != nil {
		sendFailureResponse(ctx, nil)
		return
	}

	var foundQuestionId string

out:
	for techTypeId, questionTypeData := range questionData {
		if techTypeId != techTypeIdQs {
			continue
		}
		for questionType, questionIdentifier := range questionTypeData {
			if questionType != questionTypeQs {
				continue
			}
			if questionsSlices, ok := questionIdentifier.([]interface{}); ok {
				if questionSlice, ok := questionsSlices[questionNumberQs-1].([]interface{}); ok {
					if questionId, ok := questionSlice[0].(string); ok {
						foundQuestionId = questionId
						break out
					}
				}
			}
		}
	}

	var question models.QuestionLimited
	questionWhereClause := &models.Question{
		Id: foundQuestionId,
	}
	db.DB.Model(&models.Question{}).Where(questionWhereClause).Find(&question)

	sendSuccessResponse(ctx, question)
}

func EvaluateAnswer(ctx *fasthttp.RequestCtx) {
	var assessmentAnswer models.AssessmentAnswer
	err := json.Unmarshal(ctx.PostBody(), &assessmentAnswer)
	if err != nil {
		fmt.Println(err)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}
	assessmentSessionWhereClause := &models.AssessmentSession{
		Id: assessmentAnswer.SessionKey,
	}

	var assessmentSession models.AssessmentSession
	db.DB.Where(assessmentSessionWhereClause).Find(&assessmentSession)
	if assessmentSession.EndTime != nil {
		sendSuccessResponse(ctx, nil)
		return
	}

	questionData, err := getParsedQuestionData(assessmentSession.QuestionData)
	if err != nil {
		sendFailureResponse(ctx, nil)
		return
	}

out:
	for techTypeId, questionTypeData := range questionData {
		if techTypeId != assessmentAnswer.TechTypeId {
			continue
		}
		for questionType, questionIdentifier := range questionTypeData {
			if questionType != assessmentAnswer.QuestionType {
				continue
			}
			if questionsSlices, ok := questionIdentifier.([]interface{}); ok {
				if questionSlice, ok := questionsSlices[assessmentAnswer.QuestionNumber-1].([]interface{}); ok {
					if questionId, ok := questionSlice[0].(string); ok {
						var count int64
						questionWhereClause := &models.Question{
							Id:            questionId,
							CorrectOption: assessmentAnswer.ChosenOption,
						}
						db.DB.Model(&models.Question{}).Where(questionWhereClause).Count(&count)

						var isCorrectAnswer string

						if count == 0 {
							isCorrectAnswer = "false"
						} else {
							isCorrectAnswer = "true"
						}

						if questionSlice[1] != isCorrectAnswer {
							questionSlice[1] = isCorrectAnswer
							questionsSlices[assessmentAnswer.QuestionNumber-1] = questionSlice
							questionData[techTypeId][questionType] = questionsSlices

							res, err := json.Marshal(questionData)
							if err != nil {
								sendFailureResponse(ctx, nil)
								return
							}
							assessmentSession.QuestionData = string(res)
							if err := db.DB.Updates(&assessmentSession).Error; err != nil {
								sendFailureResponse(ctx, err)
								return
							}
						}

						break out
					}
				}
			}
		}
	}
	sendSuccessResponse(ctx, nil)
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
		isEmailSent := utils.SendMail([]string{candidateEmail}, fmt.Sprintf("Please use this link- %s", sessionId))
		assessmentSession := &models.AssessmentSession{
			Id:                sessionId,
			CandidateEmail:    candidateEmail,
			QuestionData:      assessmentSessionCreate.QuestionData,
			TimeAllowedInMins: assessmentSessionCreate.TimeAllowedInMins,
			IsEmailSent:       isEmailSent,
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
	isEmailSent := utils.SendMail([]string{assessmentSession.CandidateEmail}, fmt.Sprintf("Please use this link- %s", assessmentSession.Id))
	assessmentSession.IsEmailSent = isEmailSent
	result := db.DB.Updates(&assessmentSession)
	if result.Error != nil {
		sendFailureResponse(ctx, result.Error.Error())
		return
	}
	sendSuccessResponse(ctx, nil)
}

func CompleteAssessmentSession(ctx *fasthttp.RequestCtx) {
	var assessmentSession models.AssessmentSession
	queryParams := ctx.QueryArgs()
	whereClause := &models.AssessmentSession{
		Id: string(queryParams.Peek("sessionKey")),
	}
	db.DB.Where(whereClause).Find(&assessmentSession)
	if assessmentSession.EndTime != nil {
		sendSuccessResponse(ctx, nil)
		return
	}

	questionData, err := getParsedQuestionData(assessmentSession.QuestionData)
	if err != nil {
		sendFailureResponse(ctx, nil)
		return
	}

	var questionIds []string
	for _, questionTypeData := range questionData {
		for _, questionIdentifier := range questionTypeData {
			if questionsSlices, ok := questionIdentifier.([]interface{}); ok {
				for _, questionSlice := range questionsSlices {
					if questionSlice, ok := questionSlice.([]interface{}); ok {
						if questionId, ok := questionSlice[0].(string); ok {
							if questionSlice[1] == "true" {
								questionIds = append(questionIds, questionId)
							}
						}
					}
				}
			}
		}
	}

	fmt.Println(getScoreByQuestionIds(questionIds))
	assessmentSession.Score = getScoreByQuestionIds(questionIds)
	fmt.Println(assessmentSession.Score)
	assessmentSession.ScoreOutOf100Percent = float32(assessmentSession.Score) / float32(assessmentSession.PossibleScore) * 100
	currentTime := time.Now()
	assessmentSession.EndTime = &currentTime
	if err := db.DB.Updates(&assessmentSession).Error; err != nil {
		sendFailureResponse(ctx, err)
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
