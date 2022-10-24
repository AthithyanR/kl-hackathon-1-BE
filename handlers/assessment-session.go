package handlers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/AthithyanR/kl-hackathon-1-BE/db"
	"github.com/AthithyanR/kl-hackathon-1-BE/models"
	"github.com/AthithyanR/kl-hackathon-1-BE/queries"
	"github.com/AthithyanR/kl-hackathon-1-BE/utils"
	"github.com/valyala/fasthttp"
)

// yet to do
func GetAssessmentSessionMeta(ctx *fasthttp.RequestCtx) {
	var assessmentSession models.AssessmentSession
	queryParams := ctx.QueryArgs()
	assessmentSessionId := string(queryParams.Peek("assessmentSessionId"))
	if assessmentSessionId == "" {
		sendSuccessResponse(ctx, nil)
		return
	}
	whereClause := &models.AssessmentSession{
		Id: assessmentSessionId,
	}
	db.DB.Where(whereClause).Find(&assessmentSession)
	// to handle can start now logic ??
	if assessmentSession.Id == "" || assessmentSession.EndTime != nil {
		sendSuccessResponse(ctx, nil)
		return
	}
	var assessmentSessionMeta models.AssessmentSessionMeta
	// addQuestionMeta(&assessmentSessionMeta, &assessmentSession)
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
	queryParams := ctx.QueryArgs()
	assessmentSessionId := string(queryParams.Peek("assessmentSessionId"))
	questionId := string(queryParams.Peek("questionId"))

	var assessmentSession models.AssessmentSession
	whereClause := models.AssessmentSession{
		Id: assessmentSessionId,
	}
	if err := db.DB.Model(&models.AssessmentSession{}).
		Where(whereClause).
		Find(&assessmentSession).
		Error; err != nil {
		sendFailureResponse(ctx, err)
		return
	}
	if assessmentSession.EndTime != nil {
		sendSuccessResponse(ctx, nil)
		return
	}
	var question models.QuestionLimited
	questionWhereClause := &models.Question{
		Id: questionId,
	}
	if err := db.DB.Model(&models.Question{}).
		Where(questionWhereClause).
		Find(&question).Error; err != nil {
		sendFailureResponse(ctx, err)
		return
	}

	sendSuccessResponse(ctx, question)
}

func EvaluateAnswer(ctx *fasthttp.RequestCtx) {
	var assessmentAnswer models.AssessmentSessionQuestion
	err := json.Unmarshal(ctx.PostBody(), &assessmentAnswer)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	if err := db.DB.Updates(&assessmentAnswer).Error; err != nil {
		sendFailureResponse(ctx, err)
		return
	}

	if err := db.DB.Exec(
		queries.UpdateScoreByAnswer,
		assessmentAnswer.AssessmentSessionId,
		assessmentAnswer.AssessmentSessionId,
	).Error; err != nil {
		sendFailureResponse(ctx, err)
		return
	}

	if err := db.DB.Exec(
		queries.UpdateScorePercentageById,
		assessmentAnswer.AssessmentSessionId,
	).Error; err != nil {
		sendFailureResponse(ctx, err)
		return
	}

	sendSuccessResponse(ctx, nil)
}

func GetAssessmentSessions(ctx *fasthttp.RequestCtx) {
	var assessmentSessions []models.AssessmentSession
	if err := db.DB.Find(&assessmentSessions).Error; err != nil {
		sendFailureResponse(ctx, err)
		return
	}
	sendSuccessResponse(ctx, assessmentSessions)
}

func GetAssessmentSessionDetailsById(ctx *fasthttp.RequestCtx) {
	var assessmentSessionDetails []models.AssessmentSessionDetails
	assessmentSessionId := ctx.UserValue("id").(string)
	if err := db.DB.Raw(
		queries.GetAssessmentSessionMetaById,
		assessmentSessionId).
		Scan(&assessmentSessionDetails).
		Error; err != nil {
		sendFailureResponse(ctx, err)
		return
	}
	sendSuccessResponse(ctx, assessmentSessionDetails)
}

func AddAssessmentSession(ctx *fasthttp.RequestCtx) {
	var assessmentSessionCreate models.AssessmentSessionUpsert
	err := json.Unmarshal(ctx.PostBody(), &assessmentSessionCreate)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	newQuestionIds, err := getQuestionIdByRandomCount(assessmentSessionCreate.RandomQuestions)
	if err != nil {
		sendFailureResponse(ctx, err)
		return
	}

	allQuestionIds := append(assessmentSessionCreate.QuestionIds, newQuestionIds...)
	possibleScore := getScoreByQuestionIds(allQuestionIds)

	for _, candidateEmail := range assessmentSessionCreate.CandidateEmails {
		sessionId := utils.CanonicId()
		subject := "Subject: Assessment Link!\n"
		mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
		bodyRaw := "<html><body>Please use this link - link - http://localhost:5173/assessment?sessionKey=%s</body></html>"
		bodyWithContent := fmt.Sprintf(bodyRaw, sessionId)

		isEmailSent := utils.SendMail(
			[]string{candidateEmail},
			subject+mime+bodyWithContent,
		)

		assessmentSession := &models.AssessmentSession{
			Id:                sessionId,
			CandidateEmail:    candidateEmail,
			TimeAllowedInMins: assessmentSessionCreate.TimeAllowedInMins,
			IsEmailSent:       isEmailSent,
			PossibleScore:     possibleScore,
			QuestionsCount:    len(allQuestionIds),
		}
		if err := db.DB.Create(&assessmentSession).Error; err != nil {
			sendFailureResponse(ctx, err)
			return
		}

		if err := upsertAssessmentSessionQuestions(sessionId, allQuestionIds); err != nil {
			sendFailureResponse(ctx, err)
			return
		}
	}
	sendSuccessResponse(ctx, nil)
}

func UpdateAssessmentSession(ctx *fasthttp.RequestCtx) {
	var assessmentSessionUpdate models.AssessmentSessionUpsert
	err := json.Unmarshal(ctx.PostBody(), &assessmentSessionUpdate)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	subject := "Subject: Assessment Link!\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	bodyRaw := "<html><body>Please use this link - link - http://localhost:5173/assessment?sessionKey=%s</body></html>"
	bodyWithContent := fmt.Sprintf(bodyRaw, assessmentSessionUpdate.Id)

	assessmentSessionUpdate.IsEmailSent = utils.SendMail(
		[]string{assessmentSessionUpdate.CandidateEmail},
		subject+mime+bodyWithContent,
	)

	newQuestionIds, err := getQuestionIdByRandomCount(assessmentSessionUpdate.RandomQuestions)
	if err != nil {
		sendFailureResponse(ctx, err)
		return
	}

	allQuestionIds := append(assessmentSessionUpdate.QuestionIds, newQuestionIds...)

	if err = upsertAssessmentSessionQuestions(assessmentSessionUpdate.Id, allQuestionIds); err != nil {
		sendFailureResponse(ctx, err)
		return
	}

	assessmentSessionUpdate.QuestionsCount = len(allQuestionIds)
	assessmentSessionUpdate.PossibleScore = getScoreByQuestionIds(allQuestionIds)
	if err := db.DB.Updates(&assessmentSessionUpdate.AssessmentSession).Error; err != nil {
		sendFailureResponse(ctx, err)
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
	if err := db.DB.Where(whereClause).Find(&assessmentSession).Error; err != nil {
		sendFailureResponse(ctx, err)
		return
	}
	if assessmentSession.EndTime != nil {
		sendSuccessResponse(ctx, nil)
		return
	}
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
		sendFailureResponse(ctx, "No ids are provided")
		return
	}

	if err := db.DB.Delete(&models.AssessmentSession{}, sessionIds).Error; err != nil {
		sendFailureResponse(ctx, err)
		return
	}
	sendSuccessResponse(ctx, nil)
}
