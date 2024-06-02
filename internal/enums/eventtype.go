package enums

type EventType int

const (
	ETUnknown EventType = iota
	ETUnitView
	ETQuestionView
	ETQuestionAnswer
	ETSurveyStart
	ETSurveyEnd
	ETUnitEnd
	ETUnitRequest
	ETUnitResponse
)
