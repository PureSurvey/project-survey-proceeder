package events

import (
	"fmt"
	"project-survey-proceeder/internal/context"
	"project-survey-proceeder/internal/enums"
	"strconv"
	"strings"
	"time"
)

type event struct {
	Date string
	Hour int

	UnitId int

	UserAgent string
	Platform  string
	IsMobile  bool

	EventType enums.EnumEventType
}

func GetEventString(ptCtx *context.ProceederContext) string {
	sb := strings.Builder{}

	curTime := time.Now().UTC()

	sb.WriteString(`{"date":"`)
	sb.WriteString(fmt.Sprintf(curTime.Format("2006-01-02")))
	sb.WriteString(`","hour":`)
	sb.WriteString(strconv.Itoa(curTime.Hour()))

	sb.WriteString(`,"unitId":`)
	sb.WriteString(strconv.Itoa(ptCtx.UnitId))

	sb.WriteString(`,"userAgent":"`)
	sb.WriteString(ptCtx.UserAgent)
	sb.WriteString(`","platform":"`)
	sb.WriteString(ptCtx.Platform)
	sb.WriteString(`"`)

	if ptCtx.IsMobile {
		sb.WriteString(`,"isMobile":1`)
	}

	sb.WriteString(`,"eventType":`)
	sb.WriteString(strconv.Itoa(int(ptCtx.EventType)))

	sb.WriteString(`}`)

	return sb.String()
}
