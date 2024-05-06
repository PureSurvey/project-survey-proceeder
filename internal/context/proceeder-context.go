package context

import (
	"project-survey-proceeder/internal/enums"
	"project-survey-proceeder/internal/events/contracts"
	"time"
)

type ProceederContext struct {
	MessageProducer contracts.IEventProducer

	RequestTimestamp time.Time

	UnitId int

	Ip string

	UserAgent string
	Platform  string
	IsMobile  bool

	Country    string
	Longtitude string
	Latitude   string

	UserCookie string

	EventType enums.EnumEventType
}
