package handlers

import (
	"math/rand"
	"time"

	"github.com/AthithyanR/kl-hackathon-1-BE/db"
	"github.com/AthithyanR/kl-hackathon-1-BE/models"
	"github.com/thoas/go-funk"
)

func getQuestionIds(techTypeId string, questionType string) ([]string, error) {
	var questionIds []string
	whereClause := &models.Question{
		TechTypeId:   techTypeId,
		QuestionType: questionType,
	}
	if err := db.DB.Model(&models.Question{}).Where(whereClause).Pluck("id", &questionIds).Error; err != nil {
		return questionIds, err
	}
	return questionIds, nil
}

func getQuestionIdByRandomCount(randomSlice []models.RandomQuestions) ([]string, error) {
	var allQuestionIds []string

	for _, randomStruct := range randomSlice {
		questionIds, err := getQuestionIds(randomStruct.TechTypeId, randomStruct.QuestionType)

		if err != nil {
			return nil, err
		}

		questionIdMap := make(map[string]bool)
		seed := rand.NewSource(time.Now().Unix())
		r := rand.New(seed)
		questionCountCopy := randomStruct.Count

		for questionCountCopy > 0 {
			pickedQuestionId := questionIds[r.Intn(randomStruct.Count)]
			if _, ok := questionIdMap[pickedQuestionId]; !ok {
				questionIdMap[pickedQuestionId] = true
				allQuestionIds = append(allQuestionIds, pickedQuestionId)
				questionCountCopy--
			}
		}
	}

	return allQuestionIds, nil
}

func upsertAssessmentSessionQuestions(assessmentSessionId string, allQuestionIds []string) error {
	// batch delete assessment session questions if not in allQuestionIds list
	if err := db.DB.Delete(
		&models.AssessmentSessionQuestion{},
		"assessment_session_id = ? AND question_id NOT IN ?",
		assessmentSessionId, allQuestionIds,
	).Error; err != nil {
		return err
	}

	// batch insert assessment session questions
	var existingQuestionIds []string
	if err := db.DB.Model(&models.AssessmentSessionQuestion{}).
		Where("assessment_session_id = ? AND question_id IN ?", assessmentSessionId, allQuestionIds).
		Pluck("question_id", &existingQuestionIds).Error; err != nil {
		return err
	}

	questionsToCreate, _ := funk.DifferenceString(allQuestionIds, existingQuestionIds)

	if len(questionsToCreate) != 0 {
		var assessmentSessionQuestions []models.AssessmentSessionQuestion
		for _, questionId := range questionsToCreate {
			assessmentSessionQuestion := models.AssessmentSessionQuestion{
				AssessmentSessionId: assessmentSessionId,
				QuestionId:          questionId,
			}
			assessmentSessionQuestions = append(assessmentSessionQuestions, assessmentSessionQuestion)
		}

		if err := db.DB.Create(&assessmentSessionQuestions).Error; err != nil {
			return err
		}
	}

	return nil
}

// func getParsedQuestionData(questionDataString string) (models.QuestionData, error) {
// 	var questionData models.QuestionData
// 	err := json.Unmarshal([]byte(questionDataString), &questionData)
// 	if err != nil {
// 		return questionData, err
// 	}
// 	return questionData, nil
// }

// func getProcessedData(questionData models.QuestionData) (models.QuestionData, []string, error) {
// 	allQuestionsId := make([]string, 0)

// 	for techTypeId, questionTypeData := range questionData {
// 		for questionType, questionIdentifier := range questionTypeData {
// 			if questionCountFloat, ok := questionIdentifier.(float64); ok {
// 				questionCount := int(questionCountFloat)
// 				questionIds, err := getAllQuestionIds(techTypeId, questionType)
// 				if err != nil {
// 					return nil, allQuestionsId, err
// 				}
// 				questionIdMap := make(map[string]bool)
// 				questionIdSlice := make(models.QuestionIdSlice, 0)
// 				seed := rand.NewSource(time.Now().Unix())
// 				r := rand.New(seed)
// 				questionCountCopy := questionCount

// 				for questionCountCopy > 0 {
// 					pickedQuestionId := questionIds[r.Intn(questionCount)]
// 					if _, ok = questionIdMap[pickedQuestionId]; !ok {
// 						questionIdMap[pickedQuestionId] = true
// 						questionIdSlice = append(questionIdSlice, []string{pickedQuestionId, "false"})
// 						allQuestionsId = append(allQuestionsId, pickedQuestionId)
// 						questionCountCopy -= 1
// 					}
// 				}

// 				questionData[techTypeId][questionType] = questionIdSlice
// 			} else {
// 				if questionsSlices, ok := questionIdentifier.([]interface{}); ok {
// 					for _, questionSlice := range questionsSlices {
// 						if questionSlice, ok := questionSlice.([]interface{}); ok {
// 							if questionId, ok := questionSlice[0].(string); ok {
// 								allQuestionsId = append(allQuestionsId, questionId)
// 							}
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}

// 	return questionData, allQuestionsId, nil
// }

func getScoreByQuestionIds(questionIds []string) int {
	var possibleScore int
	db.DB.Model(&models.Question{}).Select("sum(marks)").Where("id IN ?", questionIds).Find(&possibleScore)
	return possibleScore
}

// func processAssessmentSessionCreateData(assessmentSessionCreate *models.AssessmentSessionCreate) error {
// 	questionData, err := getParsedQuestionData(assessmentSessionCreate.QuestionData)
// 	if err != nil {
// 		return err
// 	}

// 	newQuestionData, allQuestionsId, err := getProcessedData(questionData)
// 	if err != nil {
// 		return err
// 	}

// 	assessmentSessionCreate.PossibleScore = getScoreByQuestionIds(allQuestionsId)
// 	assessmentSessionCreate.QuestionsCount = len(allQuestionsId)

// 	res, err := json.Marshal(newQuestionData)
// 	if err != nil {
// 		return err
// 	}
// 	assessmentSessionCreate.QuestionData = string(res)
// 	return nil
// }

// func processAssessmentSessionData(assessmentSession *models.AssessmentSession) error {
// 	questionData, err := getParsedQuestionData(assessmentSession.QuestionData)
// 	if err != nil {
// 		return err
// 	}

// 	newQuestionData, allQuestionsId, err := getProcessedData(questionData)
// 	if err != nil {
// 		return err
// 	}

// 	assessmentSession.PossibleScore = getScoreByQuestionIds(allQuestionsId)
// 	assessmentSession.QuestionsCount = len(allQuestionsId)

// 	res, err := json.Marshal(newQuestionData)
// 	if err != nil {
// 		return err
// 	}
// 	assessmentSession.QuestionData = string(res)
// 	return nil
// }

// func addQuestionMeta(assessmentSessionMeta *models.AssessmentSessionMeta, assessmentSession *models.AssessmentSession) error {
// 	assessmentSessionMeta.QuestionsMeta = make(models.QuestionsMeta)
// 	questionData, err := getParsedQuestionData(assessmentSession.QuestionData)
// 	if err != nil {
// 		return err
// 	}

// 	for techTypeId, questionTypeData := range questionData {
// 		assessmentSessionMeta.QuestionsMeta[techTypeId] = make(map[string]int)
// 		for questionType, questionIdentifier := range questionTypeData {
// 			if questionsSlices, ok := questionIdentifier.([]interface{}); ok {
// 				assessmentSessionMeta.QuestionsMeta[techTypeId][questionType] = len(questionsSlices)
// 			}
// 		}
// 	}

// 	return nil
// }
