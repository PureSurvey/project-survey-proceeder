package objects

import "project-survey-proceeder/enums"

type Question struct {
	Id             int
	Type           enums.EnumQuestionType
	SurveyId       int
	OrderNumber    int
	QuestionLineId int
}

func NewQuestion(id int, qType enums.EnumQuestionType, surveyId int,
	orderNumber int, questionLineId int) *Question {
	return &Question{
		Id:             id,
		Type:           qType,
		SurveyId:       surveyId,
		OrderNumber:    orderNumber,
		QuestionLineId: questionLineId,
	}
}
