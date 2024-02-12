package context

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fastjson"
	"project-survey-proceeder/enums"
)

func FillContextFromRequest(prCtx *ProceederContext, ctx *fasthttp.RequestCtx) error {
	if ctx.PostBody() == nil {
		return fmt.Errorf("empty body")
	}

	if err := fastjson.ValidateBytes(ctx.PostBody()); err != nil {
		return err
	}

	if prCtx.ParserPool == nil {
		prCtx.ParserPool = &fastjson.ParserPool{}
	}

	parser := prCtx.ParserPool.Get()
	defer prCtx.ParserPool.Put(parser)

	val, err := parser.Parse(string(ctx.PostBody()))
	if err != nil {
		return err
	}

	prCtx.UnitId = val.GetInt(`unitId`)
	prCtx.EventType = enums.EnumEventType(val.GetInt(`et`))

	prCtx.UserAgent = string(ctx.Request.Header.UserAgent())

	if prCtx.UserAgent != `` {
		ua := prCtx.UserAgentPool.Get()
		defer prCtx.UserAgentPool.Put(ua)

		ua.Parse(prCtx.UserAgent)
		prCtx.IsMobile = ua.Mobile()
		prCtx.Platform = ua.Platform()
	}

	return nil
}
