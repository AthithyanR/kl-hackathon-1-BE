package models

type Question struct {
	Id            string `json:"id"`
	TechTypeId    string `json:"techTypeId"`
	QuestionType  string `json:"questionType"`
	Question      string `json:"question"`
	Option1       string `json:"option1"`
	Option2       string `json:"option2"`
	Option3       string `json:"option3"`
	Option4       string `json:"option4"`
	CorrectOption string `json:"correctOption"`
	Marks         int    `json:"marks"`
}

type QuestionList struct {
	Id            string `json:"id"`
	TechTypeId    string `json:"techTypeId"`
	TechType      string `json:"techType"`
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

type QuestionLimited struct {
	Id       string `json:"id"`
	Question string `json:"question"`
	Option1  string `json:"option1"`
	Option2  string `json:"option2"`
	Option3  string `json:"option3"`
	Option4  string `json:"option4"`
}
