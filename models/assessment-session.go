package models

import (
	"time"
)

type QuestionData = map[string]map[string]any

type AssessmentSession struct {
	Id                   string     `json:"id"`
	CandidateEmail       string     `json:"candidateEmail"`
	QuestionData         string     `json:"questionData"`
	TimeAllowedInMins    int        `json:"timeAllowedInMins"`
	Score                int        `json:"score"`
	PossibleScore        int        `json:"possibleScore"`
	ScoreOutOf100Percent float32    `json:"scoreOutOf100Percent" gorm:"column:score_out_of_100_percent"`
	IsEmailSent          bool       `json:"isEmailSent"`
	StartTime            *time.Time `json:"startTime" gorm:"type:TIMESTAMP NULL"`
	EndTime              *time.Time `json:"endTime" gorm:"type:TIMESTAMP NULL"`
	QuestionsCount       int        `json:"questionsCount"`
}

type AssessmentSessionCreate struct {
	AssessmentSession
	CandidateEmails []string `json:"candidateEmails"`
}

type QuestionsMeta = map[string]map[string]int

type AssessmentSessionMeta struct {
	QuestionsMeta     QuestionsMeta `json:"questionsMeta"`
	CandidateEmail    string        `json:"candidateEmail"`
	TimeAllowedInMins int           `json:"timeAllowedInMins"`
	QuestionsCount    int           `json:"questionsCount"`
	StartTime         *time.Time    `json:"startTime"`
}

type QuestionIdSlice = [][]string
