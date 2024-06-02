package filler

import (
	"fmt"
	"github.com/ferluci/fast-realip"
	"github.com/valyala/fasthttp"
	"project-survey-proceeder/internal/context"
	geolocationcontracts "project-survey-proceeder/internal/geolocation/contracts"
	"project-survey-proceeder/internal/pools"
	"time"
)

type Base struct {
	UserAgentPool      *pools.UserAgentPool
	GeolocationService geolocationcontracts.IGeolocationService
}

func NewBaseFiller(pool *pools.UserAgentPool, geoservice geolocationcontracts.IGeolocationService) *Base {
	return &Base{UserAgentPool: pool, GeolocationService: geoservice}
}

func (f *Base) FillFromRequest(prCtx *context.ProceederContext, ctx *fasthttp.RequestCtx) error {
	if ctx.QueryArgs() == nil {
		return fmt.Errorf("nil query args")
	}

	prCtx.RequestTimestamp = time.Now().UTC()

	f.fillUserAgent(prCtx, ctx)
	f.fillIp(prCtx, ctx)
	f.fillCountry(prCtx)

	return nil
}

func (f *Base) fillUserAgent(prCtx *context.ProceederContext, ctx *fasthttp.RequestCtx) {
	prCtx.UserAgent = string(ctx.Request.Header.UserAgent())

	if prCtx.UserAgent != `` {
		ua := f.UserAgentPool.Get()
		defer f.UserAgentPool.Put(ua)

		ua.Parse(prCtx.UserAgent)
		prCtx.IsMobile = ua.Mobile()
		prCtx.Platform = ua.Platform()
	}
}

func (f *Base) fillIp(prCtx *context.ProceederContext, ctx *fasthttp.RequestCtx) {
	prCtx.Ip = realip.FromRequest(ctx)
}

func (f *Base) fillCountry(prCtx *context.ProceederContext) {
	country, err := f.GeolocationService.GetCountryByIp(prCtx.Ip)
	if err != nil {
		return
	}

	prCtx.Country = country
}
