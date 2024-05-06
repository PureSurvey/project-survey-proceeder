package enums

type EnumEventType int

const (
	Unknown EnumEventType = iota
	SurveyStart
	SurveyClose
)
