package contracts

import "github.com/valyala/fasthttp"

type IRequestHandler interface {
	HandleRequests(ctx *fasthttp.RequestCtx)
}
