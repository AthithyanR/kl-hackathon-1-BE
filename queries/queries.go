package queries

// assessment session queries
var (
	UpdateScoreByAnswer = `
		update assessment_sessions set score = (
			select ifnull(sum(marks), 0) from questions q left join
			assessment_session_questions ases on ases.question_id = q.id where
			ases.assessment_session_id = ?
			and ases.chosen_option = q.correct_option
			) where id = ?;
	`
	UpdateScorePercentageById = `
		update assessment_sessions set
			score_out_of_100_percent = (score / possible_score) * 100 
			where id = ?;
	`
	GetAssessmentSessionMetaById = `
		select
			tech_type_id,
			question_type,
			question_id
		from
			assessment_session_questions_list asql
		where
			asql.assessment_session_id = ?
		order by asql.tech_type;
	`
)

var ()
