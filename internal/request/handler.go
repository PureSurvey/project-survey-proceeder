package request

import (
	"github.com/valyala/fasthttp"
	"google.golang.org/protobuf/proto"
	"project-survey-proceeder/internal/context"
	contextcontracts "project-survey-proceeder/internal/context/contracts"
	"project-survey-proceeder/internal/dbcache"
	"project-survey-proceeder/internal/events/contracts"
	"project-survey-proceeder/internal/events/model/pb"
	surveymarkupcontracts "project-survey-proceeder/internal/surveymarkup/contracts"
	targetingcontracts "project-survey-proceeder/internal/targeting/contracts"
	"project-survey-proceeder/internal/utils"
)

type Handler struct {
	dbRepo              *dbcache.Repo
	unitContextFiller   contextcontracts.IRequestFiller
	eventContextFiller  contextcontracts.IRequestFiller
	targetingService    targetingcontracts.ITargetingService
	surveyMarkupService surveymarkupcontracts.ISurveyMarkupService
	eventProducer       contracts.IEventProducer
}

func NewHandler(dbRepo *dbcache.Repo,
	unitContextFiller contextcontracts.IRequestFiller, eventContextFiller contextcontracts.IRequestFiller,
	targetingService targetingcontracts.ITargetingService,
	surveyMarkupService surveymarkupcontracts.ISurveyMarkupService,
	eventProducer contracts.IEventProducer) *Handler {
	return &Handler{
		dbRepo:              dbRepo,
		unitContextFiller:   unitContextFiller,
		eventContextFiller:  eventContextFiller,
		targetingService:    targetingService,
		surveyMarkupService: surveyMarkupService,
		eventProducer:       eventProducer,
	}
}

func (h *Handler) Handle(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Path()) {
	case "/unit":
		h.handleUnitRequest(ctx)
	case "/event":
		h.handleSurveyEvent(ctx)
	default:
		ctx.Error("Unsupported path", fasthttp.StatusNotFound)
	}
}

func (h *Handler) handleUnitRequest(ctx *fasthttp.RequestCtx) {
	prCtx := &context.ProceederContext{}

	err := h.unitContextFiller.FillFromRequest(prCtx, ctx)
	if err != nil {
		ctx.Error("", fasthttp.StatusNoContent)
		return
	}

	if prCtx.UnitId == 0 {
		ctx.Error("", fasthttp.StatusNoContent)
		return
	}

	unit := h.dbRepo.GetUnitById(prCtx.UnitId)
	if unit == nil {
		ctx.Error("", fasthttp.StatusNoContent)
		return
	}

	surveys := h.dbRepo.GetSurveysByUnitId(unit.Id)
	if surveys == nil {
		ctx.Error("", fasthttp.StatusNoContent)
		return
	}

	var matchedSurveyIds []int
	for _, survey := range surveys {
		if survey.IsActiveOnDate(prCtx.RequestTimestamp) && h.targetingService.IsMatched(survey, prCtx) {
			matchedSurveyIds = append(matchedSurveyIds, survey.Id)
		}
	}

	if len(matchedSurveyIds) == 0 {
		ctx.Error("", fasthttp.StatusNoContent)
		return
	}

	markup, err := h.surveyMarkupService.GetMarkup(unit.Id, matchedSurveyIds, "en")
	if err != nil {
		ctx.Error("", fasthttp.StatusNoContent)
		return
	}

	ctx.SetBody([]byte(markup))
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func (h *Handler) handleSurveyEvent(ctx *fasthttp.RequestCtx) {
	prCtx := &context.ProceederContext{}

	err := h.eventContextFiller.FillFromRequest(prCtx, ctx)
	if err != nil {
		ctx.Error("Invalid request", fasthttp.StatusBadRequest)
		return
	}

	completionEvent := &pb.CompletionEvent{
		EventType:  int32(prCtx.EventType),
		Timestamp:  prCtx.RequestTimestamp.Unix(),
		SurveyId:   int32(prCtx.SurveyId),
		QuestionId: int32(prCtx.QuestionId),
		OptionIds:  utils.ConvertInts[int32](prCtx.OptionIds),
		Geo:        prCtx.Country,
		Lang:       prCtx.Language,
		Gender:     0,
	}

	bytes, _ := proto.Marshal(completionEvent)
	h.eventProducer.AsyncSendMessage(bytes)

	ctx.SetStatusCode(fasthttp.StatusNoContent)
}
