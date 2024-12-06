package request

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"google.golang.org/protobuf/proto"
	"project-survey-proceeder/internal/context"
	contextcontracts "project-survey-proceeder/internal/context/contracts"
	"project-survey-proceeder/internal/dbcache"
	"project-survey-proceeder/internal/enums"
	"project-survey-proceeder/internal/events/contracts"
	"project-survey-proceeder/internal/events/model/pb"
	"project-survey-proceeder/internal/respondents"
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
	respondentsService  *respondents.Service
}

func NewHandler(dbRepo *dbcache.Repo,
	unitContextFiller contextcontracts.IRequestFiller, eventContextFiller contextcontracts.IRequestFiller,
	targetingService targetingcontracts.ITargetingService,
	surveyMarkupService surveymarkupcontracts.ISurveyMarkupService,
	eventProducer contracts.IEventProducer, usersService *respondents.Service) *Handler {
	return &Handler{
		dbRepo:              dbRepo,
		unitContextFiller:   unitContextFiller,
		eventContextFiller:  eventContextFiller,
		targetingService:    targetingService,
		surveyMarkupService: surveyMarkupService,
		eventProducer:       eventProducer,
		respondentsService:  usersService,
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
	h.respondentsService.UpdateContext(prCtx)

	unit := h.dbRepo.GetUnitById(prCtx.UnitId)
	if unit == nil {
		ctx.Error("", fasthttp.StatusNoContent)
		return
	}

	user := h.dbRepo.GetUserById(unit.UserId)
	if user == nil || !user.IsSubscribed {
		return
	}

	surveys := h.dbRepo.GetSurveysByUnitId(unit.Id)
	if surveys == nil {
		ctx.Error("", fasthttp.StatusNoContent)
		return
	}

	completedUnitSurveys := h.respondentsService.GetCompletedSurveysByUnit(prCtx)

	var matchedSurveyIds []int
	for _, survey := range surveys {
		if len(matchedSurveyIds)+completedUnitSurveys == unit.MaximumSurveysPerDevice {
			break
		}

		takes := h.respondentsService.GetSurveyTakesByUnit(prCtx, survey.Id)
		if takes < unit.SurveyTakesPerDevice && survey.IsActiveOnDate(prCtx.RequestTimestamp) && h.targetingService.IsMatched(survey, prCtx) {
			matchedSurveyIds = append(matchedSurveyIds, survey.Id)
		}
	}

	if len(matchedSurveyIds) == 0 {
		ctx.Error("", fasthttp.StatusNoContent)
		return
	}

	markup, err := h.surveyMarkupService.GetMarkup(unit.Id, matchedSurveyIds, prCtx.Language)
	if err != nil {
		ctx.Error("", fasthttp.StatusNoContent)
		return
	}

	var gzippedBody []byte
	gzippedBody = fasthttp.AppendGzipBytes(gzippedBody, []byte(markup))
	ctx.SetBody(gzippedBody)
	ctx.Response.Header.Set("Content-Encoding", "gzip")
	ctx.Response.Header.Set("Set-Cookie", fmt.Sprintf("psid=%v; Max-Age=84600; samesite=none;secure=1;", prCtx.UserCookie))
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func (h *Handler) handleSurveyEvent(ctx *fasthttp.RequestCtx) {
	prCtx := &context.ProceederContext{}

	err := h.eventContextFiller.FillFromRequest(prCtx, ctx)
	if err != nil {
		ctx.Error("Invalid request", fasthttp.StatusBadRequest)
		return
	}
	h.respondentsService.UpdateContext(prCtx)

	if prCtx.EventType == enums.ETSurveyStart {
		go h.respondentsService.IncrementSurveyTakeByUnit(prCtx)
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
