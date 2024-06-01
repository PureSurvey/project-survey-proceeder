package model

import "project-survey-proceeder/internal/enums"

type Event struct {
	EventType                 enums.EventType
	UnitId                    int
	ValidTo                   int64
	ValidSurveys              []int
	ValidQuestions            []int
	ValidQuestionsWithAnswers map[int][]int
}
