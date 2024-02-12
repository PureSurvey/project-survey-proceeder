package request

import (
	"github.com/valyala/fasthttp"
	"project-survey-proceeder/context"
	"project-survey-proceeder/events"
)

type RequestHandler struct {
	ProceederContext *context.ProceederContext
}

func (h *RequestHandler) HandleRequest(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Path()) {
	case "/sev":
		handleSurveyEvent(h.ProceederContext, ctx)
	default:
		ctx.Error("Unsupported path", fasthttp.StatusNotFound)
	}
}

func handleSurveyEvent(prCtx *context.ProceederContext, ctx *fasthttp.RequestCtx) {
	err := context.FillContextFromRequest(prCtx, ctx)
	if err != nil {
		ctx.Error("Invalid request", fasthttp.StatusBadRequest)
		return
	}

	eventString := events.GetEventString(prCtx)

	err = prCtx.MessageProducer.SendMessage([]byte(eventString))
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
}
