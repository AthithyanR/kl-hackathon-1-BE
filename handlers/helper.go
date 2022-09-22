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

func processAssessmentSessionCreateData(assessmentSessionCreate *models.AssessmentSessionCreate) error {
	var questionData models.QuestionData
	err := json.Unmarshal([]byte(assessmentSessionCreate.QuestionData), &questionData)
	if err != nil {
		return err
	}

	allQuestionsId := make([]string, 0)

	for techType, questionTypeData := range questionData {
		for questionType, questionIdentifier := range questionTypeData {
			if questionCountFloat, ok := questionIdentifier.(float64); ok {
				questionCount := int(questionCountFloat)
				questionIds, err := getAllQuestionIds(techType, questionType)
				if err != nil {
					return err
				}
				questionIdMap := make(map[string]bool)
				seed := rand.NewSource(time.Now().Unix())
				r := rand.New(seed)
				questionCountCopy := questionCount

				for questionCountCopy > 0 {
					pickedQuestionId := questionIds[r.Intn(questionCount)]
					if _, ok = questionIdMap[pickedQuestionId]; !ok {
						questionIdMap[pickedQuestionId] = true
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

	db.DB.Model(&models.Question{}).Select("sum(marks)").Where("id IN ?", allQuestionsId).Find(&assessmentSessionCreate.PossibleScore)
	res, err := json.Marshal(questionData)

	if err != nil {
		return err
	}

	assessmentSessionCreate.QuestionData = string(res)
	return nil
}
