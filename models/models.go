package models

import (
	"github.com/dgrijalva/jwt-go"
)

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

type User struct {
	Id       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ClaimValues struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

type Claims struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	jwt.StandardClaims
}

type InterviewSession struct {
	Id                           string  `json:"id"`
	CandidateEmail               string  `json:"candidateEmail"`
	SessionKey                   string  `json:"sessionKey"`
	TimeAllowedInMins            int     `json:"timeAllowedInMins"`
	QuestionIds                  string  `json:"questionIds"`
	AnsweredQuestionIds          string  `json:"answeredQuestionIds"`
	CorrectlyAnsweredQuestionIds string  `json:"correctlyAnsweredQuestionIds"`
	WronglyAnsweredQuestionIds   string  `json:"wronglyAnsweredQuestionIds"`
	ScoreOutOf100Percent         float32 `json:"scoreOutOf100Percent" gorm:"column:score_out_of_100_percent"`
	IsInterviewDone              int     `json:"isInterviewDone"`
}
