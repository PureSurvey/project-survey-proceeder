package request

import (
	"github.com/valyala/fasthttp"
	"project-survey-proceeder/internal/context"
	"project-survey-proceeder/internal/events"
	"project-survey-proceeder/internal/services/contracts"
)

type Handler struct {
	ProceederContext *context.ProceederContext
	ServiceProvider  contracts.IServiceProvider
}

func (h *Handler) Handle(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Path()) {
	case "/unit":
		h.handleUnitRequest(ctx)
	case "/sev":
		h.handleSurveyEvent(ctx)
	default:
		ctx.Error("Unsupported path", fasthttp.StatusNotFound)
	}
}

func (h *Handler) handleUnitRequest(ctx *fasthttp.RequestCtx) {
	filler := h.ServiceProvider.GetContextFiller()

	err := filler.FillFromRequest(h.ProceederContext, ctx)
	if err != nil {
		ctx.Error("", fasthttp.StatusNoContent)
		return
	}

	// eventString := events.GetEventString(prCtx)

	// err = prCtx.MessageProducer.SendMessage([]byte(eventString))
	// if err != nil {
	// ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	// }

	ctx.SetStatusCode(fasthttp.StatusOK)
}

func (h *Handler) handleSurveyEvent(ctx *fasthttp.RequestCtx) {
	filler := h.ServiceProvider.GetContextFiller()

	err := filler.FillFromRequest(h.ProceederContext, ctx)
	if err != nil {
		ctx.Error("Invalid request", fasthttp.StatusBadRequest)
		return
	}

	eventString := events.GetEventString(h.ProceederContext)

	err = h.ProceederContext.MessageProducer.SendMessage([]byte(eventString))
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	}

	ctx.SetStatusCode(fasthttp.StatusNoContent)
}
