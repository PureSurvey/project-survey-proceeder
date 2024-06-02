package context

import (
	"project-survey-proceeder/internal/enums"
	"project-survey-proceeder/internal/events/contracts"
	"time"
)

type ProceederContext struct {
	MessageProducer contracts.IEventProducer

	RequestTimestamp time.Time

	UnitId     int
	SurveyId   int
	QuestionId int
	OptionIds  []int

	Ip string

	UserAgent string
	Platform  string
	IsMobile  bool

	Country    string
	Language   string
	Longtitude string
	Latitude   string

	UserCookie string

	EventType      enums.EventType
	MismatchReason enums.MismatchReason
}

func (pc *ProceederContext) IsMismatched() bool {
	return pc.MismatchReason != enums.MRUnknown
}
