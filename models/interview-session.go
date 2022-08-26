package models

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
