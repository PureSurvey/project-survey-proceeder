package objects

import "project-survey-proceeder/enums"

type Question struct {
	Id          int
	Type        enums.EnumQuestionType
	SurveyId    int
	OrderNumber int
}
