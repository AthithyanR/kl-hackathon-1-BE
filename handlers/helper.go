package handlers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/AthithyanR/kl-hackathon-1-BE/db"
	"github.com/AthithyanR/kl-hackathon-1-BE/models"
)

func getAllQuestionIds(techTypeId string, questionType string) ([]string, error) {
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

func getParsedQuestionData(questionDataString string) (models.QuestionData, error) {
	var questionData models.QuestionData
	err := json.Unmarshal([]byte(questionDataString), &questionData)
	if err != nil {
		return questionData, err
	}
	return questionData, nil
}

func getProcessedData(questionData models.QuestionData) (models.QuestionData, []string, error) {
	allQuestionsId := make([]string, 0)

	for techTypeId, questionTypeData := range questionData {
		for questionType, questionIdentifier := range questionTypeData {
			if questionCountFloat, ok := questionIdentifier.(float64); ok {
				questionCount := int(questionCountFloat)
				questionIds, err := getAllQuestionIds(techTypeId, questionType)
				if err != nil {
					return nil, allQuestionsId, err
				}
				questionIdMap := make(map[string]bool)
				questionIdSlice := make(models.QuestionIdSlice, 0)
				seed := rand.NewSource(time.Now().Unix())
				r := rand.New(seed)
				questionCountCopy := questionCount

				for questionCountCopy > 0 {
					pickedQuestionId := questionIds[r.Intn(questionCount)]
					if _, ok = questionIdMap[pickedQuestionId]; !ok {
						questionIdMap[pickedQuestionId] = true
						questionIdSlice = append(questionIdSlice, []string{pickedQuestionId, "false"})
						allQuestionsId = append(allQuestionsId, pickedQuestionId)
						questionCountCopy -= 1
					}
				}

				questionData[techTypeId][questionType] = questionIdSlice
			} else {
				if questionsSlices, ok := questionIdentifier.([]interface{}); ok {
					for _, questionSlice := range questionsSlices {
						if questionsSliceNew, ok := questionSlice.([]interface{}); ok {
							if questionId, ok := questionsSliceNew[0].(string); ok {
								allQuestionsId = append(allQuestionsId, questionId)
							}
						}
					}
				}
			}
		}
	}

	return questionData, allQuestionsId, nil
}

func getPossibleScore(questionIds []string) int {
	var possibleScore int
	db.DB.Model(&models.Question{}).Select("sum(marks)").Where("id IN ?", questionIds).Find(&possibleScore)
	return possibleScore
}

func processAssessmentSessionCreateData(assessmentSessionCreate *models.AssessmentSessionCreate) error {
	questionData, err := getParsedQuestionData(assessmentSessionCreate.QuestionData)
	if err != nil {
		return err
	}

	newQuestionData, allQuestionsId, err := getProcessedData(questionData)
	if err != nil {
		return err
	}

	assessmentSessionCreate.PossibleScore = getPossibleScore(allQuestionsId)
	assessmentSessionCreate.QuestionsCount = len(allQuestionsId)

	res, err := json.Marshal(newQuestionData)
	if err != nil {
		return err
	}
	assessmentSessionCreate.QuestionData = string(res)
	return nil
}

func processAssessmentSessionData(assessmentSession *models.AssessmentSession) error {
	questionData, err := getParsedQuestionData(assessmentSession.QuestionData)
	if err != nil {
		return err
	}

	newQuestionData, allQuestionsId, err := getProcessedData(questionData)
	if err != nil {
		return err
	}

	assessmentSession.PossibleScore = getPossibleScore(allQuestionsId)
	assessmentSession.QuestionsCount = len(allQuestionsId)

	res, err := json.Marshal(newQuestionData)
	if err != nil {
		return err
	}
	assessmentSession.QuestionData = string(res)
	return nil
}

func addQuestionMeta(assessmentSessionMeta *models.AssessmentSessionMeta, assessmentSession *models.AssessmentSession) error {
	assessmentSessionMeta.QuestionsMeta = make(models.QuestionsMeta)
	questionData, err := getParsedQuestionData(assessmentSession.QuestionData)
	if err != nil {
		return err
	}

	for techTypeId, questionTypeData := range questionData {
		assessmentSessionMeta.QuestionsMeta[techTypeId] = make(map[string]int)
		for questionType, questionIdentifier := range questionTypeData {
			fmt.Println(questionIdentifier)
			if questionsSlices, ok := questionIdentifier.([]interface{}); ok {
				assessmentSessionMeta.QuestionsMeta[techTypeId][questionType] = len(questionsSlices)
			}
		}
	}

	return nil
}
