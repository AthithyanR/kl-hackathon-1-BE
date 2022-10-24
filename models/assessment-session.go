package models

import (
	"time"

	"github.com/AthithyanR/kl-hackathon-1-BE/customTypes"
)

type AssessmentSession struct {
	Id                   string              `json:"id"`
	CandidateEmail       string              `json:"candidateEmail"`
	TimeAllowedInMins    int                 `json:"timeAllowedInMins"`
	Score                int                 `json:"score"`
	PossibleScore        int                 `json:"possibleScore"`
	ScoreOutOf100Percent float32             `json:"scoreOutOf100Percent" gorm:"column:score_out_of_100_percent"`
	IsEmailSent          customTypes.BitBool `json:"isEmailSent"`
	StartTime            *time.Time          `json:"startTime" gorm:"type:TIMESTAMP NULL"`
	EndTime              *time.Time          `json:"endTime" gorm:"type:TIMESTAMP NULL"`
	QuestionsCount       int                 `json:"questionsCount"`
}

type RandomQuestions struct {
	TechTypeId   string `json:"techTypeId"`
	QuestionType string `json:"questionType"`
	Count        int    `json:"count"`
}

type AssessmentSessionUpsert struct {
	AssessmentSession
	QuestionIds     []string          `json:"questionIds"`
	RandomQuestions []RandomQuestions `json:"randomQuestions"`
	CandidateEmails []string          `json:"candidateEmails"`
}

type QuestionsMeta = map[string]map[string]int

type AssessmentSessionMeta struct {
	QuestionsMeta     QuestionsMeta `json:"questionsMeta"`
	CandidateEmail    string        `json:"candidateEmail"`
	TimeAllowedInMins int           `json:"timeAllowedInMins"`
	QuestionsCount    int           `json:"questionsCount"`
	StartTime         *time.Time    `json:"startTime"`
}

type AssessmentSessionDetails struct {
	TechTypeId   string `json:"techTypeId"`
	QuestionType string `json:"questionType"`
	QuestionId   string `json:"questionId"`
}

type AssessmentSessionDetailsMeta = map[string]map[string][]string
