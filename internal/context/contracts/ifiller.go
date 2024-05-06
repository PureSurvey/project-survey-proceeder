package contracts

import (
	"github.com/valyala/fasthttp"
	"project-survey-proceeder/internal/context"
)

type IRequestFiller interface {
	FillFromRequest(prCtx *context.ProceederContext, ctx *fasthttp.RequestCtx) error
}
