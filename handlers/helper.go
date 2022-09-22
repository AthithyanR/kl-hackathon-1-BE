package handlers

import (
	"encoding/json"
	"math/rand"
	"time"

	"github.com/AthithyanR/kl-hackathon-1-BE/db"
	"github.com/AthithyanR/kl-hackathon-1-BE/models"
)

func getAllQuestionIds(techType string, questionType string) ([]string, error) {
	var questionIds []string
	whereClause := &models.Question{
		TechTypeId:   techType,
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

	for techType, questionTypeData := range questionData {
		for questionType, questionIdentifier := range questionTypeData {
			if questionCountFloat, ok := questionIdentifier.(float64); ok {
				questionCount := int(questionCountFloat)
				questionIds, err := getAllQuestionIds(techType, questionType)
				if err != nil {
					return nil, allQuestionsId, err
				}
				questionIdMap := make(map[string]bool)
				seed := rand.NewSource(time.Now().Unix())
				r := rand.New(seed)
				questionCountCopy := questionCount

				for questionCountCopy > 0 {
					pickedQuestionId := questionIds[r.Intn(questionCount)]
					if _, ok = questionIdMap[pickedQuestionId]; !ok {
						questionIdMap[pickedQuestionId] = false
						allQuestionsId = append(allQuestionsId, pickedQuestionId)
						questionCountCopy -= 1
					}
				}

				questionData[techType][questionType] = questionIdMap
			} else if questionsMap, ok := questionIdentifier.(map[string]interface{}); ok {
				for questionId := range questionsMap {
					allQuestionsId = append(allQuestionsId, questionId)
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
