package filler

import (
	"fmt"
	"github.com/ferluci/fast-realip"
	"github.com/valyala/fasthttp"
	"project-survey-proceeder/internal/context"
	"project-survey-proceeder/internal/enums"
	geolocationcontracts "project-survey-proceeder/internal/geolocation/contracts"
	"project-survey-proceeder/internal/pools"
	"time"
)

type BaseFiller struct {
	UserAgentPool      *pools.UserAgentPool
	GeolocationService geolocationcontracts.IGeolocationService
}

func NewBaseFiller(pool *pools.UserAgentPool, geoservice geolocationcontracts.IGeolocationService) *BaseFiller {
	return &BaseFiller{UserAgentPool: pool, GeolocationService: geoservice}
}

func (f *BaseFiller) FillFromRequest(prCtx *context.ProceederContext, ctx *fasthttp.RequestCtx) error {
	if ctx.QueryArgs() == nil {
		return fmt.Errorf("nil query args")
	}

	prCtx.RequestTimestamp = time.Now().UTC()

	f.fillUserAgent(prCtx, ctx)
	f.fillIp(prCtx, ctx)
	f.fillCountry(prCtx)

	f.fillEntityIds(prCtx, ctx)
	f.fillEventType(prCtx, ctx)

	return nil
}

func (f *BaseFiller) fillUserAgent(prCtx *context.ProceederContext, ctx *fasthttp.RequestCtx) {
	prCtx.UserAgent = string(ctx.Request.Header.UserAgent())

	if prCtx.UserAgent != `` {
		ua := f.UserAgentPool.Get()
		defer f.UserAgentPool.Put(ua)

		ua.Parse(prCtx.UserAgent)
		prCtx.IsMobile = ua.Mobile()
		prCtx.Platform = ua.Platform()
	}
}

func (f *BaseFiller) fillIp(prCtx *context.ProceederContext, ctx *fasthttp.RequestCtx) {
	prCtx.Ip = realip.FromRequest(ctx)
}

func (f *BaseFiller) fillCountry(prCtx *context.ProceederContext) {
	country, err := f.GeolocationService.GetCountryByIp(prCtx.Ip)
	if err != nil {
		return
	}

	prCtx.Country = country
}

func (f *BaseFiller) fillEntityIds(prCtx *context.ProceederContext, ctx *fasthttp.RequestCtx) {
	unitId, err := ctx.QueryArgs().GetUint(`unitId`)
	if err != nil {
		return
	}

	prCtx.UnitId = unitId
}

func (f *BaseFiller) fillEventType(prCtx *context.ProceederContext, ctx *fasthttp.RequestCtx) {
	et, err := ctx.QueryArgs().GetUint(`et`)
	if err != nil {
		return
	}

	prCtx.EventType = enums.EnumEventType(et)
}
