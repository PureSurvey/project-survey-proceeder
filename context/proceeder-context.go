package context

import (
	"github.com/valyala/fastjson"
	"project-survey-proceeder/contracts"
	"project-survey-proceeder/enums"
	"project-survey-proceeder/pools"
)

type ProceederContext struct {
	MessageProducer contracts.IMessageProducer
	ParserPool      *fastjson.ParserPool
	UserAgentPool   *pools.UserAgentPool

	UnitId int

	UserAgent string
	Platform  string
	IsMobile  bool

	EventType enums.EnumEventType
}
