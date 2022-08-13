package main

type TechType struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (TechType) TableName() string {
	return "tech_type"
}

type Question struct {
	Id            string `json:"id"`
	TechType      string `json:"techType"`
	QuestionType  string `json:"questionType"`
	Question      string `json:"question"`
	Option1       string `json:"option1"`
	Option2       string `json:"option2"`
	Option3       string `json:"option3"`
	Option4       string `json:"option4"`
	CorrectOption string `json:"correctOption"`
}
