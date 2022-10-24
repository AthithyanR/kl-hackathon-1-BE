package models

type AssessmentSessionQuestion struct {
	AssessmentSessionId string `json:"assessmentSessionId" gorm:"primaryKey"`
	QuestionId          string `json:"questionId" gorm:"primaryKey"`
	ChosenOption        string `json:"chosenOption"`
}

type AssessmentSessionQuestionsList struct {
	AssessmentSessionId   string `json:"assessmentSessionId"`
	QuestionId            string `json:"questionId"`
	TechTypeId            string `json:"techTypeId"`
	TechType              string `json:"techType"`
	QuestionType          string `json:"questionType"`
	ChosenOption          string `json:"chosenOption"`
	IsChosenOptionCorrect bool   `json:"isChosenOptionCorrect"`
}

func (AssessmentSessionQuestionsList) TableName() string {
	return "assessment_session_questions_list"
}
