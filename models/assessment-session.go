package models

type QuestionData = map[string]map[string]any

type AssessmentSession struct {
	Id                   string  `json:"id"`
	CandidateEmail       string  `json:"candidateEmail"`
	QuestionData         string  `json:"questionData"`
	TimeAllowedInMins    int     `json:"timeAllowedInMins"`
	Score                int     `json:"score"`
	PossibleScore        int     `json:"possibleScore"`
	ScoreOutOf100Percent float32 `json:"scoreOutOf100Percent" gorm:"column:score_out_of_100_percent"`
	IsEmailSent          int     `json:"isEmailSent"`
	StartTime            *string `json:"startTime"`
	EndTime              *string `json:"endTime"`
}

type AssessmentSessionCreate struct {
	AssessmentSession
	CandidateEmails []string `json:"candidateEmails"`
}
