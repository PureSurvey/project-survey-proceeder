package filler

import (
	"github.com/valyala/fasthttp"
	"project-survey-proceeder/internal/context"
	"project-survey-proceeder/internal/context/contracts"
	"project-survey-proceeder/internal/enums"
	geolocationcontracts "project-survey-proceeder/internal/geolocation/contracts"
	"project-survey-proceeder/internal/pools"
	"project-survey-proceeder/internal/trackers"
	"slices"
	"strconv"
	"time"
)

type Event struct {
	*Base
	decryptor *trackers.Decryptor
}

func NewEventFiller(pool *pools.UserAgentPool, geoservice geolocationcontracts.IGeolocationService, decryptor *trackers.Decryptor) contracts.IRequestFiller {
	return &Event{Base: NewBaseFiller(pool, geoservice), decryptor: decryptor}
}

func (f *Event) FillFromRequest(prCtx *context.ProceederContext, ctx *fasthttp.RequestCtx) error {
	err := f.Base.FillFromRequest(prCtx, ctx)
	if err != nil {
		return err
	}

	f.fillEntityIds(prCtx, ctx)

	return nil
}

func (f *Event) fillEntityIds(prCtx *context.ProceederContext, ctx *fasthttp.RequestCtx) {
	eventTypeInt, err := ctx.QueryArgs().GetUint(trackers.PEventType)
	if err != nil {
		return
	}
	eventType := enums.EventType(eventTypeInt)

	unitId, err := ctx.QueryArgs().GetUint(trackers.PUnitId)
	if err != nil {
		return
	}

	surveyId, _ := ctx.QueryArgs().GetUint(trackers.PSurveyId)
	if err != nil {
		return
	}

	questionId, _ := ctx.QueryArgs().GetUint(trackers.PQuestionId)
	if err != nil {
		return
	}

	optionsStringIds := ctx.QueryArgs().PeekMulti(trackers.PAnsweredOptions)
	optionsIds, _ := convert(optionsStringIds)
	if err != nil {
		return
	}

	validToString := string(ctx.QueryArgs().Peek(trackers.PValidTo))
	validTo, err := strconv.ParseInt(validToString, 10, 64)
	if err != nil {
		return
	}

	encryptedEvent := string(ctx.QueryArgs().Peek(trackers.PEncryptedEvent))
	event, err := f.decryptor.DecryptEvent(encryptedEvent)
	if err != nil {
		return
	}

	eventTypeValid := eventType != enums.ETUnknown && eventType == event.EventType
	unitIdValid := unitId == event.UnitId
	surveyIdValid := event.ValidSurveys == nil || slices.Contains(event.ValidSurveys, surveyId)
	questionIdValid := event.ValidQuestions == nil || slices.Contains(event.ValidQuestions, questionId)
	optionsIdsValid := event.ValidQuestionsWithAnswers == nil ||
		(event.ValidQuestionsWithAnswers[questionId] != nil && all(optionsIds, func(id int) bool { return slices.Contains(event.ValidQuestionsWithAnswers[questionId], id) }))

	if !eventTypeValid || !unitIdValid || !surveyIdValid || !questionIdValid || !optionsIdsValid {
		prCtx.MismatchReason = enums.MRInvalidTracker
	}

	if time.Unix(validTo, 0).Before(prCtx.RequestTimestamp) {
		prCtx.MismatchReason = enums.MROutdatedTracker
	}

	prCtx.EventType = enums.EventType(eventTypeInt)
	prCtx.UnitId = unitId
	prCtx.SurveyId = surveyId
	prCtx.QuestionId = questionId
	prCtx.OptionIds = optionsIds
}

func convert(arr [][]byte) ([]int, error) {
	var result []int
	for _, el := range arr {
		val, err := strconv.Atoi(string(el))
		if err != nil {
			return nil, err
		}
		result = append(result, val)
	}

	return result, nil
}

func all[T any](ts []T, pred func(T) bool) bool {
	for _, t := range ts {
		if !pred(t) {
			return false
		}
	}
	return true
}
