package models

type QuestionList struct {
	Id            string `json:"id"`
	TechTypeId    string `json:"techTypeId"`
	TechTypeName  string `json:"techTypeName"`
	QuestionType  string `json:"questionType"`
	Question      string `json:"question"`
	Option1       string `json:"option1"`
	Option2       string `json:"option2"`
	Option3       string `json:"option3"`
	Option4       string `json:"option4"`
	CorrectOption string `json:"correctOption"`
	Marks         int    `json:"marks"`
}

func (QuestionList) TableName() string {
	return "questions_list"
}
